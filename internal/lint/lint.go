package lint

import (
	"github.com/dewep-online/devtool/internal/global"
	"github.com/dewep-online/devtool/pkg/exec"
	"github.com/deweppro/go-sdk/console"
)

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("lint", "")
		setter.ExecFunc(func(_ []string) {
			global.SetupEnv()
			console.Infof("--- LINT ---")

			exec.CommandPack("bash",
				"golangci-lint --version",
				"golangci-lint -v run ./...",
			)
		})
	})
}
