package utils

// HandlePanicError is a utility function that will panic if the given error is not nil
// This is useful when you want to crash the program if an error occurs, but do not
// want to write an `if err != nil { panic(err) }` block every time.
func HandlePanicError(err error) {
	if err != nil {
		panic(err)
	}
}
