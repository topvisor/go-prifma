package proxy

import (
	"errors"
	"fmt"
	"net"
	"regexp"
	"strings"
)

type conditionType int

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
		return -1, errors.New(fmt.Sprintf("unavailable condition type: \"%s\"", conditionTypeStr))
	}
}

type condition interface {
	Test(str string) bool
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
		return nil, errors.New(fmt.Sprintf("parse condition from string error: \"%s\"", conditionStr))
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

func (t *Condition) Test(str string) bool {
	tester, err := t.getTester()
	if err != nil {
		return false
	}

	return tester.Test(str)
}

func (t *Condition) getTester() (condition, error) {
	var err error

	if t.tester != nil {
		switch t.Type {
		case ConditionTypeSrcIpCIDR:
			if t.tester, err = parseConditionCIDRFromString(t.Value); err != nil {
				return nil, err
			}
		case ConditionTypeDstDomainRegexp:
			if t.tester, err = parseConditionRegexpFromString(t.Value); err != nil {
				return nil, err
			}
		default:
			return nil, errors.New(fmt.Sprintf("unavailable condition type: \"%v\"", t.Type))
		}
	}

	return t.tester, nil
}

type conditionCIDR net.IPNet

func parseConditionCIDRFromString(conditionCIDSStr string) (*conditionCIDR, error) {
	_, ipNet, err := net.ParseCIDR(conditionCIDSStr)
	if err != nil {
		return nil, err
	}

	return (*conditionCIDR)(ipNet), err
}

func (t *conditionCIDR) Test(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}

	return (*net.IPNet)(t).Contains(ip)
}

type conditionRegexp struct {
	regexp *regexp.Regexp
}

func parseConditionRegexpFromString(conditionRegexpStr string) (*conditionRegexp, error) {
	compiledRegexp, err := regexp.Compile(conditionRegexpStr)
	if err != nil {
		return nil, err
	}

	return &conditionRegexp{compiledRegexp}, err
}

func (t *conditionRegexp) Test(str string) bool {
	return t.regexp.MatchString(str)
}
