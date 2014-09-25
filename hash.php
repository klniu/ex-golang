<?php
// php -f hash.php
// Author: liudng@gmail.com
// 2014-9-24

$totalStr = 200*10000;
$myArr = array();
$totalLength = 0;
$input = str_split("白日依山尽黄河入海流欲穷千里目更上一层楼危楼高百尺可以摘星辰不感高声语恐惊天上人");
for($i = 0; $i < $totalStr; $i++) {
    $tempStr = "";
    $tempLength = rand(3, 16);
    for($j = 0; $j < $tempLength; $j++) {
        $tempStr .= $input[rand(0, count($input)-1)];
    }
    $totalLength += $tempLength;
    $myArr[] = $tempStr;
}

//测试开始
$begin = microtime();
$Dict = array();
foreach($myArr as $v) {
	$keyArr = str_split($v);
    sort($keyArr);
	$key = implode('', $keyArr);
    if(!isset($Dict[$key])) {
        $Dict[$key] = $v;
    }
}

$end = microtime();

printf("共%d条数据, 其中%d条不重复数据\n", $totalStr, count($Dict));
printf("%s\n%s\n", $begin, $end);
printf("共计%s万条数据，数据总长度%s，完成过滤共耗时%s秒\n", $totalStr/10000, $totalLength, ($end-$begin)/1000);