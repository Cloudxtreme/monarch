package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
)

var (
	httpClient = &http.Client{}
)

type DependentService struct {
	Ip   string
	Port int
}

func (d *DependentService) call(m Monarch) (int, TesterResponse) {

	var response TesterResponse
	jsonStr, _ := json.Marshal(m)
	reader := bytes.NewReader(jsonStr)
	fullUrl := "http://" + d.Ip + ":" + strconv.Itoa(d.Port) + "/work"

	req, err := http.NewRequest("POST", fullUrl, reader)
	req.Header.Add("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
	defer resp.Body.Close()

	if err != nil {
		fmt.Println("Error: ", err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal(body, &response)

	response.addHop(hostname)
	response.setEndTime()
	response.RoundTripTime = response.EndTime - response.BackendTime

	fmt.Println("Post result: ", response.MonarchNm)

	return resp.StatusCode, response
}
