package controller

import (
	"context"
	"fmt"
	"slices"
	"strings"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	marinacorev1 "github.com/joshmeranda/marina/api/v1"
)

const (
	UserServiceAccountFinalizer = "marina.io.serviceaccount/finalizer"
	UserRoleBindingFinalizer    = "marina.io.rolebinding/finalizer"
	UserSelfRoleFinalizerFormat = "marina.io.selfrole.%s/finalizer"

	UserRoleBindingLabelUser = "app.marina.io/user"
	UserRoleBindingLabelRole = "app.marina.io/role"
)

func serviceAccountForUser(user *marinacorev1.User) *corev1.ServiceAccount {
	return &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      user.Name,
			Namespace: user.Namespace,
		},
	}
}

func selfRoleForUser(user *marinacorev1.User) *rbacv1.Role {
	return &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      user.Name + "-self",
			Namespace: user.Namespace,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups:     []string{"core.marina.io"},
				Resources:     []string{"users"},
				Verbs:         []string{"get", "list", "watch", "update"},
				ResourceNames: []string{user.Name},
			},
			{
				APIGroups: []string{"core.marina.io"},
				Resources: []string{"terminals"},
				Verbs:     []string{"get", "list", "watch", "create", "delete"},
			},
		},
	}
}

func userRoleBindingForRole(user *marinacorev1.User, role string) *rbacv1.RoleBinding {
	return &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      user.Name + "-" + role,
			Namespace: user.Namespace,
			Labels: map[string]string{
				UserRoleBindingLabelUser: user.Name,
				UserRoleBindingLabelRole: role,
			},
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      rbacv1.ServiceAccountKind,
				Name:      user.Name,
				Namespace: user.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     "Role",
			Name:     role,
			APIGroup: "rbac.authorization.k8s.io",
		},
	}
}

// UserReconciler reconciles a User object
type UserReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=core.marina.io,resources=users,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core.marina.io,resources=users/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=core.marina.io,resources=users/finalizers,verbs=update
// +kubebuilder:rbac:groups=*,resources=serviceaccounts,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=get;list;watch;create;update;patch;delete

func (r *UserReconciler) reconcileServiceAccount(ctx context.Context, user *marinacorev1.User) error {
	logger := log.FromContext(ctx)
	serviceAccount := serviceAccountForUser(user)

	if user.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(user, UserServiceAccountFinalizer) {
			if err := r.Delete(ctx, serviceAccount); err != nil {
				logger.Error(err, "could not delete service account", "serviceaccount", client.ObjectKeyFromObject(serviceAccount))
				return err
			}

			controllerutil.RemoveFinalizer(user, UserServiceAccountFinalizer)
		}

		return nil
	}

	_ = controllerutil.AddFinalizer(user, UserServiceAccountFinalizer)

	if err := r.Create(ctx, serviceAccount); err != nil {
		return client.IgnoreAlreadyExists(err)
	}

	logger.Info("created service account", "serviceaccount", client.ObjectKeyFromObject(serviceAccount))

	return nil
}

func (r *UserReconciler) reconcileNewRoles(ctx context.Context, user *marinacorev1.User, existingBindings *rbacv1.RoleBindingList) error {
	logger := log.FromContext(ctx)
	newRoles := []string{}

	for _, role := range user.Spec.Roles {
		if !slices.ContainsFunc(existingBindings.Items, func(binding rbacv1.RoleBinding) bool { return binding.Labels[UserRoleBindingLabelRole] == role }) {
			newRoles = append(newRoles, role)
		}
	}

	for _, role := range newRoles {
		binding := userRoleBindingForRole(user, role)

		if err := r.Create(ctx, binding); err != nil {
			return client.IgnoreAlreadyExists(err)
		}

		logger.Info("created role binding", "rolebinding", client.ObjectKeyFromObject(binding))
	}

	return nil
}

