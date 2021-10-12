敏感词匹配

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

`/words/add` 增加敏感词
```json
{
	"wordsName": "default_ad",
	"words": ["测试测试1111111"]
}
```
响应
```json
{
    "code":200,
    "msg":"Success"
}
```

`/words/delete` 删除敏感词 

POST /words/delete

request 
```json
{
    "id":6958
}
```

resp
```json
{
    "code": 200,
    "msg": "Success"
}
```

`/words/fuzzySearch` 

POST /words/fuzzySearch?pageNum=1&pageSize=10

查询敏感词

request 
```json
{
"text": "111111"
}
```

resp
```json
{
    "code":200,
    "msg":"success",
    "pageNo":1,
    "pageSize":10,
    "TotalSize":1,
    "data":[
        {
            "id":6958,
            "text":"111111111",
            "bizType":"default",
            "sceneType":"ad",
            "extent":"",
            "valid":true,
            "updateTime":1633944469,
            "createTime":1633944469
        }
    ]
}
```

GET /words-info/query HTTP/1.1
`/words-info/query` 查询敏感词 
```json
{
    "code":200,
    "msg":"",
    "data":[
        {
            "id":1,
            "name":"default_politics",
            "bizType":"default",
            "sceneType":"politics",
            "parentName":"",
            "updateTime":1630484668,
            "createTime":1628131513
        },
        {
            "id":2,
            "name":"default_abuse",
            "bizType":"default",
            "sceneType":"abuse",
            "parentName":"",
            "updateTime":1631851279,
            "createTime":1628131513
        },
        {
            "id":3,
            "name":"default_ad",
            "bizType":"default",
            "sceneType":"ad",
            "parentName":"",
            "updateTime":1633942700,
            "createTime":1628489847
        },
        {
            "id":4,
            "name":"default_im_custom",
            "bizType":"default",
            "sceneType":"im_custom",
            "parentName":"",
            "updateTime":1630484668,
            "createTime":1629688521
        },
        {
            "id":5,
            "name":"default_mobile",
            "bizType":"default",
            "sceneType":"mobile",
            "parentName":"",
            "updateTime":1630484668,
            "createTime":1629688521
        }
    ]
}
```

`/words-info/names` 查询敏感词库的名称

GET /words-info/names 
```json
{
    "code":200,
    "msg":"",
    "data":[
        {
            "id":1,
            "name":"default_politics"
        },
        {
            "id":2,
            "name":"default_abuse"
        },
        {
            "id":3,
            "name":"default_ad"
        },
        {
            "id":4,
            "name":"default_im_custom"
        },
        {
            "id":5,
            "name":"default_mobile"
        }
    ]
}
```

`/record/add` 查询操作记录

GET 

`/record/search` 搜索操作记录 

POST /record/search?pageNum=1&pageSize=10 
```json
{
    "code":200,
    "msg":"",
    "pageNo":1,
    "pageSize":10,
    "TotalSize":44,
    "data":[
        {
            "id":3,
            "text":"12345678",
            "bizType":"",
            "sceneType":"abuse",
            "operator":"HaiTang",
            "operateType":"add",
            "createTime":0
        },
        {
            "id":4,
            "text":"1111111111",
            "bizType":"",
            "sceneType":"politics",
            "operator":"HaiTang",
            "operateType":"add",
            "createTime":1629784385
        },
        {
            "id":5,
            "text":"1111111111",
            "bizType":"",
            "sceneType":"politics",
            "operator":"HaiTang",
            "operateType":"add",
            "createTime":1629784386
        },
        {
            "id":6,
            "text":"ccccvvvvvvvv",
            "bizType":"",
            "sceneType":"abuse",
            "operator":"HaiTang",
            "operateType":"add",
            "createTime":1629805535
        },
        {
            "id":7,
            "text":"11111111111",
            "bizType":"",
            "sceneType":"politics",
            "operator":"HaiTang",
            "operateType":"add",
            "createTime":1629881434
        },
        {
            "id":8,
            "text":"12345678",
            "bizType":"",
            "sceneType":"politics",
            "operator":"HaiTang",
            "operateType":"add",
            "createTime":1629883180
        },
        {
            "id":9,
            "text":"111111111111",
            "bizType":"",
            "sceneType":"politics",
            "operator":"HaiTang",
            "operateType":"add",
            "createTime":1629950082
        },
        {
            "id":10,
            "text":"111111",
            "bizType":"",
            "sceneType":"politics",
            "operator":"HaiTang",
            "operateType":"add",
            "createTime":1629956535
        },
        {
            "id":11,
            "text":"111111",
            "bizType":"",
            "sceneType":"politics",
            "operator":"HaiTang",
            "operateType":"add",
            "createTime":1629958798
        },
        {
            "id":12,
            "text":"caonima",
            "bizType":"default",
            "sceneType":"abuse",
            "operator":"HaiTang",
            "operateType":"remove",
            "createTime":1629958911
        }
    ]
}
```
