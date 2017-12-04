package dbwork

import (
	"database/sql"
	//"github.com/ziutek/mymysql/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func Connect(login string) (error, *sql.DB) {
	db, err := sql.Open("mysql", login + ":" + Groups[login] + "@(localhost:3306)/study_projects")
	if err != nil {
		return err, nil
	}
	return nil, db
}

var Groups = map[string]string {
	"teacher" : "teacher",
	"student" : "student",
	"nonauth" : "password",
}