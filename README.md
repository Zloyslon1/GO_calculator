# Привет! Это мой простейший калькулятор в рамках первого пет-проекта

## Обо мне
Я начинающий разработчик на **Go** (звездочки делают жирным).
Учусь писать чистый код и использовать Git.

### Мои навыки:
- Go (Golang)
- VS Code
- Терминал

### Пример моего кода:

```go
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
```

