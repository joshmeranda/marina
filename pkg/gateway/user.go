package gateway

import (
	"context"

	marinav1 "github.com/joshmeranda/marina-operator/api/v1"
	"github.com/joshmeranda/marina/pkg/apis/user"
	"google.golang.org/protobuf/types/known/emptypb"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
)

var _ user.UserServiceServer = &Gateway{}

func (g *Gateway) allRolesExist(roles []string) bool {
	for _, roleName := range roles {
		var role rbacv1.Role
		err := g.kubeClient.Get(context.Background(), types.NamespacedName{
			Name:      roleName,
			Namespace: "marina-system",
		}, &role)

		if errors.IsNotFound(err) {
			return false
		}
	}

	return true
}

func (g *Gateway) CreateUser(ctx context.Context, req *user.UserCreateRequest) (*emptypb.Empty, error) {
	user := marinav1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.User.Name,
			Namespace: "marina-system",
		},
		Spec: marinav1.UserSpec{
			Name:  req.User.Name,
			Roles: req.User.Roles,
		},
	}

	if !g.allRolesExist(req.User.Roles) {
		return nil, errors.NewBadRequest("one or more roles do not exist")
	}

	if err := g.kubeClient.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (g *Gateway) DeleteUser(ctx context.Context, req *user.UserDeleteRequest) (*emptypb.Empty, error) {
	user := marinav1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: "marina-system",
		},
	}

	if err := g.kubeClient.Delete(ctx, &user); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (g *Gateway) UpdateUser(ctx context.Context, req *user.UserUpdateRequest) (*emptypb.Empty, error) {
	var user marinav1.User
	if err := g.kubeClient.Get(ctx, types.NamespacedName{
		Name:      req.User.Name,
		Namespace: "marina-system",
	}, &user); err != nil {
		return nil, err
	}

	if req.User.Name != "" {
		user.Spec.Name = req.User.Name
	}

	if req.User.Roles != nil {
		user.Spec.Roles = req.User.Roles
	}

	if err := g.kubeClient.Update(ctx, &user); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
