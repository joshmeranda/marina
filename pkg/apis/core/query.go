package core

import (
	"fmt"
	"regexp"
	"slices"
)

func (q *StringQuery) Matches(s string) (bool, error) {
	switch q.MatchOp {
	case StringMatchOp_Equal:
		return q.Value == s, nil
	case StringMatchOp_NotEqual:
		return q.Value != s, nil
	case StringMatchOp_RegexMatch:
		return regexp.MatchString(q.Value, s)
	case StringMatchOp_RegexNotMatch:
		matches, err := regexp.MatchString(q.Value, s)
		if err != nil {
			return false, err
		}
		return !matches, nil
	default:
		return false, fmt.Errorf("unknown string match op: %v", q.MatchOp)
	}
}

func (q *CollectionQuery) Matches(c []string) (bool, error) {
	switch q.MatchOp {
	case CollectionMatchOp_ContainsAllOf:
		for _, v := range q.Values {
			if !slices.Contains(c, v) {
				return false, nil
			}
		}

		return true, nil
	case CollectionMatchOp_ContainsAnyOf:
		for _, v := range q.Values {
			if slices.Contains(c, v) {
				return true, nil
			}
		}
		return false, nil
	default:
		return false, fmt.Errorf("unknown collection match op: %v", q.MatchOp)
	}
}
