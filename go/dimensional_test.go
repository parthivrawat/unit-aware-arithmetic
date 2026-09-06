package dimensional

import (
	"errors"
	"math"
	"strings"
	"testing"
)

const tolerance = 1e-9

func TestDimension(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		d := Dimension{Length: 1, Time: -2}
		if d.Length != 1 || d.Time != -2 || d.Mass != 0 {
			t.Errorf("Dimension creation failed")
		}
	})

	t.Run("Multiply", func(t *testing.T) {
		d1 := Dimension{Length: 1}
		d2 := Dimension{Time: -1}
		result := d1.Multiply(d2)
		if result.Length != 1 || result.Time != -1 {
			t.Errorf("Dimension multiplication failed")
		}
	})

	t.Run("Divide", func(t *testing.T) {
		d1 := Dimension{Length: 1, Time: -1}
		d2 := Dimension{Time: 1}
		result := d1.Divide(d2)
		if result.Length != 1 || result.Time != -2 {
			t.Errorf("Dimension division failed")
		}
	})

	t.Run("Power", func(t *testing.T) {
		d := Dimension{Length: 1, Time: -1}
		result := d.Power(2)
		if result.Length != 2 || result.Time != -2 {
			t.Errorf("Dimension power failed")
		}
	})

	t.Run("IsDimensionless", func(t *testing.T) {
		d := Dimension{}
		if !d.IsDimensionless() {
			t.Errorf("Empty dimension should be dimensionless")
		}
		d2 := Dimension{Length: 1}
		if d2.IsDimensionless() {
			t.Errorf("Dimension with length should not be dimensionless")
		}
	})
}

func TestUnit(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		u := NewUnit("meter", "m", Dimension{Length: 1}, 1.0, 0.0)
		if u.Name != "meter" || u.Symbol != "m" {
			t.Errorf("Unit creation failed")
		}
	})

	t.Run("PredefinedUnits", func(t *testing.T) {
		if Meter.Symbol != "m" {
			t.Errorf("Meter unit not defined correctly")
		}
		if Kilogram.Symbol != "kg" {
			t.Errorf("Kilogram unit not defined correctly")
		}
		if Second.Symbol != "s" {
			t.Errorf("Second unit not defined correctly")
		}
	})
}

func TestQuantityBasics(t *testing.T) {
	t.Run("Creation", func(t *testing.T) {
		q := NewQuantity(5.0, Meter)
		if q.Value() != 5.0 || q.Unit().Symbol != "m" {
			t.Errorf("Quantity creation failed")
		}
	})

	t.Run("String", func(t *testing.T) {
		q := NewQuantity(5.0, Meter)
		if q.String() != "5 m" {
			t.Errorf("Quantity string representation failed: got %s", q.String())
		}
	})
}

