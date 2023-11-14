package gosite

import (
	"context"
	"net/url"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/osspkg/devtool/internal/global"
	"github.com/osspkg/devtool/pkg/exec"
	"github.com/osspkg/devtool/pkg/files"
	"github.com/osspkg/devtool/pkg/modules"
	"go.osspkg.com/goppy/sdk/console"
	"go.osspkg.com/goppy/sdk/iofile"
)

var (
	rexHEAD = regexp.MustCompile(`(?mU)ref\: refs/heads/(\w+)\s+HEAD`)
	rexMOD  = regexp.MustCompile(`(?mU)module (.*)\n`)
)

type Data struct {
	Branch  string
	Repo    string
	Root    string
	Modules []string
}

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("gosite", "")
		setter.ExecFunc(func(_ []string) {
			global.SetupEnv()
			console.Infof("--- READ CONFIG ---")

			confpath := files.CurrentDir() + "/.gosite.yaml"
			if !files.Exist(confpath) {
				console.Fatalf("File .gosite.yaml not found")
			}

			var configs []string
			result := make(map[string]*Data, 100)

			err := iofile.FileCodec(confpath).Decode(&configs)
			console.FatalIfErr(err, "Decode config")

			tempdir := files.CurrentDir() + "/.tmp"
			defer os.RemoveAll(tempdir) //nolint: errcheck
			for _, config := range configs {
				os.RemoveAll(tempdir) //nolint: errcheck
				console.FatalIfErr(os.MkdirAll(tempdir, 0755), "Create temp dir")

				var b []byte
				b, err = exec.SingleCmd(context.TODO(), "bash", "git ls-remote --symref "+config+" HEAD")
				console.FatalIfErr(err, "Get remote HEAD")
				_strs := rexHEAD.FindStringSubmatch(string(b))
				if len(_strs) != 2 {
					console.Fatalf("HEAD not found")
				}
				HEAD := _strs[1]

				_, err = exec.SingleCmd(context.TODO(), "bash", "git clone --branch "+HEAD+" --single-branch "+config+" .tmp")
				console.FatalIfErr(err, "Clone remote HEAD")
				os.RemoveAll(tempdir + "/.git") //nolint: errcheck

				var mods map[string]*modules.Mod
				mods, err = modules.Detect(tempdir)
				console.FatalIfErr(err, "Detect go.mod files")

				var dataMod *Data
				if dm, ok := result[config]; ok {
					dataMod = dm
				} else {
					dataMod = &Data{
						Branch:  HEAD,
						Repo:    strings.TrimSuffix(config, ".git"),
						Modules: make([]string, 0, 10),
					}
					result[config] = dataMod
				}

				for _, mod := range mods {
					b, err = os.ReadFile(mod.File)
					console.FatalIfErr(err, "Read go.mod [%s]", mod.File)
					_strs = rexMOD.FindStringSubmatch(string(b))
					if len(_strs) != 2 {
						console.Fatalf("Module not found in %s", mod.File)
					}
					module := _strs[1]
					dataMod.Modules = append(dataMod.Modules, module)
				}
				for i, module := range dataMod.Modules {
					if i == 0 {
						dataMod.Root = module
						continue
					}
					if len(dataMod.Root) > len(module) {
						dataMod.Root = module
					}
				}

				if len(dataMod.Modules) == 0 {
					delete(result, config)
				}
			}

			index := make(map[string][]string)
			for _, data := range result {
				var u *url.URL
				u, err = url.Parse("http://" + data.Root)
				console.FatalIfErr(err, "Decode module url [%s]", data.Root)
				domain := u.Host
				if _, ok := index[domain]; !ok {
					index[domain] = make([]string, 0, 10)
				}

				sort.Strings(data.Modules)
				for _, mod := range data.Modules {
					err = os.MkdirAll(mod, 0755)
					console.FatalIfErr(err, "Create site dir [%s]", mod)
					index[domain] = append(index[domain], mod)

					tmpl := strings.ReplaceAll(htmlPageTemplate, "{%module%}", mod)
					tmpl = strings.ReplaceAll(tmpl, "{%root%}", data.Root)
					tmpl = strings.ReplaceAll(tmpl, "{%repo%}", data.Repo)
					tmpl = strings.ReplaceAll(tmpl, "{%head%}", data.Branch)

					err = os.WriteFile(mod+"/index.html", []byte(tmpl), 0755)
					console.FatalIfErr(err, "Write HTML [%s]", mod+"/index.html")
				}
			}
			for domain, links := range index {
				linksHtml := ""
				for _, link := range links {
					linkName := strings.TrimPrefix(link, domain)
					linkName = strings.Trim(linkName, "/")
					linksHtml += "\n<li><a href=\"//" + link + "\">" + linkName + "</a></li>"
				}

				tmpl := strings.ReplaceAll(htmlIndexPage, "{%domain%}", domain)
				tmpl = strings.ReplaceAll(tmpl, "{%links%}", linksHtml)

				err = os.WriteFile(domain+"/index.html", []byte(tmpl), 0755)
				console.FatalIfErr(err, "Write HTML [%s]", domain+"/index.html")
			}

		})
	})
}

const (
	htmlPageTemplate = `
<!DOCTYPE html>
<html lang="en" dir="ltr">

<head>
    <title>{%module%}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, height=device-height, minimum-scale=1.0, initial-scale=1.0">
    <meta name="go-import" content="{%root%} git {%repo%}">
    <meta name="go-source" content="{%root%} {%repo%} {%repo%}/tree/{%head%}{/dir} {%repo%}/tree/{%head%}{/dir}/{file}#L{line}">
</head>

<body>
    <aside>
        <a href="/">Back Home</a>
    </aside>
    <hr>
    <div>
        <h1>{%module%}</h1>
    </div>

    <div>
        <b>Install command:</b>
        <pre>go get {%module%}</pre>
    </div>

    <div>
        <b>Import in source code:</b>
        <pre>import "{%module%}"</pre>
    </div>
        
    <div>
        <b>Repository:</b>
        <a href="{%repo%}">{%repo%}</a>
    </div>
</body>

</html>
`
	htmlIndexPage = `
<!DOCTYPE html>
<html lang="en" dir="ltr">

<head>
    <title>{%domain%}</title>
    <meta charset="utf-8">
    <meta name="viewport" content="width=device-width, height=device-height, minimum-scale=1.0, initial-scale=1.0">
</head>

<body>
    <div>
        <h1>{%domain%}</h1>
    </div>
	<hr>
    <aside>
        <ul>
            {%links%}
        </ul>
    </aside>
</body>

</html>
`
)
