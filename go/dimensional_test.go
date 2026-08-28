package dimensional

import (
	"math"
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
		if q.Value != 5.0 || q.Unit.Symbol != "m" {
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
		if result.Value != 8 {
			t.Errorf("Add result incorrect: got %f, want 8", result.Value)
		}
	})

	t.Run("AddCompatibleUnits", func(t *testing.T) {
		q1 := NewQuantity(1, Meter)
		q2 := NewQuantity(100, Centimeter)
		result, err := q1.Add(q2)
		if err != nil {
			t.Errorf("Add failed: %v", err)
		}
		if math.Abs(result.Value-2.0) > tolerance {
			t.Errorf("Add result incorrect: got %f, want 2.0", result.Value)
		}
	})

	t.Run("AddIncompatibleUnits", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Second)
		_, err := q1.Add(q2)
		if err == nil {
			t.Errorf("Add should fail for incompatible units")
		}
	})

	t.Run("SubtractSameUnit", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Meter)
		result, err := q1.Subtract(q2)
		if err != nil {
			t.Errorf("Subtract failed: %v", err)
		}
		if result.Value != 2 {
			t.Errorf("Subtract result incorrect: got %f, want 2", result.Value)
		}
	})

	t.Run("MultiplyByScalar", func(t *testing.T) {
		q := NewQuantity(5, Meter)
		result := q.Multiply(3.0)
		if result.Value != 15 {
			t.Errorf("Multiply by scalar failed: got %f, want 15", result.Value)
		}
	})

	t.Run("MultiplyQuantities", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Meter)
		result := q1.Multiply(q2)
		if result.Value != 15 {
			t.Errorf("Multiply quantities failed: got %f, want 15", result.Value)
		}
		if result.Unit.Dimension.Length != 2 {
			t.Errorf("Multiply dimension incorrect: got %d, want 2", result.Unit.Dimension.Length)
		}
	})

	t.Run("DivideByScalar", func(t *testing.T) {
		q := NewQuantity(10, Meter)
		result := q.Divide(2.0)
		if result.Value != 5 {
			t.Errorf("Divide by scalar failed: got %f, want 5", result.Value)
		}
	})

	t.Run("DivideQuantities", func(t *testing.T) {
		q1 := NewQuantity(100, Meter)
		q2 := NewQuantity(10, Second)
		result := q1.Divide(q2)
		if result.Value != 10 {
			t.Errorf("Divide quantities failed: got %f, want 10", result.Value)
		}
		if result.Unit.Dimension.Length != 1 || result.Unit.Dimension.Time != -1 {
			t.Errorf("Divide dimension incorrect")
		}
	})

	t.Run("Power", func(t *testing.T) {
		q := NewQuantity(3, Meter)
		result := q.Power(2)
		if result.Value != 9 {
			t.Errorf("Power failed: got %f, want 9", result.Value)
		}
		if result.Unit.Dimension.Length != 2 {
			t.Errorf("Power dimension incorrect: got %d, want 2", result.Unit.Dimension.Length)
		}
	})

	t.Run("Negate", func(t *testing.T) {
		q := NewQuantity(5, Meter)
		result := q.Negate()
		if result.Value != -5 {
			t.Errorf("Negate failed: got %f, want -5", result.Value)
		}
	})

	t.Run("Abs", func(t *testing.T) {
		q := NewQuantity(-5, Meter)
		result := q.Abs()
		if result.Value != 5 {
			t.Errorf("Abs failed: got %f, want 5", result.Value)
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

	t.Run("CompareIncompatibleUnits", func(t *testing.T) {
		q1 := NewQuantity(5, Meter)
		q2 := NewQuantity(3, Second)
		_, err := q1.LessThan(q2)
		if err == nil {
			t.Errorf("Comparison should fail for incompatible units")
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
		if math.Abs(result.Value-100) > tolerance {
			t.Errorf("Conversion incorrect: got %f, want 100", result.Value)
		}
	})

	t.Run("ConvertMass", func(t *testing.T) {
		q := NewQuantity(1, Kilogram)
		result, err := q.To(Gram)
		if err != nil {
			t.Errorf("Conversion failed: %v", err)
		}
		if math.Abs(result.Value-1000) > tolerance {
			t.Errorf("Conversion incorrect: got %f, want 1000", result.Value)
		}
	})

	t.Run("ConvertTime", func(t *testing.T) {
		q := NewQuantity(1, Hour)
		result, err := q.To(Second)
		if err != nil {
			t.Errorf("Conversion failed: %v", err)
		}
		if math.Abs(result.Value-3600) > tolerance {
			t.Errorf("Conversion incorrect: got %f, want 3600", result.Value)
		}
	})

	t.Run("ConvertImperialToMetric", func(t *testing.T) {
		q := NewQuantity(1, Mile)
		result, err := q.To(Kilometer)
		if err != nil {
			t.Errorf("Conversion failed: %v", err)
		}
		if math.Abs(result.Value-1.60934) > 1e-5 {
			t.Errorf("Conversion incorrect: got %f, want 1.60934", result.Value)
		}
	})

	t.Run("ConvertIncompatibleUnits", func(t *testing.T) {
		q := NewQuantity(5, Meter)
		_, err := q.To(Second)
		if err == nil {
			t.Errorf("Conversion should fail for incompatible units")
		}
	})

	t.Run("ConvertTemperatureCelsiusToKelvin", func(t *testing.T) {
		q := NewQuantity(0, Celsius)
		result, err := q.To(Kelvin)
		if err != nil {
			t.Errorf("Conversion failed: %v", err)
		}
		if math.Abs(result.Value-273.15) > tolerance {
			t.Errorf("Conversion incorrect: got %f, want 273.15", result.Value)
		}
	})

	t.Run("ConvertTemperatureFahrenheitToCelsius", func(t *testing.T) {
		q := NewQuantity(32, Fahrenheit)
		result, err := q.To(Celsius)
		if err != nil {
			t.Errorf("Conversion failed: %v", err)
		}
		if math.Abs(result.Value) > 0.01 {
			t.Errorf("Conversion incorrect: got %f, want 0", result.Value)
		}
	})
}

func TestPhysicsExamples(t *testing.T) {
	t.Run("VelocityCalculation", func(t *testing.T) {
		distance := NewQuantity(100, Meter)
		time := NewQuantity(9.58, Second)
		velocity := distance.Divide(time)
		if math.Abs(velocity.Value-10.438) > 0.001 {
			t.Errorf("Velocity calculation incorrect: got %f, want 10.438", velocity.Value)
		}
	})

	t.Run("AccelerationCalculation", func(t *testing.T) {
		velocity := NewQuantity(10, MeterPerSecond)
		time := NewQuantity(2, Second)
		acceleration := velocity.Divide(time)
		if acceleration.Value != 5 {
			t.Errorf("Acceleration calculation incorrect: got %f, want 5", acceleration.Value)
		}
	})

	t.Run("ForceCalculation", func(t *testing.T) {
		mass := NewQuantity(10, Kilogram)
		acceleration := NewQuantity(9.8, MeterPerSecondSquared)
		force := mass.Multiply(acceleration)
		if math.Abs(force.Value-98) > tolerance {
			t.Errorf("Force calculation incorrect: got %f, want 98", force.Value)
		}
		if force.Unit.Dimension.Mass != 1 || force.Unit.Dimension.Length != 1 || force.Unit.Dimension.Time != -2 {
			t.Errorf("Force dimension incorrect")
		}
	})

	t.Run("KineticEnergy", func(t *testing.T) {
		mass := NewQuantity(2, Kilogram)
		velocity := NewQuantity(10, MeterPerSecond)
		ke := mass.Multiply(velocity.Power(2)).Multiply(0.5)
		if ke.Value != 100 {
			t.Errorf("Kinetic energy calculation incorrect: got %f, want 100", ke.Value)
		}
	})
}

func TestEdgeCases(t *testing.T) {
	t.Run("ZeroQuantity", func(t *testing.T) {
		q := NewQuantity(0, Meter)
		if q.Value != 0 {
			t.Errorf("Zero quantity failed")
		}
	})

	t.Run("NegativeQuantity", func(t *testing.T) {
		q := NewQuantity(-5, Meter)
		if q.Value != -5 {
			t.Errorf("Negative quantity failed")
		}
	})

	t.Run("VeryLargeQuantity", func(t *testing.T) {
		q := NewQuantity(1e100, Meter)
		if q.Value != 1e100 {
			t.Errorf("Very large quantity failed")
		}
	})

	t.Run("VerySmallQuantity", func(t *testing.T) {
		q := NewQuantity(1e-100, Meter)
		if q.Value != 1e-100 {
			t.Errorf("Very small quantity failed")
		}
	})
}
