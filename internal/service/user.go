package service

import (
	"context"
	"errors"
	"pack/internal/dao"
	"pack/internal/model"
	"pack/internal/model/entity"
	"pack/utility/util"
)

type sUser struct{}

func User() *sUser {
	return &sUser{}
}

func (s *sUser) CreateUser(ctx context.Context, user *model.User) (err error) {
	var userInfo *entity.User
	// 查看用户是否存在
	if err, _ = s.GetUserInfo(ctx, user.Username); err == nil {
		return errors.New("user already exists")
	}

	// 生成salt
	salt, err := util.GenerateSalt()
	if err != nil {
		return
	}

	// 对password进行加密
	password := util.HashPassword(user.Password, salt)

	userInfo = &entity.User{
		Username: user.Username,
		Password: password,
		Salt:     salt,
		Nickname: user.NickName,
		Email:    user.Email,
	}

	if _, err = dao.User.Ctx(ctx).Data(userInfo).Insert(); err != nil {
		return
	}

	return
}

func (s *sUser) GetUserInfo(ctx context.Context, username string) (err error, user interface{}) {
	userInfo := entity.User{}
	err = dao.User.Ctx(ctx).Where("username=?", username).Scan(&userInfo)
	if err != nil {
		return
	}
	return err, userInfo
}
