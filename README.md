# Esurfing-go

golang版的广东电信天翼校园（ZSM验证）登入认证客户端
(原版: [ESurfingDialer](https://github.com/Rsplwe/ESurfingDialer))

Feature:
- 更低的内存占用
- 更好的路由器部署支持
- 支持多用户
- 支持绑定特定网卡或IP地址

### 运行环境

- 不限
- 内存 >= 5M

### 运行
直接运行: 读取运行目录下的`config.json`

可指定配置文件名：
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
    "network_bind_address": "100.2.180.34"
  }
]
```
字段说明:
- `network_check_interval_ms` 检查网络的时间间隔 单位:毫秒(ms)
- `max_retry` 登录验证的最大重试次数 负数如`-1`将允许无限次重试 `0`代表不重试
- `retry_delay_ms` 登录失败后重试的时间间隔 单位:毫秒(ms)
- `out_bound_type` 绑定类型 `"ip"`为绑定IP `"id"`为绑定接口 留空使用则使用系统默认
- `network_interface_id` 绑定的接口ID
- `network_bind_address` 绑定的IP地址

### 特别感谢
感谢 [Rsplwe](https://github.com/Rsplwe) 大佬的帮助

### 修正openwrt环境下日志时区问题
安装:
```shell
opkg update
opkg install zoneinfo-asia
```
时区设置为上海:
```shell
ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
```


