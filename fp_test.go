package bls

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"
)

func TestFp(t *testing.T) {
	field := NewFp()
	zero := &Fe{0}
	one := &Fe{1}
	t.Run("Encoding & Decoding", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			bytes := []byte{0}
			fe := &Fe{}
			fe.FromBytes(bytes)
			if !field.Equal(fe, zero) {
				t.Errorf("bad encoding\n")
			}
		})
		t.Run("2", func(t *testing.T) {
			in := []byte{254, 253}
			fe := &Fe{}
			fe.FromBytes(in)
			if bytes.Equal(in, fe.Bytes()) {
				t.Errorf("bad encoding\n")
			}
		})
		t.Run("3", func(t *testing.T) {
			a, _ := field.RandElement(&Fe{}, rand.Reader)
			b := &Fe{}
			b.FromBytes(a.Bytes())
			if !field.Equal(a, b) {
				t.Errorf("bad encoding or decoding\n")
			}
		})
		t.Run("4", func(t *testing.T) {
			a, _ := field.RandElement(&Fe{}, rand.Reader)
			b := &Fe{}
			if _, err := b.SetString(a.String()); err != nil {
				t.Errorf("bad encoding or decoding\n")
			}
			if !field.Equal(a, b) {
				t.Errorf("bad encoding or decoding\n")
			}
		})
		t.Run("5", func(t *testing.T) {
			a, _ := field.RandElement(&Fe{}, rand.Reader)
			b := &Fe{}
			b.SetBig(a.Big())
			if !field.Equal(a, b) {
				t.Errorf("bad encoding or decoding\n")
			}
		})
	})
	t.Run("Addition", func(t *testing.T) {
		var a, b, c, u, v *Fe
		for i := 0; i < n; i++ {
			u = &Fe{}
			v = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			b, _ = field.RandElement(&Fe{}, rand.Reader)
			c, _ = field.RandElement(&Fe{}, rand.Reader)
			field.Add(u, a, b)
			field.Add(u, u, c)
			field.Add(v, b, c)
			field.Add(v, v, a)
			if !field.Equal(u, v) {
				t.Fatalf("Additive associativity does not hold")
			}
			field.Add(u, a, b)
			field.Add(v, b, a)
			if !field.Equal(u, v) {
				t.Fatalf("Additive commutativity does not hold")
			}
			field.Add(u, a, zero)
			if !field.Equal(u, a) {
				t.Fatalf("Additive identity does not hold")
			}
			field.Neg(u, a)
			field.Add(u, u, a)
			if !field.Equal(u, zero) {
				t.Fatalf("Bad Negation\na:%s", a.String())
			}
		}
	})
	t.Run("Doubling", func(t *testing.T) {
		var a, u, v *Fe
		for j := 0; j < n; j++ {
			u = &Fe{}
			v = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			field.Double(u, a)
			field.Add(v, a, a)
			if !field.Equal(u, v) {
				t.Fatalf("Bad doubling\na: %s\nu: %s\nv: %s\n", a, u, v)
			}
		}
	})
	t.Run("Subtraction", func(t *testing.T) {
		var a, b, c, u, v *Fe
		for j := 0; j < n; j++ {
			u = &Fe{}
			v = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			b, _ = field.RandElement(&Fe{}, rand.Reader)
			c, _ = field.RandElement(&Fe{}, rand.Reader)
			field.Sub(u, a, c)
			field.Sub(u, u, b)
			field.Sub(v, a, b)
			field.Sub(v, v, c)
			if !field.Equal(u, v) {
				t.Fatalf("Additive associativity does not hold\na: %s\nb: %s\nc: %s\nu: %s\nv:%s\n", a, b, c, u, v)
			}
			field.Sub(u, a, zero)
			if !field.Equal(u, a) {
				t.Fatalf("Additive identity does not hold\na: %s\nu: %s\n", a, u)
			}
			field.Sub(u, a, b)
			field.Sub(v, b, a)
			field.Add(u, u, v)
			if !field.Equal(u, zero) {
				t.Fatalf("Additive commutativity does not hold\na: %s\nb: %s\nu: %s\nv: %s", a, b, u, v)
			}
			field.Sub(u, a, b)
			field.Sub(v, b, a)
			field.Neg(v, v)
			if !field.Equal(u, u) {
				t.Fatalf("Bad Negation\na:%s", a.String())
			}
		}
	})
	t.Run("Montgomerry", func(t *testing.T) {
		var a, b, c, u, v, w *Fe
		for j := 0; j < n; j++ {
			u = &Fe{}
			v = &Fe{}
			w = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			b, _ = field.RandElement(&Fe{}, rand.Reader)
			c, _ = field.RandElement(&Fe{}, rand.Reader)
			field.Mont(u, zero)
			if !field.Equal(u, zero) {
				t.Fatalf("Bad Montgomerry encoding")
			}
			field.Demont(u, zero)
			if !field.Equal(u, zero) {
				t.Fatalf("Bad Montgomerry decoding")
			}
			field.Mont(u, one)
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad Montgomerry encoding")
			}
			field.Demont(u, field.One())
			if !field.Equal(u, one) {
				t.Fatalf("Bad Montgomerry decoding")
			}
			field.Mul(u, a, zero)
			if !field.Equal(u, zero) {
				t.Fatalf("Bad zero element")
			}
			field.Mul(u, a, one)
			field.Mul(u, u, r2)
			if !field.Equal(u, a) {
				t.Fatalf("Multiplication identity does not hold")
			}
			field.Mul(u, r2, one)
			if !field.Equal(u, field.One()) {
				t.Fatalf("Multiplication identity does not hold, expected to equal r1")
			}
			field.Mul(u, a, b)
			field.Mul(u, u, c)
			field.Mul(v, b, c)
			field.Mul(v, v, a)
			if !field.Equal(u, v) {
				t.Fatalf("Multiplicative associativity does not hold")
			}
			field.Add(u, a, b)
			field.Mul(u, c, u)
			field.Mul(w, a, c)
			field.Mul(v, b, c)
			field.Add(v, v, w)
			if !field.Equal(u, v) {
				t.Fatalf("Distributivity does not hold")
			}
			field.Square(u, a)
			field.Mul(v, a, a)
			if !field.Equal(u, v) {
				t.Fatalf("Bad squaring")
			}
		}
	})
	t.Run("Exponentiation", func(t *testing.T) {
		var a, u, v *Fe
		for j := 0; j < n; j++ {
			u = &Fe{}
			v = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			field.Exp(u, a, big.NewInt(0))
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad exponentiation, expected to equal r1")
			}
			field.Exp(u, a, big.NewInt(1))
			if !field.Equal(u, a) {
				t.Fatalf("Bad exponentiation, expected to equal a")
			}
			field.Mul(u, a, a)
			field.Mul(u, u, u)
			field.Mul(u, u, u)
			field.Exp(v, a, big.NewInt(8))
			if !field.Equal(u, v) {
				t.Fatalf("Bad exponentiation")
			}
			p := new(big.Int).SetBytes(modulus.Bytes())
			field.Exp(u, a, p)
			if !field.Equal(u, a) {
				t.Fatalf("Bad exponentiation, expected to equal itself")
			}
			field.Exp(u, a, p.Sub(p, big.NewInt(1)))
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad exponentiation, expected to equal r1")
			}
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		var a, u, v *Fe
		for j := 0; j < n; j++ {
			u = &Fe{}
			v = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			field.InvMontUp(u, a)
			field.Mul(u, u, a)
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad inversion, expected to equal r1")
			}
			field.Mont(u, a)
			field.InvMontDown(v, u)
			field.Mul(v, v, u)
			if !field.Equal(v, one) {
				t.Fatalf("Bad inversion, expected to equal 1")
			}
			p := new(big.Int).SetBytes(modulus.Bytes())
			field.Exp(u, a, p.Sub(p, big.NewInt(2)))
			field.InvMontUp(v, a)
			if !field.Equal(v, u) {
				t.Fatalf("Bad inversion 1")
			}
			field.InvEEA(u, a)
			field.Mul(u, u, a)
			field.Mul(u, u, r2)
			if !field.Equal(u, one) {
				t.Fatalf("Bad inversion 2")
			}
		}
	})
	t.Run("Sqrt", func(t *testing.T) {
		r := &Fe{}
		if field.Sqrt(r, nonResidue1) {
			t.Fatalf("bad sqrt 1")
		}
		for j := 0; j < n; j++ {
			a, _ := field.RandElement(&Fe{}, rand.Reader)
			aa, rr, r := &Fe{}, &Fe{}, &Fe{}
			field.Square(aa, a)
			if !field.Sqrt(r, aa) {
				t.Fatalf("bad sqrt 2")
			}
			field.Square(rr, r)
			if !field.Equal(rr, aa) {
				t.Fatalf("bad sqrt 3")
			}
		}
	})
}

