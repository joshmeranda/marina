package marina_test

import (
	marina "github.com/joshmeranda/marina/pkg"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AccessList", func() {
	var list marina.AccessList

	BeforeEach(func() {
		list = marina.AccessList{}
		list.SetAccessFor("bbagins", marina.AccessTypeAllow)
		list.SetAccessFor("fbaggins", marina.AccessTypeDeny)
	})

	It("returns expected access levels", func() {
		Expect(list.GetAccessFor("bbagins")).To(Equal(marina.AccessTypeAllow))
		Expect(list.GetAccessFor("fbaggins")).To(Equal(marina.AccessTypeDeny))
		Expect(list.GetAccessFor("gandalf")).To(Equal(marina.AccessTypeUnknown))
	})

	It("sets access levels", func() {
		// set new value
		list.SetAccessFor("gandalf", marina.AccessTypeAllow)
		Expect(list.GetAccessFor("gandalf")).To(Equal(marina.AccessTypeAllow))

		// change existing value
		list.SetAccessFor("bbaggins", marina.AccessTypeDeny)
		Expect(list.GetAccessFor("bbaggins")).To(Equal(marina.AccessTypeDeny))
	})
})

var _ = Describe("UserAccessList", func() {
	var list marina.UserAccessList

	BeforeEach(func() {
		list = marina.UserAccessList{}
		list.SetAccessForUser("bbagins", marina.AccessTypeAllow)
		list.SetAccessForUser("fbaggins", marina.AccessTypeDeny)

		list.SetAccessForGroup("shire", marina.AccessTypeAllow)
		list.SetAccessForGroup("fellowship", marina.AccessTypeAllow)
		list.SetAccessForGroup("wizard", marina.AccessTypeAllow)
		list.SetAccessForGroup("mordor", marina.AccessTypeDeny)
	})

	It("returns expected access levels", func() {
		Expect(list.UserList.GetAccessFor("bbagins")).To(Equal(marina.AccessTypeAllow))
		Expect(list.UserList.GetAccessFor("fbaggins")).To(Equal(marina.AccessTypeDeny))

		Expect(list.GroupList.GetAccessFor(("shire"))).To(Equal(marina.AccessTypeAllow))
		Expect(list.GroupList.GetAccessFor(("fellowship"))).To(Equal(marina.AccessTypeAllow))
		Expect(list.GroupList.GetAccessFor(("wizard"))).To(Equal(marina.AccessTypeAllow))
		Expect(list.GroupList.GetAccessFor(("mordor"))).To(Equal(marina.AccessTypeDeny))
	})

	When("user and group is unknown", func() {
		It("denies access", func() {
			Expect(list.GetAccessFor("tbombadil", []string{"old-forest"})).ToNot(Equal(marina.AccessTypeAllow))
		})
	})

	When("user is unknown and group is allowed", func() {
		It("allows access", func() {
			Expect(list.GetAccessFor("sgamgee", []string{"shire"})).To(Equal(marina.AccessTypeAllow))
		})
	})

	When("user is unknown and group is denied", func() {
		It("denies access", func() {
			Expect(list.GetAccessFor("morgoth", []string{"mordor"})).ToNot(Equal(marina.AccessTypeAllow))
		})
	})

	When("user is allowed and group is unknown", func() {
		It("allows access", func() {
			Expect(list.GetAccessFor("bbagins", []string{"ring-bearer"})).To(Equal(marina.AccessTypeAllow))
		})
	})

	When("user and group is allowed", func() {
		It("allows access", func() {
			Expect(list.GetAccessFor("bbagins", []string{"shire"})).To(Equal(marina.AccessTypeAllow))
		})
	})

	When("user is allowed and group is denied", func() {
		It("allows access", func() {
			Expect(list.GetAccessFor("bbagins", []string{"mordor"})).To(Equal(marina.AccessTypeAllow))
		})
	})

	When("user is denied and group is unknown", func() {
		It("denies access", func() {
			Expect(list.GetAccessFor("fbaggins", []string{"ring-bearer"})).ToNot(Equal(marina.AccessTypeAllow))
		})
	})

	When("user is denied and group is allowed", func() {
		It("denies access", func() {
			Expect(list.GetAccessFor("fbaggins", []string{"shire"})).ToNot(Equal(marina.AccessTypeAllow))
		})
	})

	When("user and group is denied", func() {
		It("denies access", func() {
			Expect(list.GetAccessFor("fbaggin", []string{"mordor"})).ToNot(Equal(marina.AccessTypeAllow))
		})
	})
})
