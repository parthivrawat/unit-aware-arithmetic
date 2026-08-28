"""
Tests for Unit-Aware Numeric Arithmetic Library
"""

import pytest
import math
from dimensional import Quantity, Unit, Dimension, units, IncompatibleUnitsError


class TestDimension:
    """Test Dimension class."""
    
    def test_dimension_creation(self):
        d = Dimension(length=1, time=-2)
        assert d.length == 1
        assert d.time == -2
        assert d.mass == 0
    
    def test_dimension_multiplication(self):
        d1 = Dimension(length=1)  # L
        d2 = Dimension(time=-1)   # T^-1
        result = d1 * d2          # L·T^-1
        assert result.length == 1
        assert result.time == -1
    
    def test_dimension_division(self):
        d1 = Dimension(length=1, time=-1)  # L·T^-1
        d2 = Dimension(time=1)             # T
        result = d1 / d2                   # L·T^-2
        assert result.length == 1
        assert result.time == -2
    
    def test_dimension_power(self):
        d = Dimension(length=1, time=-1)  # L·T^-1
        result = d ** 2                   # L²·T^-2
        assert result.length == 2
        assert result.time == -2
    
    def test_is_dimensionless(self):
        assert Dimension().is_dimensionless()
        assert not Dimension(length=1).is_dimensionless()


class TestUnit:
    """Test Unit class."""
    
    def test_unit_creation(self):
        u = Unit("meter", "m", Dimension(length=1))
        assert u.name == "meter"
        assert u.symbol == "m"
        assert u.dimension.length == 1
    
    def test_unit_equality(self):
        u1 = Unit("meter", "m", Dimension(length=1))
        u2 = Unit("meter", "m", Dimension(length=1))
        assert u1 == u2
    
    def test_predefined_units(self):
        assert units.meter.symbol == "m"
        assert units.kilogram.symbol == "kg"
        assert units.second.symbol == "s"


class TestQuantityBasics:
    """Test basic Quantity operations."""
    
    def test_quantity_creation(self):
        q = Quantity(5.0, units.meter)
        assert q.value == 5.0
        assert q.unit == units.meter
    
    def test_quantity_repr(self):
        q = Quantity(5.0, units.meter)
        assert str(q) == "5.0 m"
    
    def test_quantity_with_integer(self):
        q = Quantity(5, units.meter)
        assert q.value == 5.0