func TestFp2(t *testing.T) {
	field := NewFp2(nil)
	t.Run("Encoding & Decoding", func(t *testing.T) {
		in := make([]byte, 96)
		for i := 0; i < 96; i++ {
			in[i] = 1
		}
		fe := &Fe2{}
		if err := field.NewElementFromBytes(fe, in); err != nil {
			panic(err)
		}
		if !bytes.Equal(in, field.ToBytes(fe)) {
			t.Errorf("bad encoding\n")
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		var a, b, c, u, v, w *Fe2
		for j := 0; j < n; j++ {
			u = &Fe2{}
			v = &Fe2{}
			w = &Fe2{}
			a, _ = field.RandElement(&Fe2{}, rand.Reader)
			b, _ = field.RandElement(&Fe2{}, rand.Reader)
			c, _ = field.RandElement(&Fe2{}, rand.Reader)
			field.Mul(u, a, b)
			field.Mul(u, u, c)
			field.Mul(v, b, c)
			field.Mul(v, v, a)
			if !field.Equal(u, v) {
				t.Fatalf("Multiplicative associativity does not hold")
			}
			field.Add(u, a, b)
			field.Mul(u, c, u)
			field.Mul(w, a, c)
			field.Mul(v, b, c)
			field.Add(v, v, w)
			if !field.Equal(u, v) {
				t.Fatalf("Distributivity does not hold")
			}
			field.Square(u, a)
			field.Mul(v, a, a)
			if !field.Equal(u, v) {
				t.Fatalf("Bad squaring")
			}
		}
	})
	t.Run("Exponentiation", func(t *testing.T) {
		var a, u, v *Fe2
		for j := 0; j < n; j++ {
			u = &Fe2{}
			v = &Fe2{}
			a, _ = field.RandElement(&Fe2{}, rand.Reader)
			field.Exp(u, a, big.NewInt(0))
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad exponentiation, expected to equal r1")
			}
			_ = v
			field.Exp(u, a, big.NewInt(1))
			if !field.Equal(u, a) {
				t.Fatalf("Bad exponentiation, expected to equal a")
			}
			field.Mul(u, a, a)
			field.Mul(u, u, u)
			field.Mul(u, u, u)
			field.Exp(v, a, big.NewInt(8))
			if !field.Equal(u, v) {
				t.Fatalf("Bad exponentiation")
			}
			// p := new(big.Int).SetBytes(modulus.Bytes())
			// field.Exp(u, a, p)
			// if !field.Equal(u, a) {
			// 	t.Fatalf("Bad exponentiation, expected to equal itself")
			// }
			// field.Exp(u, a, p.Sub(p, big.NewInt(1)))
			// if !field.Equal(u, field.One()) {
			// 	t.Fatalf("Bad exponentiation, expected to equal one")
			// }
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		var a, u *Fe2
		for j := 0; j < n; j++ {
			u = &Fe2{}
			a, _ = field.RandElement(&Fe2{}, rand.Reader)
			field.Inverse(u, a)
			field.Mul(u, u, a)
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad inversion, expected to equal r1")
			}
		}
	})
	t.Run("Sqrt", func(t *testing.T) {
		r := &Fe2{}
		if field.Sqrt(r, nonResidue2) {
			t.Fatalf("bad sqrt 1")
		}
		for j := 0; j < n; j++ {
			a, _ := field.RandElement(&Fe2{}, rand.Reader)
			aa, rr, r := &Fe2{}, &Fe2{}, &Fe2{}
			field.Square(aa, a)
			if !field.Sqrt(r, aa) {
				t.Fatalf("bad sqrt 2")
			}
			field.Square(rr, r)
			if !field.Equal(rr, aa) {
				t.Fatalf("bad sqrt 3")
			}
		}
	})
}

