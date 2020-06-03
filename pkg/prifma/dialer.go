package prifma

import (
	"context"
	"errors"
	"github.com/topvisor/go-prifma/pkg/utils"
	"net"
)

var ErrOutgoingIpNotDefined = errors.New("outgoing ip address wasn't defined")

type Dialer interface {
	GetIpV4() net.IP
	GetIpV6() net.IP
	GetLocalIp(hostname string) (net.IP, error)

	SetIpV4(ip net.IP)
	SetIpV6(ip net.IP)

	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

func NewDialer() *DefaultDialer {
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

func (t *DefaultDialer) GetLocalIp(hostname string) (net.IP, error) {
	var localIp net.IP

	switch true {
	case t.IpV4 != nil && t.IpV6 != nil:
		dstIpV4, _, err := utils.LookupIp(hostname)
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
		return nil, ErrOutgoingIpNotDefined
	}

	return localIp, nil
}

func (t *DefaultDialer) SetIpV4(ip net.IP) {
	t.IpV4 = ip.To4()
}

func (t *DefaultDialer) SetIpV6(ip net.IP) {
	t.IpV6 = ip.To16()
}

func (t *DefaultDialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	host, _, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}

	localIp, err := t.GetLocalIp(host)
	if err != nil {
		return nil, err
	}

	t.Dialer.LocalAddr = &net.TCPAddr{
		IP: localIp,
	}

	return t.Dialer.DialContext(ctx, network, address)
}
