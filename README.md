wgctl 命令行管理wireguard配置的工具，用于快速生成wireguard配置

一、使用要求和配置
必须安装docker和docker-compose 配置文件在config/config.yaml

二、启动

进入项目目录使用./wgctl init,会生成docker-compose文件,无需更改的话,直接docker-compose up -d
端口和网络需修改的话，配置文件也要做相应修改

三、示例

1.查看配置列表 wgctl ps 查看当前可用wireguard配置列表(-s表示查看服务端)

2.查看配置 wgctl show id或名称 (-p以二维码展示)

2.创建配置 wgctl create 随机分配名称 默认加入root服务端

wgctl create --name test -j 9527 创建名称为test 加入9527

3.删除配置 wgctl delete id或名称 (-s表示删除服务端）

备注: 可以把wgctl文件放到/usr/local/bin/目录下 就不用每次在项目目录以脚本执行了