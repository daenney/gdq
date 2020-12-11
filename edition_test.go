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
		}{
			{edition: AGDQ2016, name: "AGDQ2016"},
			{edition: SGDQ2016, name: "SGDQ2016"},
			{edition: AGDQ2017, name: "AGDQ2017"},
			{edition: SGDQ2017, name: "SGDQ2017"},
			{edition: HRDQ2017, name: "HRDQ2017"},
			{edition: AGDQ2018, name: "AGDQ2018"},
			{edition: SGDQ2018, name: "SGDQ2018"},
			{edition: GDQX2018, name: "GDQX2018"},
			{edition: AGDQ2019, name: "AGDQ2019"},
			{edition: SGDQ2019, name: "SGDQ2019"},
			{edition: GDQX2019, name: "GDQX2019"},
			{edition: AGDQ2020, name: "AGDQ2020"},
			{edition: FrostFatales2020, name: "FrostFatales2020"},
			{edition: SGDQ2020, name: "SGDQ2020"},
			{edition: CRDQ2020, name: "CRDQ2020"},
			{edition: FleetFatales2020, name: "FleetFatales2020"},
			{edition: AGDQ2021, name: "AGDQ2021"},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				assertEqual(t, tt.edition.String(), tt.name)
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
