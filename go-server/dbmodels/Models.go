package requests

import (
	"database/sql"
	//"github.com/davecgh/go-spew/spew"
)

type Student struct {
	RecordBookNum	int    `db:"Record_book_num"`
	Birthday		string `db:"Birthday"`
	GroupName		string `db:"Group_name"`
	LastName		string `db:"Last_Name"`
}

func Req1(db *sql.DB) interface{} {
	rows, err := db.Query("SELECT * FROM student")
	if err != nil {
		return err
	}

	students := make([]Student, 0)

	for rows.Next() {
		var s Student
		err = rows.Scan(&s.RecordBookNum, &s.Birthday, &s.GroupName, &s.LastName)
		students = append(students, s)
		if err != nil {
			return err
		}

	}
	return students
}