package main

import (
	"fmt"
	"testing"
)

type MyStruct struct {
	a int
}

var (
	t *testing.T
)

func TestPassingByReference(tst *testing.T) {
	t = tst
	m := MyStruct{a: 1}
	incByValue(m)
	assertEquals(m.a, 1)

	incByRef(&m)
	assertEquals(m.a, 2)
}

func incByRef(m *MyStruct) {
	m.a++
	assertEquals(m.a, 2)
}

func incByValue(m MyStruct) {
	m.a++
	assertEquals(m.a, 2)
}

func TestSayHello(t *testing.T) {
	t.Run("test 1", func(t *testing.T) {
		if SayHello("World") != "Hello World!" {
			t.Fail()
		}
	})
}

func ExampleSayHello() {
	msg := SayHello("World")
	fmt.Println(msg)
	// Output: Hello World!
}

func BenchmarkSayHello(b *testing.B) {
	for i := 0; i < b.N; i++ {
		SayHello("World")
	}
}

func assertEquals(a, b int) {
	if a != b {
		fmt.Printf("%v is not equal to %v\n", a, b)
		t.Fail()
	}
}
