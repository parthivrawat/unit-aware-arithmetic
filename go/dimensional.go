// Package dimensional provides type-safe dimensional arithmetic with unit tracking.
//
// This library prevents unit-conversion bugs through dimensional analysis,
// tracking physical dimensions (length, mass, time, etc.) and preventing
// incompatible operations like adding meters to seconds.
package dimensional

import (
	"fmt"
	"math"
)

// Dimension represents the dimensional formula of a unit (e.g., L^1 T^-2 for acceleration).
type Dimension struct {
	Length      int // L
	Mass        int // M
	Time        int // T
	Current     int // I
	Temperature int // Θ
	Amount      int // N
	Luminosity  int // J
}

// Multiply multiplies two dimensions (adds exponents).
func (d Dimension) Multiply(other Dimension) Dimension {
	return Dimension{
		Length:      d.Length + other.Length,
		Mass:        d.Mass + other.Mass,
		Time:        d.Time + other.Time,
		Current:     d.Current + other.Current,
		Temperature: d.Temperature + other.Temperature,
		Amount:      d.Amount + other.Amount,
		Luminosity:  d.Luminosity + other.Luminosity,
	}
}

// Divide divides two dimensions (subtracts exponents).
func (d Dimension) Divide(other Dimension) Dimension {
	return Dimension{
		Length:      d.Length - other.Length,
		Mass:        d.Mass - other.Mass,
		Time:        d.Time - other.Time,
		Current:     d.Current - other.Current,
		Temperature: d.Temperature - other.Temperature,
		Amount:      d.Amount - other.Amount,
		Luminosity:  d.Luminosity - other.Luminosity,
	}
}

// Power raises a dimension to a power.
func (d Dimension) Power(exponent int) Dimension {
	return Dimension{
		Length:      d.Length * exponent,
		Mass:        d.Mass * exponent,
		Time:        d.Time * exponent,
		Current:     d.Current * exponent,
		Temperature: d.Temperature * exponent,
		Amount:      d.Amount * exponent,
		Luminosity:  d.Luminosity * exponent,
	}
}

// IsDimensionless checks if this is a dimensionless quantity.
func (d Dimension) IsDimensionless() bool {
	return d.Length == 0 && d.Mass == 0 && d.Time == 0 &&
		d.Current == 0 && d.Temperature == 0 &&
		d.Amount == 0 && d.Luminosity == 0
}

// Equals checks if two dimensions are equal.
func (d Dimension) Equals(other Dimension) bool {
	return d.Length == other.Length &&
		d.Mass == other.Mass &&
		d.Time == other.Time &&
		d.Current == other.Current &&
		d.Temperature == other.Temperature &&
		d.Amount == other.Amount &&
		d.Luminosity == other.Luminosity
}

// Unit represents a unit of measurement with its dimension and conversion factor.
type Unit struct {
	Name      string
	Symbol    string
	Dimension Dimension
	ToBase    float64 // Conversion factor to base unit
	Offset    float64 // Offset for affine conversions (e.g., temperature)
}

// NewUnit creates a new unit.
func NewUnit(name, symbol string, dimension Dimension, toBase, offset float64) Unit {
	return Unit{
		Name:      name,
		Symbol:    symbol,
		Dimension: dimension,
		ToBase:    toBase,
		Offset:    offset,
	}
}

// IncompatibleUnitsError is returned when attempting incompatible unit operations.
type IncompatibleUnitsError struct {
	Message string
}

func (e *IncompatibleUnitsError) Error() string {
	return e.Message
}

// Quantity represents a numeric value with an associated unit.
type Quantity struct {
	Value float64
	Unit  Unit
}

// NewQuantity creates a new quantity.
func NewQuantity(value float64, unit Unit) Quantity {
	return Quantity{Value: value, Unit: unit}
}

// String returns a string representation of the quantity.
func (q Quantity) String() string {
	return fmt.Sprintf("%g %s", q.Value, q.Unit.Symbol)
}

// Add adds two quantities (must have compatible dimensions).
func (q Quantity) Add(other Quantity) (Quantity, error) {
	if !q.Unit.Dimension.Equals(other.Unit.Dimension) {
		return Quantity{}, &IncompatibleUnitsError{
			Message: fmt.Sprintf("Cannot add %s and %s: incompatible dimensions",
				q.Unit.Symbol, other.Unit.Symbol),
		}
	}

	otherInSelfUnit, err := other.To(q.Unit)
	if err != nil {
		return Quantity{}, err
	}

	return Quantity{Value: q.Value + otherInSelfUnit.Value, Unit: q.Unit}, nil
}

// Subtract subtracts two quantities (must have compatible dimensions).
func (q Quantity) Subtract(other Quantity) (Quantity, error) {
	if !q.Unit.Dimension.Equals(other.Unit.Dimension) {
		return Quantity{}, &IncompatibleUnitsError{
			Message: fmt.Sprintf("Cannot subtract %s from %s: incompatible dimensions",
				other.Unit.Symbol, q.Unit.Symbol),
		}
	}

	otherInSelfUnit, err := other.To(q.Unit)
	if err != nil {
		return Quantity{}, err
	}

	return Quantity{Value: q.Value - otherInSelfUnit.Value, Unit: q.Unit}, nil
}

