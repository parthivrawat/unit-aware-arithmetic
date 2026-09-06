package dimensional

import (
	"fmt"
	"testing"
)

// Package-level sinks prevent the compiler from optimizing away benchmark loops.
var result Quantity
var resultFloat float64

func BenchmarkRawFloatAdd(b *testing.B) {
	a, c := 100.0, 0.5
	for i := 0; i < b.N; i++ {
		resultFloat = a + c
	}
}

func BenchmarkQuantityAdd(b *testing.B) {
	q1 := NewQuantity(100.0, Meter)
	q2 := NewQuantity(50.0, Centimeter)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		result, err = q1.Add(q2)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRawFloatMultiply(b *testing.B) {
	a, c := 10.0, 9.8
	for i := 0; i < b.N; i++ {
		resultFloat = a * c
	}
}

func BenchmarkQuantityMultiply(b *testing.B) {
	q1 := NewQuantity(10.0, Kilogram)
	q2 := NewQuantity(9.8, MeterPerSecondSquared)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		result, err = q1.Multiply(q2)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRawFloatDivide(b *testing.B) {
	a, c := 100.0, 10.0
	for i := 0; i < b.N; i++ {
		resultFloat = a / c
	}
}

func BenchmarkQuantityDivide(b *testing.B) {
	q1 := NewQuantity(100.0, Meter)
	q2 := NewQuantity(10.0, Second)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		result, err = q1.Divide(q2)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkRawFloatConvert(b *testing.B) {
	a := 100.0
	for i := 0; i < b.N; i++ {
		resultFloat = a * 0.001
	}
}

func BenchmarkQuantityConvert(b *testing.B) {
	q1 := NewQuantity(100.0, Meter)
	target := Kilometer
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var err error
		result, err = q1.To(target)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkQuantityNew(b *testing.B) {
	for i := 0; i < b.N; i++ {
		result = NewQuantity(100.0, Meter)
	}
}

func BenchmarkRawFloatLiteral(b *testing.B) {
	for i := 0; i < b.N; i++ {
		resultFloat = 100.0
	}
}

func ExampleQuantity_To() {
	q := NewQuantity(100.0, Meter)
	converted, _ := q.To(Kilometer)
	fmt.Println(converted)
	// Output: 0.1 km
}
