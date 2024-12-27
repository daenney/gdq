package main

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/daenney/gdq/v3"
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
		fmt.Fprintln(w.tw, "Start Time\tTitle\tEstimate\tRunners\tHosts\tCommentators")
	} else if w.platform && w.category {
		fmt.Fprintln(w.tw, "Start Time\tTitle\tEstimate\tRunners\tHosts\tCommentators\tConsole\tCategory")
	} else if w.platform {
		fmt.Fprintln(w.tw, "Start Time\tTitle\tEstimate\tRunners\tHosts\tCommentators\tConsole")
	} else if w.category {
		fmt.Fprintln(w.tw, "Start Time\tTitle\tEstimate\tRunners\tHosts\tCommentators\tCategory")
	}
	return w
}

func (w *writer) Flush() {
	w.tw.Flush()
}

func (w *writer) Write(run *gdq.Run) {
	if !w.platform && !w.category {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\t%s\n",
			run.Start.Local().Format(time.Stamp),
			run.Title,
			run.Estimate,
			names(run.Runners),
			names(run.Hosts),
			names(run.Commentators),
		)
		return
	}
	if w.platform && w.category {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			run.Start.Local().Format(time.Stamp),
			run.Title,
			run.Estimate,
			names(run.Runners),
			names(run.Hosts),
			names(run.Commentators),
			run.Platform,
			run.Category,
		)
		return
	}
	if w.platform {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			run.Start.Local().Format(time.Stamp),
			run.Title,
			run.Estimate,
			names(run.Runners),
			names(run.Hosts),
			names(run.Commentators),
			run.Platform,
		)
		return
	}
	if w.category {
		fmt.Fprintf(w.tw, "%s\t%s\t%s\t%s\t%s\t%s\t%s\n",
			run.Start.Local().Format(time.Stamp),
			run.Title,
			run.Estimate,
			names(run.Runners),
			names(run.Hosts),
			names(run.Commentators),
			run.Category,
		)
		return
	}
}

func names(ts []gdq.Talent) string {
	res := make([]string, 0, len(ts))
	for _, t := range ts {
		res = append(res, t.Name)
	}
	return strings.Join(res, ", ")
}
