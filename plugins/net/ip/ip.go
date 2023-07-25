package ip

import (
	"bytes"
	"encoding/binary"
	"errors"
	"github.com/AyakuraYuki/go-live-broadcast-downloader/plugins/file"
	"net"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// Hostname returns system hostname
func Hostname() string {
	name, err := os.Hostname()
	if err != nil {
		return file.ReadFile("/etc/hostname")
	}
	return name
}

// PrivateIPv4 returns private IP in IPv4
func PrivateIPv4() (net.IP, error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return nil, err
	}
	for _, addr := range addrs {
		ipNet, ok := addr.(*net.IPNet)
		if !ok || ipNet.IP.IsLoopback() {
			continue
		}
		ip := ipNet.IP.To4()
		if IsPrivateIPv4(ip) {
			return ip, nil
		}
	}
	return nil, errors.New("no private ipv4 address")
}

// IsPrivateIPv4 returns if an ip address is in IPv4 and is private ip
func IsPrivateIPv4(ip net.IP) bool {
	return ip != nil &&
		(ip[0] == 10 || ip[0] == 172 && (ip[1] >= 16 && ip[1] < 32) || ip[0] == 192 && ip[1] == 168)
}

func Lower16BitPrivateIP() (uint16, error) {
	ip, err := PrivateIPv4()
	if err != nil {
		return 0, err
	}
	return uint16(ip[2])<<8 + uint16(ip[3]), nil
}

// IP2Int converts ip to int number
func IP2Int(ip net.IP) uint32 {
	if len(ip) == 16 {
		return binary.BigEndian.Uint32(ip[12:16])
	}
	return binary.BigEndian.Uint32(ip)
}

// Int2IP converts a given number to ip
func Int2IP(nn uint32) net.IP {
	ip := make(net.IP, 4)
	binary.BigEndian.PutUint32(ip, nn)
	return ip
}

// PrivateIP2Int PrivateIP2Int
func PrivateIP2Int() uint32 {
	ip, err := PrivateIPv4()
	if err != nil {
		return 0
	}
	return IP2Int(ip)
}

// LocalIP get public loc ip of network.
// get env BINDHOSTIP first
// using net interface name first if provided. or it will using env INAME
// using eth0 as default.
func LocalIP(optionalIName ...string) string {
	if ip := os.Getenv("BINDHOSTIP"); ip != "" {
		return ip
	}
	ifaces, err := net.Interfaces()
	if err != nil {
		return ""
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return ""
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String()
		}
	}
	var name string
	if len(optionalIName) != 0 && optionalIName[0] != "" {
		name = optionalIName[0]
	} else if name = os.Getenv("INAME"); name == "" {
		name = "eth0"
		if runtime.GOOS == "darwin" {
			name = "en0"
		}
	}
	n, err := net.InterfaceByName(name)
	if err != nil {
		return ""
	}
	addrs, err := n.Addrs()
	if err != nil {
		return ""
	}
	for i := range addrs {
		if ipnet, ok := addrs[i].(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ip := ipnet.IP.To4(); ip != nil {
				return ip.String()
			}
		}
	}
	return ""
}

func StringIpToInt(ipstring string) int {
	ipSegs := strings.Split(ipstring, ".")
	var ipInt = 0
	var pos uint = 24
	for _, ipSeg := range ipSegs {
		tempInt, _ := strconv.Atoi(ipSeg)
		tempInt = tempInt << pos
		ipInt = ipInt | tempInt
		pos -= 8
	}
	return ipInt
}

// IPIntToString IPIntToString
func IPIntToString(ipInt int) string {
	ipSegs := make([]string, 4)
	var size = len(ipSegs)
	buffer := bytes.NewBufferString("")
	for i := 0; i < size; i++ {
		tempInt := ipInt & 0xFF
		ipSegs[size-i-1] = strconv.Itoa(tempInt)
		ipInt = ipInt >> 8
	}
	for i := 0; i < size; i++ {
		buffer.WriteString(ipSegs[i])
		if i < size-1 {
			buffer.WriteString(".")
		}
	}
	return buffer.String()
}
