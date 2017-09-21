<?php
$select_arr = array(
    "select_1"   => "`project`.`Record_book_num`, `teachers`.`Last_name`, `student`.`last_name`, `project`.`thema`, `mark`",
    "select_2"   => "pr_id, thema, Record_book_num, Last_name",
    "select_3_5" => "*",
    "select_4" => "t_id, last_name, Birthdate, Pulpit_num, St_work_date",
	"select_6"   => "teachers.*"
);

$from_arr = array(
    "from_1"   => "protocol JOIN protocol_string USING(p_id)
JOIN project USING(pr_id)
JOIN teachers USING(t_id)
JOIN student ON project.Record_book_num = student.Record_book_num",
    "from_2"   => "project JOIN student USING(Record_book_num)",
    "from_3_5" => "teachers",
    "from_4"   => "teachers LEFT JOIN project USING(t_id)",
    "from_6"   => "(
SELECT COUNT(t_id) as cnt, t_id
FROM teachers LEFT JOIN project USING(t_id)
WHERE Record_book_num IS NOT NULL
GROUP BY t_id) t_c JOIN teachers USING(t_id)"
);

$where_arr = array(
    "where_1" => "YEAR(project.finish_date) = :chosen_year",
    "where_2" => "subject = :subj",
    "where_3" => "Birthdate = (SELECT MAX(Birthdate) FROM teachers)",
    "where_4" => "pr_id IS NULL",
    "where_5" => "t_id NOT IN (
    SELECT t_id
    FROM " . $from_arr['from_4'] .
    " WHERE YEAR(finish_date)=:chosen_year)",
	"where_6" => "cnt = (
SELECT MAX(cnt) as m_cnt
FROM (
SELECT COUNT(t_id) as cnt
FROM teachers LEFT JOIN project USING(t_id)
WHERE Record_book_num IS NOT NULL
GROUP BY t_id) t_count)"
);

$order_arr = array(
    "order_1" => "Record_book_num"
);
?>