package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/daenney/gdq"
)

func main() {
	host := flag.String("host", "", "show events matching this host")
	runner := flag.String("runner", "", "show events matching this runner")
	title := flag.String("title", "", "show events matching this title")
	category := flag.Bool("show-category", false, "show category in the output")
	platform := flag.Bool("show-platform", false, "show platform in the output")
	format := flag.String("format", "table", "one of table or json")
	edition := flag.String("edition", "", "GDQ edition to query. This can be a string or a schedule number and when omitted will result in the current/upcoming schedule being used")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(), "When using filters, each filter is applied and the resulting filtered schedule is then filtered with the next filter. This means filters are additive, so you can't say show me events for this host or this runner.")
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(), "All filters use a case insensitive substring match. This means that passing a filter of '-runner e' will find all events where any runner has the letter 'e' in their handle.")
	}

	flag.Parse()

	var ed gdq.Edition
	if *edition == "" {
		ed = gdq.Latest
	} else {
		v, ok := gdq.GetEdition(*edition)
		if !ok {
			num, err := strconv.ParseUint(*edition, 10, 64)
			if err != nil {
				log.Fatalf("Could not find an edition matching: %s\n", *edition)
			}
			ed = gdq.Edition(uint(num))
		} else {
			ed = v
		}
	}

	schedule, err := gdq.GetSchedule(ed, newHTTPClient())
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
		switch strings.ToLower(*format) {
		case "table":
			w := newWriter(*category, *platform)
			for _, event := range schedule.Events {
				w.Write(event)
			}
			w.Flush()
		case "json":
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "    ")
			if err := enc.Encode(schedule.Events); err != nil {
				log.Fatalln(err)
			}
		default:
			log.Fatalf("unrecognised value for format flag: %s\n", *format)
		}
	}
}
