package main

import (
    "fmt"
    "reflect"
    "testing"
    "unsafe"
)

// ValueOf返回一个初始化为i接口保管的具体值的Value，ValueOf(nil)返回Value零值
func TestValueOf(t *testing.T) {
   fmt.Println("valueOf 1: ", reflect.ValueOf(1))
}

// Zero返回一个持有类型typ的零值的Value。注意持有零值的Value和Value零值是两回事。
// Value零值表示不持有任何值。例如Zero(TypeOf(42))返回一个Kind为Int、值为0的Value。Zero的返回值不能设置也不会寻址。
func TestZeroValue(t *testing.T) {
    fmt.Println("int zero value: ", reflect.Zero(reflect.TypeOf(1)))
    fmt.Println("string zero value: ", reflect.Zero(reflect.TypeOf("")))
    fmt.Println("slice zero value: ", reflect.Zero(reflect.TypeOf([]int{})))
    fmt.Println("map zero value: ", reflect.Zero(reflect.TypeOf(map[string]int{})))
}

// New返回一个Value类型值，该值持有一个指向类型为typ的新申请的零值的指针，返回值的Type为PtrTo(typ)。
func TestNew(t *testing.T) {
    fmt.Println("New int Type: ", reflect.New(reflect.TypeOf(1)).Type().String())
}

// NewAt返回一个Value类型值，该值持有一个指向类型为typ、地址为p的值的指针。
func TestNewAt(t *testing.T) {
    i := 5
    fmt.Println("NewAt: ", reflect.NewAt(reflect.TypeOf(1), unsafe.Pointer(&i)).Elem())
}

// Indirect returns the value that v points to.
func TestIndirect(t *testing.T) {
    i := 5
    fmt.Println(reflect.Indirect(reflect.ValueOf(&i)).Type().String())
}

// 创建slice
func TestMakeSlice(t *testing.T) {
    v := reflect.MakeSlice(reflect.TypeOf([]int{}), 10, 10)
    fmt.Println(v)
}

// 创建map
func TestMakeMap(t *testing.T) {
    v := reflect.MakeMap(reflect.TypeOf(map[string]int{}))
    fmt.Println(v)
}

// 创建channel
func TestMakeChan(t *testing.T) {
    fmt.Println(reflect.MakeChan(reflect.TypeOf(make(chan int)), 10))
}

// MakeFunc(typ Type, fn func(args []value) (results []value)) 返回一个具有给定类型、包装函数fn的函数的Value封装。
// 当被调用时，该函数会：

// 1. 将提供给它的参数转化为Value切片
// 2. 执行results := fn(args)
// 3. 将results中每一个result依次排列作为返回值

// 函数fn的实现可以假设参数Value切片匹配typ类型指定的参数数目和类型。如果typ表示一个可变参数函数类型，
// 参数切片中最后一个Value本身必须是一个包含所有可变参数的切片。fn返回的结果Value切片也必须匹配typ类型指定的结果数目和类型。

//Value.Call方法允许程序员使用Value调用一个有类型约束的函数；反过来，MakeFunc方法允许程序员使用Value实现一个有类型约束的函数。
func TestMakeFunc(t *testing.T) {
    // swap是传递给MakeFunc的实现。它必须在reflect.Values方面起作用，以便可以在不事先知道什么类型的情况下编写代码, 类似c++的泛型？
    swap := func(in []reflect.Value) []reflect.Value {
        return []reflect.Value{in[1], in[0]}
    }
    // makeSwap期望 fptr 是指向nil函数的指针。它将指针设置为使用MakeFunc创建的新函数。调用函数时，reflect将参数转换为Values，
    // 调用swap，然后将swap的结果切片转换为新函数返回的值。
    makeSwap := func(fptr interface{}) {
        // fptr 是指向函数的指针。获取函数值本身（可能为nil）作为reflect.Value，以便我们可以查询其类型，然后设置该值。
        fn := reflect.ValueOf(fptr).Elem()
        // 使用正确的类型创建一个函数
        v := reflect.MakeFunc(fn.Type(), swap)
        // 赋值到fn表示的value中
        fn.Set(v)
    }

    // swap int
    var intSwap  func(int, int) (int, int)
    makeSwap(&intSwap)
    fmt.Println(intSwap(0, 1))

    // swap float
    var floatSwap func(float64, float64) (float64, float64)
    makeSwap(&floatSwap)
    fmt.Println(floatSwap(2.72, 3.14))
}

