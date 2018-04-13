package gdb

import (
    "testing"
    "github.com/rd-pn/go-db"
    "fmt"
)

func TestSelectOrm_GenerateSql(t *testing.T) {
    sql, args, _ := gdb.NewOrm().Select().Columns("date").TableName("api-table").
        Filter("fund_code", "abc").OrderBy("date DESC").Limit(1).GenerateSql()
    fmt.Println(sql)
    fmt.Println(args)
}
