"""
Tests for Unit-Aware Numeric Arithmetic Library
"""

import pytest
import math
from dataclasses import FrozenInstanceError
from dimensional import (
    Quantity,
    Unit,
    Dimension,
    units,
    IncompatibleUnitsError,
    AffineUnitArithmeticError,
)


class TestDimension:
    """Test Dimension class."""

    def test_dimension_creation(self):
        d = Dimension(length=1, time=-2)
        assert d.length == 1
        assert d.time == -2
        assert d.mass == 0

    def test_dimension_multiplication(self):
        d1 = Dimension(length=1)  # L
        d2 = Dimension(time=-1)  # T^-1
        result = d1 * d2  # L·T^-1
        assert result.length == 1
        assert result.time == -1

    def test_dimension_division(self):
        d1 = Dimension(length=1, time=-1)  # L·T^-1
        d2 = Dimension(time=1)  # T
        result = d1 / d2  # L·T^-2
        assert result.length == 1
        assert result.time == -2

    def test_dimension_power(self):
        d = Dimension(length=1, time=-1)  # L·T^-1
        result = d**2  # L²·T^-2
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

    def test_convert_to_base_and_from_base_linear(self):
        assert math.isclose(units.centimeter.convert_to_base(150.0), 1.5)
        assert math.isclose(units.centimeter.convert_from_base(1.5), 150.0)

    def test_convert_to_base_and_from_base_affine(self):
        # 0 °C = 273.15 K
        assert math.isclose(units.celsius.convert_to_base(0.0), 273.15)
        assert math.isclose(units.celsius.convert_from_base(273.15), 0.0)
        # 32 °F = 273.15 K
        assert math.isclose(
            units.fahrenheit.convert_to_base(32.0), 273.15, rel_tol=1e-9
        )
        assert math.isclose(
            units.fahrenheit.convert_from_base(273.15), 32.0, abs_tol=1e-9
        )

    def test_unit_is_compatible_with(self):
        assert units.meter.is_compatible_with(units.foot)
        assert not units.meter.is_compatible_with(units.second)


class TestQuantityBasics:
    """Test basic Quantity operations."""

    def test_quantity_creation(self):
        q = Quantity(5.0, units.meter)
        assert q.value == 5.0
        assert q.unit == units.meter

    def test_quantity_repr(self):
        q = Quantity(5.0, units.meter)
        assert repr(q) == "Quantity(5.0, Unit('m'))"

    def test_quantity_str(self):
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

    def test_rtruediv_scalar_divided_by_quantity(self):
        q = Quantity(2, units.meter)
        result = 1 / q
        assert math.isclose(result.value, 0.5)
        assert result.unit.dimension.length == -1

    def test_rtruediv_frequency(self):
        period = Quantity(4, units.second)
        result = 1 / period
        assert math.isclose(result.value, 0.25)
        assert result.unit == units.hertz

    def test_rtruediv_invalid_operand(self):
        q = Quantity(2, units.meter)
        with pytest.raises(TypeError):
            "string" / q

    def test_power_integer(self):
        q = Quantity(3, units.meter)
        result = q**2
        assert result.value == 9
        assert result.unit.dimension.length == 2

    def test_multiplication_affine_units(self):
        q1 = Quantity(2, units.celsius)
        q2 = Quantity(3, units.celsius)
        with pytest.raises(AffineUnitArithmeticError):
            q1 * q2

    def test_multiplication_affine_fahrenheit(self):
        q1 = Quantity(2, units.fahrenheit)
        q2 = Quantity(1, units.fahrenheit)
        with pytest.raises(AffineUnitArithmeticError):
            q1 * q2

    def test_division_affine_units(self):
        q1 = Quantity(100, units.meter)
        q2 = Quantity(2, units.celsius)
        with pytest.raises(AffineUnitArithmeticError):
            q1 / q2

    def test_power_affine_unit(self):
        q = Quantity(2, units.celsius)
        with pytest.raises(AffineUnitArithmeticError):
            q**2

    def test_scalar_multiplication_affine_allowed(self):
        q = Quantity(2, units.celsius)
        result = q * 3
        assert result.value == 6
        assert result.unit == units.celsius

    def test_scalar_division_affine_allowed(self):
        q = Quantity(10, units.celsius)
        result = q / 2
        assert result.value == 5
        assert result.unit == units.celsius

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


