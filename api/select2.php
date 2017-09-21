<!DOCTYPE html>
<html>
<head>
	<meta charset="utf-8">
	<link rel="stylesheet" type="text/css" href="http://localhost/Stylesheets/req1style.css">
</head>
<body>
<?php
$subj = $_GET['c_subj'];
$select_id = 2;
include 'dbconnect.php';
include 'select_declaration.php';
$sql = "SELECT " . $select_arr['select_2'] .
                   " FROM "  . $from_arr['from_2'] .
                   " WHERE " . $where_arr['where_2'];

try {
	$result = $pdo->prepare($sql);
	$result->execute(array(':subj' => $subj));
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
			<td>ID проекта</td>
			<td>Тема</td>
			<td>Номер зачётной книжки</td>
			<td>Фамилия студента</td>
		</tr>
	</thead>
	<?php foreach ($answer as $line): ?>
		<tr> 
			<?php echo "<td>$line[pr_id]</td>";
				echo "<td>$line[thema]</td>";
				echo "<td>$line[Record_book_num]</td>";
				echo "<td>$line[Last_name]</td>";
?>
		</tr>
	<?php endforeach; ?>
</table>
</body>
</html>
<?php exit(); ?>

?>