// Package ancestral считает ансестологическую карту (квадрат Пифагора)
// и число родового пути (ЧРП) по дате рождения.
//
// Правила взяты из материалов курса «Метод родологических расчётов»,
// модуль «Ансестологическая карта».
package ancestral

import (
	"fmt"
	"strings"
	"time"
)

// Result — итог расчёта по одной дате рождения.
type Result struct {
	Date time.Time

	// Working — рабочие числа в том порядке, в котором их записывают
	// под датой рождения, например 25.7.23.5.
	Working []int

	// ChRP — число родового пути (рабочее число №2).
	ChRP int

	// Counts[d] — сколько раз цифра d (1..9) стоит в карте.
	// Counts[0] не используется.
	Counts [10]int

	// ExtraNine — дети с 2020 года: к карте добавляется вторая константа 9.
	// Она уже входит в Counts и показывается в скобках: «9 (9)».
	ExtraNine bool

	// At21 — карта после 21 года (только для родившихся с 2000 года).
	// Ко всем цифрам, кроме 5, добавляется по одной. Для описания
	// доверителя и его программ используется Counts, а не At21.
	At21 *[10]int
}

// Calc считает карту по дате рождения.
// Поддерживаются даты с 1800 по 2099 год.
func Calc(date time.Time) (Result, error) {
	y := date.Year()
	if y < 1800 || y > 2099 {
		return Result{}, fmt.Errorf("год %d вне диапазона 1800–2099", y)
	}

	res := Result{Date: date}
	digits := dateDigits(date)

	n1 := sum(digits)
	n2 := reduceOnce(n1)

	switch {
	case y >= 2000:
		// После 2000 года прибавляется константа 19.
		n3 := n1 + 19
		n4 := reduceFull(n3)
		res.Working = []int{n1, n2, 19, n3, n4}
		res.ExtraNine = y >= 2020
	case y >= 1900:
		// Обычный расчёт: из №1 вычитаем первую цифру даты, умноженную на 2.
		n3 := n1 - 2*digits[0]
		n4 := n3
		if n3 < 10 || n3 > 12 {
			n4 = digitSum(n3)
		}
		res.Working = []int{n1, n2, n3, n4}
	default:
		// 1800–1899: к рабочему числу №2 прибавляется константа 9.
		n3 := n2 + 9
		n4 := reduceFull(n3)
		res.Working = []int{n1, n2, 9, n3, n4}
	}
	res.ChRP = n2

	for _, d := range digits {
		res.Counts[d]++
	}
	for _, w := range res.Working {
		for _, d := range numDigits(w) {
			res.Counts[d]++
		}
	}
	if res.ExtraNine {
		res.Counts[9]++
	}

	if y >= 2000 {
		at21 := res.Counts
		for d := 1; d <= 9; d++ {
			if d != 5 {
				at21[d]++
			}
		}
		res.At21 = &at21
	}
	return res, nil
}

// dateDigits возвращает цифры даты ДДММГГГГ без нулей.
func dateDigits(date time.Time) []int {
	s := date.Format("02012006")
	var out []int
	for _, r := range s {
		if r != '0' {
			out = append(out, int(r-'0'))
		}
	}
	return out
}

// numDigits возвращает цифры числа без нулей.
func numDigits(n int) []int {
	var out []int
	for _, r := range fmt.Sprint(n) {
		if r != '0' {
			out = append(out, int(r-'0'))
		}
	}
	return out
}

func sum(xs []int) int {
	t := 0
	for _, x := range xs {
		t += x
	}
	return t
}

func digitSum(n int) int { return sum(numDigits(n)) }

// reduceOnce складывает цифры числа один раз. Результат 10, 11 и 12
// остаётся двузначным — это допустимые значения числа родового пути.
func reduceOnce(n int) int {
	if n < 10 {
		return n
	}
	return digitSum(n)
}

// reduceFull складывает цифры, пока не получится однозначное число.
func reduceFull(n int) int {
	for n >= 10 {
		n = digitSum(n)
	}
	return n
}

// Grid возвращает карту по строкам в том виде, как её рисуют на бумаге:
//
//	1 4 7
//	2 5 8
//	3 6 9
//
// Отсутствующая цифра обозначена «-». У детей с 2020 года в ячейке девяток
// вторая константа показана в скобках: «9 (9)».
func (r Result) Grid() [3][3]string {
	var g [3][3]string
	for row := 0; row < 3; row++ {
		for col := 0; col < 3; col++ {
			digit := col*3 + row + 1
			g[row][col] = cell(digit, r.Counts[digit], r.ExtraNine && digit == 9)
		}
	}
	return g
}

func cell(digit, n int, extra bool) string {
	if n == 0 {
		return "-"
	}
	ch := string(rune('0' + digit))
	if extra {
		return strings.Repeat(ch, n-1) + " (" + ch + ")"
	}
	return strings.Repeat(ch, n)
}

