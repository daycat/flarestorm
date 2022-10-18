package main

import (
	_ "embed"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/netip"
	"strings"
	"sync"
	"time"
)

//go:embed range.txt
var rangefile string
var (
	mu    sync.Mutex
	colos = make(map[string][]string)
)

func Hosts(cidr string) string {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		panic(err)
	}

	var ips = []string{}
	count := 1
	for addr := prefix.Addr(); prefix.Contains(addr); addr = addr.Next() {
		ips = append(ips, addr.String())
		count += 1
	}

	return ips[rand.Intn(253-1)+1]
	//return ips[17]
}

func GetURL(ip string, wg *sync.WaitGroup, colos map[string][]string) string {
	client := http.Client{
		Timeout: 3 * time.Second,
	}
	url := "http://" + ip + "/cdn-cgi/trace"
	resp, err := client.Get(url)
	if err != nil {
		wg.Done()
		return "0"
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wg.Done()
		return "0"
	}
	sb := string(body)
	splitted := strings.Fields(sb)
	var colosx string
	if len(splitted) >= 6 && strings.Contains(splitted[6], "colo=") {
		colosx = strings.Split(splitted[6], "colo=")[1]
	} else {
		wg.Done()
		return sb
	}
	mu.Lock()
	colos[colosx] = append(colos[colosx], ip)
	mu.Unlock()
	wg.Done()
	return sb
}

func init() {
	log.SetFlags(log.Ldate | log.Ltime)
}

func main() {
	var (
		wg       sync.WaitGroup
		ipranges []string
		testip   []string
		timeNow  = time.Now()
	)
	log.Println("----------flarestorm v0.1----------")
	ipranges = strings.Fields(rangefile)
	for i := 0; i < len(ipranges); i++ {
		testip = append(testip, Hosts(ipranges[i]))
	}
	log.Println("Parsing done in", time.Since(timeNow))
	log.Println("Now testing IP")
	for i := 0; i < len(testip); i++ {
		wg.Add(1)
		go GetURL(testip[i], &wg, colos)
	}
	wg.Wait()
	log.Println("Test finished in ", time.Since(timeNow))
	for key, value := range colos {
		fmt.Print("----------", key, " has ", len(value), " subnets", "----------\n")
		fmt.Print(value, "\n")
	}
}
