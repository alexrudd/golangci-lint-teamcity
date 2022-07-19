# Golangci-lint TeamCity

This cli tool takes golangci-lint's json output and converts it to TeamCity service messages.

Usage:

```sh
go install github.com/alexrudd/golangci-lint-teamcity@latest
golangci-lint run --out-format=json | golangci-lint-teamcity
```
