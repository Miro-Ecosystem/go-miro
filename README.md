# Miro API in Go

![tests](https://github.com/Miro-Ecosystem/go-miro/workflows/tests/badge.svg)
[![codecov](https://codecov.io/gh/Miro-Ecosystem/go-miro/branch/master/graph/badge.svg)](https://codecov.io/gh/Miro-Ecosystem/go-miro)

Go written [Miro](https://miro.com/app/dashboard/) API client.

*Note that Widget is not supported yet*

## Installation

Include this is your code as below:

```go
import "github.com/Miro-Ecosystem/go-miro/miro"
```

Using `go get`:

```console
$ go get github.com/Miro-Ecosystem/go-miro
```

## Usage

Using the client:

```go
client := miro.NewClient("access token")
```

API's are very simple and easy to understand.

```go
client.Board.Get("10")
```

## Supported APIs

See the [Supported APIs](docs/API.md) for supported APIs.

## Copyright and License

Please see the LICENSE file for the included license information.
Copyright 2020 by Keisuke Yamashita.
