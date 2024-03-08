/*
 *  Copyright (c) 2022-2024 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package lint

import (
	"path/filepath"

	"github.com/osspkg/devtool/internal/global"
	"github.com/osspkg/devtool/pkg/exec"
	"github.com/osspkg/devtool/pkg/files"
	"go.osspkg.com/goppy/console"
)

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("lint", "")
		setter.ExecFunc(func(_ []string) {
			global.SetupEnv()
			console.Infof("--- LINT ---")

			cmds := make([]string, 0, 50)
			cmds = append(cmds, "golangci-lint --version")
			if files.Exist(files.CurrentDir() + "/go.work") {
				cmds = append(cmds, "go work use -r .", "go work sync")
				mods, err := files.Detect("go.mod")
				console.FatalIfErr(err, "detects go.mod in workspace")
				for _, mod := range mods {
					dir := filepath.Dir(mod)
					cmds = append(cmds,
						"cd "+dir+" && golangci-lint -v run ./...",
					)
				}
			} else {
				cmds = append(cmds,
					"golangci-lint -v run ./...",
				)
			}

			exec.CommandPack("bash", cmds...)
		})
	})
}
