package router

import (
	"github.com/Forceattack012/todoapidemo/todo"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

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

// Convert MyContext to Gin Context
func NewGinContext(handler func(todo.Context)) gin.HandlerFunc {
	return func(c *gin.Context) {
		handler(NewMyContext(c))
	}
}

type MyRouter struct {
	*gin.Engine
}

func NewMyRouter() *MyRouter {
	r := gin.Default()

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{
		"http://localhost:8080",
	}
	config.AllowHeaders = []string{
		"Origin",
		"Authorization",
		"TransactionID",
	}
	return &MyRouter{r}
}

func (r *MyRouter) POST(relativePath string, handler func(todo.Context)) {
	r.Engine.POST(relativePath, NewGinContext(handler))
}
