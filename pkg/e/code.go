package e

const (
	SUCCESS        = 200
	INVALID_PARAMS = 300
	ERROR          = 400
)

var MsgFlags = map[int]string{
	SUCCESS:        "ok",
	INVALID_PARAMS: "请求参数错误",
	ERROR:          "fail",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}

func GetErrorCode(err error) int {
	for key, val := range MsgFlags {
		if err.Error() == val {
			return key
		}
	}
	return 500 //未找到具体错误
}
