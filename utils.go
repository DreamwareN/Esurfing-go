package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"net"
	"net/http"
	"strings"
	"time"
)

func GetInterfaceIP(interfaceName string) (string, error) {
	iFace, err := net.InterfaceByName(interfaceName)
	if err != nil {
		return "", fmt.Errorf("interface not found: %v", err)
	}

	if iFace.Flags&net.FlagUp == 0 {
		return "", fmt.Errorf("interface %s is down", interfaceName)
	}

	addresses, err := iFace.Addrs()
	if err != nil {
		return "", fmt.Errorf("can not get addresses from interface %s: %v", interfaceName, err)
	}

	for _, addr := range addresses {
		var ip net.IP
		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		default:
			continue
		}

		if ip == nil || ip.IsLoopback() || ip.IsLinkLocalUnicast() {
			continue
		}

		ipv4 := ip.To4()
		if ipv4 == nil {
			continue
		}

		return ipv4.String(), nil
	}

	return "", fmt.Errorf("no available ipv4 address at interface %s", interfaceName)
}

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func GenerateRandomString(length int) string {
	if length <= 0 {
		return ""
	}

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[rand.N(len(charset))]
	}
	return string(b)
}

func GenerateRandomMAC() string {
	mac := make([]byte, 6)

	for i := range mac {
		mac[i] = byte(rand.N(256))
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

func FormatEConfig(data []byte) ([]byte, error) {
	str1 := strings.Split(string(data), ConfigStartTag)
	str2 := strings.Split(str1[1], ConfigEndTag)

	str3 := strings.ReplaceAll(str2[0], "&width=0", "")
	str4 := strings.ReplaceAll(str3, "&adtype=0", "")

	return []byte(str4), nil
}

func NewHttpTransport(c *Config) (http.RoundTripper, error) {
	if c.BindInterface != "" {
		ip, err := GetInterfaceIP(c.BindInterface)
		if err != nil {
			return nil, errors.New(fmt.Errorf("failed to get interface IP: %w", err).Error())
		}

		localAddr := &net.TCPAddr{IP: net.ParseIP(ip)}
		return &http.Transport{
			DialContext: (&net.Dialer{
				LocalAddr: localAddr,
				Resolver:  GetResolver(c),
			}).DialContext,
		}, nil
	} else {
		return &http.Transport{
			DialContext: (&net.Dialer{
				Resolver: GetResolver(c),
			}).DialContext,
		}, nil
	}
}

func GetResolver(c *Config) *net.Resolver {
	if c.DnsAddress == "" {
		return net.DefaultResolver
	}

	return &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: 5 * time.Second,
			}
			return d.DialContext(ctx, "udp", c.DnsAddress)
		},
	}
}
