package main

import (
	"time"
)

type TesterResponse struct {
	MonarchCty    string `json:"cty"`
	MonarchHse    string `json:"hse"`
	MonarchNm     string `json:"nm"`
	MonarchYrs    string `json:"yrs"`
	Hops          []Hop  `json:"hops"`
	BackendTime   int64  `json:"backendTime"`
	EndTime       int64  `json:"endTime"`
	RoundTripTime int64  `json:"roundTripTime"`
}

type Hop struct {
	Id        int    `json:"id"`
	Host      string `json:"host"`
	TimeStamp int64  `json:"timeStamp"`
}

//  adds a Hop to the array of hops with regard to order
func (t *TesterResponse) addHop(hostname string) {
	var id int
	if len(t.Hops) == 0 {
		id = 0
	} else {
		id = t.Hops[len(t.Hops)-1].Id + 1
	}
	hop := Hop{id, hostname, t.setTime()}
	t.Hops = append(t.Hops, hop)
}

func (t *TesterResponse) setEndTime() {
	t.EndTime = t.setTime()
}

func (t *TesterResponse) setBackendTime() {
	t.BackendTime = t.setTime()
}

// sets time in milliseconds
func (t *TesterResponse) setTime() int64 {
	return time.Now().UnixNano() / 1000000
}
