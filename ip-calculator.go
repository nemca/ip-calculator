package main

import (
	"encoding/binary"
	"fmt"
	"log"
	"net"
	"os"
	"path/filepath"
)

const result = "%-10s %v\n"

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("  Usage:   %s CIDR\n", filepath.Base(os.Args[0]))
		fmt.Printf("  Example: %s 192.168.34.27/24\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	ip, ipnet, err := net.ParseCIDR(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	bitmask, _ := ipnet.Mask.Size()
	networkIP, broadcastIP, wildcardIP, netmask := networkRange(ipnet)
	size := networkSize(ipnet.Mask)

	fmt.Printf(result, "Address:", ip)
	fmt.Printf(result, "Bitmask:", bitmask)
	fmt.Printf(result, "Netmask:", netmask)
	fmt.Printf(result, "Wildcard:", wildcardIP)
	fmt.Printf(result, "Network:", ipnet.IP)
	fmt.Printf(result, "HostMin:", networkIPInc(networkIP))
	fmt.Printf(result, "HostMax:", networkIPDec(broadcastIP))
	fmt.Printf(result, "Broadcast:", broadcastIP)
	fmt.Printf(result, "Hosts:", size)
}

// Calculates the first and last IP addresses in an IPNet
func networkRange(network *net.IPNet) (net.IP, net.IP, net.IP, net.IP) {
	netIP := network.IP.To4()
	networkIP := netIP.Mask(network.Mask)
	broadcastIP := net.IPv4(0, 0, 0, 0).To4()
	wildcardIP := net.IPv4(0, 0, 0, 0).To4()
	networkMask := net.IPv4(0, 0, 0, 0).To4()
	for i := 0; i < len(broadcastIP); i++ {
		broadcastIP[i] = netIP[i] | ^network.Mask[i]
		wildcardIP[i] = net.IPv4bcast[i] | ^network.Mask[i]
		networkMask[i] = ^wildcardIP[i]
	}
	return networkIP, broadcastIP, wildcardIP, networkMask
}

// Given a netmask, calculates the number of available hosts
func networkSize(mask net.IPMask) int32 {
	m := net.IPv4Mask(0, 0, 0, 0)
	for i := 0; i < net.IPv4len; i++ {
		m[i] = ^mask[i]
	}
	return int32(binary.BigEndian.Uint32(m)) + 1
}

func networkIPInc(ip net.IP) net.IP {
	minIPNum := ipToInt(ip.To4()) + 1
	return intToIP(minIPNum)
}

func networkIPDec(ip net.IP) net.IP {
	maxIPNum := ipToInt(ip.To4()) - 1
	return intToIP(maxIPNum)
}

// Converts a 4 bytes IP into a 32 bit integer
func ipToInt(ip net.IP) int32 {
	return int32(binary.BigEndian.Uint32(ip.To4()))
}

// Converts 32 bit integer into a 4 bytes IP address
func intToIP(n int32) net.IP {
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, uint32(n))
	return net.IP(b)
}
