package repository

import (
	"study-gin-gorm/dao/mysql"
	"study-gin-gorm/model"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	DB *gorm.DB
}

func NewCategoryRepository() CategoryRepository {
	return CategoryRepository{DB: mysql.GetDB()}
}

func (c CategoryRepository) Create(name string) (*model.Category, error) {
	category := model.Category{
		Name: name,
	}

	if err := c.DB.Create(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (c CategoryRepository) Update(category model.Category, name string) (*model.Category, error) {
	if err := c.DB.Model(&category).Update("name", name).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (c CategoryRepository) SelectByID(id int) (*model.Category, error) {
	var category model.Category
	if err := c.DB.First(&category, id).Error; err != nil {
		return nil, err
	}

	return &category, nil
}

func (c CategoryRepository) DeleteByID(id int) error {
	if err := c.DB.Delete(model.Category{}, id).Error; err != nil {
		return err
	}

	return nil
}
