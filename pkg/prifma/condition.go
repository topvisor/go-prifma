package prifma

import (
	"fmt"
	"github.com/topvisor/prifma/pkg/utils"
	"golang.org/x/net/idna"
	"net/http"
	"strings"
)

type Condition interface {
	Test(req *http.Request) bool
}

func NewCondition(key string, typ string, val string) (Condition, error) {
	tester, err := NewConditionTester(typ, val)
	if err != nil {
		return nil, err
	}

	switch true {
	case key == "src_ip":
		return NewConditionSrcIp(tester), nil
	case key == "dst_domain":
		return NewConditionDstDomain(tester), nil
	case strings.HasPrefix(key, "header_"):
		return NewConditionHeader(tester, key), nil
	}

	return nil, fmt.Errorf("unavailable condition key: '%s'", key)
}

type ConditionSrcIp struct {
	Tester ConditionTester
}

func NewConditionSrcIp(tester ConditionTester) Condition {
	return &ConditionSrcIp{
		Tester: tester,
	}
}

func (t *ConditionSrcIp) Test(req *http.Request) bool {
	return t.Tester.Test(req.RemoteAddr)
}

type ConditionDstDomain struct {
	Tester ConditionTester
}

func NewConditionDstDomain(tester ConditionTester) Condition {
	return &ConditionDstDomain{
		Tester: tester,
	}
}

func (t *ConditionDstDomain) Test(req *http.Request) bool {
	host, err := idna.ToUnicode(utils.GetRequestHostname(req))
	if err != nil {
		return false
	}

	return t.Tester.Test(host)
}

type ConditionHeader struct {
	Tester ConditionTester
	Name   string
}

func NewConditionHeader(tester ConditionTester, key string) Condition {
	name := strings.Replace(key, "header_", "", 1)
	name = strings.ReplaceAll(name, "_", "-")

	return &ConditionHeader{
		Tester: tester,
		Name:   name,
	}
}

func (t *ConditionHeader) Test(req *http.Request) bool {
	header := req.Header.Get(t.Name)

	return t.Tester.Test(header)
}
