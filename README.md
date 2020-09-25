# urllib
golang urllib Simple HTTP Client

### 分支
- master  FastHttp   追求极致的性能
- v1      net/http   拥有完善的功能

### characteristic
- simple http client 
- automatic transcoding utf-8

### Install
```
go get github.com/dollarkillerx/urllib
```
### Use
- GET
```go
	// get
	httpCode, bytes, err := urllib.Get("http://www.baidu.com").Byte()
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("HttpCode: %d \n", httpCode)
	fmt.Println(string(bytes))
```
- GET & Query
```go
	httpCode, bytes, err = urllib.Get("http://www.baidu.com").
		Queries("q","122").
		Queries("h","1213").Byte()   // 生成URL： http://www.baidu.com?q=122&h=1213
```
- POST FROM 
```go
	httpCode, bytes, err = urllib.Post("http://www.baidu.com").
		Params("q","122").
		Params("h","1213").Byte()   // URL： http://www.baidu.com  表单 q=122 h=1213
```
- POST JSON
```go
	httpCode, bytes, err = urllib.Post("http://www.baidu.com").
		SetJson([]byte(`{"MSG":"121321"}`)).Byte()   
```
- PUT & DELETE
```go 
	httpCode, bytes, err = urllib.Delete("http://www.baidu.com").Byte()   
	httpCode, bytes, err = urllib.Put("http://www.baidu.com").Byte()   
```
- 设置代理
```go
	httpCode, bytes, err = urllib.Post("http://www.baidu.com").HttpProxy("http://xxxx.c").Byte()
```
- set timeout
````go
	httpCode, bytes, err = urllib.Post("http://www.baidu.com").SetTimeout(3).Byte()
````
- 设置重试次数
```go 
	httpCode, bytes, err = urllib.Delete("http://www.baidu.com").ByteRetry(3)
```
- set Header & 返回方式Body
```go
	body, err := urllib.Post("http://www.baidu.com").
		SetHeader("xxx", "xxx").
		SetHeader("xxx", "xxx").Body()
	defer body.Body.Close()
```
- Auth
```go 
	urllib.Post("http://www.baidu.com").SetAuth("user","passwd").Body()
```
- 自定义UserAgent & 随即 UserAgent
```go
	urllib.Post("http://www.baidu.com").SetUserAgent("chrome").Body()
	urllib.Post("http://www.baidu.com").RandUserAgent().Body()
```
- Set Cookie
```go
	urllib.Post("http://www.baidu.com").SetCookie().Body()
```
- 禁止自动重定向
```go
urllib.Get("http://www.baidu.com").NoRedirect().Body()
```

