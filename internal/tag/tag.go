package tag

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/osspkg/devtool/internal/global"
	"github.com/osspkg/devtool/pkg/exec"
	"github.com/osspkg/devtool/pkg/files"
	"github.com/osspkg/devtool/pkg/modules"
	"github.com/osspkg/devtool/pkg/repo"
	"github.com/osspkg/devtool/pkg/ver"
	"go.osspkg.com/algorithms/graph/kahn"
	"go.osspkg.com/goppy/sdk/console"
	"golang.org/x/mod/modfile"
)

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("tag", "")
		setter.Flag(func(fs console.FlagsSetter) {
			fs.Bool("minor", "update minor version (default - patch)")
		})
		setter.ExecFunc(func(_ []string, minor bool) {
			global.SetupEnv()
			console.Infof("--- READ CONFIG ---")

			var (
				allMods map[string]*modules.Mod
				currmod *modules.Mod
				err     error
				b       []byte
				f       *modfile.File
				fi      os.FileInfo
				HEAD    string
			)

			console.Infof("--- GET ALL MODULES ---")

			allMods, err = modules.Detect(files.CurrentDir())
			console.FatalIfErr(err, "Detect go.mod files")

			var root *modules.Mod
			for _, m := range allMods {
				if root == nil {
					root = m
					continue
				}
				if len(root.Name) > len(m.Name) {
					root = m
				}
			}
			for _, m := range allMods {
				m.Prefix = strings.Trim(strings.TrimPrefix(m.Name, root.Name), "/")
				if len(m.Prefix) > 0 {
					m.Prefix += "/"
				}
				b, err = exec.SingleCmd(context.TODO(), "bash", "git tag -l \""+m.Prefix+"v*\"")
				console.FatalIfErr(err, "Get tags for: %s", m.Name)
				m.Version = ver.Max(strings.Split(string(b), "\n")...)
			}

			console.Infof("--- DETECT CHANGES ---")

			HEAD, err = repo.HEAD("")
			console.FatalIfErr(err, "Get git HEAD")
			b, err = exec.SingleCmd(context.TODO(), "bash", "git diff --name-only --ignore-submodules --diff-filter=ACMRTUXB origin/"+HEAD+" -- \"*.go\" \"*.mod\" \"*.sum\"")
			console.FatalIfErr(err, "Detect changed files")
			changedFiles := strings.Split(string(b), "\n")
			for _, file := range changedFiles {
				if len(file) == 0 {
					continue
				}
				dir := filepath.Dir(file)
				isRoot := dir == "."
				for {
					currmod, err = modules.Read(dir + "/go.mod")
					if err != nil && !isRoot && strings.Contains(err.Error(), "no such file") {
						dir = filepath.Dir(dir)
						if dir != "." {
							continue
						}
						break
					}
					break
				}
				if err != nil {
					continue
				}

				for _, m := range allMods {
					if m.Name == currmod.Name && !m.Changed {
						m.Changed = true
						if minor {
							m.Version.Minor++
							m.Version.Patch = 0
						} else {
							m.Version.Patch++
						}
					}
				}
			}

			console.Infof("--- UPDATE MODULES ---")

			graph := kahn.New()
			for _, m := range allMods {
				_, err = os.Stat(m.File)
				console.FatalIfErr(err, "Get info go.mod file: %s", m.File)
				b, err = os.ReadFile(m.File)
				console.FatalIfErr(err, "Read go.mod file: %s", m.File)
				_, err = modfile.Parse(m.File, b, func(path, version string) (string, error) {
					if _, ok := allMods[path]; ok {
						console.FatalIfErr(graph.Add(path, m.Name), "Create graph from: %s", m.File)
					}
					return version, nil
				})
				console.FatalIfErr(err, "Parse go.mod file: %s", m.File)
			}
			console.FatalIfErr(graph.Build(), "Build graph")
			for _, s := range graph.Result() {
				m, ok := allMods[s]
				if !ok {
					continue
				}
				fmt.Println(">", m.Name)
				fi, err = os.Stat(m.File)
				console.FatalIfErr(err, "Get info go.mod file: %s", m.File)
				b, err = os.ReadFile(m.File)
				console.FatalIfErr(err, "Read go.mod file: %s", m.File)
				f, err = modfile.Parse(m.File, b, func(path, version string) (string, error) {
					if mm, ok := allMods[path]; ok && mm.Version.String() != version {
						if !m.Changed {
							if minor {
								m.Version.Minor++
								m.Version.Patch = 0
							} else {
								m.Version.Patch++
							}
							m.Changed = true
						}
						return mm.Version.String(), nil
					}
					return version, nil
				})
				console.FatalIfErr(err, "Parse go.mod file: %s", m.File)
				b, err = f.Format()
				console.FatalIfErr(err, "Format go.mod file: %s", m.File)
				err = os.WriteFile(m.File, b, fi.Mode())
				console.FatalIfErr(err, "Update go.mod file: %s", m.File)
			}

			console.Infof("--- GIT COMMITTED ---")

			cmds := make([]string, 0, 50)
			cmds = append(cmds,
				"git add .",
				"git commit -m \"update vendors\"",
			)
			for _, m := range allMods {
				if !m.Changed {
					continue
				}
				cmds = append(cmds, "git tag "+m.Prefix+m.Version.String())
			}
			cmds = append(cmds, "git push", "git push --tags")
			exec.CommandPack("bash", cmds...)
		})
	})
}
