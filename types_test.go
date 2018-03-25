package gdb

import (
    "testing"
    "strings"
)

type PN struct {
    ID   int64  `db:"id"`
    Name string `db:"name"`
    Age  int64  `db:"age"`
    Date string `db:"date"`
}

var t *PN
var _ DBInterface = t

func (p *PN) TableName() string {
    return "user"
}

func TestGenerateUpdate(t *testing.T) {
    data := map[string]interface{}{
        "name": "zhangsan",
        "age":  25,
    }
    conditions := map[string]interface{}{
        "id":   1,
        "date": 19920408,
    }
    sql, _, err := GenerateUpdate(&PN{}, data, conditions)
    if err != nil {
        t.Fatalf("generate update msg happened some err")
        return
    }
    if !strings.EqualFold(sql, "UPDATE user SET name=$1, age=$2 WHERE id=$3 AND date=$4") {
        t.Fatalf("generate update msg, the sql is not expected")
        return
    }
}
