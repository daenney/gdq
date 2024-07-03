package gdq

import (
	"fmt"
	"strings"
	"time"
)

type eventResp struct {
	PK     uint `json:"pk"`
	Fields struct {
		Short          string    `json:"short"`
		Name           string    `json:"name"`
		Date           time.Time `json:"datetime"`
		DonationsTotal float64   `json:"amount"`
		DonationsCount uint      `json:"count"`
		DonationsMax   float64   `json:"max"`
		DonationsAvg   float64   `json:"avg"`
	} `json:"fields"`
}

func (e eventResp) toEvent() Event {
	return Event{
		ID:    e.PK,
		Short: e.Fields.Short,
		Name:  e.Fields.Name,
		Year:  e.Fields.Date.Year(),
		Donations: Donations{
			Total:   e.Fields.DonationsTotal,
			Max:     e.Fields.DonationsMax,
			Count:   e.Fields.DonationsCount,
			Average: e.Fields.DonationsAvg,
		},
	}
}

type eventsResp []eventResp

// Event is the schedule ID of a GDQ event
type Event struct {
	ID        uint      `json:"id"`
	Short     string    `json:"short"`
	Name      string    `json:"name"`
	Year      int       `json:"year"`
	Donations Donations `json:"donations"`
}

// Donations is the donation summary of a GDQ event
type Donations struct {
	Total   float64 `json:"total"`
	Max     float64 `json:"max"`
	Count   uint    `json:"count"`
	Average float64 `json:"average"`
}

const (
	agdq = "Awesome Games Done Quick"
	sgdq = "Summer Games Done Quick"
)

// All the GDQ events, sorted by Event.ID
var (
	SGDQ2012          = Event{ID: 1, Short: "SGDQ2012", Name: sgdq, Year: 2012}
	AGDQ2012          = Event{ID: 2, Short: "SGDQ2012", Name: agdq, Year: 2012}
	SGDQ2011          = Event{ID: 3, Short: "SGDQ2011", Name: sgdq, Year: 2011}
	JRDQ              = Event{ID: 4, Short: "JRDQ", Name: "Japan Relief Done Quick", Year: 2011}
	AGDQ2011          = Event{ID: 5, Short: "AGDQ2011", Name: agdq, Year: 2011}
	CGDQ              = Event{ID: 6, Short: "CGDQ", Name: "Classic Games Done Quick", Year: 2010}
	AGDQ2013          = Event{ID: 7, Short: "AGDQ2013", Name: agdq, Year: 2013}
	SGDQ2013          = Event{ID: 8, Short: "SGDQ2013", Name: sgdq, Year: 2013}
	AGDQ2014          = Event{ID: 9, Short: "AGDQ2014", Name: agdq, Year: 2014}
	SGDQ2014          = Event{ID: 10, Short: "SGDQ2014", Name: sgdq, Year: 2014}
	AGDQ2015          = Event{ID: 12, Short: "AGDQ2015", Name: agdq, Year: 2015}
	SPOOK             = Event{ID: 13, Short: "SPOOK", Name: "Speedrun Spooktacular", Year: 2012}
	SGDQ2015          = Event{ID: 16, Short: "SGDQ2015", Name: sgdq, Year: 2015}
	AGDQ2016          = Event{ID: 17, Short: "AGDQ2016", Name: agdq, Year: 2016}
	SGDQ2016          = Event{ID: 18, Short: "SGDQ2016", Name: sgdq, Year: 2016}
	AGDQ2017          = Event{ID: 19, Short: "AGDQ2017", Name: agdq, Year: 2017}
	SGDQ2017          = Event{ID: 20, Short: "SGDQ2017", Name: sgdq, Year: 2017}
	HRDQ              = Event{ID: 21, Short: "HRDQ", Name: "Harvey Relief Done Quick", Year: 2017}
	AGDQ2018          = Event{ID: 22, Short: "AGDQ2018", Name: agdq, Year: 2018}
	SGDQ2018          = Event{ID: 23, Short: "SGDQ2018", Name: sgdq, Year: 2018}
	GDQX2018          = Event{ID: 24, Short: "GDQX2018", Name: "Games Done Quick Express", Year: 2018}
	AGDQ2019          = Event{ID: 25, Short: "AGDQ2019", Name: agdq, Year: 2019}
	SGDQ2019          = Event{ID: 26, Short: "SGDQ2019", Name: sgdq, Year: 2019}
	GDQX2019          = Event{ID: 27, Short: "GDQX2019", Name: "Games Done Quick Express", Year: 2019}
	AGDQ2020          = Event{ID: 28, Short: "AGDQ2020", Name: agdq, Year: 2020}
	FrostFatales2020  = Event{ID: 29, Short: "FrostFatales2020", Name: "Frost Fatales", Year: 2020}
	SGDQ2020          = Event{ID: 30, Short: "SGDQ2020", Name: sgdq, Year: 2020}
	CRDQ              = Event{ID: 31, Short: "CRDQ", Name: "Corona Relief Done Quick", Year: 2020}
	THPSLaunch        = Event{ID: 32, Short: "THPSLaunch", Name: "Tony Hawk's Pro Skater 1 + 2 Launch Celebration", Year: 2020}
	FleetFatales2020  = Event{ID: 33, Short: "FleetFatales2020", Name: "Fleet Fatales", Year: 2020}
	AGDQ2021          = Event{ID: 34, Short: "AGDQ2021", Name: agdq + " Online", Year: 2021}
	SGDQ2021          = Event{ID: 35, Short: "SGDQ2021", Name: sgdq + " Online", Year: 2021}
	FlamesFatales2021 = Event{ID: 36, Short: "FlamesFatales2021", Name: "Flames Fatales", Year: 2021}
	AGDQ2022          = Event{ID: 37, Short: "AGDQ2022", Name: agdq + " Online", Year: 2022}
	FrostFatales2022  = Event{ID: 38, Short: "FrostFatales2022", Name: "Frost Fatales", Year: 2022}
	SGDQ2022          = Event{ID: 39, Short: "SGDQ2022", Name: sgdq, Year: 2022}
	FlamesFatales2022 = Event{ID: 40, Short: "FlamesFatales2022", Name: "Flames Fatales", Year: 2022}
	AGDQ2023          = Event{ID: 41, Short: "AGDQ2023", Name: agdq, Year: 2023}
	SGDQ2024          = Event{ID: 48, Short: "SGDQ2024", Name: sgdq, Year: 2024}
)

