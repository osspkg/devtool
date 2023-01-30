package build

import (
	"strings"

	"github.com/dewep-online/devtool/internal/global"
	"github.com/dewep-online/devtool/pkg/exec"
	"github.com/dewep-online/devtool/pkg/files"
	"github.com/deweppro/go-app/console"
)

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("build", "")
		setter.Flag(func(flagsSetter console.FlagsSetter) {
			flagsSetter.StringVar("arch", "amd64,arm64", "")
		})
		setter.ExecFunc(func(_ []string, arch string) {
			console.Infof("setup tools")
			buildDir := global.GetBuildDir()
			global.SetupEnv()

			cmds := []string{
				"go mod tidy",
				"go mod download",
				"go generate ./...",
			}

			mainFiles, err := files.Detect("main.go")
			console.FatalIfErr(err, "detect main.go")

			aarch64 := files.Exist("/usr/bin/aarch64-linux-gnu-gcc")

			for _, main := range mainFiles {
				appName := files.Folder(main)
				archList := strings.Split(arch, ",")

				for _, arch = range archList {
					cmds = append(cmds, "rm -rf "+buildDir+"/"+appName+"_"+arch)

					switch arch {
					case "arm64":
						if aarch64 {
							cmds = append(cmds, "GODEBUG=netdns=9 GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=arm64 CC=aarch64-linux-gnu-gcc go build -a -o "+buildDir+"/"+appName+"_"+arch+" "+main)
						} else {
							cmds = append(cmds, "GODEBUG=netdns=9 GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -a -o "+buildDir+"/"+appName+"_"+arch+" "+main)
						}
					case "amd64":
						cmds = append(cmds, "GODEBUG=netdns=9 GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o "+buildDir+"/"+appName+"_"+arch+" "+main)
					default:
						console.Fatalf("use only arch=[amd64,arm64]")
					}
				}
			}

			exec.CommandPack("bash", cmds...)
		})
	})
}
