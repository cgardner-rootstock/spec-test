// Package tempconv provides functions for converting temperatures between
// Celsius, Fahrenheit, and Kelvin.
package tempconv

// AbsoluteZeroC is the lowest possible temperature in Celsius.
const AbsoluteZeroC = -273.15

// AbsoluteZeroF is the lowest possible temperature in Fahrenheit.
const AbsoluteZeroF = -459.67

// CelsiusToFahrenheit converts a temperature in Celsius to Fahrenheit.
func CelsiusToFahrenheit(c float64) float64 {
	return c*9/5 + 32
}

// CelsiusToKelvin converts a temperature in Celsius to Kelvin.
func CelsiusToKelvin(c float64) float64 {
	return c - AbsoluteZeroC
}

// FahrenheitToCelsius converts a temperature in Fahrenheit to Celsius.
func FahrenheitToCelsius(f float64) float64 {
	return (f - 32) * 5 / 9
}

// FahrenheitToKelvin converts a temperature in Fahrenheit to Kelvin.
func FahrenheitToKelvin(f float64) float64 {
	return FahrenheitToCelsius(f) - AbsoluteZeroC
}

// KelvinToCelsius converts a temperature in Kelvin to Celsius.
func KelvinToCelsius(k float64) float64 {
	return k + AbsoluteZeroC
}

// KelvinToFahrenheit converts a temperature in Kelvin to Fahrenheit.
func KelvinToFahrenheit(k float64) float64 {
	return CelsiusToFahrenheit(KelvinToCelsius(k))
}
