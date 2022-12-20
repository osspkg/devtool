package setup

import (
	"os"
	"strings"

	"github.com/dewep-online/devtool/internal/global"
	"github.com/dewep-online/devtool/pkg/exec"
	"github.com/dewep-online/devtool/pkg/files"
	"github.com/deweppro/go-app/console"
)

func Cmd() console.CommandGetter {
	return console.NewCommand(func(setter console.CommandSetter) {
		setter.Setup("setup", "")
		setter.ExecFunc(func(_ []string) {
			global.SetupEnv()

			toolDir, initDir, scriptsDir := global.GetToolsDir(), global.GetInitDir(), global.GetScriptsDir()

			console.FatalIfErr(os.MkdirAll(initDir, 0755), "create init dir")
			console.FatalIfErr(os.MkdirAll(toolDir, 0755), "create tools dir")
			console.FatalIfErr(os.MkdirAll(scriptsDir, 0755), "create scripts dir")

			console.Infof("update .gitignore")
			console.FatalIfErr(files.Rewrite(files.CurrentDir()+"/.gitignore", func(s string) string {
				if !strings.Contains(s, global.ToolsDir) {
					s += global.ToolsDir + "/\n"
				}
				if !strings.Contains(s, global.BuildDir) {
					s += global.BuildDir + "/\n"
				}
				return s
			}), "create tools dir")

			console.Infof("install tools")
			for name, install := range tools1 {
				if !files.Exist(toolDir + "/" + name) {
					console.FatalIfErr(exec.Command("bash", install), "install tool [%s]", name)
				}
			}

			gover := global.GoVersion()
			tools, ok := tools2[gover]
			if ok {
				for name, install := range tools {
					if !files.Exist(toolDir + "/" + name) {
						console.FatalIfErr(exec.Command("bash", install), "install tool [%s]", name)
					}
				}
			}

			console.Infof("create configs")
			for name, config := range configs {
				if !files.Exist(files.CurrentDir() + "/" + name) {
					console.FatalIfErr(os.WriteFile(files.CurrentDir()+"/"+name, []byte(config), 0755), "create config [%s]", name)
				}
			}

			console.Infof("create services and deb scripts")
			postinstData, postrmData, preinstData, prermData := bashPrefix, bashPrefix, bashPrefix, bashPrefix

			mainFiles, err := files.Detect("main.go")
			console.FatalIfErr(err, "detect main.go")
			for _, main := range mainFiles {
				appName := files.Folder(main)
				if !files.Exist(initDir + "/" + appName + ".service") {
					tmpl := strings.ReplaceAll(systemctlConfig, "{%app_name%}", appName)
					console.FatalIfErr(os.WriteFile(initDir+"/"+appName+".service", []byte(tmpl), 0755), "create init config [%s]", appName)
				}

				postinstData += strings.ReplaceAll(postinst, "{%app_name%}", appName)
				preinstData += strings.ReplaceAll(preinstDir, "{%app_name%}", appName)
				preinstData += strings.ReplaceAll(preinst, "{%app_name%}", appName)
				prermData += strings.ReplaceAll(prerm, "{%app_name%}", appName)
			}

			if !files.Exist(scriptsDir + "/postinst.sh") {
				console.FatalIfErr(os.WriteFile(scriptsDir+"/postinst.sh", []byte(postinstData), 0755), "create postinst")
			}
			if !files.Exist(scriptsDir + "/postrm.sh") {
				console.FatalIfErr(os.WriteFile(scriptsDir+"/postrm.sh", []byte(postrmData), 0755), "create postrm")
			}
			if !files.Exist(scriptsDir + "/preinst.sh") {
				console.FatalIfErr(os.WriteFile(scriptsDir+"/preinst.sh", []byte(preinstData), 0755), "create preinst")
			}
			if !files.Exist(scriptsDir + "/prerm.sh") {
				console.FatalIfErr(os.WriteFile(scriptsDir+"/prerm.sh", []byte(prermData), 0755), "create prerm")
			}

		})
	})
}

var tools1 = map[string]string{
	"goveralls": "go install github.com/mattn/goveralls@latest",
	"static":    "go install github.com/deweppro/go-static/cmd/static@latest",
	"easyjson":  "go install github.com/mailru/easyjson/...@latest",
}

var tools2 = map[string]map[string]string{
	"go1.19": {
		"golangci-lint": "go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.50.0",
	},
	"go1.18": {
		"golangci-lint": "go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0",
	},
	"go1.17": {
		"golangci-lint": "go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.38.0",
	},
}

var configs = map[string]string{
	".golangci.yml": golangciLintConfig,
}

var golangciLintConfig = `
run:
  concurrency: 1
  deadline: 2m
  issues-exit-code: 1
  tests: true
  skip-files:
    - easyjson

issues:
  exclude-use-default: false

output:
  format: colored-line-number
  print-issued-lines: true
  print-linter-name: true

linters-settings:
  govet:
    check-shadowing: true
  golint:
    min-confidence: 0.8
  gofmt:
    simplify: true
  errcheck:
    check-type-assertions: true
    check-blank: true
  gocyclo:
    min-complexity: 30
  misspell:
    locale: US
  gosimple:
    go: "1.16"
    checks: ["all"]
  prealloc:
    simple: true
    range-loops: true
    for-loops: false

linters:
  disable-all: true
  enable:
    - govet
    - gofmt
    - errcheck
    - misspell
    - gocyclo
    - ineffassign
    - goimports
    - gosimple
    - prealloc
  fast: false

`
var systemctlConfig = `[Unit]
After=network.target

[Service]
User=root
Group=root
Restart=on-failure
RestartSec=30s
Type=simple
ExecStart=/usr/bin/{%app_name%} --config=/etc/{%app_name%}/config.yaml
KillMode=process
KillSignal=SIGTERM

[Install]
WantedBy=default.target
`

var (
	bashPrefix = "#!/bin/bash\n\n"
	postinst   = `
if [ -f "/etc/systemd/system/{%app_name%}.service" ]; then
    systemctl start {%app_name%}
    systemctl enable {%app_name%}
    systemctl daemon-reload
fi
`
	preinstDir = `
if ! [ -d /var/lib/{%app_name%}/ ]; then
    mkdir /var/lib/{%app_name%}
fi
`
	preinst = `
if [ -f "/etc/systemd/system/{%app_name%}.service" ]; then
    systemctl stop {%app_name%}
    systemctl disable {%app_name%}
    systemctl daemon-reload
fi
`
	prerm = `
if [ -f "/etc/systemd/system/{%app_name%}.service" ]; then
    systemctl stop {%app_name%}
    systemctl disable {%app_name%}
    systemctl daemon-reload
fi
`
)
