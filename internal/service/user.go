package service

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"pack/internal/dao"
	"pack/internal/model/entity"
)

type sUser struct{}

func User() *sUser {
	return &sUser{}
}

func (s *sUser) CreateUser(ctx context.Context, user entity.User) (err error) {
	// 这里需要对password进行加密解密
	if _, err = dao.User.Ctx(ctx).Data(user).Insert(); err != nil {
		return
	}

	return
}

func (s *sUser) GetUserInfo(ctx context.Context, username string) (err error, user string) {
	ret, err := dao.User.Ctx(ctx).Where("username=?", username).One()
	if err != nil {
		return
	}
	return err, ret.Json()
}

func (s *sUser) UpdateUser(ctx context.Context, user entity.User) (err error) {
	if ret, err := dao.User.Ctx(ctx).Data(user).Where(user.Username).Update(user); err != nil {
		return
	} else {
		g.Log().Debugf(ctx, "update user success: %s", ret)
	}

	return
}
