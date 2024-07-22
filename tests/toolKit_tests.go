package toolkits

import (
	"bytes"
	toolkit "github.com/Platform48/function_toolkit"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/mock"
	"net/http"
	"net/http/httptest"
)

type MockJson struct {
	mock.Mock
}

func (m *MockJson) Marshal(v any) ([]byte, error) {
	args := m.Called(v)
	return args.Get(0).([]byte), args.Error(1)
}

// MockResponseRecorder is a custom implementation of http.ResponseWriter.
type MockResponseRecorder struct {
	mock.Mock
	http.ResponseWriter
}

// Write is the mocked version of the Write method.
func (m *MockResponseRecorder) Write(buf []byte) (int, error) {
	args := m.Called(buf)
	return args.Int(0), args.Error(1)
}

var _ = Describe("Toolkit", func() {
	var rq *http.Request
	var rr *httptest.ResponseRecorder
	var method string
	var body bytes.Buffer
	var ctx toolkit.FunctionContext
	var newCtx toolkit.FunctionContext
	var outBuffer bytes.Buffer

	BeforeEach(func() {
		rq = httptest.NewRequest(method, "/", &body)
		rr = httptest.NewRecorder()
		ctx = toolkit.FuncCtx(rr, rq)
		outBuffer.Reset()
		Expect(ctx.SpanId).ToNot(BeEmpty())
		Expect(ctx.Response).ToNot(BeNil())
		Expect(ctx.Request).ToNot(BeNil())
		Expect(ctx.Logger).ToNot(BeNil())
		Expect(ctx.Context).ToNot(BeNil())
	})
	When("WithCtx Success", func() {
		BeforeEach(func() {
			newCtx = ctx.WithCtx(ctx.Context)
		})
		It("should return a copy of the context", func() {
			Expect(newCtx.Context).To(Equal(ctx.Context))
			Expect(newCtx.SpanId).To(Equal(ctx.SpanId))
			Expect(newCtx.Response).To(Equal(ctx.Response))
			Expect(newCtx.Request).To(Equal(ctx.Request))
			Expect(newCtx.Logger).To(Equal(ctx.Logger))
		})
	})
	When("Log is called with log level info", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the info level with info ", func() {
			ctx.Log(toolkit.LogLevelInfo, "info")
			Expect(outBuffer.String()).To(ContainSubstring("info"))
		})
	})
	When("Log is called with log level debug", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the debug level with debug ", func() {
			ctx.Log(toolkit.LogLevelDebug, "debug")
			Expect(outBuffer.String()).To(ContainSubstring("debug"))
		})
	})
	When("Log is called with log level warn", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the warn level with warn ", func() {
			ctx.Log(toolkit.LogLevelWarn, "warn")
			Expect(outBuffer.String()).To(ContainSubstring("warn"))
		})
	})
	When("Log is called with log level error", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the warn level with error ", func() {
			ctx.Log(toolkit.LogLevelError, "error")
			Expect(outBuffer.String()).To(ContainSubstring("error"))
		})
	})

	When("Logf is called with log level info", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the info level with info ", func() {
			ctx.Logf(toolkit.LogLevelInfo, "formatted %s", "info")
			Expect(outBuffer.String()).To(ContainSubstring("formatted info"))
		})
	})
	When("Logf is called with log level debug", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the debug level with debug ", func() {
			ctx.Logf(toolkit.LogLevelDebug, "formatted %s", "debug")
			Expect(outBuffer.String()).To(ContainSubstring("formatted debug"))
		})
	})
	When("Logf is called with log level warn", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the warn level with warn ", func() {
			ctx.Logf(toolkit.LogLevelWarn, "formatted %s", "warn")
			Expect(outBuffer.String()).To(ContainSubstring("formatted warn"))
		})
	})
	When("Logf is called with log level error", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the warn level with error ", func() {
			ctx.Logf(toolkit.LogLevelError, "formatted %s", "error")
			Expect(outBuffer.String()).To(ContainSubstring("formatted error"))
		})
	})
	When("Info is called", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the info level", func() {
			ctx.Info("foo bar")
			Expect(outBuffer.String()).To(ContainSubstring("foo bar"))
		})
	})
	When("Infof is called", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the info level", func() {
			ctx.Infof("formatted %s", "foo bar")
			Expect(outBuffer.String()).To(ContainSubstring("formatted foo bar"))
		})
	})
	When("Warn is called", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the warn level", func() {
			ctx.Warn("foo bar")
			Expect(outBuffer.String()).To(ContainSubstring("foo bar"))
		})
	})
	When("Warnf is called", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the warn level", func() {
			ctx.Warnf("formatted %s", "foo bar")
			Expect(outBuffer.String()).To(ContainSubstring("formatted foo bar"))
		})
	})
	When("Error is called", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the error level", func() {
			ctx.Error("foo bar")
			Expect(outBuffer.String()).To(ContainSubstring("foo bar"))
		})
	})
	When("Errorf is called", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the error level", func() {
			ctx.Errorf("formatted %s", "foo bar")
			Expect(outBuffer.String()).To(ContainSubstring("formatted foo bar"))
		})
	})
	When("Debug is called", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the debug level", func() {
			ctx.Debug("foo bar")
			Expect(outBuffer.String()).To(ContainSubstring("foo bar"))
		})
	})
	When("Debugf is called", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("should write to the debug level", func() {
			ctx.Debugf("formatted %s", "foo bar")
			Expect(outBuffer.String()).To(ContainSubstring("formatted foo bar"))
		})
	})
	When("The logf level isn't specified", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("", func() {
			ctx.Logf(123123, "msg")
			Expect(outBuffer.String()).To(ContainSubstring("msg"))
		})
	})
	When("The log level isn't specified", func() {
		BeforeEach(func() {
			outBuffer = bytes.Buffer{}
			ctx = toolkit.FuncCtx(rr, rq)
			logger := zerolog.New(&outBuffer).With().Timestamp().Str("spanId", "["+"testSpanId"+"]").Logger()
			ctx.Logger = &logger
		})
		It("", func() {
			ctx.Log(123123, "msg")
			Expect(outBuffer.String()).To(ContainSubstring("msg"))
		})
	})
})
