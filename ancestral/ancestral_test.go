package ancestral

import (
	"strconv"
	"strings"
	"testing"
	"time"
)

func d(day, month, year int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

func working(r Result) string {
	parts := make([]string, len(r.Working))
	for i, w := range r.Working {
		parts[i] = strconv.Itoa(w)
	}
	return strings.Join(parts, ".")
}

// Примеры из PDF «Ансестологическая карта». Ячейки идут по строкам:
// 1 4 7 / 2 5 8 / 3 6 9, «-» означает, что цифры нет.
func TestCardsFromCourse(t *testing.T) {
	tests := []struct {
		name    string
		date    time.Time
		working string
		card    [9]string // 1 4 7 2 5 8 3 6 9 по строкам
		nine    bool
	}{
		{"обычная", d(1, 12, 1983), "25.7.23.5",
			[9]string{"111", "-", "7", "222", "55", "8", "33", "-", "9"}, false},
		{"№3 равно 10", d(9, 1, 1980), "28.10.10.10",
			[9]string{"11111", "-", "-", "2", "-", "88", "-", "-", "99"}, false},
		{"после 2000", d(4, 5, 2012), "14.5.19.33.6",
			[9]string{"111", "44", "-", "22", "55", "-", "33", "6", "9"}, false},
		{"с 2020", d(5, 11, 2021), "12.3.19.31.4",
			[9]string{"111111", "4", "-", "222", "5", "-", "33", "-", "99"}, true},
		{"№1 однозначное", d(2, 10, 2020), "7.7.19.26.8",
			[9]string{"11", "-", "77", "2222", "-", "8", "-", "6", "99"}, true},
		{"№4 двузначное", d(19, 9, 2018), "30.3.19.49.4",
			[9]string{"111", "44", "-", "2", "-", "8", "33", "-", "9999"}, false},
		{"1800-е", d(14, 2, 1800), "16.7.9.16.7",
			[9]string{"1111", "4", "77", "2", "-", "8", "-", "66", "9"}, false},
	}
	// порядок ячеек в массиве card -> цифра
	order := [9]int{1, 4, 7, 2, 5, 8, 3, 6, 9}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Calc(tt.date)
			if err != nil {
				t.Fatal(err)
			}
			if got := working(r); got != tt.working {
				t.Errorf("рабочие числа: получили %s, ждали %s", got, tt.working)
			}
			for i, digit := range order {
				want := tt.card[i]
				got := strings.Repeat(string(rune('0'+digit)), r.Counts[digit])
				if r.Counts[digit] == 0 {
					got = "-"
				}
				if got != want {
					t.Errorf("цифра %d: получили %q, ждали %q", digit, got, want)
				}
			}
			if r.ExtraNine != tt.nine {
				t.Errorf("вторая константа 9: получили %v, ждали %v", r.ExtraNine, tt.nine)
			}
		})
	}
}

// Контроль из PDF: цифр в дате и рабочих числах должно быть столько же,
// сколько в таблице (13 цифр в примере 01.12.1983).
func TestDigitCount(t *testing.T) {
	r, _ := Calc(d(1, 12, 1983))
	total := 0
	for i := 1; i <= 9; i++ {
		total += r.Counts[i]
	}
	if total != 13 {
		t.Errorf("всего цифр в карте %d, ждали 13", total)
	}
}

// Карта после 21 года: +1 ко всем цифрам, кроме 5.
func TestAfter21(t *testing.T) {
	r, _ := Calc(d(4, 5, 2012))
	if r.At21 == nil {
		t.Fatal("для 2012 года ждали карту после 21 года")
	}
	want := map[int]int{1: 4, 2: 3, 3: 3, 4: 3, 5: 2, 6: 2, 7: 1, 8: 1, 9: 2}
	for digit, n := range want {
		if r.At21[digit] != n {
			t.Errorf("цифра %d: получили %d, ждали %d", digit, r.At21[digit], n)
		}
	}
	if r2, _ := Calc(d(1, 12, 1983)); r2.At21 != nil {
		t.Error("для 1983 года карты после 21 года быть не должно")
	}
}

// Цепочки рабочих чисел из PDF «Определение предназначения».
func TestWorkingChainsFromCourse(t *testing.T) {
	tests := []struct {
		date    time.Time
		working string
		chrp    int
	}{
		{d(23, 10, 1966), "28.10.24.6", 10},
		{d(12, 2, 1975), "27.9.25.7", 9},
		{d(1, 12, 1991), "24.6.22.4", 6},
		{d(5, 10, 1992), "27.9.17.8", 9},
	}
	for _, tt := range tests {
		r, err := Calc(tt.date)
		if err != nil {
			t.Fatal(err)
		}
		if got := working(r); got != tt.working {
			t.Errorf("%s: получили %s, ждали %s", tt.date.Format("02.01.2006"), got, tt.working)
		}
		if r.ChRP != tt.chrp {
			t.Errorf("%s: ЧРП %d, ждали %d", tt.date.Format("02.01.2006"), r.ChRP, tt.chrp)
		}
	}
}

func TestOutOfRange(t *testing.T) {
	for _, y := range []int{1799, 2100} {
		if _, err := Calc(d(1, 1, y)); err == nil {
			t.Errorf("год %d: ждали ошибку", y)
		}
	}
}
