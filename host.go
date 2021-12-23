package gdq

import (
	"encoding/json"
)

type host struct {
	Handle string `json:"handle"`
	Runs   []uint `json:"runs"`
}

type hostResp struct {
	Fields struct {
		StartRun uint   `json:"start_run"`
		EndRun   uint   `json:"end_run"`
		Handle   string `json:"name"`
	} `json:"fields"`
}

func (h *host) UnmarshalJSON(b []byte) error {
	var v hostResp
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}

	h.Handle = v.Fields.Handle

	if v.Fields.StartRun > v.Fields.EndRun {
		start := v.Fields.StartRun
		end := v.Fields.EndRun

		v.Fields.StartRun = end
		v.Fields.EndRun = start
	}

	if v.Fields.StartRun == v.Fields.EndRun {
		h.Runs = []uint{v.Fields.StartRun}
	} else {
		seq := make([]uint, v.Fields.EndRun-v.Fields.StartRun+1)
		for i := range seq {
			seq[i] = v.Fields.StartRun + uint(i)
		}
		h.Runs = seq
	}

	return nil
}
