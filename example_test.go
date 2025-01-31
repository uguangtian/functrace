// go test -tags trace
package functrace_test

import (
	"github.com/uguangtian/functrace"
)

func a() {
	defer functrace.Trace()()
	b()
}

func b() {
	defer functrace.Trace()()
	c()
}

func c() {
	defer functrace.Trace()()
	d()
}

func d() {
	defer functrace.Trace()()
}

func ExampleTrace() {
	a()
	// Output:
	// g[01]:	->github.com/bigwhite/functrace_test.a
	// g[01]:		->github.com/bigwhite/functrace_test.b
	// g[01]:			->github.com/bigwhite/functrace_test.c
	// g[01]:				->github.com/bigwhite/functrace_test.d
	// g[01]:				<-github.com/bigwhite/functrace_test.d
	// g[01]:			<-github.com/bigwhite/functrace_test.c
	// g[01]:		<-github.com/bigwhite/functrace_test.b
	// g[01]:	<-github.com/bigwhite/functrace_test.a
}
