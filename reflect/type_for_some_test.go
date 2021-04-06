package main

import (
    "fmt"
    "reflect"
	"testing"
)

// Methods applicable only to some types, depending on Kind.
// The methods allowed for each kind are:
//
//	Int*, Uint*, Float*, Complex*: Bits
//	Array: Elem, Len
//	Chan: ChanDir, Elem
//	Func: In, NumIn, Out, NumOut, IsVariadic.
//	Map: Key, Elem
//	Ptr: Elem
//	Slice: Elem
//	Struct: Field, FieldByIndex, FieldByName, FieldByNameFunc, NumField

// 返回该类型的字位数。如果该类型的Kind不是Int、Uint、Float或Complex，会panic
func TestBits(t *testing.T) {
    fmt.Println("int bits: ", reflect.TypeOf(1).Bits())
    defer func() {
        if err := recover(); err != nil {
            fmt.Printf("panic err: %v", err)
        }
    }()
    fmt.Println("string bits will panic", reflect.TypeOf("abc").Bits())
}

// 返回一个channel类型的方向，如非通道类型将会panic
func TestChanDir(t *testing.T) {
    c := make(<-chan int, 10)
    fmt.Println(reflect.TypeOf(c).ChanDir())
}

// 如果函数类型的最后一个输入参数是"..."形式的参数，IsVariadic返回真
// 如果这样，t.In(t.NumIn() - 1)返回参数的隐式的实际类型（声明类型的切片）
// 如非函数类型将panic
func TestIsVariadic(t *testing.T) {
    fmt.Println("fmt.Println is variadic?", reflect.TypeOf(fmt.Println).IsVariadic())
}

// 返回该类型的元素类型，如果该类型的Kind不是Array、Chan、Map、Ptr或Slice，会panic
func TestElem(t *testing.T) {
    fmt.Println("map[int]string Elem: ", reflect.TypeOf(map[int]string{}).Elem().String())
    fmt.Println("[]int Elem: ", reflect.TypeOf([]int{}).Elem().String())
    fmt.Println("chan struct{} Elem: ", reflect.TypeOf(make(chan struct{})).Elem().String())
    i := 1
    fmt.Println("ptr to int Elem: ", reflect.TypeOf(&i).Elem().String())
    a := [3]string{"a", "b", "c"}
    fmt.Println("array to string Elem: ", reflect.TypeOf(a).Elem().String())
}

// 返回struct类型的第i个字段的类型，如非结构体或者i不在[0, NumField())内将会panic
func TestField(t *testing.T) {
   type s struct {
       a int
       b string
   }
   fmt.Println("s first field: ", reflect.TypeOf(s{}).Field(0))
   fmt.Println("s second field: ", reflect.TypeOf(s{}).Field(1))
   fmt.Println("s fields: ", reflect.TypeOf(s{}).NumField()) // 2
}

// 返回索引序列指定的嵌套字段的类型，
// 等价于用索引中每个值链式调用本方法，如非结构体将会panic
func TestFieldByIndex(t *testing.T) {
    type s struct {
        a int
        b string
    }
    type s1 struct {
        s
        c float32
        d uint
    }
    fmt.Println("s1 field by index {0 1}: ", reflect.TypeOf(s1{}).FieldByIndex([]int{0,1}))
    fmt.Println("s fields: ", reflect.TypeOf(s1{}).NumField()) // 3
}

// 返回该类型名为name的字段（会查找匿名字段及其子字段），
// 布尔值说明是否找到，如非结构体将panic
func TestFieldByName(t *testing.T) {
    type s1 struct {
        c int
    }
    type s struct {
        s1
        a int
        b string
    }
    // 存在
    sf, b := reflect.TypeOf(s{}).FieldByName("a")
    fmt.Println("s a field: ", sf, b) // true
    // 不存在
    sf, b = reflect.TypeOf(s{}).FieldByName("d")
    fmt.Println("s a field: ", sf, b) // false
    // 匿名字段
    sf, b = reflect.TypeOf(s{}).FieldByName("s1")
    fmt.Println("s a field: ", sf, b) // true
    // 子字段
    sf, b = reflect.TypeOf(s{}).FieldByName("c")
    fmt.Println("s a field: ", sf, b) // true
}

// 返回该类型第一个字段名满足函数match的字段，布尔值说明是否找到，如非结构体将会panic
// FieldByNameFunc 以广度优先的顺序考虑结构本身中的字段，然后考虑任何嵌入式结构中的字段
// 停在最浅的嵌套深度，其中包含一个或多个满足匹配功能的字段。如果该深度处的多个字段满足匹配功能，则它们会互相抵消，
// 并且FieldByNameFunc不返回匹配项。 此行为反映了Go在包含嵌入式字段的结构中对名称查找的处理。
func TestFileByNameFunc(t *testing.T) {
    type s struct {
        abort int
        q     string
        p     int
    }
    sf, b := reflect.TypeOf(s{}).FieldByNameFunc(func(name string) bool {
        if len(name) == 5 {
            return true
        }
        return false
    })
    fmt.Println("s abort field: ", sf, b) // true

    // p, q 都满足，但是在同一层，因此返回false
    sf, b = reflect.TypeOf(s{}).FieldByNameFunc(func(s string) bool {
        if len(s) == 1 {
            return true
        }
        return false
    })
    fmt.Println("s find field: ", sf, b) // fasle
}

// 返回func类型的第i个参数的类型，如非函数或者i不在[0, NumIn())内将会panic
func TestIn(t *testing.T) {
    f := func (int, string){}
    fmt.Println("f first arg type: ", reflect.TypeOf(f).In(0).String())
}

// 返回map类型的键的类型。如非映射类型将panic
func TestKey(t *testing.T) {
    fmt.Println("map[int]string key type: ", reflect.TypeOf(map[int]string{}).Key().String())
}

// 返回array类型的长度，如非数组类型将panic
func TestLen(t *testing.T) {
    defer func(){
        if err := recover(); err != nil {
            fmt.Printf("panic err: %v", err)
        }
    }()
    fmt.Println("array len: ", reflect.TypeOf([3]int{1,2,3}).Len())
    // panic for slice
    fmt.Println("slice len: ", reflect.TypeOf([]int{1,2,3}).Len())
}

func TestNumInOut(t *testing.T) {
    f := func(int, string) (float32, error) {
        return 0.0, nil
    }
    fmt.Println("in args: ", reflect.TypeOf(f).NumIn())
    fmt.Println("out args: ", reflect.TypeOf(f).NumOut())
    fmt.Println("out[0] arg: ", reflect.TypeOf(f).Out(0).String())
}




