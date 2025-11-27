package job

/**
用来日志查询，显示到xxl-job-admin后台
*/

type LogHandler func(req *LogReq) *LogRes

// 默认返回
func defaultLogHandler(req *LogReq) *LogRes {
	return &LogRes{Code: SuccessCode, Msg: "", Content: LogResContent{
		FromLineNum: req.FromLineNum,
		ToLineNum:   2,
		LogContent:  "这是日志默认返回，说明没有设置LogHandler",
		IsEnd:       true,
	}}
}
