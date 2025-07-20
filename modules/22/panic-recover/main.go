package main

func myPanic1() {
	panic("This is a panic1")
}

func myPanic2() {
	panic("This is a panic2")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			if r == "This is a panic1" {
				println("Recovered from panic1:", r)
			}
			if r == "This is a panic2" {
				println("Recovered from panic2:", r)
			}
		}
	}()

	println("Before panic")
	myPanic2()
	println("After panic")
}
