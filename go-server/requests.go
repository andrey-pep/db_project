package main

import (
	//"github.com/davecgh/go-spew/spew"
	"net/http"
	"./models"
	"crypto/md5"
	"encoding/hex"
	"strconv"
)


func (c *MainController) MarksInput(arr []models.MarksUpdate) (int, error) {
	update := `update protocol_string set mark = ? 
	WHERE thema = ? and Record_book_num = ?`
	affected := 0
	for _, v := range arr {
		mark, _ := strconv.Atoi(v.Mark)
		res, err := c.DataBase.Exec(update, mark, v.Theme, v.RecordBookNum)
		if (err != nil) {
			return affected, err
		}
		_, err = res.RowsAffected()
		if (err != nil) {
			return affected, err
		}
		affected++
	}
	return affected, nil;
}

func (c *MainController) Otchet() (error, []*models.Otchet) {
	Select := `SELECT o_id, t_id, subject_name, o_group, avg_mark, o_year FROM avg_marks;`
	rows, err := c.DataBase.Query(Select)
	if err != nil {
		return err, nil
	}
	results := make([]*models.Otchet, 0)
	for rows.Next() {
		res := new(models.Otchet)
		err := rows.Scan(&res.OId, &res.TId, &res.SubjectName, &res.OGroup, &res.AvgMark, &res.OYear)
		if err != nil {
			return err, nil
		}
		results = append(results, res)
	}
	return nil, results	
}

func (c *MainController) SelectStudents(r *http.Request) (error, []*models.Student) {
	Select := `Select student.Record_book_num, Birthday, Group_name, Last_Name, subject, project.thema FROM student JOIN protocol_string USING(Record_book_num)
		JOIN project USING(pr_id)
		Where Group_Name = ? and mark IS NULL`
	rows, err := c.DataBase.Query(Select, r.URL.Query().Get("g_index"))
	if err != nil {
		return err, nil
	}
	results := make([]*models.Student, 0)
	for rows.Next() {
		res := new(models.Student)
		err := rows.Scan(&res.RecordBookNum, &res.Birthday, &res.GroupName, &res.LastName, &res.Subject, &res.Thema)
		if err != nil {
			return err, nil
		}
		results = append(results, res)
	}
	return nil, results
}

func (c *MainController) MakeOtchet(r *http.Request) error {
	query := `CALL Otchet(` + r.URL.Query().Get("YP") + `, "` + r.URL.Query().Get("SP") + `")`

	res, err := c.DataBase.Exec(query)

	return err
}

func (c *MainController) Req3(r *http.Request) (error, interface{}) {
	Select := `Select * FROM teachers WHERE Birthdate = (SELECT MAX(Birthdate) FROM teachers)`
	rows, err := c.DataBase.Query(Select)
	if err != nil {
		return err, nil
	}
	results := make([]*models.ResultThree, 0)
	for rows.Next() {
		res := new(models.ResultThree)
		err := rows.Scan(&res.TId, &res.LastName, &res.Birthdate, &res.PulpitNum, &res.StWorkTime)
		if err != nil {
			return err, nil
		}
		results = append(results, res)
	}
	return nil, results
}

func (c *MainController) Req4(r *http.Request) (error, interface{}) {
	Select := `Select t_id, last_name, Birthdate, Pulpit_num, St_work_date`
	Where := ` Where pr_id IS NULL`
	From := ` From teachers LEFT JOIN project USING(t_id)`
	rows, err := c.DataBase.Query(Select + From + Where)
	if err != nil {
		return err, nil
	}
	results := make([]*models.ResultFour, 0)
	for rows.Next() {
		res := new(models.ResultFour)
		err := rows.Scan(&res.TId, &res.LastName, &res.Birthdate, &res.PulpitNum, &res.StWorkTime)
		if err != nil {
			return err, nil
		}
		results = append(results, res)
	}
	return nil, results
}

func (c *MainController) Req5(r *http.Request) (error, interface{}) {
	Select := `Select teachers.* `
	From := `FROM teachers `
	Where := `WHERE t_id NOT IN (
    SELECT t_id
    FROM teachers LEFT JOIN project USING(t_id)
     WHERE YEAR(finish_date)= ?)`
	rows, err := c.DataBase.Query(Select + From + Where, r.URL.Query().Get("g_index"))
	if err != nil {
		return err, nil
	}
	results := make([]*models.ResultFive, 0)
	for rows.Next() {
		res := new(models.ResultFive)
		err := rows.Scan(&res.TId, &res.LastName, &res.Birthdate, &res.PulpitNum, &res.StWorkTime)
		if err != nil {
			return err, nil
		}
		results = append(results, res)
	}
	return nil, results
}

func (c *MainController) Req6(r *http.Request) (error, interface{}) {
	Select := `Select teachers.* `
	From := `FROM (
		SELECT COUNT(t_id) as cnt, t_id
		FROM teachers LEFT JOIN project USING(t_id)
		WHERE Record_book_num IS NOT NULL
		GROUP BY t_id) t_c JOIN teachers USING(t_id)`
	Where := `WHERE cnt = (
		SELECT MAX(cnt) as m_cnt
		FROM (
		SELECT COUNT(t_id) as cnt
		FROM teachers LEFT JOIN project USING(t_id)
		WHERE Record_book_num IS NOT NULL
		GROUP BY t_id) t_count)`
	rows, err := c.DataBase.Query(Select + From + Where)
	if err != nil {
		return err, nil
	}
	results := make([]*models.ResultSix, 0)
	for rows.Next() {
		res := new(models.ResultSix)
		err := rows.Scan(&res.TId, &res.LastName, &res.Birthdate, &res.PulpitNum, &res.StWorkTime)
		if err != nil {
			return err, nil
		}
		results = append(results, res)
	}
	return nil, results
}


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

func (c *MainController) SelectUser(r *http.Request) error {
	hasher := md5.New()
	hasher.Write([]byte(r.URL.Query().Get("password")))
	hashedPass := hex.EncodeToString(hasher.Sum(nil))
	rows, err := c.DataBase.Query("Select login, group_type, password from users where login = ? AND password = ?", r.URL.Query().Get("login"), hashedPass)
	if err != nil {
		return err
	}
	for rows.Next() {
		err = rows.Scan(&c.Login, &c.UsrGroup, &c.UsrPass)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *MainController) CheckIfUserExists() (error, bool) {
	rows, err := c.DataBase.Query("Select login, group_type from users where login = ? AND password = ?", c.Login, c.UsrPass )
	if err != nil {
		return err, false
	}
	r := new(models.User)
	for rows.Next() {
		err = rows.Scan(&r.Login, &c.UsrGroup)
		if err != nil {
			return err, false
		}
	}
	if r.Login != "" {
		return nil, true
	} else { return nil, false }
	return nil, true
}