package build

import (
	"strings"

	"github.com/dewep-online/devtool/internal/global"
	"github.com/dewep-online/devtool/pkg/exec"
	"github.com/dewep-online/devtool/pkg/files"
	"github.com/deweppro/go-sdk/console"
)

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("build", "")
		setter.Flag(func(flagsSetter console.FlagsSetter) {
			flagsSetter.StringVar("arch", "amd64,arm64", "")
		})
		setter.ExecFunc(func(_ []string, arch string) {
			global.SetupEnv()
			console.Infof("--- BUILD ---")

			pack := make([]string, 0)
			buildDir := global.GetBuildDir()

			mainFiles, err := files.Detect("main.go")
			console.FatalIfErr(err, "detect main.go")

			for _, main := range mainFiles {
				appName := files.Folder(main)
				archList := strings.Split(arch, ",")

				for _, arch = range archList {
					pack = append(pack, "rm -rf "+buildDir+"/"+appName+"_"+arch)

					chunk := []string{
						"GODEBUG=netdns=9",
						"GO111MODULE=on",
						"CGO_ENABLED=1",
					}

					switch arch {
					case "arm64":
						chunk = append(chunk, "GOOS=linux", "GOARCH=arm64")

						if exist("/usr/bin/aarch64-linux-gnu-gcc") {
							chunk = append(chunk, "CC=aarch64-linux-gnu-gcc")
						}

					case "amd64":
						chunk = append(chunk, "GOOS=linux", "GOARCH=amd64")

					default:
						console.Fatalf("use only arch=[amd64,arm64]")
					}

					chunk = append(chunk, "go build -ldflags=\"-s -w\" -a -o "+buildDir+"/"+appName+"_"+arch+" "+main)
					pack = append(pack, strings.Join(chunk, " "))
				}
			}

			exec.CommandPack("bash", pack...)
		})
	})
}

func exist(filename string) bool {
	return files.Exist(filename)
}
