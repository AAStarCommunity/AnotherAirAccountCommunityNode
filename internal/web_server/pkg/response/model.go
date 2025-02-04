package response

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func GetResponse() *Response {
	return &Response{
		httpCode: http.StatusOK,
		Result: &Result{
			Code:    0,
			Message: "",
			Data:    nil,
			Cost:    "",
		},
	}
}
func BadRequest(ctx *gin.Context, data ...any) {
	GetResponse().withDataAndHttpCode(http.StatusBadRequest, ctx, data)
}

func NotFound(ctx *gin.Context, data ...any) {
	GetResponse().withDataAndHttpCode(http.StatusOK, ctx, data)
}

// Success response when business operation is successful
func Success(ctx *gin.Context, data ...any) {
	if data != nil {
		GetResponse().SuccessWithData(ctx, data[0])
		return
	}
	GetResponse().Success(ctx)
}

func Created(ctx *gin.Context, data ...any) {
	if data != nil {
		GetResponse().withDataAndHttpCode(http.StatusCreated, ctx, data[0])
		return
	}
	GetResponse().SetHttpCode(http.StatusCreated).Success(ctx)
}

func InternalServerError(ctx *gin.Context, data ...any) {
	if data != nil {
		GetResponse().withDataAndHttpCode(http.StatusInternalServerError, ctx, data[0])
		return
	}
	GetResponse().SetHttpCode(http.StatusInternalServerError).Fail(ctx)
}

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Cost    string      `json:"cost"`
}

type Response struct {
	httpCode int
	Result   *Result
}

// Fail response when an error occurs
func (r *Response) Fail(ctx *gin.Context) *Response {
	r.SetCode(http.StatusInternalServerError)
	r.json(ctx)
	return r
}

// FailCode response with custom error code
func (r *Response) FailCode(ctx *gin.Context, code int, msg ...string) *Response {
	r.SetHttpCode(code)
	r.SetCode(code)
	if msg != nil {
		r.WithMessage(msg[0])
	}
	r.json(ctx)
	return r
}

// Success response when operation is successful
func (r *Response) Success(ctx *gin.Context) *Response {
	r.SetCode(http.StatusOK)
	r.json(ctx)
	return r
}

func (r *Response) SuccessWithDataAndCode(code int, ctx *gin.Context, data interface{}) *Response {
	r.SetHttpCode(http.StatusOK)
	r.SetCode(code)
	if data != nil {
		r.WithData(data)
	}
	r.json(ctx)
	return r
}

// SuccessWithData response with data when operation is successful
func (r *Response) SuccessWithData(ctx *gin.Context, data interface{}) *Response {
	r.SetCode(http.StatusOK)
	r.WithData(data)
	r.json(ctx)
	return r
}

func (r *Response) withDataAndHttpCode(code int, ctx *gin.Context, data interface{}) *Response {
	r.SetHttpCode(code)
	if data != nil {
		r.WithData(data)
	}
	r.json(ctx)
	return r
}

// SetCode sets the return code
func (r *Response) SetCode(code int) *Response {
	r.Result.Code = code
	return r
}

// SetHttpCode sets the HTTP status code
func (r *Response) SetHttpCode(code int) *Response {
	r.httpCode = code
	return r
}

type defaultRes struct {
	Result any `json:"result"`
}

// WithData sets the return data
func (r *Response) WithData(data interface{}) *Response {
	switch data.(type) {
	case string, int, bool:
		r.Result.Data = &defaultRes{Result: data}
	default:
		r.Result.Data = data
	}
	return r
}

// WithMessage sets the return custom error message
func (r *Response) WithMessage(message string) *Response {
	r.Result.Message = message
	return r
}

// json returns a HandlerFunc for the gin framework
func (r *Response) json(ctx *gin.Context) {
	r.Result.Cost = time.Since(ctx.GetTime("requestStartTime")).String()
	ctx.AbortWithStatusJSON(r.httpCode, r.Result)
}
