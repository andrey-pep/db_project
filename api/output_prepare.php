<html>
<body bgcolor = silver text = green>
<table align=center border=2 width=80%>
<tbody><tr>

<?php
$cnt = $result->rowCount();
$table_data = $result->fetchAll();
print_r($table_data);
$arr = $table_data['0'];
print_r(array_keys($arr));
foreach (array_keys($arr) as $key) {
    echo "<td align=center>" . $key . "</td>";
}
echo "</tr>";
for ($i=0; $i < $cnt; $i++) {
    echo "<tr>";
    foreach (array_keys($table_data[$i]) as $key) {
        echo "<td align=center>" . $table_data[$i]["$key"] .  "</td>";
    }
    echo "</tr>";
}

?>
</tbody></table>
</body></html>
