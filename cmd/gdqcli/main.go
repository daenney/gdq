package main

import (
	"flag"
	"fmt"
	"log"
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

func main() {

	host := flag.String("host", "", "show events matching this host")
	runner := flag.String("runner", "", "show events matching this runner")
	title := flag.String("title", "", "show events matching this title")
	category := flag.Bool("show-category", false, "show category in the output")
	platform := flag.Bool("show-platform", false, "show platform in the output")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprint(flag.CommandLine.Output(), "\n")
		fmt.Fprint(flag.CommandLine.Output(), "When using filters, each filter is applied and the resulting filtered schedule is then filtered with the next filter. This means filters are additive, so you can't say show me events for this host or this runner.\n")
	}

	flag.Parse()

	schedule, err := gdq.GetSchedule(gdq.AGDQ2021, nil)
	if err != nil {
		log.Fatalln(err)
	}

	if *runner != "" {
		schedule = schedule.ForRunner(*runner)
	}
	if *host != "" {
		schedule = schedule.ForHost(*host)
	}
	if *title != "" {
		schedule = schedule.ForTitle(*title)
	}

	if schedule != nil && len(schedule.Events) > 0 {
		w := newWriter(*category, *platform)
		for _, event := range schedule.Events {
			w.Write(event)
		}
		w.Flush()
	}
}
