package controller

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	marinacorev1 "github.com/joshmeranda/marina/api/v1"
)

var _ = Describe("User Controller", Ordered, func() {
	var reconciler *UserReconciler
	var namespace *corev1.Namespace
	var ctx context.Context
	var user *marinacorev1.User

	BeforeAll(func() {
		ctx = context.Background()

		reconciler = &UserReconciler{
			Client: k8sClient,
		}

		namespace = &corev1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name:      "marina-system",
				Namespace: "marina-system",
			},
		}

		err := k8sClient.Create(context.Background(), namespace)
		if !errors.IsAlreadyExists(err) {
			Expect(err).NotTo(HaveOccurred())
		}

		roles := []rbacv1.Role{
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "SomeRole",
					Namespace: namespace.Name,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "AnotherRole",
					Namespace: namespace.Name,
				},
			},
			{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "NewRole",
					Namespace: namespace.Name,
				},
			},
		}

		for _, role := range roles {
			err := k8sClient.Create(ctx, &role)
			if !errors.IsAlreadyExists(err) {
				Expect(err).NotTo(HaveOccurred())
			}
		}

		user = &marinacorev1.User{
			ObjectMeta: metav1.ObjectMeta{Name: "user-test", Namespace: namespace.Name},
			Spec: marinacorev1.UserSpec{
				Name:     "bilbo",
				Password: []byte("baggins"),
				Roles:    []string{"SomeRole", "AnotherRole"},
			},
		}

		err = k8sClient.Create(ctx, user)
		Expect(err).NotTo(HaveOccurred())
	})

	When("user with roles is created", func() {
		BeforeAll(func() {
			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: user.Namespace,
					Name:      user.Name,
				},
			}
			result, err := reconciler.Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(ctrl.Result{}))
		})

		It("should update user roles with self-role", func() {
			err := k8sClient.Get(ctx, client.ObjectKeyFromObject(user), user)
			Expect(err).NotTo(HaveOccurred())
			Expect(user.Spec.Roles).To(ContainElements("SomeRole", "AnotherRole", user.Name+"-self"))
		})

		It("should create user service account", func() {
			var serviceaccount corev1.ServiceAccount
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name,
				Namespace: user.Namespace,
			}, &serviceaccount)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create user role", func() {
			var role rbacv1.Role
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name + "-" + "self",
				Namespace: user.Namespace,
			}, &role)
			Expect(err).NotTo(HaveOccurred())
		})

		It("should create role bindings for user roles", func() {
			var roleBinding rbacv1.RoleBinding
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name + "-" + "SomeRole",
				Namespace: user.Namespace,
			}, &roleBinding)
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name + "-" + "AnotherRole",
				Namespace: user.Namespace,
			}, &roleBinding)
			Expect(err).NotTo(HaveOccurred())

			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name + "-" + user.Name + "-" + "self",
				Namespace: user.Namespace,
			}, &roleBinding)
			Expect(err).NotTo(HaveOccurred())
		})
	})

	When("user is updated", func() {
		BeforeAll(func() {
			user.Spec.Roles = []string{"SomeRole", "NewRole", "user-test-self"}

			err := k8sClient.Update(ctx, user)
			Expect(err).ToNot(HaveOccurred())

			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: user.Namespace,
					Name:      user.Name,
				},
			}

			result, err := reconciler.Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(ctrl.Result{}))
		})

		It("should create rolebindings for new roles", func() {
			var roleBinding rbacv1.RoleBinding
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name + "-" + "NewRole",
				Namespace: user.Namespace,
			}, &roleBinding)
			Expect(err).ToNot(HaveOccurred())
		})

		It("should delete rolebindings for removed roles", func() {
			var roleBinding rbacv1.RoleBinding
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name + "-" + "AnotherRole",
				Namespace: user.Namespace,
			}, &roleBinding)
			Expect(err).To(HaveOccurred())
			Expect(roleBinding).To(BeZero())
		})
	})

	When("user is deleted", func() {
		BeforeAll(func() {
			err := k8sClient.Delete(ctx, user)
			Expect(err).NotTo(HaveOccurred())

			req := ctrl.Request{
				NamespacedName: types.NamespacedName{
					Namespace: user.Namespace,
					Name:      user.Name,
				},
			}
			result, err := reconciler.Reconcile(ctx, req)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(Equal(ctrl.Result{}))
		})

		It("should delete user service account", func() {
			var serviceaccount corev1.ServiceAccount
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name,
				Namespace: user.Namespace,
			}, &serviceaccount)
			Expect(err).To(HaveOccurred())
			Expect(serviceaccount).To(BeZero())
		})

		It("should delete role bindings for user roles", func() {
			var roleBinding rbacv1.RoleBinding
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name + "-" + "SomeRole",
				Namespace: user.Namespace,
			}, &roleBinding)
			Expect(err).To(HaveOccurred())
			Expect(roleBinding).To(BeZero())

			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name + "-" + "AnotherRole",
				Namespace: user.Namespace,
			}, &roleBinding)
			Expect(err).To(HaveOccurred())
			Expect(roleBinding).To(BeZero())

			err = k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name + "-" + user.Name + "-" + "self",
				Namespace: user.Namespace,
			}, &roleBinding)
			Expect(err).To(HaveOccurred())
			Expect(roleBinding).To(BeZero())
		})

		It("should delete user roles", func() {
			var role rbacv1.Role
			err := k8sClient.Get(ctx, types.NamespacedName{
				Name:      user.Name + "-" + "self",
				Namespace: user.Namespace,
			}, &role)
			Expect(err).To(HaveOccurred())
			Expect(role).To(BeZero())
		})
	})
})
