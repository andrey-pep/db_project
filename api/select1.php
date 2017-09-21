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
$sql = "SELECT " . $select_arr['select_1'] . 
" FROM " . $from_arr['from_1'] . 
" WHERE " . $where_arr['where_1'] . 
" ORDER BY " . $order_arr['order_1'] . ";";

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
			<td>Номер зачётки</td>
			<td>Фамилия преподавателя</td>
			<td>Фамилия студента</td>
			<td>Тема проекта</td>
			<td>Оценка</td>
		</tr>
	</thead>
	<?php foreach ($answer as $line): ?>
		<tr> 
			<?php echo "<td>$line[Record_book_num]</td>";
				echo "<td>$line[Last_name]</td>";
				echo "<td>$line[last_name]</td>";
				echo "<td>$line[thema]</td>";
				echo "<td>$line[mark]</td>";
?>
		</tr>
	<?php endforeach; ?>
</table>
</body>
</html>
<?php exit(); ?>

?>