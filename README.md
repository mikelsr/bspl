**Blindingly Simple Protocol Language ([BSPL](https://confluence.oceanobservatories.org/download/attachments/28809860/AAMAS-11-IBIOP.pdf))** Go parser.

[![Build Status](https://travis-ci.com/mikelsr/bspl.svg?token=736yMuj6XUy7yCEvSpBB&branch=master)](https://travis-ci.com/mikelsr/bspl)
[![codecov](https://codecov.io/gh/mikelsr/bspl/branch/master/graph/badge.svg?token=ZKX6HOVW00)](https://codecov.io/gh/mikelsr/bspl)
[![License: MPL 2.0](https://img.shields.io/badge/License-MPL%202.0-brightgreen.svg)](https://opensource.org/licenses/MPL-2.0)
[![Go Version](https://img.shields.io/github/go-mod/go-version/mikelsr/bspl)](https://github.com/mikelsr/bspl/blob/master/go.mod)

This repository also contains interfaces for a BSPL reasoner (`reason` package) and an implementation of some components of that reasoner (`implementation` package).
This implementations are used in [another project](https://github.com/mikelsr/nahs).

## Modules

* `parser`: Standalone BSPL parser implemented using [a toy lexer](https://github.com/mikelsr/gauzaez) I wrote a while ago.

* `proto`: Go structures to form a BSPL protocol, e.g., `Protocol`, `Role` and `Action`.

* `reason`: Interface definition for implementing a reasoner and protocol instances.

* `implementation`: Draft implementation to use in another project.

Production use of this project is not advised as it is far from ready.

## Other folders

* `config`: Contains the automaton fed to the lexer to process a BSPL protocol.

* `test`: Test resources.

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