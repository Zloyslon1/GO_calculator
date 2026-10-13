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

	printFinance(res.Finance())
}

func printFinance(f ancestral.Finance) {
	fmt.Println()
	fmt.Println("Финансовый канал (цифры 4, 5, 6 — столбец «Деньги»):")
	fmt.Printf("  Цифр в канале: %d  (4×%d, 5×%d, 6×%d)\n", f.Total, f.Count4, f.Count5, f.Count6)
	fmt.Println("  " + f.Model)

	stab := "неустойчивый"
	if f.Stable {
		stab = "устойчивый"
	}
	levels := map[int]string{
		0: "негативной денежной программы нет",
		1: "наименьшая",
		2: "средняя",
		3: "высокая",
	}
	fmt.Printf("  Модель: %s. Влияние родовых программ: %s.\n", stab, levels[f.Level])

	if len(f.Entries) > 0 {
		fmt.Printf("  Точки входа (%d):\n", len(f.Entries))
		for _, e := range f.Entries {
			digits := strings.Repeat(fmt.Sprint(e.Digit), e.Count)
			fmt.Printf("    %-5s (%s) — %s\n", digits, e.Part, e.Hint)
		}
	}
	if len(f.Missing) > 0 {
		fmt.Println("  Чего не хватает (что нарабатывать):")
		for _, m := range f.Missing {
			fmt.Printf("    нет %d (%s) — %s\n", m.Digit, m.Part, m.Hint)
		}
	}
}
