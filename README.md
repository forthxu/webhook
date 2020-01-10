# webhook git钩子

分解project参数，dir = 配置的目录/baidu.com/www/，进入目录后执行拉取和推送命令
"cd " + dir + " && git checkout release && git push release release && git show -2 --name-status"

## php版本

php流程建站后钩子

[http://127.0.0.1/githooks.php?project=www.baidu.com](http://127.0.0.1/githooks.php?project=www.baidu.com)

githooks.php文件放到目录，开放php操作系统目录权限

## golang版本

/bin/bash run.sh start

[http://127.0.0.1:9999/release/?project=www.baidu.com&token=abc](http://127.0.0.1:9999/release/?project=www.baidu.com&token=abc)


## Next

目前需要登陆服务器预先 按克隆下来，再添加release远程地址
这部分改为web登陆操作
