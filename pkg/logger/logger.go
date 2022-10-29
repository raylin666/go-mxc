package logger

import (
	"github.com/raylin666/go-utils/logger"
	"go.uber.org/zap"
	"reflect"
	"time"
)

const (
	LogApp     = "app"
	LogSQL     = "sql"
	LogRequest = "request"
)

type Logger struct {
	*zap.Logger
}

func NewJSONLogger(opts ...logger.Option) (*Logger, error) {
	zaplogger, err := logger.NewJSONLogger(opts...)
	return &Logger{zaplogger}, err
}

func (log *Logger) UseApp() *zap.Logger {
	return log.Logger.Named(LogApp)
}

func (log *Logger) UseSQL() *zap.Logger {
	return log.Logger.Named(LogSQL)
}

func (log *Logger) UseRequest() *zap.Logger {
	return log.Logger.Named(LogRequest)
}

type RequestLogFormat struct {
	TraceId           string              `json:"trace_id"`
	ClientIp          string              `json:"client_ip"`
	Method            string              `json:"method"`
	Path              string              `json:"path"`
	RequestProto      string              `json:"request_proto"`
	RequestReferer    string              `json:"request_referer"`
	RequestUa         string              `json:"request_ua"`
	RequestPostData   string              `json:"request_post_data"`
	RequestBodyData   string              `json:"request_body_data"`
	RequestHeaderData map[string][]string `json:"request_header_data"`
	HttpCode          int                 `json:"http_code"`
	BusinessCode      int                 `json:"business_code"`
	BusinessMessage   string              `json:"business_message"`
	RequestTime       time.Time           `json:"request_time"`
	ResponseTime      time.Time           `json:"response_time"`
	CostSeconds       float64             `json:"cost_seconds"`
}

// RequestLog 打印请求日志
func (log *Logger) RequestLog(rlf *RequestLogFormat, err error) {
	var types = reflect.TypeOf(rlf)
	var values = reflect.ValueOf(rlf)
	var zaplog = log.UseRequest()
	for i := 0; i < types.Elem().NumField(); i++ {
		zaplog = zaplog.With(zap.Any(types.Elem().Field(i).Tag.Get("json"), values.Elem().Field(i).Interface()))
	}

	zaplog.With(zap.Error(err)).Info("REQUEST LOG")
}
