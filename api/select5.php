<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<link rel="stylesheet" type="text/css" href="http://localhost/Stylesheets/req1style.css">
</head>
<body>
<?php
$y = $_GET['c_year'];

$select_id = 1;
include 'dbconnect.php';
include 'select_declaration.php';
 $sql = "SELECT " . $select_arr['select_3_5'] .
                   " FROM  " . $from_arr['from_3_5'] .
                   " WHERE " . $where_arr['where_5'];

try {
	$result = $pdo->prepare($sql);
	$result->execute(array(':chosen_year' => $y));
}
catch(PDOException $e) {
	echo "Ошибка при извлечении данных " . $e -> GetMessage();
	exit();
}


$answer = $result -> fetchAll();
echo "<h1>Результаты запроса</h1>";

?>
<table border=1 align=center width=80%>
	<thead>
		<tr>
			<td>ID преподавателя</td>
			<td>Фамилия преподавателя</td>
			<td>Дата рождения преподавателя</td>
			<td>Номер кафедры</td>
			<td>Дата начала работы</td>
		</tr>
	</thead>
	<?php foreach ($answer as $line): ?>
		<tr> 
			<?php echo "<td>$line[t_id]</td>";
				echo "<td>$line[last_name]</td>";
				echo "<td>$line[Birthdate]</td>";
				echo "<td>$line[Pulpit_num]</td>";
				echo "<td>$line[St_work_date]</td>";
?>
		</tr>
	<?php endforeach; ?>
</table>
</body>
</html>
<?php exit(); ?>

?>