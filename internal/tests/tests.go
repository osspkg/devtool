/*
 *  Copyright (c) 2022-2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package tests

import (
	"os"

	"github.com/osspkg/devtool/internal/global"
	"github.com/osspkg/devtool/pkg/exec"
	"go.osspkg.com/goppy/sdk/console"
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
