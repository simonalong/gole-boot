package job

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	clientHttp "github.com/simonalong/gole-boot/client/http"
	"github.com/simonalong/gole-boot/errorx"
	listener2 "github.com/simonalong/gole-boot/event"
	_ "github.com/simonalong/gole-boot/otel"
	serverHttp "github.com/simonalong/gole-boot/server/http"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/listener"
	"github.com/simonalong/gole/logger"
	"github.com/simonalong/gole/util"
)

var cfg Config
var jobClient *Client
var registered = false

type Client struct {
	regList *taskList //注册任务列表
	runList *taskList //正在执行任务列表
	mu      sync.RWMutex

	logHandler LogHandler //日志查询handler
}

func init() {
	config.Load()

	if config.Loaded && config.GetValueBoolDefault("gole.job.enable", false) {
		// 核查http是否开启，没有开启，则不启动
		if !config.GetValueBoolDefault("gole.server.http.enable", false) {
			logger.Group("job").Fatalf("http 服务未开启，无法启动cbb-job的客户端服务端调度功能，请先开启：gole.server.http.enable")
			return
		}

		err := config.GetValueObject("gole.job", &cfg)
		if err != nil {
			logger.Group("job").Warnf("读取分布式 job 配置异常, %v", err.Error())
			return
		}
		parameterFix()
		initConnect()
	}
}

func initConnect() {
	_jobClient, err := newClient()
	if err != nil {
		logger.Group("job").Errorf("初始化cbb-job客户端异常：%v", err)
		return
	}
	jobClient = _jobClient
	// 注册执行器
	go jobClient.registry()

	// http服务启动之前注册
	listener.AddListener(listener2.EventOfServerHttpRunStart, func(event listener.BaseEvent) {
		jobClient.addHttpHandler()
	})
}

func parameterFix() {
	if cfg.ServerAddress == "" {
		cfg.ServerAddress = "http://cbb-mid-srv-job:18080"
	}

	if cfg.Timeout == 0 {
		cfg.Timeout = 5 * time.Second
	}

	if cfg.ExecutorName == "" {
		if val := config.GoleCfg.Application.Name; val != "" {
			cfg.ExecutorName = val
		} else {
			logger.Fatalf("gole.application.name 不可为空")
		}
	}

	if cfg.ExecutorAddress == "" {
		ip, _ := util.GetPublicIP()
		cfg.ExecutorAddress = fmt.Sprintf("%v:%v", ip, config.GetValueIntDefault("gole.server.http.port", 8080))
	}
}

func newClient() (*Client, error) {
	taskClient := &Client{}
	taskClient.regList = &taskList{
		data: make(map[string]*Task),
	}
	taskClient.runList = &taskList{
		data: make(map[string]*Task),
	}
	return taskClient, nil
}

func AddJobHandler(jobHandlerName string, task TaskFunc) {
	if jobClient == nil {
		logger.Group("job").Fatalf("请开启cbb-job的配置：cbb.job.enable")
	}
	var t = &Task{}
	t.fn = task
	jobClient.regList.Set(jobHandlerName, t)
	return
}

func (client *Client) addHttpHandler() {
	serverHttp.AddRoute(serverHttp.HmPost, "/job/executor/run", client.runTask)
	serverHttp.AddRoute(serverHttp.HmPost, "/job/executor/kill", client.killTask)
	serverHttp.AddRoute(serverHttp.HmPost, "/job/executor/log", client.taskLog)
	serverHttp.AddRoute(serverHttp.HmGet, "/job/executor/beat", client.beat)
	serverHttp.AddRoute(serverHttp.HmPost, "/job/executor/idleBeat", client.idleBeat)
}

