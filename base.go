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
