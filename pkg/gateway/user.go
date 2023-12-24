package gateway

import (
	"context"

	marinav1 "github.com/joshmeranda/marina-operator/api/v1"
	"github.com/joshmeranda/marina/pkg/apis/user"
	"golang.org/x/crypto/bcrypt"
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
			Namespace: g.namespace,
		}, &role)

		if errors.IsNotFound(err) {
			return false
		}
	}

	return true
}

func (g *Gateway) CreateUser(ctx context.Context, req *user.UserCreateRequest) (*emptypb.Empty, error) {
	if !g.allRolesExist(req.User.Roles) {
		return nil, errors.NewBadRequest("one or more roles do not exist")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := marinav1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.User.Name,
			Namespace: g.namespace,
		},
		Spec: marinav1.UserSpec{
			Name:     req.User.Name, // todo: remove redundant field
			Password: hash,
			Roles:    req.User.Roles,
		},
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
			Namespace: g.namespace,
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
		Namespace: g.namespace,
	}, &user); err != nil {
		return nil, err
	}

	if req.User.Name != "" {
		user.Spec.Name = req.User.Name
	}

	if req.User.Password == nil {
		user.Spec.Password = req.User.Password
	}

	if req.User.Roles != nil {
		user.Spec.Roles = req.User.Roles
	}

	if err := g.kubeClient.Update(ctx, &user); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
