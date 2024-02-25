package fserr

import (
	"fmt"
	"testing"
)

// 针对于 New、Wrap 两种构建错误的方式进行性能测试对比，结果如下（堆栈深度为10、100、1000），golang版本1.22:
// New
// fserr
// 801 ns/op
// 1650 ns/op
// 6262 ns/op
//
// Wrap
// fserr
// 826.2 ns/op
// 1623 ns/op
// 6201 ns/op
//
// WithCode
// 2077 ns/op
// 2765 ns/op
// 7628 ns/op

func yesErrors(at, depth int, how func() error) error {
	if at >= depth {
		return how()
	}
	return yesErrors(at+1, depth, how)
}

// GlobalE is an exported global to store the result of benchmark results,
// preventing the compiler from optimising the benchmark functions away.
var GlobalE interface{}

func BenchmarkNew(b *testing.B) {
	for _, r := range []int{
		10,
		100,
		1000,
	} {
		name := fmt.Sprintf("fserr-stack-%d", r)
		b.Run(name, func(b *testing.B) {
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = yesErrors(0, r, func() error {
					return New("success")
				})
			}
			b.StopTimer()
			GlobalE = err
		})
	}
}

func BenchmarkWrap(b *testing.B) {
	for _, r := range []int{
		10,
		100,
		1000,
	} {
		name := fmt.Sprintf("fserr-stack-%d", r)
		b.Run(name, func(b *testing.B) {
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = yesErrors(0, r, func() error {
					return Wrap(New("origin"), "success")
				})
			}
			b.StopTimer()
			GlobalE = err
		})
	}
}

func BenchmarkWithCode(b *testing.B) {
	NewOK(ErrDb, "success")
	for _, r := range []int{
		10,
		100,
		1000,
	} {
		name := fmt.Sprintf("fserr-stack-%d", r)
		b.Run(name, func(b *testing.B) {
			var err error
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				err = yesErrors(0, r, func() error {
					return WithCode(nil, ErrDb)
				})
			}
			b.StopTimer()
			GlobalE = err
		})
	}
}
