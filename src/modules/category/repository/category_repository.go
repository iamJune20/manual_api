package repository

import (
	"database/sql"
	"time"

	"github.com/iamJune20/dds/src/modules/category/model"
)

type categoryRepository struct {
	db *sql.DB
}

func NewCategoryRepository(db *sql.DB) *categoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) FindAll() (model.Categories, error) {

	query := `SELECT * FROM "category" WHERE "category_publish" = 'Yes' AND "delete_at" IS NULL`

	var categories model.Categories

	rows, err := r.db.Query(query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var category model.Category

		err = rows.Scan(&category.Code, &category.Name, &category.Desc, &category.Icon, &category.ManualCode, &category.CreatedAt, &category.UpdatedAt, &category.DeleteAt, &category.Publish)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}
func (r *categoryRepository) FindByManualCode(manual_code string) (model.Categories, error) {

	query := `SELECT * FROM "category" WHERE "category_publish" = 'Yes' AND "delete_at" IS NULL AND "manual_code" = $1`

	var categories model.Categories

	rows, err := r.db.Query(query, manual_code)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var category model.Category

		err = rows.Scan(&category.Code, &category.Name, &category.Desc, &category.Icon, &category.ManualCode, &category.CreatedAt, &category.UpdatedAt, &category.DeleteAt, &category.Publish)

		if err != nil {
			return nil, err
		}

		categories = append(categories, category)
	}

	return categories, nil
}
func (r *categoryRepository) FindByID(category_code string) (*model.Category, error) {

	query := `SELECT * FROM "category" WHERE "category_publish" = 'Yes' AND "delete_at" IS NULL AND "category_code" = $1`

	var category model.Category

	statement, err := r.db.Prepare(query)

	if err != nil {
		return nil, err
	}

	defer statement.Close()

	err = statement.QueryRow(category_code).Scan(&category.Code, &category.Name, &category.Desc, &category.Icon, &category.ManualCode, &category.CreatedAt, &category.UpdatedAt, &category.DeleteAt, &category.Publish)

	if err != nil {
		return nil, err
	}

	return &category, nil
}

func (r *categoryRepository) Save(category *model.Category) (string, error) {

	query := `INSERT INTO category (category_name, category_desc,category_icon,manual_code,create_at,update_at, category_publish) VALUES ($1, $2, $3, $4, $5, $6,'Yes') RETURNING category_code`

	var Code string
	err := r.db.QueryRow(query, category.Name, category.Desc, category.Icon, category.ManualCode, category.CreatedAt, category.UpdatedAt).Scan(&Code)
	// fmt.Printf("Error : %v", err)
	if err != nil {
		return "Data gagal disimpan", err
	}

	return Code, err
}

func (r *categoryRepository) Update(category_code string, category *model.Category) (string, error) {
	query := `UPDATE category SET category_name=$1, category_desc=$2, category_icon=$3, manual_code=$4, update_at=$5 WHERE category_code=$6 AND "category_publish" = 'Yes' AND "delete_at" IS NULL`

	statement, err := r.db.Prepare(query)

	if err != nil {
		return "Data gagal di ubah", err
	}

	defer statement.Close()

	_, err = statement.Exec(category.Name, category.Desc, category.Icon, category.ManualCode, category.UpdatedAt, category_code)

	if err != nil {
		return "Data gagal di ubah", err
	}

	return "Data berhasil diubah", err
}

func (r *categoryRepository) Delete(category_code string) (string, error) {
	query := `UPDATE category SET delete_at=$1, category_publish='No' WHERE category_code=$2 AND "category_publish" = 'Yes' AND "delete_at" IS NULL`

	statement, err := r.db.Prepare(query)

	if err != nil {
		return "Data gagal di hapus", err
	}

	defer statement.Close()

	_, err = statement.Exec(time.Now(), category_code)

	if err != nil {
		return "Data gagal di hapus", err
	}

	return "Data berhasil dihapus", err
}
