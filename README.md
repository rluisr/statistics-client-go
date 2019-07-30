statistics-client-go
====================

TBD

Usage
-----

```go
package main

import (
	scg "github.com/rluisr/statistics-client-go"
)

func main() {
	s := scg.Client("opst", "http://localhost:8000")
	_ = s.Register()
}
```
