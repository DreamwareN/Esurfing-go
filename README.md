# Esurfing-go

广东电信天翼校园（ZSM验证）登入认证客户端
(原版已经被删库了)

### 特性

- 多种运行环境支持 openwrt/windows/macOS/linux/termux
- 支持多账号
- 自定义网卡设备绑定
- 自定义此程序的DNS服务器

### 如何使用

请下载对应平台的的对应架构的二进制文件

(最好添加权限chmod +x Esurfing-go)

直接运行（默认加载运行目录下的config.json）

指定配置文件
```shell
./Esurfing-go -c /path/to/your/config/file
```

### 配置文件示例
```json
[
  {
    "username": "10001234",
    "password": "12345678",
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

多用户配置参考配置文件的格式(是个```[]```) 这不明白的去问ai吧