// Multiply multiplies a quantity by another quantity or scalar.
func (q Quantity) Multiply(other interface{}) Quantity {
	switch v := other.(type) {
	case float64:
		return Quantity{Value: q.Value * v, Unit: q.Unit}
	case int:
		return Quantity{Value: q.Value * float64(v), Unit: q.Unit}
	case Quantity:
		newValue := q.Value * v.Value
		newDimension := q.Unit.Dimension.Multiply(v.Unit.Dimension)
		newSymbol := fmt.Sprintf("%s·%s", q.Unit.Symbol, v.Unit.Symbol)
		newUnit := NewUnit(
			fmt.Sprintf("%s %s", q.Unit.Name, v.Unit.Name),
			newSymbol,
			newDimension,
			q.Unit.ToBase*v.Unit.ToBase,
			0,
		)
		return Quantity{Value: newValue, Unit: newUnit}
	default:
		panic(fmt.Sprintf("Cannot multiply Quantity by %T", other))
	}
}

// Divide divides a quantity by another quantity or scalar.
func (q Quantity) Divide(other interface{}) Quantity {
	switch v := other.(type) {
	case float64:
		return Quantity{Value: q.Value / v, Unit: q.Unit}
	case int:
		return Quantity{Value: q.Value / float64(v), Unit: q.Unit}
	case Quantity:
		newValue := q.Value / v.Value
		newDimension := q.Unit.Dimension.Divide(v.Unit.Dimension)
		newSymbol := fmt.Sprintf("%s/%s", q.Unit.Symbol, v.Unit.Symbol)
		newUnit := NewUnit(
			fmt.Sprintf("%s per %s", q.Unit.Name, v.Unit.Name),
			newSymbol,
			newDimension,
			q.Unit.ToBase/v.Unit.ToBase,
			0,
		)
		return Quantity{Value: newValue, Unit: newUnit}
	default:
		panic(fmt.Sprintf("Cannot divide Quantity by %T", other))
	}
}

// Power raises a quantity to an integer power.
func (q Quantity) Power(exponent int) Quantity {
	newValue := math.Pow(q.Value, float64(exponent))
	newDimension := q.Unit.Dimension.Power(exponent)
	newSymbol := fmt.Sprintf("%s^%d", q.Unit.Symbol, exponent)
	newUnit := NewUnit(
		fmt.Sprintf("%s to the power %d", q.Unit.Name, exponent),
		newSymbol,
		newDimension,
		math.Pow(q.Unit.ToBase, float64(exponent)),
		0,
	)
	return Quantity{Value: newValue, Unit: newUnit}
}

// Negate negates a quantity.
func (q Quantity) Negate() Quantity {
	return Quantity{Value: -q.Value, Unit: q.Unit}
}

// Abs returns the absolute value.
func (q Quantity) Abs() Quantity {
	return Quantity{Value: math.Abs(q.Value), Unit: q.Unit}
}

// Equals checks if two quantities are equal (with tolerance).
func (q Quantity) Equals(other Quantity, tolerance float64) bool {
	if !q.Unit.Dimension.Equals(other.Unit.Dimension) {
		return false
	}

	otherInSelfUnit, err := other.To(q.Unit)
	if err != nil {
		return false
	}

	return math.Abs(q.Value-otherInSelfUnit.Value) < tolerance
}

// LessThan checks if this quantity is less than another.
func (q Quantity) LessThan(other Quantity) (bool, error) {
	if !q.Unit.Dimension.Equals(other.Unit.Dimension) {
		return false, &IncompatibleUnitsError{
			Message: fmt.Sprintf("Cannot compare %s and %s", q.Unit.Symbol, other.Unit.Symbol),
		}
	}

	otherInSelfUnit, err := other.To(q.Unit)
	if err != nil {
		return false, err
	}

	return q.Value < otherInSelfUnit.Value, nil
}

// GreaterThan checks if this quantity is greater than another.
func (q Quantity) GreaterThan(other Quantity) (bool, error) {
	if !q.Unit.Dimension.Equals(other.Unit.Dimension) {
		return false, &IncompatibleUnitsError{
			Message: fmt.Sprintf("Cannot compare %s and %s", q.Unit.Symbol, other.Unit.Symbol),
		}
	}

	otherInSelfUnit, err := other.To(q.Unit)
	if err != nil {
		return false, err
	}

	return q.Value > otherInSelfUnit.Value, nil
}

