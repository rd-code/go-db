package gdb

import (
    "reflect"
    "time"
    "database/sql"
    "encoding/json"
    "errors"
    "strings"
)

var unknowTypeErr = errors.New("unknown type")

func QueryMap(sqlStr string, args ...interface{}) (res []map[string]interface{}, err error) {
    var rows *sql.Rows

    if rows, err = query(sqlStr, args...); err != nil {
        return
    }
    defer rows.Close()

    var columns []string
    if columns, err = rows.Columns(); err != nil {
        return
    }
    var columnTypes []*sql.ColumnType
    if columnTypes, err = rows.ColumnTypes(); err != nil {
        return
    }

    for rows.Next() {
        items := make([]interface{}, len(columnTypes))
        for i := range columnTypes {
            columnType := columnTypes[i]
            var item interface{}
            switch columnType.ScanType().Kind() {
            case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
                reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
                item = &sql.NullInt64{}
            case reflect.Float64, reflect.Float32:
                item = &sql.NullFloat64{}
            case reflect.Bool:
                item = &sql.NullBool{}
            case reflect.String:
                item = &sql.NullString{}
            case reflect.Interface:
                if strings.HasPrefix(columnType.DatabaseTypeName(), "FLOAT") {
                    item = &sql.NullFloat64{}
                } else {
                    item = &[]byte{}
                }
            default:
                if columnType.ScanType() == timeType {
                    item = &sql.NullString{}
                } else {
                    err = unknowTypeErr
                    return
                }
            }
            items[i] = item
        }
        if err = rows.Scan(items...); err != nil {
            return
        }

        data := make(map[string]interface{})
        for i := range items {
            columnType := columnTypes[i]
            if columnType.ScanType() == timeType {
                item := items[i].(*sql.NullString)
                if item.Valid {
                    var t time.Time
                    if t, err = time.Parse(time.RFC3339, item.String); err != nil {
                        return
                    }
                    data[columns[i]] = t
                }
            } else {
                switch columnType.ScanType().Kind() {

                case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
                    reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
                    item := items[i].(*sql.NullInt64)
                    if item.Valid {
                        data[columns[i]] = item.Int64
                    }
                case reflect.Float64, reflect.Float32:
                    item := items[i].(*sql.NullFloat64)
                    if item.Valid {
                        data[columns[i]] = item.Float64
                    }
                case reflect.Bool:
                    item := items[i].(*sql.NullBool)
                    if item.Valid {
                        data[columns[i]] = item.Bool
                    }
                case reflect.String:
                    item := items[i].(*sql.NullString)
                    if item.Valid {
                        data[columns[i]] = item.String
                    }
                case reflect.Interface:
                    if strings.HasPrefix(columnType.DatabaseTypeName(), "FLOAT") {
                        item := items[i].(*sql.NullFloat64)
                        if item.Valid {
                            data[columns[i]] = item.Float64
                        }
                    } else {
                        item := items[i].(*[]byte)
                        if columnType.DatabaseTypeName() == "JSONB" {
                            var t interface{}
                            if t, err = convertByteToJson(*item); err != nil {
                                return
                            }
                            if t != nil {
                                data[columns[i]] = t
                            }
                        } else {
                            if len(*item) > 0 {
                                data[columns[i]] = item
                            }
                        }
                    }
                }
            }
        }
        res = append(res, data)
    }
    return
}

func convertByteToJson(array []byte) (interface{}, error) {
    if len(array) == 0 {
        return nil, nil
    }
    var res interface{}
    res = map[string]interface{}{}
    err := json.Unmarshal(array, &res)
    if err == nil {
        return res, nil
    }
    res = []interface{}{}
    err = json.Unmarshal(array, &res)
    return res, err
}

func GetMap(sqlStr string, args ...interface{}) (res map[string]interface{}, err error) {
    var items []map[string]interface{}
    if items, err = QueryMap(sqlStr, args...); err != nil {
        return
    }
    if len(items) == 0 {
        return
    }
    res = items[0]
    return
}
