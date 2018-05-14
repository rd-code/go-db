package gdb

import (
    "database/sql"
    "errors"
)

var invalidColumnNumErr = errors.New("invalid number of columns for baseType")

func queryBase(sqlStr string, args ...interface{}) (rows *sql.Rows, err error) {
    if rows, err = query(sqlStr, args...); err != nil {
        return
    }

    var columns []string
    if columns, err = rows.Columns(); err != nil {
        rows.Close()
        return
    }
    if len(columns) != 1 {
        err = invalidColumnNumErr
        rows.Close()
        return
    }
    return
}

//从数据库查询字符串，单列
func QueryString(sqlStr string, args ...interface{}) (res []string, err error) {
    var rows *sql.Rows
    if rows, err = queryBase(sqlStr, args...); err != nil {
        return
    }
    defer rows.Close()

    var columns []string
    if columns, err = rows.Columns(); err != nil {
        return
    }
    if len(columns) != 1 {
        err = invalidColumnNumErr
        return
    }

    for rows.Next() {
        var t = &sql.NullString{}
        if err = rows.Scan(t); err != nil {
            return
        }
        if t.Valid {
            res = append(res, t.String)
        }
    }
    return
}

//查询一条记录
func GetString(sqlStr string, args ...interface{}) (res string, ok bool, err error) {
    var items []string
    if items, err = QueryString(sqlStr, args...); err != nil {
        return
    }
    if len(items) == 0 {
        return
    }
    ok = true
    res = items[0]
    return
}

//从数据库查询整数信息
func QueryInt(sqlStr string, args ...interface{}) (res []int64, err error) {
    var rows *sql.Rows
    if rows, err = queryBase(sqlStr, args...); err != nil {
        return
    }
    defer rows.Close()

    var columns []string
    if columns, err = rows.Columns(); err != nil {
        return
    }
    if len(columns) != 1 {
        err = invalidColumnNumErr
        return
    }

    for rows.Next() {
        var t = &sql.NullInt64{}
        if err = rows.Scan(t); err != nil {
            return
        }
        if t.Valid {
            res = append(res, t.Int64)
        }
    }
    return
}

//查询一条整数记录
func GetInt(sqlStr string, args ...interface{}) (res int64, ok bool, err error) {
    var items []int64
    if items, err = QueryInt(sqlStr, args...); err != nil {
        return
    }
    if len(items) == 0 {
        return
    }
    ok = true
    res = items[0]
    return
}

//从数据库查询浮点数信息
func QueryFloat(sqlStr string, args ...interface{}) (res []float64, err error) {
    var rows *sql.Rows
    if rows, err = queryBase(sqlStr, args...); err != nil {
        return
    }
    defer rows.Close()

    var columns []string
    if columns, err = rows.Columns(); err != nil {
        return
    }
    if len(columns) != 1 {
        err = invalidColumnNumErr
        return
    }

    for rows.Next() {
        var t = &sql.NullFloat64{}
        if err = rows.Scan(t); err != nil {
            return
        }
        if t.Valid {
            res = append(res, t.Float64)
        }
    }
    return
}

//查询一条浮点数记录
func GetFloat(sqlStr string, args ...interface{}) (res float64, ok bool, err error) {
    var items []float64
    if items, err = QueryFloat(sqlStr, args...); err != nil {
        return
    }
    if len(items) == 0 {
        return
    }
    ok = true
    res = items[0]
    return
}

//从数据库查询浮点数信息
func QueryBool(sqlStr string, args ...interface{}) (res []bool, err error) {
    var rows *sql.Rows
    if rows, err = queryBase(sqlStr, args...); err != nil {
        return
    }
    defer rows.Close()

    var columns []string
    if columns, err = rows.Columns(); err != nil {
        return
    }
    if len(columns) != 1 {
        err = invalidColumnNumErr
        return
    }

    for rows.Next() {
        var t = &sql.NullBool{}
        if err = rows.Scan(t); err != nil {
            return
        }
        if t.Valid {
            res = append(res, t.Bool)
        }
    }
    return
}

//查询一条布尔记录
func GetBool(sqlStr string, args ...interface{}) (res bool, ok bool, err error) {
    var items []bool
    if items, err = QueryBool(sqlStr, args...); err != nil {
        return
    }
    if len(items) == 0 {
        return
    }
    ok = true
    res = items[0]
    return
}
