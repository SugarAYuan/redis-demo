package model

import "redisdemo/engine"

type Test struct {
	Id   int64  `gorm:"id"`
	Name string `gorm:"name"`
}

func (t *Test) TableName() string {
	return "test_table"
}

func (t *Test) InsertOne() error {
	return engine.MysqlDB.Table(t.TableName()).Create(t).Error
}
