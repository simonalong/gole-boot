package tdengine

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/simonalong/gole-boot/constants"
	_ "github.com/simonalong/gole-boot/otel"
	"github.com/simonalong/gole/bean"
	"github.com/simonalong/gole/config"
	"github.com/simonalong/gole/global"
	"github.com/simonalong/gole/logger"
	baseTime "github.com/simonalong/gole/time"
	"github.com/simonalong/gole/util"
	"github.com/simonalong/gole/validate"
	orm "github.com/simonalong/tdorm"
	ormConstants "github.com/simonalong/tdorm/constants"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

const DEFAULT_TDORM_LOGGER_GROUP = "tdorm"

var initLock sync.Mutex

type ConfigOfTdengine struct {
	// 连接类型，支持三种：original, restful, websocket
	ConnectType string `match:"value={websocket} isBlank" errMsg:"连接类型配置错误，类型(#current)不支持"`
	Host        string `match:"isUnBlank"`
	Username    string `match:"isUnBlank"`
	Password    string `match:"isUnBlank"`
	DbName      string
	Port        int `match:"value=0" accept:"false"`
}

func init() {
	config.Load()

	if !config.Loaded || !config.GetValueBoolDefault("gole.tdengine.enable", false) {
		return
	}
}

func GetClient() (*orm.TdClient, error) {
	if bean.GetBean(constants.BeanNameTdengine) != nil {
		return bean.GetBean(constants.BeanNameTdengine).(*orm.TdClient), nil
	}
	initLock.Lock()
	defer initLock.Unlock()
	if bean.GetBean(constants.BeanNameTdengine) != nil {
		return bean.GetBean(constants.BeanNameTdengine).(*orm.TdClient), nil
	}
	tdClient, err := NewClient()
	if err != nil {
		return nil, err
	}
	bean.AddBean(constants.BeanNameTdengine, tdClient)
	return tdClient, nil
}

func GetClientWithName(dbName string) (*orm.TdClient, error) {
	tdengineNameOfDbName := fmt.Sprintf("%s_%s", constants.BeanNameTdengine, dbName)
	if bean.GetBean(tdengineNameOfDbName) != nil {
		return bean.GetBean(tdengineNameOfDbName).(*orm.TdClient), nil
	}
	initLock.Lock()
	defer initLock.Unlock()
	if bean.GetBean(tdengineNameOfDbName) != nil {
		return bean.GetBean(tdengineNameOfDbName).(*orm.TdClient), nil
	}
	tdClient, err := NewClientWithName(dbName)
	if err != nil {
		return nil, err
	}
	bean.AddBean(tdengineNameOfDbName, tdClient)
	return tdClient, nil
}

func NewClient() (*orm.TdClient, error) {
	if !config.GetValueBoolDefault("gole.tdengine.enable", false) {
		logger.Error("tdengine配置开关为关闭，请开启")
		return nil, errors.New("tdengine配置开关为关闭，请开启")
	}

	var cfgOfTdengine ConfigOfTdengine
	err := config.GetValueObject("gole.tdengine", &cfgOfTdengine)
	if err != nil {
		logger.Warn("读取tdengine配置异常")
		return nil, err
	}

	appendConfig(&cfgOfTdengine)

	if success, _, errMsg := validate.Check(cfgOfTdengine); !success {
		logger.Fatalf("tdengine配置异常：%v", errMsg)
		return nil, nil
	}

	// 设置tdorm的日志
	logger.SetGroupLevel(DEFAULT_TDORM_LOGGER_GROUP, config.GetValueString("gole.tdengine.logger.level"))

	var tdClient *orm.TdClient
	if cfgOfTdengine.ConnectType == "" || cfgOfTdengine.ConnectType == "original" {
		tdClient = orm.NewConnectOriginal(cfgOfTdengine.Host, cfgOfTdengine.Port, cfgOfTdengine.Username, cfgOfTdengine.Password, cfgOfTdengine.DbName)
	} else if cfgOfTdengine.ConnectType == "restful" {
		tdClient = orm.NewConnectRest(cfgOfTdengine.Host, cfgOfTdengine.Port, cfgOfTdengine.Username, cfgOfTdengine.Password, cfgOfTdengine.DbName)
	} else if cfgOfTdengine.ConnectType == "websocket" {
		tdClient = orm.NewConnectWebsocket(cfgOfTdengine.Host, cfgOfTdengine.Port, cfgOfTdengine.Username, cfgOfTdengine.Password, cfgOfTdengine.DbName)
	}

	if tdClient == nil {
		return nil, errors.New("tdengine创建失败")
	}

	// 支持opentelemetry
	if config.GetValueBoolDefault("gole.opentelemetry.enable", true) {
		tdClient.AddHook(&OtelTdHook{
			ConnectType: cfgOfTdengine.ConnectType,
			DbName:      cfgOfTdengine.DbName,
			tracer:      global.Tracer,
		})
	}

	// 添加测量点
	if config.GetValueBoolDefault("gole.meter.tdengine.enable", false) {
		tdClient.AddHook(&MeterTdHook{})
	}

	initMeter()
	return tdClient, nil
}

