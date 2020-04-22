**Blindingly Simple Protocol Language ([BSPL](https://confluence.oceanobservatories.org/download/attachments/28809860/AAMAS-11-IBIOP.pdf))** Go parser.

[![Build Status](https://travis-ci.com/mikelsr/bspl.svg?token=736yMuj6XUy7yCEvSpBB&branch=master)](https://travis-ci.com/mikelsr/bspl)

This repository also contains interfaces for a BSPL reasoner (`reason` package) and an implementation of some components of that reasoner (`implementation` package).
This implementations are used in [another project](https://github.com/mikelsr/nahs).

* `parser`: Standalone BSPL parser implemented using [a toy lexer](https://github.com/mikelsr/gauzaez) I wrote a while ago.

* `reason`: Interface definition for implementing a reasoner and protocol instances.

* `implementation`: Draft implementation to use in another project.

Production use of this project is not advised as it is far from ready.

## Usage example

1. Define a valid protocol in a file with path `path`.

2. Open the file and pass the reader to bspl.Parse()

```go
package main

import (
        "fmt"
        "os"

        "github.com/mikelsr/bspl"
)

func main() {
	source, err := os.Open(path)
	if err != nil {
		panic(err)
	}
        protocol, err := bspl.Parse(source)
        if err != nil {
		panic(err)
        }
        fmt.Println(protocol)
}
```

3. Done!