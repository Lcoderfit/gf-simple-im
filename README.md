# GoFrame Project

https://goframe.org


1.swagger注解类型
```text
swagger 参考文档
https://hub.fastgit.org/swaggo/swag/blob/master/README_zh-CN.md#mime-types

一般使用下面这六种选项

// @summary 接口概要简介
// @description 接口描述
// @tags 接口所属分类标签
// @param 传入接口的参数（如果没有传入参数则不写此项）
// @success 成功时返回的状态码和数据信息
// @router 路由 请求方法


// @param entity body model.UserApiSignUpReq true "注册请求"
// @param password fromData string true "用户密码"
// @param passport query string true "用户帐号"


object表示参数类型是一个对象，response.JsonResponse是项目中自定义的结构体（返回的数据类型）
// @success 200 {object} response.JsonResponse "执行结果: `true/false`"
// @success 200 {object} response.JsonResponse "执行结果"
// @success 200 {object} model.User "执行结果"
// @success 200 {string} string "执行结果" 


// ShowAccount godoc
// @Summary Show a account
// @Description get string by ID
// @ID get-string-by-int
// @Accept  json
// @Produce  json
// @Param id path int true "Account ID"
// @Success 200 {object} model.Account
// @Header 200 {string} Token "qwerty"
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Failure default {object} httputil.DefaultError
// @Router /accounts/{id} [get]
func (c *Controller) ShowAccount(ctx *gin.Context) {
    id := ctx.Param("id")
    aid, err := strconv.Atoi(id)
    if err != nil {
        httputil.NewError(ctx, http.StatusBadRequest, err)
        return
    }
    account, err := model.AccountOne(aid)
    if err != nil {
        httputil.NewError(ctx, http.StatusNotFound, err)
        return
    }
    ctx.JSON(http.StatusOK, account)
}

// ListAccounts godoc
// @Summary List accounts
// @Description get accounts
// @Accept  json
// @Produce  json
// @Param q query string false "name search by q"
// @Success 200 {array} model.Account
// @Header 200 {string} Token "qwerty"
// @Failure 400,404 {object} httputil.HTTPError
// @Failure 500 {object} httputil.HTTPError
// @Failure default {object} httputil.DefaultError
// @Router /accounts [get]
func (c *Controller) ListAccounts(ctx *gin.Context) {
    q := ctx.Request.URL.Query().Get("q")
    accounts, err := model.AccountsAll(q)
    if err != nil {
        httputil.NewError(ctx, http.StatusNotFound, err)
        return
    }
    ctx.JSON(http.StatusOK, accounts)
}

```

## API操作

Example [celler/controller](https://github.com/swaggo/swag/tree/master/example/celler/controller)

| 注释                 | 描述                                                                                                    |
| -------------------- | ------------------------------------------------------------------------------------------------------- |
| description          | 操作行为的详细说明。                                                                                    |
| description.markdown | 应用程序的简短描述。该描述将从名为`endpointname.md`的文件中读取。                                       |
| id                   | 用于标识操作的唯一字符串。在所有API操作中必须唯一。                                                     |
| tags                 | 每个API操作的标签列表，以逗号分隔。                                                                     |
| summary              | 该操作的简短摘要。                                                                                      |
| accept               | API可以使用的MIME类型的列表。值必须如“[Mime类型](#mime-types)”中所述。                                  |
| produce              | API可以生成的MIME类型的列表。值必须如“[Mime类型](#mime-types)”中所述。                                  |
| param                | 用空格分隔的参数。`param name`,`param type`,`data type`,`is mandatory?`,`comment` `attribute(optional)` |
| security             | 每个API操作的[安全性](#security)。                                                                      |
| success              | 以空格分隔的成功响应。`return code`,`{param type}`,`data type`,`comment`                                |
| failure              | 以空格分隔的故障响应。`return code`,`{param type}`,`data type`,`comment`                                |
| response             | 与success、failure作用相同                                                                               |
| header               | 以空格分隔的头字段。 `return code`,`{param type}`,`data type`,`comment`                                 |
| router               | 以空格分隔的路径定义。 `path`,`[httpMethod]`                                                            |
| x-name               | 扩展字段必须以`x-`开头，并且只能使用json值。                                                            |

## Mime类型

`swag` 接受所有格式正确的MIME类型, 即使匹配 `*/*`。除此之外，`swag`还接受某些MIME类型的别名，如下所示：

| Alias                 | MIME Type                         |
| --------------------- | --------------------------------- |
| json                  | application/json                  |
| xml                   | text/xml                          |
| plain                 | text/plain                        |
| html                  | text/html                         |
| mpfd                  | multipart/form-data               |
| x-www-form-urlencoded | application/x-www-form-urlencoded |
| json-api              | application/vnd.api+json          |
| json-stream           | application/x-json-stream         |
| octet-stream          | application/octet-stream          |
| png                   | image/png                         |
| jpeg                  | image/jpeg                        |
| gif                   | image/gif                         |

## 参数类型

- query
- path
- header
- body
- formData

## 数据类型

- string (string)
- integer (int, uint, uint32, uint64)
- number (float32)
- boolean (bool)
- user defined struct



# 二、boot包中需要加上
```text
import (
	_ "gf-simple-im/packed"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/swagger"
)

// 用于应用初始化
func init() {
	s := g.Server()
	s.Plugin(&swagger.Swagger{})
}
```

# 三、gf swagger
```text
// 自动根据注解生成swagger.json文档
gf swagger
```

# 四、replace手动替换包版本-指定gf包版本为v1.15.5
```text
module gf-simple-im

go 1.13

require (
	github.com/cpuguy83/go-md2man/v2 v2.0.1 // indirect
	github.com/go-openapi/jsonreference v0.19.6 // indirect
	github.com/go-openapi/spec v0.20.3 // indirect
	github.com/go-openapi/swag v0.19.15 // indirect
	github.com/gogf/gf v1.16.4
	github.com/gogf/swagger v1.3.0
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/swaggo/swag v1.7.0 // indirect
	golang.org/x/net v0.0.0-20210716203947-853a461950ff // indirect
	golang.org/x/sys v0.0.0-20210630005230-0f9fa26af87c // indirect
	golang.org/x/tools v0.1.5 // indirect
)

replace github.com/gogf/gf => github.com/gogf/gf v1.15.5

```