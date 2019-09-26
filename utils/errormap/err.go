package errormap

// Error 错误结构体
type Error struct {
	Code   int    `json:"code"`
	Detail string `json:"detail"`
}
