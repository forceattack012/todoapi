package todo

import "gorm.io/gorm"

type GormStore struct {
	db *gorm.DB
}

func NewGormStore(db *gorm.DB) *GormStore {
	return &GormStore{db: db}
}

func (s *GormStore) New(todo *Todo) error {
	return s.db.Create(todo).Error
}

func (s *GormStore) GetById(todo *Todo, id int) error {
	return s.db.Find(todo, id).Error
}

func (s *GormStore) GetList(todos *[]Todo) error {
	return s.db.Find(todos).Error
}

func (s *GormStore) Remove(todo *Todo, id int) error {
	return s.db.Delete(todo, id).Error
}

func (s *GormStore) Update(todo *Todo) error {
	return s.db.Updates(todo).Error
}
