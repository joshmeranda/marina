package core_test

import (
	"fmt"
	"testing"

	"github.com/joshmeranda/marina/gateway/api/core"
)

func TestStringQuery(t *testing.T) {
	type TestCase struct {
		Name    string
		Value   string
		Query   core.StringQuery
		Matches bool
		Err     error
	}

	testCases := []*TestCase{
		{
			Name:  "MatchEqual",
			Value: "abc",
			Query: core.StringQuery{
				Value:   "abc",
				MatchOp: core.StringMatchOp_Equal,
			},
			Matches: true,
		},
		{
			Name:  "NoMatchEqual",
			Value: "abc",
			Query: core.StringQuery{
				Value:   "ab",
				MatchOp: core.StringMatchOp_Equal,
			},
		},

		{
			Name:  "MatchNotEqual",
			Value: "abc",
			Query: core.StringQuery{
				Value:   "ab",
				MatchOp: core.StringMatchOp_NotEqual,
			},
			Matches: true,
		},
		{
			Name:  "NoMatchNotEqual",
			Value: "abc",
			Query: core.StringQuery{
				Value:   "abc",
				MatchOp: core.StringMatchOp_NotEqual,
			},
		},

		{
			Name:  "MatchRegexMatch",
			Value: "abc",
			Query: core.StringQuery{
				Value:   "a.{1}c",
				MatchOp: core.StringMatchOp_RegexMatch,
			},
			Matches: true,
		},
		{
			Name:  "NoMatchRegexMatch",
			Value: "abc",
			Query: core.StringQuery{
				Value:   "a.{2}c",
				MatchOp: core.StringMatchOp_RegexMatch,
			},
		},
		{
			Name: "RegexMatchInvalid",
			Query: core.StringQuery{
				Value:   "[",
				MatchOp: core.StringMatchOp_RegexMatch,
			},
			Err: fmt.Errorf("error parsing regexp: missing closing ]: `[`"),
		},

		{
			Name:  "MatchRegexNoMatch",
			Value: "abc",
			Query: core.StringQuery{
				Value:   "a.{2}c",
				MatchOp: core.StringMatchOp_RegexNotMatch,
			},
			Matches: true,
		},
		{
			Name:  "NoMatchRegexMatch",
			Value: "abc",
			Query: core.StringQuery{
				Value:   "a.{1}c",
				MatchOp: core.StringMatchOp_RegexNotMatch,
			},
		},
		{
			Name: "RegexNoMatchInvalid",
			Query: core.StringQuery{
				Value:   "[",
				MatchOp: core.StringMatchOp_RegexNotMatch,
			},
			Err: fmt.Errorf("error parsing regexp: missing closing ]: `[`"),
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			matches, err := testCase.Query.Matches(testCase.Value)

			if matches != testCase.Matches {
				t.Fatalf("expected testCase.Matches == %v, but found %v", testCase.Matches, matches)
			}

			if err == nil {
				if testCase.Err != nil {
					t.Fatalf("expected err to be '%s' but was nil", testCase.Err)
				}
			} else {
				if testCase.Err == nil {
					t.Fatalf("expected err to be nil but found '%s'", err)
				} else if testCase.Err.Error() != err.Error() {
					t.Fatalf("expected err to be '%s' but was '%s'", testCase.Err, err)
				}
			}
		})
	}
}

func TestCollectionQuery(t *testing.T) {
	type TestCase struct {
		Name    string
		Value   []string
		Query   core.CollectionQuery
		Matches bool
		Err     error
	}

	testCases := []*TestCase{
		{
			Name:  "MatchAllOf",
			Value: []string{"a", "b", "c"},
			Query: core.CollectionQuery{
				Values:  []string{"b", "c"},
				MatchOp: core.CollectionMatchOp_ContainsAllOf,
			},
			Matches: true,
		},
		{
			Name:  "NoMatchAllOf",
			Value: []string{"a", "b", "c"},
			Query: core.CollectionQuery{
				Values:  []string{"b", "d"},
				MatchOp: core.CollectionMatchOp_ContainsAllOf,
			},
		},

		{
			Name:  "MatchAnyOf",
			Value: []string{"a", "b", "c"},
			Query: core.CollectionQuery{
				Values:  []string{"b", "d"},
				MatchOp: core.CollectionMatchOp_ContainsAnyOf,
			},
			Matches: true,
		},
		{
			Name:  "NoMatchAnyOf",
			Value: []string{"a", "b", "c"},
			Query: core.CollectionQuery{
				Values:  []string{"x", "y", "z"},
				MatchOp: core.CollectionMatchOp_ContainsAnyOf,
			},
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.Name, func(t *testing.T) {
			matches, err := testCase.Query.Matches(testCase.Value)

			if matches != testCase.Matches {
				t.Fatalf("expected testCase.Matches == %v, but found %v", testCase.Matches, matches)
			}

			if err == nil {
				if testCase.Err != nil {
					t.Fatalf("expected err to be '%s' but was nil", testCase.Err)
				}
			} else {
				if testCase.Err == nil {
					t.Fatalf("expected err to be nil but found '%s'", err)
				} else if testCase.Err.Error() != err.Error() {
					t.Fatalf("expected err to be '%s' but was '%s'", testCase.Err, err)
				}
			}
		})
	}
}
