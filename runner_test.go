package gdq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToRunner(t *testing.T) {
	r := runnerResp{}
	ru := r.toRunner()
	assert.Equal(t, Runner{}, ru)
}

func TestRunnersToString(t *testing.T) {
	r := Runners{
		{Handle: "one"},
		{Handle: "two"},
	}
	s := r.String()
	assert.Equal(t, "one, two", s)
}
