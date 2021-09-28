# Go-lnd

Just a simple library that calls LND API and returns a Json.

[![Donate](https://img.shields.io/badge/Donate-Bitcoin-green.svg)](https://coinos.io/lukedevj)
[![Donate](https://shields.io/badge/Package--green?logo=go)](https://pkg.go.dev/github.com/lukedevj/go-lnd)


```go
package main

import (
  "github.com/lukedevj/go-lnd"
  "fmt"
)

func main() {
    client := lnd.Client{
      Host:     "127.0.0.1:8080",
      Cert:     "/home/dev/.lnd/tls.cert",
      Macaroon: "/home/dev/.lnd/chain/bitcoin/mainnet/invoices.macaroon",
    }
    invoice, _ := client.CreateInvoice(1, "Hello, word")
    fmt.Println(invoice)
    
}
```
