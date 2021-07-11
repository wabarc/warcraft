package main

import (
	"context"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"sync"

	"github.com/wabarc/warcraft"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\n")
		fmt.Fprintf(os.Stderr, "  warcraft [options] [url1] ... [urlN]\n")

		flag.PrintDefaults()
	}
	var basePrint = func() {
		fmt.Print("A CLI tool help download webpage as warc file using wget.\n\n")
		flag.Usage()
		fmt.Fprint(os.Stderr, "\n")
	}

	flag.Parse()

	args := flag.Args()

	if len(args) < 1 {
		basePrint()
		os.Exit(0)
	}

}

func main() {
	uris := flag.Args()
	warc := warcraft.New()

	pwd, _ := os.Getwd()

	var wg sync.WaitGroup
	for _, uri := range uris {
		wg.Add(1)
		go func(uri string) {
			in, err := url.Parse(uri)
			if err != nil {
				fmt.Fprintf(os.Stderr, "parse %s failed: %v\n", uri, err)
				return
			}

			if path, err := warc.Download(context.Background(), in); err != nil {
				fmt.Fprintf(os.Stderr, "warcraft: %v\n", err)
			} else {
				fmt.Fprintf(os.Stdout, "%s  %s\n", strings.TrimLeft(path, pwd), uri)
			}
			wg.Done()
		}(uri)
	}
	wg.Wait()
}
