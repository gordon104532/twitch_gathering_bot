package ErrorHandle

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// InternalError 內部錯誤格式
type InternalError struct {
	code string
	text string
}

// ErrorCode 錯誤碼
func (e InternalError) ErrorCode() string {
	e.code = strings.TrimSpace(e.code)
	if e.code == "" {
		return "0"
	}
	return e.code
}

// ErrorText 錯誤文字
func (e InternalError) ErrorText() string {
	return e.text
}

// Error 錯誤文字
func (e InternalError) Error() string {
	return "[" + e.ErrorCode() + "]" + " " + e.ErrorText()
}

//---------

type ApiResp struct {
	Code   string      `json:"error_code"`
	Text   string      `json:"error_text"`
	Result interface{} `json:"result"`
}

func (err ApiResp) Error() string {
	return err.Text
}

type WrapperHandle func(c *gin.Context) (interface{}, error)

func ErrorWrapper(handle WrapperHandle) gin.HandlerFunc {
	return func(c *gin.Context) {
		data, r := handle(c)
		if r != nil {
			resp := r.(ApiResp)
			if resp.Code == "not json" {
				return
			}
			c.JSON(http.StatusOK, resp)

		} else {
			c.JSON(http.StatusOK, ApiResp{
				Code:   "0",
				Text:   "",
				Result: data,
			})
		}

	}
}
