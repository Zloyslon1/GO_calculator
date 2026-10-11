// Команда ancestral считает ансестологическую карту и число родового пути.
//
//	go run ./cmd/ancestral 01.12.1983 "Иванова Лариса Петровна"
package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"main.go/ancestral"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "Использование: ancestral ДД.ММ.ГГГГ [ФИО]")
		os.Exit(2)
	}
	date, err := time.Parse("02.01.2006", os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Не получилось прочитать дату, нужен вид ДД.ММ.ГГГГ:", err)
		os.Exit(2)
	}
	name := strings.Join(os.Args[2:], " ")

	res, err := ancestral.Calc(date)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка:", err)
		os.Exit(1)
	}

	if name != "" {
		fmt.Println(name)
	}
	fmt.Println(date.Format("02.01.2006"))
	parts := make([]string, len(res.Working))
	for i, w := range res.Working {
		parts[i] = fmt.Sprint(w)
	}
	fmt.Println(strings.Join(parts, "."))
	fmt.Println()
	fmt.Println("Число родового пути (ЧРП):", res.ChRP)
	fmt.Println()
	fmt.Println("Ансестологическая карта:")
	for _, row := range res.Grid() {
		fmt.Printf("  %-7s %-7s %s\n", row[0], row[1], row[2])
	}
	if res.At21 != nil {
		fmt.Println()
		fmt.Println("После 21 года к цифрам добавится по одной (кроме 5), описывать карту нужно без них.")
	}
}
