package user

import (
	marinav1 "github.com/joshmeranda/marina/api/v1"
)

func (q *UserQuery) Matches(user *marinav1.User) (bool, error) {

	if q.Name != nil {
		matches, err := q.Name.Matches(user.Name)
		if err != nil {
			return false, err
		}

		if !matches {
			return false, nil
		}
	}

	if q.Roles != nil {
		matches, err := q.Roles.Matches(user.Spec.Roles)
		if err != nil {
			return false, err
		}

		if !matches {
			return false, nil
		}
	}

	return true, nil
}
