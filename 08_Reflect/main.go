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

//传入[]People, 在函数体内进行追加
//ps := make([]People, 0)
//slice2_test(&ps)
func slice2_test(p interface{}) error {

	sT := reflect.TypeOf(p)
	if sT.Kind() != reflect.Ptr {
		fmt.Println("Must be reflect.Ptr")
	}

	//取数组的值
	sVE := reflect.ValueOf(p).Elem()
	//取数组中元素的类型	 <<-------------------
	//p是[]People. 一次Elem()取元素类型p[0],二次Elem()取元素指针指向的元素
	sEE := sT.Elem().Elem()
	fmt.Println(sEE)
	fmt.Printf("%v\n", sVE)

	//根据reflect.Type创建reflect.Value对象
	sON := reflect.New(sEE)
	//对象的值
	sONE := sON.Elem()
	sONEId := sONE.FieldByName("ID")
	sONEName := sONE.FieldByName("Name")
	//sONEId.SetInt(10)
	sONEId.SetString("10")
	sONEName.SetString("李四")

	// 创建一个新数组并把元素的值追加进去
	newArr := make([]reflect.Value, 0)
	newArr = append(newArr, sON.Elem())
	resArr := reflect.Append(sVE, newArr...)
	//实际上也可以直接把新元素追加到输入切片中，不用额外定义新数组
	//resArr := reflect.Append(sVE, sONE)

	// 最终结果给原数组
	sVE.Set(resArr)

	return nil
}
func basic_test(p interface{}) {
	t := reflect.TypeOf(p)
	fmt.Println(t)
	fmt.Println(t.Elem())
	fmt.Println(t.Elem().Elem())

	v := reflect.ValueOf(p)
	fmt.Println(v)
	fmt.Println(v.Elem())
}

func main() {
	ps := make([]People, 0)
	p1 := People{"kako", "18", "1"}
	ps = append(ps, p1)
	//basic_test(&ps)
	slice2_test(&ps)
	fmt.Println(ps)
	/*
			ps := People{"kako", "18", "1"}
			parse2(ps)
				ps := People{}
				new := slice_test(ps)
				fmt.Println(new)
				fmt.Printf("%T\n", new)
		ps := make([]People, 0)
		fmt.Printf("before %p\n", ps)

		slice2_test(&ps)
		fmt.Printf("after %p\n", ps)
	*/

}
