package utils

import (
	"fmt"
	"github.com/daycat/flarestorm/templates"
	"github.com/schollz/progressbar/v3"
	"io/ioutil"
	"net"
	"net/http"
	"net/netip"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	mu sync.Mutex
)

func Hosts(cidr string) string {
	prefix, err := netip.ParsePrefix(cidr)
	if err != nil {
		panic(err)
	}
	addr := prefix.Addr().Next().String()
	// addr := prefix.Addr().String()
	return addr
}

func tcping(ip string) (time.Duration, error) {
	timeThen := time.Now()
	_, err := net.DialTimeout("tcp", ip+":80", 3*time.Second)
	rtt := time.Since(timeThen)
	return rtt, err
}

func Mkreq(ip string, wg *sync.WaitGroup, mu *sync.Mutex, rs map[string]templates.Loc) {
	//fmt.Println(ICMP(ip))
	rtt, err := tcping(ip)
	if err != nil {
		wg.Done()
		return
	}
	// fmt.Println(rtt)
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	url := "http://" + ip + "/cdn-cgi/trace"
	resp, err := client.Get(url)
	if err != nil {
		wg.Done()
		return
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		wg.Done()
		return
	}
	// fmt.Println(time.Since(timeThen) / 7)
	sb := string(body)
	fields := strings.Fields(sb)
	var location string
	if len(fields) >= 6 && strings.Contains(fields[6], "colo=") {
		location = strings.Split(fields[6], "colo=")[1]
	} else {
		wg.Done()
		return
	}
	// appends ip and info to rs map
	mu.Lock()
	entry, ok := rs[location]
	if ok != true {
		rs[location] = templates.Loc{Name: location, Avgping: 0, Avgspeed: 0}
	}
	entry.Addresses = append(entry.Addresses, templates.IP{IP: ip, Ping: rtt, Speed: -1})
	rs[location] = entry
	mu.Unlock()
	wg.Done()
	return
}

func GetAllLocs(rangetxt string) map[string]templates.Loc {
	var (
		testip []string
		wg     sync.WaitGroup
		rs     map[string]templates.Loc
		tping  time.Duration
		tnum   int64
	)
	sem := make(chan int, 500)
	IP_Range := strings.Fields(rangetxt)
	for i := 0; i < len(IP_Range); i++ {
		testip = append(testip, Hosts(IP_Range[i]))
	}
	rs = make(map[string]templates.Loc)
	bar := progressbar.Default(int64(len(testip) - 1))
	for i := 0; i < len(testip)-1; i++ {
		wg.Add(1)
		_ = bar.Add(1)
		sem <- 1
		go func() {
			Mkreq(testip[i], &wg, &mu, rs)
			<-sem
		}()
	}
	wg.Wait()
	fmt.Println("The following locations are available:")
	for key := range rs {
		tnum, tping = 0, 0
		for i := 0; i < len(rs[key].Addresses); i++ {
			tnum++
			tping += rs[key].Addresses[i].Ping
		}
		entry := rs[key]
		entry.Avgping = float64(tping.Milliseconds() / tnum)
		rs[key] = entry
		fmt.Print(key, " | ", rs[key].Avgping, "ms, ")
	}
	_ = os.WriteFile("rs.txt", []byte(fmt.Sprint(rs)), 0644)
	return rs
}
