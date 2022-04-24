# bleve-server
[![GoDoc](https://godoc.org/github.com/gomarkdown/markdown?status.svg)](https://pkg.go.dev/github.com/bartmika/bleve-server)
[![Go Report Card](https://goreportcard.com/badge/github.com/bartmika/bleve-server)](https://goreportcard.com/report/github.com/bartmika/bleve-server)
[![License](https://img.shields.io/github/license/bartmika/bleve-server)](https://github.com/bartmika/bleve-server/blob/master/LICENSE)
![Go version](https://img.shields.io/github/go-mod/go-version/bartmika/bleve-server)

RPC server over the text Golang text indexing library [`bleve`](https://github.com/blevesearch/bleve).

## Purpose
Currently [`bleve`](https://github.com/blevesearch/bleve) does not support read/write access between more then programs according to [issue #1571](https://github.com/blevesearch/bleve/issues/1571); so that means, if you have one program accessing `bleve` then another program cannot access it.

As a result, this stand-alone server was created to allow multiple programs to access without problem. For your Golang program to access `bleve`, simply make `remote procedure calls` to this server.

## Installation
1. Clone the library to your computer.

  ```bash
  git clone https://github.com/bartmika/bleve-server.git
  cd bleve-server
  ```

2. Install the library dependencies.

  ```bash
  go mod tidy
  ```

3. Before you start the server, make sure you setup the following environment variable.

  ```bash
  export BLEVE_SERVER_ADDRESS=127.0.0.1:8001
  ```

4. Verify the server starts up.

  ```bash
  go run main.go serve
  ```

## Usage

```text
RPC server over a single running bleve instance

Usage:
  bleve-server [flags]
  bleve-server [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  index       Submit data to index
  query       Perform full-text search of indexed data
  register    Register a bleve index
  serve       Run the rpc server.
  version     Print the version number

Flags:
      --appAddress string   The applications address. (default "127.0.0.1:8001")
  -h, --help                help for bleve-server

Use "bleve-server [command] --help" for more information about a command.
```

## Usage Examples
### Command Line Interface
#### Start Server Command

1. Open up a terminal and register the *index* file before you begin. If *index* do not make sense then please read the [`bleve` developer documentation](https://github.com/blevesearch/bleve).

  ```bash
  go run main.go register --filename=dune.bleve
  ```

2. You have successfully created an *index*. Please note `bleve-server` supports multiple-index files open concurrently. Every time you want a new *index* you'll need to rerun the `regsiter` command.

3. Start `bleve-server` in your terminal:

  ```bash
  go run main.go serve
  ```

4. Now you are ready to make `rpc` calls to the server.

#### Index Data Command

While your `bleve-server` is running in terminal, open up another terminal run the following commands:

```bash
go run main.go index --filename=dune.bleve --identifier=123456789 --data="The spice extends life"
go run main.go index --filename=dune.bleve --identifier=987654321 --data="The spice is vital for space travel"
```

#### Query Command
Start by getting an individual result.

```bash
go run main.go query --filename=dune.bleve --search="life"

# OUTPUT:
# UUIDs: &[123456789]
```

Try another query to get an individual result.

```bash
go run main.go query --filename=dune.bleve --search="space travel"

# OUTPUT:
# UUIDs: &[987654321]
```

Try a query to get multiple results.

```bash
go run main.go query --filename=dune.bleve --search="spice"

# OUTPUT:
# UUIDs: &[123456789 987654321]
```

### Application

Here is a sample file of accessing `bleve-server` over `rpc` in your code:

```go
package main

import (
	"fmt"
	"log"
	"os"
	"time"

	remote "github.com/bartmika/bleve-server/pkg/rpc_client"
)

// Assume you ran the following commands in your `bleve-server` project:
//    go run main.go register --filename=dune.bleve
//    go run main.go serve

func main() {
	// Load up our `environment variables` from our operating system.
	addr := os.Getenv("BLEVE_SERVER_ADDRESS") // Example Value: 127.0.0.1:8001

	// Initialize the RPC client.
	rpc := remote.New(addr, 3, 15*time.Second)

	// Index the following data...

	err := rpc.Index("dune.bleve", "123456789", []byte("The spice extends life"))
	if err != nil {
		log.Fatal("doIndex err:", err)
	}
	err := rpc.Index("dune.bleve", "987654321", []byte("The spice is vital for space travel"))
	if err != nil {
		log.Fatal("doIndex err:", err)
	}

	// Try querrying...

	uuids, err := s.Query("dune.bleve", "life")
	if err != nil {
		log.Fatal("doQuery err:", err)
	}
	fmt.Println("UUIDs:", uuids) // OUTPUT: [123456789]

	uuids, err = s.Query("dune.bleve", "space travel")
	if err != nil {
		log.Fatal("doQuery err:", err)
	}
	fmt.Println("UUIDs:", uuids) // OUTPUT: [987654321]

	uuids, err = s.Query("dune.bleve", "spice")
	if err != nil {
		log.Fatal("doQuery err:", err)
	}
	fmt.Println("UUIDs:", uuids) // OUTPUT: [123456789, 987654321]
}
```

## Contributing

Found a bug? Want a feature to improve the package? Please create an [issue](https://github.com/bartmika/bleve-server/issues).

## License
Made with ❤️ by [Bartlomiej Mika](https://bartlomiejmika.com).   
The project is licensed under the [ISC License](LICENSE).

Resource used:

* [levesearch/bleve](https://github.com/blevesearch/bleve) is modern text indexing library for go.
* [spf13/cobra](https://github.com/spf13/cobra) is a commander for modern Go CLI interactions.
* [spf13/viper](https://github.com/spf13/viper) is a Go configuration with fangs.
