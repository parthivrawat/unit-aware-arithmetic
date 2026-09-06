# Unit-Aware Arithmetic (Go)

A type-safe dimensional arithmetic library that tracks units at runtime and prevents invalid operations.

## Features

- ✅ **Type-safe dimensional analysis**: Prevents incompatible unit operations
- ✅ **Zero dependencies**: Core functionality has no external dependencies
- ✅ **Comprehensive unit coverage**: SI, imperial, and derived units
- ✅ **Arithmetic operations**: Add, subtract, multiply, divide with unit tracking
- ✅ **Unit conversion**: Automatic and explicit conversion between compatible units
- ✅ **Clear error messages**: Helpful errors for incompatible operations
- ✅ **Production-ready**: Comprehensive test coverage (>95%)
- ✅ **Idiomatic Go**: Follows Go best practices and conventions

## Installation

```bash
go get github.com/parthivrawat/unit-aware-arithmetic/go/v2
```

## Quick Start

```go
package main

import (
	"fmt"
	dim "github.com/parthivrawat/unit-aware-arithmetic/go/v2"
)

func main() {
	// Create quantities with units
	distance := dim.NewQuantity(100, dim.Meter)
	time := dim.NewQuantity(9.58, dim.Second)

	// Arithmetic operations with automatic unit tracking
	speed, err := distance.Divide(time)
	if err != nil {
		panic(err)
	}
	fmt.Println(speed) // 10.438413361169102 m/s

	// Unit conversion
	distanceKm, _ := distance.To(dim.Kilometer)
	fmt.Println(distanceKm) // 0.1 km

	// Type-safe operations - this will return an error!
	_, err := distance.Add(time)
	if err != nil {
		fmt.Println(err) // Cannot add m and s: incompatible dimensions
	}
}
```

## Usage Examples

### Basic Arithmetic

```go
import dim "github.com/parthivrawat/unit-aware-arithmetic/go/v2"

// Addition (same dimension required)
d1 := dim.NewQuantity(5, dim.Meter)
d2 := dim.NewQuantity(3, dim.Meter)
total, _ := d1.Add(d2) // 8.0 m

// Multiplication creates derived units
area, _ := dim.NewQuantity(5, dim.Meter).Multiply(dim.NewQuantity(3, dim.Meter))
fmt.Println(area) // 15 m·m

// Scalar multiplication keeps the unit
double := dim.NewQuantity(5, dim.Meter).MultiplyScalar(2.0)
fmt.Println(double) // 10 m

// Division creates derived units
velocity, _ := dim.NewQuantity(100, dim.Meter).Divide(dim.NewQuantity(10, dim.Second))
fmt.Println(velocity) // 10 m/s
```

### Unit Conversion

```go
import dim "github.com/parthivrawat/unit-aware-arithmetic/go/v2"

// Length conversion
distance := dim.NewQuantity(1, dim.Mile)
distanceKm, _ := distance.To(dim.Kilometer)
fmt.Println(distanceKm) // 1.60934 km

// Temperature conversion
tempC := dim.NewQuantity(0, dim.Celsius)
tempK, _ := tempC.To(dim.Kelvin)
fmt.Println(tempK) // 273.15 K

// Mass conversion
massLb := dim.NewQuantity(10, dim.Pound)
massKg, _ := massLb.To(dim.Kilogram)
fmt.Println(massKg) // 4.53592 kg
```

### Physics Calculations

```go
import dim "github.com/parthivrawat/unit-aware-arithmetic/go/v2"

// Calculate velocity
distance := dim.NewQuantity(100, dim.Meter)
time := dim.NewQuantity(9.58, dim.Second)
velocity, _ := distance.Divide(time)
fmt.Printf("Velocity: %v\n", velocity)

// Calculate force (F = ma)
mass := dim.NewQuantity(10, dim.Kilogram)
acceleration := dim.NewQuantity(9.8, dim.MeterPerSecondSquared)
force, _ := mass.Multiply(acceleration)
fmt.Printf("Force: %v\n", force)

// Calculate kinetic energy (KE = 1/2 * m * v²)
mass2 := dim.NewQuantity(2, dim.Kilogram)
velocity2 := dim.NewQuantity(10, dim.MeterPerSecond)
velocity2Squared, _ := velocity2.Power(2)
ke, _ := mass2.Multiply(velocity2Squared)
ke = ke.MultiplyScalar(0.5)
fmt.Printf("Kinetic Energy: %v\n", ke)
```

### Comparison Operations

```go
import dim "github.com/parthivrawat/unit-aware-arithmetic/go/v2"

d1 := dim.NewQuantity(100, dim.Centimeter)
d2 := dim.NewQuantity(1, dim.Meter)

// Automatic conversion for comparison
fmt.Println(d1.Equals(d2, 1e-9)) // true
isLess, _ := d1.LessThan(dim.NewQuantity(2, dim.Meter))
fmt.Println(isLess) // true
```

## Available Units

### Length
- `Meter`, `Kilometer`, `Centimeter`, `Millimeter`
- `Inch`, `Foot`, `Yard`, `Mile`

### Mass
- `Kilogram`, `Gram`, `Milligram`, `Tonne`
- `Pound`, `Ounce`

### Time
- `Second`, `Minute`, `Hour`, `Day`

### Temperature
- `Kelvin`, `Celsius`, `Fahrenheit`

### Current
- `Ampere`, `Milliampere`

### Amount of Substance
- `Mole`

### Angle (dimensionless)
- `Radian`, `Degree`, `Arcminute`, `Arcsecond`

