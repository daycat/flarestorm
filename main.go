package main

import (
	_ "embed"
	"github.com/daycat/flarestorm/utils"
)

//go:embed range.txt
var rangeTxt string

func main() {
	utils.GetAllLocs(rangeTxt)
	/**
	rs["HKG"] = templates.Loc{Name: "HKG", Avgping: 1, Avgspeed: 2}
	entry, _ := rs["HKG"]
	entry.Addresses = append(entry.Addresses)
	entry.Addresses = append(entry.Addresses, templates.IP{"192.168.3.2", 1, 2})
	rs["HKG"] = entry
	utils.Mkreq(rs["HKG"].Addresses[0].IP, &wg)
	fmt.Println(rs["HKG"].Addresses)
	**/
}
