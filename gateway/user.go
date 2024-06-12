package gateway

import (
	"context"
	"fmt"

	marinav1 "github.com/joshmeranda/marina/api/v1"
	"github.com/joshmeranda/marina/gateway/api/user"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/protobuf/types/known/emptypb"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var _ user.UserServiceServer = &Gateway{}

const (
	DefaultRoleName = "marina-user"
)

func (g *Gateway) allRolesExist(kubeClient client.Client, roles []string) bool {
	for _, roleName := range roles {
		var role rbacv1.Role
		err := kubeClient.Get(context.Background(), types.NamespacedName{
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
	kubeClient, err := g.clientFromContext(ctx, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to impersonate user: %w", err)
	}

	if !g.allRolesExist(kubeClient, req.User.Roles) {
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

	if err := kubeClient.Create(ctx, &user); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (g *Gateway) GetUser(ctx context.Context, req *user.UserGetRequest) (*user.User, error) {
	kubeClient, err := g.clientFromContext(ctx, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to impersonate user: %w", err)
	}

	var u marinav1.User
	if err := kubeClient.Get(ctx, types.NamespacedName{
		Name:      req.Name,
		Namespace: g.namespace,
	}, &u); err != nil {
		return nil, err
	}

	return &user.User{
		Name:     u.Name,
		Password: []byte{},
		Roles:    u.Spec.Roles,
	}, nil
}

func (g *Gateway) DeleteUser(ctx context.Context, req *user.UserDeleteRequest) (*emptypb.Empty, error) {
	user := marinav1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name:      req.Name,
			Namespace: g.namespace,
		},
	}

	kubeClient, err := g.clientFromContext(ctx, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to impersonate user: %w", err)
	}

	if err := kubeClient.Delete(ctx, &user); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (g *Gateway) UpdateUser(ctx context.Context, req *user.UserUpdateRequest) (*emptypb.Empty, error) {
	kubeClient, err := g.clientFromContext(ctx, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to impersonate user: %w", err)
	}

	if !g.allRolesExist(kubeClient, req.User.Roles) {
		return nil, errors.NewBadRequest("one or more roles do not exist")
	}

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

	if req.User.Password != nil {
		hash, err := bcrypt.GenerateFromPassword([]byte(req.User.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		user.Spec.Password = hash
	}

	if req.User.Roles != nil {
		user.Spec.Roles = req.User.Roles
	}

	if err := kubeClient.Update(ctx, &user); err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (g *Gateway) ListUser(ctx context.Context, req *user.UserListRequest) (*user.UserListResponse, error) {
	kubeClient, err := g.clientFromContext(ctx, client.Options{})
	if err != nil {
		return nil, fmt.Errorf("failed to impersonate user: %w", err)
	}

	var list marinav1.UserList
	if err := kubeClient.List(ctx, &list, client.InNamespace(g.namespace)); err != nil {
		return nil, err
	}

	response := &user.UserListResponse{
		Users: make([]*user.User, 0, len(list.Items)),
	}

	for _, foundUser := range list.Items {
		matches, err := req.Query.Matches(&foundUser)
		if err != nil {
			return nil, fmt.Errorf("failed to apply qurery to user: %w", err)
		}

		if matches {
			response.Users = append(response.Users, &user.User{
				Name:  foundUser.Name,
				Roles: foundUser.Spec.Roles,
			})
		}
	}

	return response, nil
}
