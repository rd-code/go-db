package gdb

import (
    "testing"
    "fmt"
    "time"
    "reflect"
)

type Data struct {
    Name string    `db:"name""`
    Age  int       `db:"age"`
    Date time.Time `db:"date;2006-01-02"`
}

func (d *Data) TableName() string {
    return "api_test"
}
func TestCacheTypeFileds(t *testing.T) {
    d := &Data{
        Name: "rd",
        Age:  12,
        Date: time.Now(),
    }
    items := cacheTypeFileds(reflect.TypeOf(&d))
    fmt.Println(len(items))
    for _, item := range items {
        fmt.Printf("%+v\n", item)
    }

    fmt.Println(NewOrm().Select().Model(d).Filter("name", "123").Filter("age", []interface{}{1, 2, 3}, NOTIN).GenerateSql())
}
