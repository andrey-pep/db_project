<?php
try {
print_r($result);
    $result->execute();
}
catch (PDOException $e){
    $output = 'ошибка при извлечении данных  '. $e->getMessage();
    echo $output."<br>";
    exit();
}
$table_data = $result->fetchAll();
?>