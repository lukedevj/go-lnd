# Go-lnd

Just a simple library that calls LND API and returns a Json.

[![Donate](https://img.shields.io/badge/Donate-Bitcoin-green.svg)](https://coinos.io/lukedevj)


```go
package main

import (
  "github.com/lukedevj/go-lnd"
  "fmt"
)

func main() {
    c := lnd.Client{
      Host:     "127.0.0.1:8080",
      Cert:     "/home/dev/.lnd/tls.cert",
      Macaroon: "/home/dev/.lnd/chain/bitcoin/regtest/invoices.macaroon",
    }
    d := map[string]interface{}{"memo": "Testing", "value": 20}
    // Create new invoice
    i, _ := c.Call("POST", "v1/invoices", d)
    fmt.Println(i)
}
```
