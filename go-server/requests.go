package main

import (
	//"github.com/davecgh/go-spew/spew"
	"net/http"
	"./models"
)

func (c *MainController) Req1(r *http.Request) (error, interface{}) {
	var Select = "Select project.Record_book_num, student.last_name, teachers.Last_name, project.thema, protocol_string.mark, student.Group_Name "
	var From = `From protocol JOIN protocol_string USING(p_id)
		JOIN project USING(pr_id)
		JOIN teachers USING(t_id)
		JOIN student ON project.Record_book_num = student.Record_book_num `
	var Where = `WHere YEAR(project.finish_date) = ? AND student.Group_name = ? `
	var Order = `Order By Record_book_num `
	rows, err := c.DataBase.Query(Select + From + Where + Order, r.URL.Query().Get("c_year"), r.URL.Query().Get("g_index"))
	if err != nil {
		return err, nil
	}
	results := make([]*models.ResultOne, 0)
	for rows.Next() {
		res := new(models.ResultOne)
		err := rows.Scan(&res.RecordBookNum, &res.SLastName, &res.TLastName, &res.PThema, &res.Mark, &res.GroupName)
		if err != nil {
			return err, nil
		}
		results = append(results, res)
	}
	return nil, results
}

func (c *MainController) Req2(r *http.Request) (error, interface{}) {
	var Select = `SELECT pr_id, thema, Record_book_num, Last_name
				  FROM project JOIN student USING(Record_book_num)
				  WHERE subject = ?`
	rows, err := c.DataBase.Query(Select, r.URL.Query().Get("c_subj"))
	if err != nil {
		return err, nil
	}
	results := make([]*models.ResultTwo, 0)
	for rows.Next() {
		result := new(models.ResultTwo)
		err := rows.Scan(&result.ProjectId, &result.Thema, &result.RecordBookNum, &result.LastName)
		if err != nil {
			return err, nil
		}
		results = append(results, result)
	}
	return nil, results
}