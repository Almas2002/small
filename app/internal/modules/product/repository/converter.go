package repository

import (
	"small/internal/models"
	"small/internal/modules/product/repository/dao"
)

func (r *repository) productFromDaoToDomain(dao *dao.Product) *models.Product {
	return &models.Product{
		Id:    dao.Id,
		Title: dao.Title,
		Price: dao.Price,
	}
}
