package user

import (
	"context"
	apiv1 "pack/api/pkg/v1"
	"pack/internal/model"
	"pack/internal/service"
)

type cUser struct{}

func User() *cUser {
	return &cUser{}
}

func (c *cUser) Register(ctx context.Context, req *apiv1.RegisterUserReq) (res *apiv1.RegisterUserRes, err error) {
	userInfo := &model.User{
		Username: req.Username,
		Password: req.Password,
		Email:    req.Email,
		NickName: req.NickName,
	}
	if err = service.User().CreateUser(ctx, userInfo); err != nil {
		return
	}

	return nil, err
}
