package repository

import (
	"github.com/iamJune20/dds/src/modules/category/model"
)

type CategoryRepository interface {
	Save(*model.Category) (string, error)
	Update(string, *model.Category) (string, error)
	Delete(string) (string, error)
	FindByID(string) (*model.Category, error)
	FindByManualCode(string) (*model.Categories, error)
	FindAll() (model.Categories, error)
}
