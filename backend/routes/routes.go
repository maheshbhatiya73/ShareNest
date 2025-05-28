package routes

import (
    "github.com/gorilla/mux"
    "ShareNestbackend/controllers"
)

func FileRoutes(router *mux.Router) {
    api := router.PathPrefix("/api").Subrouter()
    api.HandleFunc("/upload", controllers.UploadFile).Methods("POST")
    api.HandleFunc("/download/{id}", controllers.DownloadFilePost).Methods("POST")
}
