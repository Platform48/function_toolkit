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
	Context        context.Context
	SpanId         string
	spanIdLogField string
	Logger         *zerolog.Logger
	Response       http.ResponseWriter
	Request        *http.Request
}

type ErrorResponse struct {
	SpanId    string `json:"spanId"`
	ErrorCode int    `json:"errorCode"`
	Message   string `json:"message", omitempty`
}
type SuccessResponse struct {
	SpanId string      `json:"spanId"`
	Data   interface{} `json:"data,omitempty"`
}

type Json map[string]any

func FuncCtx(w http.ResponseWriter, r *http.Request) FunctionContext {
	spanId := shortid.MustGenerate()
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	logger := zerolog.New(os.Stdout).With().Ctx(context.Background()).Timestamp().Str("spanId", "["+spanId+"]").Logger()
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
		SpanId:         spanId,
		spanIdLogField: spanIdLogField,
		Logger:         &logger,
		Response:       w,
		Request:        r,
		Context:        r.Context(),
	}

}
func (this FunctionContext) WithCtx(ctx context.Context) FunctionContext {
	return FunctionContext{
		SpanId:         this.SpanId,
		spanIdLogField: this.spanIdLogField,
		Logger:         this.Logger,
		Response:       this.Response,
		Request:        this.Request,
		Context:        ctx,
	}
}

func (this FunctionContext) GetParameter(name string) string {
	return this.Request.URL.Query().Get(name)
}
func (this FunctionContext) HasParameter(name string) bool {
	return this.Request.URL.Query().Has(name)
}
func (this FunctionContext) GetBody() ([]byte, error) {
	var result []byte = make([]byte, this.Request.ContentLength)
	_, err := this.Request.Body.Read(result)

	return result, err
}
func (this FunctionContext) GetJsonBody(result *any) error {
	bytes, err := this.GetBody()
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, result)
}
func (this FunctionContext) GetHeader(name string) string {
	return this.Request.Header.Get(name)
}
func (this FunctionContext) SetResponseHeader(name string, value string) {
	this.Response.Header().Set(name, value)
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
	e.Msgf(this.spanIdLogField+format, args...)
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

func (this FunctionContext) FailResponse(errorCode int, explanation string) {
	this.ErrResponse(errorCode, nil, explanation)
}
func (this FunctionContext) ErrResponse(errorCode int, err error, explanation string) {
	w := this.Response

	this.Errorf("Exception occured (Error code %v) \"%s\": %s", errorCode, explanation, err.Error())

	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	result, err := json.Marshal(ErrorResponse{
		SpanId:    this.SpanId,
		Message:   explanation,
		ErrorCode: errorCode,
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

func (this FunctionContext) OkResponse(format string, data []byte) {
	w := this.Response

	this.Info("Finished processing the request")

	w.Header().Set("Content-Type", format)

	w.WriteHeader(http.StatusOK)
	_, err := w.Write(data)
	if err != nil {
		this.Logger.Panic().Msg(this.spanIdLogField + "Could not send response to user: " + err.Error())
	}
}
func (this FunctionContext) OkResponseJson(object interface{}) {

	resp := SuccessResponse{
		SpanId: this.spanIdLogField,
		Data:   object,
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		this.Logger.Panic().Msg(this.spanIdLogField + "Could not serialize object: " + err.Error())
	}

	this.OkResponse("application/json; charset=utf-8", bytes)
}

func (this Json) AsMap() map[string]interface{} {
	return this
}
