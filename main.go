package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/daenney/gdqbot/gdq"
)

func main() {
	schedule, err := gdq.GetSchedule(gdq.AGDQ2021, nil)
	if err != nil {
		log.Fatalln(err)
	}

	if schedule != nil && len(schedule.Events) > 0 {
		w := tabwriter.NewWriter(os.Stdout, 0, 0, 1, ' ', tabwriter.Debug)
		fmt.Fprintln(w, "Start Time\tTitle\tEstimate\tRunners\tHosts")
		schedule.RLock()
		defer schedule.RUnlock()
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
