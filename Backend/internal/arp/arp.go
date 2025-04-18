package arp

import (
	"fmt"
	"log/slog"
	"os/exec"
	"strings"

	"github.com/mascarenhasmelson/Grafana-Network-Monitoring-System/Backend/internal/utils"
	"time"
)

func ScanIface(iface string) string {
	fmt.Println("Executing scan on interface:", iface)

	cmd := exec.Command("sudo", "arp-scan", "-glNx", "-I", iface)

	out, err := cmd.Output()
	// fmt.Println(out)
	if err != nil {
		fmt.Println("Error executing command:", err)
		return ""
	}

	return string(out)
}


func parseOutput(text, iface string) []utils.Host {
	var foundHosts = []utils.Host{}

	p := strings.Split(text, "\n")

	for _, host := range p {
		if host != "" {
			var oneHost utils.Host
			p := strings.Split(host, "	")
			oneHost.Iface = iface
			oneHost.IP = p[0]
			oneHost.Mac = p[1]
			oneHost.Hw = p[2]
			oneHost.Date = time.Now().Format("2006-01-02 15:04:05")
			oneHost.Now = 1
			foundHosts = append(foundHosts, oneHost)
			// fmt.Println(foundHosts, "found")
			fmt.Println("iface", oneHost.Iface)
			fmt.Println("ip", oneHost.IP)
			fmt.Println("mac", oneHost.Mac)
			fmt.Println("harware", oneHost.Hw)
			fmt.Println("time", oneHost.Date)

		}
	}

	return foundHosts
}

// Scan all interfaces
func Scan(ifaces string) []utils.Host {
	// var strs []string
	var text string
	var p []string
	var foundHosts = []utils.Host{}
	// arpArgs = args

	if ifaces != "" {

		p = strings.Split(ifaces, " ")

		for _, iface := range p {
			slog.Debug("Scanning interface " + iface)
			text = ScanIface(iface)
			slog.Debug("Found IPs: \n" + text)

			foundHosts = append(foundHosts, parseOutput(text, iface)...)
			// fmt.Println(foundHosts)
			// fmt.Println("___________________")
			// fmt.Println(text, iface)
		}
	}

	// for _, s := range strs {
	// 	slog.Debug("Scanning string " + s)
	// 	text = scanStr(s)
	// 	slog.Debug("Found IPs: \n" + text)
	// 	p = strings.Split(s, " ")

	// 	foundHosts = append(foundHosts, parseOutput(text, p[len(p)-1])...)
	// }

	return foundHosts
}
