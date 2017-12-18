package models

import (
	//"github.com/davecgh/go-spew/spew"
)

type Student struct {
	RecordBookNum	int    `db:"Record_book_num"`
	Birthday		string `db:"Birthday"`
	GroupName		string `db:"Group_name"`
	LastName		string `db:"Last_Name"`
	Subject         string
	Thema           string 
}

type ResultTwo struct {
	ProjectId	     int
	Thema       	 string
	RecordBookNum    int
	LastName         string
}

type MarksUpdate struct {
	RecordBookNum string	`json:"record_book_num"`
	Birthdate   string		`json:"birthdate"`
	GroupName   string		`json:"group_name"`
	LastName    string		`json:"last_name"`
	Subject 	string
	Theme 		string
	Mark 		string
}

type ResultOne struct {
	RecordBookNum    int
	SLastName        string
	TLastName        string
	PThema           string
	Mark             *int
	GroupName        string
}

type ResultThree struct {
	TId			int
	LastName 	string
	Birthdate	string
	PulpitNum	int
	StWorkTime	string
}

type ResultFour struct {
	ResultThree
}

type ResultSix struct {
	ResultThree
}

type ResultFive struct {
	ResultThree
}

type User struct {
	Login	  string
	GroupName string
}

type Otchet struct {
	OId			int
	TId			int
	SubjectName	string
	OGroup		string
	AvgMark		float64
	OYear		int
}
