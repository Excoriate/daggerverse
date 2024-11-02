package main

import (
	"fmt"
)

func main() {
	fmt.Println("Hello, Dagger!")
}

func fibonacci(n int) string {
	if n < 0 {
		return "Please provide a positive integer.\n"
	}
	if n == 0 {
		return "Please provide a positive integer.\n"
	}
	if n == 1 {
		return "0 \n"
	}
	a, b := 0, 1
	for i := 0; i < n; i++ {
		fmt.Printf("%d ", a)
		next := a + b
		a = b
		b = next
	}
	fmt.Println()
	return fmt.Sprintf("%d", a)
}
