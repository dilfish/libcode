// Copyright 2018 Sean.ZH

package tools

import (
	"fmt"
	"github.com/miekg/dns"
	"net"
)

// DIG like dig
func DIG(domain, remote, clientIP string) ([]dns.RR, error) {
	m := new(dns.Msg)
	m.SetQuestion(domain, dns.TypeA)
	o := new(dns.OPT)
	o.Hdr.Name = "."
	o.Hdr.Rrtype = dns.TypeOPT
	e := new(dns.EDNS0_SUBNET)
	e.Code = dns.EDNS0SUBNET
	e.Family = 1
	e.SourceNetmask = 32
	e.SourceScope = 0
	e.Address = net.ParseIP(clientIP).To4()
	o.Option = append(o.Option, e)
	m.Extra = append(m.Extra, o)
	c := new(dns.Client)
	r, _, err := c.Exchange(m, remote+":53")
	if err != nil {
		return nil, err
	}
	return r.Answer, nil
}

// IP2Num ipv4 to a uint32 number
func IP2Num(ipstr string) uint32 {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	if ip == nil {
		return 0
	}
	n := uint32(0)
	n = n + uint32(ip[0])*256*256*256
	n = n + uint32(ip[1])*256*256
	n = n + uint32(ip[2])*256
	n = n + uint32(ip[3])
	return n
}

// Num2IP uint32 number to ipv4
func Num2IP(ipnum uint32) string {
	c1 := ipnum / 256 / 256 / 256
	c2 := (ipnum / 256 / 256) % 256
	c3 := (ipnum / 256) % 256
	c4 := ipnum % 256
	return net.IPv4(byte(c1), byte(c2), byte(c3), byte(c4)).String()
}

// Num2IPv6 makes two uint64 into a string ipv6
func Num2IPv6(network, intf uint64) string {
	c1 := intf >> 56
	c2 := (intf >> 48) % 256
	c3 := (intf >> 40) % 256
	c4 := (intf >> 32) % 256
	c5 := (intf >> 24) % 256
	c6 := (intf >> 16) % 256
	c7 := (intf >> 8) % 256
	c8 := intf % 256

	n1 := network >> 56
	n2 := (network >> 48) % 256
	n3 := (network >> 40) % 256
	n4 := (network >> 32) % 256
	n5 := (network >> 24) % 256
	n6 := (network >> 16) % 256
	n7 := (network >> 8) % 256
	n8 := network % 256

	netStr := fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x", n1, n2, n3, n4, n5, n6, n7, n8)
	intfStr := fmt.Sprintf("%02x%02x:%02x%02x:%02x%02x:%02x%02x", c1, c2, c3, c4, c5, c6, c7, c8)
	return netStr + ":" + intfStr
}

// IPv62Num makes a string to 2 uint64 number
func IPv62Num(ipstr string) (uint64, uint64) {
	ip := net.ParseIP(ipstr)
	if ip == nil {
		return 0, 0
	}
	ip = ip.To16()
	if ip == nil {
		return 0, 0
	}
	inet := uint64(0)
	iint := uint64(0)
	for i := uint(0); i < 8; i++ {
		inet = inet + (uint64(ip[i]) << (8 * (7 - i)))
	}
	for i := uint(0); i < 8; i++ {
		iint = iint + (uint64(ip[i+8] << (8 * (7 - i))))
	}
	return inet, iint
}