func TestArithmetic(t *testing.T) {
	t.Run("AddSameUnit", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Meter)
		result, err := q1.Add(q2)
		if err != nil {
			t.Errorf("Add failed: %v", err)
		}
		if result.Value() != 8 {
			t.Errorf("Add result incorrect: got %f, want 8", result.Value())
		}
	})

	t.Run("AddCompatibleUnits", func(t *testing.T) {
		q1 := NewQuantity(1, Meter)
		q2 := NewQuantity(100, Centimeter)
		result, err := q1.Add(q2)
		if err != nil {
			t.Errorf("Add failed: %v", err)
		}
		if math.Abs(result.Value()-2.0) > tolerance {
			t.Errorf("Add result incorrect: got %f, want 2.0", result.Value())
		}
	})

	t.Run("AddIncompatibleUnits", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Second)
		_, err := q1.Add(q2)
		if err == nil {
			t.Errorf("Add should fail for incompatible units")
		} else if !strings.Contains(err.Error(), "Cannot add") {
			t.Errorf("Add error message incorrect: got %q", err.Error())
		}
	})

	t.Run("SubtractSameUnit", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Meter)
		result, err := q1.Subtract(q2)
		if err != nil {
			t.Errorf("Subtract failed: %v", err)
		}
		if result.Value() != 2 {
			t.Errorf("Subtract result incorrect: got %f, want 2", result.Value())
		}
	})

	t.Run("SubtractIncompatibleUnits", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Second)
		_, err := q1.Subtract(q2)
		if err == nil {
			t.Errorf("Subtract should fail for incompatible units")
		} else if !strings.Contains(err.Error(), "Cannot subtract") {
			t.Errorf("Subtract error message incorrect: got %q", err.Error())
		}
	})

	t.Run("MultiplyByScalar", func(t *testing.T) {
		q := NewQuantity(5, Meter)
		result := q.MultiplyScalar(3.0)
		if result.Value() != 15 {
			t.Errorf("MultiplyScalar failed: got %f, want 15", result.Value())
		}
		if result.Unit().Symbol != "m" {
			t.Errorf("MultiplyScalar should keep the unit: got %s", result.Unit().Symbol)
		}
	})

	t.Run("MultiplyQuantities", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Meter)
		result := must(q1.Multiply(q2))
		if result.Value() != 15 {
			t.Errorf("Multiply quantities failed: got %f, want 15", result.Value())
		}
		if result.Unit().Dimension.Length != 2 {
			t.Errorf("Multiply dimension incorrect: got %d, want 2", result.Unit().Dimension.Length)
		}
	})

	t.Run("MultiplyDerivedDimension", func(t *testing.T) {
		mass := NewQuantity(2, Kilogram)
		velocity := NewQuantity(10, MeterPerSecond)
		momentum := must(mass.Multiply(velocity))
		if momentum.Unit().Dimension.Mass != 1 ||
			momentum.Unit().Dimension.Length != 1 ||
			momentum.Unit().Dimension.Time != -1 {
			t.Errorf("Multiply derived dimension incorrect")
		}
		if momentum.Value() != 20 {
			t.Errorf("Multiply derived value incorrect: got %f, want 20", momentum.Value())
		}
	})

	t.Run("DivideByScalar", func(t *testing.T) {
		q := NewQuantity(10, Meter)
		result := q.DivideScalar(2.0)
		if result.Value() != 5 {
			t.Errorf("DivideScalar failed: got %f, want 5", result.Value())
		}
		if result.Unit().Symbol != "m" {
			t.Errorf("DivideScalar should keep the unit: got %s", result.Unit().Symbol)
		}
	})

	t.Run("DivideQuantities", func(t *testing.T) {
		q1 := NewQuantity(100, Meter)
		q2 := NewQuantity(10, Second)
		result := must(q1.Divide(q2))
		if result.Value() != 10 {
			t.Errorf("Divide quantities failed: got %f, want 10", result.Value())
		}
		if result.Unit().Dimension.Length != 1 || result.Unit().Dimension.Time != -1 {
			t.Errorf("Divide dimension incorrect")
		}
	})

	t.Run("DivideScalarByZero", func(t *testing.T) {
		q := NewQuantity(10, Meter)
		result := q.DivideScalar(0.0)
		if !math.IsInf(result.Value(), 1) {
			t.Errorf("DivideScalar by zero should produce +Inf, got %f", result.Value())
		}
	})

	t.Run("DivideByZeroQuantity", func(t *testing.T) {
		q := NewQuantity(10, Meter)
		result := must(q.Divide(NewQuantity(0.0, Second)))
		if !math.IsInf(result.Value(), 1) {
			t.Errorf("Divide by zero quantity should produce +Inf, got %f", result.Value())
		}
		if result.Unit().Dimension.Length != 1 || result.Unit().Dimension.Time != -1 {
			t.Errorf("Divide by zero quantity dimension incorrect")
		}
	})

	t.Run("Power", func(t *testing.T) {
		q := NewQuantity(3, Meter)
		result := must(q.Power(2))
		if result.Value() != 9 {
			t.Errorf("Power failed: got %f, want 9", result.Value())
		}
		if result.Unit().Dimension.Length != 2 {
			t.Errorf("Power dimension incorrect: got %d, want 2", result.Unit().Dimension.Length)
		}
	})

	t.Run("Negate", func(t *testing.T) {
		q := NewQuantity(5, Meter)
		result := q.Negate()
		if result.Value() != -5 {
			t.Errorf("Negate failed: got %f, want -5", result.Value())
		}
	})

	t.Run("Abs", func(t *testing.T) {
		q := NewQuantity(-5, Meter)
		result := q.Abs()
		if result.Value() != 5 {
			t.Errorf("Abs failed: got %f, want 5", result.Value())
		}
	})
}

