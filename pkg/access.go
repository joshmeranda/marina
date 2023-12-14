package marina

import (
	"slices"
)

type AccessType int

const (
	AccessTypeUnknown AccessType = iota
	AccessTypeAllow   AccessType = iota
	AccessTypeDeny    AccessType = iota
)

// AccessList is a list of names are that either Allowed or Denied.
type AccessList struct {
	AllowList []string `json:"allowList,omitempty"`
	DenyList  []string `json:"denyList,omitempty"`
}

func (l AccessList) GetAcecssFor(name string) AccessType {
	if slices.Contains(l.AllowList, name) {
		return AccessTypeAllow
	}

	if slices.Contains(l.DenyList, name) {
		return AccessTypeDeny
	}

	return AccessTypeUnknown
}

type UserAccessList struct {
	UserList AccessList `json:"userList,omitempty"`
	OrgList  AccessList `json:"orgList,omitempty"`
}

func (l UserAccessList) GetAccessFor(name string, orgs []string) AccessType {
	userAccess := l.UserList.GetAcecssFor(name)

	if userAccess == AccessTypeAllow {
		return AccessTypeAllow
	}

	if userAccess == AccessTypeDeny {
		return AccessTypeDeny
	}

	for _, org := range orgs {
		orgAccess := l.OrgList.GetAcecssFor(org)

		if orgAccess == AccessTypeAllow {
			return AccessTypeAllow
		}

		if orgAccess == AccessTypeDeny {
			return AccessTypeDeny
		}
	}

	return AccessTypeUnknown
}
