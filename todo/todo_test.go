package todo

import "testing"

func TestNewTask(t *testing.T) {
	handler := NewTodoHanlder(&TestDB{})
	c := &TestContext{}

	want := "not allowed"

	handler.NewTask(c)

	if want != c.v["error"] {
		t.Errorf("want %s but get %s", want, c.v["error"])
	}
}

type TestDB struct{}

// GetById implements Storer
func (*TestDB) GetById(*Todo, int) error {
	panic("unimplemented")
}

// GetList implements Storer
func (*TestDB) GetList(*[]Todo) error {
	panic("unimplemented")
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
}

func (c *TestContext) Bind(v interface{}) error {
	*v.(*Todo) = Todo{
		Title: "sleep",
	}
	return nil
}
func (c *TestContext) JSON(statusCode int, v interface{}) {
	c.status = statusCode
	c.v = v.(map[string]interface{})
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
