<?php
try {
print_r($result);
    $result->execute();
}
catch (PDOException $e){
    $output = '������ ��� ���������� ������  '. $e->getMessage();
    echo $output."<br>";
    exit();
}
$table_data = $result->fetchAll();
?>