/*
 *  Copyright (c) 2022-2024 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package appgoppy

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/osspkg/devtool/internal/global"
	"go.osspkg.com/goppy/console"
	"go.osspkg.com/goppy/iofile"
	"go.osspkg.com/goppy/shell"
)

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("goppy", "Generate new app with goppy sdk")
		setter.ExecFunc(func(_ []string) {
			currdir := iofile.CurrentDir()

			data := make(map[string]interface{}, 100)
			data["go_version"] = strings.TrimLeft(global.GoVersion(), "go")
			data["app_module"] = console.Input("Input project name", nil, "app")
			data["app_name"] = func() string {
				vv := strings.Split(data["app_module"].(string), "/")
				return vv[len(vv)-1]
			}()

			for _, blocks := range modules {
				for _, name := range blocks {
					data["mod_"+name] = false
				}
			}

			userInput("Add modules", modules, "q", func(s string) {
				data["mod_"+s] = true
			})

			for _, folder := range folders {
				console.FatalIfErr(os.MkdirAll(currdir+"/"+folder, 0744), "Create folder")
			}

			for filename, tmpl := range templates {
				if strings.Contains(filename, "{{") {
					for key, value := range data {
						filename = strings.ReplaceAll(
							filename,
							"{{"+key+"}}",
							fmt.Sprintf("%+v", value),
						)
					}
				}
				writeFile(currdir+"/"+filename, tmpl, data)
			}

			sh := shell.New()
			sh.SetDir(currdir)
			sh.SetShell("bash")
			sh.SetWriter(os.Stdout)
			err := sh.CallPackageContext(context.TODO(),
				"gofmt -w .",
				"go mod tidy",
				"devtool setup-lib",
				"devtool setup-app",
			)
			console.FatalIfErr(err, "Call commands")
		})
	})
}

var modules = [][]string{
	{
		"metrics",
		"geoip",
		"oauth",
		"auth_jwt",
	},
	{
		"db_mysql",
		"db_sqlite",
		"db_postgre",
	},
	{
		"web_server",
		"web_client",
	},
	{
		"websocket_server",
		"websocket_client",
	},
	{
		"unixsocket_server",
		"unixsocket_client",
	},
	{
		"dns_server",
		"dns_client",
	},
}

var folders = []string{
	"app",
	"config",
	"cmd",
}

var templates = map[string]string{
	".gitignore":               tmplGitIgnore,
	"README.md":                tmplReadMe,
	"go.mod":                   tmplGoMod,
	"docker-compose.yaml":      tmplDockerFile,
	"cmd/{{app_name}}/main.go": tmplMainGO,
	"app/plugin.go":            tmplAppGo,
}

func writeFile(filename, t string, data map[string]interface{}) {
	tmpl, err := template.New("bt").Parse(t)
	console.FatalIfErr(err, "Parse template")
	tmpl.Option("missingkey=error")

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, data)
	console.FatalIfErr(err, "Build template")

	console.FatalIfErr(os.MkdirAll(filepath.Dir(filename), 0744), "Create folder")
	console.FatalIfErr(os.WriteFile(filename, buf.Bytes(), 0664), "Write %s", filename)
}

func userInput(msg string, mods [][]string, exit string, call func(s string)) {
	fmt.Printf("--- %s ---\n", msg)

	list := make(map[string]string, len(mods)*4)
	i := 0
	for _, blocks := range mods {
		for _, name := range blocks {
			i++
			fmt.Printf("(%d) %s, ", i, name)
			list[fmt.Sprintf("%d", i)] = name
		}
		fmt.Printf("\n")
	}
	fmt.Printf("and (%s) Done: \n", exit)

	scan := bufio.NewScanner(os.Stdin)
	for {
		if scan.Scan() {
			r := scan.Text()
			if r == exit {
				fmt.Printf("\u001B[1A\u001B[K--- Done ---\n\n")
				return
			}
			if name, ok := list[r]; ok {
				call(name)
				fmt.Printf("\033[1A\033[K + %s\n", name)
				continue
			}
			fmt.Printf("\u001B[1A\u001B[KBad answer! Try again: ")
		}
	}
}
