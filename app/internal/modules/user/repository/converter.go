package repository

import (
	"small/internal/models"
	"small/internal/modules/user/repository/dao"
)

func (r *repository) userFromDaoToDomain(user *dao.User) *models.User {
	return &models.User{
		Id:    user.Id,
		Phone: user.Phone,
		Email: user.Email,
	}
}

func (r *repository) usersFromDaoToDomain(daoUsers []*dao.User) []*models.User {
	users := make([]*models.User, 0, len(daoUsers))

	for _, item := range daoUsers {
		users = append(users, r.userFromDaoToDomain(item))
	}
	return users
}
