package utils

import (
	"crypto/rand"
	"errors"
	"fmt"
	"net"
	"strings"
)

func GetInterfaceIP(interfaceName string) (string, error) {
	iFace, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", fmt.Errorf("interface unexist: %v", err)
	}

	addresses, err := iFace.Addrs()
	if err != nil {
		return "", fmt.Errorf("can not get interface IP: %v", err)
	}

	for _, addr := range addresses {
		ipNet, ok := addr.(*net.IPNet)
		if !ok {
			continue
		}
		ip := ipNet.IP
		if ip.IsLoopback() || ip.IsLinkLocalUnicast() {
			continue
		}
		if ip.To4() != nil {
			return ip.String(), nil
		}
	}

	return "", fmt.Errorf("no available ip at interface %s", interfaceName)
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return ""
	}

	for i := 0; i < length; i++ {
		bytes[i] = charset[int(bytes[i])%len(charset)]
	}

	return string(bytes)
}

func GenerateRandomMAC() string {
	mac := make([]byte, 6)

	_, err := rand.Read(mac)
	if err != nil {
		return ""
	}

	mac[0] = (mac[0] & 0xfe) | 0x02

	return net.HardwareAddr(mac).String()
}

func DecodeAlgoID(data []byte) (algoID string, key string, err error) {
	dataLen := len(data)
	if dataLen < 4 {
		return "", "", errors.New("data Error: insufficient header length")
	}

	len1 := int(data[3])
	pos := 4

	if pos+len1 > dataLen {
		return "", "", errors.New("data Error: key length exceeds data size")
	}
	keyBytes := data[pos : pos+len1]
	pos += len1

	if pos >= dataLen {
		return "", "", errors.New("data Error: missing algoID header")
	}

	len2 := int(data[pos])
	pos++

	if pos+len2 > dataLen {
		return "", "", errors.New("data Error: algoID length exceeds data size")
	}
	algoIDBytes := data[pos : pos+len2]

	return string(algoIDBytes), string(keyBytes), nil
}

const ConfigStartTag = "<!--//config.campus.js.chinatelecom.com "
const ConfigEndTag = "//config.campus.js.chinatelecom.com-->"

func ParseEConfig(data []byte) ([]byte, error) {
	str1 := strings.Split(string(data), ConfigStartTag)
	str2 := strings.Split(str1[1], ConfigEndTag)
	
	str3 := strings.ReplaceAll(str2[0], "&width=0", "")
	str4 := strings.ReplaceAll(str3, "&adtype=0", "")

	return []byte(str4), nil
}
