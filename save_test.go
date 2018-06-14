package gdb

import (
    "testing"
    "strconv"
    "fmt"
)

/**
 * DESCRIPTION:
 *
 * @author rd
 * @create 2018-06-14 19:30
 **/

type Test struct {
    Name string `db:"name"`
    Age  int    `db:"age"`
}

func (*Test)TableName()string  {
    return "api_test"
}

func TestSave(t *testing.T) {
    items := make([]Test, 2)
    for i := 0; i < 2; i++ {
        items[i] = Test{
            Name: "--" + strconv.Itoa(i),
            Age:  i + 20,
        }
    }
    fmt.Println(GenerateAdd(items))
}
