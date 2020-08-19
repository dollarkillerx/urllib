# urllib
golang urllib Simple HTTP Client

### Install
```
go get github.com/dollarkillerx/urllib
```
### Use
- 简单操作
```
Get("http://www.baidu.com").Body()

Get("http://www.baidu.com").
    Params("url","hello").
    Params("url","hello")   // 生成的URL为: http://www.baidu.com?url=hello&url=hello

post := Post("http://168.1xxxxx/cg")
	post.Params("username", "root")
	post.Params("password", "we")
	body, err := post.Body()
```
- send json
``` 
Post("http://168.1xxxxx/cg").
    SetJson(marshal).   // send json
    ByteRetry(3)        // 设置获取Byte模式重试次数3
```
- 设置代理
``` 
	bytes, err := Get("http://www.google.com").
		HttpProxy("http://proxy.com").
		SetTimeout(time.Second * 3).
		Byte()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(bytes)
```
- set timeout
```` 
	bytes, err := Get("http://www.google.com").SetTimeout(time.Second * 3).Byte()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(bytes)
````
