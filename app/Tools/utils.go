package Tools

import (
	"log"
	"net/netip"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

var (
	// ClientList contains IPs/networks authorized to do purge/ban
	ClientList []netip.Prefix
)

// ReadDotEnvFile reads environment variables from .env file
func ReadDotEnvFile(f string) {
	err := godotenv.Load(f)
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// InitLog ensure log file exists and set appropriate flags (remove timestamp at start of line).
func InitLog(p string) {
	logFile, err := os.OpenFile(p, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
}

// InitAllowedIPList initialize the list of client authorized to do purge/ban requests
func InitAllowedIPList(l string) []netip.Prefix {
	list := []netip.Prefix{}
	if l != "" {
		sliceData := strings.Split(l, ",")
		for i := 0; i < len(sliceData); i++ {
			t, err := netip.ParsePrefix(sliceData[i])
			if err != nil {
				panic(err)
			}
			list = append(list, t)
		}
		return list
	}
	return list
}

// IPAllowed check if the IP is authorized to do BAN/PURGE requests
func IPAllowed(ip string) bool {
	ipAddr, err := netip.ParseAddr(ip)
	if err != nil {
		log.Printf("Ip address wrong format %s", err)
	}
	for i := 0; i < len(ClientList); i++ {
		if ClientList[i].Contains(ipAddr) {
			return true
		}
	}
	return false
}
