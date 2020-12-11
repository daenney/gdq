package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/daenney/gdq"
)

type writer struct {
	tw       *tabwriter.Writer
	category bool
	platform bool
}

func newWriter(category bool, platform bool) *writer {
	w := &writer{tw: tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)}
	w.category = category
	w.platform = platform
	if !w.platform && !w.category {
		fmt.Fprintln(w.tw, "Start Time\tTitle\tEstimate\tRunners\tHosts")
	} else if w.platform && w.category {
		fmt.Fprintln(w.tw, "Start Time\tTitle\tEstimate\tRunners\tHosts\tPlatform\tCategory")
	} else if w.platform {
		fmt.Fprintln(w.tw, "Start Time\tTitle\tEstimate\tRunners\tHosts\tPlatform")
	} else if w.category {
		fmt.Fprintln(w.tw, "Start Time\tTitle\tEstimate\tRunners\tHosts\tCategory")

	}
	return w
}

func (w *writer) Flush() {
	w.tw.Flush()
}

func (w *writer) Write(event *gdq.Event) {
	if !w.platform && !w.category {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\n",
			event.Start.Local().Format(time.Stamp),
			event.Title,
			event.Estimate,
			strings.Join(event.Runners, ", "),
			strings.Join(event.Hosts, ", "),
		)
		return
	}
	if w.platform && w.category {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			event.Start.Local().Format(time.Stamp),
			event.Title,
			event.Estimate,
			strings.Join(event.Runners, ", "),
			strings.Join(event.Hosts, ", "),
			event.Platform,
			event.Category,
		)
		return
	}
	if w.platform {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			event.Start.Local().Format(time.Stamp),
			event.Title,
			event.Estimate,
			strings.Join(event.Runners, ", "),
			strings.Join(event.Hosts, ", "),
			event.Platform,
		)
		return
	}
	if w.category {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			event.Start.Local().Format(time.Stamp),
			event.Title,
			event.Estimate,
			strings.Join(event.Runners, ", "),
			strings.Join(event.Hosts, ", "),
			event.Category,
		)
		return
	}
}
