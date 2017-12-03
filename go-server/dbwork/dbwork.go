package dbwork

import (
	"database/sql"
	//"github.com/ziutek/mymysql/mysql"
	_ "github.com/go-sql-driver/mysql"
)

func Connect() *sql.DB {
	db, err := sql.Open("mysql", "root:ghjybr7@(localhost:3306)/study_projects")
	if err != nil {
		panic(err)
	}
	return db
}