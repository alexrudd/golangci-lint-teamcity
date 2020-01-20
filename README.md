# Golangci-lint TeamCity

This cli tool takes golangci-lint's json output and converts it to TeamCity service messages.

Usage:

```sh
go get github.com/alexrudd/golangci-lint-teamcity
golangci-lint run --out-format=json | golangci-lint-teamcity
```
