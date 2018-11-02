package tools

import (
	"github.com/miekg/dns"
	"net"
)

func DIG(domain, remote, client_ip string) ([]dns.RR, error) {
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
	e.Address = net.ParseIP(client_ip).To4()
	o.Option = append(o.Option, e)
	m.Extra = append(m.Extra, o)
	c := new(dns.Client)
	r, _, err := c.Exchange(m, remote+":53")
	if err != nil {
		return nil, err
	}
	return r.Answer, nil
}

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

func Num2IP(ipnum uint32) string {
	c1 := ipnum / 256 / 256 / 256
	c2 := (ipnum / 256 / 256) % 256
	c3 := (ipnum / 256) % 256
	c4 := ipnum % 256
	return net.IPv4(byte(c1), byte(c2), byte(c3), byte(c4)).String()
}
