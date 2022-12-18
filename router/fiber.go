package router

import (
	"strings"

	"github.com/Forceattack012/todoapidemo/todo"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

type FiberRouter struct {
	*fiber.App
}

func NewFiberRouter() *FiberRouter {
	r := fiber.New()

	config := cors.Config{
		AllowOrigins: "http://localhost:8080",
		AllowHeaders: strings.Join([]string{"Origin", "Authorization", "TransactionID"}, ","),
	}

	r.Use(cors.New(config))
	r.Use(logger.New())

	return &FiberRouter{r}
}

func (r *FiberRouter) POST(path string, handler func(todo.Context)) {
	r.App.Post(path, func(c *fiber.Ctx) error {
		handler(NewfiberContext(c))
		return nil
	})
}

func (r *FiberRouter) PUT(path string, handler func(todo.Context)) {
	r.App.Put(path, func(c *fiber.Ctx) error {
		handler(NewfiberContext(c))
		return nil
	})
}

func (r *FiberRouter) PATCH(path string, handler func(todo.Context)) {
	r.App.Patch(path, func(c *fiber.Ctx) error {
		handler(NewfiberContext(c))
		return nil
	})
}

func (r *FiberRouter) GET(path string, handler func(todo.Context)) {
	r.App.Get(path, func(c *fiber.Ctx) error {
		handler(NewfiberContext(c))
		return nil
	})
}

func (r *FiberRouter) DELETE(path string, handler func(todo.Context)) {
	r.App.Delete(path, func(c *fiber.Ctx) error {
		handler(NewfiberContext(c))
		return nil
	})
}

type fiberContext struct {
	*fiber.Ctx
}

func NewfiberContext(ctx *fiber.Ctx) *fiberContext {
	return &fiberContext{ctx}
}

func (m *fiberContext) Bind(v interface{}) error {
	return m.Ctx.BodyParser(v)
}
func (m *fiberContext) JSON(statusCode int, v interface{}) {
	m.Ctx.Status(statusCode).JSON(v)
}

func (m *fiberContext) TransactionID() string {
	return m.Ctx.GetReqHeaders()["TransactionID"]
}
func (m *fiberContext) Audience() string {
	return m.Get("aud")
}
func (m *fiberContext) GetParam(v interface{}) string {
	return m.Ctx.Params(v.(string))
}
