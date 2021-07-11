// Copyright 2021 Wayback Archiver. All rights reserved.
// Use of this source code is governed by the MIT
// license that can be found in the LICENSE file.

package warcraft // import "github.com/wabarc/warcraft"

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"testing"

	"github.com/wabarc/helper"
)

func TestDownload(t *testing.T) {
	if _, err := findWgetExecPath(); err != nil {
		t.Skip(err.Error(), ", skipped")
	}

	_, mux, server := helper.MockServer()
	defer server.Close()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, Golang.")
	})

	uri := server.URL
	in, err := url.Parse(uri)
	if err != nil {
		t.Fatal(err)
	}

	dir, err := ioutil.TempDir("", "warcraft")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	warc := New()
	path, err := warc.Download(context.TODO(), in)
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(path)

	if !helper.Exists(path) {
		t.Log(path)
		t.Errorf(`download warc file failed`)
	}
}