// To converts to another unit (must have compatible dimensions).
func (q Quantity) To(targetUnit Unit) (Quantity, error) {
	if !q.Unit.Dimension.Equals(targetUnit.Dimension) {
		return Quantity{}, &IncompatibleUnitsError{
			Message: fmt.Sprintf("Cannot convert %s to %s: incompatible dimensions",
				q.Unit.Symbol, targetUnit.Symbol),
		}
	}

	var newValue float64
	// Handle affine conversions (e.g., temperature)
	if q.Unit.Offset != 0 || targetUnit.Offset != 0 {
		// Convert to base unit first (remove offset)
		baseValue := (q.Value + q.Unit.Offset) * q.Unit.ToBase
		// Convert from base to target (apply offset)
		newValue = baseValue/targetUnit.ToBase - targetUnit.Offset
	} else {
		// Simple linear conversion
		newValue = q.Value * (q.Unit.ToBase / targetUnit.ToBase)
	}

	return Quantity{Value: newValue, Unit: targetUnit}, nil
}

// Predefined units
var (
	// Dimensionless
	Dimensionless = NewUnit("dimensionless", "", Dimension{}, 1.0, 0.0)

	// Length
	Meter      = NewUnit("meter", "m", Dimension{Length: 1}, 1.0, 0.0)
	Kilometer  = NewUnit("kilometer", "km", Dimension{Length: 1}, 1000.0, 0.0)
	Centimeter = NewUnit("centimeter", "cm", Dimension{Length: 1}, 0.01, 0.0)
	Millimeter = NewUnit("millimeter", "mm", Dimension{Length: 1}, 0.001, 0.0)

	// Imperial length
	Inch = NewUnit("inch", "in", Dimension{Length: 1}, 0.0254, 0.0)
	Foot = NewUnit("foot", "ft", Dimension{Length: 1}, 0.3048, 0.0)
	Yard = NewUnit("yard", "yd", Dimension{Length: 1}, 0.9144, 0.0)
	Mile = NewUnit("mile", "mi", Dimension{Length: 1}, 1609.34, 0.0)

	// Mass
	Kilogram  = NewUnit("kilogram", "kg", Dimension{Mass: 1}, 1.0, 0.0)
	Gram      = NewUnit("gram", "g", Dimension{Mass: 1}, 0.001, 0.0)
	Milligram = NewUnit("milligram", "mg", Dimension{Mass: 1}, 1e-6, 0.0)
	Tonne     = NewUnit("tonne", "t", Dimension{Mass: 1}, 1000.0, 0.0)

	// Imperial mass
	Pound = NewUnit("pound", "lb", Dimension{Mass: 1}, 0.453592, 0.0)
	Ounce = NewUnit("ounce", "oz", Dimension{Mass: 1}, 0.0283495, 0.0)

	// Time
	Second = NewUnit("second", "s", Dimension{Time: 1}, 1.0, 0.0)
	Minute = NewUnit("minute", "min", Dimension{Time: 1}, 60.0, 0.0)
	Hour   = NewUnit("hour", "h", Dimension{Time: 1}, 3600.0, 0.0)
	Day    = NewUnit("day", "d", Dimension{Time: 1}, 86400.0, 0.0)

	// Temperature
	Kelvin     = NewUnit("kelvin", "K", Dimension{Temperature: 1}, 1.0, 0.0)
	Celsius    = NewUnit("celsius", "°C", Dimension{Temperature: 1}, 1.0, 273.15)
	Fahrenheit = NewUnit("fahrenheit", "°F", Dimension{Temperature: 1}, 5.0/9.0, 459.67)

	// Current
	Ampere      = NewUnit("ampere", "A", Dimension{Current: 1}, 1.0, 0.0)
	Milliampere = NewUnit("milliampere", "mA", Dimension{Current: 1}, 0.001, 0.0)

	// Derived units

	// Force (kg·m/s²)
	Newton = NewUnit("newton", "N", Dimension{Length: 1, Mass: 1, Time: -2}, 1.0, 0.0)

	// Energy (kg·m²/s²)
	Joule     = NewUnit("joule", "J", Dimension{Length: 2, Mass: 1, Time: -2}, 1.0, 0.0)
	Kilojoule = NewUnit("kilojoule", "kJ", Dimension{Length: 2, Mass: 1, Time: -2}, 1000.0, 0.0)

	// Power (kg·m²/s³)
	Watt     = NewUnit("watt", "W", Dimension{Length: 2, Mass: 1, Time: -3}, 1.0, 0.0)
	Kilowatt = NewUnit("kilowatt", "kW", Dimension{Length: 2, Mass: 1, Time: -3}, 1000.0, 0.0)

	// Pressure (kg/(m·s²))
	Pascal     = NewUnit("pascal", "Pa", Dimension{Length: -1, Mass: 1, Time: -2}, 1.0, 0.0)
	Kilopascal = NewUnit("kilopascal", "kPa", Dimension{Length: -1, Mass: 1, Time: -2}, 1000.0, 0.0)

	// Velocity (m/s)
	MeterPerSecond    = NewUnit("meter per second", "m/s", Dimension{Length: 1, Time: -1}, 1.0, 0.0)
	KilometerPerHour  = NewUnit("kilometer per hour", "km/h", Dimension{Length: 1, Time: -1}, 1000.0/3600.0, 0.0)

	// Acceleration (m/s²)
	MeterPerSecondSquared = NewUnit("meter per second squared", "m/s²", Dimension{Length: 1, Time: -2}, 1.0, 0.0)
)
