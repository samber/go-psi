
# Pressure stall informations (PSI) and starvation notifier

[![tag](https://img.shields.io/github/tag/samber/go-psi.svg)](https://github.com/samber/go-psi/releases)
![Go Version](https://img.shields.io/badge/Go-%3E%3D%201.18.0-%23007d9c)
[![GoDoc](https://godoc.org/github.com/samber/go-psi?status.svg)](https://pkg.go.dev/github.com/samber/go-psi)
![Build Status](https://github.com/samber/go-psi/actions/workflows/test.yml/badge.svg)
[![Go report](https://goreportcard.com/badge/github.com/samber/go-psi)](https://goreportcard.com/report/github.com/samber/go-psi)
[![Coverage](https://img.shields.io/codecov/c/github/samber/go-psi)](https://codecov.io/gh/samber/go-psi)
[![Contributors](https://img.shields.io/github/contributors/samber/go-psi)](https://github.com/samber/go-psi/graphs/contributors)
[![License](https://img.shields.io/github/license/samber/go-psi)](./LICENSE)

A few readings to getting started with PSI:
- https://docs.kernel.org/accounting/psi.html
- https://facebookmicrosites.github.io/psi/docs/overview
- https://unixism.net/2019/08/linux-pressure-stall-information-psi-by-example/

## 🚀 Install

```sh
go get github.com/samber/go-psi
```

This library is v1 and follows SemVer strictly. No breaking changes will be made to exported APIs before v2.0.0.

Requires Linux kernel >= 4.20.

## 💡 Usage

GoDoc: [https://pkg.go.dev/github.com/samber/go-psi](https://pkg.go.dev/github.com/samber/go-psi)

### Retrieve current PSI state

```go
import "github.com/samber/go-psi"

// Get PSI for a single resource: psi.Memory or psi.CPU or psi.IO.
stats, err := psi.PSIStatsForResource(psi.Memory)

// Get all PSI stats.
all, err := psi.AllPSIStats()
```

### Get PSI change notifications

```go
import "github.com/samber/go-psi"

onChange, done, err := psi.Notify(psi.Memory)

for {
    last, _ := <-onChange
    fmt.Printf("\nMemory:\n%s\n", last)
}

// when you're done, just stop the notifier
<-done
```

### Get PSI starvation alerts

```go
import "github.com/samber/go-psi"

onAlert, done, err := psi.NotifyStarvation(psi.CPU, psi.Avg10, 3, 4)
for {
    alert, _ := <-onAlert
    fmt.Printf("\nALERT %t\nCPU: %f%%\n", alert.Starved, alert.Current)
}

// when you're done, just stop the notifier
<-done
```

## 🤝 Contributing

- Ping me on twitter [@samuelberthe](https://twitter.com/samuelberthe) (DMs, mentions, whatever :))
- Fork the [project](https://github.com/samber/go-psi)
- Fix [open issues](https://github.com/samber/go-psi/issues) or request new features

Don't hesitate ;)

```bash
# Install some dev dependencies
make tools

# Run tests
make test
# or
make watch-test
```

## 👤 Contributors

![Contributors](https://contrib.rocks/image?repo=samber/go-psi)

## 💫 Show your support

Give a ⭐️ if this project helped you!

[![GitHub Sponsors](https://img.shields.io/github/sponsors/samber?style=for-the-badge)](https://github.com/sponsors/samber)

## 📝 License

Copyright © 2023 [Samuel Berthe](https://github.com/samber).

This project is [MIT](./LICENSE) licensed.
