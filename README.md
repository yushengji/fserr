# goerr
Golang错误包，无缝切换原生errors包至goerr。
## 安装
```shell
go get github.com/yushengji/goerr
```
## 特性
- 支持堆栈打印
- 新建一个字符串错误
- 包装已有错误
- 业务错误码支持
- 覆盖常用错误包的工具函数支持
## 新建错误
使用New得到一个错误，使用%+v占位符格式化打印即可看到堆栈信息，如果只想看到错误信息也可以使用%s或.Error
```go
package main

import (
    "fmt"
    "github.com/yushengji/goerr"
)

func main() {
    err := goerr.New("this is a error")
    fmt.Sprintf("%+v", err)
}
```
## 包装错误
使用Wrap包装已有错误，并添加额外信息，例如下层业务报错信息为：db error，使用Wrap可以给用户提供更加人性化的提示信息，例如：系统繁忙，请重试
```go
package main

import (
    "fmt"
    "github.com/yushengji/goerr"
)

func main() {
    err := goerr.Wrap(goerr.New("db error"), "系统繁忙，请重试")
    fmt.Sprintf("%+v", err)
}
```
## 业务错误码
在这里使用NewXX绑定了ErrBasic错误码（它是一个数字）与HTTP错误码，使用WithCode可以将错误码绑定到错误上。
对于业务错误码还可以使用goerr.ParseCode获取业务错误码、HTTP状态码，错误信息将使用错误码附带的错误信息，也就是示例中的basic error，
如果想要修改错误信息，则可以直接在WithCode后添加WithMessage参数即可覆盖。
```go
package main

import (
    "fmt"
    "github.com/yushengji/goerr"
)

func main() {
	goerr.NewInternalError(goerr.ErrBasic, "basic error")
    err := goerr.WithCode(goerr.New("inner error"), goerr.ErrBasic)
    fmt.Sprintf("%+v", err)
	codeErr := goerr.ParseCode(err)
	fmt.Sprintf("http code is %d", codeErr.HttpCode)
}
```
## 性能
11th i7 16G Golang 1.22版本下，新建错误堆栈层数为10层性能如下：

| 操作 | New        | Wrap        | WithCode   |
|----|------------|-------------|------------|
| 时间 |  801 ns/op | 826.2 ns/op | 2077 ns/op |
