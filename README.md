# personnummer [![GitHub Workflow Status](https://img.shields.io/github/workflow/status/personnummer/go/test)](https://github.com/personnummer/go/actions) [![GoDoc](https://godoc.org/github.com/personnummer/go?status.svg)](https://godoc.org/github.com/personnummer/go) [![Go Report Card](https://goreportcard.com/badge/github.com/personnummer/go)](https://goreportcard.com/report/github.com/personnummer/go)

Validate Swedish social security numbers.

## Installation

```
go get -u github.com/personnummer/go
```

## Example

```go
package main

import (
	personnummer "github.com/personnummer/go"
)

func main() {
	personnummer.Valid("198507099805")
	//=> true

	// works with co-ordination numbers
	personnummer.Valid("198507699802")
	//=> true
}
```

## License

MIT