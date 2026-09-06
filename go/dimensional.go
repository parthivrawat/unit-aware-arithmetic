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

// IsCompatibleWith reports whether two dimensions are equal (i.e., compatible).
func (d Dimension) IsCompatibleWith(other Dimension) bool {
	return d.Equals(other)
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

// ToBaseValue converts a value in this unit to the base (zero-offset) unit.
func (u Unit) ToBaseValue(v float64) float64 {
	return (v + u.Offset) * u.ToBase
}

// FromBaseValue converts a value from the base (zero-offset) unit to this unit.
func (u Unit) FromBaseValue(v float64) float64 {
	return v/u.ToBase - u.Offset
}

// IncompatibleUnitsError is returned when attempting incompatible unit operations.
type IncompatibleUnitsError struct {
	Message string
}

func (e *IncompatibleUnitsError) Error() string {
	return e.Message
}

// AffineUnitArithmeticError is raised when attempting arithmetic that is
// undefined for affine (offset) units such as Celsius or Fahrenheit.
type AffineUnitArithmeticError struct {
	Message string
}

func (e *AffineUnitArithmeticError) Error() string {
	return e.Message
}

// Quantity represents a numeric value with an associated unit.
// Its fields are unexported; use Value() and Unit() to access them.
type Quantity struct {
	value float64
	unit  Unit
}

// NewQuantity creates a new quantity.
func NewQuantity(value float64, unit Unit) Quantity {
	return Quantity{value: value, unit: unit}
}

// Value returns the numeric value of the quantity.
func (q Quantity) Value() float64 {
	return q.value
}

// Unit returns the unit of the quantity.
func (q Quantity) Unit() Unit {
	return q.unit
}

// String returns a string representation of the quantity.
func (q Quantity) String() string {
	return fmt.Sprintf("%g %s", q.value, q.unit.Symbol)
}

// ensureCompatible returns an IncompatibleUnitsError if the two quantities' dimensions differ.
func (q Quantity) ensureCompatible(other Quantity, message string) *IncompatibleUnitsError {
	if !q.unit.Dimension.IsCompatibleWith(other.unit.Dimension) {
		return &IncompatibleUnitsError{Message: message}
	}
	return nil
}

// Add adds two quantities (must have compatible dimensions).
func (q Quantity) Add(other Quantity) (Quantity, error) {
	if err := q.ensureCompatible(other, fmt.Sprintf("Cannot add %s and %s: incompatible dimensions",
		q.unit.Symbol, other.unit.Symbol)); err != nil {
		return Quantity{}, err
	}

	otherInSelfUnit, err := other.To(q.unit)
	if err != nil {
		return Quantity{}, err
	}

	return Quantity{value: q.value + otherInSelfUnit.value, unit: q.unit}, nil
}

// Subtract subtracts two quantities (must have compatible dimensions).
func (q Quantity) Subtract(other Quantity) (Quantity, error) {
	if err := q.ensureCompatible(other, fmt.Sprintf("Cannot subtract %s from %s: incompatible dimensions",
		other.unit.Symbol, q.unit.Symbol)); err != nil {
		return Quantity{}, err
	}

	otherInSelfUnit, err := other.To(q.unit)
	if err != nil {
		return Quantity{}, err
	}

	return Quantity{value: q.value - otherInSelfUnit.value, unit: q.unit}, nil
}

// Multiply multiplies a quantity by another quantity, producing a derived unit.
// It returns an *AffineUnitArithmeticError if either unit is affine
// (has a non-zero Offset, e.g., Celsius or Fahrenheit).
func (q Quantity) Multiply(other Quantity) (Quantity, error) {
	if q.unit.Offset != 0 || other.unit.Offset != 0 {
		return Quantity{}, &AffineUnitArithmeticError{
			Message: fmt.Sprintf("Cannot multiply affine units %s and %s; convert to an absolute (zero-offset) unit first.",
				q.unit.Symbol, other.unit.Symbol),
		}
	}
	baseValue := q.unit.ToBaseValue(q.value) * other.unit.ToBaseValue(other.value)
	newDimension := q.unit.Dimension.Multiply(other.unit.Dimension)
	if canon, ok := canonicalUnits[newDimension]; ok {
		return Quantity{value: canon.FromBaseValue(baseValue), unit: canon}, nil
	}
	newSymbol := fmt.Sprintf("%s·%s", q.unit.Symbol, other.unit.Symbol)
	newUnit := NewUnit(
		fmt.Sprintf("%s %s", q.unit.Name, other.unit.Name),
		newSymbol,
		newDimension,
		q.unit.ToBase*other.unit.ToBase,
		0,
	)
	return Quantity{value: newUnit.FromBaseValue(baseValue), unit: newUnit}, nil
}

// MultiplyScalar multiplies a quantity by a scalar factor.
func (q Quantity) MultiplyScalar(f float64) Quantity {
	return Quantity{value: q.value * f, unit: q.unit}
}

// Divide divides a quantity by another quantity, producing a derived unit.
// It returns an *AffineUnitArithmeticError if either unit is affine
// (has a non-zero Offset, e.g., Celsius or Fahrenheit).
func (q Quantity) Divide(other Quantity) (Quantity, error) {
	if q.unit.Offset != 0 || other.unit.Offset != 0 {
		return Quantity{}, &AffineUnitArithmeticError{
			Message: fmt.Sprintf("Cannot divide affine units %s and %s; convert to an absolute (zero-offset) unit first.",
				q.unit.Symbol, other.unit.Symbol),
		}
	}
	baseValue := q.unit.ToBaseValue(q.value) / other.unit.ToBaseValue(other.value)
	newDimension := q.unit.Dimension.Divide(other.unit.Dimension)
	if canon, ok := canonicalUnits[newDimension]; ok {
		return Quantity{value: canon.FromBaseValue(baseValue), unit: canon}, nil
	}
	newSymbol := fmt.Sprintf("%s/%s", q.unit.Symbol, other.unit.Symbol)
	newUnit := NewUnit(
		fmt.Sprintf("%s per %s", q.unit.Name, other.unit.Name),
		newSymbol,
		newDimension,
		q.unit.ToBase/other.unit.ToBase,
		0,
	)
	return Quantity{value: newUnit.FromBaseValue(baseValue), unit: newUnit}, nil
}

// DivideScalar divides a quantity by a scalar factor.
func (q Quantity) DivideScalar(f float64) Quantity {
	return Quantity{value: q.value / f, unit: q.unit}
}

// Power raises a quantity to an integer power.
// It returns an *AffineUnitArithmeticError if the unit is affine
// (has a non-zero Offset, e.g., Celsius or Fahrenheit).
func (q Quantity) Power(exponent int) (Quantity, error) {
	if q.unit.Offset != 0 {
		return Quantity{}, &AffineUnitArithmeticError{
			Message: fmt.Sprintf("Cannot raise affine unit %s to a power; convert to an absolute (zero-offset) unit first.",
				q.unit.Symbol),
		}
	}
	baseValue := math.Pow(q.unit.ToBaseValue(q.value), float64(exponent))
	newDimension := q.unit.Dimension.Power(exponent)
	if canon, ok := canonicalUnits[newDimension]; ok {
		return Quantity{value: canon.FromBaseValue(baseValue), unit: canon}, nil
	}
	newSymbol := fmt.Sprintf("%s^%d", q.unit.Symbol, exponent)
	newUnit := NewUnit(
		fmt.Sprintf("%s to the power %d", q.unit.Name, exponent),
		newSymbol,
		newDimension,
		math.Pow(q.unit.ToBase, float64(exponent)),
		0,
	)
	return Quantity{value: newUnit.FromBaseValue(baseValue), unit: newUnit}, nil
}

// Negate negates a quantity.
func (q Quantity) Negate() Quantity {
	return Quantity{value: -q.value, unit: q.unit}
}

// Abs returns the absolute value.
func (q Quantity) Abs() Quantity {
	return Quantity{value: math.Abs(q.value), unit: q.unit}
}

// Equals checks if two quantities are equal within an absolute tolerance.
func (q Quantity) Equals(other Quantity, tolerance float64) bool {
	if !q.unit.Dimension.IsCompatibleWith(other.unit.Dimension) {
		return false
	}

	otherInSelfUnit, err := other.To(q.unit)
	if err != nil {
		return false
	}

	return math.Abs(q.value-otherInSelfUnit.value) < tolerance
}

// IsClose checks if two quantities are close within a relative or absolute tolerance.
func (q Quantity) IsClose(other Quantity, relTol, absTol float64) (bool, error) {
	if err := q.ensureCompatible(other, fmt.Sprintf("Cannot compare %s and %s", q.unit.Symbol, other.unit.Symbol)); err != nil {
		return false, err
	}

	otherInSelfUnit, err := other.To(q.unit)
	if err != nil {
		return false, err
	}

	if q.value == otherInSelfUnit.value {
		return true, nil
	}

	if math.IsNaN(q.value) || math.IsNaN(otherInSelfUnit.value) ||
		math.IsInf(q.value, 0) || math.IsInf(otherInSelfUnit.value, 0) {
		return false, nil
	}

	diff := math.Abs(q.value - otherInSelfUnit.value)
	max := math.Max(math.Abs(q.value), math.Abs(otherInSelfUnit.value))
	return diff <= absTol || diff <= relTol*max, nil
}

// LessThan checks if this quantity is less than another.
func (q Quantity) LessThan(other Quantity) (bool, error) {
	if err := q.ensureCompatible(other, fmt.Sprintf("Cannot compare %s and %s", q.unit.Symbol, other.unit.Symbol)); err != nil {
		return false, err
	}

	otherInSelfUnit, err := other.To(q.unit)
	if err != nil {
		return false, err
	}

	return q.value < otherInSelfUnit.value, nil
}

// GreaterThan checks if this quantity is greater than another.
func (q Quantity) GreaterThan(other Quantity) (bool, error) {
	if err := q.ensureCompatible(other, fmt.Sprintf("Cannot compare %s and %s", q.unit.Symbol, other.unit.Symbol)); err != nil {
		return false, err
	}

	otherInSelfUnit, err := other.To(q.unit)
	if err != nil {
		return false, err
	}

	return q.value > otherInSelfUnit.value, nil
}

// LessThanOrEqual checks if this quantity is less than or equal to another (with tolerance).
func (q Quantity) LessThanOrEqual(other Quantity, tolerance float64) (bool, error) {
	less, err := q.LessThan(other)
	if err != nil {
		return false, err
	}

	isClose, err := q.IsClose(other, 0.0, tolerance)
	if err != nil {
		return false, err
	}

	return less || isClose, nil
}

// GreaterThanOrEqual checks if this quantity is greater than or equal to another (with tolerance).
func (q Quantity) GreaterThanOrEqual(other Quantity, tolerance float64) (bool, error) {
	greater, err := q.GreaterThan(other)
	if err != nil {
		return false, err
	}

	isClose, err := q.IsClose(other, 0.0, tolerance)
	if err != nil {
		return false, err
	}

	return greater || isClose, nil
}

// To converts to another unit (must have compatible dimensions).
func (q Quantity) To(targetUnit Unit) (Quantity, error) {
	if !q.unit.Dimension.IsCompatibleWith(targetUnit.Dimension) {
		return Quantity{}, &IncompatibleUnitsError{
			Message: fmt.Sprintf("Cannot convert %s to %s: incompatible dimensions",
				q.unit.Symbol, targetUnit.Symbol),
		}
	}

	baseValue := q.unit.ToBaseValue(q.value)
	return Quantity{value: targetUnit.FromBaseValue(baseValue), unit: targetUnit}, nil
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
	Mile = NewUnit("mile", "mi", Dimension{Length: 1}, 1609.344, 0.0)

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
	MeterPerSecond   = NewUnit("meter per second", "m/s", Dimension{Length: 1, Time: -1}, 1.0, 0.0)
	KilometerPerHour = NewUnit("kilometer per hour", "km/h", Dimension{Length: 1, Time: -1}, 1000.0/3600.0, 0.0)

	// Acceleration (m/s²)
	MeterPerSecondSquared = NewUnit("meter per second squared", "m/s²", Dimension{Length: 1, Time: -2}, 1.0, 0.0)

	// Area (m²)
	SquareMeter = NewUnit("square meter", "m²", Dimension{Length: 2}, 1.0, 0.0)

	// Volume (m³)
	CubicMeter = NewUnit("cubic meter", "m³", Dimension{Length: 3}, 1.0, 0.0)

	// Frequency (1/s)
	Hertz     = NewUnit("hertz", "Hz", Dimension{Time: -1}, 1.0, 0.0)
	Kilohertz = NewUnit("kilohertz", "kHz", Dimension{Time: -1}, 1e3, 0.0)
	Megahertz = NewUnit("megahertz", "MHz", Dimension{Time: -1}, 1e6, 0.0)

	// Angle (dimensionless in SI)
	Radian    = NewUnit("radian", "rad", Dimension{}, 1.0, 0.0)
	Degree    = NewUnit("degree", "°", Dimension{}, math.Pi/180.0, 0.0)
	Arcminute = NewUnit("arcminute", "′", Dimension{}, math.Pi/10800.0, 0.0)
	Arcsecond = NewUnit("arcsecond", "″", Dimension{}, math.Pi/648000.0, 0.0)

	// Area
	SquareKilometer = NewUnit("square kilometer", "km²", Dimension{Length: 2}, 1e6, 0.0)
	Hectare         = NewUnit("hectare", "ha", Dimension{Length: 2}, 1e4, 0.0)

	// Volume
	Liter      = NewUnit("liter", "L", Dimension{Length: 3}, 1e-3, 0.0)
	Milliliter = NewUnit("milliliter", "mL", Dimension{Length: 3}, 1e-6, 0.0)

	// Velocity
	MilePerHour = NewUnit("mile per hour", "mph", Dimension{Length: 1, Time: -1}, 0.44704, 0.0)

	// Amount of substance
	Mole = NewUnit("mole", "mol", Dimension{Amount: 1}, 1.0, 0.0)

	// Energy
	Calorie     = NewUnit("calorie", "cal", Dimension{Length: 2, Mass: 1, Time: -2}, 4.184, 0.0)
	Kilocalorie = NewUnit("kilocalorie", "kcal", Dimension{Length: 2, Mass: 1, Time: -2}, 4184.0, 0.0)
	WattHour    = NewUnit("watt hour", "Wh", Dimension{Length: 2, Mass: 1, Time: -2}, 3600.0, 0.0)

	// Electricity
	Volt = NewUnit("volt", "V", Dimension{Length: 2, Mass: 1, Time: -3, Current: -1}, 1.0, 0.0)
	Ohm  = NewUnit("ohm", "Ω", Dimension{Length: 2, Mass: 1, Time: -3, Current: -2}, 1.0, 0.0)
)

// canonicalUnits maps a Dimension to its canonical SI/base unit, if one is
// known. Results of multiplication, division, and power operations whose
// dimension has an entry here are expressed in that unit.
var canonicalUnits map[Dimension]Unit

func init() {
	canonicalUnits = make(map[Dimension]Unit)
	for _, u := range []Unit{
		Dimensionless, Meter, Kilogram, Second, Kelvin, Ampere,
		Newton, Joule, Watt, Pascal, MeterPerSecond, MeterPerSecondSquared,
		SquareMeter, CubicMeter, Hertz, Mole, Volt, Ohm,
	} {
		canonicalUnits[u.Dimension] = u
	}
}
