package prifma

import (
	"fmt"
	"net"
	"regexp"
)

type ConditionTester interface {
	Test(val string) bool
}

func NewConditionTester(typ string, val string) (tester ConditionTester, err error) {
	isNegation := false

	if typ[0] == '!' {
		isNegation = true
		typ = typ[1:]
	}

	switch typ {
	case "=":
		tester, err = NewConditionTesterEquals(val)
	case "~":
		tester, err = NewConditionTesterRegexp(val)
	case "cidr":
		tester, err = NewConditionTesterCIDR(val)
	default:
		err = fmt.Errorf("unavailable condition type - '%s'", typ)
	}

	if err == nil && isNegation {
		tester = NewConditionTesterNegation(tester)
	}

	return tester, err
}

type ConditionTesterNegation struct {
	Tester ConditionTester
}

func (t *ConditionTesterNegation) Test(val string) bool {
	return !t.Tester.Test(val)
}

func NewConditionTesterNegation(tester ConditionTester) *ConditionTesterNegation {
	return &ConditionTesterNegation{tester}
}

type ConditionTesterEquals struct {
	Value string
}

func NewConditionTesterEquals(val string) (*ConditionTesterEquals, error) {
	return &ConditionTesterEquals{val}, nil
}

func (t *ConditionTesterEquals) Test(val string) bool {
	return t.Value == val
}

type ConditionTesterRegexp struct {
	Regexp *regexp.Regexp
}

func NewConditionTesterRegexp(val string) (*ConditionTesterRegexp, error) {
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

func NewConditionTesterCIDR(val string) (*ConditionTesterCIDR, error) {
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
