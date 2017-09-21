<?php
$chosen_year  = $_GET['c_year'];
$chosen_month = $_GET['c_month'];
$subj         = $GET['c_subj'];
$sql = "";
echo $_GET['c_year'];
switch($select_id) {
    case 1: $sql = "SELECT " . $select_arr['select_1'] . " FROM " . $from_arr['from_1'] . " WHERE " . $where_arr['where_1'] . " ORDER BY " . $order_arr['order_1'] . ";";
	try { 
    $result = $pdo->prepare($sql); 
    $result->BindValue(":chosen_year", $chosen_year, PDO::PARAM_INT);
}
catch (PDOException $e){
    $output = 'ошибка при извлечении данных  '. $e->getMessage();
    echo $output."<br>";
    exit();
}
                    break; 
    case 2: $sql = "SELECT " . $select_arr['select_2'] .
                   " FROM "  . $from_arr['from_2'] .
                   " WHERE " . $where_arr['where_2'];
            $result = $pdo->prepare($sql);
            $result->BindValue(":chosen_year", $chosen_year, PDO::PARAM_INT);
            $result->BindValue(":subj", $subj, PDO::PARAM_STR);
                    break;
    case 3: $sql = "SELECT " . $select_arr['select_3_5'] .
                   " FROM "  . $from_arr['from_3_5'] .
                   " WHERE " . $where_arr['where_3'];
            $result = $pdo->prepare($sql);
                    break; 
    case 4: $sql = "SELECT " . $select_arr['select_4_6'] .
                   " FROM "  . $from_arr['from_4'] .
                   " WHERE " . $where_arr['where_4'];
            $result = $pdo->prepare($sql);
                    break; 
    case 5: $sql = "SELECT " . $select_arr['select_3_5'] .
                   " FROM  " . $from_arr['from_3_5'] .
                   " WHERE " . $where_arr['where_5'];
            $result = $pdo->prepare($sql);
            $result->BindValue(":chosen_year", $chosen_year, PDO::PARAM_INT);
                    break; 
    case 6: $sql = "SELECT " . $select_arr['select_4_6'] .
                   " FROM "  . $from_arr['from_6'];
                    break; 
    default: break;
}
?>