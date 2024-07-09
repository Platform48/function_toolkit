package toolkits

import (
	"bytes"
	"encoding/json"
	"errors"
	tk "function_toolkit"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"net/http"
	"net/http/httptest"
)

func generateOkJson(w http.ResponseWriter, r *http.Request) {
	ctx := tk.FuncCtx(w, r)

	ctx.Info("Started!")
	ctx.OkResponseJson(tk.Json{"Foo": "Bar", "Heh": 1234})
}

var errorGuard = errors.New("Test Error")

func generateErrJson(w http.ResponseWriter, r *http.Request) {
	ctx := tk.FuncCtx(w, r)
	ctx.Info("Started!")
	ctx.ErrResponse(http.StatusBadRequest, errorGuard, "Test Error Message")
}

var _ = Describe("Toolkit OkJson", func() {
	var rq *http.Request
	var rr *httptest.ResponseRecorder
	var method string
	var body bytes.Buffer

	BeforeEach(func() {
		rq = nil
		rr = httptest.NewRecorder()
		method = "POST"
		body.Reset()
	})

	JustBeforeEach(func() {
		rq = httptest.NewRequest(method, "/", &body)
		rq.Header.Set("Content-Type", "text/plain")
		rr = httptest.NewRecorder()
		handler := http.HandlerFunc(generateOkJson)
		handler.ServeHTTP(rr, rq)
	})

	When("the request is valid", func() {
		It("should return an OK answer", func() {
			var res tk.SuccessResponse
			err := json.NewDecoder(rr.Body).Decode(&res)
			Expect(rr.Code).To(Equal(http.StatusOK))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Data).To(Equal(tk.Json{"Foo": "Bar", "Heh": 1234.}.AsMap()))
		})
	})
})

var _ = Describe("Toolkit ErrJson", func() {
	var rq *http.Request
	var rr *httptest.ResponseRecorder
	var method string
	var body bytes.Buffer

	BeforeEach(func() {
		rq = nil
		rr = httptest.NewRecorder()
		method = "POST"
		body.Reset()
	})

	JustBeforeEach(func() {
		rq = httptest.NewRequest(method, "/", &body)
		rq.Header.Set("Content-Type", "text/plain")
		rr = httptest.NewRecorder()
		handler := http.HandlerFunc(generateErrJson)
		handler.ServeHTTP(rr, rq)
	})

	When("the request is invalid", func() {
		It("should return an error", func() {
			var res tk.ErrorResponse
			err := json.NewDecoder(rr.Body).Decode(&res)
			Expect(rr.Code).To(Equal(http.StatusBadRequest))
			Expect(err).ToNot(HaveOccurred())
			Expect(res.Message).To(Equal("Test Error Message"))
			Expect(res.ErrorCode).To(Equal(http.StatusBadRequest))
		})
	})
})