// Append 向切片类型的Value值s中添加一系列值，x等Value值持有的值必须能直接赋值给s持有的切片的元素类型
// AppendSlice 类似Append函数，但接受一个切片类型的Value值。将切片t的每一个值添加到s。
func TestAppend(t *testing.T) {
    fmt.Println(reflect.Append(reflect.ValueOf([]int{}), reflect.ValueOf(1), reflect.ValueOf(2)))
    fmt.Println(reflect.AppendSlice(reflect.ValueOf([]int{}), reflect.ValueOf([]int{1,2,3,4,5})))
}

// Copy 将src中的值拷贝到dst，直到src被耗尽或者dst被装满，要求这二者都是slice或array，且元素类型相同。
func TestCopy(t *testing.T) {
    reflect.Copy(reflect.ValueOf([]int{1,2,3}), reflect.ValueOf([]int{4,5,6}))
}

// 用来判断两个值是否深度一致：除了类型相同；在可以时（主要是基本类型）会使用==；但还会比较array、slice的成员，map的键值对，
// 结构体字段进行深入比对。map的键值对，对键只使用==，但值会继续往深层比对。DeepEqual函数可以正确处理循环的类型。
// 函数类型只有都会nil时才相等；空切片不等于nil切片；还会考虑array、slice的长度、map键值对数。
func TestDeepEqual(t *testing.T) {
    f := func(){}
    f1 := f
    fmt.Println(reflect.DeepEqual(f, f1)) // false
}

// IsValid返回v是否持有一个值。如果v是Value零值会返回假，此时v除了IsValid、String、Kind之外的方法都会导致panic。
// 绝大多数函数和方法都永远不返回Value零值。如果某个函数/方法返回了非法的Value，它的文档必须显式的说明具体情况。
func TestValueIsValid(t *testing.T) {
    fmt.Println(reflect.ValueOf(1).IsValid()) // true
    fmt.Println(reflect.ValueOf(nil).IsValid()) // false
}

// IsNil报告v持有的值是否为nil。v持有的值的分类必须是通道、函数、接口、映射、指针、切片之一；否则IsNil函数会导致panic。
// 注意IsNil并不总是等价于go语言中值与nil的常规比较。例如：如果v是通过使用某个值为nil的接口调用ValueOf函数创建的，v.IsNil()返回真，
// 但是如果v是Value零值，会panic。
func TestValueIsNil(t *testing.T) {
    defer func(){
        if err := recover(); err != nil {
            fmt.Println("panic error: ", err)
        }
    }()
    fmt.Println(reflect.ValueOf((*error)(nil)).IsNil()) // true
    fmt.Println(reflect.ValueOf(nil).IsNil()) // panic
}

// Kind返回v持有的值的分类，如果v是Value零值，返回值为Invalid
func TestValueKind(t *testing.T){
    fmt.Println(reflect.ValueOf(1).Kind())
}

// 返回v持有的值的类型的Type表示
func TestValueType(t *testing.T) {
    fmt.Println(reflect.ValueOf(1).Type().String())
}

// Convert将v持有的值转换为类型为t的值，并返回该值的Value封装。如果go转换规则不支持这种转换，会panic
func TestValueConvert(t *testing.T) {
    fmt.Println(reflect.ValueOf(1).Convert(reflect.TypeOf(1.1))) // panic
}

// Elem返回v持有的接口保管的值的Value封装，或者v持有的指针指向的值的Value封装。
// 如果v的Kind不是Interface或Ptr会panic；如果v持有的值为nil，会返回Value零值
func TestValueElem(t *testing.T) {
    i := 1
    var ii interface{} = i
    var iii interface{} = &ii

    fmt.Println(reflect.ValueOf(iii).Elem().Elem())
    fmt.Println(reflect.ValueOf(&i).Elem())
}

// 返回v持有的布尔值，如果v的Kind不是Bool会panic
func TestValueBool(t *testing.T) {
    fmt.Println(reflect.ValueOf(false).Bool())
}

// 返回v持有的有符号整数（表示为int64），如果v的Kind不是Int、Int8、Int16、Int32、Int64会panic
func TestValueInt(t *testing.T) {
    fmt.Println(reflect.ValueOf(1).Int())
}

