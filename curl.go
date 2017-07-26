/**
1.可设置代理
2.可设置 cookie
3.自动保存并应用响应的 cookie
4.自动为重新向的请求添加 cookie
*/
package curl

import (
	"net/http"
	"net/url"
	"io/ioutil"
	"strings"
	"fmt"
	"bytes"
	"mime/multipart"
	"os"
	"io"
)

type Browser struct {
	header map[string]string;
	cookies []*http.Cookie;
	client *http.Client;
}

//初始化
func NewBrowser() *Browser {
	hc := &Browser{};
	hc.header = make(map[string]string)
	hc.client = &http.Client{};
	//为所有重定向的请求增加cookie
	hc.client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		if len(via) > 0 {
			for _,v := range hc.GetCookie() {
				req.AddCookie(v);
			}
		}
		return nil
	}
	return hc;
}

//设置代理地址
func (self *Browser) SetProxyUrl(proxyUrl string)  {
	proxy := func(_ *http.Request) (*url.URL, error) {
		return url.Parse(proxyUrl);
	};
	transport := &http.Transport{Proxy:proxy};
	self.client.Transport = transport;
}

func (self *Browser) AddHeader(header map[string]string)  {
	for k,v := range header{
		self.header[k] = v
	}
}

//设置请求cookie
func (self *Browser) AddCookie(cookies []*http.Cookie)  {
	self.cookies = append(self.cookies, cookies...);
}

//获取当前所有的cookie
func (self *Browser) GetCookie() ([]*http.Cookie) {
	return self.cookies;
}

//发送Get请求
func (self *Browser) Get(requestUrl string) ([]byte, int) {
	request,_ := http.NewRequest("GET", requestUrl, nil);
	self.setHeader(request)
	self.setRequestCookie(request);
	response,err := self.client.Do(request);
	if err!=nil{
		fmt.Println(err);
		return nil,0;
	}
	defer response.Body.Close();

	//保存响应的 cookie
	respCks := response.Cookies();
	for _,v := range respCks {
		val,_ := request.Cookie(v.Name)
		if(val == nil){
			request.AddCookie(v)
		}
	}

	data, _ := ioutil.ReadAll(response.Body)
	return data, response.StatusCode;
}

//发送Post请求
func (self *Browser) Post(requestUrl string, params map[string]string) ([]byte, int) {
	postData := self.encodeParams(params);
	request,_ := http.NewRequest("POST", requestUrl, strings.NewReader(postData));
	header := map[string]string{"Content-Type":"application/x-www-form-urlencoded"}
	self.AddHeader(header)
	self.setHeader(request)
	self.setRequestCookie(request);

	response,err := self.client.Do(request);
	if err!=nil{
		fmt.Println(err);
		return nil,0;
	}
	defer response.Body.Close();

	//保存响应的 cookie
	respCks := response.Cookies();
	for _,v := range respCks {
		val,_ := request.Cookie(v.Name)
		if(val == nil){
			request.AddCookie(v)
		}
	}

	data, _ := ioutil.ReadAll(response.Body)
	return data,response.StatusCode;
}

//上传文件
func (self *Browser) UploadFile(requestUrl, fieldName, filename string, params map[string]string)  ([]byte, error) {
	bodyBuf := &bytes.Buffer{}
	bodyWriter := multipart.NewWriter(bodyBuf)

	//关键的一步操作
	fileWriter, err := bodyWriter.CreateFormFile(fieldName, filename)
	if err != nil {
		fmt.Println("error writing to buffer")
		return nil,err
	}

	//打开文件句柄操作
	fh, err := os.Open(filename)
	if err != nil {
		fmt.Println("error opening file")
		return nil,err
	}
	defer fh.Close()

	//iocopy
	_, err = io.Copy(fileWriter, fh)
	if err != nil {
		return nil,err
	}

	for key, val := range params {
		_ = bodyWriter.WriteField(key, val)
	}

	contentType := bodyWriter.FormDataContentType()
	bodyWriter.Close()

	request,_ := http.NewRequest("POST", requestUrl, bodyBuf);
	header := map[string]string{"Content-Type":contentType}
	self.AddHeader(header)
	self.setHeader(request)
	self.setRequestCookie(request);

	response,err := self.client.Do(request);
	if err!=nil{
		fmt.Println(err);
		return nil,err
	}
	defer response.Body.Close();

	//保存响应的 cookie
	respCks := response.Cookies();
	for _,v := range respCks {
		val,_ := request.Cookie(v.Name)
		if(val == nil){
			request.AddCookie(v)
		}
	}

	data, _ := ioutil.ReadAll(response.Body)
	return data,nil
}


//为请求设置header
func (self *Browser) setHeader(request *http.Request)  {
	for k,v := range self.header{
		if len(request.Header.Get(k))==0{
			request.Header.Set(k,v)
		}else {
			request.Header.Del(k)
			request.Header.Set(k,v)
		}
	}
}

//为请求设置 cookie
func (self *Browser) setRequestCookie(request *http.Request)  {
	for _,v := range self.cookies{
		val,_ := request.Cookie(v.Name)
		if(val == nil){
			request.AddCookie(v)
		}
	}
}

//参数 encode
func (self *Browser) encodeParams(params map[string]string) string {
	paramsData := url.Values{};
	for k,v := range params {
		paramsData.Set(k,v);
	}
	return paramsData.Encode();
}
