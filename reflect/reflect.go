package main

import (
    "fmt"
    "reflect"
)

// 反射的三大原则
// 1. 从interface{}变量可以反射出反射对象
// 2. 从反射对象可以获取interface{}变量
// 3. 要修改反射对象，其值必须可修改

func main() {
    principle1()
    principle2()
    principle32()
    implements()
    //principle3()
    runtimeCall()
}

// 1.1   interface{} ---> reflect.TypeOf()  ---> reflect.Type obj
// 1.2   interface{} ---> reflect.ValueOf() ---> reflect.Value obj
// TypeOf,ValueOf 两个方法的入参都是interface{}, 这里会发生类型转换。拿到一个变量的类型和值，就拿到了所有信息。
// 编译期：当我们想要将一个变量转换成反射对象时，Go 语言会在编译期间完成类型转换，
// 将变量的类型和值转换成了 interface{} 并等待运行期间使用 reflect 包获取接口中存储的信息。
func principle1() {
    author := "mark"
    fmt.Println("TypeOf: ", reflect.TypeOf(author))
    fmt.Println("ValueOf: ", reflect.ValueOf(author))
}

// 2.1  interface{} <--- reflect.Value.Interface() <--- reflect.Value obj
// 结合原则1,总结如下：
// value ---> typecast ---> interface{} ---> reflect.TypeOf()/ValueOf() ---> reflect.Type/Value obj
// value <--- typecast <--- interface{} <--- reflect.Value.Interface()  <--- reflect.Type/Value obj
func principle2() {
    v := reflect.ValueOf(1)
    if i, ok := v.Interface().(int); ok {
        fmt.Println("get i: ", i)
    }
}

// 3. 要修改反射对象，其值必须可修改
// 由于 Go 语言的函数调用都是传值的，所以我们得到的反射对象跟最开始的变量没有任何关系，
// 那么直接修改反射对象无法改变原始变量，程序为了防止错误就会崩溃。
// panic: reflect: reflect.Value.SetInt using unaddressable value
func principle3() {
    i := 1
    v := reflect.ValueOf(i)
    v.SetInt(10)
    fmt.Println(i)
}

// 由于 Go 语言的函数调用都是值传递的，所以我们只能只能用迂回的方式改变原变量：
// 先获取指针对应的 reflect.Value，再通过 reflect.Value.Elem 方法得到可以被设置的变量
func principle32() {
    i := 1
    v := reflect.ValueOf(&i)
    // Elem() 获取指针指向的变量
    v.Elem().SetInt(10)
    fmt.Println(i)
}

// 在 Go 语言中获取结构体的反射类型 reflect.Type 还是比较容易的，但是想要获得接口类型需要通过以下方式
// reflect.TypeOf((*<interface>)(nil)).Elem()
// 判断一个类型是否实现了某个接口：
type CustomError struct{}

func (*CustomError) Error() string {
    return ""
}

func implements() {
    typeOfError := reflect.TypeOf((*error)(nil)).Elem()
    customErrorPtr := reflect.TypeOf(&CustomError{})
    customError := reflect.TypeOf(CustomError{})

    fmt.Println(customErrorPtr.Implements(typeOfError)) // #=> true
    fmt.Println(customError.Implements(typeOfError)) // #=> false
}

// 运行时执行函数：
// 1. 通过 reflect.ValueOf 获取函数 Add 对应的反射对象；
// 2. 调用 reflect.rtype.NumIn 获取函数的入参个数；
// 3. 多次调用 reflect.ValueOf 函数逐一设置 argv 数组中的各个参数；
// 4. 调用反射对象 Add 的 reflect.Value.Call 方法并传入参数列表；
// 5. 获取返回值数组、验证数组的长度以及类型并打印其中的数据；

// 其中 reflect.Value.Call 执行如下(具体可以看源码)：
// 0. 确定当前反射对象的类型是函数以及可见性
// 1. 检查输入参数以及类型的合法性；
// 2. 将传入的 reflect.Value 参数数组设置到栈上；
// 3. 通过函数指针和输入参数调用函数；
// 4. 从栈上获取函数的返回值；
func Add(a, b int) int { return a + b }

func runtimeCall() {
    v := reflect.ValueOf(Add)
    if v.Kind() != reflect.Func {
        return
    }
    t := v.Type()
    argv := make([]reflect.Value, t.NumIn())
    for i := range argv {
        if t.In(i).Kind() != reflect.Int {
            return
        }
        argv[i] = reflect.ValueOf(i)
    }
    result := v.Call(argv)
    if len(result) != 1 || result[0].Kind() != reflect.Int {
        return
    }
    fmt.Println(result[0].Int()) // # 0 + 1 => 1
}
