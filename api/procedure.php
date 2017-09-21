<?php
$y = $_GET['YP'];
$s=$_GET['SP'];
static $executed_proc;
include 'dbconnect.php';

    $sql = "SELECT o_id, t_id, subject_name, o_group, avg_mark, o_year FROM avg_marks
	        WHERE subject_name = :subject AND o_year = :o_year;";

	try {
		$results = $pdo->prepare($sql);
		$results->execute(array(':subject' => $s, ':o_year' => $y));
	}
	catch(PDOException $e) {
		echo "Ошибка при извлечении данных " . $e -> GetMessage();
		exit();
	}
	$answer = $results -> fetchAll();
    $cnt = $results->rowCount();
	if ($cnt > 0) {
		echo "<script>alert('Процедура с такими параметрами уже выполнялась!');</script>";
		exit();
	}

	$sql2 = "CALL Otchet(:yy,:ss);";
	try {
		$results2 = $pdo->prepare($sql2);
		$results2->execute(array(':yy' => $y,':ss' => $s));
	}
	catch(PDOException $e) {
		echo "Ошибка при извлечении данных " . $e -> GetMessage();
		exit();
	}
	echo "<script>alert('Процедура выполнена');</script>";
	exit(); ?>