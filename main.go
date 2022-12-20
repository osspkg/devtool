package main

import (
	"os"

	"github.com/dewep-online/devtool/internal/setup"

	"github.com/dewep-online/devtool/internal/build"

	"github.com/dewep-online/devtool/internal/lint"
	"github.com/dewep-online/devtool/internal/tests"
	"github.com/deweppro/go-app/console"
)

func main() {
	console.ShowDebug(true)

	app := console.New("devtool", "help devtool")

	app.RootCommand(console.NewCommand(func(setter console.CommandSetter) {
		setter.ExecFunc(func(_ []string) {
			console.Infof("os env:\n%s", func() string {
				out := ""
				for _, s := range os.Environ() {
					out += s + "\n"
				}
				return out
			}())
		})
	}))

	app.AddCommand(
		setup.Cmd(),
		lint.Cmd(),
		tests.Cmd(),
		build.Cmd(),
	)

	app.Exec()
}
