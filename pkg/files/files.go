/*
 *  Copyright (c) 2022-2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package files

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"go.osspkg.com/goppy/sdk/console"
	"gopkg.in/yaml.v3"
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

func DetectByExt(ext string) ([]string, error) {
	curDir := CurrentDir()
	files := make([]string, 0)
	err := filepath.Walk(curDir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || filepath.Ext(info.Name()) != ext {
			return nil
		}
		files = append(files, path)
		return nil
	})
	return files, err
}

func Rewrite(filename string, cb func(string) string) error {
	if !Exist(filename) {
		if err := os.WriteFile(filename, []byte(""), 0755); err != nil {
			return err
		}
	}
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

func YamlRead(filename string, v interface{}) error {
	b, err := os.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.Unmarshal(b, v)
}

func YamlWrite(filename string, v interface{}) error {
	b, err := yaml.Marshal(v)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, b, 0755)
}
