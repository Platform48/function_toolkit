package toolkit

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

var isLocalDeployment = (0 == (len(os.Getenv("FUNCTION_NAME")) + len(os.Getenv("FUNCTION_REGION")) + len(os.Getenv("FUNCTION_IDENTITY")) + len(os.Getenv("K_SERVICE")) + len(os.Getenv("K_CONFIGURATION")) + len(os.Getenv("GOOGLE_FUNCTION_TARGET")) + len(os.Getenv("GOOGLE_CLOUD_PROJECT"))))

type FunctionContext struct {
	Context         context.Context
	SpanId          string
	spanIdLogField  string
	Logger          *zerolog.Logger
	Response        http.ResponseWriter
	Request         *http.Request
	stackFrameLevel int
}

// ErrorResponseStruct used internally to return data in an invalid json response. Exported to allow for manually building responses
type ErrorResponseStruct struct {
	SpanId    string `json:"spanId"`
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message,omitempty"`
}

// SuccessResponseStruct used internally to return data in a successful json response. Exported to allow for manually building responses
type SuccessResponseStruct struct {
	SpanId string      `json:"spanId"`
	Data   interface{} `json:"data,omitempty"`
}

// FuncCtx Creates a context from the given request reader and response writer. Generates a new span id and context.Context from the request.
func FuncCtx(w http.ResponseWriter, r *http.Request) FunctionContext {
	spanId := shortid.MustGenerate()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stdout).With().Timestamp().Str("spanId", "["+spanId+"]").Logger()
	if isLocalDeployment {
		logger = logger.Output(zerolog.ConsoleWriter{
			Out:           os.Stdout,
			PartsOrder:    []string{zerolog.TimestampFieldName, zerolog.LevelFieldName, "spanId", zerolog.CallerFieldName, zerolog.MessageFieldName},
			FieldsExclude: []string{"spanId"},
		})
	}

	var spanIdLogField = "[" + spanId + "] "
	if isLocalDeployment {
		spanIdLogField = ""
	}

	return FunctionContext{
		SpanId:          spanId,
		spanIdLogField:  spanIdLogField,
		Logger:          &logger,
		Response:        w,
		Request:         r,
		Context:         r.Context(),
		stackFrameLevel: 1,
	}
}

// WithCtx generates a copy of this ctx object with the given `context.Context` as its context.
func (this FunctionContext) WithCtx(ctx context.Context) FunctionContext {
	return FunctionContext{
		SpanId:   this.SpanId,
		Logger:   this.Logger,
		Response: this.Response,
		Request:  this.Request,
		Context:  ctx,

		spanIdLogField:  this.spanIdLogField,
		stackFrameLevel: 1,
	}
}

// Info logs a message to the console at the INFO level
func (this FunctionContext) Info(message string) {
	this.Logger.Info().Ctx(this.Context).Caller(this.stackFrameLevel).Msg(this.spanIdLogField + message)
}

// Warn logs a message to the console at the WARN level
func (this FunctionContext) Warn(message string) {
	this.Logger.Warn().Ctx(this.Context).Caller(this.stackFrameLevel).Msg(this.spanIdLogField + message)
}

// Error logs a message to the console at the ERROR level
func (this FunctionContext) Error(message string) {
	this.Logger.Error().Ctx(this.Context).Caller(this.stackFrameLevel).Msg(this.spanIdLogField + message)
}

// Debug logs a message to the console at the DEBUG level
func (this FunctionContext) Debug(message string) {
	this.Logger.Debug().Ctx(this.Context).Caller(this.stackFrameLevel).Msg(this.spanIdLogField + message)
}

// Log logs a message to the console at the given log level
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
	e.Ctx(this.Context).Caller(this.stackFrameLevel).Msg(this.spanIdLogField + message)
}

// Logf Formats a message with the given format and logs it to the console at the given log level
func (this FunctionContext) Logf(level int, format string, args ...interface{}) {
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
	e.Ctx(this.Context).Caller(this.stackFrameLevel).Msgf(this.spanIdLogField+format, args...)
}

