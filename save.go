package gdb

import (
    "database/sql"
    "reflect"
    "errors"
    "strings"
    "fmt"
)

var invalidTypeErr = errors.New("cannot handle type")

/**
实现保存功能
 */
func Save(data interface{}, columns ...string) (result sql.Result, err error) {
    rt := reflect.TypeOf(data)
    switch rt.Kind() {
    case reflect.Struct:
    case reflect.Slice:
    default:
        err = invalidTypeErr
        return
    }
    return
}

func multiSave(items []DBInterface, columns ...string) (result sql.Result, err error) {
    if len(items) == 0 {
        return
    }
    typeFiled := cacheTypeFileds(reflect.TypeOf(items[0]))
    sb := &strings.Builder{}
    if _, err = sb.WriteString("INSERT INTO "); err != nil {
        return
    }
    if _, err = sb.WriteString(items[0].TableName()); err != nil {
        return
    }
    if _, err = sb.WriteString(" ("); err != nil {
        return
    }
    if _, err = sb.WriteString(strings.Join(columns, ", ")); err != nil {
        return
    }
    if _, err = sb.WriteString(") VALUES"); err != nil {
        return
    }

    var template string
    {
        tmp := make([]string, len(columns))
        for i := range columns {
            tmp[i] = "$%d"
        }
        template = "(" + strings.Join(tmp, ", ") + ")"
    }

    count := 1
    args := make([]interface{}, 0, len(items)*len(columns))
    marks := make([]string, 0, len(items))

    for _, item := range items {
        rv := reflect.ValueOf(item)
        flags := make([]interface{}, 0, len(columns))
        for _, column := range columns {
            flags = append(flags, count)
            v := rv.FieldByIndex(typeFiled[column].index)
            args = append(args, v.Interface())
            count += 1
        }
        marks = append(marks, fmt.Sprintf(template, flags...))
    }
    if _, err = sb.WriteString(strings.Join(marks, ", ")); err != nil {
        return
    }
    var db *sql.DB
    if db, err = DB(); err != nil {
        return
    }
    result, err = db.Exec(sb.String(), args...)
    return
}