class TestIsClose:
    """Test tolerance-based approximate equality."""

    def test_is_close_default_rel_tol(self):
        q1 = Quantity(1, units.meter)
        q2 = Quantity(1.0000001, units.meter)
        # diff is 1e-7, so rel_tol must exceed ~1e-7 to be "close"
        assert not q1.is_close(q2, rel_tol=1e-9)
        assert q1.is_close(q2, rel_tol=1e-6)

    def test_is_close_abs_tol(self):
        q1 = Quantity(1e-9, units.meter)
        q2 = Quantity(2e-9, units.meter)
        assert q1.is_close(q2, rel_tol=0.0, abs_tol=1.5e-9)

    def test_is_close_relative(self):
        q1 = Quantity(1e-9, units.meter)
        q2 = Quantity(2e-9, units.meter)
        # diff is 1e-9; relative tolerance must be >= 0.5 to cover it
        assert not q1.is_close(q2, rel_tol=1e-9)
        assert q1.is_close(q2, rel_tol=0.6)

    def test_is_close_different_values(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(3, units.meter)
        assert not q1.is_close(q2)

    def test_is_close_incompatible_dimensions(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(5, units.second)
        assert not q1.is_close(q2)

    def test_is_close_near_zero_with_abs_tol(self):
        q1 = Quantity(0, units.meter)
        q2 = Quantity(1e-12, units.meter)
        assert q1.is_close(q2, abs_tol=1e-9)


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
        ke = 0.5 * mass * (velocity**2)
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


class TestCanonicalUnitLookup:
    """Test canonicalization of derived dimensions into base/SI units."""

    def test_area_from_length_multiplication(self):
        result = Quantity(5, units.meter) * Quantity(3, units.meter)
        assert math.isclose(result.value, 15)
        assert result.unit == units.square_meter

    def test_area_from_power(self):
        result = Quantity(3, units.meter) ** 2
        assert math.isclose(result.value, 9)
        assert result.unit == units.square_meter

    def test_velocity_from_division(self):
        result = Quantity(100, units.meter) / Quantity(10, units.second)
        assert math.isclose(result.value, 10)
        assert result.unit == units.meter_per_second

    def test_force_from_mass_times_acceleration(self):
        result = Quantity(10, units.kilogram) * Quantity(
            9.8, units.meter_per_second_squared
        )
        assert math.isclose(result.value, 98)
        assert result.unit == units.newton

    def test_kinetic_energy_canonicalizes_to_joule(self):
        result = (
            0.5
            * Quantity(2, units.kilogram)
            * (Quantity(10, units.meter_per_second) ** 2)
        )
        assert math.isclose(result.value, 100)
        assert result.unit == units.joule

    def test_power_from_kilojoule_per_second(self):
        result = Quantity(1, units.kilojoule) / Quantity(1, units.second)
        assert math.isclose(result.value, 1000)
        assert result.unit == units.watt

    def test_force_from_joule_per_meter(self):
        result = Quantity(1, units.joule) / Quantity(1, units.meter)
        assert math.isclose(result.value, 1)
        assert result.unit == units.newton

    def test_energy_from_force_times_length(self):
        result = Quantity(1, units.newton) * Quantity(1, units.meter)
        assert math.isclose(result.value, 1)
        assert result.unit == units.joule

    def test_hertz_from_dimensionless_per_second(self):
        result = Quantity(1, units.dimensionless) / Quantity(1, units.second)
        assert math.isclose(result.value, 1)
        assert result.unit == units.hertz

    def test_voltage_canonicalizes_to_volt(self):
        result = Quantity(1, units.ampere) * Quantity(1, units.ohm)
        assert math.isclose(result.value, 1)
        assert result.unit == units.volt

    def test_resistance_canonicalizes_to_ohm(self):
        result = Quantity(1, units.volt) / Quantity(1, units.ampere)
        assert math.isclose(result.value, 1)
        assert result.unit == units.ohm

    def test_volume_canonicalizes_to_cubic_meter(self):
        result = Quantity(2, units.meter) ** 3
        assert math.isclose(result.value, 8)
        assert result.unit == units.cubic_meter

    def test_frequency_canonicalizes_to_hertz_from_khz(self):
        result = Quantity(1, units.kilohertz) / Quantity(1000, units.dimensionless)
        assert math.isclose(result.value, 1)
        assert result.unit == units.hertz

    def test_dimensionless_result_uses_dimensionless_unit(self):
        # The canonical unit for Dimension() must stay `dimensionless`,
        # not radian (which is also dimensionless with to_base == 1).
        result = Quantity(2, units.meter) / Quantity(1, units.meter)
        assert math.isclose(result.value, 2.0)
        assert result.unit == units.dimensionless


class TestHashingAndImmutability:
    """Test that Quantity is hashable and immutable."""

    def test_quantity_in_set(self):
        q1 = Quantity(5, units.meter)
        q2 = Quantity(3, units.second)
        s = {q1, q2}
        assert len(s) == 2
        assert q1 in s

    def test_quantity_as_dict_key(self):
        q = Quantity(5, units.meter)
        d = {q: "length"}
        assert d[q] == "length"

    def test_hash_consistent_with_equality(self):
        # Equal quantities must have equal hashes
        q1 = Quantity(1, units.meter)
        q2 = Quantity(100, units.centimeter)
        assert q1 == q2
        assert hash(q1) == hash(q2)

    def test_unit_hashable(self):
        u1 = Unit("meter", "m", Dimension(length=1))
        u2 = Unit("meter", "m", Dimension(length=1))
        assert hash(u1) == hash(u2)
        assert len({u1, u2, units.meter}) == 1

    def test_quantity_immutable(self):
        q = Quantity(5, units.meter)
        with pytest.raises(FrozenInstanceError):
            q.value = 10
        with pytest.raises(FrozenInstanceError):
            q.unit = units.second


class TestFormatting:
    """Test __format__ support."""

    def test_format_empty_spec(self):
        q = Quantity(5.0, units.meter)
        assert f"{q}" == "5.0 m"

    def test_format_precision(self):
        q = Quantity(5.0, units.meter)
        assert f"{q:.2f}" == "5.00 m"

    def test_format_scientific(self):
        q = Quantity(12345, units.meter)
        assert f"{q:.2e}" == "1.23e+04 m"


class TestNewUnits:
    """Test the newly added units and conversions."""

    def test_angle_radian_and_degree(self):
        assert math.isclose(Quantity(180, units.degree).to(units.radian).value, math.pi)
        assert math.isclose(
            Quantity(1, units.radian).to(units.degree).value,
            180.0 / math.pi,
            rel_tol=1e-9,
        )

    def test_angle_arcminute_and_arcsecond(self):
        assert math.isclose(Quantity(60, units.arcminute).to(units.degree).value, 1.0)
        assert math.isclose(
            Quantity(60, units.arcsecond).to(units.arcminute).value, 1.0
        )

    def test_frequency_kilohertz_and_megahertz(self):
        assert math.isclose(Quantity(1, units.kilohertz).to(units.hertz).value, 1000.0)
        assert math.isclose(Quantity(1, units.megahertz).to(units.hertz).value, 1e6)
        assert math.isclose(Quantity(1000, units.hertz).to(units.kilohertz).value, 1.0)

    def test_area_square_kilometer_and_hectare(self):
        assert math.isclose(
            Quantity(1, units.square_kilometer).to(units.square_meter).value,
            1e6,
        )
        assert math.isclose(
            Quantity(1, units.hectare).to(units.square_meter).value, 1e4
        )

    def test_volume_liter_and_milliliter(self):
        assert math.isclose(Quantity(1, units.liter).to(units.cubic_meter).value, 0.001)
        assert math.isclose(Quantity(1, units.milliliter).to(units.liter).value, 0.001)

    def test_velocity_mile_per_hour(self):
        assert math.isclose(
            Quantity(1, units.mile_per_hour).to(units.meter_per_second).value,
            0.44704,
        )
        assert math.isclose(
            Quantity(1, units.mile_per_hour).to(units.kilometer_per_hour).value,
            1.609344,
        )

    def test_chemistry_mole(self):
        q = Quantity(2, units.mole)
        assert math.isclose(q.to(units.mole).value, 2.0)

    def test_energy_calorie_kilocalorie_watt_hour(self):
        assert math.isclose(Quantity(1, units.calorie).to(units.joule).value, 4.184)
        assert math.isclose(
            Quantity(1, units.kilocalorie).to(units.joule).value, 4184.0
        )
        assert math.isclose(Quantity(1, units.watt_hour).to(units.joule).value, 3600.0)

    def test_electricity_volt_and_ohm(self):
        # V = A · Ω
        voltage = Quantity(1, units.ampere) * Quantity(1, units.ohm)
        assert math.isclose(voltage.to(units.volt).value, 1.0)
        assert voltage.unit == units.volt

        # Ω = V / A
        resistance = Quantity(1, units.volt) / Quantity(1, units.ampere)
        assert math.isclose(resistance.to(units.ohm).value, 1.0)
        assert resistance.unit == units.ohm

    def test_arithmetic_derived_molar_and_frequency(self):
        # Molar concentration from mol / volume
        concentration = Quantity(1, units.mole) / Quantity(1, units.liter)
        assert concentration.unit.dimension.amount == 1
        assert concentration.unit.dimension.length == -3
        assert math.isclose(concentration.value, 1.0)
        assert concentration.unit.symbol == "mol/L"

        # 1 / 0.001 s = 1000 Hz
        freq = Quantity(1, units.dimensionless) / Quantity(0.001, units.second)
        assert math.isclose(freq.to(units.hertz).value, 1000.0)
        assert freq.unit == units.hertz


class TestValueIn:
    """Test the value_in convenience alias."""

    def test_value_in_compatible_unit(self):
        q = Quantity(1, units.kilometer)
        assert math.isclose(q.value_in(units.meter), 1000)

    def test_value_in_same_unit(self):
        q = Quantity(5, units.second)
        assert q.value_in(units.second) == 5.0

    def test_value_in_incompatible_unit(self):
        q = Quantity(5, units.meter)
        with pytest.raises(IncompatibleUnitsError):
            q.value_in(units.kilogram)


if __name__ == "__main__":
    pytest.main([__file__, "-v"])
