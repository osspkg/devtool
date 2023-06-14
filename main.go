package main

import (
	"os"

	"github.com/osspkg/devtool/internal/build"
	"github.com/osspkg/devtool/internal/global"
	"github.com/osspkg/devtool/internal/lint"
	"github.com/osspkg/devtool/internal/setup"
	"github.com/osspkg/devtool/internal/tests"
	"github.com/osspkg/go-sdk/console"
)

func main() {
	console.ShowDebug(true)

	app := console.New("devtool", "help devtool")

	app.RootCommand(console.NewCommand(func(setter console.CommandSetter) {
		setter.ExecFunc(func(_ []string) {
			global.SetupEnv()
			console.Rawf("os env:\n%s", func() string {
				out := ""
				for _, s := range os.Environ() {
					out += s + "\n"
				}
				return out
			}())
		})
	}))

	app.AddCommand(
		setup.CmdApp(),
		setup.CmdLib(),
		lint.Cmd(),
		tests.Cmd(),
		build.Cmd(),
	)

	app.Exec()
}
