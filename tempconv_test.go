package tempconv_test

import (
	"math"
	"testing"

	"tempconv"
)

const tol = 1e-9

func assertNear(t *testing.T, got, want, tol float64) {
	t.Helper()
	if math.Abs(got-want) > tol {
		t.Errorf("got %v, want %v (tolerance %v)", got, want, tol)
	}
}

func TestCelsiusToFahrenheit(t *testing.T) {
	cases := []struct {
		input, want float64
	}{
		{-273.15, -459.67},  // absolute zero
		{-40.00, -40.00},    // -40 intersection
		{0.00, 32.00},       // freezing point
		{37.00, 98.60},      // body temperature
		{100.00, 212.00},    // boiling point
	}
	for _, c := range cases {
		assertNear(t, tempconv.CelsiusToFahrenheit(c.input), c.want, tol)
	}
}

func TestCelsiusToKelvin(t *testing.T) {
	cases := []struct {
		input, want float64
	}{
		{-273.15, 0.00},     // absolute zero
		{-40.00, 233.15},    // -40 intersection
		{0.00, 273.15},      // freezing point
		{37.00, 310.15},     // body temperature
		{100.00, 373.15},    // boiling point
	}
	for _, c := range cases {
		assertNear(t, tempconv.CelsiusToKelvin(c.input), c.want, tol)
	}
}

func TestFahrenheitToCelsius(t *testing.T) {
	cases := []struct {
		input, want float64
	}{
		{-459.67, -273.15},  // absolute zero
		{-40.00, -40.00},    // -40 intersection
		{32.00, 0.00},       // freezing point
		{98.60, 37.00},      // body temperature
		{212.00, 100.00},    // boiling point
	}
	for _, c := range cases {
		assertNear(t, tempconv.FahrenheitToCelsius(c.input), c.want, tol)
	}
}

func TestFahrenheitToKelvin(t *testing.T) {
	cases := []struct {
		input, want float64
	}{
		{-459.67, 0.00},     // absolute zero
		{-40.00, 233.15},    // -40 intersection
		{32.00, 273.15},     // freezing point
		{98.60, 310.15},     // body temperature
		{212.00, 373.15},    // boiling point
	}
	for _, c := range cases {
		assertNear(t, tempconv.FahrenheitToKelvin(c.input), c.want, tol)
	}
}

func TestKelvinToCelsius(t *testing.T) {
	cases := []struct {
		input, want float64
	}{
		{0.00, -273.15},     // absolute zero
		{233.15, -40.00},    // -40 intersection
		{273.15, 0.00},      // freezing point
		{310.15, 37.00},     // body temperature
		{373.15, 100.00},    // boiling point
	}
	for _, c := range cases {
		assertNear(t, tempconv.KelvinToCelsius(c.input), c.want, tol)
	}
}

func TestKelvinToFahrenheit(t *testing.T) {
	cases := []struct {
		input, want float64
	}{
		{0.00, -459.67},     // absolute zero
		{233.15, -40.00},    // -40 intersection
		{273.15, 32.00},     // freezing point
		{310.15, 98.60},     // body temperature
		{373.15, 212.00},    // boiling point
	}
	for _, c := range cases {
		assertNear(t, tempconv.KelvinToFahrenheit(c.input), c.want, tol)
	}
}