func TestFp6(t *testing.T) {
	field := NewFp6(nil)
	// zero := field.Zero()
	t.Run("Encoding & Decoding", func(t *testing.T) {
		in := make([]byte, 288)
		for i := 0; i < 288; i++ {
			in[i] = 1
		}
		fe := &Fe6{}
		if err := field.NewElementFromBytes(fe, in); err != nil {
			panic(err)
		}
		if !bytes.Equal(in, field.ToBytes(fe)) {
			t.Errorf("bad encoding\n")
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		var a, b, c, u, v, w *Fe6
		for j := 0; j < n; j++ {
			u = &Fe6{}
			v = &Fe6{}
			w = &Fe6{}
			a, _ = field.RandElement(&Fe6{}, rand.Reader)
			b, _ = field.RandElement(&Fe6{}, rand.Reader)
			c, _ = field.RandElement(&Fe6{}, rand.Reader)
			field.Mul(u, a, b)
			field.Mul(u, u, c)
			field.Mul(v, b, c)
			field.Mul(v, v, a)
			if !field.Equal(u, v) {
				t.Fatalf("Multiplicative associativity does not hold")
			}
			field.Add(u, a, b)
			field.Mul(u, c, u)
			field.Mul(w, a, c)
			field.Mul(v, b, c)
			field.Add(v, v, w)
			if !field.Equal(u, v) {
				t.Fatalf("Distributivity does not hold")
			}
			field.Square(u, a)
			field.Mul(v, a, a)
			if !field.Equal(u, v) {
				t.Fatalf("Bad squaring")
			}
		}
	})
	t.Run("Exponentiation", func(t *testing.T) {
		var a, u, v *Fe6
		for j := 0; j < n; j++ {
			u = &Fe6{}
			v = &Fe6{}
			a, _ = field.RandElement(&Fe6{}, rand.Reader)
			field.Exp(u, a, big.NewInt(0))
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad exponentiation, expected to equal r1")
			}
			_ = v
			field.Exp(u, a, big.NewInt(1))
			if !field.Equal(u, a) {
				t.Fatalf("Bad exponentiation, expected to equal a")
			}
			field.Mul(u, a, a)
			field.Mul(u, u, u)
			field.Mul(u, u, u)
			field.Exp(v, a, big.NewInt(8))
			if !field.Equal(u, v) {
				t.Fatalf("Bad exponentiation")
			}
			// p := new(big.Int).SetBytes(modulus.Bytes())
			// field.Exp(u, a, p)
			// if !field.Equal(u, a) {
			// 	t.Fatalf("Bad exponentiation, expected to equal itself")
			// }
			// field.Exp(u, a, p.Sub(p, big.NewInt(1)))
			// if !field.Equal(u, field.One()) {
			// 	t.Fatalf("Bad exponentiation, expected to equal one")
			// }
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		var a, u *Fe6
		for j := 0; j < n; j++ {
			u = &Fe6{}
			a, _ = field.RandElement(&Fe6{}, rand.Reader)
			field.Inverse(u, a)
			field.Mul(u, u, a)
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad inversion, expected to equal r1")
			}
		}
	})
	t.Run("MulBy01", func(t *testing.T) {
		fq2 := field.f
		var a, b, u *Fe6
		for j := 0; j < n; j++ {
			a, _ = field.RandElement(&Fe6{}, rand.Reader)
			b, _ = field.RandElement(&Fe6{}, rand.Reader)
			u, _ = field.RandElement(&Fe6{}, rand.Reader)
			fq2.Copy(&b[2], fq2.Zero())
			field.Mul(u, a, b)
			field.MulBy01(a, &b[0], &b[1])
			if !field.Equal(a, u) {
				t.Fatal("Bad mul by 01")
			}
		}
	})
	t.Run("MulBy1", func(t *testing.T) {
		fq2 := field.f
		var a, b, u *Fe6
		for j := 0; j < n; j++ {
			a, _ = field.RandElement(&Fe6{}, rand.Reader)
			b, _ = field.RandElement(&Fe6{}, rand.Reader)
			u, _ = field.RandElement(&Fe6{}, rand.Reader)
			fq2.Copy(&b[2], fq2.Zero())
			fq2.Copy(&b[0], fq2.Zero())
			field.Mul(u, a, b)
			field.MulBy1(a, &b[1])
			if !field.Equal(a, u) {
				t.Fatal("Bad mul by 1")
			}
		}
	})
}

