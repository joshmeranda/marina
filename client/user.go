package client

import (
	"context"

	"github.com/joshmeranda/marina/gateway/api/user"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ user.UserServiceClient = &Client{}

func (c *Client) CreateUser(ctx context.Context, in *user.UserCreateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_, err := c.userClient.CreateUser(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Client) GetUser(ctx context.Context, in *user.UserGetRequest, opts ...grpc.CallOption) (*user.User, error) {
	user, err := c.userClient.GetUser(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (c *Client) DeleteUser(ctx context.Context, in *user.UserDeleteRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_, err := c.userClient.DeleteUser(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Client) UpdateUser(ctx context.Context, in *user.UserUpdateRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	_, err := c.userClient.UpdateUser(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Client) ListUser(ctx context.Context, in *user.UserListRequest, opts ...grpc.CallOption) (*user.UserListResponse, error) {
	resp, err := c.userClient.ListUser(ctx, in, opts...)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
