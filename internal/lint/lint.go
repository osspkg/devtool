package lint

import (
	"github.com/dewep-online/devtool/internal/global"
	"github.com/dewep-online/devtool/pkg/exec"
	"github.com/dewep-online/devtool/pkg/files"
	"github.com/deweppro/go-app/console"
)

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("lint", "")
		setter.ExecFunc(func(_ []string) {
			console.Infof("setup tools")
			global.SetupEnv()

			exec.CommandPack("bash",
				"go mod tidy",
				"go mod download",
				"go generate ./...",
				"golangci-lint --version",
				"golangci-lint -v run ./...",
			)

			mainFiles, err := files.Detect("main.go")
			console.FatalIfErr(err, "detect main.go")

			for _, main := range mainFiles {
				exec.CommandPack("bash", "go build -a -race -o /tmp/bin.test "+main)
			}
		})
	})
}
