# go-curl
golang实现的curl包
# 使用

var browser *curl.Browser;
browser = curl.NewBrowser();
//设置代理
browser.SetProxyUrl("http://221.174.153.63:9999");

response,status := browser.Get(url);
response：响应的内容
status: 响应的状态码