func TestComparison(t *testing.T) {
	t.Run("EqualsSameUnit", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(5, Meter)
		if !q1.Equals(q2, tolerance) {
			t.Errorf("Equals failed for same values")
		}
	})

	t.Run("EqualsCompatibleUnits", func(t *testing.T) {
		q1 := NewQuantity(1, Meter)
		q2 := NewQuantity(100, Centimeter)
		if !q1.Equals(q2, tolerance) {
			t.Errorf("Equals failed for compatible units")
		}
	})

	t.Run("LessThan", func(t *testing.T) {
		q1 := NewQuantity(3, Meter)
		q2 := NewQuantity(5, Meter)
		result, err := q1.LessThan(q2)
		if err != nil {
			t.Errorf("LessThan failed: %v", err)
		}
		if !result {
			t.Errorf("LessThan should return true")
		}
	})

	t.Run("GreaterThan", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Meter)
		result, err := q1.GreaterThan(q2)
		if err != nil {
			t.Errorf("GreaterThan failed: %v", err)
		}
		if !result {
			t.Errorf("GreaterThan should return true")
		}
	})

	t.Run("LessThanOrEqual", func(t *testing.T) {
		q1 := NewQuantity(3, Meter)
		q2 := NewQuantity(5, Meter)
		result, err := q1.LessThanOrEqual(q2, tolerance)
		if err != nil {
			t.Errorf("LessThanOrEqual failed: %v", err)
		}
		if !result {
			t.Errorf("LessThanOrEqual should return true for 3 m <= 5 m")
		}

		equal, err := q1.LessThanOrEqual(NewQuantity(3, Meter), tolerance)
		if err != nil {
			t.Errorf("LessThanOrEqual failed: %v", err)
		}
		if !equal {
			t.Errorf("LessThanOrEqual should return true for equal values")
		}

		greater, err := NewQuantity(10, Meter).LessThanOrEqual(q2, tolerance)
		if err != nil {
			t.Errorf("LessThanOrEqual failed: %v", err)
		}
		if greater {
			t.Errorf("LessThanOrEqual should return false for 10 m <= 5 m")
		}
	})

	t.Run("GreaterThanOrEqual", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Meter)
		result, err := q1.GreaterThanOrEqual(q2, tolerance)
		if err != nil {
			t.Errorf("GreaterThanOrEqual failed: %v", err)
		}
		if !result {
			t.Errorf("GreaterThanOrEqual should return true for 5 m >= 3 m")
		}

		equal, err := q1.GreaterThanOrEqual(NewQuantity(5, Meter), tolerance)
		if err != nil {
			t.Errorf("GreaterThanOrEqual failed: %v", err)
		}
		if !equal {
			t.Errorf("GreaterThanOrEqual should return true for equal values")
		}

		less, err := NewQuantity(1, Meter).GreaterThanOrEqual(q2, tolerance)
		if err != nil {
			t.Errorf("GreaterThanOrEqual failed: %v", err)
		}
		if less {
			t.Errorf("GreaterThanOrEqual should return false for 1 m >= 3 m")
		}
	})

	t.Run("CompareIncompatibleUnits", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Second)
		_, err := q1.LessThan(q2)
		if err == nil {
			t.Errorf("Comparison should fail for incompatible units")
		}
		_, err = q1.LessThanOrEqual(q2, tolerance)
		if err == nil {
			t.Errorf("LessThanOrEqual should fail for incompatible units")
		}
		_, err = q1.GreaterThanOrEqual(q2, tolerance)
		if err == nil {
			t.Errorf("GreaterThanOrEqual should fail for incompatible units")
		}
	})
}

