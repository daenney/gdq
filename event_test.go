package gdq

import (
	"errors"
	"testing"
	"time"

	"github.com/anaskhan96/soup"
)

func TestDurationMarshal(t *testing.T) {
	d := duration{0}
	j, err := d.MarshalJSON()
	assertEqual(t, err, nil)
	assertEqual(t, string(j), "0")

	dn := duration{1*time.Hour + 1*time.Minute + 1*time.Second}
	jn, err := dn.MarshalJSON()
	assertEqual(t, err, nil)
	assertEqual(t, string(jn), "3661000000000")
}

func TestDurationString(t *testing.T) {
	t.Run("no hours, no minutes", func(t *testing.T) {
		assertEqual(t, duration{0}.String(), "0 minutes")
	})
	t.Run("one hour, no minutes", func(t *testing.T) {
		assertEqual(t, duration{1 * time.Hour}.String(), "1 hour")
	})
	t.Run("one hour, one minute", func(t *testing.T) {
		assertEqual(t, duration{1*time.Hour + 1*time.Minute}.String(), "1 hour and 1 minute")
	})
	t.Run("one hour, two minutes", func(t *testing.T) {
		assertEqual(t, duration{1*time.Hour + 2*time.Minute}.String(), "1 hour and 2 minutes")
	})
	t.Run("one minute", func(t *testing.T) {
		assertEqual(t, duration{1 * time.Minute}.String(), "1 minute")
	})
	t.Run(" two minutes", func(t *testing.T) {
		assertEqual(t, duration{2 * time.Minute}.String(), "2 minutes")
	})
}

func TestNicksToSlice(t *testing.T) {
	t.Run("empty nick", func(t *testing.T) {
		assertEqual(t, len(nicksToSlice(" ")), 0)
	})
	t.Run("single nick", func(t *testing.T) {
		assertEqual(t, nicksToSlice("a ")[0], "a")
	})
	t.Run("multiple nicks", func(t *testing.T) {
		n := nicksToSlice("a,b,c")
		assertEqual(t, len(n), 3)
	})
}

func TestToDateTime(t *testing.T) {
	t.Run("empty date time", func(t *testing.T) {
		assertEqual(t, toDateTime(""), time.Time{})
	})
	t.Run("nonesense date time", func(t *testing.T) {
		assertEqual(t, toDateTime("2001:0db8:85a3:0000:0000:8a2e:0370:7334Z"), time.Time{})
	})
}

func TestToDuration(t *testing.T) {
	t.Run("empty", func(t *testing.T) {
		assertEqual(t, toDuration("").Duration, time.Duration(0))
	})
	t.Run("H:m", func(t *testing.T) {
		assertEqual(t, toDuration("1:2").Duration, time.Duration(0))
	})
	t.Run("a:b:c", func(t *testing.T) {
		assertEqual(t, toDuration("a:b:c").Duration, time.Duration(0))
	})
	t.Run("H:m:s", func(t *testing.T) {
		assertEqual(t, toDuration("01:2:30").Duration, time.Duration(1*time.Hour+2*time.Minute+30*time.Second))
	})
}

