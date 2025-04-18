package dnsfunc

import (
	"net"
	"slices"
	"strconv"
	"strings"
	"github.com/mascarenhasmelson/Grafana-Network-Monitoring-System/Backend/internal/utils"
)

func getHostByID(idStr string, hosts []utils.Host) (oneHost utils.Host) {

	id, _ := strconv.Atoi(idStr)

	for _, host := range hosts {
		if host.ID == id {
			oneHost = host
			break
		}
	}

	return oneHost
}

func UpdateDNS(host utils.Host) (name, dns string) {

	dnsNames, _ := net.LookupAddr(host.IP)

	if len(dnsNames) > 0 {
		name = dnsNames[0]
		dns = strings.Join(dnsNames, " ")
	}

	return name, dns
}

func getHostsByMAC(mac string, hosts []utils.Host) (foundHosts []utils.Host) {

	for _, host := range hosts {
		if host.Mac == mac {

			foundHosts = append(foundHosts, host)
		}
	}

	return foundHosts
}

// host := getHostByID(idStr, allHosts)
// 	_, host.DNS = updateDNS(host)
// 	frontend = host
func getAllIfaces(hosts []utils.Host) (ifaces []string) {

	for _, host := range hosts {
		if !slices.Contains(ifaces, host.Iface) {
			ifaces = append(ifaces, host.Iface)
		}
	}
	return ifaces
}
