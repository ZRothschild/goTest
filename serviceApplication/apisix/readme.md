### apisix

docker-compose -p docker-apisix up -d


###  示例一

#### route 	457957755790033602   /get/order 

```json

{
  "uri": "/get/order",
  "name": "获取订单",
  "desc": "获取订单第0个版本",
  "methods": [
    "GET",
    "POST",
    "PUT",
    "DELETE",
    "PATCH",
    "HEAD",
    "OPTIONS",
    "CONNECT",
    "TRACE",
    "PURGE"
  ],
  "host": "www.order.com",
  "plugins": {
    "api-breaker": {
      "_meta": {
        "disable": false
      },
      "break_response_body": "{     \"code\": 0,     \"msg\": \"请求成功\",     \"data\": \"https://pnt-console-1316774134.cos.ap-singapore.myqcloud.com/icon/29201.png\" }",
      "break_response_code": 403,
      "healthy": {  #请求成功三次，且状态码为 200 就会返回 break_response_code 与  break_response_body 持续时间  300
        "http_statuses": [
          200
        ],
        "successes": 3
      },
      "max_breaker_sec": 300,
      "unhealthy": {  #请求失败三次，且状态码为 500 502 就会返回 break_response_code 与  break_response_body 持续时间  300
        "failures": 3,
        "http_statuses": [
          500,
          502
        ]
      }
    },
    "prometheus": {
      "_meta": {
        "disable": false
      }
    }
  },
  "service_id": "457956160360678082",
  "labels": {
    "API_VERSION": "v0"
  },
  "status": 1
}

```

###  service  457956160360678082

```json
{
  "name": "订单服务",
  "desc": "订单系统",
  "upstream_id": "457957964951585474",
  "plugins": {
    "prometheus": {
      "_meta": {
        "disable": false
      }
    }
  },
  "hosts": [
    "www.test.com"
  ]
}
```


###  upstream_id  457957964951585474

```json
{
  "nodes": [
    {
      "host": "10.250.221.93",
      "port": 9081,
      "weight": 1
    },
    {
      "host": "10.250.221.93",
      "port": 9082,
      "weight": 1
    }
  ],
  "timeout": {
    "connect": 6,
    "send": 6,
    "read": 6
  },
  "type": "roundrobin",
  "scheme": "http",
  "pass_host": "pass",
  "name": "订单上游",
  "desc": "订单上游",
  "keepalive_pool": {
    "idle_timeout": 60,
    "requests": 1000,
    "size": 320
  }
}


```


####   curl -i -X GET "http://127.0.0.1:9080/get/userinfo" -H "Host: www.order.com"


### 示例二

```json

{
  "uri": "/anything/*",
  "name": "anything",
  "methods": [
    "GET",
    "POST"
  ],
  "host": "example.com",
  "plugins": {
    "prometheus": { # 监控开启
      "_meta": {
        "disable": false
      }
    },
    "request-id": {
      "_meta": {
        "disable": false
      }
    },
    "request-validation": { # 数据body 校验
      "_meta": {
        "disable": false
      },
      "body_schema": {
        "properties": {
          "array_payload": {
            "default": [
              200,
              302
            ],
            "items": {
              "maximum": 599,
              "minimum": 200,
              "type": "integer"
            },
            "minItems": 1,
            "type": "array",
            "uniqueItems": true
          },
          "boolean_payload": {
            "type": "boolean"
          },
          "regex_payload": {
            "maxLength": 32,
            "minLength": 1,
            "pattern": "[[^[a-zA-Z0-9_]+$]]",
            "type": "string"
          }
        },
        "required": [
          "boolean_payload",
          "array_payload",
          "regex_payload"
        ],
        "type": "object"
      }
    }
  },
  "upstream": {
    "nodes": [
      {
        "host": "10.250.221.93",
        "port": 9081,
        "weight": 1
      }
    ],
    "timeout": {
      "connect": 6,
      "send": 6,
      "read": 6
    },
    "type": "roundrobin",
    "scheme": "http",
    "pass_host": "pass",
    "keepalive_pool": {
      "idle_timeout": 60,
      "requests": 1000,
      "size": 320
    }
  },
  "status": 1
}

```

####  curl --header "Content-Type: application/json" --request POST --data '{"boolean-payload":"true","required_payload":"xxx"}' "http://127.0.0.1:9080/anything/order" -H "Host: example.com"


curl "http://127.0.0.1:9180/apisix/admin/upstreams/12" -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1" -X PUT -d '{"uri":"/get/userinfo","name":"重复路由3","methods":["GET","POST","PUT","DELETE","PATCH","HEAD","OPTIONS","CONNECT","TRACE","PURGE"],"upstream":{"nodes":[{"host":"10.250.211.93","port":9082,"weight":1}],"timeout":{"connect":6,"send":6,"read":6},"type":"roundrobin","scheme":"http","pass_host":"pass","keepalive_pool":{"idle_timeout":60,"requests":1000,"size":320}},"labels":{"API_VERSION":"v3"},"status":0}'


curl "http://127.0.0.1:9180/apisix/admin/upstreams/121" -H "X-API-KEY: edd1c9f034335f136f87ad84b625c8f1" -X PUT -d '{"type":"roundrobin","nodes":{"httpbin.org:80":1}}'