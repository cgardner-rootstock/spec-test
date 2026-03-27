# Temperature Conversion Module — Plan

## Overview

A Go module (`tempconv.go`) exposing six pure conversion functions between Celsius, Fahrenheit, and Kelvin, plus a companion test file (`tempconv_test.go`) with table-driven tests covering every conversion direction and a shared set of physically meaningful edge cases.

---

## Conversion Formulas

| Direction         | Formula                          |
|-------------------|----------------------------------|
| C → F             | `F = C × 9/5 + 32`              |
| C → K             | `K = C + 273.15`                 |
| F → C             | `C = (F − 32) × 5/9`            |
| F → K             | `K = (F − 32) × 5/9 + 273.15`  |
| K → C             | `C = K − 273.15`                 |
| K → F             | `F = (K − 273.15) × 9/5 + 32`  |

---

## Function Signatures

```go
package tempconv

// CelsiusToFahrenheit converts a temperature in Celsius to Fahrenheit.
func CelsiusToFahrenheit(c float64) float64

// CelsiusToKelvin converts a temperature in Celsius to Kelvin.
func CelsiusToKelvin(c float64) float64

// FahrenheitToCelsius converts a temperature in Fahrenheit to Celsius.
func FahrenheitToCelsius(f float64) float64

// FahrenheitToKelvin converts a temperature in Fahrenheit to Kelvin.
func FahrenheitToKelvin(f float64) float64

// KelvinToCelsius converts a temperature in Kelvin to Celsius.
func KelvinToCelsius(k float64) float64

// KelvinToFahrenheit converts a temperature in Kelvin to Fahrenheit.
func KelvinToFahrenheit(k float64) float64
```

All functions accept and return `float64`. No error is returned; the functions are mathematically defined for all real-number inputs and callers are trusted to supply physically meaningful values.

---

## Edge Cases (applied to every conversion direction)

| Name              | Celsius    | Fahrenheit | Kelvin   | Notes                                          |
|-------------------|------------|------------|----------|------------------------------------------------|
| Absolute zero     | −273.15    | −459.67    | 0        | Lower bound of thermodynamic temperature       |
| -40 intersection  | −40.00     | −40.00     | 233.15   | Only point where °C and °F are numerically equal |
| Freezing point    | 0.00       | 32.00      | 273.15   | Water freezes at standard pressure             |
| Body temperature  | 37.00      | 98.60      | 310.15   | Normal human body temperature                  |
| Boiling point     | 100.00     | 212.00     | 373.15   | Water boils at standard pressure               |

Each of the six functions gets a table-driven test exercising all five rows above, giving **30 test cases** in total.

---

## Test Strategy

- **Package**: `tempconv_test` (black-box, exercises only exported API).
- **Structure**: one `TestXxx` function per conversion direction, each driven by a `[]struct{ input, want float64 }` table.
- **Tolerance**: floating-point results are compared with a tolerance of `1e-9` to accommodate rounding in the formulas (e.g. 5/9 is irrational in binary).
- **Helper**: a shared `assertNear(t, got, want, tol float64)` function to avoid repetition.

---

## File Layout

```
tempconv.go        # package tempconv — six conversion functions
tempconv_test.go   # package tempconv_test — table-driven tests
go.mod             # module declaration (module tempconv, go 1.22)
plan.md            # this file
```

---

## Implementation Notes

- Use the named constants `AbsoluteZeroC = -273.15` and `AbsoluteZeroF = -459.67` inside the package for clarity; they also serve as self-documenting anchors in the code.
- No third-party dependencies; `testing` from the standard library is sufficient.
