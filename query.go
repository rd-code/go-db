package gdb

import (
    "strings"
    "reflect"
    "fmt"
    "database/sql"
    "time"
    "errors"
)

//操作 =,in,like 等等
type Operation int

const (
    EQUAL    Operation = iota
    NOTEQUAl
    IN
    NOTIN
    LIKE
    LT
    LTE
    GT
    GTE
)

//逻辑关系 add or
type Logic string

const (
    AND Logic = "AND"
    OR        = "OR"
)

type DBInterface interface {
    TableName() string
}

type Orm struct {
    columns []string
}

func NewOrm() *Orm {
    return &Orm{}
}

func (o *Orm) Select() *SelectOrm {
    return &SelectOrm{}
}

type SelectOrm struct {
    columns   []string
    model     DBInterface
    tableName string
    limit     int
    offset    int
    filter    []Conditions
    orderBy   string
    groupBy   []string
}

type Conditions struct {
    key   string
    logic Logic
    op    Operation
    value interface{}
}

func (c *Conditions) Generate(count *int) (logic string, res string, args []interface{}, err error) {
    sb := strings.Builder{}
    if _, err = sb.WriteString(c.key); err != nil {
        return
    }
    switch c.op {
    case EQUAL:
        if _, err = sb.WriteString(fmt.Sprintf("=$%d", *count)); err != nil {
            return
        }
        args = append(args, c.value)
    case NOTEQUAl:
        if _, err = sb.WriteString(fmt.Sprintf("!=$%d", *count)); err != nil {
            return
        }
        args = append(args, c.value)
    case LIKE:
        if _, err = sb.WriteString(fmt.Sprintf(" LIKE $%d", *count)); err != nil {
            return
        }
        args = append(args, c.value)
    case LT:
        if _, err = sb.WriteString(fmt.Sprintf("<$%d", *count)); err != nil {
            return
        }
        args = append(args, c.value)
    case LTE:
        if _, err = sb.WriteString(fmt.Sprintf("<=$%d", *count)); err != nil {
            return
        }
        args = append(args, c.value)
    case GT:
        if _, err = sb.WriteString(fmt.Sprintf(">$%d", *count)); err != nil {
            return
        }
        args = append(args, c.value)
    case GTE:
        if _, err = sb.WriteString(fmt.Sprintf(">=$%d", *count)); err != nil {
            return
        }
        args = append(args, c.value)
    case IN, NOTIN:
        if c.op == IN {
            if _, err = sb.WriteString(" IN ("); err != nil {
                return
            }

        } else {
            if _, err = sb.WriteString(" NOT IN ("); err != nil {
                return
            }
        }
        items := c.value.([]interface{})
        marks := make([]string, len(items))
        for i, item := range items {
            marks[i] = fmt.Sprintf("$%d", *count)
            args = append(args, item)
            *count += 1
        }
        *count -= 1
        if _, err = sb.WriteString(strings.Join(marks, ", ")); err != nil {
            return
        }
        if _, err = sb.WriteString(")"); err != nil {
            return
        }

    }
    *count += 1
    res = sb.String()

    logic = string(c.logic)

    return
}

func (so *SelectOrm) Columns(columns ...string) *SelectOrm {
    so.columns = columns
    return so
}

func (so *SelectOrm) Model(model DBInterface) *SelectOrm {
    so.model = model
    return so
}

func (so *SelectOrm) TableName(tableName string) *SelectOrm {
    so.tableName = tableName
    return so
}

func (so *SelectOrm) OrderBy(orderBy string) *SelectOrm {
    so.orderBy = orderBy
    return so
}

func (so *SelectOrm) GroupBy(columns ...string) *SelectOrm {
    so.groupBy = columns
    return so
}

func (so *SelectOrm) Limit(limit int) *SelectOrm {
    so.limit = limit
    return so
}
func (so *SelectOrm) Offset(offset int) *SelectOrm {
    so.offset = offset
    return so
}

var InvalidOperationNumErr = errors.New("the operate number is invalid")

func (so *SelectOrm) Filter(key string, value interface{}, operation ...Operation) *SelectOrm {
    if so.filter == nil {
        so.filter = make([]Conditions, 0, 1)
    }
    var op Operation
    if len(operation) == 0 {
        op = EQUAL
    } else if len(operation) == 1 {
        op = operation[0]
    }
    so.filter = append(so.filter, Conditions{
        key:   key,
        logic: AND,
        op:    op,
        value: value,
    })
    return so
}

func (so *SelectOrm) FilterOr(key string, value interface{}, operation ...Operation) *SelectOrm {
    if so.filter == nil {
        so.filter = make([]Conditions, 0, 1)
    }
    var op Operation
    if len(operation) == 0 {
        op = EQUAL
    } else if len(operation) == 1 {
        op = operation[0]
    }
    so.filter = append(so.filter, Conditions{
        key:   key,
        logic: OR,
        op:    op,
        value: value,
    })
    return so
}

