package main

import (
	"fmt"
)

func main() {
	var a, b float32
	var operator string
	fmt.Scan(&a, &b, &operator)
	switch operator {
	case "+":
		fmt.Println(a + b)
	case "-":
		fmt.Println(a - b)
	case "*":
		fmt.Println(a * b)
	case "/":
		if b != 0 {
			fmt.Println(a / b)
		} else {
			fmt.Println("Делить на ноль нельзя!")
		}
	case "%":
		if b != 0 {
			fmt.Println(float32(int(a) % int(b)))
		} else {
			fmt.Println("Делить на ноль нельзя!")
		}
	default:
		fmt.Println("Неизвестная операция")

	}
}
