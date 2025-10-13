# Esurfing-go

广东电信天翼校园（ZSM验证）登入认证客户端
(原版已经被删库了...)

### 特性

- 支持多种运行环境 openwrt/windows/macOS/linux/termux
- 支持多账号
- 自定义网卡设备绑定
- 自定义程序DNS服务器

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
    "check_interval":"0",
    "retry_interval":"0",
    "bind_device":"eth1",
    "dns_address": "119.29.29.29:53"
  }
]
```

`check_interval`检查网络状态的时间间隔。单位毫秒。留空或为0则使用默认值3000毫秒检查一次

`retry_interval`登录失败重试的时间间隔。单位毫秒。留空或为0则使用默认值10000毫秒。填负数为不尝试重试

`bind_device`绑定的网卡设备名称，比如linux中常见的`eth0` `enp0s1`openwrt的`wan0`。留空则使用系统设置

`dns_address`这个一般留空即可。当系统DNS使用DOH的时候有用。在没有经过登录验证的情况下，DOH是不通的，这会在登陆验证的时候无法解析必要的域名导致登陆失败。(请注意要带上端口号)

多用户配置参考配置文件的格式(是个`[]`) 这不明白的去问ai吧