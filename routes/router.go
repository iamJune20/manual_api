package routes

import (
	"net/http"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
	"github.com/iamJune20/dds/src/controllers"
)

func middlewareOne(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")
		errBearer := validation.Validate(bearer,
			validation.Required,
		)
		if errBearer != nil {
			w.WriteHeader(http.StatusUnauthorized)
			controllers.ReturnError(w, 401, "SignIn terlebih dahulu")
			return
		}
		// fmt.Printf("Bearer = %v ", tokenString)
		var cek = controllers.ValidateToken(bearer)
		if cek == "" {
			next.ServeHTTP(w, r)
		} else {
			w.WriteHeader(http.StatusUnauthorized)
			controllers.ReturnError(w, 401, cek)
			return
		}

	})
}

func Router() *mux.Router {

	router := mux.NewRouter()
	api := router.PathPrefix("/api").Subrouter()
	api.Use(middlewareOne)

	router.HandleFunc("/public/signin", controllers.Signin).Methods("POST")
	router.HandleFunc("/public/refresh", controllers.Refresh).Methods("POST")

	// api.HandleFunc("/apps", controllers.GetApps).Methods("GET")
	router.HandleFunc("/public/apps", controllers.GetApps).Methods("GET")
	router.HandleFunc("/public/app/{code}", controllers.GetApp).Methods("GET")

	// with auth
	api.HandleFunc("/app", controllers.InsertApp).Methods("POST")
	api.HandleFunc("/app/{code}", controllers.UpdateApp).Methods("PUT")
	api.HandleFunc("/app/{code}", controllers.DeleteApp).Methods("DELETE")

	router.HandleFunc("/public/categories", controllers.GetCategories).Methods("GET")
	router.HandleFunc("/public/categoriesByManual/{code}", controllers.GetCategoriesByManualCode).Methods("GET")
	router.HandleFunc("/public/category/{code}", controllers.GetCategory).Methods("GET")

	// with auth
	api.HandleFunc("/category/{manual_code}", controllers.InsertCategory).Methods("POST")
	api.HandleFunc("/category/{code}", controllers.UpdateCategory).Methods("PUT")
	api.HandleFunc("/category/{code}", controllers.DeleteCategory).Methods("DELETE")

	router.HandleFunc("/public/manuals", controllers.GetManuals).Methods("GET")
	router.HandleFunc("/public/manualsByApp/{code}", controllers.GetManualsByAppCode).Methods("GET")
	router.HandleFunc("/public/manual/{code}", controllers.GetManual).Methods("GET")

	// with auth
	api.HandleFunc("/manual/{app_code}", controllers.InsertManual).Methods("POST")
	api.HandleFunc("/manual/{code}", controllers.UpdateManual).Methods("PUT")
	api.HandleFunc("/manual/{code}", controllers.DeleteManual).Methods("DELETE")

	router.HandleFunc("/public/contents", controllers.GetContents).Methods("GET")
	router.HandleFunc("/public/content/{code}", controllers.GetContent).Methods("GET")
	router.HandleFunc("/public/contentsIn/{manual_code}/{category_code}", controllers.GetContentInManualAndCategory).Methods("GET")

	// with auth
	api.HandleFunc("/content/{manual_code}/{category_code}", controllers.InsertContent).Methods("POST")
	api.HandleFunc("/content/{code}", controllers.UpdateContent).Methods("PUT")
	api.HandleFunc("/content/{code}", controllers.DeleteContent).Methods("DELETE")

	router.HandleFunc("/public/searchAll", controllers.GetSearchAll).Methods("POST")
	router.HandleFunc("/public/searchOne/{app_code}", controllers.GetSearchOne).Methods("POST")

	return router
}
