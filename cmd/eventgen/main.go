package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/daenney/gdq/v3"
)

func main() {
	flag.CommandLine.SetOutput(os.Stdout)
	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), "Usage of %s:\n", os.Args[0])
		fmt.Fprintln(flag.CommandLine.Output(), "This generates the events.")
		flag.PrintDefaults()
	}

	flag.Parse()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	c := gdq.New(http.DefaultClient)
	evs, err := c.Events(ctx)
	if err != nil {
		panic(err)
	}
	slices.SortFunc(evs, func(a, b *gdq.Event) int {
		if a.ID < b.ID {
			return 1
		}
		if a.ID > b.ID {
			return -1
		}
		return 0
	})
	slices.Reverse(evs)

	var output strings.Builder
	writeLine("// Code generated by eventgen; DO NOT EDIT.", &output)
	writeLine("", &output)
	writeLine("package gdq", &output)
	writeLine("", &output)

	writeLine("// All the GDQ events, sorted by Event.ID", &output)
	writeLine("var (", &output)
	for _, ev := range evs {
		s := short(ev.Short)
		n := name(ev.Name)
		writeLine(fmt.Sprintf("\t%s = Event{ ID: %d, Short: \"%s\", Name: \"%s\", Year: %d}", s, ev.ID, s, n, ev.Year), &output)
	}
	writeLine(")", &output)
	writeLine("", &output)

	writeLine("var eventsByName = map[string]Event{", &output)
	for _, ev := range evs {
		s := short(ev.Short)
		writeLine(fmt.Sprintf("\t\"%s\": %s,", strings.ToLower(s), s), &output)
	}
	writeLine("}", &output)
	writeLine("", &output)

	writeLine("var eventsByID = map[uint]Event{", &output)
	for _, ev := range evs {
		s := short(ev.Short)
		writeLine(fmt.Sprintf("\t%d: %s,", ev.ID, s), &output)
	}
	writeLine("}", &output)

	fmt.Fprintln(flag.CommandLine.Output(), output.String())
}

func writeLine(str string, builder *strings.Builder) {
	builder.WriteString(str + "\n")
}

func short(s string) string {
	if s == "spook" {
		return "Spook"
	}
	if s == "thpslaunch" {
		return "THPSLaunch"
	}
	if strings.HasPrefix(s, "sgdq") || strings.HasPrefix(s, "agdq") || strings.HasSuffix(s, "q") {
		return strings.ToUpper(s)
	}
	if strings.HasPrefix(s, "frostfatales") {
		year := strings.TrimPrefix(s, "frostfatales")
		return "FrostFatales" + year
	}
	if strings.HasPrefix(s, "flamefatales") {
		year := strings.TrimPrefix(s, "flamefatales")
		return "FlameFatales" + year
	}
	if strings.HasPrefix(s, "fleetfatales") {
		year := strings.TrimPrefix(s, "fleetfatales")
		return "FleetFatales" + year
	}
	return s
}

func name(s string) string {
	// strip the year
	elems := strings.Fields(s)
	last := elems[len(elems)-1]
	_, err := strconv.Atoi(last)
	if err == nil {
		elems = elems[:len(elems)-1]
	}

	return strings.Join(elems, " ")
}