func TestFp12(t *testing.T) {
	field := NewFp12(nil)
	t.Run("Encoding & Decoding", func(t *testing.T) {
		in := make([]byte, 576)
		for i := 0; i < 288; i++ {
			in[i] = 1
		}
		fe := &Fe12{}
		if err := field.NewElementFromBytes(fe, in); err != nil {
			panic(err)
		}
		if !bytes.Equal(in, field.ToBytes(fe)) {
			t.Errorf("bad encoding\n")
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		var a, b, c, u, v, w *Fe12
		for j := 0; j < n; j++ {
			u = &Fe12{}
			v = &Fe12{}
			w = &Fe12{}
			a, _ = field.RandElement(&Fe12{}, rand.Reader)
			b, _ = field.RandElement(&Fe12{}, rand.Reader)
			c, _ = field.RandElement(&Fe12{}, rand.Reader)
			field.Mul(u, a, b)
			field.Mul(u, u, c)
			field.Mul(v, b, c)
			field.Mul(v, v, a)
			if !field.Equal(u, v) {
				t.Fatalf("Multiplicative associativity does not hold")
			}
			field.Add(u, a, b)
			field.Mul(u, c, u)
			field.Mul(w, a, c)
			field.Mul(v, b, c)
			field.Add(v, v, w)
			if !field.Equal(u, v) {
				t.Fatalf("Distributivity does not hold")
			}
			field.Square(u, a)
			field.Mul(v, a, a)
			if !field.Equal(u, v) {
				t.Fatalf("Bad squaring")
			}
		}
	})
	t.Run("Exponentiation", func(t *testing.T) {
		var a, u, v *Fe12
		for j := 0; j < n; j++ {
			u = &Fe12{}
			v = &Fe12{}
			a, _ = field.RandElement(&Fe12{}, rand.Reader)
			field.Exp(u, a, big.NewInt(0))
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad exponentiation, expected to equal r1")
			}
			_ = v
			field.Exp(u, a, big.NewInt(1))
			if !field.Equal(u, a) {
				t.Fatalf("Bad exponentiation, expected to equal a")
			}
			field.Mul(u, a, a)
			field.Mul(u, u, u)
			field.Mul(u, u, u)
			field.Exp(v, a, big.NewInt(8))
			if !field.Equal(u, v) {
				t.Fatalf("Bad exponentiation")
			}
			// p := new(big.Int).SetBytes(modulus.Bytes())
			// field.Exp(u, a, p)
			// if !field.Equal(u, a) {
			// 	t.Fatalf("Bad exponentiation, expected to equal itself")
			// }
			// field.Exp(u, a, p.Sub(p, big.NewInt(1)))
			// if !field.Equal(u, field.One()) {
			// 	t.Fatalf("Bad exponentiation, expected to equal one")
			// }
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		var a, u *Fe12
		for j := 0; j < n; j++ {
			u = &Fe12{}
			a, _ = field.RandElement(&Fe12{}, rand.Reader)
			field.Inverse(u, a)
			field.Mul(u, u, a)
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad inversion, expected to equal r1")
			}
		}
	})
	t.Run("MulBy014", func(t *testing.T) {
		fq2 := field.f.f
		var a, b, u *Fe12
		for j := 0; j < n; j++ {
			a, _ = field.RandElement(&Fe12{}, rand.Reader)
			b, _ = field.RandElement(&Fe12{}, rand.Reader)
			u, _ = field.RandElement(&Fe12{}, rand.Reader)
			fq2.Copy(&b[0][2], fq2.Zero())
			fq2.Copy(&b[1][0], fq2.Zero())
			fq2.Copy(&b[1][2], fq2.Zero())
			field.Mul(u, a, b)
			field.MulBy014Assign(a, &b[0][0], &b[0][1], &b[1][1])
			if !field.Equal(a, u) {
				t.Fatal("Bad mul by 014")
			}
		}
	})
	t.Run("MulBy034", func(t *testing.T) {
		fq2 := field.f.f
		var a, b, u *Fe12
		for j := 0; j < n; j++ {
			a, _ = field.RandElement(&Fe12{}, rand.Reader)
			b, _ = field.RandElement(&Fe12{}, rand.Reader)
			u, _ = field.RandElement(&Fe12{}, rand.Reader)
			fq2.Copy(&b[0][1], fq2.Zero())
			fq2.Copy(&b[0][2], fq2.Zero())
			fq2.Copy(&b[1][2], fq2.Zero())
			field.Mul(u, a, b)
			field.MulBy034Assign(a, &b[0][0], &b[1][0], &b[1][1])
			if !field.Equal(a, u) {
				t.Fatal("Bad mul by 034")
			}
		}
	})
}

