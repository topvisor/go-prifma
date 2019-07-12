package proxy

import (
	"fmt"
	"golang.org/x/net/idna"
	"net"
	"net/http"
	"regexp"
	"strings"
)

type conditionType byte

const (
	ConditionTypeSrcIpCIDR conditionType = iota
	ConditionTypeDstDomainRegexp
)

func ConditionTypeFromString(conditionTypeStr string) (conditionType, error) {
	switch conditionTypeStr {
	case "srcIpCIDR":
		return ConditionTypeSrcIpCIDR, nil
	case "dstDomainRegexp":
		return ConditionTypeDstDomainRegexp, nil
	default:
		return -1, fmt.Errorf("unavailable condition type: \"%s\"", conditionTypeStr)
	}
}

type condition interface {
	Test(req *http.Request) bool
}

type Condition struct {
	Type  conditionType
	Value string

	tester condition
}

func ParseConditionFromString(conditionStr string) (*Condition, error) {
	var err error
	condition := new(Condition)

	conditionStrs := strings.SplitN(conditionStr, ":", 2)
	if len(conditionStrs) != 2 {
		return nil, fmt.Errorf("parse condition from string error: \"%s\"", conditionStr)
	}

	if condition.Type, err = ConditionTypeFromString(conditionStrs[0]); err != nil {
		return nil, err
	}

	condition.Value = conditionStrs[1]

	if _, err = condition.getTester(); err != nil {
		return nil, err
	}

	return condition, nil
}

func (t *Condition) Test(req *http.Request) bool {
	tester, err := t.getTester()
	if err != nil {
		return false
	}

	return tester.Test(req)
}

func (t *Condition) getTester() (condition, error) {
	var err error

	if t.tester != nil {
		switch t.Type {
		case ConditionTypeSrcIpCIDR:
			if t.tester, err = parseConditionSrcIpCIDRFromString(t.Value); err != nil {
				return nil, err
			}
		case ConditionTypeDstDomainRegexp:
			if t.tester, err = parseConditionDstDomainRegexpFromString(t.Value); err != nil {
				return nil, err
			}
		default:
			return nil, fmt.Errorf("unavailable condition type: \"%v\"", t.Type)
		}
	}

	return t.tester, nil
}

type conditionCIDR net.IPNet

func parseConditionSrcIpCIDRFromString(conditionCIDSStr string) (*conditionCIDR, error) {
	_, ipNet, err := net.ParseCIDR(conditionCIDSStr)
	if err != nil {
		return nil, err
	}

	return (*conditionCIDR)(ipNet), err
}

func (t *conditionCIDR) Test(req *http.Request) bool {
	ip := net.ParseIP(req.RemoteAddr)
	if ip == nil {
		return false
	}

	return (*net.IPNet)(t).Contains(ip)
}

type conditionRegexp struct {
	regexp *regexp.Regexp
}

func parseConditionDstDomainRegexpFromString(conditionRegexpStr string) (*conditionRegexp, error) {
	compiledRegexp, err := regexp.Compile(conditionRegexpStr)
	if err != nil {
		return nil, err
	}

	return &conditionRegexp{compiledRegexp}, err
}

func (t *conditionRegexp) Test(req *http.Request) bool {
	host, _, err := net.SplitHostPort(req.Host)
	if err != nil {
		return false
	}

	host, err = idna.ToUnicode(host)
	if err != nil {
		return false
	}

	return t.regexp.MatchString(host)
}
