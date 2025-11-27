package job

//func successMsg(c *gin.Context, msg string) {
//	c.AbortWithStatusJSON(SuccessCode, &res{
//		Code: SuccessCode,
//		Msg:  msg,
//	})
//}

//func successRsp(c *gin.Context, dataRsp any) {
//	c.AbortWithStatusJSON(SuccessCode, dataRsp)
//}
//
//func errorRsp(c *gin.Context, dataRsp any) {
//	c.AbortWithStatusJSON(FailureCode, dataRsp)
//}
//
//func errorMsg(c *gin.Context, errMsg string) {
//	c.AbortWithStatusJSON(FailureCode, &res{
//		Code: FailureCode,
//		Msg:  errMsg,
//	})
//}
//
//func doneMsg(c *gin.Context, code int, msg string) {
//	data := res{
//		Code: code,
//		Msg:  msg,
//	}
//	c.AbortWithStatusJSON(SuccessCode, data)
//}