### Derived Units
- `Newton` (force)
- `Joule`, `Kilojoule`, `Calorie`, `Kilocalorie`, `WattHour` (energy)
- `Watt`, `Kilowatt` (power)
- `Pascal`, `Kilopascal` (pressure)
- `MeterPerSecond`, `KilometerPerHour`, `MilePerHour` (velocity)
- `MeterPerSecondSquared` (acceleration)
- `SquareMeter`, `SquareKilometer`, `Hectare` (area)
- `CubicMeter`, `Liter`, `Milliliter` (volume)
- `Hertz`, `Kilohertz`, `Megahertz` (frequency)
- `Volt`, `Ohm` (electricity)

## Error Handling

The library provides clear error messages for invalid operations:

```go
import dim "github.com/parthivrawat/unit-aware-arithmetic/go/v2"

distance := dim.NewQuantity(100, dim.Meter)
time := dim.NewQuantity(10, dim.Second)

// This will return an error
result, err := distance.Add(time)
if err != nil {
	fmt.Println(err) // "Cannot add m and s: incompatible dimensions"
}
```

### Affine (Temperature) Units

Multiplication, division, and exponentiation are undefined for affine units —
units with a non-zero offset such as `Celsius` and `Fahrenheit`. `Multiply`,
`Divide`, and `Power` return an `*AffineUnitArithmeticError` when any
involved `Quantity` has `Unit.Offset != 0`. Scalar operations
(`MultiplyScalar`, `DivideScalar`) are unaffected.

```go
import dim "github.com/parthivrawat/unit-aware-arithmetic/go/v2"

// Returns an error: Cannot multiply affine units °C and °C; ...
_, err := dim.NewQuantity(2, dim.Celsius).Multiply(dim.NewQuantity(3, dim.Celsius))

// Convert to an absolute (zero-offset) unit first, e.g., Kelvin
k, _ := dim.NewQuantity(2, dim.Celsius).To(dim.Kelvin)
_, _ = k.Multiply(k) // OK
```

## Testing

Run the test suite:

```bash
go test -v
```

Run with coverage:

```bash
go test -cover
go test -coverprofile=coverage.out
go tool cover -html=coverage.out
```

## API Reference

### `Dimension`

Represents the dimensional formula of a unit.

**Methods:**
- `Multiply(other Dimension) Dimension`
- `Divide(other Dimension) Dimension`
- `Power(exponent int) Dimension`
- `IsDimensionless() bool`
- `Equals(other Dimension) bool`

### `Unit`

Represents a unit of measurement with its dimension and conversion factor.

**Constructor:**
```go
NewUnit(name, symbol string, dimension Dimension, toBase, offset float64) Unit
```

### `Quantity`

A numeric value with an associated unit. `Quantity` is immutable by value; its
fields are unexported and accessed through `Value()` and `Unit()`.

**Constructor:**
```go
NewQuantity(value float64, unit Unit) Quantity
```

**Accessors:**
- `Value() float64`
- `Unit() Unit`

**Arithmetic Methods:**
- `Add(other Quantity) (Quantity, error)`
- `Subtract(other Quantity) (Quantity, error)`
- `Multiply(other Quantity) (Quantity, error)`
- `MultiplyScalar(f float64) Quantity`
- `Divide(other Quantity) (Quantity, error)`
- `DivideScalar(f float64) Quantity`
- `Power(exponent int) (Quantity, error)`
- `Negate() Quantity`
- `Abs() Quantity`

**Comparison Methods:**
- `Equals(other Quantity, tolerance float64) bool` — `true` if the values are equal within an absolute `tolerance`.
- `IsClose(other Quantity, relTol, absTol float64) (bool, error)` — `true` if the values are close within a relative `relTol` or an absolute `absTol` (or both).
- `LessThan(other Quantity) (bool, error)`
- `GreaterThan(other Quantity) (bool, error)`
- `LessThanOrEqual(other Quantity, tolerance float64) (bool, error)`
- `GreaterThanOrEqual(other Quantity, tolerance float64) (bool, error)`

`==` is strict struct equality; use `Equals` or `IsClose` for approximate comparisons.

**Conversion:**
- `To(targetUnit Unit) (Quantity, error)`

**Utility:**
- `String() string`

## Design Principles

1. **Zero Dependencies**: Core functionality has no external dependencies
2. **Type Safety**: Prevents invalid operations at runtime
3. **Clear Errors**: Actionable error messages with context
4. **Production Ready**: Comprehensive test coverage
5. **Performance**: Efficient implementation with minimal overhead
6. **Idiomatic Go**: Follows Go conventions and best practices

## License

MIT License

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- Code follows Go conventions
- Documentation is updated
- Examples are provided

## Changelog

### 2.0.0
- Expanded unit catalog:
  - Angle (dimensionless): `Radian`, `Degree`, `Arcminute`, `Arcsecond`
  - Frequency: `Kilohertz`, `Megahertz` (joins existing `Hertz`)
  - Area: `SquareKilometer`, `Hectare` (joins existing `SquareMeter`)
  - Volume: `Liter`, `Milliliter` (joins existing `CubicMeter`)
  - Velocity: `MilePerHour`
  - Amount of substance: `Mole`
  - Energy: `Calorie`, `Kilocalorie`, `WattHour`
  - Electricity: `Volt`, `Ohm`
- Extended canonical-unit resolution: `Mole`, `Volt`, and `Ohm` are now
  preferred result units for their dimensions (e.g., W/A → V, V/A → Ω)

### 1.0.0 (2026-08-28)
- Initial release
- Support for SI and imperial units
- Comprehensive dimensional analysis
- Temperature conversion with offset handling
- Full test coverage (44 tests)
- Idiomatic Go implementation
