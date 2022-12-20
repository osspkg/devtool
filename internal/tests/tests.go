package tests

import (
	"os"

	"github.com/dewep-online/devtool/internal/global"
	"github.com/dewep-online/devtool/pkg/exec"
	"github.com/deweppro/go-app/console"
)

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("test", "")
		setter.ExecFunc(func(_ []string) {
			console.Infof("setup tools")
			toolsDir := global.GetToolsDir()
			global.SetupEnv()

			coverallsToken := os.Getenv("COVERALLS_TOKEN")

			cmds := []string{
				"go mod tidy",
				"go mod download",
				"go generate ./...",
				"go clean -testcache",
			}

			if len(coverallsToken) > 0 {
				cmds = append(cmds,
					"go test -v -race -run Unit -covermode=atomic -coverprofile=coverage.out ./...",
					toolsDir+"/goveralls -coverprofile=coverage.out -repotoken "+coverallsToken,
				)
			} else {
				cmds = append(cmds, "go test -v -race -run Unit ./...")
			}

			exec.CommandPack("bash", cmds...)
		})
	})
}