// 如果v持有值的类型不能无溢出的表示x，会返回真。如果v的Kind不是Int、Int8、Int16、Int32、Int64会panic
func TestValueOverflowInt(t *testing.T) {
    fmt.Println(reflect.ValueOf(int8(1)).OverflowInt(111111111111111111))
}

// 返回v持有的无符号整数（表示为uint64），如v的Kind不是Uint、Uintptr、Uint8、Uint16、Uint32、Uint64会panic
func TestValueUint(t *testing.T) {
    fmt.Println(reflect.ValueOf(uint(1)).Uint())
}

// 如果v持有值的类型不能无溢出的表示x，会返回真。如果v的Kind不是Uint、Uintptr、Uint8、Uint16、Uint32、Uint64会panic
func TestValueOverflowUint(t *testing.T) {
    fmt.Println(reflect.ValueOf(uint8(1)).OverflowUint(11111111))
}

// Float 返回v持有的浮点数（表示为float64），如果v的Kind不是Float32、Float64会panic
// OverflowFloat 如果v持有值的类型不能无溢出的表示x，会返回真。如果v的Kind不是Float32、Float64会panic
func TestValueFloat(t *testing.T) {
    fmt.Println(reflect.ValueOf(1.1).Float())
    fmt.Println(reflect.ValueOf(float32(1.1)).OverflowFloat(2.2))
}

// Complex 返回v持有的复数（表示为complex64），如果v的Kind不是Complex64、Complex128会panic
// OverflowComplex 如果v持有值的类型不能无溢出的表示x，会返回真。如果v的Kind不是Complex64、Complex128会panic
func TestValueComplex(t *testing.T) {
    fmt.Println(reflect.ValueOf(complex(1.1, 2.2)).Complex())
    fmt.Println(reflect.ValueOf(complex(1.1,2.2)).OverflowComplex(complex(1.2, 2.3)))
}

// 将v持有的值作为一个指针返回。本方法返回值不是unsafe.Pointer类型，以避免程序员不显式导入unsafe包却得到unsafe.Pointer类型表示的指针。
// 如果v的Kind不是Chan、Func、Map、Ptr、Slice或UnsafePointer会panic。
//
// 如果v的Kind是Func，返回值是底层代码的指针，但并不足以用于区分不同的函数；只能保证当且仅当v持有函数类型零值nil时，返回值为0。
// 如果v的Kind是Slice，返回值是指向切片第一个元素的指针。如果持有的切片为nil，返回值为0；如果持有的切片没有元素但不是nil，返回值不会是0
func TestValuePointer(t *testing.T) {
    fmt.Println(reflect.ValueOf([]int{2, 1}).Pointer())
}

// 返回v持有的[]byte类型值。如果v持有的值的类型不是[]byte会panic
func TestValueBytes(t *testing.T) {
    fmt.Println(reflect.ValueOf([]byte("abcd")).Bytes())
}

// 返回v持有的值的字符串表示。因为go的String方法的惯例，Value的String方法比较特别。和其他获取v持有值的方法不同：v的Kind是String时，
// 返回该字符串；v的Kind不是String时也不会panic而是返回格式为"<T value>"的字符串，其中T是v持有值的类型
func TestValueString(t *testing.T) {
    fmt.Println(reflect.ValueOf("abcd").String())
    fmt.Println(reflect.ValueOf([]int{1,2,3,4,5}).String())
}

//返回v持有的接口类型值的数据。如果v的Kind不是Interface会panic
func TestValueInterfaceData(t *testing.T) {
    i := 5
    var ii interface{} = i
    var iii interface{} = &ii
    fmt.Println(reflect.ValueOf(iii).Elem().InterfaceData())
}
//Slice 返回v[i:j]（v持有的切片的子切片的Value封装）；如果v的Kind不是Array、Slice或String会panic。如果v是一个不可寻址的数组，或者索引出界，也会panic
//Slice3 是Slice的3参数版本，返回v[i:j:k] ；如果v的Kind不是Array、Slice或String会panic。如果v是一个不可寻址的数组，或者索引出界，也会panic。
//Len 返回v持有值的长度，如果v的Kind不是Array、Chan、Slice、Map、String会panic
//Cap 返回v持有值的容量，如果v的Kind不是Array、Chan、Slice会panic
//Index 返回v持有值的第i个元素。如果v的Kind不是Array、Chan、Slice、String，或者i出界，会panic
func TestValueSlice(t *testing.T) {
    fmt.Println(reflect.ValueOf([]int{1,2,3,4,5,6}).Slice(2,5))
    fmt.Println(reflect.ValueOf([]int{1,2,3,4,5,6}).Slice3(0, 2, 3).Cap())
    fmt.Println(reflect.ValueOf([]int{1,2,3,4,5,6}).Slice(2,5).Len())
    fmt.Println(reflect.ValueOf([]int{1,2,3,4,5,6}).Index(5))
}

