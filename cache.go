package gdb

import (
    "sync/atomic"
    "sync"
    "reflect"
    "strings"
)

const (
    DB_TAG    = "db"
    EMPTY_TAG = "-"
)

var fieldCache struct {
    value atomic.Value // map[reflect.Type][]field
    mu    sync.Mutex   // used only by writers
}

type field struct {
    name string
    tag  string

    valid  bool
    index  []int
    typ    reflect.Type
    format string
}

func cacheTypeFileds(t reflect.Type) (map[string]field) {
    m, _ := fieldCache.value.Load().(map[reflect.Type]map[string]field)
    f := m[t]
    if f != nil {
        return f
    }
    f = typeFileds(t)
    if f == nil {
        f = map[string]field{}
    }
    fieldCache.mu.Lock()
    m, _ = fieldCache.value.Load().(map[reflect.Type]map[string]field)
    newM := make(map[reflect.Type]map[string]field, len(m)+1)
    for k, v := range m {
        newM[k] = v
    }
    newM[t] = f
    fieldCache.value.Store(newM)
    fieldCache.mu.Unlock()

    return f
}

func typeFileds(t reflect.Type) (map[string]field) {
    for {
        if t.Kind() == reflect.Ptr || t.Kind() == reflect.Slice {
            t = t.Elem()
        } else {
            break
        }
    }
    if t.Kind() != reflect.Struct {
        return nil
    }
    num := t.NumField()
    res := make(map[string]field)
    var ok bool
    var tag string
    for i := 0; i < num; i++ {
        f := t.Field(i)
        var tmp field
        if tag, ok = getTag(&f); ok {
            tmp.valid = true
            if f.Type == timeType {
                if strings.Contains(tag, ";") {
                    items := strings.Split(tag, ";")
                    tag, tmp.format = items[0], items[1]
                }
            }
            tmp.name = f.Name
            tmp.tag = tag
            tmp.typ = f.Type
            tmp.index = f.Index
        } else {
            tmp.valid = false
        }
        res[tag] = tmp
    }
    return res
}

func getTag(f *reflect.StructField) (tag string, ok bool) {
    tag = f.Tag.Get(DB_TAG)
    if len(tag) == 0 {
        tag = f.Name
        ok = true
    } else if tag == EMPTY_TAG {
        ok = false
    } else {
        ok = true
    }
    return
}