func NewClientWithName(dbName string) (*orm.TdClient, error) {
	if !config.GetValueBoolDefault("gole.tdengine.enable", false) {
		logger.Fatalf("tdengine配置开关为关闭，请开启")
		return nil, errors.New("tdengine配置开关为关闭，请开启")
	}

	if dbName == "" {
		logger.Fatalf("dbName不可为空，请检查配置，或者使用函数 NewClient()")
		return nil, errors.New("dbName不可为空，请检查配置，或者使用函数 NewClient()")
	}

	if config.GetValueString(fmt.Sprintf("gole.tdengine.%v.host", dbName)) == "" {
		logger.Fatalf("gole.tdengine.%v.host配置为空", dbName)
		return nil, errors.New(fmt.Sprintf("gole.tdengine.%v.host配置为空", dbName))
	}

	var cfgOfTdengine ConfigOfTdengine
	err := config.GetValueObject("gole.tdengine."+dbName, &cfgOfTdengine)
	if err != nil {
		logger.Warn("读取tdengine配置异常")
		return nil, err
	}

	appendConfig(&cfgOfTdengine)

	var pOrm *orm.TdClient
	if cfgOfTdengine.ConnectType == "original" {
		pOrm = orm.NewConnectOriginal(cfgOfTdengine.Host, cfgOfTdengine.Port, cfgOfTdengine.Username, cfgOfTdengine.Password, cfgOfTdengine.DbName)
	} else if cfgOfTdengine.ConnectType == "restful" {
		pOrm = orm.NewConnectRest(cfgOfTdengine.Host, cfgOfTdengine.Port, cfgOfTdengine.Username, cfgOfTdengine.Password, cfgOfTdengine.DbName)
	} else if cfgOfTdengine.ConnectType == "" || cfgOfTdengine.ConnectType == "websocket" {
		pOrm = orm.NewConnectWebsocket(cfgOfTdengine.Host, cfgOfTdengine.Port, cfgOfTdengine.Username, cfgOfTdengine.Password, cfgOfTdengine.DbName)
	}

	if pOrm == nil {
		return nil, errors.New("tdengine创建失败")
	}

	// 支持opentelemetry
	if config.GetValueBoolDefault("gole.opentelemetry.enable", true) {
		pOrm.AddHook(&OtelTdHook{
			ConnectType: cfgOfTdengine.ConnectType,
			DbName:      cfgOfTdengine.DbName,
			tracer:      global.Tracer,
		})
	}

	bean.AddBean(constants.BeanNameTdengine+dbName, &pOrm)
	return pOrm, nil
}

func appendConfig(cfgOfTdengine *ConfigOfTdengine) {
	// 兼容旧的配置部分
	if cfgOfTdengine.Username == "" && config.GetValueString("gole.tdengine.user-name") != "" {
		cfgOfTdengine.Username = config.GetValueString("gole.tdengine.user-name")
	}
}

type OtelTdHook struct {
	ConnectType string
	DbName      string
	tracer      trace.Tracer
}

func (tdHook *OtelTdHook) Before(thc *orm.TdHookContext) (*orm.TdHookContext, error) {
	spanName := thc.Sql
	if thc.Sql != "" {
		spanName = strings.SplitN(thc.Sql, " ", 2)[0]
	}
	ctx, _ := tdHook.tracer.Start(global.GetGlobalContext(), "tdengine: "+spanName, trace.WithSpanKind(trace.SpanKindClient))
	thc.Context = ctx
	return thc, nil
}

