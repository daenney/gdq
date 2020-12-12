<h1 align="center">
🏃 GDQ 🎮
</h1>
<h4 align="center">A Go library and CLI for Games Done Quick</h4>
<p align="center">
    <a href="https://github.com/daenney/gdq/actions?query=workflow%3ACI"><img src="https://github.com/daenney/gdq/workflows/CI/badge.svg" alt="Build Status"></a>
    <a href="https://codecov.io/gh/daenney/gdq"><img src="https://codecov.io/gh/daenney/gdq/branch/main/graph/badge.svg" alt="Coverage Status"></a>
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
