package prifma

import (
	"fmt"
	"net"
	"regexp"
)

type ConditionTester interface {
	Test(val string) bool
}

func NewConditionTester(typ string, val string) (ConditionTester, error) {
	switch typ {
	case "=":
		return NewConditionTesterEquals(val)
	case "~":
		return NewConditionTesterRegexp(val)
	case "cidr":
		return NewConditionTesterCIDR(val)
	}

	return nil, fmt.Errorf("unavailable condition type - '%s'", typ)
}

type ConditionTesterEquals struct {
	Value string
}

func NewConditionTesterEquals(val string) (ConditionTester, error) {
	return &ConditionTesterEquals{val}, nil
}

func (t *ConditionTesterEquals) Test(val string) bool {
	return t.Value == val
}

type ConditionTesterRegexp struct {
	Regexp *regexp.Regexp
}

func NewConditionTesterRegexp(val string) (ConditionTester, error) {
	regex, err := regexp.Compile(val)
	if err != nil {
		return nil, err
	}

	return &ConditionTesterRegexp{regex}, nil
}

func (t *ConditionTesterRegexp) Test(val string) bool {
	return t.Regexp.MatchString(val)
}

type ConditionTesterCIDR struct {
	Net *net.IPNet
}

func NewConditionTesterCIDR(val string) (ConditionTester, error) {
	_, ipNet, err := net.ParseCIDR(val)
	if err != nil {
		return nil, err
	}

	return &ConditionTesterCIDR{ipNet}, nil
}

func (t *ConditionTesterCIDR) Test(val string) bool {
	ip := net.ParseIP(val)
	if ip == nil {
		return false
	}

	return t.Net.Contains(ip)
}
