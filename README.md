# GDQ [![Go Reference](https://pkg.go.dev/badge/github.com/daenney/gdq.svg)](https://pkg.go.dev/github.com/daenney/gdq)

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
