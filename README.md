# warcraft

`warcraft` is a toolkit to help download webpage as `warc` file using wget.

## Installation

The simplest, cross-platform way is to download from [GitHub Releases](https://github.com/wabarc/warcraft/releases) and place the executable file in your PATH.

Via Golang package get command

```sh
go get -u github.com/wabarc/warcraft/cmd/warcraft
```

From [gobinaries.com](https://gobinaries.com):

```sh
$ curl -sf https://gobinaries.com/wabarc/warcraft | sh
```

## Usage

Command-line:

```sh
$ warcraft
A CLI tool help download webpage as warc file using wget.

Usage:

  warcraft [options] [url1] ... [urlN]
```

Go package:
```go
import (
        "fmt"

        "github.com/wabarc/warcraft"
)

func main() {
        if b, err := warcraft.NewWarcraft(nil).Download(url); err != nil {
            fmt.Fprintf(os.Stderr, "warcraft: %v\n", err)
        } else {
            fmt.Fprintf(os.Stdout, "%s  %s\n", url, string(b))
        }
}
```

## License

This software is released under the terms of the MIT. See the [LICENSE](https://github.com/wabarc/warcraft/blob/main/LICENSE) file for details.
