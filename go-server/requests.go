package main

import (
	"github.com/davecgh/go-spew/spew"
	"net/http"
)

func (c *MainController) Req1(r *http.Request) interface{} {
	var Select = "Select project.Record_book_num, teachers.Last_name, student.last_name, project.thema, protocol_string.mark "
	var From = `From protocol JOIN protocol_string USING(p_id)
		JOIN project USING(pr_id)
		JOIN teachers USING(t_id)
		JOIN student ON project.Record_book_num = student.Record_book_num `
	var Where = `WHere YEAR(project.finish_date) = ? `
	var Order = `Order By Record_book_num `
	rows, err := c.DataBase.Query(Select + From + Where + Order, r.URL.Query().Get("c_year"))
	spew.Dump(c.DataBase.Query(Select + From + Where + Order, r.URL.Query().Get("c_year")))
	spew.Dump(r.Form)
	spew.Dump(rows)
	if err != nil {
		return err
	}
	return rows
}