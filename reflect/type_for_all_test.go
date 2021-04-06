package main

//test Type Methods applicable to all types.

import (
    "fmt"
    "reflect"
    "testing"
)

type alignS struct {
    int
    string
    float32
}

// 返回当从内存中申请一个该类型值时，会对齐的字节数
func TestAlign(t *testing.T) {
    fmt.Println("int align: ", reflect.TypeOf(1).Align(),
                    "field align: ", reflect.TypeOf(1).FieldAlign()) // 8 8
    fmt.Println("struct align: ", reflect.TypeOf(alignS{}).Align(),
                    "field align", reflect.TypeOf(alignS{}).FieldAlign())// 8 8
}

// 返回该类型方法集中的第i个方法，i不在[0, NumMethod())范围内时，将导致panic
// 对于非接口类型T或* T，返回的Method的Type和Func字段描述了一个函数，该函数的第一个参数是接收方，并且只能访问导出的方法。
// 对于接口类型，返回的Method的Type字段给出方法签名，没有接收方，而Func字段为nil。
func TestMethod(t *testing.T) {
    defer func(){
        if err := recover(); err != nil {
            fmt.Printf("Method panic: %v\n", err)
            // panic recover 可以无限套？
            defer func(){
                if err := recover(); err != nil {
                    fmt.Printf("panic again: %v\n", err)
                }
            }()
            fmt.Println("panic1: ", reflect.TypeOf(&numM{}).Method(2))
        }
    }()
    fmt.Println("numM String method: ", reflect.TypeOf(&numM{}).Method(0).Type.String())
    fmt.Println("method interface method: ", reflect.TypeOf((*method)(nil)).Elem().Method(0).Type.String())
    fmt.Println("interface would not panic: ", reflect.TypeOf((*method)(nil)).Elem().Method(100))
    fmt.Println("after panic")
    fmt.Println("panic1: ", reflect.TypeOf(&numM{}).Method(2))
}

// MethodByName 非接口类型只能访问导出方法
func TestMethodByName(t *testing.T) {
    // struct export method
    m, b := reflect.TypeOf(&numM{}).MethodByName("String")
    fmt.Println("numM String method: ", m, b)
    // struct unexported method
    m, b = reflect.TypeOf(&numM{}).MethodByName("string")
    fmt.Println("numM String method: ", m, b)
    // interface unexported method
    m, b = reflect.TypeOf((*method)(nil)).Elem().MethodByName("a")
    fmt.Println("numM String method: ", m, b)
    // interface export method, Method.PkgPath is empty string
    m, b = reflect.TypeOf((*method)(nil)).Elem().MethodByName("A")
    fmt.Println("numM String method: ", m, b)
}

// NumMethod 返回该类型的方法集中方法的数目
// 匿名字段的方法会被计算；主体类型的方法会屏蔽匿名字段的同名方法；
// 匿名字段导致的歧义方法会滤除
// NumMethod仅对接口类型计算未导出的方法
type numM struct {
    child
    a int
    b string
}

// +1
func (n *numM) String() string{
    return n.b
}

// NumMethod仅对接口类型计算未导出的方法 +0
func (n *numM) string() string {
    return n.b
}

type child struct {
    s string
}

// 主体类型的方法会屏蔽匿名字段的同名方法； +0
func (c *child) String() string{
    return c.s
}

// NumMethod仅对接口类型计算未导出的方法 +0
func (c *child) string() string {
    return c.s
}

// 匿名字段的方法会被计算 +1
func (c *child) A() string {
    return c.s
}

// NumMethod仅对接口类型计算未导出的方法 +2
type method interface {
    A(string) int
    a(string) int
    b(int) string
}

func TestNumMethod(t *testing.T) {
    fmt.Println("*numM methods: ", reflect.TypeOf(&numM{}).NumMethod()) // 2
    fmt.Println("numM methods: ", reflect.TypeOf(numM{}).NumMethod()) // 0
    fmt.Println("*interface method: ", reflect.TypeOf((*method)(nil)).NumMethod()) // 0
    fmt.Println("interface method: ", reflect.TypeOf((*method)(nil)).Elem().NumMethod()) // 2
}

// Name返回该类型在自身包内的类型名，如果是未命名类型会返回""
func TestName(t *testing.T) {
    fmt.Println("int name str: ", reflect.TypeOf(1).Name()) // # => int
    fmt.Println("slice name str: ", reflect.TypeOf([]int{}).Name()) // # => ""
    fmt.Println("map name str: ", reflect.TypeOf(map[string]interface{}{}).Name()) // # => ""
    fmt.Println("struct alignS name str: ", reflect.TypeOf(alignS{}).Name()) // # => alignS
}

// PkgPath返回类型的包路径，即明确指定包的import路径，如"encoding/base64"
// 如果类型为内建类型(string, error)或未命名类型(*T, struct{}, []int)，会返回""
type A string
type B struct{}