// MapIndex 返回v持有值里key持有值为键对应的值的Value封装。如果v的Kind不是Map会panic。
// 如果未找到对应值或者v持有值是nil映射，会返回Value零值。key的持有值必须可以直接赋值给v持有值类型的键类型
func TestValueMapIndex(t *testing.T) {
    fmt.Println(reflect.ValueOf(map[string]string{"key":"value"}).MapIndex(reflect.ValueOf("key")))
    fmt.Println(reflect.ValueOf(map[string]string{"key":"value"}).MapIndex(reflect.ValueOf("key2")))
    fmt.Println(reflect.ValueOf(map[string]string(nil)).MapIndex(reflect.ValueOf("key")))
}

//MapKeys 返回一个包含v持有值中所有键的Value封装的切片，该切片未排序。如果v的Kind不是Map会panic。如果v持有值是nil，返回空切片（非nil）。
func TestValueMapKeys(t *testing.T) {
    fmt.Println(reflect.ValueOf(map[string]string{"k1":"v1", "k2":"v2"}).MapKeys())
    fmt.Println(reflect.ValueOf(map[string]string(nil)).MapKeys())
}

// NumField 返回v持有的结构体类型值的字段数，如果v的Kind不是Struct会panic
// Field 返回结构体的第i个字段（的Value封装）。如果v的Kind不是Struct或i出界会panic
// FieldByIndex 返回索引序列指定的嵌套字段的Value表示，等价于用索引中的值链式调用本方法，如v的Kind非Struct将会panic
// FieldByName 返回该类型名为name的字段（的Value封装）（会查找匿名字段及其子字段），如果v的Kind不是Struct会panic；如果未找到会返回Value零值
// FieldByNameFunc 返回该类型第一个字段名满足match的字段（的Value封装）（会查找匿名字段及其子字段），如果v的Kind不是Struct会panic；如果未找到会返回Value零值
func TestValueStruct(t *testing.T) {
    type s struct {
        aa int
        b string
    }
    fmt.Println(reflect.ValueOf(s{1, "abc"}).NumField())
    fmt.Println(reflect.ValueOf(s{1, "abc"}).Field(1))
    fmt.Println(reflect.ValueOf(s{1, "abc"}).FieldByIndex([]int{1}))
    fmt.Println(reflect.ValueOf(s{1, "abc"}).FieldByName("b"))
    fmt.Println(reflect.ValueOf(s{1, "abc"}).FieldByNameFunc(func (name string) bool {
        return len(name) == 1
    }))
}

// Recv 方法从v持有的通道接收并返回一个值（的Value封装）。如果v的Kind不是Chan会panic。方法会阻塞直到获取到值。
// 如果返回值x对应于某个发送到v持有的通道的值，ok为真；如果因为通道关闭而返回，x为Value零值而ok为假。

// TryRecv尝试从v持有的通道接收一个值，但不会阻塞。如果v的Kind不是Chan会panic。如果方法成功接收到一个值，会返回该值（的Value封装）和true；
// 如果不能无阻塞的接收到值，返回Value零值和false；如果因为通道关闭而返回，返回值x是持有通道元素类型的零值的Value和false。

// Send 方法向v持有的通道发送x持有的值。如果v的Kind不是Chan，或者x的持有值不能直接赋值给v持有通道的元素类型，会panic。
// TrySend 尝试向v持有的通道发送x持有的值，但不会阻塞。如果v的Kind不是Chan会panic。如果成功发送会返回真，否则返回假。x
// 的持有值必须可以直接赋值给v持有通道的元素类型。

