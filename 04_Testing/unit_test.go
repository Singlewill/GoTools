package main

import (
	"testing"
)

func TestAdd(t *testing.T) {
	if ans := Add(1, 2); ans != 3 {
		t.Errorf("1 + 2 expected be 3, but %d got", ans)
	}

	if ans := Add(-10, -20); ans != -30 {
		t.Errorf("-10 + -20 expected be -30, but %d got", ans)
	}
}

//子测试
//go test -run TestMul/pos
//go test -run TestMul/neg
func TestMul(t *testing.T) {
	t.Run("pos", func(t *testing.T) {
		if Mul(2, 3) != 6 {
			t.Fatal("fail") //相比t.Errorf, t.Fatal遇错即即终止进程
		}

	})
	t.Run("neg", func(t *testing.T) {
		if Mul(2, -3) != -6 {
			t.Fatal("fail")
		}
	})
}

//////////////////////////////////////////////////
//测试用例的方式，结合了子测试
//推荐用法

var cases = []struct {
	Name string
	Num1 int
	Num2 int
	Ret  int
}{
	{"test1", 1, 1, 1},
	{"test2", 1, 2, 2},
	{"test3", 3, 4, 12},
	{"test4", 10, 4, 40},
}

func TestMul2(t *testing.T) {
	for _, c := range cases {
		//t.Helper()
		t.Run(c.Name, func(t *testing.T) {
			if ans := Mul(c.Num1, c.Num2); ans != c.Ret {
				t.Fatalf("%d * %d expected %d, but %d got", c.Num1, c.Num2, c.Ret, ans)
			}
		})

	}
}
