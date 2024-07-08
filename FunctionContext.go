package main

import (
	"context"
	"encoding/json"
	"github.com/rs/zerolog"
	"github.com/teris-io/shortid"
	"net/http"
	"os"
)

const (
	LogLevelDebug = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type FunctionContext struct {
	Context        context.Context
	SpanId         string
	spanIdLogField string
	Logger         *zerolog.Logger
	Response       http.ResponseWriter
	Request        http.Request
}

func FuncCtx() FunctionContext {
	spanId := shortid.MustGenerate()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stdout).With().Ctx(context.Background()).Timestamp().Str("spanId", "["+spanId+"]").Logger()

	logger.Info().Msg("Info Message")
	logger.Error().Msg("Error Message")
	logger.Warn().Msg("Warn Message")
	logger.Debug().Msg("Debug Message")
	logger.Trace().Msg("Trace Message")

	return FunctionContext{
		SpanId:         spanId,
		spanIdLogField: "[" + spanId + "] ",
		Logger:         &logger,
	}

}

func (this FunctionContext) Info(message string) {
	this.Logger.Info().Msg(this.spanIdLogField + message)
}
func (this FunctionContext) Warn(message string) {
	this.Logger.Warn().Msg(this.spanIdLogField + message)
}
func (this FunctionContext) Error(message string) {
	this.Logger.Error().Msg(this.spanIdLogField + message)
}
func (this FunctionContext) Debug(message string) {
	this.Logger.Debug().Msg(this.spanIdLogField + message)
}

func (this FunctionContext) Log(level int, message string) {
	var e *zerolog.Event
	switch level {
	case LogLevelDebug:
		e = this.Logger.Debug()
		break
	case LogLevelInfo:
		e = this.Logger.Info()
		break
	case LogLevelWarn:
		e = this.Logger.Warn()
		break
	case LogLevelError:
		e = this.Logger.Error()
		break
	}
	e.Msg(this.spanIdLogField + message)
}

func (this FunctionContext) Infof(format string, args ...interface{}) {
	this.Logger.Info().Msgf(this.spanIdLogField+format, args...)
}
func (this FunctionContext) Warnf(format string, args ...interface{}) {
	this.Logger.Warn().Msgf(this.spanIdLogField+format, args...)
}
func (this FunctionContext) Errorf(format string, args ...interface{}) {
	this.Logger.Error().Msgf(this.spanIdLogField+format, args...)
}
func (this FunctionContext) Debugf(format string, args ...interface{}) {
	this.Logger.Debug().Msgf(this.spanIdLogField+format, args...)
}

type errorResponse struct {
	SpanId    string `json:"spanId"`
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message", omitempty`
	//Reason    string `json:"reason", omitempty`
}
type successResponse struct {
	SpanId string      `json:"spanId"`
	Data   interface{} `json:"data,omitempty"`
}

func (this FunctionContext) ErrResponse(errorCode int, err error, explanation string) {
	w := this.Response

	this.Errorf("Fatal exception occured (Error code %v) \"%s\": %s", errorCode, explanation, err.Error())

	w.Header().Set("SpanId", this.SpanId)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	result, err := json.Marshal(errorResponse{
		SpanId:  this.SpanId,
		Message: explanation,
	})

	if err != nil {
		this.Logger.Panic().Msg(this.spanIdLogField + "Could not serialize the error to JSON: " + err.Error())
	}
	w.WriteHeader(errorCode)
	_, err = w.Write(result)
	if err != nil {
		this.Logger.Panic().Msg(this.spanIdLogField + "Could not send response to user: " + err.Error())
	}
}

type Json map[string]any

func (this FunctionContext) OkResponseRaw(format string, data []byte) {
	w := this.Response

	this.Info("Finished processing the request")

	w.Header().Set("SpanId", this.SpanId)
	w.Header().Set("Content-Type", format)

	w.WriteHeader(200)
	_, err := w.Write(data)
	if err != nil {
		this.Logger.Panic().Msg(this.spanIdLogField + "Could not send response to user: " + err.Error())
	}
}
func (this FunctionContext) OkResponseJson(object interface{}) {
	this.Info("Generating JSON response")

	resp := successResponse{
		SpanId: this.spanIdLogField,
		Data:   object,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		this.Logger.Panic().Msg(this.spanIdLogField + "Could not serialize object: " + err.Error())
	}

	this.OkResponseRaw("application/json; charset=utf-8", bytes)
}
