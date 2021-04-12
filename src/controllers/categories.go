package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	"github.com/iamJune20/dds/config"
	"github.com/iamJune20/dds/src/modules/category/repository"
	appRe "github.com/iamJune20/dds/src/modules/manual/repository"

	"github.com/iamJune20/dds/src/modules/category/model"
	_ "github.com/lib/pq"
)

type ResponseGetCategories struct {
	Status  int              `json:"status"`
	Message string           `json:"message"`
	Data    model.Categories `json:"data"`
}

type ResponseGetCategory struct {
	Status  int             `json:"status"`
	Message string          `json:"message"`
	Data    *model.Category `json:"data"`
}

func GetCategories(w http.ResponseWriter, r *http.Request) {
	db := config.CreateConnection()
	categoryRepository := repository.NewCategoryRepository(db)

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	out, err := categoryRepository.FindAll()

	if err != nil {
		log.Fatalf("Error: %v", err)
		return
	}

	var response ResponseGetCategories
	response.Status = 200
	response.Message = "Success"
	response.Data = out

	json.NewEncoder(w).Encode(response)
}

func GetCategory(w http.ResponseWriter, r *http.Request) {
	db := config.CreateConnection()
	categoryRepository := repository.NewCategoryRepository(db)

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	out, err := categoryRepository.FindByID(params["code"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		var response ResponseError
		response.Status = 404
		response.Message = "Code Not Found"
		json.NewEncoder(w).Encode(response)
	} else {
		var response ResponseGetCategory
		response.Status = 200
		response.Message = "Success"
		response.Data = out
		json.NewEncoder(w).Encode(response)
	}
}

func GetCategoriesByManualCode(w http.ResponseWriter, r *http.Request) {
	db := config.CreateConnection()
	categoryRepository := repository.NewCategoryRepository(db)

	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	params := mux.Vars(r)

	out, err := categoryRepository.FindByManualCode(params["code"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		var response ResponseError
		response.Status = 404
		response.Message = "Code Not Found"
		json.NewEncoder(w).Encode(response)
	} else {
		var response ResponseGetCategories
		response.Status = 200
		response.Message = "Success"
		response.Data = out
		json.NewEncoder(w).Encode(response)
	}
}

func InsertCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := config.CreateConnection()
	categoryRepository := repository.NewCategoryRepository(db)

	var category *model.Category
	category = model.NewCategory()
	category.Name = r.FormValue("Name")
	errName := validation.Validate(category.Name,
		validation.Required,
	)
	if errName != nil {
		w.WriteHeader(http.StatusBadRequest)
		ReturnError(w, 400, "Nama kategory harus diisi")
		return
	}
	category.Desc = r.FormValue("Desc")
	errDesc := validation.Validate(category.Desc,
		validation.Required,
	)
	if errDesc != nil {
		w.WriteHeader(http.StatusBadRequest)
		ReturnError(w, 400, "Deskripsi kategory harus diisi")
		return
	}

	category.Icon = r.FormValue("Icon")
	errIcon := validation.Validate(category.Icon,
		validation.Required,
	)
	if errIcon != nil {
		w.WriteHeader(http.StatusBadRequest)
		ReturnError(w, 400, "Icon kategory harus diisi")
		return
	}

	params := mux.Vars(r)
	appRepository := appRe.NewManualRepository(db)
	app, err := appRepository.FindByID(params["manual_code"])

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		var response ResponseError
		response.Status = 404
		response.Message = "Code Not Found"
		json.NewEncoder(w).Encode(response)
		return
	} else {
		if app == nil {
			w.WriteHeader(http.StatusNotFound)
			var response ResponseError
			response.Status = 404
			response.Message = "Code Not Found"
			json.NewEncoder(w).Encode(response)
			return
		} else {
			category.ManualCode = params["manual_code"]
		}
	}
	// fmt.Printf("AppCode : %v", category.AppCode)
	out, err := categoryRepository.Save(category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ReturnError(w, 400, out)
		return
	} else {
		res := response{
			Status:  201,
			Message: "Success",
			Code:    out,
		}
		json.NewEncoder(w).Encode(res)
	}
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := config.CreateConnection()
	categoryRepository := repository.NewCategoryRepository(db)

	params := mux.Vars(r)

	var category *model.Category
	category = model.NewCategory()
	category.Name = r.FormValue("Name")
	errName := validation.Validate(category.Name,
		validation.Required,
	)
	if errName != nil {
		w.WriteHeader(http.StatusBadRequest)
		ReturnError(w, 400, "Nama kategory harus diisi")
		return
	}
	category.Desc = r.FormValue("Desc")
	errDesc := validation.Validate(category.Desc,
		validation.Required,
	)
	if errDesc != nil {
		w.WriteHeader(http.StatusBadRequest)
		ReturnError(w, 400, "Deskripsi kategory harus diisi")
		return
	}

	category.Icon = r.FormValue("Icon")
	errIcon := validation.Validate(category.Icon,
		validation.Required,
	)
	if errIcon != nil {
		w.WriteHeader(http.StatusBadRequest)
		ReturnError(w, 400, "Icon kategory harus diisi")
		return
	}

	category.ManualCode = r.FormValue("ManualCode")
	errAppCode := validation.Validate(category.ManualCode,
		validation.Required,
	)
	if errAppCode != nil {
		w.WriteHeader(http.StatusBadRequest)
		ReturnError(w, 400, "ManualCode harus diisi")
		return
	}
	appRepository := appRe.NewManualRepository(db)
	app, err := appRepository.FindByID(category.ManualCode)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		var response ResponseError
		response.Status = 404
		response.Message = "Code Not Found"
		json.NewEncoder(w).Encode(response)
		return
	} else {
		if app == nil {
			w.WriteHeader(http.StatusNotFound)
			var response ResponseError
			response.Status = 404
			response.Message = "Code Not Found"
			json.NewEncoder(w).Encode(response)
		}
	}

	out, err := categoryRepository.Update(params["code"], category)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ReturnError(w, 400, out)
		return
	} else {
		res := response{
			Status:  200,
			Message: "Success",
			Code:    out,
		}
		json.NewEncoder(w).Encode(res)
	}
}

func DeleteCategory(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Context-Type", "application/x-www-form-urlencoded")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	db := config.CreateConnection()
	categoryRepository := repository.NewCategoryRepository(db)

	params := mux.Vars(r)
	errCode := validation.Validate(params["code"],
		validation.Required,
	)
	if errCode != nil {
		w.WriteHeader(http.StatusNotFound)
		ReturnError(w, 404, "Kode tidak ditemukan")
		return
	}
	out, err := categoryRepository.Delete(params["code"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		ReturnError(w, 400, "Gagal menghapus data")
		return
	} else {
		res := response{
			Status:  200,
			Message: out,
			Code:    params["code"],
		}
		json.NewEncoder(w).Encode(res)
	}
}
