// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package warcraft // import "github.com/wabarc/warcraft"

import (
	"context"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/pkg/errors"
	"github.com/wabarc/helper"
)

// Warcraft represents warcraft config.
type Warcraft struct {
	BasePath string // base path of warc file, defaults to current directory
}

// New a Warcraft struct
func New() *Warcraft {
	pwd, _ := os.Getwd()

	return &Warcraft{
		BasePath: pwd,
	}
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
		"--no-config", "--no-directories", "--no-verbose", "--no-netrc", "--no-check-certificate",
		"--no-hsts", "--no-parent", "--timestamping", "--adjust-extension", "--convert-links",
		"--span-hosts", "--delete-after", "--tries=3", "--compression=auto", "-e robots=off",
		"--page-requisites", "--warc-tempdir=" + warc.BasePath, "--warc-file=" + name,
		u.String(),
	}
	cmd := exec.CommandContext(ctx, binPath, args...)
	cmd.Dir = warc.BasePath
	// _, err = cmd.StdoutPipe()
	// if err != nil {
	// 	return "", err
	// }
	// cmd.Stderr = cmd.Stdout

	// We must start the cmd before calling cmd.Wait, as otherwise the two
	// can run into a data race.
	if err := cmd.Start(); err != nil {
		return "", errors.Wrap(err, "starts wget failed")
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
