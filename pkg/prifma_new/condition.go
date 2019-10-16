package prifma_new

import (
	"fmt"
	"golang.org/x/net/idna"
	"net"
	"net/http"
	"strings"
)

type Condition interface {
	GetHash() ConditionHash
	Test(req *http.Request) bool
}

func NewCondition(key string, typ string, val string) (Condition, error) {
	tester, err := NewConditionTester(typ, val)
	if err != nil {
		return nil, err
	}

	base := ConditionBase{
		Tester: tester,
		ConditionHash: ConditionHash{
			Key:   key,
			Type:  typ,
			Value: val,
		},
	}

	switch true {
	case key == "src_ip":
		return &ConditionSrcIp{base}, nil
	case key == "dst_domain":
		return &ConditionDstDomain{base}, nil
	case strings.HasPrefix(key, "header_"):
		return &ConditionHeader{base}, nil
	}

	return nil, fmt.Errorf("unavailable condition key: '%s'", key)
}

type ConditionHash struct {
	Key   string
	Type  string
	Value string
}

func (t ConditionHash) GetHash() ConditionHash {
	return t
}

type ConditionBase struct {
	Tester ConditionTester

	ConditionHash
}

type ConditionSrcIp struct {
	ConditionBase
}

func (t *ConditionSrcIp) Test(req *http.Request) bool {
	return t.Tester.Test(req.RemoteAddr)
}

type ConditionDstDomain struct {
	ConditionBase
}

func (t *ConditionDstDomain) Test(req *http.Request) bool {
	host, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		return false
	}

	host, err = idna.ToUnicode(host)
	if err != nil {
		return false
	}

	return t.Tester.Test(host)
}

type ConditionHeader struct {
	ConditionBase
}

func (t *ConditionHeader) Test(req *http.Request) bool {
	name := t.GetName()
	header := req.Header.Get(name)

	return t.Tester.Test(header)
}

func (t *ConditionHeader) GetName() string {
	name := strings.Replace(t.Key, "header_", "", 1)
	name = strings.ReplaceAll(name, "_", "-")

	return name
}
