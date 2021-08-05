// +build !trace

package functrace

//use for what?  func Trace() func() {
func shadeTrace() func() {
	return func() {

	}
}
