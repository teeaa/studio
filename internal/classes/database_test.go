package classes

import (
	"testing"
	"time"

	mocket "github.com/selvatico/go-mocket"
)

func TestGetClassById(t *testing.T) {
	setup()
	commonReply := []map[string]interface{}{{
		"id":         1,
		"name":       "Class #1",
		"start_date": time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
		"end_date":   time.Date(2019, 8, 31, 0, 0, 0, 0, time.UTC),
		"capacity":   20,
	}}
	mocket.Catcher.Reset().NewMock().WithQuery(`SELECT * FROM "classes"  WHERE ("classes"."id" = 1) ORDER BY "classes"."id" ASC LIMIT 1`).WithReply(commonReply)

	class, err := GetClassByID(1)
	if err != nil {
		t.Error("Error getting class by id")
	}

	compare := Class{
		1,
		"Class #1",
		time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 8, 31, 0, 0, 0, 0, time.UTC),
		20,
	}
	if compare != class {
		t.Error("Retrieved class data didn't match expectations:", compare, class)
	}
}

func TestGetClassFromReq(t *testing.T) {
	w, r, _ := makeRequest(nil, map[string]string{"id": "1"})

	// var class Class
	class, err := db.getClassFromReq(w, r)
	if err != nil {
		t.Error("Error getting class by request vars")
	}
	compare := &Class{
		1,
		"Class #1",
		time.Date(2019, 6, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2019, 8, 31, 0, 0, 0, 0, time.UTC),
		20,
	}
	if *compare != *class {
		t.Error("Retrieved class data didn't match expectations:", *compare, *class)
	}
}
