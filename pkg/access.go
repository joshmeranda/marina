package marina

type AccessType int

const (
	AccessTypeUnknown AccessType = iota
	AccessTypeAllow   AccessType = iota
	AccessTypeDeny    AccessType = iota
)

// AccessList is a list of names are that either Allowed or Denied.
type AccessList struct {
	inner map[string]AccessType
}

func (l *AccessList) SetAccessFor(name string, accessType AccessType) {
	if accessType == AccessTypeUnknown {
		delete(l.inner, name)
	}

	if l.inner == nil {
		l.inner = make(map[string]AccessType)
	}

	l.inner[name] = accessType
}

func (l *AccessList) GetAccessFor(name string) AccessType {
	accessType, ok := l.inner[name]
	if !ok {
		return AccessTypeUnknown
	}

	return accessType
}

type UserAccessList struct {
	UserList  AccessList `json:"userList,omitempty"`
	GroupList AccessList `json:"orgList,omitempty"`
}

func (l *UserAccessList) SetAccessForUser(name string, accessType AccessType) {
	l.UserList.SetAccessFor(name, accessType)
}

func (l *UserAccessList) SetAccessForGroup(name string, accessType AccessType) {
	l.GroupList.SetAccessFor(name, accessType)
}

func (l *UserAccessList) GetAccessFor(name string, groups []string) AccessType {
	userAccess := l.UserList.GetAccessFor(name)

	if userAccess == AccessTypeAllow {
		return AccessTypeAllow
	}

	if userAccess == AccessTypeDeny {
		return AccessTypeDeny
	}

	for _, group := range groups {
		groupAccess := l.GroupList.GetAccessFor(group)

		if groupAccess == AccessTypeAllow {
			return AccessTypeAllow
		}

		if groupAccess == AccessTypeDeny {
			return AccessTypeDeny
		}
	}

	return AccessTypeUnknown
}
