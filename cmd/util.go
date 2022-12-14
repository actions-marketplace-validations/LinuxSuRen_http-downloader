package cmd

import (
	"context"
	"github.com/spf13/cobra"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"
)

func getOrDefault(key, def string, data map[string]string) (result string) {
	var ok bool
	if result, ok = data[key]; !ok {
		result = def
	}
	return
}

type contextRoundTripper string

func getRoundTripper(ctx context.Context) (tripper http.RoundTripper) {
	if ctx == nil {
		return
	}
	roundTripper := ctx.Value(contextRoundTripper("roundTripper"))

	switch v := roundTripper.(type) {
	case *http.Transport:
		tripper = v
	}
	return
}

func execCommandInDir(name, dir string, arg ...string) (err error) {
	command := exec.Command(name, arg...)
	if dir != "" {
		command.Dir = dir
	}

	//var stdout []byte
	//var errStdout error
	stdoutIn, _ := command.StdoutPipe()
	stderrIn, _ := command.StderrPipe()
	err = command.Start()
	if err != nil {
		return err
	}

	// cmd.Wait() should be called only after we finish reading
	// from stdoutIn and stderrIn.
	// wg ensures that we finish
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		_, _ = copyAndCapture(os.Stdout, stdoutIn)
		wg.Done()
	}()

	_, _ = copyAndCapture(os.Stderr, stderrIn)

	wg.Wait()

	err = command.Wait()
	return
}

func copyAndCapture(w io.Writer, r io.Reader) ([]byte, error) {
	var out []byte
	buf := make([]byte, 1024, 1024)
	for {
		n, err := r.Read(buf[:])
		if n > 0 {
			d := buf[:n]
			out = append(out, d...)
			_, err := w.Write(d)
			if err != nil {
				return out, err
			}
		}
		if err != nil {
			// Read returns io.EOF at the end of file, which is not an error for us
			if err == io.EOF {
				err = nil
			}
			return out, err
		}
	}
}

// CompletionFunc is the function for command completion
type CompletionFunc func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective)

// ArrayCompletion return a completion  which base on an array
func ArrayCompletion(array ...string) CompletionFunc {
	return func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return array, cobra.ShellCompDirectiveNoFileComp
	}
}
