// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package warcraft // import "github.com/wabarc/warcraft"

import (
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

// Download
// wget --delete-after --no-directories --warc-file=google --recursive --level=1 URI
func (warc *Warcraft) Download(u *url.URL) (string, error) {
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

	name := filepath.Join(warc.BasePath, strings.TrimSuffix(helper.FileName(u.String(), ""), ".html"))
	args := []string{
		"--delete-after", "--no-directories",
		"--recursive", "--level=1",
		"--warc-file=" + name,
		u.String(),
	}
	cmd := exec.Command(binPath, args...)
	if err := cmd.Start(); err != nil {
		return "", err
	}
	if err := cmd.Wait(); err != nil {
		return "", err
	}

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
