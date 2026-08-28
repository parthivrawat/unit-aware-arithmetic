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
go get github.com/parthivrawat/unit-aware-arithmetic/go
```

## Quick Start

```go
package main

import (
	"fmt"
	dim "github.com/parthivrawat/unit-aware-arithmetic/go"
)

func main() {
	// Create quantities with units
	distance := dim.NewQuantity(100, dim.Meter)
	time := dim.NewQuantity(9.58, dim.Second)

	// Arithmetic operations with automatic unit tracking
	speed := distance.Divide(time)
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
import dim "github.com/parthivrawat/unit-aware-arithmetic/go"

// Addition (same dimension required)
d1 := dim.NewQuantity(5, dim.Meter)
d2 := dim.NewQuantity(3, dim.Meter)
total, _ := d1.Add(d2) // 8.0 m

// Multiplication creates derived units
area := dim.NewQuantity(5, dim.Meter).Multiply(dim.NewQuantity(3, dim.Meter))
fmt.Println(area) // 15 m·m

// Division creates derived units
velocity := dim.NewQuantity(100, dim.Meter).Divide(dim.NewQuantity(10, dim.Second))
fmt.Println(velocity) // 10 m/s
```

### Unit Conversion

```go
import dim "github.com/parthivrawat/unit-aware-arithmetic/go"

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
import dim "github.com/parthivrawat/unit-aware-arithmetic/go"

// Calculate velocity
distance := dim.NewQuantity(100, dim.Meter)
time := dim.NewQuantity(9.58, dim.Second)
velocity := distance.Divide(time)
fmt.Printf("Velocity: %v\n", velocity)

// Calculate force (F = ma)
mass := dim.NewQuantity(10, dim.Kilogram)
acceleration := dim.NewQuantity(9.8, dim.MeterPerSecondSquared)
force := mass.Multiply(acceleration)
fmt.Printf("Force: %v\n", force)

// Calculate kinetic energy (KE = 1/2 * m * v²)
mass2 := dim.NewQuantity(2, dim.Kilogram)
velocity2 := dim.NewQuantity(10, dim.MeterPerSecond)
ke := mass2.Multiply(velocity2.Power(2)).Multiply(0.5)
fmt.Printf("Kinetic Energy: %v\n", ke)
```

### Comparison Operations

```go
import dim "github.com/parthivrawat/unit-aware-arithmetic/go"

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

### Derived Units
- `Newton` (force)
- `Joule`, `Kilojoule` (energy)
- `Watt`, `Kilowatt` (power)
- `Pascal`, `Kilopascal` (pressure)
- `MeterPerSecond`, `KilometerPerHour` (velocity)
- `MeterPerSecondSquared` (acceleration)

## Error Handling

The library provides clear error messages for invalid operations:

```go
import dim "github.com/parthivrawat/unit-aware-arithmetic/go"

distance := dim.NewQuantity(100, dim.Meter)
time := dim.NewQuantity(10, dim.Second)

// This will return an error
result, err := distance.Add(time)
if err != nil {
	fmt.Println(err) // "Cannot add m and s: incompatible dimensions"
}
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

A numeric value with an associated unit.

**Constructor:**
```go
NewQuantity(value float64, unit Unit) Quantity
```

**Arithmetic Methods:**
- `Add(other Quantity) (Quantity, error)`
- `Subtract(other Quantity) (Quantity, error)`
- `Multiply(other interface{}) Quantity` (accepts Quantity, float64, or int)
- `Divide(other interface{}) Quantity` (accepts Quantity, float64, or int)
- `Power(exponent int) Quantity`
- `Negate() Quantity`
- `Abs() Quantity`

**Comparison Methods:**
- `Equals(other Quantity, tolerance float64) bool`
- `LessThan(other Quantity) (bool, error)`
- `GreaterThan(other Quantity) (bool, error)`

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

### 1.0.0 (2026-08-28)
- Initial release
- Support for SI and imperial units
- Comprehensive dimensional analysis
- Temperature conversion with offset handling
- Full test coverage (44 tests)
- Idiomatic Go implementation
