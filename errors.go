package gdq

// Error represents an error returned by this package
type err string

func (e err) Error() string { return string(e) }

const (
	// ErrMissingSchedule means we couldn't find a schedule at all
	ErrMissingSchedule err = "missing schedule"
	// ErrInvalidSchedule means the schedule we found isn't conforming to our expectations
	ErrInvalidSchedule err = "invalid schedule"
	// ErrUnexpectedData means we encountered a row in a different format
	ErrUnexpectedData err = "row did not contain the data we expected"
)
