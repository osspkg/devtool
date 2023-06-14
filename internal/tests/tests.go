package tests

import (
	"os"

	"github.com/osspkg/devtool/internal/global"
	"github.com/osspkg/devtool/pkg/exec"
	"github.com/osspkg/go-sdk/console"
)

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("test", "")
		setter.ExecFunc(func(_ []string) {
			global.SetupEnv()
			console.Infof("--- TESTS ---")

			pack := []string{
				"go clean -testcache",
				"go test -v -race -run Unit -covermode=atomic -coverprofile=coverage.out ./...",
			}

			coverallsToken := os.Getenv("COVERALLS_TOKEN")
			if len(coverallsToken) > 0 {
				pack = append(pack, "goveralls -coverprofile=coverage.out -repotoken "+coverallsToken)
			}

			exec.CommandPack("bash", pack...)
		})
	})
}