class TestArithmetic:
    """Test arithmetic operations."""
    
    def test_addition_same_unit(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(3, units.meter)
        result = q1 + q2
        assert result.value == 8
        assert result.unit == units.meter
    
    def test_addition_compatible_units(self):
        q1 = Quantity(1, units.meter)
        q2 = Quantity(100, units.centimeter)
        result = q1 + q2
        assert math.isclose(result.value, 2.0)
        assert result.unit == units.meter
    
    def test_addition_incompatible_units(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(3, units.second)
        with pytest.raises(IncompatibleUnitsError):
            q1 + q2
    
    def test_subtraction_same_unit(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(3, units.meter)
        result = q1 - q2
        assert result.value == 2
        assert result.unit == units.meter
    
    def test_subtraction_compatible_units(self):
        q1 = Quantity(1, units.kilometer)
        q2 = Quantity(500, units.meter)
        result = q1 - q2
        assert math.isclose(result.value, 0.5)
        assert result.unit == units.kilometer
    
    def test_subtraction_incompatible_units(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(3, units.kilogram)
        with pytest.raises(IncompatibleUnitsError):
            q1 - q2
    
    def test_multiplication_by_scalar(self):
        q = Quantity(5, units.meter)
        result = q * 3
        assert result.value == 15
        assert result.unit == units.meter
    
    def test_multiplication_scalar_left(self):
        q = Quantity(5, units.meter)
        result = 3 * q
        assert result.value == 15
        assert result.unit == units.meter
    
    def test_multiplication_quantities(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(3, units.meter)
        result = q1 * q2
        assert result.value == 15
        assert result.unit.dimension.length == 2  # m²
    
    def test_division_by_scalar(self):
        q = Quantity(10, units.meter)
        result = q / 2
        assert result.value == 5
        assert result.unit == units.meter
    
    def test_division_quantities(self):
        q1 = Quantity(100, units.meter)
        q2 = Quantity(10, units.second)
        result = q1 / q2
        assert result.value == 10
        assert result.unit.dimension.length == 1
        assert result.unit.dimension.time == -1
    
    def test_power_integer(self):
        q = Quantity(3, units.meter)
        result = q ** 2
        assert result.value == 9
        assert result.unit.dimension.length == 2
    
    def test_negation(self):
        q = Quantity(5, units.meter)
        result = -q
        assert result.value == -5
        assert result.unit == units.meter
    
    def test_absolute_value(self):
        q = Quantity(-5, units.meter)
        result = abs(q)
        assert result.value == 5
        assert result.unit == units.meter


class TestComparison:
    """Test comparison operations."""
    
    def test_equality_same_unit(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(5, units.meter)
        assert q1 == q2
    
    def test_equality_different_compatible_units(self):
        q1 = Quantity(1, units.meter)
        q2 = Quantity(100, units.centimeter)
        assert q1 == q2
    
    def test_inequality_different_values(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(3, units.meter)
        assert q1 != q2
    
    def test_less_than(self):
        q1 = Quantity(3, units.meter)
        q2 = Quantity(5, units.meter)
        assert q1 < q2
    
    def test_less_than_or_equal(self):
        q1 = Quantity(3, units.meter)
        q2 = Quantity(5, units.meter)
        q3 = Quantity(3, units.meter)
        assert q1 <= q2
        assert q1 <= q3
    
    def test_greater_than(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(3, units.meter)
        assert q1 > q2
    
    def test_greater_than_or_equal(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(3, units.meter)
        q3 = Quantity(5, units.meter)
        assert q1 >= q2
        assert q1 >= q3
    
    def test_comparison_incompatible_units(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(3, units.second)
        with pytest.raises(IncompatibleUnitsError):
            q1 < q2


class TestConversion:
    """Test unit conversion."""
    
    def test_convert_length(self):
        q = Quantity(1, units.meter)
        result = q.to(units.centimeter)
        assert math.isclose(result.value, 100)
        assert result.unit == units.centimeter
    
    def test_convert_mass(self):
        q = Quantity(1, units.kilogram)
        result = q.to(units.gram)
        assert math.isclose(result.value, 1000)
        assert result.unit == units.gram
    
    def test_convert_time(self):
        q = Quantity(1, units.hour)
        result = q.to(units.second)
        assert math.isclose(result.value, 3600)
        assert result.unit == units.second
    
    def test_convert_imperial_to_metric(self):
        q = Quantity(1, units.mile)
        result = q.to(units.kilometer)
        assert math.isclose(result.value, 1.60934, rel_tol=1e-5)
        assert result.unit == units.kilometer
    
    def test_convert_incompatible_units(self):
        q = Quantity(5, units.meter)
        with pytest.raises(IncompatibleUnitsError):
            q.to(units.second)
    
    def test_convert_temperature_celsius_to_kelvin(self):
        q = Quantity(0, units.celsius)
        result = q.to(units.kelvin)
        assert math.isclose(result.value, 273.15)
        assert result.unit == units.kelvin
    
    def test_convert_temperature_fahrenheit_to_celsius(self):
        q = Quantity(32, units.fahrenheit)
        result = q.to(units.celsius)
        assert math.isclose(result.value, 0, abs_tol=0.01)
        assert result.unit == units.celsius


class TestPhysicsExamples:
    """Test real-world physics examples."""
    
    def test_velocity_calculation(self):
        distance = Quantity(100, units.meter)
        time = Quantity(9.58, units.second)
        velocity = distance / time
        assert math.isclose(velocity.value, 10.438, rel_tol=1e-3)
        assert velocity.unit.dimension.length == 1
        assert velocity.unit.dimension.time == -1
    
    def test_acceleration_calculation(self):
        velocity = Quantity(10, units.meter_per_second)
        time = Quantity(2, units.second)
        acceleration = velocity / time
        assert acceleration.value == 5
        assert acceleration.unit.dimension.length == 1
        assert acceleration.unit.dimension.time == -2
    
    def test_force_calculation(self):
        # F = ma
        mass = Quantity(10, units.kilogram)
        acceleration = Quantity(9.8, units.meter_per_second_squared)
        force = mass * acceleration
        assert math.isclose(force.value, 98)
        assert force.unit.dimension.mass == 1
        assert force.unit.dimension.length == 1
        assert force.unit.dimension.time == -2
    
    def test_kinetic_energy(self):
        # KE = 1/2 * m * v²
        mass = Quantity(2, units.kilogram)
        velocity = Quantity(10, units.meter_per_second)
        ke = 0.5 * mass * (velocity ** 2)
        assert ke.value == 100
        assert ke.unit.dimension.mass == 1
        assert ke.unit.dimension.length == 2
        assert ke.unit.dimension.time == -2


class TestEdgeCases:
    """Test edge cases and error handling."""
    
    def test_zero_quantity(self):
        q = Quantity(0, units.meter)
        assert q.value == 0
    
    def test_negative_quantity(self):
        q = Quantity(-5, units.meter)
        assert q.value == -5
    
    def test_very_large_quantity(self):
        q = Quantity(1e100, units.meter)
        assert q.value == 1e100
    
    def test_very_small_quantity(self):
        q = Quantity(1e-100, units.meter)
        assert q.value == 1e-100
    
    def test_type_error_on_invalid_addition(self):
        q = Quantity(5, units.meter)
        with pytest.raises(TypeError):
            q + 5
    
    def test_type_error_on_invalid_multiplication(self):
        q = Quantity(5, units.meter)
        with pytest.raises(TypeError):
            q * "string"


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