func (tdHook *OtelTdHook) After(thc *orm.TdHookContext) error {
	span := trace.SpanFromContext(thc.Context)
	defer span.End()

	var attrs []attribute.KeyValue
	attrs = append(attrs, attribute.Key("connect.type").String(thc.GetConnectTypeStr()))
	attrs = append(attrs, attribute.Key("db.name").String(thc.DbName))
	attrs = append(attrs, attribute.Key("start.time").String(baseTime.TimeToStringYmdHmsS(thc.Start)))
	attrs = append(attrs, attribute.Key("db.run.type").String(thc.RunType))
	attrs = append(attrs, attribute.Key("db.Sql").String(thc.Sql))
	attrs = append(attrs, attribute.Key("db.sql.args").StringSlice(getArgsSlice(thc)))
	attrs = append(attrs, attribute.Key("execute.time").String(baseTime.ParseDurationForView(thc.ExecuteTime)))

	switch thc.RunType {
	case ormConstants.EXE:
		if thc.ResultOfExe == nil {
			return nil
		}
		rowsAffected, err := thc.ResultOfExe.RowsAffected()
		if err != nil {
			attrs = append(attrs, attribute.Key("db.err").String(err.Error()))
			span.SetAttributes(attrs...)

			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return err
		}
		attrs = append(attrs, attribute.Key("db.sql.result").Int64(rowsAffected))
	case ormConstants.INSERT, ormConstants.SAVE:
		attrs = append(attrs, attribute.Key("db.sql.result").Int64(thc.ResultOfInsert))
	case ormConstants.BATCH_INSERT:
		attrs = append(attrs, attribute.Key("db.sql.result").Int64(thc.ResultOfInsert))
	}

	if thc.Err != nil {
		attrs = append(attrs, attribute.Key("db.err").String(thc.Err.Error()))
		span.SetAttributes(attrs...)

		span.RecordError(thc.Err)
		span.SetStatus(codes.Error, thc.Err.Error())
		return thc.Err
	} else {
		span.SetAttributes(attrs...)
	}
	return nil
}

func getArgsSlice(thc *orm.TdHookContext) []string {
	var argStrSlice []string
	switch thc.RunType {
	case ormConstants.EXE:
		if thc.ArgsOfExe == nil {
			return nil
		}
		for _, arg := range thc.ArgsOfExe {
			if reflect.TypeOf(arg) == reflect.TypeOf(time.Time{}) {
				argStrSlice = append(argStrSlice, baseTime.TimeToStringYmdHmsS(arg.(time.Time)))
			} else {
				argStrSlice = append(argStrSlice, util.ToString(arg))
			}
		}
	case ormConstants.QUERY:
		if thc.ArgsOfQuery == nil {
			return nil
		}
		for _, arg := range thc.ArgsOfQuery {
			if reflect.TypeOf(arg) == reflect.TypeOf(time.Time{}) {
				argStrSlice = append(argStrSlice, baseTime.TimeToStringYmdHmsS(arg.(time.Time)))
			} else {
				argStrSlice = append(argStrSlice, util.ToString(arg))
			}
		}
	case ormConstants.INSERT, ormConstants.BATCH_INSERT:
		if thc.FieldsArgsOfInsert == nil {
			return nil
		}
		for _, key := range thc.FieldsArgsOfInsert.Keys() {
			arg, _ := thc.FieldsArgsOfInsert.Get(key)
			if reflect.TypeOf(arg) == reflect.TypeOf(time.Time{}) {
				argStrSlice = append(argStrSlice, baseTime.TimeToStringYmdHmsS(arg.(time.Time)))
			} else {
				argStrSlice = append(argStrSlice, util.ToString(arg))
			}
		}
	case ormConstants.SAVE:
		if thc.TagsArgsOfInsert == nil {
			return nil
		}
		for _, key := range thc.TagsArgsOfInsert.Keys() {
			arg, _ := thc.TagsArgsOfInsert.Get(key)
			if reflect.TypeOf(arg) == reflect.TypeOf(time.Time{}) {
				argStrSlice = append(argStrSlice, baseTime.TimeToStringYmdHmsS(arg.(time.Time)))
			} else {
				argStrSlice = append(argStrSlice, util.ToString(arg))
			}
		}

		if thc.FieldsArgsOfInsert == nil {
			return argStrSlice
		}
		for _, key := range thc.FieldsArgsOfInsert.Keys() {
			arg, _ := thc.FieldsArgsOfInsert.Get(key)
			if reflect.TypeOf(arg) == reflect.TypeOf(time.Time{}) {
				argStrSlice = append(argStrSlice, baseTime.TimeToStringYmdHmsS(arg.(time.Time)))
			} else {
				argStrSlice = append(argStrSlice, util.ToString(arg))
			}
		}
	}
	return argStrSlice
}
