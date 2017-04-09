# go-curl
golang实现的curl包
# 使用

var browser *curl.Browser;
browser = curl.NewBrowser();

//设置header
header := map[string]string{
		"Referer":"http://fuli.asia/luyilu/2017/0127/2910.html",
	}
browser.AddHeader(header)

//设置代理
browser.SetProxyUrl("http://221.174.153.63:9999");

response,status := browser.Get(url);
response：响应的内容
status: 响应的状态码

