<?php
$project = $_GET['project'];

$projectArr= explode('.',$project);
$count = count($projectArr);
$dir = $projectArr[$count-2].'.'.$projectArr[$count-1];
unset($projectArr[$count-1]);
unset($projectArr[$count-2]);
$target = implode('.',$projectArr);
$path = '/web/wwwroot/'.$dir.'/'.$target;
if(!is_dir($path)){
        $projectArr= explode('.',$project);
        $target = array_shift($projectArr);
        $dir = implode('.',$projectArr);
        $path = '/web/wwwroot/'.$dir.'/'.$target;
}
echo $path."<br/>";
system('cd '.$path.' && /usr/bin/git pull && /usr/bin/git push release release >> /tmp/push.log');
