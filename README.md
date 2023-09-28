## 框架文档

## 路由

#### 沿用 gin 框架的设计，但是对其中 controller 的入参设计和返回进行了封装，入参定义合适的结构体承接，返回按照 return 方式返回

* 增加 controller 请求参数按照 struct 结构体定义


#### 使用反射实现函数的转变，只在启动时执行，运行时无影响，没有性能损耗

## 日志

#### 在分布式链路跟踪中有两个重要的概念：跟踪（trace）和 跨度（ span）。trace 是请求在分布式系统中的整个链路视图，span 则代表整个链路中不同服务内部的视图，span 组合在一起就是整个 trace 的视图。
#### 参考 Google Dapper 论文 `https://bigbully.github.io/Dapper-translation/`

## 错误处理

* 增加 controller 返回 error 按照自定义 ERR_xxx 常量进行定义，可自行定义 ERR_xxx 常量映射的业务错误码和 http status，错误信息