package errormap

var (
	ErrParam           = 400001
	ErrQueryUserByName = 600001
	ErrMakeAccessToken = 600002
)

var ErrMap = map[int]string{
	ErrParam:           "参数错误",
	ErrQueryUserByName: "获取用户失败",
	ErrMakeAccessToken: "生成token失败",
}

// GetError
func GetError(code int) Error {
	return Error{code, ErrMap[code]}
}
