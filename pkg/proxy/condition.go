package proxy

import (
	"fmt"
	"golang.org/x/net/idna"
	"net"
	"net/http"
	"regexp"
	"strings"
)

// conditionType determines condition's value parser and request's parameter which tested by condition
type conditionType byte

// ConditionType determines condition's value parser and request's parameter which tested by condition
//
// ConditionTypeSrcIpCIDR - request's source ip; it's tested by CIDR subnet mask
// ConditionTypeDstDomainRegexp - request's destination domain; it's tested by regular expression
const (
	ConditionTypeSrcIpCIDR conditionType = iota
	ConditionTypeDstDomainRegexp
)

// conditionTypeFromString parses conditionType from string.
func conditionTypeFromString(conditionTypeStr string) (*conditionType, error) {
	switch conditionTypeStr {
	case "srcIpCIDR":
		conditionType := ConditionTypeSrcIpCIDR
		return &conditionType, nil
	case "dstDomainRegexp":
		conditionType := ConditionTypeDstDomainRegexp
		return &conditionType, nil
	default:
		return nil, fmt.Errorf("unavailable condition type: \"%s\"", conditionTypeStr)
	}
}

type condition interface {
	test(req *http.Request) bool
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

	conditionType, err := conditionTypeFromString(conditionStrs[0])
	if err != nil {
		return nil, err
	}

	condition.Type = *conditionType
	condition.Value = conditionStrs[1]

	if _, err = condition.getTester(); err != nil {
		return nil, err
	}

	return condition, nil
}

func (t *Condition) test(req *http.Request) bool {
	tester, err := t.getTester()
	if err != nil {
		return false
	}

	return tester.test(req)
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

func (t *conditionCIDR) test(req *http.Request) bool {
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

func (t *conditionRegexp) test(req *http.Request) bool {
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
