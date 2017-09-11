package benchmarks

import "testing"

func BenchmarkFilterByFuncOrSwitch(b *testing.B) {
	type t struct {
		rx    string
		begin string
		level int
		f     func(string) bool
	}

	foo := new(t)
	foo.f = func(s string) bool {
		return s == foo.begin
	}

	testSwitch := func() {
		x := "hello"
		switch {
		case len(foo.rx) > 0:
			_ = x == foo.rx
		case foo.level > 0 && len(foo.begin) > 0:
			_ = x == foo.begin
		case len(foo.begin) > 0:
			_ = x == foo.begin
		default:
		}
	}

	b.Run("test switch", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			testSwitch()
		}
	})

	b.Run("Test f", func(b *testing.B) {
		x := "hello"
		for i := 0; i < b.N; i++ {
			_ = foo.f(x)
		}
	})
}

func BenchmarkBoolAndNil(b *testing.B) {
	type tmp struct {
		Bool bool
		Nil  *struct{}
	}
	b.Run("Bench != nil", func(b *testing.B) {
		foo := new(tmp)

		for i := 0; i < b.N; i++ {
			if foo.Nil != nil {
				_ = 1 + 1
			}
		}
	})

	b.Run("Bench bool", func(b *testing.B) {
		bar := new(tmp)
		for i := 0; i < b.N; i++ {
			if bar.Bool {
				_ = 1 + 1
			}
		}
	})
}
