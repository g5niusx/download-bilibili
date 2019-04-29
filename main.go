package main

import (
	"download-bilibili/engine"
	"download-bilibili/model"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
)

const configFile = "config.json"

// 全局 cookie，由配置文件读取
var cookie = ""
var path = ""

func main() {
	configByte, i := ioutil.ReadFile(configFile)
	if i != nil {
		log.Fatalf("读取配置文件[%s]异常:%v \n", configFile, i)
	}
	config := model.Config{}
	unmarshal := json.Unmarshal(configByte, &config)
	if unmarshal != nil {
		log.Fatalf("反序列化[%s]异常:%v \n", string(configByte), unmarshal)
	}
	var upCode = config.UpCode
	var cookieStr = config.CookieString
	if cookieStr == "" {
		log.Fatalln("config.json 没有配置 cookie!!!")
	}
	cookie = cookieStr
	path = config.Path

	var url = fmt.Sprintf(`https://space.bilibili.com/ajax/member/getSubmitVideos?mid=%d&pagesize=30&tid=0&page=1&keyword=&order=pubdate`, upCode)
	log.Printf("当前下载的 up 主网址为: %s \n", url)
	bytes := engine.Get(&engine.BiLiBiLiEngine{Url: url, Cookie: cookieStr})
	liResponse := model.BiLiBiLiResponse{}
	e := json.Unmarshal(bytes, &liResponse)
	if e == nil {
		pageDatas := getDownloadUrl(&liResponse)
		var count = 0
		chanstrings := make(chan string, len(pageDatas))
		for _, pageData := range pageDatas {
			pages := pageData.Pages
			for _, page := range pages {
				var kanUrl = `https://www.kanbilibili.com/api/video/%d/download?cid=%d&quality=80&page=%d`
				kanUrl = fmt.Sprintf(kanUrl, pageData.Aid, page.Cid, page.Page)
				result := engine.Get(&engine.BiLiBiLiEngine{Url: kanUrl, Cookie: cookieStr})
				kanResponse := model.KanResponse{}
				e := json.Unmarshal(result, &kanResponse)
				if e != nil {
					log.Fatalf("下载[%s]异常:%v", kanUrl, e)
				}
				downLoadUrl := kanResponse.KanData.DUrl[0].Url
				format := kanResponse.KanData.Format
				if strings.Contains(format, "flv") {
					format = "flv"
				}
				log.Printf("单 p 的下载视频为: %s \v", downLoadUrl)
				var fileName = page.Part + "." + format
				if len(pages) == 1 {
					// 当是单 p 的时候，当前的主标题作为文件名
					fileName = pageData.Title + "." + format
				}
				count++
				go download(downLoadUrl, fileName, pageData.Title, chanstrings)
			}
		}
		log.Printf("预计下载的总视频为:%d\n", count)
		var i = 0
	loop:
		for {
			select {
			case path, err := <-chanstrings:
				i++
				if err {
					log.Printf("[%s]下载成功", path)
				} else {
					log.Printf("[%s]下载失败", path)
				}
				if i >= count {
					break loop
				}
			}
		}

	} else {
		log.Fatal(e)
	}
}

func getDownloadUrl(response *model.BiLiBiLiResponse) []model.PageData {
	if !response.Status {
		log.Fatal("b 站返回失败 ～\n")
		return nil
	}
	if &response.Data != nil && response.Data.VList != nil {
		pageDatas := make([]model.PageData, len(response.Data.VList))
		for i, video := range response.Data.VList {
			// 获取单个视频里面的p数和地址
			pageData := getPageDetails(video)
			pageDatas[i] = pageData
		}
		return pageDatas
	}
	return nil
}

func getPageDetails(video model.Video) (pageData model.PageData) {
	// video 的 api 接口地址，返回单个视频下面的 p 数
	var videoAPIUrl = `https://api.bilibili.com/x/web-interface/view?aid=%d`
	videoAPIUrl = fmt.Sprintf(videoAPIUrl, video.Aid)
	bytes := engine.Get(&engine.BiLiBiLiEngine{Url: videoAPIUrl, Cookie: cookie})
	pageResponse := model.PageResponse{}
	e := json.Unmarshal(bytes, &pageResponse)
	if e == nil {
		if pageResponse.Code == 0 {
			pageData = pageResponse.Data
			pageData.Aid = video.Aid
		} else {
			log.Fatalf("获取[%s]失败: %v\n", videoAPIUrl, e)
		}
	} else {
		log.Fatal(e)
	}
	return pageData
}

func download(url, fileName, dirName string, chanstrings chan string) {
	defer func() {
		chanstrings <- fileName
	}()
	dirName = path + "/" + dirName
	fileName = dirName + "/" + fileName
	liEngine := engine.BiLiBiLiEngine{Url: url, Cookie: cookie}
	result := engine.Get(&liEngine)
	mkdir := os.Mkdir(dirName, os.ModePerm)
	log.Printf("创建目录: %s \n", dirName)
	if mkdir != nil {
		log.Fatalf("创建目录[%s]异常:%v \n", dirName, mkdir)
	}
	log.Printf("创建文件: %s \n", fileName)
	err := ioutil.WriteFile(fileName, result, os.ModeAppend)
	if err != nil {
		log.Fatalf("写入文件[%s]异常:%v \n", fileName, err)
	}
}
