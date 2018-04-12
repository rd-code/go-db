package gdb

import (
    "strings"
    "fmt"
    "database/sql"
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
func GenerateAdd(model interface{}, columns ...string) {

}
