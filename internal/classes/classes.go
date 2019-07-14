package classes

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/teeaa/studio/internal/helpers"
)

func getClasses(w http.ResponseWriter, r *http.Request) {
	var classes []Class
	err := db.Find(&classes).Error

	if err != nil {
		log.Error("Error fetching classes from db: ", err)
		helpers.ResponseJSON(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	json.NewEncoder(w).Encode(&classes)
}

func addClass(w http.ResponseWriter, r *http.Request) {
	var class Class
	err := json.NewDecoder(r.Body).Decode(&class)
	if err != nil {
		log.Warn("Error parsing JSON when creating new class: ", err)
		helpers.ResponseJSON(w, http.StatusBadRequest, "Invalid request body for class")
		return
	}

	err = db.Create(&class).Error
	if err != nil {
		log.Error("Error inserting class to db: ", err)
		helpers.ResponseJSON(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&class)
}

func getClass(w http.ResponseWriter, r *http.Request) {
	class, err := db.getClassFromReq(w, r)
	if err != nil {
		return
	}

	json.NewEncoder(w).Encode(&class)
}

func updateClass(w http.ResponseWriter, r *http.Request) {
	class, err := db.getClassFromReq(w, r)
	if err != nil {
		return
	}

	err = json.NewDecoder(r.Body).Decode(&class)
	if err != nil {
		log.Warn("Error parsing JSON when creating new class: ", err)
		helpers.ResponseJSON(w, http.StatusBadRequest, "Invalid request body for class")
		return
	}

	err = db.Save(&class).Error
	if err != nil {
		log.Error("Error saving class to db: ", err)
		helpers.ResponseJSON(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	json.NewEncoder(w).Encode(&class)
}
func deleteClass(w http.ResponseWriter, r *http.Request) {
	class, err := db.getClassFromReq(w, r)
	if err != nil {
		return
	}

	err = db.Delete(&class).Error
	if err != nil {
		log.Error("Error deleting class from db: ", err)
		helpers.ResponseJSON(w, http.StatusInternalServerError, "Something went wrong")
		return
	}

	helpers.ResponseJSON(w, 200, "Class removed")
}

// Routes set routes for /classes
func Routes(gormDB *gorm.DB, router *mux.Router) {
	db = Database{gormDB}
	db.AutoMigrate(&Class{})
	router.HandleFunc("", getClasses).Methods("GET")
	router.HandleFunc("", addClass).Methods("POST")
	router.HandleFunc("/{id}", getClass).Methods("GET")
	router.HandleFunc("/{id}", updateClass).Methods("PUT")
	router.HandleFunc("/{id}", deleteClass).Methods("DELETE")
}
