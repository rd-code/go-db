package gdb

import (
    "database/sql"
    "reflect"
    "errors"
)

var invalidTypeErr = errors.New("cannot handle type")

/**
实现保存功能
 */
func Save(data interface{}, columns ...string) (result sql.Result, err error) {
    var saveSql string
    var args []interface{}
    if saveSql, args, err = GenerateAdd(data, columns...); err != nil {
        return
    }
    var db *sql.DB
    if db, err = DB(); err != nil {
        return
    }
    return db.Exec(saveSql, args...)

}

func getStructValue(rt reflect.Value) (res reflect.Value, err error) {
    for i := 0; i < 10; i++ {
        if rt.Kind() == reflect.Struct {
            res = rt
            return
        }
        rt = rt.Elem()
    }
    err = NotStructErr
    return
}
