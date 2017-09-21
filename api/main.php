<?php
include 'dbconnect.php';
$pr_flag = $_GET['pr_flag'];
if($pr_flag==1) {
    include 'prepare_procedure.php';
} else {
    include 'select_declaration.php';
    include 'select_prepare.php';
    include 'select_exec.php';
}
include 'output_prepare.php';
include 'output_write.php';
?>