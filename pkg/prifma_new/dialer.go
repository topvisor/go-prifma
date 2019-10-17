package prifma_new

import (
	"context"
	"errors"
	"github.com/topvisor/prifma/pkg/utils"
	"net"
)

type Dialer interface {
	GetIpV4() net.IP
	GetIpV6() net.IP

	SetIpV4(ip net.IP)
	SetIpV6(ip net.IP)

	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

func NewDialer() Dialer {
	return new(DefaultDialer)
}

type DefaultDialer struct {
	IpV4   net.IP
	IpV6   net.IP
	Dialer net.Dialer
}

func (t *DefaultDialer) GetIpV4() net.IP {
	return t.IpV4
}

func (t *DefaultDialer) GetIpV6() net.IP {
	return t.IpV6
}

func (t *DefaultDialer) SetIpV4(ip net.IP) {
	t.IpV4 = ip
}

func (t *DefaultDialer) SetIpV6(ip net.IP) {
	t.IpV6 = ip
}

func (t *DefaultDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	var localIp net.IP

	switch true {
	case t.IpV4 != nil && t.IpV6 != nil:
		host, _, err := net.SplitHostPort(address)
		if err != nil {
			return nil, err
		}

		dstIpV4, _, err := utils.LookupIp(host)
		if err != nil {
			return nil, err
		}

		if dstIpV4 != nil {
			localIp = t.IpV4
		} else {
			localIp = t.IpV6
		}
	case t.IpV4 != nil:
		localIp = t.IpV4
	case t.IpV6 != nil:
		localIp = t.IpV6
	}

	if localIp == nil {
		return nil, errors.New("outgoing ip address wasn't defined")
	}

	t.Dialer.LocalAddr = &net.TCPAddr{
		IP: localIp,
	}

	return t.Dialer.DialContext(ctx, network, address)
}