func (client *Client) registry() {
	t := time.NewTimer(time.Second * 0) //初始立即执行
	defer t.Stop()
	req := &Registry{
		RegistryGroup: "EXECUTOR",
		RegistryKey:   cfg.ExecutorName,
		RegistryValue: cfg.ExecutorAddress,
	}
	for {
		<-t.C
		t.Reset(time.Second * time.Duration(20)) //20秒心跳防止过期
		func() {
			result, err := client.call(http.MethodPost, "/api/job/group/registry", req)
			if err != nil {
				logger.Group("job").Errorf("执行器[%v]注册失败，异常：%v", cfg.ExecutorName, err.Error())
				registered = false
				return
			}

			res := &res{}
			_ = json.Unmarshal(result, &res)
			if !errorx.IsOkCode(res.Code) {
				logger.Group("job").Errorf("执行器[%v]注册失败, code=%v, msg=%v", cfg.ExecutorName, res.Code, res.Msg)
				registered = false
				return
			}
			if !registered {
				logger.Group("job").Infof("执行器[%v]注册成功: %v", cfg.ExecutorName, util.ToJsonString(res))
			}
			registered = true
		}()
	}
}

func (client *Client) call(method, action string, bodyReq any) ([]byte, error) {
	header := http.Header{
		"Content-Type": []string{"application/json"},
	}
	code, _, result, err := clientHttp.Call(method, cfg.ServerAddress+action, header, nil, bodyReq)
	if err != nil {
		logger.Group("job").Errorf("请求异常: url=%v, head=%v, bodyReq=%v, code=%v, 异常：%v", cfg.ServerAddress+action, util.ToJsonString(header), util.ToJsonString(bodyReq), code, err)
		return nil, err
	}
	return result, nil
}

func (client *Client) runTask(c *gin.Context) (any, error) {
	client.mu.Lock()
	defer client.mu.Unlock()

	req := RunReq{}
	_, err := util.DataToEntity(c.Request.Body, &req)
	if err != nil {
		logger.Group("job").Errorf("参数解析错误: %v", util.ToJsonString(req))
		return getCallRsp(&req, FailureCode, fmt.Sprintf("params err：%v", util.ToJsonString(req))), errorx.SC_BAD_REQUEST
	}
	logger.Group("job").Debugf("任务[%v]准备执行 %v，参数:%v", req.JobID, req.ExecutorHandler, util.ToJsonString(req))

	if !client.regList.Exists(req.ExecutorHandler) {
		logger.Group("job").Errorf("任务[%v]没有注册：%v", req.JobID, req.ExecutorHandler)
		return getCallRsp(&req, FailureCode, fmt.Sprintf("任务[%v]没有注册：%v", req.JobID, req.ExecutorHandler)), errorx.SC_BAD_REQUEST
	}

	//阻塞策略处理
	if client.runList.Exists(util.ToString(req.JobID)) {
		if req.ExecutorBlockStrategy == coverEarly { //覆盖之前调度
			oldTask := client.runList.Get(util.ToString(req.JobID))
			if oldTask != nil {
				logger.Group("job").Warnf("任务[%v]已经在运行了:%v，根据阻塞策略【%v】取消前一次调用处理，运行本次调用", req.JobID, req.ExecutorHandler, blockStrategyStr(req.ExecutorBlockStrategy))
				oldTask.Cancel()
				client.runList.Del(util.ToString(oldTask.Id))
			}
		} else { //单机串行,丢弃后续调度 都丢弃和报错
			logger.Group("job").Errorf("任务[%v]已经在运行了:%v，根据阻塞策略【%v】进行丢弃本次调用处理", req.JobID, req.ExecutorHandler, blockStrategyStr(req.ExecutorBlockStrategy))
			return getCallRsp(&req, FailureCode, fmt.Sprintf("任务[%v]已经在运行了:%v，根据阻塞策略【%v】进行丢弃本次调用处理", req.JobID, req.ExecutorHandler, blockStrategyStr(req.ExecutorBlockStrategy))), errorx.SC_SERVER_ERROR
		}
	}

	cxt := context.Background()
	task := client.regList.Get(req.ExecutorHandler)
	if req.ExecutorTimeout > 0 {
		task.Ext, task.Cancel = context.WithTimeout(cxt, time.Duration(req.ExecutorTimeout)*time.Second)
	} else {
		task.Ext, task.Cancel = context.WithCancel(cxt)
	}
	task.Id = req.JobID
	task.Name = req.ExecutorHandler
	task.Param = &req

	client.runList.Set(util.ToString(task.Id), task)
	go task.Run(func(code int, msg string) {
		client.callback(task, code, msg)
	})
	return "", nil
}