func (so *SelectOrm) In(key string, value ...interface{}) *SelectOrm {
    return so.Filter(key, value, IN)
}

var invalidSelectErr = errors.New("invalid select error")

func (so *SelectOrm) GenerateSql() (sql string, args []interface{}, err error) {
    if so.model == nil && (len(so.tableName) == 0 || len(so.columns) == 0) {
        err = invalidSelectErr
        return
    }
    sb := strings.Builder{}
    if _, err = sb.WriteString("SELECT "); err != nil {
        return
    }
    var columns []string
    if len(so.columns) != 0 {
        columns = so.columns
    } else {
        columns = GetColumns(so.model)
    }
    if _, err = sb.WriteString(strings.Join(columns, ", ")); err != nil {
        return
    }
    if _, err = sb.WriteString(" FROM "); err != nil {
        return
    }
    var tableName string
    if len(so.tableName) != 0 {
        tableName = so.tableName
    } else {
        tableName = so.model.TableName()
    }
    if _, err = sb.WriteString(tableName); err != nil {
        return
    }

    count := 1
    if len(so.filter) != 0 {
        if _, err = sb.WriteString(" WHERE "); err != nil {
            return
        }
        condition := so.filter[0]

        var mark string
        var tmpArgs []interface{}
        if _, mark, tmpArgs, err = condition.Generate(&count); err != nil {
            return
        }
        if _, err = sb.WriteString(mark); err != nil {
            return
        }
        args = append(args, tmpArgs...)

        for i := 1; i < len(so.filter); i++ {
            condition = so.filter[i]
            var logic string
            if logic, mark, tmpArgs, err = condition.Generate(&count); err != nil {
                return
            }
            if _, err = sb.WriteString(" "); err != nil {
                return
            }
            if _, err = sb.WriteString(logic); err != nil {
                return
            }
            if _, err = sb.WriteString(" "); err != nil {
                return
            }
            if _, err = sb.WriteString(mark); err != nil {
                return
            }
            args = append(args, tmpArgs...)
        }
    }

    if len(so.groupBy) != 0 {
        if _, err = sb.WriteString(" GROUP BY "); err != nil {
            return
        }
        if _, err = sb.WriteString(strings.Join(so.groupBy, ", ")); err != nil {
            return
        }
    }
    if len(so.orderBy) != 0 {
        if _, err = sb.WriteString(" ORDER BY "); err != nil {
            return
        }
        if _, err = sb.WriteString(so.orderBy); err != nil {
            return
        }
    }

    if so.limit > 0 {
        if _, err = sb.WriteString(fmt.Sprintf(" LIMIT $%d", count)); err != nil {
            return
        }
        args = append(args, so.limit)
        count += 1
    }
    if so.offset > 0 {
        if _, err = sb.WriteString(fmt.Sprintf(" OFFSET $%d", count)); err != nil {
            return
        }
        args = append(args, so.offset)
        count += 1
    }

    sql = sb.String()
    return
}

func GetColumns(model DBInterface) (columns []string) {
    fields := cacheTypeFileds(reflect.TypeOf(model))
    columns = make([]string, 0, len(fields))
    for _, f := range fields {
        if f.valid {
            columns = append(columns, f.tag)
        }
    }
    return
}

func (so *SelectOrm) QueryString() ([]string, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return nil, err
    }
    return QueryString(sql, args...)
}

func (so *SelectOrm) GetString() (string, bool, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return "", false, err
    }
    return GetString(sql, args...)
}

func (so *SelectOrm) QueryInt() ([]int64, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return nil, err
    }
    return QueryInt(sql, args...)
}

func (so *SelectOrm) GetInt() (int64, bool, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return 0, false, err
    }
    return GetInt(sql, args...)
}

func (so *SelectOrm) QueryFloat() ([]float64, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return nil, err
    }
    return QueryFloat(sql, args...)
}

func (so *SelectOrm) GetFloat() (float64, bool, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return 0, false, err
    }
    return GetFloat(sql, args...)
}

func (so *SelectOrm) QueryBool() ([]bool, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return nil, err
    }
    return QueryBool(sql, args...)
}

func (so *SelectOrm) GetBool() (bool, bool, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return false, false, err
    }
    return GetBool(sql, args...)
}

func (so *SelectOrm) QueryMap() ([]map[string]interface{}, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return nil, err
    }
    return QueryMap(sql, args...)
}

func (so *SelectOrm) GetMap() (map[string]interface{}, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return nil, err
    }
    return GetMap(sql, args...)
}

func (so *SelectOrm) Query() ([]interface{}, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return nil, err
    }
    return Query(so.model, sql, args...)
}

func (so *SelectOrm) Get() (interface{}, error) {
    sql, args, err := so.GenerateSql()
    if err != nil {
        return nil, err
    }
    return Get(so.model, sql, args...)
}

