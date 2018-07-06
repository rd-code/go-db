package gdb

import (
    "testing"
    "fmt"
)

func TestSelectOrm_GenerateSql(t *testing.T) {
    sql, args, _ := NewOrm().Select().Columns("date").TableName("api-table").
        Filter("fund_code", "abc").Filter("aa", []int{1, 2, 3}, IN).Filter("dd", "44").
        OrderBy("date DESC").Limit(1).GenerateSql()
    fmt.Println(sql)
    fmt.Println(args)
}
