# Esurfing-go

golang版的广东电信天翼校园（ZSM验证）登入认证客户端
(原版: [ESurfingDialer](https://github.com/Rsplwe/ESurfingDialer))

[![Go Version](https://img.shields.io/badge/Go-1.24.3-blue)](https://golang.org/)
[![License](https://img.shields.io/badge/License-Apache2.0-green)](LICENSE)

### Feature:
- 内存占用 < 10MB
- 多平台支持 支持 openwrt/windows/mac/linux
- 路由器部署支持
    - 支持多账号(可配合单线多拨使用)
    - 网卡/IP 绑定功能
- 自定义DNS

### 运行
请下载对应平台的的对应架构的二进制文件

- **一般**的路由器都是arm芯片,下载缀为arm64的文件,老一点的处理器下普通arm的就行,有些远古指令集自己编译吧 :)

(不清楚的在issue问吧 这么多设备我也列不出来)

直接运行（默认加载 config.json）
```shell
chmod +x Esurfing-go
./Esurfing-go
```

指定配置文件
```shell
./Esurfing-go name_of_your_config_file.json
```

### 配置文件示例
```json
[
  {
    "username": "用户名",
    "password": "密码",
    "network_check_interval_ms": 1000,
    "max_retry": 100,
    "retry_delay_ms": 1000,
    "out_bound_type": "ip",
    "network_interface_id": "eth0",
    "network_bind_address": "192.168.1.100",
    "use_custom_dns": true,
    "dns_address": "202.96.xxx.xxx"
  }
]
```
### 配置参数

| field                       | 类型      | 默认值  | 说明                                        |
|-----------------------------|---------|------|-------------------------------------------|
| `username`                  | string  | 必填   | 用户名                                       |
| `password`                  | string  | 必填   | 密码                                        |
| `network_check_interval_ms` | int     | 1000 | 网络状态检测间隔(毫秒)                              |
| `max_retry`                 | int     | 0    | 登录失败后最大重试次数(-1=无限重试，0=不重试)                |
| `retry_delay_ms`            | int     | 1000 | 登录失败重试间隔(毫秒)                              |
| `out_bound_type`            | string  | 无    | 出口绑定类型:`ip`-IP绑定 / `id`-网卡绑定 / 留空-系统默认    |
| `network_interface_id`      | string  | 无    | 绑定的网络接口名称(如 eth0)                         |
| `network_bind_address`      | string  | 无    | 绑定的 IP 地址(如100.2.xxx.xxx/192.168.xxx.xxx) |
| `use_custom_dns`            | boolean | 无    | 是否启用自定义dns                                |
| `dns_address`               | string  | 无    | 自定义dns服务器地址                               |

### openwrt golang日志时区修复
```shell
opkg update
opkg install zoneinfo-asia
```
```shell
#时区设置为上海
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
```

### 致谢
特别感谢 [Rsplwe](https://github.com/Rsplwe) 的关键贡献。

