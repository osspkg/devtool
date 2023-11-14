package modules

import (
	"os"
	"regexp"
	"strings"

	"github.com/osspkg/devtool/pkg/files"
	"github.com/osspkg/devtool/pkg/ver"
	"go.osspkg.com/goppy/sdk/errors"
)

var rex = regexp.MustCompile(`(?mU)module (.*)\n`)
var SkipErr = errors.New("skip read module")

type Mod struct {
	Name    string
	File    string
	Prefix  string
	Version *ver.Ver
	Changed bool
}

func Detect(dir string) (map[string]*Mod, error) {
	list := make(map[string]*Mod, 20)
	mods, err := files.DetectInDir(dir, "go.mod")
	if err != nil {
		return nil, err
	}
	var b []byte
	for _, mod := range mods {
		if b, err = os.ReadFile(mod); err != nil {
			return nil, err
		}

		temp := rex.FindStringSubmatch(string(b))
		if len(temp) != 2 {
			continue
		}
		module := temp[1]
		if !strings.Contains(module, "/") {
			continue
		}
		list[module] = &Mod{
			Name: module,
			File: mod,
		}
	}

	return list, nil
}

func Read(filepath string) (*Mod, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	temp := rex.FindStringSubmatch(string(b))
	if len(temp) != 2 {
		return nil, SkipErr
	}
	module := temp[1]
	if !strings.Contains(module, "/") {
		return nil, SkipErr
	}
	return &Mod{
		Name: module,
		File: filepath,
	}, nil
}
