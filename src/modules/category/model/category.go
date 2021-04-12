package model

import (
	"time"

	alias "github.com/iamJune20/dds/helper"
)

type Category struct {
	Code       string
	Name       string
	Desc       string
	Icon       string
	ManualCode string
	CreatedAt  time.Time
	UpdatedAt  time.Time
	DeleteAt   alias.NullTime
	Publish    string
}

type Categories []Category

func NewCategory() *Category {
	return &Category{
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func UpdateCategory() *Category {
	return &Category{
		UpdatedAt: time.Now(),
	}
}
