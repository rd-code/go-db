package gdb

import (
    "strings"
    "fmt"
    "database/sql"
    "reflect"
)

//实现更新操作
//data 需要更新数据
//conditions 更新条件
func Update(model DBInterface, data, conditions map[string]interface{}) (sql.Result, error) {
    sql, args, err := GenerateUpdate(model, data, conditions)
    if err != nil {
        return nil, err
    }
    db, err := DB()
    if err != nil {
        return nil, err
    }
    return db.Exec(sql, args...)
}

//生成数据库的update语句以及相应参数
func GenerateUpdate(model DBInterface, data, conditions map[string]interface{}) (sql string, args []interface{}, err error) {
    builder := &strings.Builder{}
    if _, err = builder.WriteString("UPDATE "); err != nil {
        return
    }
    if _, err = builder.WriteString(model.TableName()); err != nil {
        return
    }
    if _, err = builder.WriteString(" SET "); err != nil {
        return
    }
    marks := make([]string, 0, len(data))
    args = make([]interface{}, 0, len(data)+len(conditions))
    count := 1
    for k, v := range data {
        marks = append(marks, fmt.Sprintf("%s=$%d", k, count))
        args = append(args, v)
        count += 1
    }
    if _, err = builder.WriteString(strings.Join(marks, ", ")); err != nil {
        return
    }
    if len(conditions) == 0 {
        return
    }
    if _, err = builder.WriteString(" WHERE "); err != nil {
        return
    }
    marks = make([]string, 0, len(conditions))
    for k, v := range conditions {
        marks = append(marks, fmt.Sprintf("%s=$%d", k, count))
        args = append(args, v)
        count += 1
    }
    if _, err = builder.WriteString(strings.Join(marks, " AND ")); err != nil {
        return
    }
    sql = builder.String()
    return
}

//生成数据库的Insert语句以及相应参数
//model需要实现DBInterface，或者model为数组，且里面的每个元素实现DBInterface
func GenerateAdd(data interface{}, columns ...string) (saveSql string, args []interface{}, err error) {
    rt := reflect.TypeOf(data)
    var items []interface{}
    switch rt.Kind() {
    case reflect.Struct, reflect.Ptr:
        if _, ok := data.(DBInterface); ok {
            items = append(items, data)
        } else {
            err = invalidTypeErr
            return
        }
    case reflect.Slice:
        rv := reflect.ValueOf(data)
        for i := 0; i < rv.Len(); i++ {
            items = append(items, rv.Index(i).Addr().Interface())
        }
    default:
        err = invalidTypeErr
        return
    }
    return generateMultiSave(items, columns...)
}

func generateMultiSave(items []interface{}, columns ...string) (saveSql string, args []interface{}, err error) {
    if len(items) == 0 {
        return
    }
    typeFiled := cacheTypeFileds(reflect.TypeOf(items[0]))
    if len(columns) == 0 {
        columns = make([]string, 0, len(typeFiled))
        for k := range typeFiled {
            columns = append(columns, k)
        }
    }
    sb := &strings.Builder{}
    if _, err = sb.WriteString("INSERT INTO "); err != nil {
        return
    }
    if _, ok := items[0].(DBInterface); !ok {
        err = invalidTypeErr
        return
    }
    if _, err = sb.WriteString(items[0].(DBInterface).TableName()); err != nil {
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
    args = make([]interface{}, 0, len(items)*len(columns))
    marks := make([]string, 0, len(items))
    for _, item := range items {
        rv := reflect.ValueOf(item)
        if rv, err = getStructValue(rv); err != nil {
            return
        }
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
    saveSql = sb.String()
    return
}
