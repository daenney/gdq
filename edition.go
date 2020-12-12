package gdq

import (
	"fmt"
	"strings"
)

// Edition is the schedule ID of a GDQ edition
type Edition uint

// All the GDQ editions for which a schedule is available
const (
	Latest   Edition = 0
	AGDQ2016 Edition = iota + 16
	SGDQ2016
	AGDQ2017
	SGDQ2017
	HRDQ2017
	AGDQ2018
	SGDQ2018
	GDQX2018
	AGDQ2019
	SGDQ2019
	GDQX2019
	AGDQ2020
	FrostFatales2020
	SGDQ2020
	CRDQ2020
	_
	FleetFatales2020
	AGDQ2021
)

//gocyclo:ignore
func (e Edition) String() string {
	switch e {
	case Latest:
		return "latest"
	case AGDQ2016:
		return "AGDQ2016"
	case SGDQ2016:
		return "SGDQ2016"
	case AGDQ2017:
		return "AGDQ2017"
	case SGDQ2017:
		return "SGDQ2017"
	case HRDQ2017:
		return "HRDQ2017"
	case AGDQ2018:
		return "AGDQ2018"
	case SGDQ2018:
		return "SGDQ2018"
	case GDQX2018:
		return "GDQX2018"
	case AGDQ2019:
		return "AGDQ2019"
	case SGDQ2019:
		return "SGDQ2019"
	case GDQX2019:
		return "GDQX2019"
	case AGDQ2020:
		return "AGDQ2020"
	case FrostFatales2020:
		return "FrostFatales2020"
	case SGDQ2020:
		return "SGDQ2020"
	case CRDQ2020:
		return "CRDQ2020"
	case FleetFatales2020:
		return "FleetFatales2020"
	case AGDQ2021:
		return "AGDQ2021"
	default:
		return fmt.Sprintf("unknown edition: %d", e)
	}
}

var editions = map[string]Edition{
	"latest":           Latest,
	"agdq2016":         AGDQ2016,
	"sgdq2016":         SGDQ2016,
	"agdq2017":         AGDQ2017,
	"sgdq2017":         SGDQ2017,
	"hrdq2017":         HRDQ2017,
	"agdq2018":         AGDQ2018,
	"sgdq2018":         SGDQ2018,
	"gdqx2018":         GDQX2018,
	"agdq2019":         AGDQ2019,
	"sgdq2019":         SGDQ2019,
	"gdqx2019":         GDQX2019,
	"agdq2020":         AGDQ2020,
	"frostfatales2020": FrostFatales2020,
	"sgdq2020":         SGDQ2020,
	"crdq2020":         CRDQ2020,
	"fleetfatales2020": FleetFatales2020,
	"agdq2021":         AGDQ2021,
}

// GetEdition tries to find an edition matching the input
func GetEdition(input string) (edition Edition, found bool) {
	edition, ok := editions[strings.ToLower(input)]
	return edition, ok
}
