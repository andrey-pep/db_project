<?php
	try {
		$pdo = new PDO("mysql:host=localhost; dbname=study_projects", 'root', '');
		$pdo -> setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);		
	} catch (PDOException $e) {
		echo "������ ����������� � ��: " . $e -> GetMessage();
		exit();
	}
	?>