package repository

import (
	"database/sql"

	content "github.com/iamJune20/dds/src/modules/content/model"
	"github.com/iamJune20/dds/src/modules/search/model"
)

type searchRepository struct {
	db *sql.DB
}

func NewSearchRepository(db *sql.DB) *searchRepository {
	return &searchRepository{db}
}
func (r *searchRepository) FindAll(valueSearch string) (model.SearchAll, error) {
	query := `SELECT DISTINCT(app_code) as app_code FROM "search_all" WHERE "search_field" LIKE $1 `
	rows, err := r.db.Query(query, "%"+valueSearch+"%")

	if err != nil {
		return nil, err
	}

	var searchAll model.SearchAll
	for rows.Next() {
		var AppCode string
		var searchOne model.SearchOne
		err = rows.Scan(&AppCode)
		if err != nil {
			return nil, err
		}
		query5 := `SELECT DISTINCT(content_code) as content_code FROM "search_all" WHERE "search_field" LIKE $1 AND "app_code" = $2`
		var contents content.Contents
		rows5, err := r.db.Query(query5, "%"+valueSearch+"%", AppCode)
		if err != nil {
			return nil, err
		}
		for rows5.Next() {
			var ContentCode string
			var content content.Content
			err = rows5.Scan(&ContentCode)
			// fmt.Printf("Isi : %v", ContentCode)
			query6 := `SELECT * FROM "content" WHERE "content_code" = $1`
			err := r.db.QueryRow(query6, ContentCode).Scan(&content.Code, &content.Title, &content.Desc, &content.CategoryCode, &content.ManualCode, &content.CreatedAt, &content.UpdatedAt, &content.DeleteAt, &content.Publish)
			if err != nil {
				return nil, err
			}
			contents = append(contents, content)
		}
		query3 := `SELECT app_code, app_name, app_logo FROM "apps" WHERE "app_publish" = 'Yes' AND "delete_at" IS NULL AND "app_code" = $1`
		statement, err := r.db.Prepare(query3)
		var appCode string
		var appName string
		var appLogo string

		if err != nil {
			return nil, err
		}
		defer statement.Close()

		err = statement.QueryRow(AppCode).Scan(&appCode, &appName, &appLogo)

		if err != nil {
			return nil, err
		}
		searchOne.AppCode = appCode
		searchOne.AppName = appName
		searchOne.AppLogo = appLogo
		searchOne.CountContent = len(contents)
		searchOne.Content = contents
		searchAll = append(searchAll, searchOne)
	}
	return searchAll, err
}

func (r *searchRepository) FindOne(AppCode string, valueSearch string) (*model.SearchOne, error) {

	var searchOne model.SearchOne
	query5 := `SELECT DISTINCT(content_code) as content_code FROM "search_all" WHERE "search_field" LIKE $1 AND "app_code" = $2`
	var contents content.Contents
	rows5, err := r.db.Query(query5, "%"+valueSearch+"%", AppCode)
	if err != nil {
		return nil, err
	}
	for rows5.Next() {
		var ContentCode string
		var content content.Content
		err = rows5.Scan(&ContentCode)
		query6 := `SELECT * FROM "content" WHERE "content_code" = $1`
		err := r.db.QueryRow(query6, ContentCode).Scan(&content.Code, &content.Title, &content.Desc, &content.CategoryCode, &content.ManualCode, &content.CreatedAt, &content.UpdatedAt, &content.DeleteAt, &content.Publish)
		if err != nil {
			return nil, err
		}
		contents = append(contents, content)
	}
	query3 := `SELECT app_code, app_name, app_logo FROM "apps" WHERE "app_publish" = 'Yes' AND "delete_at" IS NULL AND "app_code" = $1`
	statement, err := r.db.Prepare(query3)
	var appCode string
	var appName string
	var appLogo string

	if err != nil {
		return nil, err
	}
	defer statement.Close()

	err = statement.QueryRow(AppCode).Scan(&appCode, &appName, &appLogo)

	if err != nil {
		return nil, err
	}
	searchOne.AppCode = appCode
	searchOne.AppName = appName
	searchOne.AppLogo = appLogo
	searchOne.CountContent = len(contents)
	searchOne.Content = contents

	// fmt.Printf("Isi : %v", searchOne)
	return &searchOne, err
}
