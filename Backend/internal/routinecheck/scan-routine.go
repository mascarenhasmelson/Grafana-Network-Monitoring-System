package routinecheck

import (
	"fmt"

	// "os"
	"time"

	"test/internal/arp"
	"test/internal/db"
	"test/internal/dnsfunc"
	"test/internal/env"
	"test/internal/utils"
)

var (
	allHosts  []utils.Host
	histHosts []utils.Host
	// appconfig []utils.Conf
	quitScan chan bool
)
var boolvar bool

func UpdateRoutines() {

	fmt.Println("restarting scan")
	if quitScan != nil {
		close(quitScan) // Close the quitScan channel if it's not nil
	}
	fmt.Println("Creating database...")
	db.Create()
	fmt.Println("DB created successfully")
	allHosts = db.Select("now")
	fmt.Println(allHosts)
	quitScan = make(chan bool)
	fmt.Println("startscan")
	StartScan(quitScan) // scan-routine.go
}

func StartScan(quit chan bool) {

	var lastDate, nowDate, plusDate time.Time
	var foundHosts []utils.Host
	fmt.Println("Scan routine started")
	for {
		select {
		case <-quit:
			fmt.Println("Received quit signal, stopping goroutine.")
			return
		default:
			nowDate = time.Now()
			fmt.Println(nowDate, "now data im here time 150")
			plusDate = lastDate.Add(time.Duration(50) * time.Second)

			// var m int
			// m++
			// fmt.Println(plusDate, "pulsedate", m)

			if nowDate.After(plusDate) {

				eVar := env.GetEnv()
				fmt.Println(eVar)
				foundHosts = arp.Scan(eVar)
				fmt.Println(foundHosts)
				compareHosts(foundHosts)
				allHosts = db.Select("now")

				lastDate = time.Now()
			}

			time.Sleep(time.Duration(1) * time.Second)
		}
	}
}

func compareHosts(foundHosts []utils.Host) {
	fmt.Println("comparehost")
	foundHostsMap := make(map[string]utils.Host)
	for _, fHost := range foundHosts {
		foundHostsMap[fHost.Mac] = fHost
		fmt.Println("hi", fHost)

		// db.Insert("history", fHost)
	}
	for _, aHost := range allHosts {
		fmt.Println("not ")
		fHost, exists := foundHostsMap[aHost.Mac]
		fmt.Println(exists, "exting")
		if exists {

			aHost.Iface = fHost.Iface
			aHost.IP = fHost.IP
			aHost.Date = fHost.Date
			aHost.Now = 1
			fmt.Println("Updating host:", aHost)
			// db.Insert("history", aHost)
			delete(foundHostsMap, aHost.Mac)

		} else {
			aHost.Now = 0
		}
		// db.Update("now", aHost)
		db.Update("now", aHost)

		aHost.Date = time.Now().Format("2006-01-02 15:04:05")

		histHosts = append(histHosts, aHost)
		fmt.Println(histHosts, "hist")

		if !boolvar {
			db.Insert("history", aHost)

		}
		fmt.Println("before end")
	}
	for _, fHost := range foundHostsMap {

		fHost.Name, fHost.DNS = dnsfunc.UpdateDNS(fHost)

		msg := fmt.Sprintf("Unknown host Names: '%s', IP: '%s', MAC: '%s', Hw: '%s', Iface: '%s'", fHost.DNS, fHost.IP, fHost.Mac, fHost.Hw, fHost.Iface)
		fmt.Println(msg)
		fmt.Println("Name", fHost.Name)
		fmt.Println("IP", fHost.IP)
		db.Insert("now", fHost)
	}

	fmt.Println("end")
}
