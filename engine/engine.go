package engine

import (
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type BiLiBiLiEngine struct {
	Url    string
	Cookie string
}

func Get(engine *BiLiBiLiEngine) (result []byte) {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", engine.Url, nil)
	parseCookie(request, engine.Cookie)
	headers := request.Header
	headers.Add("User-Agent", "绕过反爬虫 ^_^ ")
	response, e := client.Do(request)
	if e != nil {
		log.Fatalf("请求[%s]出错:%v \n", engine.Url, e)
	}
	defer response.Body.Close()
	bytes, _ := ioutil.ReadAll(response.Body)
	result = bytes
	return result
}

// cookie 转换
func parseCookie(request *http.Request, str string) {
	split := strings.Split(str, ";")
	for _, value := range split {
		i := strings.Split(value, "=")
		request.AddCookie(addCookie(i[0], i[1]))
	}
}

// 设置 cookie
func addCookie(key, value string) *http.Cookie {
	return &http.Cookie{Name: key, Value: value}
}