func BenchmarkFp(t *testing.B) {
	var a, b, c Fe
	var field = NewFp()
	field.RandElement(&a, rand.Reader)
	field.RandElement(&b, rand.Reader)
	t.Run("Addition", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Add(&c, &a, &b)
		}
	})
	t.Run("Subtraction", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Sub(&c, &a, &b)
		}
	})
	t.Run("Doubling", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Double(&c, &a)
		}
	})
	t.Run("Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Mul(&c, &a, &b)
		}
	})
	t.Run("Squaring", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Square(&c, &a)
		}
	})
	t.Run("Inversion", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.InvMontUp(&c, &a)
		}
	})
	t.Run("Exponentiation", func(t *testing.B) {
		e := new(big.Int).SetBytes(modulus.Bytes())
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Exp(&c, &a, e)
		}
	})
}

func BenchmarkFp2(t *testing.B) {
	var a, b, c Fe2
	var field = NewFp2(nil)
	field.RandElement(&a, rand.Reader)
	field.RandElement(&b, rand.Reader)
	t.Run("Addition", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Add(&c, &a, &b)
		}
	})
	t.Run("Subtraction", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Sub(&c, &a, &b)
		}
	})
	t.Run("Doubling", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Double(&c, &a)
		}
	})
	t.Run("Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Mul(&c, &a, &b)
		}
	})
	t.Run("Squaring", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Square(&c, &a)
		}
	})
	t.Run("Inversion", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Inverse(&c, &a)
		}
	})
	t.Run("Exponentiation", func(t *testing.B) {
		e := new(big.Int).SetBytes(modulus.Bytes())
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Exp(&c, &a, e)
		}
	})
}
