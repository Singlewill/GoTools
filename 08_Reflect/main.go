package main

import (
	"fmt"
	"reflect"
)

type People struct {
	Name string
	Age  string
	ID   string
}

type People2 struct {
	Name string
	Age  string
}

//reflect.Value 展开的数据类型
func value_test() {
	v := reflect.ValueOf(3) // a reflect.Value
	x := v.Interface()      // an interface{}
	i := x.(int)            // an int
	fmt.Printf("%d\n", i)   // "3"
}

//普普通通在函数内部改外部的struct
//由于要在函数内部改修传入变量，因此必须传入指针类型
func parse(p interface{}) {
	v := reflect.ValueOf(p)
	t := reflect.TypeOf(p)
	if t.Kind() == reflect.Ptr {
		// 传入的inStructPtr是指针，需要.Elem()取得指针指向的value
		t = t.Elem()
		v = v.Elem()
	} else {
		panic("inStructPtr must be ptr to struct")
	}

	for i := 0; i < t.NumField(); i++ {
		//v_v.Elem().Field(i).Set(reflect.ValueOf("haha"))

		f := v.Field(i)

		f.Set(reflect.ValueOf("1234"))
	}

}

//断言测试
func parse2(p interface{}) {
	switch t := p.(type) {
	default:
		fmt.Printf("unexpected type %T\n", t)
		break
	case People:
		fmt.Println("this is people")
	case People2:
		fmt.Println("this is people2")

	}
}

//根据参数类型，定义新的变量进行赋值
func parse3(p interface{}) {
	//定义新变量
	sy, ok := p.(People)
	if !ok {
		fmt.Println("Type MisMactch")
		return
	}

	v_v := reflect.ValueOf(&sy).Elem()
	t_t := reflect.TypeOf(&sy).Elem()

	fmt.Printf("%d\n", t_t.NumField())

	for i := 0; i < t_t.NumField(); i++ {
		f := v_v.Field(i)

		f.Set(reflect.ValueOf("1234"))
	}

	fmt.Printf("%v\n", sy)

}

//传入People,传出[]People
func slice_test(p interface{}) interface{} {
	//定义切片
	t := reflect.TypeOf(p)
	slice_t := reflect.SliceOf(t) //[]People
	s := reflect.MakeSlice(slice_t, 0, 0)

	//定义切片元素
	sy, ok := p.(People)
	if !ok {
		fmt.Println("Type MisMactch")
		return s
	}
	//元素赋值
	v_v := reflect.ValueOf(&sy).Elem()
	t_t := reflect.TypeOf(&sy).Elem()
	for i := 0; i < t_t.NumField(); i++ {
		f := v_v.Field(i)
		f.Set(reflect.ValueOf("1234"))
	}

	//将元素追加到切片中
	s = reflect.Append(s, reflect.ValueOf(sy))

	return s.Interface()
}

func main() {
	ps := People{"kako", "18", "1"}
	parse2(ps)

	/*
		ps := People{}
		new := slice_test(ps)
		fmt.Println(new)
		fmt.Printf("%T\n", new)
	*/

}
