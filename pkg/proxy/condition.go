package proxy

import (
	"fmt"
	"golang.org/x/net/idna"
	"net"
	"net/http"
	"regexp"
	"strings"
)

// conditionType determines a condition's value parser and a request's parameter which tested by condition
type conditionType byte

// ConditionType determines a condition's value parser and a request's parameter which tested by condition
//
// ConditionTypeSrcIpCIDR - request's source ip; it's tested by a CIDR subnet mask
// ConditionTypeDstDomainRegexp - request's destination domain; it's tested by a regular expression
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

// condition is a interface of different conditions types
type condition interface {
	test(req *http.Request) bool
}

// Condition is a tester of requests. It uses the type to parse the value and
// test a request's parameter by the parsed value
type Condition struct {
	Type  conditionType
	Value string

	tester condition
}

// parseConditionFromString parses the Condition from the string in the format "type:value".
// For example: "srcIpCIDR:1.2.3.4/32", "dstDomainRegexp:(?i)(www.)?example.com"
func parseConditionFromString(conditionStr string) (*Condition, error) {
	var err error
	condition := new(Condition)

	conditionStrs := strings.SplitN(conditionStr, ":", 2)
	if len(conditionStrs) != 2 {
		return nil, fmt.Errorf("parse condition from string Error: \"%s\"", conditionStr)
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

// test checks the request by condition
func (t *Condition) test(req *http.Request) bool {
	tester, err := t.getTester()
	if err != nil {
		return false
	}

	return tester.test(req)
}

// getTester generates a condition of the specified type
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

// conditionSrcIpCIDR is a condition which tested request's source ip by CIDR subnet mask
type conditionSrcIpCIDR net.IPNet

// parseConditionSrcIpCIDRFromString parses conditionSrcIpCIDR from string which contains a CIDR subnet mask
func parseConditionSrcIpCIDRFromString(conditionCIDRStr string) (*conditionSrcIpCIDR, error) {
	_, ipNet, err := net.ParseCIDR(conditionCIDRStr)
	if err != nil {
		return nil, err
	}

	return (*conditionSrcIpCIDR)(ipNet), err
}

// test checks the request's source ip by the CIDR subnet mask
func (t *conditionSrcIpCIDR) test(req *http.Request) bool {
	ip := net.ParseIP(req.RemoteAddr)
	if ip == nil {
		return false
	}

	return (*net.IPNet)(t).Contains(ip)
}

// conditionDstDomainRegexp is a condition which tested request's destination domain by a regular expression
type conditionDstDomainRegexp struct {
	regexp *regexp.Regexp
}

// parseConditionDstDomainRegexpFromString parses conditionDstDomainRegexp from string which contains a regular expression
func parseConditionDstDomainRegexpFromString(conditionRegexpStr string) (*conditionDstDomainRegexp, error) {
	compiledRegexp, err := regexp.Compile(conditionRegexpStr)
	if err != nil {
		return nil, err
	}

	return &conditionDstDomainRegexp{compiledRegexp}, err
}

// test checks the request's destination domain by the regular expression
func (t *conditionDstDomainRegexp) test(req *http.Request) bool {
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
