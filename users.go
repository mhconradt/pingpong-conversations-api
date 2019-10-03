package main

import (
	"context"
	"fmt"
	status "github.com/mhconradt/grpc-statuses"
	proto "github.com/mhconradt/pingpong-conversations-api/proto"
	"github.com/mhconradt/proto/user"
)

func (c *Conversations) CreateUser(ctx context.Context, req *proto.CreateUserRequest, res *proto.CreateUserResponse) error {
	if u, err := user.DefaultCreateUser(ctx, req.User, pool); err != nil {
		res.Status = status.BadRequest(err.Error())
		fmt.Println(err)
	} else {
		res.User = u
		res.Status = status.Success()
	}
	return nil
}

func (c *Conversations) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest, res *proto.UpdateUserResponse) error {
	if u, err := user.DefaultPatchUser(ctx, req.User, req.UpdateMask, pool); err != nil {
		res.Status = status.BadRequest(err.Error())
		fmt.Println(err)
	} else {
		res.User = u
		res.Status = status.Success()
	}
	return nil
}

func (c *Conversations) GetUser(ctx context.Context, req *proto.GetUserRequest, res *proto.GetUserResponse) error {
	q := &user.User{Id: req.Id}
	u, err := user.DefaultReadUser(ctx, q, pool)
	if err != nil {
		res.Status = status.NotFound("user with id: %v not found", req.Id)
		fmt.Println(err)
		return nil
	}
	if u, err = user.DefaultApplyFieldMaskUser(ctx, new(user.User), u, req.GetMask, "", pool); err != nil {
		res.Status = status.InternalServerError("failed to apply field mask to result")
	} else {
		res.User = u
		res.Status = status.Success()
	}
	return nil
}

func (c *Conversations) DeactivateUser(ctx context.Context, req *proto.DeactivateUserRequest, res *proto.DeactivateUserResponse) error {
	orm, err := (&user.User{Id: req.Id}).ToORM(ctx)
	if err != nil {
		res.Status = status.BadRequest(err.Error())
		fmt.Println(err)
	}
	if err := pool.Model(orm).Update("status", 0).Error; err != nil {
		res.Status = status.InternalServerError("failed to deactivate user")
		fmt.Println(err)
	} else {
		res.User = &user.User{}
		res.Status = status.Success()
	}
	return nil
}