// Финансовый (родовой денежный) канал — это средний столбец карты, то есть
// цифры 4, 5 и 6 (столбец «Деньги»). Правила взяты из модуля 3, урок
// «Родовой финансовый канал», материалы «Родовой финансовый канал» и
// «Точки входа в финансовый канал».
//
// Канал представляют в виде дома: 4 — крыша, 5 — стены, 6 — фундамент.
// Сила и устойчивость канала зависят от того, сколько всего цифр 4, 5, 6
// в карте, а точки входа — от того, какие именно ячейки заполнены.
// В расчёте используется только дата рождения; ФИО доверителя ни на что
// не влияет и служит лишь подписью.

// FinanceCell — одна ячейка финансового канала (4, 5 или 6).
type FinanceCell struct {
	Digit int    // 4, 5 или 6
	Count int    // сколько раз цифра встречается в карте (0 — ячейка пустая)
	Part  string // часть «дома»: крыша, стены или фундамент
	// Hint — для заполненной ячейки это точка входа (через что легче
	// приходят деньги), для пустой — что доверителю нужно нарабатывать.
	Hint string
}

// Finance — разбор финансового канала по карте.
type Finance struct {
	Count4, Count5, Count6 int // сколько раз 4, 5, 6 встречаются в карте
	Total                  int // всего цифр в канале (Count4+Count5+Count6)

	Stable bool // устойчивая модель канала (2 или 4 цифры)
	Ideal  bool // идеальный, донорский канал (4 цифры и все три ячейки заполнены)

	// Level — степень влияния негативных родовых программ на деньги:
	// 0 — программы нет (идеальный канал), 1 — наименьшая (3 цифры или 4 при пустой ячейке),
	// 2 — средняя (2 цифры), 3 — высокая (1 цифра, 5 и более, либо пусто).
	Level int

	Model string // краткое описание модели канала

	Entries []FinanceCell // точки входа: заполненные ячейки среди 4, 5, 6
	Missing []FinanceCell // пустые ячейки — что нужно нарабатывать
}

var financeParts = map[int]string{4: "крыша", 5: "стены", 6: "фундамент"}

// точка входа для заполненной ячейки
var financeEntryHint = map[int]string{
	4: "тело и физические возможности: спорт, бьюти, телесные практики, медитации",
	5: "общение и сообщество: торговля, посредничество, право, блогерство",
	6: "труд и мастерство: экспертиза, работа с ресурсами, бухгалтерия",
}

// что нарабатывать, если ячейки нет
var financeMissingHint = map[int]string{
	4: "заняться физическим и ментальным здоровьем: ЗОЖ, лёгкие нагрузки, медитации",
	5: "проявляться в социуме: сообщество, общение, блогерство",
	6: "больше трудиться",
}

// Finance считает разбор финансового канала по готовой карте.
// Используются рабочие цифры доверителя (r.Counts), без добавок «после 21
// года»: описание канала делают именно по основной карте.
func (r Result) Finance() Finance {
	f := Finance{
		Count4: r.Counts[4],
		Count5: r.Counts[5],
		Count6: r.Counts[6],
	}
	f.Total = f.Count4 + f.Count5 + f.Count6
	counts := map[int]int{4: f.Count4, 5: f.Count5, 6: f.Count6}

	for _, d := range []int{4, 5, 6} {
		c := FinanceCell{Digit: d, Count: counts[d], Part: financeParts[d]}
		if counts[d] > 0 {
			c.Hint = financeEntryHint[d]
			f.Entries = append(f.Entries, c)
		} else {
			c.Hint = financeMissingHint[d]
			f.Missing = append(f.Missing, c)
		}
	}

	allCells := f.Count4 > 0 && f.Count5 > 0 && f.Count6 > 0
	f.Stable = f.Total == 2 || f.Total == 4
	f.Ideal = f.Total == 4 && allCells

	switch {
	case f.Ideal:
		f.Level = 0 // идеальный канал — негативной денежной программы нет
	case f.Total == 4:
		f.Level = 1 // четыре цифры, но одна ячейка пустая: программа есть, слабая
	case f.Total == 3:
		f.Level = 1
	case f.Total == 2:
		f.Level = 2
	default: // 1 цифра, 5 и более, либо пустой канал
		f.Level = 3
	}

	switch {
	case f.Total == 0:
		f.Model = "Пустой финансовый канал: деньги даются тяжело, их нужно нарабатывать"
	case f.Total == 1:
		f.Model = "Очень слабый, неустойчивый канал: нужен постоянный «допинг» и много усилий"
	case f.Total == 2:
		f.Model = "Устойчивый, но слабый канал: средний достаток"
	case f.Total == 3:
		f.Model = "Неустойчивый канал («заработал — потерял»), но человек хорошо зарабатывает"
	case f.Ideal:
		f.Model = "Идеальный, донорский канал: устойчивый и полный"
	case f.Total == 4:
		f.Model = "Четыре цифры, но с перекосом: одна из ячеек пустует"
	default: // 5 и более
		f.Model = "Перегруженный канал: считается как пустой"
	}

	return f
}