func TestConversion(t *testing.T) {
	t.Run("ConvertLength", func(t *testing.T) {
		q := NewQuantity(1, Meter)
		result, err := q.To(Centimeter)
		if err != nil {
			t.Errorf("Conversion failed: %v", err)
		}
		if math.Abs(result.Value()-100) > tolerance {
			t.Errorf("Conversion incorrect: got %f, want 100", result.Value())
		}
	})

	t.Run("ConvertMass", func(t *testing.T) {
		q := NewQuantity(1, Kilogram)
		result, err := q.To(Gram)
		if err != nil {
			t.Errorf("Conversion failed: %v", err)
		}
		if math.Abs(result.Value()-1000) > tolerance {
			t.Errorf("Conversion incorrect: got %f, want 1000", result.Value())
		}
	})

	t.Run("ConvertTime", func(t *testing.T) {
		q := NewQuantity(1, Hour)
		result, err := q.To(Second)
		if err != nil {
			t.Errorf("Conversion failed: %v", err)
		}
		if math.Abs(result.Value()-3600) > tolerance {
			t.Errorf("Conversion incorrect: got %f, want 3600", result.Value())
		}
	})

	t.Run("ConvertImperialToMetric", func(t *testing.T) {
		q := NewQuantity(1, Mile)
		result, err := q.To(Kilometer)
		if err != nil {
			t.Errorf("Conversion failed: %v", err)
		}
		if math.Abs(result.Value()-1.60934) > 1e-5 {
			t.Errorf("Conversion incorrect: got %f, want 1.60934", result.Value())
		}
	})

	t.Run("ConvertIncompatibleUnits", func(t *testing.T) {
		q := NewQuantity(5, Meter)
		_, err := q.To(Second)
		if err == nil {
			t.Errorf("Conversion should fail for incompatible units")
		} else if !strings.Contains(err.Error(), "Cannot convert") {
			t.Errorf("Conversion error message incorrect: got %q", err.Error())
		}
	})

	t.Run("ConvertTemperatureCelsiusToKelvin", func(t *testing.T) {
		q := NewQuantity(0, Celsius)
		result, err := q.To(Kelvin)
		if err != nil {
			t.Errorf("Conversion failed: %v", err)
		}
		if math.Abs(result.Value()-273.15) > tolerance {
			t.Errorf("Conversion incorrect: got %f, want 273.15", result.Value())
		}
	})

	t.Run("ConvertTemperatureFahrenheitToCelsius", func(t *testing.T) {
		q := NewQuantity(32, Fahrenheit)
		result, err := q.To(Celsius)
		if err != nil {
			t.Errorf("Conversion failed: %v", err)
		}
		if math.Abs(result.Value()) > 0.01 {
			t.Errorf("Conversion incorrect: got %f, want 0", result.Value())
		}
	})
}

