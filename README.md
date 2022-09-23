wgctl 命令行管理wireguard配置的工具，用于快速生成wireguard配置

一、使用要求和配置
```
必须安装docker和docker-compose 配置文件在config/config.yaml
```

二、启动
```
进入项目目录使用./wgctl init,会生成docker-compose文件,无需更改的话,直接docker-compose up -d
端口和网络需修改的话，配置文件也要做相应修改
```

三、示例

1.查看当前可用wireguard配置列表 (-s表示查看服务端)
```
wgctl ps
```

2.查看配置 (-p以二维码展示)
```
wgctl show xxx
```

2.创建配置,不加--name随机分配名称 默认加入root服务端
```
wgctl create
```
创建名称为test 加入9527服务配置
```
wgctl create --name test -j 9527 
```

3.删除配置 (-s表示删除服务端）
```
wgctl delete xxx
```

4.创建服务端配置 

三个字段都不能为空
```
wgctl create --new --name test -s 172.26.0.0/24
```

备注: 可以把wgctl文件放到/usr/local/bin/目录下 就不用每次在项目目录以脚本执行了
