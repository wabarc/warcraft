// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package warcraft // import "github.com/wabarc/warcraft"

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/wabarc/helper"
)

var userAgent = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/85.0.4183.83 Safari/537.36"

// Warcraft represents warcraft config.
type Warcraft struct {
	BasePath string // base path of warc file, defaults to current directory
	Verbose  bool

	userAgent string
}

// New a Warcraft struct
func New() *Warcraft {
	pwd, _ := os.Getwd()

	return &Warcraft{
		BasePath:  pwd,
		userAgent: userAgent,
	}
}

// UserAgent set User-Agent for wget
func (warc *Warcraft) UserAgent(s string) *Warcraft {
	if s != "" {
		warc.userAgent = s
	}
	return warc
}

// Download webpage as warc via wget
func (warc *Warcraft) Download(ctx context.Context, u *url.URL) (string, error) {
	if warc.BasePath == "" {
		pwd, err := os.Getwd()
		if err != nil {
			return "", err
		}
		warc.BasePath = pwd
	}
	if !helper.IsDir(warc.BasePath) {
		return "", errors.New(warc.BasePath + " is invalid")
	}
	if err := helper.Writable(warc.BasePath); err != nil {
		return "", errors.Wrap(err, "no writable")
	}

	binPath, err := findWgetExecPath()
	if err != nil {
		return "", err
	}

	// WGET CLI Docs: https://www.gnu.org/software/wget/manual/wget.html
	name := filepath.Join(warc.BasePath, strings.TrimSuffix(helper.FileName(u.String(), ""), ".html"))
	args := []string{
		"--no-config", "--no-directories", "--no-netrc", "--no-check-certificate", "--no-hsts", "--no-parent",
		"--adjust-extension", "--convert-links", "--delete-after", "--span-hosts", "--random-wait",
		"-e robots=off", "--page-requisites", "--header=Accept-Encoding: *",
		"--quiet=" + warc.quiet(), "--user-agent=" + warc.userAgent,
		fmt.Sprintf("--referer=%s://%s", u.Scheme, u.Hostname()),
		"--warc-tempdir=" + warc.BasePath,
		"--warc-file=" + name,
		u.String(),
	}
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Dir = warc.BasePath
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return "", err
	}
	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// Ignore server issued error response
			if exitError.ExitCode() != 8 {
				return "", exitError
			}
		}
	}
	if warc.Verbose {
		readOutput(stdout)
	}

	// First wait for the process to be finished.
	// Don't care about this error in any scenario.
	_ = cmd.Wait()

	// For WARC Archive version 1.0
	dst := name + ".warc"
	if helper.Exists(dst) {
		return dst, nil
	}
	dst += ".gz"

	return dst, nil
}

func (warc *Warcraft) quiet() string {
	if warc.Verbose {
		return "off"
	}
	return "on"
}

func findWgetExecPath() (string, error) {
	var locations []string
	switch runtime.GOOS {
	case "darwin":
		locations = []string{
			// Mac
			"wget",
			"/usr/local/bin/wget",
		}
	case "windows":
		locations = []string{
			// Windows
			"wget",
			"wget.exe", // in case PATHEXT is misconfigured
		}
	default:
		locations = []string{
			// Unix-like
			"wget",
			"/usr/bin/wget",
		}
	}

	for _, path := range locations {
		found, err := exec.LookPath(path)
		if err == nil {
			return found, nil
		}
	}

	return "", errors.New("wget not found")
}

func readOutput(rc io.ReadCloser) {
	for {
		out := make([]byte, 1024)
		_, err := rc.Read(out)
		fmt.Print(string(out))
		if err != nil {
			break
		}
	}
}
