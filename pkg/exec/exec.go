/*
 *  Copyright (c) 2022-2023 Mikhail Knyazhev <markus621@yandex.ru>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package exec

import (
	"bufio"
	"context"
	"os"
	"os/exec"
	"sync"

	"github.com/osspkg/devtool/pkg/files"
	"go.osspkg.com/goppy/sdk/console"
	"go.osspkg.com/goppy/sdk/syscall"
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

func SingleCmd(ctx context.Context, shell string, command string) ([]byte, error) {
	console.Infof(command)
	cmd := exec.CommandContext(ctx, shell, "-c", command)
	cmd.Env = os.Environ()
	cmd.Dir = files.CurrentDir()

	return cmd.CombinedOutput()
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
