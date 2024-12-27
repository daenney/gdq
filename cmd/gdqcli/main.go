package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/daenney/gdq/v3"
)

func main() {
	host := flag.String("host", "", "show runs matching this host")
	runner := flag.String("runner", "", "show runs matching this runner")
	title := flag.String("title", "", "show runs matching this title")
	category := flag.Bool("show-category", false, "show category in the output")
	platform := flag.Bool("show-platform", false, "show platform in the output")
	format := flag.String("format", "table", "one of table or json")
	event := flag.String("event", "", "GDQ event to query. This can be a string or a event number and when omitted will result in the current/upcoming schedule being used")
	showVersion := flag.Bool("version", false, "show CLI version and build info")
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(), "When using filters, each filter is applied and the resulting filtered schedule is then filtered with the next filter. This means filters are additive, so you can't say show me runs for this host or this runner.")
		fmt.Fprintln(flag.CommandLine.Output())
		fmt.Fprintln(flag.CommandLine.Output(), "All filters use a case insensitive substring match. This means that passing a filter of '-runner e' will find all runs where any runner has the letter 'e' in their handle.")
	}

	flag.Parse()

	if *showVersion {
		fmt.Fprintf(os.Stdout, "{\"version\": \"%s\", \"commit\": \"%s\", \"date\": \"%s\"}\n", version, commit, date)
		os.Exit(0)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	g := gdq.New(newHTTPClient())

	var ev *gdq.Event
	if *event == "" {
		v, err := g.Events(ctx)
		if err != nil {
			log.Fatalln(err)
		}
		ev = v[0]
	} else {
		v, ok := gdq.GetEventByName(*event)
		if !ok {
			num, err := strconv.ParseUint(*event, 10, 64)
			if err != nil {
				log.Fatalf("Could not find an event matching: %s\n", *event)
			}
			v, ok = gdq.GetEventByID(uint(num))
			if !ok {
				ev = &gdq.Event{ID: uint(num), Short: "unknown", Name: "unknown", Year: 0}
			} else {
				ev = v
			}
		} else {
			ev = v
		}
	}

	schedule, err := g.Schedule(ctx, ev.ID)
	if err != nil {
		log.Fatalln(err)
	}

	if len(schedule.Runs) == 0 {
		log.Printf("No runs for event with ID %d: (%s)\n", ev.ID, ev.String())
		os.Exit(0)
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

	if schedule != nil && len(schedule.Runs) > 0 {
		switch strings.ToLower(*format) {
		case "table":
			w := newWriter(*category, *platform)
			for _, run := range schedule.Runs {
				w.Write(run)
			}
			w.Flush()
		case "json":
			enc := json.NewEncoder(os.Stdout)
			enc.SetIndent("", "    ")
			if err := enc.Encode(schedule.Runs); err != nil {
				log.Fatalln(err)
			}
		default:
			log.Fatalf("unrecognised value for format flag: %s\n", *format)
		}
	}
}