// Close 关闭v持有的通道，如果v的Kind不是Chan会panic
func TestValueChan(t *testing.T) {
    c := make(chan int, 1)
    reflect.ValueOf(c).Send(reflect.ValueOf(100))
    if v, ok := reflect.ValueOf(c).Recv(); ok {
        fmt.Println("recv from channel: ", v)
    }
    reflect.ValueOf(c).TrySend(reflect.ValueOf(99))
    if v, ok := reflect.ValueOf(c).TryRecv(); ok {
        fmt.Println("recv from channel: ", v)
    }
    reflect.ValueOf(c).Close()
    if _, ok := reflect.ValueOf(c).TryRecv(); !ok {
        fmt.Println("channel has closed")
    }
}

// Call方法使用输入的参数in调用v持有的函数。如果v的Kind不是Func会panic。它返回函数所有输出结果的Value封装的切片。
// 和go代码一样，每一个输入实参的持有值都必须可以直接赋值给函数对应输入参数的类型。如果v持有值是可变参数函数，
// Call方法会自行创建一个代表可变参数的切片，将对应可变参数的值都拷贝到里面。
func TestValueFuncCall(t *testing.T) {
    swap := func(a, b int) (int, int){
        return b, a
    }
    results := reflect.ValueOf(swap).Call([]reflect.Value{reflect.ValueOf(2), reflect.ValueOf(3)})
    fmt.Printf("b: %d, a : %d\n", results[0].Int(), results[1].Int())
}

// CallSlice调用v持有的可变参数函数，会将切片类型的in[len(in)-1]（的成员）分配给v的最后的可变参数。
// 如果v的Kind不是Func或者v的持有值不是可变参数函数，会panic。它返回函数所有输出结果的Value封装的切片。
// 和go代码一样，每一个输入实参的持有值都必须可以直接赋值给函数对应输入参数的类型。
func Sum(nums ...int) int{
    s := 0
    for _, n := range nums {
        s += n
    }
    return s
}
func TestValueFuncCallSlice(t *testing.T) {

    fmt.Println(Sum(1,2,3,4,5))
    fmt.Println(reflect.TypeOf(Sum).NumIn())
    result := reflect.ValueOf(Sum).CallSlice([]reflect.Value{reflect.ValueOf([]int{1,2,3,4,5})})
    fmt.Println("get sum: ", result[0].Int())
}

// NumMethod 返回v持有值的方法集的方法数目。
// Method 返回v持有值类型的第i个方法的已绑定（到v的持有值的）状态的函数形式的Value封装。返回值调用Call方法时不应包含接收者；
// 返回值持有的函数总是使用v的持有者作为接收者（即第一个参数）。如果i出界，或者v的持有值是接口类型的零值（nil），会panic。
// MethodByName 返回v的名为name的方法的已绑定（到v的持有值的）状态的函数形式的Value封装。返回值调用Call方法时不应包含接收者；
// 返回值持有的函数总是使用v的持有者作为接收者（即第一个参数）。如果未找到该方法，会返回一个Value零值。
type sm struct {
    a int
    b string
}

type s struct {
    sm
    c float32
}
func (s *s) String() string{
    return ""
}
func (s *s) Sum() float32 {
    return float32(s.a) + s.c
}

func TestValueMethod(t *testing.T) {
    fmt.Println(reflect.ValueOf(&s{}).NumMethod())
    fmt.Println(reflect.ValueOf(&s{}).Method(1))
    fmt.Println(reflect.ValueOf(&s{}).MethodByName("Sum"))
}

// CanAddr 返回是否可以获取v持有值的指针。可以获取指针的值被称为可寻址的。如果一个值是切片或可寻址数组的元素、可寻址结构体的字段、
// 或从指针解引用得到的，该值即为可寻址的。
// Addr 函数返回一个持有指向v持有者的指针的Value封装。如果v.CanAddr()返回假，调用本方法会panic。
// Addr 一般用于获取结构体字段的指针或者切片的元素（的Value封装）以便调用需要指针类型接收者的方法。
// UnsafeAddr 返回指向v持有数据的地址的指针（表示为uintptr）以用作高级用途，如果v不可寻址会panic。
// CanInterface 如果CanInterface 返回真，v可以不导致panic的调用Interface方法。
// Interface 本方法返回v当前持有的值（表示为/保管在interface{}类型) .如果v是通过访问非导出结构体字段获取的，会导致panic。

