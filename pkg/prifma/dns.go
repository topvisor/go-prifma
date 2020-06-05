package prifma

import "net"

type DNSIps struct {
	IpV4 net.IP
	IpV6 net.IP
}

type DNS interface {
	SetCache(host string, ips DNSIps)
	ClearCache()

	LookupIp(host string) (DNSIps, error)
}

func NewDNS() DefaultDNS {
	return make(DefaultDNS)
}

type DefaultDNS map[string]DNSIps

func (t DefaultDNS) SetCache(host string, ips DNSIps) {
	t[host] = ips
}

func (t DefaultDNS) ClearCache() {
	for k := range t {
		delete(t, k)
	}
}

func (t DefaultDNS) LookupIp(host string) (DNSIps, error) {
	ips, ok := t[host]

	if !ok {
		ipsRaw, err := net.LookupIP(host)
		if err != nil {
			return ips, err
		}

		for _, ip := range ipsRaw {
			if ipV4Tmp := ip.To4(); ipV4Tmp != nil {
				if ips.IpV4 == nil {
					ips.IpV4 = ipV4Tmp
				}
			} else {
				if ips.IpV6 == nil {
					ips.IpV6 = ip
				}
			}

			if ips.IpV4 != nil && ips.IpV6 != nil {
				break
			}
		}

		t[host] = ips
	}

	return ips, nil
}
