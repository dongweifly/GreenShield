# GreenShield
相对完善一点的敏感词过滤服务

1. 高性能的敏感词匹配；支持DFA，正则，组合词的方式；
2. 敏感词支持文件配置，数据库配置，配置修改后不需要重启，定时自动加载；
3. 接口丰富，使用简单；Client可以指定使用不同的词库和类型来匹配；
4. 完整的工程化配套支持；支持本地，开发，预发，线上测试环境；支持线上服务的健康监测；

## Support 

## Usage

**BUILD**
```bash
 go mod tidy && go build -v
```

**START**
```bash
green_shield -c /opt/green_shield/conf/conf_dev.toml -s /opt/green_shield/sensitive-words
```

也可以使用已经写好的脚本的来编译部署

```bash 
# 默认安装到/opt/green_shield/green_shield, 如果需要修改安装目录，修改脚本中APP_HOME="/opt/green_shield/"即可；
sudo ./build.sh install 
```

```bash
cd your_install_path
# Usage: ./deploy.sh green_shield {start|stop|online|offline|restart} [dev|pre|prod]
# dev 使用测试环境的配置文件，conf_dev.toml; pre 使用预发环境
./deploy.sh green_shield start dev
```

**HTTP API**
服务请求
`/filter/` 敏感词匹配 POST

request
```json
{
  "requestId" : "123456789",
  "text" : "111111, 他妈的, 九十九月九",
  "bizType" :"default",
  "sceneType" : "ad,abuse,politics"
}
```
如果使用文件的方式，文件名字是根据bizType_sceneType来的，对于上面的例子来说，相当于请求default_ad, default_abuse, default_politics
三个文件中的敏感词；

resp
```json
{
  "code": 200,
  "msg": "success",
  "result": {
    "requestId": "123456789",
    "suggestion": "block",
    "reason": "abuse",
    "sensitiveWords": [
      "他妈的"
    ],
    "desensitization": "111111, ***, 九十九月九"
  }
}
```

[详细的接口文档](./api.md)

## RoadMap

## Contact
[欢迎提交issue](https://github.com/dongweifly/GreenShield/issues)

dongwei.fly@gmail.com
