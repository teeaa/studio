package classes

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	log "github.com/sirupsen/logrus"
	"github.com/teeaa/studio/internal/helpers"
)

// Database wrapper
type Database struct {
	*gorm.DB
}

var db Database

// SetupExternally to set db from imports
func SetupExternally(database Database) {
	db = database
}

// GetClassByID get class from class db
func GetClassByID(classID uint64) (Class, error) {
	var class Class
	err := db.First(&class, classID).Error

	if err != nil {
		log.Error("Error retrieving class from db:", err)
		return Class{}, err
	}

	return class, nil
}

// Get class from database by id in request and handle error situations
func (db *Database) getClassFromReq(w http.ResponseWriter, r *http.Request) (*Class, error) {
	var class Class
	vars := mux.Vars(r)
	classID, err := strconv.Atoi(vars["id"])
	if err != nil {
		log.Warnf("Requested class id (%s) is not an integer: %s", vars["id"], err)
		helpers.ResponseJSON(w, http.StatusBadRequest, "Invalid class ID")
		return nil, err
	}

	err = db.First(&class, classID).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			log.Warnf("Requested class by id %d does not exist", classID)
			helpers.ResponseJSON(w, http.StatusNotFound, "Class does not exist")
		} else {
			log.Error("Error fetching class from db: ", err)
			helpers.ResponseJSON(w, http.StatusInternalServerError, "Something went wrong")
		}
		return nil, err
	}
	return &class, nil
}
