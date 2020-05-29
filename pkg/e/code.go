package e

const (
	SUCCESS        = 200
	INVALID_PARAMS = 300
	ERROR          = 400
	// app error code
	JOB_GET_FAILED     = 10000
	HISTORY_GET_FAILED = 10001
	FOLDER_GET_FAILED  = 10002
)

var MsgFlags = map[int]string{
	SUCCESS:            "ok",
	INVALID_PARAMS:     "请求参数错误",
	ERROR:              "fail",
	JOB_GET_FAILED:     "任务列表获取失败",
	HISTORY_GET_FAILED: "历史列表获取失败",
	FOLDER_GET_FAILED:  "获取文件夹列表失败",
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
