package utils

// HandleError is a utility function that will panic if the given error is not nil
// This is useful when you want to crash the program if an error occurs, but do not
// want to write an `if err != nil { panic(err) }` block every time.
func HandleError(err error, isPanic bool) {
	if err != nil {
		if isPanic {
			panic(err)
		} else {
			println(err.Error())
		}
	}

}
