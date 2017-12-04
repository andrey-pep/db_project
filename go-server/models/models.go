package models

import (
	//"github.com/davecgh/go-spew/spew"
)

type Student struct {
	RecordBookNum	int    `db:"Record_book_num"`
	Birthday		string `db:"Birthday"`
	GroupName		string `db:"Group_name"`
	LastName		string `db:"Last_Name"`
}

type ResultTwo struct {
	ProjectId	     int
	Thema       	 string
	RecordBookNum    int
	LastName         string
}

type ResultOne struct {
	RecordBookNum    int
	SLastName        string
	TLastName        string
	PThema           string
	Mark             int
	GroupName        string
}

type User struct {
	Login	  string
	GroupName string
}
