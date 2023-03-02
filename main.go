package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func ipToInt(ip string) uint32 {

	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		log.Fatalf("Invalid IP address: %s", ip)
	}
	ipBytes := parsedIP.To4()
	if ipBytes == nil {
		log.Fatalf("Invalid IPv4 address: %s", ip)
	}
	return (uint32(ipBytes[0]) << 24) + (uint32(ipBytes[1]) << 16) + (uint32(ipBytes[2]) << 8) + uint32(ipBytes[3])
}

func intToIP(num uint32) string {
	first := (num >> 24) & 0xff
	second := (num >> 16) & 0xff
	third := (num >> 8) & 0xff
	fourth := num & 0xff
	return fmt.Sprintf("%d.%d.%d.%d", first, second, third, fourth)
}

func calculateSubnetMask(startIP string, endIP string) uint32 {
	startNum := ipToInt(startIP)
	endNum := ipToInt(endIP)

	xorNum := startNum ^ endNum
	xorBin := strconv.FormatInt(int64(xorNum), 2)

	prefixLen := 0
	for i := 0; i < len(xorBin); i++ {
		if xorBin[i:i+1] == "1" {
			prefixLen = i + 1
		}
	}

	subnetMaskDecimal := uint32(32 - prefixLen)

	return subnetMaskDecimal
}

func main() {
	file, err := os.Open("ips-range.csv")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	for {
		row, err := reader.Read()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			log.Fatal(err)
		}

		startIP := row[0]
		endIP := row[1]

		subnetMask := calculateSubnetMask(startIP, endIP)

		output, err := os.OpenFile("subnet_masks.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}
		defer output.Close()

		output.WriteString(fmt.Sprintf("%s/%d\n", startIP, subnetMask))
	}
}