func TestPkgPath(t *testing.T) {
    // internal type
    fmt.Println("int pkgPath: ", reflect.TypeOf(1).PkgPath()) // ""
    fmt.Println("no named type pkgPath: ", reflect.TypeOf(struct{}{}).PkgPath()) // ""
    fmt.Println("named type pkgPath: ", reflect.TypeOf(alignS{}).PkgPath()) // main
    // alias
    fmt.Println("alias to internal type type pkgPath: ", reflect.TypeOf(A("")).PkgPath()) // main
    fmt.Println("alias to undefined type pkgPath: ", reflect.TypeOf(B{}).PkgPath()) // main
}

// 返回要保存一个该类型的值需要多少字节；类似unsafe.Sizeof
func TestSize(t *testing.T) {
    fmt.Println("int size: ",reflect.TypeOf(1).Size()) // 8
    fmt.Println("float64 str: ", reflect.TypeOf(0.01).String(), "size: ", reflect.TypeOf(0.01).Size()) // 8
    fmt.Println("slice size: ",reflect.TypeOf([]int{}).Size()) // 24
    fmt.Println("struct alignS size: ", reflect.TypeOf(alignS{}).Size()) // 32
    fmt.Println("string size: ", reflect.TypeOf("abc").Size()) // 16
    fmt.Println("map[int]string size: ", reflect.TypeOf(map[int]string{}).Size()) // 8
    fmt.Println("chan size: ", reflect.TypeOf(make(chan int)).Size()) // 8
}

// 返回类型的字符串表示。该字符串可能会使用短包名（如用base64代替"encoding/base64"）
// 也不保证每个类型的字符串表示不同。如果要比较两个类型是否相等，请直接用Type类型比较。
func TestString(t *testing.T) {
    fmt.Println("int str: ", reflect.TypeOf(1).String()) // int
    fmt.Println("slice str: ", reflect.TypeOf([]int{}).String()) // []int
    fmt.Println("struct str: ", reflect.TypeOf(alignS{}).String()) // main.alignS
}

// Kind返回该接口的具体分类
func TestKind(t *testing.T) {
    fmt.Println("int kind str: ", reflect.TypeOf(1).Kind()) // int
    fmt.Println("slice kind str: ", reflect.TypeOf([]int{}).Kind())  // slice
    fmt.Println("map kind str: ", reflect.TypeOf(map[string]interface{}{}).Kind()) // map
    fmt.Println("struct kind str: ", reflect.TypeOf(alignS{}).Kind()) // struct
}

func (a *alignS) String() string {
    return "align"
}
// 如果该类型实现了u代表的接口，会返回真
func TestImplements(t *testing.T) {
    fmt.Println("aligns implements fmt.Stringer: ",
        reflect.TypeOf(&alignS{}).Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem())) // true
    fmt.Println("aligns implements fmt.Stringer: ",
        reflect.TypeOf(alignS{}).Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem())) // false
}

// 如果该类型的值可以直接赋值给u代表的类型，返回真
func TestAssignable(t *testing.T) {
    fmt.Println("int can assign to float? ", reflect.TypeOf(1).AssignableTo(reflect.TypeOf(0.01))) // false
    fmt.Println("int can assign to bool?", reflect.TypeOf(1).AssignableTo(reflect.TypeOf(false))) // false
    fmt.Println("int can assign to string?", reflect.TypeOf(1).AssignableTo(reflect.TypeOf("abc"))) // false
}

// 如该类型的值可以转换为u代表的类型，返回真
func TestConvertible(t *testing.T) {
    fmt.Println("int can convert to float?", reflect.TypeOf(1).ConvertibleTo(reflect.TypeOf(0.01))) // true
    fmt.Println("int can convert to bool?", reflect.TypeOf(1).ConvertibleTo(reflect.TypeOf(false))) // false
    fmt.Println("int can convert to string?", reflect.TypeOf(1).ConvertibleTo(reflect.TypeOf("abc"))) // true
}

// Comparable reports whether values of this type are comparable.
func TestComparable(t *testing.T) {
    fmt.Println("int comparable?", reflect.TypeOf(1).Comparable()) // true
    fmt.Println("map comparable?", reflect.TypeOf(map[int]string{}).Comparable()) // false
    fmt.Println("slice comparable?", reflect.TypeOf([]int{}).Comparable()) // false
    fmt.Println("chan comparable?", reflect.TypeOf((chan int)(nil)).Comparable()) // true
}

func TestSomeTypeOfFunc(t *testing.T) {
    fmt.Println(reflect.PtrTo(reflect.TypeOf(1)).String())
    fmt.Println(reflect.SliceOf(reflect.TypeOf(1)).String())
    fmt.Println(reflect.MapOf(reflect.TypeOf(1), reflect.TypeOf("abc")).String())
    fmt.Println(reflect.ChanOf(reflect.SendDir, reflect.TypeOf("string")).String())
}


