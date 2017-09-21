<!DOCTYPE html>
<html>
<head>
	<title>Результат</title>
	<meta charset="utf-8">
	<link rel="stylesheet" type="text/css" href="/Stylesheets/req1style.css">
</head>
<body>
<?php
include 'dbconnect.php';


$sql = "SELECT o_id, t_id, subject_name, o_group, avg_mark, o_year FROM avg_marks;";

try {
	$results = $pdo->prepare($sql);
	$results->execute();
}
catch(PDOException $e) {
	echo "Ошибка при извлечении данных " . $e -> GetMessage();
	exit();
}

$answer = $results -> fetchAll();
if(count($answer) == 0) {
	echo "<h1>В отчете нет записей.</h1>";
	exit();
}
echo "<h1>Отчет:</h1>";
?>
<table>
	<thead>
		<tr>
			<td>ID строки отчёта</td>
			<td>ID преподавателя</td>
			<td>Название предмета</td>
			<td>Группа</td>
			<td>Средняя оценка</td>
			<td>Год</td>
		</tr>
	</thead>
	<?php foreach ($answer as $line): ?>
		<tr> 
			<?php echo "<td>$line[o_id]</td>";
				echo "<td>$line[t_id]</td>";
				echo "<td>$line[subject_name]</td>";
				echo "<td>$line[o_group]</td>"; 
				echo "<td>$line[avg_mark]</td>";
				echo "<td>$line[o_year]</td>";
				?>
		</tr>
	<?php endforeach; ?>
</table>
</body>
</html>
<?php exit(); ?>