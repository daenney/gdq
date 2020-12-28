package gdq

import (
	"testing"
)

func TestEvent(t *testing.T) {
	t.Run("latest", func(t *testing.T) {
		assertEqual(t, Event(0).String(), "latest")
	})
	t.Run("unknown", func(t *testing.T) {
		assertEqual(t, Event(^uint(0)).String(), "unknown event: 18446744073709551615")
	})
	t.Run("named events", func(t *testing.T) {
		var tests = []struct {
			event Event
			name  string
			id    uint
		}{
			{event: AGDQ2016, name: "AGDQ2016", id: 17},
			{event: SGDQ2016, name: "SGDQ2016", id: 18},
			{event: AGDQ2017, name: "AGDQ2017", id: 19},
			{event: SGDQ2017, name: "SGDQ2017", id: 20},
			{event: HRDQ2017, name: "HRDQ2017", id: 21},
			{event: AGDQ2018, name: "AGDQ2018", id: 22},
			{event: SGDQ2018, name: "SGDQ2018", id: 23},
			{event: GDQX2018, name: "GDQX2018", id: 24},
			{event: AGDQ2019, name: "AGDQ2019", id: 25},
			{event: SGDQ2019, name: "SGDQ2019", id: 26},
			{event: GDQX2019, name: "GDQX2019", id: 27},
			{event: AGDQ2020, name: "AGDQ2020", id: 28},
			{event: FrostFatales2020, name: "FrostFatales2020", id: 29},
			{event: SGDQ2020, name: "SGDQ2020", id: 30},
			{event: CRDQ2020, name: "CRDQ2020", id: 31},
			{event: FleetFatales2020, name: "FleetFatales2020", id: 33},
			{event: AGDQ2021, name: "AGDQ2021", id: 34},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assertEqual(t, tt.event.String(), tt.name)
				assertEqual(t, uint(tt.event), tt.id)
			})
		}
	})
}

func TestGetEvent(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		_, ok := GetEvent("")
		assertEqual(t, ok, false)
	})
	t.Run("unknown", func(t *testing.T) {
		_, ok := GetEvent("AGDQ0")
		assertEqual(t, ok, false)
	})
	t.Run("known", func(t *testing.T) {
		e, ok := GetEvent("AgdQ2016")
		assertEqual(t, ok, true)
		assertEqual(t, e, AGDQ2016)
	})
}