func TestValueAddrInterface(t *testing.T) {
    fmt.Println(reflect.ValueOf([]int{1,2,3}[0]).CanAddr()) // false
    fmt.Println(reflect.ValueOf(&([]int{1,2,3}[0])).Elem().CanAddr()) // true
    s1 := []int{1,2,3}
    fmt.Println(reflect.ValueOf(&s1[0]).Elem().CanAddr()) // true
    fmt.Println(reflect.ValueOf([3]int{1,2,3}[2]).CanAddr()) // false
    fmt.Println(reflect.ValueOf(struct{name string}{name: "abc"}.name).CanAddr()) // false
    fmt.Println(reflect.ValueOf(&s{}).Elem().CanAddr()) // true

    fmt.Println(reflect.ValueOf(&s1[0]).Elem().Addr()) // 0xc00001a0f0
    fmt.Println(reflect.ValueOf(&s1[0]).Elem().UnsafeAddr()) // 824633827568

    fmt.Println(reflect.ValueOf(1).CanInterface()) // true
    fmt.Println(reflect.ValueOf(1).Interface()) // 1
}

// CanSet 如果v持有的值可以被修改，CanSet就会返回真。只有一个Value持有值可以被寻址同时又不是来自非导出字段时，它才可以被修改。
// 如果CanSet返回假，调用Set或任何限定类型的设置函数（如SetBool、SetInt64）都会panic。
// SetBool 设置v的持有值。如果v的Kind不是Bool或者v.CanSet()返回假，会panic。
// SetInt 设置v的持有值。如果v的Kind不是Int、Int8、Int16、Int32、Int64之一或者v.CanSet()返回假，会panic。
// SetUint 设置v的持有值。如果v的Kind不是Uint、Uintptr、Uint8、Uint16、Uint32、Uint64或者v.CanSet()返回假，会panic。
// SetFloat 设置v的持有值。如果v的Kind不是Float32、Float64或者v.CanSet()返回假，会panic。
// SetComplex 设置v的持有值。如果v的Kind不是Complex64、Complex128或者v.CanSet()返回假，会panic。
// SetBytes 设置v的持有值。如果v持有值不是[]byte类型或者v.CanSet()返回假，会panic。
// SetString设置v的持有值。如果v的Kind不是String或者v.CanSet()返回假，会panic。
// SetPointer 设置v的持有值。如果v的Kind不是UnsafePointer或者v.CanSet()返回假，会panic。
// SetCap 设定v持有值的容量。如果v的Kind不是Slice或者n出界（小于长度或超出容量），将导致panic
// SetLen 设定v持有值的长度。如果v的Kind不是Slice或者n出界（小于零或超出容量），将导致panic
// SetMapIndex 用来给v的映射类型持有值添加/修改键值对，如果val是Value零值，则是删除键值对。如果v的Kind不是Map，
// 或者v的持有值是nil，将会panic。key的持有值必须可以直接赋值给v持有值类型的键类型。val的持有值必须可以直接赋值给v持有值类型的值类型。
// Set 将v的持有值修改为x的持有值。如果v.CanSet()返回假，会panic。x的持有值必须能直接赋给v持有值的类型

func TestValueSet(t *testing.T) {
    fmt.Println(reflect.ValueOf([]int{1,2}).CanSet()) // false
    fmt.Println(reflect.ValueOf(1).CanSet()) // false
    s := []int{1,2}
    fmt.Println(reflect.ValueOf(s[0]).CanSet()) // false
    fmt.Println(reflect.ValueOf(&s[0]).Elem().CanSet()) // true

    b := false
    reflect.ValueOf(&b).Elem().SetBool(true)
    fmt.Println(b)

    i := 1
    reflect.ValueOf(&i).Elem().SetInt(2)
    fmt.Println(i)

    bs := []byte("abc")
    reflect.ValueOf(&bs).Elem().SetBytes([]byte("efg"))
    fmt.Println(string(bs))

    // reflect.ValueOf(&[]int{1,2}).Elem().SetCap(5) // panic cap out of range
    m := map[int]string{
        1: "a",
        2: "b",
    }
    reflect.ValueOf(m).SetMapIndex(reflect.ValueOf(3), reflect.ValueOf("abc")) // add key 3
    reflect.ValueOf(m).SetMapIndex(reflect.ValueOf(2), reflect.ValueOf(nil)) // del key 2
    fmt.Println(m)

    reflect.ValueOf(&i).Elem().Set(reflect.ValueOf(4))
    fmt.Println(i) // 4
}
