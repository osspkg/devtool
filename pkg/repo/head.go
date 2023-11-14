package repo

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/osspkg/devtool/pkg/exec"
)

var rexHEAD = regexp.MustCompile(`(?mU)ref\: refs/heads/(\w+)\s+HEAD`)

func HEAD(url string) (string, error) {
	if len(url) == 0 {
		b, err := exec.SingleCmd(context.TODO(), "bash", "git remote get-url origin")
		if err != nil {
			return "", err
		}
		url = strings.Trim(string(b), "\n")
	}
	b, err := exec.SingleCmd(context.TODO(), "bash", "git ls-remote --symref "+url+" HEAD")
	if err != nil {
		return "", err
	}
	_strs := rexHEAD.FindStringSubmatch(string(b))
	if len(_strs) != 2 {
		return "", fmt.Errorf("HEAD not found")
	}
	return _strs[1], nil
}
