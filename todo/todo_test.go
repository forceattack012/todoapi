package todo

import (
	"testing"
	"time"
)

func TestNewTask(t *testing.T) {
	handler := NewTodoHanlder(&TestDB{})
	c := &TestContext{}

	want := "not allowed"

	handler.NewTask(c)

	if want != c.v["error"] {
		t.Errorf("want %s but get %s", want, c.v["error"])
	}
}

func TestGetTodoList(t *testing.T) {
	handler := NewTodoHanlder(&TestDB{})
	c := &TestContext{}

	want := &Todo{
		Title: "test",
		ID:    0,
	}
	handler.GetTodoList(c)

	if c.todos[0].Title != want.Title {
		t.Errorf("want %s but get %s", want.Title, c.todos[0].Title)
	}
}

type TestDB struct {
}

// GetById implements Storer
func (*TestDB) GetById(*Todo, int) error {
	panic("unimplemented")
}

// GetList implements Storer
func (t *TestDB) GetList(todos *[]Todo) error {
	myTodo := Todo{
		Title:     "test",
		ID:        0,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	*todos = append(*todos, myTodo)
	return nil
}

// Remove implements Storer
func (*TestDB) Remove(*Todo, int) error {
	panic("unimplemented")
}

// Update implements Storer
func (*TestDB) Update(*Todo) error {
	panic("unimplemented")
}

func (t *TestDB) New(todo *Todo) error {
	return nil
}

type TestContext struct {
	status int
	v      map[string]interface{}
	todos  []Todo
}

func (c *TestContext) Bind(v interface{}) error {
	*v.(*Todo) = Todo{
		Title: "sleep",
	}
	return nil
}
func (c *TestContext) JSON(statusCode int, v interface{}) {
	c.status = statusCode
	if _, ok := v.(map[string]interface{}); ok {
		c.v = v.(map[string]interface{})
	} else {
		todos := v.([]Todo)
		c.todos = append(todos)
	}
}
func (c *TestContext) TransactionID() string {
	return ""
}
func (c *TestContext) Audience() string {
	return ""
}
func (c *TestContext) GetParam(interface{}) string {
	return ""
}
