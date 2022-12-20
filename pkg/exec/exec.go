package exec

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/exec"
	"sync"

	"github.com/dewep-online/devtool/pkg/files"
	"github.com/deweppro/go-app/application/sys"
	"github.com/deweppro/go-app/console"
)

func CommandPack(shell string, command ...string) {
	for _, s := range command {
		console.FatalIfErr(Command(shell, s), "run command")
	}
}

func Command(shell string, command string) (err error) {
	wg := sync.WaitGroup{}
	ctx, cncl := context.WithCancel(context.Background())
	go sys.OnSyscallStop(func() {
		cncl()
	})

	wg.Add(1)
	go func() {
		err = runCmd(ctx, shell, command)
		wg.Done()
	}()

	wg.Wait()
	return
}

func runCmd(ctx context.Context, shell string, command string) error {
	console.Infof(command)
	cmd := exec.CommandContext(ctx, shell, "-c", command)
	cmd.Env = os.Environ()
	cmd.Dir = files.CurrentDir()

	stdout, err := cmd.StdoutPipe()
	console.FatalIfErr(err, "stdout init")
	defer stdout.Close() //nolint: errcheck
	stderr, err := cmd.StderrPipe()
	console.FatalIfErr(err, "stderr init")
	defer stderr.Close() //nolint: errcheck

	console.FatalIfErr(cmd.Start(), "start command")

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			fmt.Println(scanner.Text())
			select {
			case <-ctx.Done():
				break
			default:
			}
		}
	}()

	return cmd.Wait()
}