func TestEventFromHTML(t *testing.T) {
	t.Run("missing data", func(t *testing.T) {
		s := soup.HTMLParse("")
		_, err := eventFromHTML(s, s)
		assertNotNil(t, err)
		if !errors.Is(err, ErrUnexpectedData) {
			t.Errorf("expected an error of %s, got %s", ErrUnexpectedData, err)
		}
	})
	t.Run("missing columns", func(t *testing.T) {
		s := soup.HTMLParse(`<html><body><table><tr>
		<td></td>
		<td></td>
		<td> <i class="fa fa-clock-o text-gdq-red" aria-hidden="true"></i> </td>
		</tr>
		<tr>
		<td> <i class="fa fa-clock-o" aria-hidden="true"></i>  </td>
		<td><i class="fa fa-microphone"></i> </td>
		</tr></table></body></html>`)
		rows := s.FindAll("tr")
		_, err := eventFromHTML(rows[0], rows[1])
		assertNotNil(t, err)
		if !errors.Is(err, ErrUnexpectedData) {
			t.Errorf("Got %v, expected %v", err, ErrUnexpectedData)
		}
	})
	t.Run("extra columns", func(t *testing.T) {
		s := soup.HTMLParse(`<html><body><table><tr>
		<td></td>
		<td></td>
		<td> <i class="fa fa-clock-o text-gdq-red" aria-hidden="true"></i> </td>
		</tr>
		<tr>
		<td> <i class="fa fa-clock-o" aria-hidden="true"></i>  </td>
		<td><i class="fa fa-microphone"></i> </td>
		<td><i class="fa fa-microphone"></i> </td>
		</tr></table></body></html>`)
		rows := s.FindAll("tr")
		_, err := eventFromHTML(rows[0], rows[1])
		assertNotNil(t, err)
		if !errors.Is(err, ErrUnexpectedData) {
			t.Errorf("Got %v, expected %v", err, ErrUnexpectedData)
		}
	})
	t.Run("single broken event", func(t *testing.T) {
		s := soup.HTMLParse(`<html><table id="runTable"><tbody><tr>
				<td></td>
				<td></td>
				<td></td>
				<td> <i class="fa fa-clock-o text-gdq-red" aria-hidden="true"></i> </td>
				</tr>
				<tr>
				<td> <i class="fa fa-clock-o" aria-hidden="true"></i>  </td>
				<td></td>
				<td><i class="fa fa-microphone"></i> </td>
				</tr></tbody></table></html>`)
		rows := s.FindAll("tr")
		ev, err := eventFromHTML(rows[0], rows[1])
		assertEqual(t, err, nil)
		assertEqual(t, ev.Title, "")
		assertEqual(t, ev.Platform, "unknown")
		assertEqual(t, ev.Category, "unknown")
		assertEqual(t, len(ev.Hosts), 0)
		assertEqual(t, len(ev.Runners), 0)
		assertEqual(t, ev.Setup.Duration, time.Duration(0))
		assertEqual(t, ev.Estimate.Duration, time.Duration(0))
		assertEqual(t, ev.Start, time.Time{})
	})
	t.Run("single valid event", func(t *testing.T) {
		s := soup.HTMLParse(`<html><table id="runTable"><tbody><tr>
				<td>2021-01-03T17:00:00Z</td>
				<td>Awesome&#039;s Sauce</td>
				<td>my_runner</td>
				<td> <i class="fa fa-clock-o text-gdq-red" aria-hidden="true"></i> 0:17:00 </td>
				</tr>
				<tr>
				<td> <i class="fa fa-clock-o" aria-hidden="true"></i> 0:46:00 </td>
				<td>Any% &mdash; PC</td>
				<td><i class="fa fa-microphone"></i> my_host</td>
				</tr></tbody></table></html>`)

		rows := s.FindAll("tr")
		ev, err := eventFromHTML(rows[0], rows[1])
		assertEqual(t, err, nil)
		assertEqual(t, ev.Title, "Awesome's Sauce")
		assertEqual(t, ev.Platform, "PC")
		assertEqual(t, ev.Category, "Any%")
		assertEqual(t, len(ev.Hosts), 1)
		assertEqual(t, ev.Hosts[0], "my_host")
		assertEqual(t, len(ev.Runners), 1)
		assertEqual(t, ev.Runners[0], "my_runner")
		assertEqual(t, ev.Setup.Duration, time.Duration(17*time.Minute))
		assertEqual(t, ev.Estimate.Duration, time.Duration(46*time.Minute))
		assertEqual(t, ev.Start, time.Date(2021, 01, 03, 17, 00, 00, 00, time.UTC))
	})
}