// Infof Formats a message with the given format and logs it to the console at the INFO level
func (this FunctionContext) Infof(format string, args ...interface{}) {
	this.Logger.Info().Ctx(this.Context).Caller(this.stackFrameLevel).Msgf(this.spanIdLogField+format, args...)
}

// Warnf Formats a message with the given format and logs it to the console at the WARN level
func (this FunctionContext) Warnf(format string, args ...interface{}) {
	this.Logger.Warn().Ctx(this.Context).Caller(this.stackFrameLevel).Msgf(this.spanIdLogField+format, args...)
}

// Errorf Formats a message with the given format and logs it to the console at the ERROR level
func (this FunctionContext) Errorf(format string, args ...interface{}) {
	this.Logger.Error().Ctx(this.Context).Caller(this.stackFrameLevel).Msgf(this.spanIdLogField+format, args...)
}

// Debugf Formats a message with the given format and logs it to the console at the DEBUG level
func (this FunctionContext) Debugf(format string, args ...interface{}) {
	this.Logger.Debug().Ctx(this.Context).Caller(this.stackFrameLevel).Msgf(this.spanIdLogField+format, args...)
}

// FailResponse Builds, sends, and logs an error response for the user with the given message and error code.
func (this FunctionContext) FailResponse(errorCode int, explanation string) {
	this.stackFrameLevel++
	this.ErrResponse(errorCode, nil, explanation)
	this.stackFrameLevel--
}

// ErrResponse Builds, sends, and logs an error response for the user with the given message and error code. Also attaches the given error object to the message in the logs
func (this FunctionContext) ErrResponse(errorCode int, err error, explanation string) {
	this.stackFrameLevel++
	w := this.Response

	if err != nil {
		this.Errorf("Exception occured (Error code %v) \"%s\": %s", errorCode, explanation, err.Error())
	} else {
		this.Errorf("Exception occured (Error code %v) \"%s\"", errorCode, explanation)
	}

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	result, err := json.Marshal(ErrorResponseStruct{
		SpanId:    this.SpanId,
		Message:   explanation,
		ErrorCode: errorCode,
	})

	if err != nil {
		this.stackFrameLevel--
		this.Logger.Panic().Ctx(this.Context).Caller(this.stackFrameLevel + 1).Msg(this.spanIdLogField + "Could not serialize the error to JSON: " + err.Error())
		return
	}
	w.WriteHeader(errorCode)
	_, err = w.Write(result)
	if err != nil {
		this.stackFrameLevel--
		this.Logger.Panic().Ctx(this.Context).Caller(this.stackFrameLevel + 1).Msg(this.spanIdLogField + "Could not send response to user: " + err.Error())
		return
	}
	this.stackFrameLevel--
}

// OkResponse Builds, sends, and logs an OK response for the user with the given Content Type and message body.
func (this FunctionContext) OkResponse(format string, data []byte) {
	this.stackFrameLevel++
	w := this.Response

	this.Info("Finished processing the request")

	w.Header().Set("Content-Type", format)

	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	if err != nil {
		this.stackFrameLevel--
		this.Logger.Panic().Ctx(this.Context).Caller(this.stackFrameLevel + 1).Msg(this.spanIdLogField + "Could not send response to user: " + err.Error())
		return
	}
	this.stackFrameLevel--
}

// OkResponseJson /*
/*
{ "foo":"bar" }
*/
// The object will generate a json response like
/*
{ "spanId": "asdf1234", "data": { "foo":"bar" } }
*/
func (this FunctionContext) OkResponseJson(data interface{}) {
	this.stackFrameLevel++
	resp := SuccessResponseStruct{
		SpanId: this.SpanId,
		Data:   data,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		this.stackFrameLevel--
		this.Logger.Panic().Ctx(this.Context).Caller(this.stackFrameLevel + 1).Msg(this.spanIdLogField + "Could not serialize data: " + err.Error())
		return
	}

	this.OkResponse("application/json; charset=utf-8", bytes)
	this.stackFrameLevel--
}
