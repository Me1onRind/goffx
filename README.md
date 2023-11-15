# goffx
goffx is a pure go implementation of Format-preserving, Feistel-based encryption (FFX) from https://github.com/emulbreh/pyffx.

Only support sha1 at present. It's no difficult to add other encryption algorithms.


## Usage
```go
package main

import (
    "fmt"

    "github.com/Me1onRind/goffx"
)

func main() {
    ie := goffx.Integer("secret-key", 4)
    fmt.Println(ie.Encrypt(1234)) // output 6103 <nil>
    fmt.Println(ie.Decrypt(6103)) // output 1234 <nil>

    se := goffx.String("secret-key", "abc", 6)
    fmt.Println(se.Encrypt("aaabbb")) // output "acbacc" <nil>
    fmt.Println(se.Decrypt("acbacc")) // output "aaabbb" <nil>
}
```
