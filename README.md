# ESurfing-go

golang版的广东电信天翼校园（ZSM验证）登入认证客户端
(原版: [ESurfingDialer](https://github.com/Rsplwe/ESurfingDialer))

更低的内存占用

更好的路由器部署支持

支持多用户 指定用户绑定指定网卡或虚拟网卡

### 运行环境

- 不限
- 内存 >= 5M

### 运行

```shell
./ESurfing-go config.json
#或者
./ESurfing-go name_of_your_config_file.json
```

如果参数留空那么会直接读取同级目录下的`config.json`

可指定配置文件名

### 配置文件
例子:

```json
[
  {
    "username": "用户名1",
    "password": "密码1",
    "network_check_interval_ms": 1000,
    "max_retry": 100,
    "retry_delay_ms": 1000,
    "network_interface_id": "eth0"
  },
  {
    "username": "用户名2",
    "password": "密码2",
    "network_check_interval_ms": 1000,
    "max_retry": 100,
    "retry_delay_ms": 1000,
    "network_interface_id": "eth1"
  }
]
```
个别字段解释
-  `network_check_interval_ms` 检查网络的时间间隔 单位:毫秒(ms)
-  `max_retry` 登录验证的最大重试次数 设置为负数代表比如`-1`将允许无限次重试
-  `retry_delay_ms` 登录验证重试的时间间隔 单位:毫秒(ms)
-  `network_interface_id` 绑定的接口ID 留空:`""`表示使用系统默认的接口

### 特别感谢
感谢 [Rsplwe](https://github.com/Rsplwe) 大佬的帮助