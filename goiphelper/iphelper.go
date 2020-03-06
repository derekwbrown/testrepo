package main

import (
	"encoding/binary"
	"fmt"
	"net"
	"github.com/DataDog/datadog-agent/pkg/util/winutil"
	"github.com/DataDog/datadog-agent/pkg/util/winutil/iphelper"
)



func int2ip(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

func testTcpTable() {
	t, err := iphelper.GetExtendedTcpV4Table()
	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}
	for pid, list := range t {
		for _, entry := range list {
			fmt.Printf("Pid (%8v)    %16v:%5v  %16v:%5v\n",
						pid,
						int2ip(iphelper.Ntohl(entry.DwLocalAddr)).String(),
						iphelper.Ntohs(uint16(entry.DwLocalPort)),
						int2ip(iphelper.Ntohl(entry.DwRemoteAddr)).String(),
						iphelper.Ntohs(uint16(entry.DwRemotePort)))
		}

	}
	return
}


func testRouteTable() {
	adapters, err := iphelper.GetAdaptersAddresses()
	t, err := iphelper.GetIPv4RouteTable()
	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}
	fmt.Printf("  Destination         Netmask              Gateway          InterfaceID    Metric\n")
	fmt.Printf("=================================================================================\n")
	for _, entry := range t {
		adapter, present := adapters[entry.DwForwardIfIndex]
		var interfaceid string
		interfaceid = string(entry.DwForwardIfIndex)
		if present {
			for _, a := range adapter.UnicastAddresses {
				as4 := a.Address.To4()
				if as4 != nil {
					interfaceid = as4.String()
					break
				}
			}
		}
		fmt.Printf("%15v   %15v    %15v    %12v    %6v\n",
					int2ip(iphelper.Ntohl(entry.DwForwardDest)).String(),
					int2ip(iphelper.Ntohl(entry.DwForwardMask)).String(),
					int2ip(iphelper.Ntohl(entry.DwForwardNextHop)).String(),
					interfaceid,
					entry.DwForwardMetric1)
	}
}

func testIfTable() {
	t, err := iphelper.GetIFTable()
	if err != nil {
		fmt.Printf("error %v\n", err)
		return
	}
	fmt.Printf("  Index         Name              Description\n")
	fmt.Printf("=============================================\n")
	for _, entry := range t {
		fmt.Printf("%5v   %20v %v\n",
					entry.DwIndex, winutil.ConvertWindowsString16(entry.WszName[:]), winutil.ConvertASCIIString(entry.BDescr[:]))
	}
}

func testAddressList() {
	t, err := iphelper.GetAdaptersAddresses()
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}
	fmt.Printf("%v\n", t)
}

func main() {
//	testIfTable()
//	testAddressList()
	testRouteTable()
}