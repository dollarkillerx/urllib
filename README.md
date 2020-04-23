# urllib
golang urllib 简单HTTP请求库

### Install
```
go get github.com/dollarkillerx/urllib
```
### Use
简单操作
```
Get("http://www.baidu.com").Body()

post := Post("http://168.1xxxxx/cg")
	post.Params("username", "root")
	post.Params("password", "we")
	body, err := post.Body()
```
设置代理
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
timeout
```` 
	bytes, err := Get("http://www.google.com").SetTimeout(time.Second * 3).Byte()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(bytes)
````
