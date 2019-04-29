package model

type BiLiBiLiResponse struct {
	// 是否成功
	Status bool  `json:"status"`
	Data   BData `json:"data"`
}
type Video struct {
	// 视频 id
	Aid uint32 `json:"aid"`
	// 视频标题
	Title string `json:"title"`
}

type BData struct {
	// 视频列表
	VList []Video `json:"vlist"`
	count uint8
	pages uint8
}

// 单个视频的返回，其中包含了多个 p 的视频
type PageResponse struct {
	// 0 查询成功
	Code uint8    `json:"code"`
	Data PageData `json:"data"`
}

type PageData struct {
	Aid   uint32
	Title string `json:"title"`
	// 包含的具体视频
	Pages []Pages `json:"pages"`
}

type Pages struct {
	// 单页视频的标题
	Part string `json:"part"`
	// 视频 id
	Cid uint32 `json:"cid"`
	// 第几页
	Page uint8 `json:"page"`
}

type KanResponse struct {
	KanData KanData `json:"data"`
}

type KanData struct {
	DUrl []KanUrl `json:"durl"`
	// 下载格式
	Format string `json:"format"`
}

type KanUrl struct {
	// kankanbilibili 真实的下载地址
	Url string `json:"url"`
}

type Config struct {
	UpCode       uint32 `json:"up_code"`
	CookieString string `json:"cookie_string"`
	Path         string `json:"path"`
}
