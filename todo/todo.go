package todo

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"github.com/Forceattack012/todoapidemo/auth"
)

type Todo struct {
	Title string `json:"text" binding:"required"`
	gorm.Model
}

func TableName() string {
	return "todos"
}

type TodoHanlder struct {
	db *gorm.DB
}

func NewTodoHanlder(db *gorm.DB) *TodoHanlder {
	return &TodoHanlder{db: db}
}

func (h *TodoHanlder) NewTask(c *gin.Context) {

	token := c.Request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(token, "Bearer ")

	if err := auth.Protect(tokenString); err != nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var todo Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	r := h.db.Create(&todo)
	if err := r.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"ID": todo.Model.ID,
	})
}

func (h *TodoHanlder) GetTodo(c *gin.Context) {
	id := c.Param("id")
	nId, err := strconv.Atoi(id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var todo Todo
	result := h.db.Find(&todo, nId)
	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}

func (h *TodoHanlder) GetTodoList(c *gin.Context) {
	var todo []Todo
	result := h.db.Find(&todo)

	if err := result.Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, todo)
}