func TestPhysicsExamples(t *testing.T) {
	t.Run("VelocityCalculation", func(t *testing.T) {
		distance := NewQuantity(100, Meter)
		time := NewQuantity(9.58, Second)
		velocity := must(distance.Divide(time))
		if math.Abs(velocity.Value()-10.438) > 0.001 {
			t.Errorf("Velocity calculation incorrect: got %f, want 10.438", velocity.Value())
		}
	})

	t.Run("AccelerationCalculation", func(t *testing.T) {
		velocity := NewQuantity(10, MeterPerSecond)
		time := NewQuantity(2, Second)
		acceleration := must(velocity.Divide(time))
		if acceleration.Value() != 5 {
			t.Errorf("Acceleration calculation incorrect: got %f, want 5", acceleration.Value())
		}
	})

	t.Run("ForceCalculation", func(t *testing.T) {
		mass := NewQuantity(10, Kilogram)
		acceleration := NewQuantity(9.8, MeterPerSecondSquared)
		force := must(mass.Multiply(acceleration))
		if math.Abs(force.Value()-98) > tolerance {
			t.Errorf("Force calculation incorrect: got %f, want 98", force.Value())
		}
		if force.Unit().Dimension.Mass != 1 || force.Unit().Dimension.Length != 1 || force.Unit().Dimension.Time != -2 {
			t.Errorf("Force dimension incorrect")
		}
	})

	t.Run("KineticEnergy", func(t *testing.T) {
		mass := NewQuantity(2, Kilogram)
		velocity := NewQuantity(10, MeterPerSecond)
		ke := must(mass.Multiply(must(velocity.Power(2)))).MultiplyScalar(0.5)
		if ke.Value() != 100 {
			t.Errorf("Kinetic energy calculation incorrect: got %f, want 100", ke.Value())
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("ZeroQuantity", func(t *testing.T) {
		q := NewQuantity(0, Meter)
		if q.Value() != 0 {
			t.Errorf("Zero quantity failed")
		}
	})

	t.Run("NegativeQuantity", func(t *testing.T) {
		q := NewQuantity(-5, Meter)
		if q.Value() != -5 {
			t.Errorf("Negative quantity failed")
		}
	})

	t.Run("VeryLargeQuantity", func(t *testing.T) {
		q := NewQuantity(1e100, Meter)
		if q.Value() != 1e100 {
			t.Errorf("Very large quantity failed")
		}
	})

	t.Run("VerySmallQuantity", func(t *testing.T) {
		q := NewQuantity(1e-100, Meter)
		if q.Value() != 1e-100 {
			t.Errorf("Very small quantity failed")
		}
	})
}

// must unwraps a (Quantity, error) result and fails the test if err is non-nil.
func must(q Quantity, err error) Quantity {
	if err != nil {
		panic(err)
	}
	return q
}

// expectError verifies err is non-nil and is an *AffineUnitArithmeticError
// whose message contains wantSubstr.
func expectError(t *testing.T, name, wantSubstr string, err error) {
	t.Helper()
	if err == nil {
		t.Errorf("%s: expected error, got nil", name)
		return
	}
	var affine *AffineUnitArithmeticError
	if !errors.As(err, &affine) {
		t.Errorf("%s: expected *AffineUnitArithmeticError, got %T", name, err)
		return
	}
	if !strings.Contains(affine.Error(), wantSubstr) {
		t.Errorf("%s: error message %q does not contain %q", name, affine.Error(), wantSubstr)
	}
}

func TestAffineUnitArithmetic(t *testing.T) {
	t.Run("MultiplyCelsiusByCelsiusErrors", func(t *testing.T) {
		_, err := NewQuantity(2, Celsius).Multiply(NewQuantity(3, Celsius))
		expectError(t, "Multiply Celsius*Celsius", "Cannot multiply affine units", err)
	})

	t.Run("DivideMeterByCelsiusErrors", func(t *testing.T) {
		_, err := NewQuantity(100, Meter).Divide(NewQuantity(2, Celsius))
		expectError(t, "Divide Meter/Celsius", "Cannot divide affine units", err)
	})

	t.Run("DivideCelsiusByMeterErrors", func(t *testing.T) {
		_, err := NewQuantity(2, Celsius).Divide(NewQuantity(100, Meter))
		expectError(t, "Divide Celsius/Meter", "Cannot divide affine units", err)
	})

	t.Run("PowerCelsiusErrors", func(t *testing.T) {
		_, err := NewQuantity(2, Celsius).Power(2)
		expectError(t, "Power Celsius", "Cannot raise affine unit", err)
	})

	t.Run("MultiplyFahrenheitByFahrenheitErrors", func(t *testing.T) {
		_, err := NewQuantity(2, Fahrenheit).Multiply(NewQuantity(1, Fahrenheit))
		expectError(t, "Multiply Fahrenheit*Fahrenheit", "Cannot multiply affine units", err)
	})

	t.Run("ScalarOpsOnAffineUnitsAllowed", func(t *testing.T) {
		q := NewQuantity(2, Celsius)
		if got := q.MultiplyScalar(3).Value(); got != 6 {
			t.Errorf("MultiplyScalar on affine unit failed: got %f, want 6", got)
		}
		if got := q.DivideScalar(2).Value(); got != 1 {
			t.Errorf("DivideScalar on affine unit failed: got %f, want 1", got)
		}
	})

	t.Run("KelvinMultiplyAllowed", func(t *testing.T) {
		result := must(NewQuantity(2, Kelvin).Multiply(NewQuantity(3, Kelvin)))
		if result.Value() != 6 {
			t.Errorf("Kelvin multiply failed: got %f, want 6", result.Value())
		}
	})
}

func TestCanonicalUnits(t *testing.T) {
	t.Run("MeterTimesMeterYieldsSquareMeter", func(t *testing.T) {
		result := must(NewQuantity(5, Meter).Multiply(NewQuantity(3, Meter)))
		if result.Value() != 15 || result.Unit().Symbol != SquareMeter.Symbol {
			t.Errorf("got %s, want 15 m²", result.String())
		}
	})

	t.Run("MeterSquaredYieldsSquareMeter", func(t *testing.T) {
		result := must(NewQuantity(3, Meter).Power(2))
		if result.Value() != 9 || result.Unit().Symbol != SquareMeter.Symbol {
			t.Errorf("got %s, want 9 m²", result.String())
		}
	})

	t.Run("MeterPerSecondCanonical", func(t *testing.T) {
		result := must(NewQuantity(100, Meter).Divide(NewQuantity(10, Second)))
		if result.Value() != 10 || result.Unit().Symbol != MeterPerSecond.Symbol {
			t.Errorf("got %s, want 10 m/s", result.String())
		}
	})

	t.Run("ForceYieldsNewton", func(t *testing.T) {
		force := must(NewQuantity(10, Kilogram).Multiply(NewQuantity(9.8, MeterPerSecondSquared)))
		if math.Abs(force.Value()-98) > tolerance || force.Unit().Symbol != Newton.Symbol {
			t.Errorf("got %s, want 98 N", force.String())
		}
	})

	t.Run("KineticEnergyYieldsJoule", func(t *testing.T) {
		mass := NewQuantity(2, Kilogram)
		velocity := NewQuantity(10, MeterPerSecond)
		energy := must(mass.Multiply(must(velocity.Power(2))))
		if energy.Value() != 200 || energy.Unit().Symbol != Joule.Symbol {
			t.Errorf("got %s, want 200 J", energy.String())
		}
		ke := energy.MultiplyScalar(0.5)
		if ke.Value() != 100 || ke.Unit().Symbol != Joule.Symbol {
			t.Errorf("got %s, want 100 J", ke.String())
		}
	})

	t.Run("KilojoulePerSecondYieldsWatt", func(t *testing.T) {
		result := must(NewQuantity(1, Kilojoule).Divide(NewQuantity(1, Second)))
		if result.Value() != 1000 || result.Unit().Symbol != Watt.Symbol {
			t.Errorf("got %s, want 1000 W", result.String())
		}
	})

	t.Run("JoulePerMeterYieldsNewton", func(t *testing.T) {
		result := must(NewQuantity(1, Joule).Divide(NewQuantity(1, Meter)))
		if result.Value() != 1 || result.Unit().Symbol != Newton.Symbol {
			t.Errorf("got %s, want 1 N", result.String())
		}
	})

	t.Run("NewtonTimesMeterYieldsJoule", func(t *testing.T) {
		result := must(NewQuantity(1, Newton).Multiply(NewQuantity(1, Meter)))
		if result.Value() != 1 || result.Unit().Symbol != Joule.Symbol {
			t.Errorf("got %s, want 1 J", result.String())
		}
	})

	t.Run("PerSecondYieldsHertz", func(t *testing.T) {
		result := must(NewQuantity(1, Dimensionless).Divide(NewQuantity(1, Second)))
		if result.Value() != 1 || result.Unit().Symbol != Hertz.Symbol {
			t.Errorf("got %s, want 1 Hz", result.String())
		}
	})

	t.Run("UnknownDimensionKeepsGeneratedSymbol", func(t *testing.T) {
		result := must(NewQuantity(1, Meter).Multiply(NewQuantity(1, Second)))
		if result.Unit().Symbol != "m·s" {
			t.Errorf("got %s, want symbol m·s", result.String())
		}
	})
}

func TestIsClose(t *testing.T) {
	t.Run("SameUnitCloseWithinRelTol", func(t *testing.T) {
		q1 := NewQuantity(1.0, Meter)
		q2 := NewQuantity(1.0+1e-7, Meter)
		isClose, err := q1.IsClose(q2, 1e-6, 0.0)
		if err != nil {
			t.Errorf("IsClose returned unexpected error: %v", err)
		}
		if !isClose {
			t.Errorf("IsClose should return true within relTol")
		}
	})

	t.Run("SameUnitNotClose", func(t *testing.T) {
		q1 := NewQuantity(1.0, Meter)
		q2 := NewQuantity(2.0, Meter)
		isClose, err := q1.IsClose(q2, 1e-9, 0.0)
		if err != nil {
			t.Errorf("IsClose returned unexpected error: %v", err)
		}
		if isClose {
			t.Errorf("IsClose should return false for values far apart")
		}
	})

	t.Run("CloseWithAbsTolForVerySmallValues", func(t *testing.T) {
		q1 := NewQuantity(1e-12, Meter)
		q2 := NewQuantity(2e-12, Meter)
		isClose, err := q1.IsClose(q2, 1e-9, 1e-9)
		if err != nil {
			t.Errorf("IsClose returned unexpected error: %v", err)
		}
		if !isClose {
			t.Errorf("IsClose should return true for small values within absTol")
		}
	})

	t.Run("IncompatibleDimensionsError", func(t *testing.T) {
		q1 := NewQuantity(1.0, Meter)
		q2 := NewQuantity(1.0, Second)
		_, err := q1.IsClose(q2, 1e-9, 1e-9)
		if err == nil {
			t.Errorf("IsClose should return error for incompatible dimensions")
		}
	})

	t.Run("EqualValuesReturnTrue", func(t *testing.T) {
		q1 := NewQuantity(1.0, Meter)
		q2 := NewQuantity(1.0, Meter)
		isClose, err := q1.IsClose(q2, 0.0, 0.0)
		if err != nil {
			t.Errorf("IsClose returned unexpected error: %v", err)
		}
		if !isClose {
			t.Errorf("IsClose should return true for equal values")
		}
	})
}

func TestNewUnits(t *testing.T) {
	t.Run("AngleUnitsDefined", func(t *testing.T) {
		if Radian.Symbol != "rad" || Degree.Symbol != "°" ||
			Arcminute.Symbol != "′" || Arcsecond.Symbol != "″" {
			t.Errorf("Angle unit symbols incorrect")
		}
		if !Radian.Dimension.IsDimensionless() {
			t.Errorf("Radian should be dimensionless")
		}
	})

	t.Run("DegreeToRadian", func(t *testing.T) {
		result := must(NewQuantity(180, Degree).To(Radian))
		if math.Abs(result.Value()-math.Pi) > tolerance {
			t.Errorf("got %f, want %f", result.Value(), math.Pi)
		}
	})

	t.Run("RadianToDegree", func(t *testing.T) {
		result := must(NewQuantity(math.Pi/2, Radian).To(Degree))
		if math.Abs(result.Value()-90) > tolerance {
			t.Errorf("got %f, want 90", result.Value())
		}
	})

	t.Run("ArcminuteArcsecond", func(t *testing.T) {
		result := must(NewQuantity(1, Degree).To(Arcminute))
		if math.Abs(result.Value()-60) > tolerance {
			t.Errorf("got %f, want 60", result.Value())
		}
		result = must(NewQuantity(1, Degree).To(Arcsecond))
		if math.Abs(result.Value()-3600) > tolerance {
			t.Errorf("got %f, want 3600", result.Value())
		}
	})

	t.Run("FrequencyConversions", func(t *testing.T) {
		result := must(NewQuantity(1, Kilohertz).To(Hertz))
		if math.Abs(result.Value()-1000) > tolerance {
			t.Errorf("got %f, want 1000", result.Value())
		}
		result = must(NewQuantity(1, Megahertz).To(Kilohertz))
		if math.Abs(result.Value()-1000) > tolerance {
			t.Errorf("got %f, want 1000", result.Value())
		}
	})

	t.Run("AreaConversions", func(t *testing.T) {
		result := must(NewQuantity(1, SquareKilometer).To(SquareMeter))
		if math.Abs(result.Value()-1e6) > tolerance {
			t.Errorf("got %f, want 1e6", result.Value())
		}
		result = must(NewQuantity(1, Hectare).To(SquareMeter))
		if math.Abs(result.Value()-1e4) > tolerance {
			t.Errorf("got %f, want 1e4", result.Value())
		}
		result = must(NewQuantity(1, SquareKilometer).To(Hectare))
		if math.Abs(result.Value()-100) > tolerance {
			t.Errorf("got %f, want 100", result.Value())
		}
	})

	t.Run("VolumeConversions", func(t *testing.T) {
		result := must(NewQuantity(1, CubicMeter).To(Liter))
		if math.Abs(result.Value()-1000) > tolerance {
			t.Errorf("got %f, want 1000", result.Value())
		}
		result = must(NewQuantity(1, Liter).To(Milliliter))
		if math.Abs(result.Value()-1000) > tolerance {
			t.Errorf("got %f, want 1000", result.Value())
		}
	})

	t.Run("MilePerHourConversion", func(t *testing.T) {
		result := must(NewQuantity(1, MilePerHour).To(MeterPerSecond))
		if math.Abs(result.Value()-0.44704) > tolerance {
			t.Errorf("got %f, want 0.44704", result.Value())
		}
		result = must(NewQuantity(60, MilePerHour).To(KilometerPerHour))
		if math.Abs(result.Value()-96.5606) > 1e-3 {
			t.Errorf("got %f, want ~96.56", result.Value())
		}
	})

	t.Run("MoleDimension", func(t *testing.T) {
		if Mole.Symbol != "mol" || Mole.Dimension.Amount != 1 {
			t.Errorf("Mole unit incorrect")
		}
		_, err := NewQuantity(1, Mole).To(Kilogram)
		if err == nil {
			t.Errorf("Mole to Kilogram should fail: incompatible dimensions")
		}
	})

	t.Run("EnergyConversions", func(t *testing.T) {
		result := must(NewQuantity(1, Calorie).To(Joule))
		if math.Abs(result.Value()-4.184) > tolerance {
			t.Errorf("got %f, want 4.184", result.Value())
		}
		result = must(NewQuantity(1, Kilocalorie).To(Calorie))
		if math.Abs(result.Value()-1000) > tolerance {
			t.Errorf("got %f, want 1000", result.Value())
		}
		result = must(NewQuantity(1, WattHour).To(Joule))
		if math.Abs(result.Value()-3600) > tolerance {
			t.Errorf("got %f, want 3600", result.Value())
		}
		result = must(NewQuantity(1, WattHour).To(Kilojoule))
		if math.Abs(result.Value()-3.6) > tolerance {
			t.Errorf("got %f, want 3.6", result.Value())
		}
	})

	t.Run("VoltDimension", func(t *testing.T) {
		d := Volt.Dimension
		if d.Mass != 1 || d.Length != 2 || d.Time != -3 || d.Current != -1 {
			t.Errorf("Volt dimension incorrect: %+v", d)
		}
	})

	t.Run("OhmDimension", func(t *testing.T) {
		d := Ohm.Dimension
		if d.Mass != 1 || d.Length != 2 || d.Time != -3 || d.Current != -2 {
			t.Errorf("Ohm dimension incorrect: %+v", d)
		}
	})

	t.Run("WattPerAmpereYieldsVolt", func(t *testing.T) {
		result := must(NewQuantity(1, Watt).Divide(NewQuantity(1, Ampere)))
		if result.Value() != 1 || result.Unit().Symbol != Volt.Symbol {
			t.Errorf("got %s, want 1 V", result.String())
		}
	})

	t.Run("VoltPerAmpereYieldsOhm", func(t *testing.T) {
		result := must(NewQuantity(1, Volt).Divide(NewQuantity(1, Ampere)))
		if result.Value() != 1 || result.Unit().Symbol != Ohm.Symbol {
			t.Errorf("got %s, want 1 Ω", result.String())
		}
	})

	t.Run("VoltTimesAmpereYieldsWatt", func(t *testing.T) {
		result := must(NewQuantity(12, Volt).Multiply(NewQuantity(2, Ampere)))
		if result.Value() != 24 || result.Unit().Symbol != Watt.Symbol {
			t.Errorf("got %s, want 24 W", result.String())
		}
	})

	t.Run("MeterCubedYieldsCubicMeter", func(t *testing.T) {
		result := must(NewQuantity(2, Meter).Power(3))
		if result.Value() != 8 || result.Unit().Symbol != CubicMeter.Symbol {
			t.Errorf("got %s, want 8 m³", result.String())
		}
	})

	t.Run("KilometerTimesKilometerYieldsSquareMeter", func(t *testing.T) {
		result := must(NewQuantity(1, Kilometer).Multiply(NewQuantity(1, Kilometer)))
		if result.Value() != 1e6 || result.Unit().Symbol != SquareMeter.Symbol {
			t.Errorf("got %s, want 1000000 m²", result.String())
		}
	})

	t.Run("AddAngleToLengthFails", func(t *testing.T) {
		// Angles are dimensionless, so they cannot be added to lengths.
		_, err := NewQuantity(1, Degree).Add(NewQuantity(1, Meter))
		if err == nil {
			t.Errorf("Add should fail for degree and meter")
		}
	})

	t.Run("MoleIsCompatibleAcrossAmountUnits", func(t *testing.T) {
		// Mole is the only Amount unit; verify dimension equality via Add.
		result := must(NewQuantity(1, Mole).Add(NewQuantity(2, Mole)))
		if result.Value() != 3 {
			t.Errorf("got %f, want 3", result.Value())
		}
	})
}