func (r *UserReconciler) reconcileRemovedRoles(ctx context.Context, user *marinacorev1.User, existingBindings *rbacv1.RoleBindingList) error {
	logger := log.FromContext(ctx)
	removedRoles := []*rbacv1.RoleBinding{}

	for _, binding := range existingBindings.Items {
		if !slices.Contains(user.Spec.Roles, binding.Labels[UserRoleBindingLabelRole]) {
			removedRoles = append(removedRoles, &binding)
		}
	}

	for _, role := range removedRoles {
		if err := r.Delete(ctx, role); err != nil {
			logger.Error(err, "error deleting role binding", "rolebinding", client.ObjectKeyFromObject(role))
			return err
		}

		logger.Info("deleted role binding", "rolebinding", client.ObjectKeyFromObject(role))
	}

	return nil
}

func (r *UserReconciler) reconcileRoleBindings(ctx context.Context, user *marinacorev1.User) error {
	logger := log.FromContext(ctx)
	isDeleting := user.GetDeletionTimestamp() != nil

	existingBindings := rbacv1.RoleBindingList{}
	if err := r.List(ctx, &existingBindings, client.InNamespace(user.Namespace), client.MatchingLabels{UserRoleBindingLabelUser: user.Name}); err != nil {
		logger.Error(err, "failed to list user role bindings", "user", client.ObjectKeyFromObject(user))
		return err
	}

	if isDeleting {
		for _, binding := range existingBindings.Items {
			if err := r.Delete(ctx, &binding); err != nil {
				logger.Error(err, "failed to delete role binding", "rolebinding", client.ObjectKeyFromObject(&binding))
				return nil
			}
		}
		_ = controllerutil.RemoveFinalizer(user, UserRoleBindingFinalizer)

		return nil
	}

	_ = controllerutil.AddFinalizer(user, UserRoleBindingFinalizer)

	if err := r.reconcileNewRoles(ctx, user, &existingBindings); err != nil {
		logger.Error(err, "failed to reconcile new roles", "user", client.ObjectKeyFromObject(user))
		return err
	}

	if err := r.reconcileRemovedRoles(ctx, user, &existingBindings); err != nil {
		logger.Error(err, "failed to reconcile removed roles", "user", client.ObjectKeyFromObject(user))
		return err
	}

	return nil
}

func (r *UserReconciler) reconcileUserSelfRole(ctx context.Context, user *marinacorev1.User) error {
	logger := log.FromContext(ctx)
	selfRole := selfRoleForUser(user)

	finalizerName := fmt.Sprintf(UserSelfRoleFinalizerFormat, strings.ReplaceAll(user.Name, "-", "."))

	if user.GetDeletionTimestamp() != nil {
		if controllerutil.ContainsFinalizer(user, finalizerName) {
			if err := r.Delete(ctx, selfRole); err != nil {
				logger.Error(err, "could not delete self role", "role", client.ObjectKeyFromObject(selfRole))
				return err
			}

			controllerutil.RemoveFinalizer(user, finalizerName)
		}

		return nil
	}

	_ = controllerutil.AddFinalizer(user, finalizerName)

	if err := r.Create(ctx, selfRole); err != nil {
		return client.IgnoreAlreadyExists(err)
	}

	logger.Info("created self role for user", "role", client.ObjectKeyFromObject(selfRole))

	user.Spec.Roles = append(user.Spec.Roles, selfRole.Name)

	return nil
}

func (r *UserReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	user := &marinacorev1.User{}

	if err := r.Get(ctx, req.NamespacedName, user); err != nil {
		logger.Error(err, "error fethcing user", "user", req.NamespacedName)
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.reconcileServiceAccount(ctx, user); err != nil {
		logger.Error(err, "error reconciling service account", "user", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if err := r.reconcileUserSelfRole(ctx, user); err != nil {
		logger.Error(err, "error reconciling self role", "user", req.NamespacedName)
		return ctrl.Result{}, err
	}

	if err := r.reconcileRoleBindings(ctx, user); err != nil {
		logger.Error(err, "error reconciling role bindings", "user", req.NamespacedName)
		return ctrl.Result{}, err

	}

	if err := r.Update(ctx, user); err != nil {
		logger.Error(err, "error updating user", "user", req.NamespacedName)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UserReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&marinacorev1.User{}).
		Owns(&corev1.ServiceAccount{}).
		Owns(&rbacv1.Role{}).
		Owns(&rbacv1.RoleBinding{}).
		Complete(r)
}
