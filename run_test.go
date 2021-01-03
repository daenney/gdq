package gdq

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToRun(t *testing.T) {
	r := runResp{}
	ru := r.toRun()
	assert.Equal(t, Run{
		Hosts:   []string{},
		Runners: Runners{},
	}, ru)
}
