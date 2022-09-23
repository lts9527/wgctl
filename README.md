# 源码打包运行
一、启动

先执行build脚本，打包镜像
进入ctl目录 直接 go run main.go 会有提示命令
先go run main.go init 会初始化一个docker-compose文件 将镜像替换成刚才打包好的 不然会拉docker hub上的 不是当前源码打包的 然后up
配置文件在config/config.yaml里

二、命令示例

1.查看配置列表 wgctl ps 查看当前可用wireguard配置列表(-s表示查看服务端)

2.查看配置 wgctl show id或名称 (-p以二维码展示)

2.创建配置 wgctl create 随机分配名称 默认加入root服务端

wgctl create --name test -j 9527 创建名称为test 加入9527

3.删除配置 wgctl delete id或名称 (-s表示删除服务端）

4.创建服务端配置 

wgctl create --new --name test -s 172.26.0.0/24  三个字段都不能为空