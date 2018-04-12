package gdb

import (
    "github.com/pkg/errors"
    "strings"
    "reflect"
    "fmt"
)

type Operation int

const (
    EQUAL    Operation = iota
    NOTEQUAl
    IN
    NOTIN
)

type DBInterface interface {
    TableName() string
}

type Orm struct {
    op      Operation
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
    filter    map[string]Conditions
}

type Conditions struct {
    op    Operation
    value interface{}
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
        so.filter = make(map[string]Conditions)
    }
    var op Operation
    if len(operation) == 0 {
        op = EQUAL
    } else if len(operation) == 1 {
        op = operation[0]
    }
    so.filter[key] = Conditions{
        op:    op,
        value: value,
    }
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
        marks := make([]string, 0, len(so.filter))
        for k, v := range so.filter {
            c := &strings.Builder{}
            if _, err = c.WriteString(k); err != nil {
                return
            }
            if v.op == EQUAL {
                if _, err = c.WriteString(fmt.Sprintf("=$%d", count)); err != nil {
                    return
                }
                args = append(args, v.value)
            } else if v.op == NOTEQUAl {
                if _, err = c.WriteString(fmt.Sprintf("!=$%d", count)); err != nil {
                    return
                }
            } else {
                if v.op == IN {
                    if _, err = c.WriteString(" IN ("); err != nil {
                        return
                    }
                } else {
                    if _, err = c.WriteString(" NOT IN ("); err != nil {
                        return
                    }
                    items := v.value.([]interface{})
                    tt := make([]string, len(items))
                    for i, item := range items {
                        tt[i] = fmt.Sprintf("$%d", count)
                        args = append(args, item)
                        count += 1
                    }
                    count -= 1
                    if _, err = c.WriteString(strings.Join(tt, ", ")); err != nil {
                        return
                    }
                    if _, err = c.WriteString(")"); err != nil {
                        return
                    }
                }
            }
            count += 1

            marks = append(marks, c.String())
        }
        if _, err = sb.WriteString(strings.Join(marks, " AND ")); err != nil {
            return
        }

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
