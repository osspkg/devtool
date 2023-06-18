/*
 *  Copyright (c) 2022-2023 Mikhail Knyzhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package lic

import (
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/osspkg/devtool/internal/global"
	"github.com/osspkg/devtool/pkg/exec"
	"github.com/osspkg/devtool/pkg/files"
	"github.com/osspkg/go-sdk/console"
)

const (
	gitFirstCommitDate = `git log --reverse --date="format:%Y" --format="format:%ad" | head -n 1`
	licFilename        = ".lict"
)

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("license", "")
		setter.ExecFunc(func(_ []string) {
			global.SetupEnv()
			console.Infof("--- LICENSE ---")

			model := &Lic{}
			if !files.Exist(licFilename) {
				model.Default()
				console.FatalIfErr(files.YamlWrite(licFilename, model), "Create `%s`", licFilename)
				console.Errorf("Update `%s`", licFilename)
				os.Exit(1)
			} else {
				console.FatalIfErr(files.YamlRead(licFilename, model), "Read `%s`", licFilename)
			}

			out, err := exec.SingleCmd(context.TODO(), "bash", gitFirstCommitDate)
			console.FatalIfErr(err, "Get fist commit")

			start, err := time.Parse("2006", strings.TrimSpace(string(out)))
			console.FatalIfErr(err, "Get fist commit date parse")

			startY := start.Year()
			currY := time.Now().Year()
			model.Years = fmt.Sprintf("%d-%d", startY, currY)
			if startY >= currY {
				model.Years = fmt.Sprintf("%d", currY)
			}

			goFiles, err := files.DetectByExt(".go")
			console.FatalIfErr(err, "Get go files")

			for _, file := range goFiles {
				err = files.Rewrite(file, func(s string) string {
					tmpl := buildTemplate(model)
					return replaceLic(s, tmpl)
				})
				console.FatalIfErr(err, "Update go file `%s`", file)
			}

			exec.CommandPack("bash",
				"go fmt ./...",
			)
		})
	})
}

type Lic struct {
	Author   string `yaml:"author"`
	LicShort string `yaml:"lic_short"`
	LicFile  string `yaml:"lic_file"`
	Years    string `yaml:"-"`
}

func (l *Lic) Default() {
	l.Author = "User <user@email>"
	l.LicShort = "MIT"
	l.LicFile = "LICENSE"
}

const template = "/*\n" +
	" *  Copyright (c) %s %s. All rights reserved.\n" +
	" *  Use of this source code is governed by a %s license that can be found in the %s file.\n" +
	" */\n\n"

func buildTemplate(m *Lic) string {
	return fmt.Sprintf(template, m.Years, m.Author, m.LicShort, m.LicFile)
}

var rexLic = regexp.MustCompile(`(?miUs)^\/\*(.*Copyright.*)\*\/`)

func replaceLic(b string, n string) string {
	if rexLic.MatchString(b) {
		return rexLic.ReplaceAllString(b, n)
	}
	return n + b
}
