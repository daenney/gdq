package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/daenney/gdq/v2"
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
		fmt.Fprintln(w.tw, "Start Time\tTitle\tEstimate\tRunners\tHosts\tConsole\tCategory")
	} else if w.platform {
		fmt.Fprintln(w.tw, "Start Time\tTitle\tEstimate\tRunners\tHosts\tConsole")
	} else if w.category {
		fmt.Fprintln(w.tw, "Start Time\tTitle\tEstimate\tRunners\tHosts\tCategory")

	}
	return w
}

func (w *writer) Flush() {
	w.tw.Flush()
}

func (w *writer) Write(run *gdq.Run) {
	if !w.platform && !w.category {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\n",
			run.Start.Local().Format(time.Stamp),
			run.Title,
			run.Estimate,
			run.Runners,
			strings.Join(run.Hosts, ", "),
		)
		return
	}
	if w.platform && w.category {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			run.Start.Local().Format(time.Stamp),
			run.Title,
			run.Estimate,
			run.Runners,
			strings.Join(run.Hosts, ", "),
			run.Console,
			run.Category,
		)
		return
	}
	if w.platform {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			run.Start.Local().Format(time.Stamp),
			run.Title,
			run.Estimate,
			run.Runners,
			strings.Join(run.Hosts, ", "),
			run.Console,
		)
		return
	}
	if w.category {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			run.Start.Local().Format(time.Stamp),
			run.Title,
			run.Estimate,
			run.Runners,
			strings.Join(run.Hosts, ", "),
			run.Category,
		)
		return
	}
}
