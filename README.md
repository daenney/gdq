<h1 align="center">
🏃 GDQ 🎮
</h1>
<h4 align="center">A Go library and CLI for Games Done Quick</h4>
<p align="center">
    <a href="https://github.com/daenney/gdq/actions?query=workflow%3ACI"><img src="https://github.com/daenney/gdq/workflows/CI/badge.svg" alt="Build Status"></a>
    <a href="https://codecov.io/gh/daenney/gdq"><img src="https://codecov.io/gh/daenney/gdq/branch/main/graph/badge.svg" alt="Coverage Status"></a>
    <a href="https://goreportcard.com/report/github.com/daenney/gdq"><img src="https://goreportcard.com/badge/github.com/daenney/gdq" alt="Go report card"></a>
    <a href="https://pkg.go.dev/github.com/daenney/gdq"><img src="https://pkg.go.dev/badge/github.com/daenney/gdq.svg" alt="GoDoc"></a>
    <a href="LICENSE"><img src="https://img.shields.io/github/license/daenney/gdq?style=flat-square" alt="License: MIT"></a>
</p>

[Games Done Quick (GDQ)](https://gamesdonequick.com/) is a regular
speedrunning event that collects money for charity. The event is incredibly
fun, especially if you enjoy seeing your favourite games torn to shreds by
amazing runners and supported with great commentary and prizes to win.

This repo contains a Go library and CLI for working with the GDQ schedule.
Since GDQ doesn't have an API for the schedule we rely on parsing the HTML of
the schedule page instead, so beware this library might break if the design
of the schedule changes.

There is also a companion [Matrix](https://matrix.org) bot over at
[GDQBot](https://github.com/daenney/gdqbot).

## Installation

There are prebuilt binaries available for every release from v1.0.0 onwards. You
can find them [over here](https://github.com/daenney/gdq/releases).

|Platform|Architecture|Binary|
|---|---|---|
|Windows|amd64|✅|
|macOS|amd64|✅|
|macOS|arm64/M1<sup id="a1">[1](#f1)</sup>|❌|
|Linux|amd64|✅|
|Linux|arm64|✅|
|Linux|armv7/amrhf|✅|
|Linux|armv6/arm</sup>|✅|

<b id="f1"><sup>1</sup></b> Pending Go 1.16 release [↩](#a1)

## Building

You can `go get` the library, or `git clone` and then run a `go build` followed
by a `go test` to ensure everything is OK.

You can build the CLI using `go build -trimpath -o gdqctl cmd/gdqcli/*.go` or
install it directly using `go install github.com/daenney/gdq/cmd/gdqcli`. See
`go help install` for where the binaries will end up.

To embed the version, commit and date at build time you'll need to add
`-X main.version=VERSION -X main.commit=SHA -X main.date=DATE` and compute
the right values yourself.

## Contributing

PRs welcome! Fork+clone the repo and send me a patch. Please ensure that:
* Make small commits that encapsulate one functional change at a time
  (implementation change, the associated tests and any doc changes)
* Every commit explains what it's trying to achieve and why
* The tests pass
