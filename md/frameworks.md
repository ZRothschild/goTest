## go语言几个最快最好运用最广的web框架比较

### 令人敬畏的Web框架

如果你为自己设计一个小应用程序，你可能不需要一个Web框架，但如果你正在进行生产，那么你肯定需要一个，一个好的应用程序。

虽然您认为自己拥有必要的知识和经验，但您是否愿意自行编写所有这些功能的代码？
您是否有时间找到生产级外部包来完成这项工作？ 您确定这将与您应用的其余部分保持一致吗？

这些是推动我们使用框架的原因，如果其他人已经做了这些繁琐且艰苦的工作，我们不想自己编写所有必要的功能。

### 介绍

[Go](https://golang.org/)是一种快速增长的开源编程语言，旨在构建简单，快速，可靠的软件。 看看哪些[伟大的公司](https://github.com/golang/go/wiki/GoUsers)使用Go来为他们的服务提供支持。

本文提供了所有必要的信息，以帮助开发人员了解有关使用Go开发Web应用程序的最佳选项的更多信息。

本文包含了最详细的框架比较，通过从尽可能多的角度比较最知名的Web框架：流行度，支持和内置功能：

**Beego**：Go编程语言的开源，高性能Web框架。

* [https://github.com/astaxie/beego](https://github.com/astaxie/beego)
* [https://beego.me](https://beego.me)

**Buffalo**：快速Web开发w/Go。

* [https://github.com/gobuffalo/buffalo](https://github.com/gobuffalo/buffalo)
* [https://gobuffalo.io](https://gobuffalo.io)

**Echo**：高性能，极简主义的Go Web框架。

* [https://github.com/labstack/echo](https://github.com/labstack/echo)
* [https://echo.labstack.com](https://echo.labstack.com)

**Gin**：用Go（Golang）编写的HTTP Web框架。它具有类似Martini的API，具有更好的性能。

* [https://github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
* [https://gin-gonic.github.io/gin](https://gin-gonic.github.io/gin)

**Iris**：Go in the Universe中最快的Web框架。MVC功能齐全。今天拥抱未来。

* [https://github.com/kataras/iris](https://github.com/kataras/iris)
* [https://iris-go.com](https://iris-go.com)

**Revel**：Go语言的高生产力，全栈Web框架。

* [https://github.com/revel/revel](https://github.com/revel/revel)
* [https://revel.github.io](https://revel.github.io)

### 声望

> 人气排序（星级）

| 框架名称 | 当前版本|Watch|Stars| Forks| Issues Open|Issues Close|框架诞生日|最近更新|作者|
| -------- | -----:   | :----: | :----: | :----: |  :----: | :----: | :----: |:----: |:----: |
|beego|stable v1.10.0| 1159 |16762|3557|544|1667|2012|2018年7月31日 GMT+8 下午9:18|[@astaxie](https://github.com/astaxie)|
|buffalo|stable 0.12.6|125 |3317|259|51|508|2014| 2018年8月15日 GMT+8 上午5:26|[@markbates](https://github.com/markbates)|
|Echo|stable 3.1|513|11474 |1010|108|665|2015|2018年8月29日 GMT+8 上午11:49|[@vishr](https://github.com/vishr)|
|Gin|stable v1.stable |915|20122|2377|126|805|2014|2018年9月1日 GMT+8 上午4:38|[@manucorporat](https://github.com/manucorporat)|
|Iris|stable v11|569 |11604|1177|10|888|2016| 2018年9月2日 GMT+8 上午12:27|[@kataras](https://github.com/kataras)|
|Revel|stable 0.19.0|546  |10208|1260|63|769|2011|2018年7月13日 GMT+8 上午3:43|[@robfig](https://github.com/robfig)|

### 学习曲线

| 框架名称 | 示例|示例参考|实时支持|
| :----:| :----:  | :----: | :----: | 
|beego | 49 | [https://beego.me/docs](https://beego.me/docs) | | 
|buffalo | 6 |[https://gobuffalo.io/docs/installation](https://gobuffalo.io/docs/installation)  | [slack](https://gophers.slack.com/messages/buffalo/) | 
|Echo | 20 |[https://echo.labstack.com/cookbook/hello-world](https://echo.labstack.com/cookbook/hello-world)  || 
| Gin| 15 |[https://github.com/gin-gonic/gin/tree/master/examples](https://github.com/gin-gonic/gin/tree/master/examples)  | [gitter](https://gitter.im/gin-gonic/gin) | 
|Iris |92  | [https://github.com/kataras/iris/tree/master/_examples](https://github.com/kataras/iris/tree/master/_examples) |[rocket chat](https://chat.iris-go.com/)  | 
|Revel | 6 | [http://revel.github.io/examples/index.html](http://revel.github.io/examples/index.html) | [gitter](https://gitter.im/revel/community)| 


这里的astaxie和kataras非常出色，其他框架将会有更多的例子，至少对我来说，如果我切换到一个新的框架，这是最快速掌握尽可能多的信息的最有资源的地方。 一个例子就像1000个单词。

### 核心功能

> 排序最多，功能较少

| 核心功能名称 | Iris|Beego|Revel| Echo|Gin|Buffalo|
| :----:| :----:  | :----: | :----: |  :----:  | :----: | :----: | 
| 路由器：命名路径参数和通配符（Router: Named Path Parameters & Wildcard*）| ![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png)|![Go Web Iris中文网](/static/resource/img/articles/true.png)  | ![Go Web Iris中文网](/static/resource/img/articles/true.png) | 
|路由器：正则表达式（Router: Regex*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png)| 
|路由器：分组（Router: Grouping*）|![Go Web Iris中文网](/static/resource/img/articles/true.png)  | ![Go Web Iris中文网](https://assets-cdn.github.com/images/icons/emoji/unicode/2754.png) | ![Go Web Iris中文网](https://assets-cdn.github.com/images/icons/emoji/unicode/2754.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png)|  ![Go Web Iris中文网](/static/resource/img/articles/true.png)|  ![Go Web Iris中文网](/static/resource/img/articles/true.png)| 
|路由器：以上所有混合无冲突（Router: All the above Mixed Without Conflict*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|路由器：自定义`HTTP`错误（Router: Custom HTTP Errors*） |![Go Web Iris中文网](/static/resource/img/articles/true.png) |  ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png) |![Go Web Iris中文网](/static/resource/img/articles/true.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/true.png)  | 
|与`net/http`100％兼容（100% compatible with net/http*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](https://assets-cdn.github.com/images/icons/emoji/unicode/2754.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png) | 
|中间件生态系统（Middleware ecosystem*）|![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | ![Go Web Iris中文网](/static/resource/img/articles/true.png)|![Go Web Iris中文网](/static/resource/img/articles/true.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png) | 
|`Sinatra`风格API（Sinatra-like API*）|![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | ![Go Web Iris中文网](/static/resource/img/articles/true.png)|![Go Web Iris中文网](/static/resource/img/articles/true.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png) | 
|服务器：自动`HTTPS`（Server: Automatic HTTPS*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](https://assets-cdn.github.com/images/icons/emoji/unicode/2754.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/true.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |
|服务器：正常关机（Server: Gracefully Shutdown*）| ![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|服务器：多监听（Server: Multi Listeners*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|完整的`HTTP/2`（Full HTTP/2*）|![Go Web Iris中文网](/static/resource/img/articles/true.png)  |![Go Web Iris中文网](https://assets-cdn.github.com/images/icons/emoji/unicode/2754.png)  | ![Go Web Iris中文网](https://assets-cdn.github.com/images/icons/emoji/unicode/2754.png) |![Go Web Iris中文网](/static/resource/img/articles/true.png) |标准(standard) | 标准(standard)| 
|子域（Subdomains*）|![Go Web Iris中文网](/static/resource/img/articles/true.png)  |次要(secondary) |次要(secondary) |次要(secondary) | 次要(secondary) |次要(secondary)  | 
|会话管理（Sessions*）| ![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png)|![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 次要(secondary) | 
|Websockets| ![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png)|![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|嵌入到应用（View (aka Templates) Embedded Into App*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|View视图引擎 STD（View Engine: STD*）| ![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png)|![Go Web Iris中文网](/static/resource/img/articles/true.png)  | ![Go Web Iris中文网](/static/resource/img/articles/false.png)| 
|View视图引擎 Pug（View Engine: Pug*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|View视图引擎 Django（View Engine: Django*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|View视图引擎 Handlebars（View Engine: Handlebars*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|View视图引擎 Amber（View Engine: Amber*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|数据渲染器：Markdown，JSON，JSONP，XML| ![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png)|![Go Web Iris中文网](/static/resource/img/articles/true.png)  | ![Go Web Iris中文网](/static/resource/img/articles/false.png)| 
|MVC| ![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) |generator |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|缓存（Caching*）| ![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png)|![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|文件服务器（File Server*）| ![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png)|![Go Web Iris中文网](/static/resource/img/articles/true.png)  | ![Go Web Iris中文网](/static/resource/img/articles/true.png) | | 
|文件服务器：嵌入到应用程序（File Server: Embedded Into App*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|响应可以在发送之前的生命周期中多次修改（Response can be Modified Many times through lifecycle before sent*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|Gzip| ![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png) |(中间件) middleware | (中间件)middleware |(中间件) middleware | 
|测试框架（Testing Framework*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|Typescript Transpiler|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|在线编辑（Online Editor*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png)  | 
|日志系统（Logging System*）| ![Go Web Iris中文网](/static/resource/img/articles/true.png)| ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/true.png)|![Go Web Iris中文网](/static/resource/img/articles/true.png)  | ![Go Web Iris中文网](/static/resource/img/articles/false.png)| 
|维护和自动更新（Maintenance & Auto-Updates*）|![Go Web Iris中文网](/static/resource/img/articles/true.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) | ![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png) |![Go Web Iris中文网](/static/resource/img/articles/false.png)  |![Go Web Iris中文网](/static/resource/img/articles/false.png) | 
|性能（Performance）|★★★★★  |★★★  | ★★ | ★★★★★|★★★★★  |★★★| 

> Go中最着名的“Web框架”并不是真正的框架，这意味着：Echo，Gin和Buffalo并不是真正的（功能齐全的）Web框架但Go社区的大多数人都认为他们是。因此他们认为它们与Iris，Beego或Revel相当，因此，我们也有义务将它们纳入此清单。
> 除了 Beego 和 Revel 之外，所有上述框架都可以适应任何中间件这是为 net/http 创建的，其中一些可以轻松地和其他人一起完成有一些黑客[甚至痛苦是这里的选择]。

### 词汇

##### 路由：命名路径参数和通配符(Router: Named Path Parameters & Wildcard)
当您可以将处理程序注册到具有动态路径的路由时。

示例命名路径参数：

```go
"/user/{username}" 匹配 "/user/me", "/user/speedwheel" etc
```

 路径参数`username`匹配值分别是 `me` 和 `speedwheel`。
 
 示例通配符：
 
 ```go
 
 "/user/{path *wildcard}" matches to
 "/user/some/path/here",
 "/user/this/is/a/dynamic/multi/level/path" etc
 ```

`path` 路径参数的值分别是`some/path/here`和`this/is/a/dynamic/multi/level/path`。

> Iris也支持一种称为宏的功能，可以描述为/user/{username：string}或/user/{username：int min（1）}。

#### 路由：正则表达式(Router: Regex)

当您可以使用带有过滤器的动态路径向具有过滤器的路径注册处理程序时，应该传递一些处理程序以执行处理程序。

示例：

 ```go
"/user/{id ^[0-9]$}" matches to "/user/42" but not to "/user/somestring"
 ```
> id路径参数的值为42。

#### 路由：分组(Router: Grouping)

当您可以将公共逻辑或中间件/处理程序注册到共享相同路径前缀的特定路由组时。

示例：

 ```go
myGroup := Group("/user", userAuthenticationMiddleware)
myGroup.Handle("GET", "/", userHandler)
myGroup.Handle("GET", "/profile", userProfileHandler)
myGroup.Handle("GET", "/signup", getUserSignupForm)
 ```

* /user
* /user/profile
* /user/signup

您甚至可以从组中创建子组：

 ```go
myGroup.Group("/messages", optionalUserMessagesMiddleware)
myGroup.Handle("GET', "/{id}", getMessageByID)
 ```
* /user/messages/{id}

#### 路由：以上所有混合无冲突(Router: All the above Mixed Without Conflict)

这是一个先进但有用的功能，我们许多人希望它由路由器或Web框架支持，目前只有Iris在Go世界中支持这一功能。

这意味着`/{path *wildcard}`和`/user/{username}`和`/user/static`和`/user/{path*wildcard}`之类的东西可以在同一个路由器中注册，它可以正确匹配而不会受到静态路径的冲突（/user/static）或通配符（/{path*wildcard}）。

#### 路由：自定义HTTP错误(Router: Custom HTTP Errors)

当您可以为“错误”状态代码注册处理程序时。 错误http状态代码是> = 400状态代码，`i.s Not Found 404`。

示例：

 ```go
OnErrorCode(404, myNotFoundHandler)
 ```
 
 上面的大多数Web框架仅支持404,405和500注册，但像Iris，Beego和Revel这样的功能完全支持任何状态代码甚至任何错误代码（Iris仅支持任何错误）。
 
 #### 100% compatible with net/http （与net/http 100％兼容）
 
 意味着你有：
 
* 框架为您提供了直接访问`*http.Request`和`http.ResponseWriter`的上下文。
* 一种将`net/http`处理程序转换为特定框架的`Handler`类型的方法。

#### 中间件生态系统 (Middleware ecosystem)

当您不必自己用中间件包装每个处理程序时，框架会为您提供一个完整的引擎来定义流，全局或每个路由或每组路由。 i.e Use(middleware), Done(middleware) etc.

#### 类似 Sinatra API （Sinatra-like API）

类似于`Sinatra`的`APIRegister`在运行时处理程序中用于特定`HTTP`方法（和路径参数）的路由。

示例：

 ```go
.Get or GET("/path", gethandler)
.Post or POST("/path", postHandler)
.Put or PUT("/path", putHandler) and etc.
 ```

#### 服务器：自动HTTPS（Server: Automatic HTTPS）

当框架的服务器支持注册和自动续订`SSL`认证以管理`SSL/TLS`传入连接（https）时。 最着名的自动https提供程序是[letsencrypt](https://letsencrypt.org/)。

#### 服务器：正常关机 (Server: Gracefully Shutdown)

按CTRL + C关闭终端应用程序时; 服务器将正常关闭，等待一些连接完成其工作（具有特定的超时）或触发自定义事件以进行清理（即数据库关闭）。

#### 服务器：多监听（Server: Multi Listeners）

当框架的服务器支持注册自定义`net.Listener`或使用多个http服务器和地址提供Web应用程序时。

#### 完整的`HTTP/2`（Full HTTP/2）

当框架支持带有`https`的`HTTP/2`和服务器推送功能时。

####  子域（Subdomains）

当您可以直接从Web应用程序注册每个x，y子域的路由。

这个框架不支持这个功能，但你仍然可以通过启动多个http服务器来实现它，这样做的缺点是主应用程序和子域没有连接，默认情况下不可能在它们之间共享逻辑。

#### 会话（Sessions）

* 支持http会话并准备在特定处理程序中使用时。
* 一些Web框架支持后端数据库来存储会话，因此您可以在服务器重新启动之间获得持久性。Buffalo使用gorilla会话，这些会话比其他实现慢一点。

示例：

 ```go
func setValue(context http_context){
    s := Sessions.New(http_context)
    s.Set("key", "my value")
}

func getValue(context http_context){
    s := Sessions.New(http_context)
    myValue := s.Get("key")
}

func logoutHandler(context http_context){
    Sessions.Destroy(http_context)
}
 ```
Wiki: [https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol#HTTP_session](https://en.wikipedia.org/wiki/Hypertext_Transfer_Protocol#HTTP_session)

#### 套接字(WebSockets)

当框架支持websocket通信协议时。 实现是不同的。

您应该搜索他们的示例以查看适合您的内容。 我尝试所有这些的同事告诉我，与其他API相比，Iris使用更简单的API实现了最具特色的`webosocket`连接。

Wiki: [https://en.wikipedia.org/wiki/WebSocket](https://en.wikipedia.org/wiki/WebSocket)

#### 视图（aka模板）嵌入到应用程序中

通常，您必须将所有模板文件与Web应用程序的可执行文件并行传输。 嵌入到应用程序意味着框架支持与`go-bindata`集成，因此最终的可执行文件包含其中的模板，表示为`[]byte`。

**什么是视图引擎？**

当框架支持模板加载，自定义和构建模板功能时，可以在关键部件上完成。

#### 视图引擎：STD

当框架支持通过标准`html/template`解析器加载模板时。

#### 视图引擎：Pug

当框架支持通过`Pug`解析器加载模板时。

#### 视图引擎：Django

当框架支持通过`Django`解析器加载模板时。

#### 视图引擎：Handlebars

当框架支持通过`Handlebars`解析器加载模板时。

#### 视图引擎：Amber

当框架支持通过`Amber`解析器加载模板时。

#### 渲染：Markdown，JSON，JSONP，XML（Renderer: Markdown, JSON, JSONP, XML）

当框架的上下文为您提供一种简单的方法来轻松地发送/和自定义各种内容类型的响应。

#### MVC

模型 - 视图 - 控制器（MVC）是用于在计算机上实现用户界面的软件架构模式。 
它将给定的应用程序划分为三个相互关联的部分。 
这样做是为了将信息的内部表示与向用户呈现和接受信息的方式分开。
MVC设计模式将这些主要组件分离，从而实现高效的代码重用和并行开发。

* Iris支持完整的MVC功能，可以在运行时注册。
* Beego仅支持方法和模型匹配，可以在运行时注册。
* Revel支持方法，路径和模型匹配，只能通过生成器（必须运行以构建Web应用程序的不同软件）注册。

Wiki: [https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller](https://en.wikipedia.org/wiki/Model%E2%80%93view%E2%80%93controller)

#### 高速缓存(Caching)

Web缓存（或HTTP缓存）是用于临时存储（缓存）Web文档（例如HTML页面和图像）的信息技术，以减少服务器滞后。
通过它的Web缓存系统文档; 如果满足某些条件，可以满足后续要求。[1] Web缓存系统可以指设备或计算机程序。

Wiki: [https://en.wikipedia.org/wiki/Web_cache](https://en.wikipedia.org/wiki/Web_cache)

#### 文件服务器(File Server)

当您可以将（物理）目录注册到将自动向客户端提供此目录文件的路由时。

#### 文件服务器：嵌入到应用程序(File Server: Embedded Into App)

通常，您必须传输所有静态文件（如资产; css，javascript文件...）以及应用程序的可执行文件。
支持此功能的框架使您有机会将所有这些数据嵌入到应用程序中，表示为`[]byte`，它们的响应时间也更快，因为服务器可以直接为它们提供服务，而无需在物理位置查找文件。

#### 响应可以在发送之前的生命周期中多次修改(Response can be Modified Many times through lifecycle before sent)

目前只有Iris通过其http_context中的内置响应编写器支持此功能。

当框架支持此功能时，您可以在发送到客户端之前检索或重置或修改写入的状态代码，正文和标题（在基于`net/http`的Web框架中，默认情况下这是不可能的，因为无法检索或更改正文和状态代码 书面）。

#### Gzip

当你在路由的处理程序中并且你可以更改响应编写器以便使用gzip压缩发送响应时，框架应该处理已发送的头文件，如果发生任何错误，它应该将响应写入恢复正常。
它也应该能够检查客户端是否支持gzip。

> gzip是一种文件格式和用于文件压缩和解压缩的软件应用程序

Wiki: [https://en.wikipedia.org/wiki/Gzip](https://en.wikipedia.org/wiki/Gzip)

#### 测试框架(Testing Framework)

当您可以使用特定的框架库测试HTTP时，它的工作就是帮助您轻松编写更好的测试。

示例（目前，只有Iris支持）:

 ```go
func TestAPI(t *testing.T) {
    app := myIrisApp() 
    tt := httptest.New(t, app)
    tt.GET("/admin").WithBasicAuth("name", "pass").Expect().
    Status(httptest.StatusOK).Body().Equal("welcome")
}
 ```
 
myIrisApp返回你想象中的Web应用程序，它有一个`/admin`的GET处理程序，受基本身份验证保护。

上面的简单测试检查/admin是否以状态OK响应，并且使用特定用户名和密码传递身份验证，并且其正文为`welcome`。

#### Typescript Transpiler

`Typescript`目标是成为ES6的超集，除了标准定义的所有新东西之外，还将添加一个静态类型系统。 
`Typescript`还有一个转换器，它将我们的`Typescript`代码（即ES6 +类型）转换为ES5或ES3 `javascript`代码，因此我们可以在今天的浏览器中使用它。

#### 在线编辑(Online Editor)

在线编辑器借助在线编辑器，您可以快速轻松地在线编译和运行代码。

#### 日志系统（Logging System）

自定义日志记录系统通过提供诸如颜色编码，格式化，日志级别分离，不同日志记录后端等有用功能来扩展本机日志包行为。

#### 维护和自动更新（Maintenance & Auto-Updates）

以非侵入方式通知用户“即时”更新框架。


> 本文翻译于 [https://github.com/speedwheel/awesome-go-web-frameworks](https://github.com/speedwheel/awesome-go-web-frameworks)



