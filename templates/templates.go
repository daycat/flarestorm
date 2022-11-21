package templates

import "time"

type Loc struct {
	Name      string  `json:"name"`
	Avgping   float64 `json:"avgping"`
	Avgspeed  float64 `json:"avgspeed"`
	Addresses []IP
}

type IP struct {
	IP    string
	Ping  time.Duration
	Speed float64
}
