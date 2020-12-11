package gdq

import (
	"testing"
)

func TestEdition(t *testing.T) {
	t.Run("latest", func(t *testing.T) {
		assertEqual(t, Edition(0).String(), "latest")
	})
	t.Run("unknown", func(t *testing.T) {
		assertEqual(t, Edition(^uint(0)).String(), "unknown edition: 18446744073709551615")
	})
	t.Run("named editions", func(t *testing.T) {
		var tests = []struct {
			edition Edition
			name    string
			id      uint
		}{
			{edition: AGDQ2016, name: "AGDQ2016", id: 17},
			{edition: SGDQ2016, name: "SGDQ2016", id: 18},
			{edition: AGDQ2017, name: "AGDQ2017", id: 19},
			{edition: SGDQ2017, name: "SGDQ2017", id: 20},
			{edition: HRDQ2017, name: "HRDQ2017", id: 21},
			{edition: AGDQ2018, name: "AGDQ2018", id: 22},
			{edition: SGDQ2018, name: "SGDQ2018", id: 23},
			{edition: GDQX2018, name: "GDQX2018", id: 24},
			{edition: AGDQ2019, name: "AGDQ2019", id: 25},
			{edition: SGDQ2019, name: "SGDQ2019", id: 26},
			{edition: GDQX2019, name: "GDQX2019", id: 27},
			{edition: AGDQ2020, name: "AGDQ2020", id: 28},
			{edition: FrostFatales2020, name: "FrostFatales2020", id: 29},
			{edition: SGDQ2020, name: "SGDQ2020", id: 30},
			{edition: CRDQ2020, name: "CRDQ2020", id: 31},
			{edition: FleetFatales2020, name: "FleetFatales2020", id: 33},
			{edition: AGDQ2021, name: "AGDQ2021", id: 34},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assertEqual(t, tt.edition.String(), tt.name)
				assertEqual(t, uint(tt.edition), tt.id)
			})
		}
	})
}

func TestGetEdition(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		_, ok := GetEdition("")
		assertEqual(t, ok, false)
	})
	t.Run("unknown", func(t *testing.T) {
		_, ok := GetEdition("AGDQ0")
		assertEqual(t, ok, false)
	})
	t.Run("known", func(t *testing.T) {
		e, ok := GetEdition("AgdQ2016")
		assertEqual(t, ok, true)
		assertEqual(t, e, AGDQ2016)
	})
}
