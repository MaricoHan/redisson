package helper

// GoSafe runs the given fn using another goroutine, recovers if fn panics.
func GoSafe(fn func(), callback func(error)) {
	go RunSafe(fn, callback)
}

// RunSafe runs the given fn, recovers if fn panics.
func RunSafe(fn func(), callback func(error)) {
	defer Recover(callback)
	fn()
}

// Recover is used with defer to do cleanup on panics.
// Use it like:
//  defer Recover(func() {})
func Recover(callback func(error)) {
	if p := recover(); p != nil {
		callback(p.(error))
	}
}
