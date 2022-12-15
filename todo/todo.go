package todo

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Todo struct {
	Title string `json:"text" binding:"required"`
	gorm.Model
}

func TableName() string {
	return "todos"
}

type Storer interface {
	New(*Todo) error
	GetById(*Todo, int) error
	GetList(*[]Todo) error
	Remove(*Todo, int) error
	Update(*Todo) error
}

type Context interface {
	Bind(interface{}) error
	JSON(int, interface{})
	TransactionID() string
	Audience() string
	GetParam(interface{}) string
}

type TodoHanlder struct {
	store Storer
}

func NewTodoHanlder(store Storer) *TodoHanlder {
	return &TodoHanlder{store: store}
}

func (h *TodoHanlder) NewTask(c Context) {
	var todo Todo
	if err := c.Bind(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if todo.Title == "sleep" {
		transactionId := c.TransactionID()
		aud := c.Audience()
		log.Println(aud, transactionId, "not allowed")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "not allowed",
		})
		return
	}

	err := h.store.New(&todo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})
}

func (h *TodoHanlder) GetTodo(c Context) {
	id := c.GetParam("id")
	nId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var todo Todo
	errDb := h.store.GetById(&todo, nId)
	if errDb != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errDb.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHanlder) GetTodoList(c Context) {
	var todo []Todo
	err := h.store.GetList(&todo)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHanlder) RemoveTask(c Context) {
	paramId := c.GetParam("id")
	id, err := strconv.Atoi(paramId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	errorRemove := h.store.Remove(&Todo{}, id)
	if errorRemove != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errorRemove.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "success",
	})
}

func (h *TodoHanlder) UpdateTask(c Context) {
	paramId := c.GetParam("id")
	id, err := strconv.Atoi(paramId)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var newTodo Todo
	if err := c.Bind(&newTodo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	var oldTodo Todo
	errorGet := h.store.GetById(&oldTodo, id)
	if errorGet != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": errorGet.Error(),
		})
		return
	}

	oldTodo.Title = newTodo.Title
	resultUpdate := h.store.Update(&oldTodo)

	if resultUpdate != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "sucess",
	})

}
