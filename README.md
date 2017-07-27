# go-curl
golang实现的curl包。可以自动处理服务的响应的cookie，为下次请求，或重定向加上cookie

# 使用
```
*var browser *curl.Browser;
browser = curl.NewBrowser();

//设置header
header := map[string]string{
		"Referer":"http://fuli.asia/luyilu/2017/0127/2910.html",
	}
browser.AddHeader(header)

//设置代理
browser.SetProxyUrl("http://221.174.153.63:9999");

//GET请求
response,err := browser.Get(url);

//POST请求
loginUrl := "http://rtpush.com/login"
loginParams := map[string]string{
    "username":"15656073550",
    "password":"e10adc3949ba59abbe56e057f20f883e",
}
content,err := b.Post(loginUrl,loginParams)

//上传文件
_, filename, _, _ := runtime.Caller(0)
f:= path.Join(path.Dir(filename), "1.html")
loginParams := map[string]string{
    "username":"15656073550",
    "password":"e10adc3949ba59abbe56e057f20f883e",
}
upUrl := "http://rtpush.com/api/v1/upload"
content,err = b.UploadFile(upUrl, "file", f, loginParams)
log.Print(string(content),err)

```