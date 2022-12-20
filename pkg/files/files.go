package files

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/deweppro/go-app/console"
)

func CurrentDir() string {
	dir, err := os.Getwd()
	console.FatalIfErr(err, "get current dir")
	return dir
}

func Detect(filename string) ([]string, error) {
	curDir := CurrentDir()
	files := make([]string, 0)
	err := filepath.Walk(curDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Name() != filename {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

func Rewrite(filename string, cb func(string) string) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	b = []byte(cb(string(b)))

	return os.WriteFile(filename, b, 0755)
}

func Exist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func Folder(filename string) string {
	dir := filepath.Dir(filename)
	tree := strings.Split(dir, string(os.PathSeparator))
	return tree[len(tree)-1]
}