func (client *Client) killTask(c *gin.Context) (any, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	req := killReq{}
	_, _ = util.DataToEntity(c.Request.Body, &req)

	if !client.runList.Exists(util.ToString(req.JobID)) {
		logger.Group("job").Errorf("任务[%v]没有运行", req.JobID)
		return "", errorx.SC_BAD_REQUEST.WithDetail(fmt.Sprintf("任务[%v]没有运行", req.JobID))
	}
	task := client.runList.Get(util.ToString(req.JobID))
	task.Cancel()
	client.runList.Del(util.ToString(req.JobID))
	logger.Group("job").Warnf("成功取消任务[%v]运行", req.JobID)
	return fmt.Sprintf("成功取消任务[%v]运行", req.JobID), nil
}

func (client *Client) taskLog(c *gin.Context) (any, error) {
	var res *LogRes
	req := LogReq{}
	_, err := util.DataToEntity(c.Request.Body, &req)
	if err != nil {
		logger.Group("job").Errorf("日志请求解析失败：%v", err.Error())
		logContent := LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   0,
			LogContent:  err.Error(),
			IsEnd:       true,
		}
		res.Code = FailureCode
		res.Msg = err.Error()
		res.Content = logContent
		return res, errorx.SC_SERVER_ERROR.WithError(err)
	}

	logger.Group("job").Info("日志请求参数:%+v", req)
	//if client.logHandler != nil {
	//	res = defaultLogHandler(&req)
	//} else {
	//	res = defaultLogHandler(&req)
	//}
	res = defaultLogHandler(&req)
	return res, nil
}

func (client *Client) beat(c *gin.Context) (any, error) {
	logger.Group("job").Info("心跳检测")
	return "", nil
}

func (client *Client) idleBeat(c *gin.Context) (any, error) {
	client.mu.Lock()
	defer client.mu.Unlock()
	req := idleBeatReq{}
	_, err := util.DataToEntity(c.Request.Body, &req)
	if err != nil {
		logger.Group("job").Errorf("参数解析错误: %v", util.ToJsonString(req))
		return "", errorx.SC_BAD_REQUEST.WithDetail(fmt.Sprintf("参数解析错误: %v", util.ToJsonString(req)))
	}

	if client.runList.Exists(util.ToString(req.JobID)) {
		logger.Group("job").Errorf("idleBeat任务[%v]正在运行", req.JobID)
		return "", errorx.SC_SERVER_ERROR.WithDetail(fmt.Sprintf("idleBeat任务[%v]正在运行", req.JobID))
	}
	logger.Group("job").Infof("忙碌检测任务参数:%v", util.ToJsonString(req))
	return "", nil
}

func (client *Client) callback(task *Task, code int, msg string) {
	client.runList.Del(util.ToString(task.Id))
	rspBody, err := client.call(http.MethodPut, "/api/job/job/callback", getCallRsp(task.Param, code, msg))
	if err != nil {
		logger.Group("job").Errorf("任务【%v】执行完发送结果：code=%v, msg=%v, 到服务端失败：%v", task.Id, code, msg, err.Error())
		return
	}
	logger.Group("job").Debugf("任务【%v】执行完发送结果：code=%v, msg=%v, 到服务端成功：%v", task.Id, code, msg, string(rspBody))
}

func getCallRsp(req *RunReq, code int, msg string) []*callElement {
	return call{
		&callElement{
			LogID:      req.LogID,
			LogDateTim: req.LogDateTime,
			HandleCode: code,
			HandleMsg:  msg,
		},
	}
}

func blockStrategyStr(blockStrategyCode string) string {
	switch blockStrategyCode {
	case coverEarly:
		return "覆盖之前调度"
	case discardLater:
		return "丢弃后续调度"
	case serialExecution:
		return "单机串行"
	default:
		return "默认"
	}
}