var InvalidFieldNumErr = errors.New("the number of struct field is invalid")
var FieldTypeErr = errors.New("the field of sturct type is invalid")
var timeType = reflect.TypeOf(time.Time{})

func Query(model interface{}, sqlStr string, args ...interface{}) (res []interface{}, err error) {
    switch model.(type) {
    case int, int8, int16, int32, int64, *int, *int8, *int16, *int32, *int64,
    uint, uint8, uint16, uint32, uint64, *uint, *uint8, *uint16, *uint32, *uint64:
        var items []int64
        if items, err = QueryInt(sqlStr, args...); err != nil {
            return
        }
        res = make([]interface{}, len(items))
        for i := range items {
            res[i] = items[i]
        }
        return
    case string, *string:
        var items []string
        if items, err = QueryString(sqlStr, args...); err != nil {
            return
        }
        res = make([]interface{}, len(items))
        for i := range items {
            res[i] = items[i]
        }
        return
    case bool, *bool:
        var items []bool
        if items, err = QueryBool(sqlStr, args...); err != nil {
            return
        }
        res = make([]interface{}, len(items))
        for i := range items {
            res[i] = items[i]
        }
        return
    case float64, *float64, float32, *float32:
        var items []float64
        if items, err = QueryFloat(sqlStr, args...); err != nil {
            return
        }
        res = make([]interface{}, len(items))
        for i := range items {
            res[i] = items[i]
        }

        return
    }

    fields := cacheTypeFileds(reflect.TypeOf(model))
    if len(fields) == 0 {
        err = InvalidFieldNumErr
        return
    }

    var rows *sql.Rows
    if rows, err = query(sqlStr, args...); err != nil {
        return
    }
    defer rows.Close()

    var columns []string
    if columns, err = rows.Columns(); err != nil {
        return
    }
    rt := reflect.TypeOf(model)
    if rt, err = GetStructType(rt); err != nil {
        return
    }
    columnField := typeFileds(rt)

    for rows.Next() {
        items := make([]interface{}, len(columns))
        for i, column := range columns {
            if field, ok := columnField[column]; ok {
                switch field.typ.Kind() {
                case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
                    reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
                    items[i] = &sql.NullInt64{}
                case reflect.String:
                    items[i] = &sql.NullString{}
                case reflect.Bool:
                    items[i] = &sql.NullBool{}
                case reflect.Float64, reflect.Float32:
                    items[i] = &sql.NullFloat64{}
                case reflect.Slice:
                    if field.typ.Elem().Kind() == reflect.Uint8 {
                        t := &[]byte{}
                        items[i] = t
                    } else {
                        items[i] = &sql.NullString{}
                    }
                default:
                    items[i] = &sql.NullString{}
                }
            } else {
                items[i] = &sql.NullString{}
            }
        }
        if err = rows.Scan(items...); err != nil {
            return
        }
        rv := reflect.New(rt)
        for i, column := range columns {
            if field, ok := columnField[column]; ok {
                switch field.typ.Kind() {
                case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
                    reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
                    t := items[i].(*sql.NullInt64)
                    if t.Valid {
                        rv.Elem().FieldByName(field.name).SetInt(t.Int64)
                    }
                case reflect.String:
                    t := items[i].(*sql.NullString)
                    if t.Valid {
                        rv.Elem().FieldByName(field.name).SetString(t.String)
                    }
                case reflect.Bool:
                    t := items[i].(*sql.NullBool)
                    if t.Valid {
                        rv.Elem().FieldByName(field.name).SetBool(t.Bool)
                    }
                case reflect.Float32, reflect.Float64:
                    t := items[i].(*sql.NullFloat64)
                    if t.Valid {
                        rv.Elem().FieldByName(field.name).SetFloat(t.Float64)
                    }
                case reflect.Slice:
                    if field.typ.Elem().Kind() == reflect.Uint8 {
                        value := *items[i].(*[]byte)
                        rv.Elem().FieldByName(field.name).SetBytes(value)
                    } else {
                        err = FieldTypeErr
                        return
                    }

                default:
                    if field.typ == timeType {
                        t := items[i].(*sql.NullString)
                        if t.Valid {
                            var tim time.Time
                            var layout = time.RFC3339
                            if len(field.format) != 0 {
                                layout = field.format
                            }
                            if tim, err = time.Parse(layout, t.String); err != nil {
                                return
                            }
                            rv.Elem().FieldByName(field.name).Set(reflect.ValueOf(tim))
                        }
                    } else {
                        err = FieldTypeErr
                        return
                    }

                }
            }
        }

        res = append(res, rv.Elem().Interface())
    }

    return
}

func Get(model interface{}, sqlStr string, args ...interface{}) (interface{}, error) {
    items, err := Query(model, sqlStr, args...)
    if err != nil {
        return nil, err
    }
    if len(items) == 0 {
        return nil, nil
    }
    return items[0], nil

}

var NotStructErr = errors.New("cannot get struct")

func GetStructType(rt reflect.Type) (res reflect.Type, err error) {
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
