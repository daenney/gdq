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

func main() {

	host := flag.String("host", "", "show events matching this host")
	runner := flag.String("runner", "", "show events matching this runner")
	title := flag.String("title", "", "show events matching this title")

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
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Start Time\tTitle\tEstimate\tRunners\tHosts")
		for _, event := range schedule.Events {
			fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
				event.Start.Local().Format(time.Stamp),
				event.Title,
				event.Estimate,
				strings.Join(event.Runners, ", "),
				strings.Join(event.Hosts, ", "),
			)
		}
		w.Flush()
	}
}
