package global

import (
	"os"
	"os/exec"
	"regexp"

	"github.com/dewep-online/devtool/pkg/files"
	"github.com/deweppro/go-sdk/console"
)

const (
	ToolsDir   = ".tools"
	BuildDir   = "build"
	InitDir    = "init"
	ScriptsDir = "scripts"
)

func GetToolsDir() string {
	return files.CurrentDir() + "/" + ToolsDir
}

func GetBuildDir() string {
	return files.CurrentDir() + "/" + BuildDir
}

func GetInitDir() string {
	return files.CurrentDir() + "/" + InitDir
}

func GetScriptsDir() string {
	return files.CurrentDir() + "/" + ScriptsDir
}

func SetupEnv() {
	console.FatalIfErr(os.Setenv("GOBIN", GetToolsDir()), "setup env")
	console.FatalIfErr(os.Setenv("PATH", GetToolsDir()+":"+os.Getenv("PATH")), "setup env")
}

var rex = regexp.MustCompile(`go(\d+)\.(\d+)`)

func GoVersion() string {
	b, err := exec.Command("bash", "-c", "go version").CombinedOutput()
	console.FatalIfErr(err, "detect go version")
	result := rex.FindAllString(string(b), 1)
	for _, s := range result {
		return s
	}
	return "unknown"
}
