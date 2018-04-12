package gdb

import (
    "testing"
    "fmt"
)

type Data struct {
    Name string `db:"name""`
    Age  int    `db:"age"`
}

func (d *Data) TableName() string {
    return "api_test"
}
func TestCacheTypeFileds(t *testing.T) {
    d := &Data{
        Name: "rd",
        Age:  12,
    }
    fmt.Println(NewOrm().Select().Model(d).Filter("name", "123").Filter("age", []interface{}{1, 2, 3}, NOTIN).GenerateSql())
}
