package utils

import (
	"context"
	"fmt"
	"github.com/daycat/flarestorm/templates"
	"net"
	"net/http"
	"time"
)

func getDialContext(ip string) func(ctx context.Context, network, address string) (net.Conn, error) {
	fakeSourceAddr := ip + ":443"
	return func(ctx context.Context, network, address string) (net.Conn, error) {
		return (&net.Dialer{}).DialContext(ctx, network, fakeSourceAddr)
	}
}

func bestloc(rs map[string]templates.Loc) (string, float64) {
	var (
		minping  float64 = 1000
		bestcolo string  = ""
	)
	for key := range rs {
		if rs[key].Avgping < minping {
			minping = rs[key].Avgping
			bestcolo = key
		}
	}
	return bestcolo, minping
}

func getspeed(ip string) int64 {
	/*var (
		ctx     context.Context
		network string
	)

	*/
	client := http.Client{
		Timeout:   5 * time.Second,
		Transport: &http.Transport{DialContext: getDialContext(ip)},
	}

	req, _ := http.NewRequest("GET", "https://speedtest.daycat.space/16m.bin", nil)
	//req.Header.Add("Host", "speedtest.daycat.space")
	timeNow := time.Now()
	resp, err := client.Do(req)
	print(resp.ContentLength)
	if err != nil {
		fmt.Println(err)
		return -1
	}
	// defer res.Body.Close()
	timeTaken := time.Since(timeNow)
	//fmt.Println(timeTaken)
	return resp.ContentLength / (timeTaken.Milliseconds() * 1000.0)
}

func dltest(colo templates.Loc) {
	var (
		sorted bool = false
	)
	// bubble sort array according to ping
	for sorted != true {
		sorted = true
		for i := 0; i < len(colo.Addresses)-1; i++ {
			if colo.Addresses[i].Ping > colo.Addresses[i+1].Ping {
				entry := colo.Addresses[i+1]
				colo.Addresses[i+1] = colo.Addresses[i]
				colo.Addresses[i] = entry
				sorted = false
			}
		}
	}
	//ump.P(colo)
	for i := 0; i < len(colo.Addresses)-1 && i < 5; i++ {
		fmt.Println("Speed of ", colo.Addresses[i].IP, ": ", getspeed(colo.Addresses[i].IP), "MB/s")
	}
}

func Speedtest(rs map[string]templates.Loc) {
	bestcolo, minping := bestloc(rs)
	fmt.Printf("\nTesting %v, with ping %vms", bestcolo, minping)
	dltest(rs[bestcolo])
}
