package gdq

type hostResp struct {
	Fields struct {
		Run    uint   `json:"start_run"`
		Handle string `json:"name"`
	} `json:"fields"`
}

type hostsResp []hostResp
