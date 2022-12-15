package todo

import "github.com/gin-gonic/gin"

type MyContext struct {
	*gin.Context
}

func NewMyContext(ctx *gin.Context) *MyContext {
	return &MyContext{Context: ctx}
}

func (m *MyContext) Bind(v interface{}) error {
	return m.Context.ShouldBindJSON(v)
}
func (m *MyContext) JSON(statusCode int, v interface{}) {
	m.Context.JSON(statusCode, v)
}

func (m *MyContext) TransactionID() string {
	return m.Request.Header.Get("TransactionID")
}
func (m *MyContext) Audience() string {
	if aud, ok := m.Get("aud"); ok {
		if s, ok := aud.(string); ok {
			return s
		}
	}

	return ""
}
func (m *MyContext) GetParam(v interface{}) string {
	return m.Context.Param(v.(string))
}

//Convert MyContext to Gin Context
func NewGinContext(handler func(Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewMyContext(c))
	}
}
