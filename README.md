# Esurfing-go
广|东|电|信|天|翼|校|园|登|入|认|证|客|户|端

### 特性

- 支持 openwrt/windows/macOS/linux/termux
- 多账号
- 网卡绑定

### 如何使用

指定配置文件(默认为运行目录的config.json)
```shell
./Esurfing-go -c /path/to/your/config/file
```

### 配置文件示例
```json
[
  {
    "username": "10001234",
    "password": "12345678",
    "check_interval":0,
    "retry_interval":0,
    "bind_device":"eth1",
    "dns_address": "119.29.29.29:53"
  }
]
```

`check_interval`检查网络状态间隔。单位毫秒。

`retry_interval`登录失败重试间隔。单位毫秒。值 <0 = 不重试

`bind_device`绑定的网卡设备名称，比如linux中常见的`eth0` `enp0s1`openwrt的`wan0`。留空则使用系统设置

`dns_address`这个一般留空即可。当系统使用Doh的时候有用。在没有经过登录验证的情况下，Doh是无法正常工作的，无法解析必要的域名导致登陆失败。一般填上DHCP获取的dns即可(请注意要带上端口号)

可按照json格式进行多用户配置
