package user_test

import (
	"testing"

	marinav1 "github.com/joshmeranda/marina/api/v1"
	core "github.com/joshmeranda/marina/gateway/api/core"
	"github.com/joshmeranda/marina/gateway/api/user"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestUserQuery(t *testing.T) {
	type TestCase struct {
		Name    string
		Value   *marinav1.User
		Query   *user.UserQuery
		Matches bool
		Err     error
	}

	bbaggins := &marinav1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "bbaggins",
		},
		Spec: marinav1.UserSpec{
			Roles: []string{"shire"},
		},
	}

	testCases := []TestCase{
		{
			Name:  "NameAndRoleMatch",
			Value: bbaggins,
			Query: &user.UserQuery{
				Name: &core.StringQuery{
					Value:   "bbaggins",
					MatchOp: core.StringMatchOp_Equal,
				},
				Roles: &core.CollectionQuery{
					Values:  []string{"shire"},
					MatchOp: core.CollectionMatchOp_ContainsAllOf,
				},
			},
			Matches: true,
		},
		{
			Name:  "OnlyNameMatches",
			Value: bbaggins,
			Query: &user.UserQuery{
				Name: &core.StringQuery{
					Value:   "bbaggins",
					MatchOp: core.StringMatchOp_Equal,
				},
				Roles: &core.CollectionQuery{
					Values:  []string{"wizard"},
					MatchOp: core.CollectionMatchOp_ContainsAllOf,
				},
			},
		},
		{
			Name:  "OnlyRoleMatches",
			Value: bbaggins,
			Query: &user.UserQuery{
				Name: &core.StringQuery{
					Value:   "fbaggins",
					MatchOp: core.StringMatchOp_Equal,
				},
				Roles: &core.CollectionQuery{
					Values:  []string{"shire"},
					MatchOp: core.CollectionMatchOp_ContainsAllOf,
				},
			},
		},
		{
			Name:  "NothingMatches",
			Value: bbaggins,
			Query: &user.UserQuery{
				Name: &core.StringQuery{
					Value:   "gandalf",
					MatchOp: core.StringMatchOp_Equal,
				},
				Roles: &core.CollectionQuery{
					Values:  []string{"wizard"},
					MatchOp: core.CollectionMatchOp_ContainsAllOf,
				},
			},
		},

		{
			Name:  "UserIsNil",
			Value: bbaggins,
			Query: &user.UserQuery{
				Roles: &core.CollectionQuery{
					Values:  []string{"shire"},
					MatchOp: core.CollectionMatchOp_ContainsAllOf,
				},
			},
			Matches: true,
		},
		{
			Name:  "RoleIsNil",
			Value: bbaggins,
			Query: &user.UserQuery{
				Name: &core.StringQuery{
					Value:   "bbaggins",
					MatchOp: core.StringMatchOp_Equal,
				},
			},
			Matches: true,
		},
		{
			Name:    "EverythingIsNil",
			Value:   bbaggins,
			Query:   &user.UserQuery{},
			Matches: true,
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
