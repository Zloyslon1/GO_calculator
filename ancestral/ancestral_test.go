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

// Пример из урока модуля 3 «Родовой финансовый канал»: дата 15.08.2001,
// карта 1111/-/7, 2/5/88, 3/6/99. В финансовом канале две цифры (5 и 6),
// значит устойчивый, но слабый канал с двумя точками входа, и нет цифры 4.
func TestFinanceFromCourse(t *testing.T) {
	r, err := Calc(d(15, 8, 2001))
	if err != nil {
		t.Fatal(err)
	}
	f := r.Finance()
	if f.Count4 != 0 || f.Count5 != 1 || f.Count6 != 1 {
		t.Fatalf("цифры канала: 4×%d 5×%d 6×%d, ждали 4×0 5×1 6×1", f.Count4, f.Count5, f.Count6)
	}
	if f.Total != 2 {
		t.Errorf("всего цифр в канале %d, ждали 2", f.Total)
	}
	if !f.Stable {
		t.Error("канал из 2 цифр должен быть устойчивым")
	}
	if f.Ideal {
		t.Error("канал из 2 цифр не идеальный")
	}
	if f.Level != 2 {
		t.Errorf("степень влияния %d, ждали 2 (средняя)", f.Level)
	}
	if len(f.Entries) != 2 {
		t.Errorf("точек входа %d, ждали 2", len(f.Entries))
	}
	if len(f.Missing) != 1 || f.Missing[0].Digit != 4 {
		t.Errorf("не хватать должно только цифры 4, получили %+v", f.Missing)
	}
}

// Классификация модели канала по количеству цифр 4, 5, 6.
func TestFinanceModels(t *testing.T) {
	tests := []struct {
		name           string
		date           time.Time
		total          int
		stable, ideal  bool
		level, entries int
	}{
		// 01.12.1983: карта 111/-/7, 222/55/8, 33/-/9 → канал только 55 (две 5).
		{"две цифры в одной ячейке", d(1, 12, 1983), 2, true, false, 2, 1},
		// 13.01.1800: в канале 44 и 55, цифры 6 нет — четыре цифры, но не идеальный канал.
		{"четыре цифры, одна ячейка пустая", d(13, 1, 1800), 4, true, false, 1, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r, err := Calc(tt.date)
			if err != nil {
				t.Fatal(err)
			}
			f := r.Finance()
			if f.Total != tt.total {
				t.Errorf("всего цифр %d, ждали %d", f.Total, tt.total)
			}
			if f.Stable != tt.stable {
				t.Errorf("устойчивость %v, ждали %v", f.Stable, tt.stable)
			}
			if f.Ideal != tt.ideal {
				t.Errorf("идеальность %v, ждали %v", f.Ideal, tt.ideal)
			}
			if f.Level != tt.level {
				t.Errorf("уровень %d, ждали %d", f.Level, tt.level)
			}
			if len(f.Entries) != tt.entries {
				t.Errorf("точек входа %d, ждали %d", len(f.Entries), tt.entries)
			}
		})
	}
}