func (e Event) String() string {
	return fmt.Sprintf("%s (%d)", e.Name, e.Year)
}

var eventsByName = map[string]Event{
	"classic":           CGDQ,
	"cgdq":              CGDQ,
	"cgdq2010":          CGDQ,
	"agdq2011":          AGDQ2011,
	"japan":             JRDQ,
	"jrdq":              JRDQ,
	"jrdq2011":          JRDQ,
	"sgdq2011":          SGDQ2011,
	"agdq2012":          AGDQ2012,
	"sgdq2012":          SGDQ2012,
	"spook":             SPOOK,
	"spook2012":         SPOOK,
	"agdq2013":          AGDQ2013,
	"sgdq2013":          SGDQ2013,
	"agdq2014":          AGDQ2014,
	"sgdq2014":          SGDQ2014,
	"agdq2015":          AGDQ2015,
	"sgdq2015":          SGDQ2015,
	"agdq2016":          AGDQ2016,
	"sgdq2016":          SGDQ2016,
	"agdq2017":          AGDQ2017,
	"sgdq2017":          SGDQ2017,
	"harvey":            HRDQ,
	"hrdq":              HRDQ,
	"hrdq2017":          HRDQ,
	"agdq2018":          AGDQ2018,
	"sgdq2018":          SGDQ2018,
	"gdqx2018":          GDQX2018,
	"agdq2019":          AGDQ2019,
	"sgdq2019":          SGDQ2019,
	"gdqx2019":          GDQX2019,
	"agdq2020":          AGDQ2020,
	"frostfatales2020":  FrostFatales2020,
	"sgdq2020":          SGDQ2020,
	"corona":            CRDQ,
	"crdq":              CRDQ,
	"crdq2020":          CRDQ,
	"thps":              THPSLaunch,
	"thpslaunch":        THPSLaunch,
	"fleetfatales2020":  FleetFatales2020,
	"agdq2021":          AGDQ2021,
	"sgdq2021":          SGDQ2021,
	"flamesfatales2021": FlamesFatales2021,
	"agdq2022":          AGDQ2022,
	"frostfatales2022":  FrostFatales2022,
	"sgdq2022":          SGDQ2022,
	"flamefatales2022":  FlamesFatales2022,
	"agdq2023":          AGDQ2023,
}

// GetEventByName tries to find an event matching the input
func GetEventByName(input string) (ev *Event, found bool) {
	e, ok := eventsByName[strings.ToLower(input)]
	return &e, ok
}

var eventsByID = map[uint]Event{
	1:  SGDQ2012,
	2:  AGDQ2012,
	3:  SGDQ2011,
	4:  JRDQ,
	5:  AGDQ2011,
	6:  CGDQ,
	7:  AGDQ2013,
	8:  SGDQ2013,
	9:  AGDQ2014,
	10: SGDQ2014,
	12: AGDQ2015,
	13: SPOOK,
	16: SGDQ2015,
	17: AGDQ2016,
	18: SGDQ2016,
	19: AGDQ2017,
	20: SGDQ2017,
	21: HRDQ,
	22: AGDQ2018,
	23: SGDQ2018,
	24: GDQX2018,
	25: AGDQ2019,
	26: SGDQ2019,
	27: GDQX2019,
	28: AGDQ2020,
	29: FrostFatales2020,
	30: SGDQ2020,
	31: CRDQ,
	32: THPSLaunch,
	33: FleetFatales2020,
	34: AGDQ2021,
	35: SGDQ2021,
	36: FlamesFatales2021,
	37: AGDQ2022,
	38: FrostFatales2022,
	39: SGDQ2022,
	40: FlamesFatales2022,
	41: AGDQ2023,
}

// GetEventByID fetches the event by ID
func GetEventByID(id uint) (ev *Event, found bool) {
	e, ok := eventsByID[id]
	return &e, ok
}
