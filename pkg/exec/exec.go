package exec

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"sync"

	"github.com/dewep-online/devtool/pkg/files"
	"github.com/deweppro/go-sdk/console"
	"github.com/deweppro/go-sdk/syscall"
)

func CommandPack(shell string, command ...string) {
	for _, s := range command {
		console.FatalIfErr(Command(shell, s), "run command")
	}
}

func Command(shell string, command string) (err error) {
	wg := sync.WaitGroup{}
	ctx, cncl := context.WithCancel(context.Background())
	go syscall.OnStop(func() {
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
			console.Rawf(scanner.Text())
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()
	go func() {
		scanner := bufio.NewScanner(stderr)
		for scanner.Scan() {
			console.Rawf(scanner.Text())
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return cmd.Wait()
}
