/*
 *  Copyright (c) 2022-2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package lint

import (
	"github.com/osspkg/devtool/internal/global"
	"github.com/osspkg/devtool/pkg/exec"
	"go.osspkg.com/goppy/sdk/console"
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
