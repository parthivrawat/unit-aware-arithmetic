"""
Unit-Aware Numeric Arithmetic Library

A type-safe dimensional arithmetic library that tracks units at runtime
and prevents invalid operations.

Zero dependencies, production-ready.
"""

from __future__ import annotations
from typing import Dict, Optional, Union, Any
from dataclasses import dataclass
import math


@dataclass(frozen=True)
class Dimension:
    """Represents the dimensional formula of a unit (e.g., L^1 T^-2 for acceleration)."""

    length: int = 0  # L
    mass: int = 0  # M
    time: int = 0  # T
    current: int = 0  # I
    temperature: int = 0  # Θ
    amount: int = 0  # N
    luminosity: int = 0  # J

    def __mul__(self, other: Dimension) -> Dimension:
        """Multiply dimensions (add exponents)."""
        return Dimension(
            length=self.length + other.length,
            mass=self.mass + other.mass,
            time=self.time + other.time,
            current=self.current + other.current,
            temperature=self.temperature + other.temperature,
            amount=self.amount + other.amount,
            luminosity=self.luminosity + other.luminosity,
        )

    def __truediv__(self, other: Dimension) -> Dimension:
        """Divide dimensions (subtract exponents)."""
        return Dimension(
            length=self.length - other.length,
            mass=self.mass - other.mass,
            time=self.time - other.time,
            current=self.current - other.current,
            temperature=self.temperature - other.temperature,
            amount=self.amount - other.amount,
            luminosity=self.luminosity - other.luminosity,
        )

    def __pow__(self, exponent: int) -> Dimension:
        """Raise dimension to a power."""
        return Dimension(
            length=self.length * exponent,
            mass=self.mass * exponent,
            time=self.time * exponent,
            current=self.current * exponent,
            temperature=self.temperature * exponent,
            amount=self.amount * exponent,
            luminosity=self.luminosity * exponent,
        )

    def is_dimensionless(self) -> bool:
        """Check if this is a dimensionless quantity."""
        return all(
            [
                self.length == 0,
                self.mass == 0,
                self.time == 0,
                self.current == 0,
                self.temperature == 0,
                self.amount == 0,
                self.luminosity == 0,
            ]
        )


class Unit:
    """Represents a unit of measurement with its dimension and conversion factor to base units."""

    def __init__(
        self,
        name: str,
        symbol: str,
        dimension: Dimension,
        to_base: float = 1.0,
        offset: float = 0.0,
    ):
        self.name = name
        self.symbol = symbol
        self.dimension = dimension
        self.to_base = to_base  # Conversion factor to base unit
        self.offset = offset  # Offset for affine conversions (e.g., Celsius)

    def convert_to_base(self, value: float) -> float:
        """Convert ``value`` expressed in this unit to the base unit,
        applying the affine offset if present (e.g., °C -> K)."""
        return (value + self.offset) * self.to_base

    def convert_from_base(self, base_value: float) -> float:
        """Convert ``base_value`` expressed in the base unit back into
        this unit, undoing the affine offset if present (e.g., K -> °C)."""
        return base_value / self.to_base - self.offset

    def is_compatible_with(self, other: "Unit") -> bool:
        """Check whether ``other`` has the same dimension as this unit."""
        return self.dimension == other.dimension

    def __repr__(self) -> str:
        return f"Unit({self.symbol!r})"

    def __eq__(self, other: Any) -> bool:
        if not isinstance(other, Unit):
            return False
        return (
            self.symbol == other.symbol
            and self.dimension == other.dimension
            and self.to_base == other.to_base
            and self.offset == other.offset
        )

    def __hash__(self) -> int:
        return hash((self.symbol, self.dimension, self.to_base, self.offset))


class IncompatibleUnitsError(Exception):
    """Raised when attempting incompatible unit operations."""

    pass


class AffineUnitArithmeticError(Exception):
    """Raised when multiplying, dividing, or powering an affine unit (e.g., °C, °F)."""

    pass


@dataclass(frozen=True, eq=False)
class Quantity:
    """A numeric value with an associated unit.

    Instances are immutable and hashable, so they can be stored in sets
    and used as dictionary keys. ``eq=False`` is used so that the
    dimension-aware :meth:`__eq__` defined below is preserved; the
    matching :meth:`__hash__` is defined explicitly.
    """

    value: float
    unit: Unit

    def __init__(self, value: Union[int, float], unit: Unit):
        object.__setattr__(self, "value", float(value))
        object.__setattr__(self, "unit", unit)

    def __repr__(self) -> str:
        return f"Quantity({self.value}, {self.unit!r})"

    def __str__(self) -> str:
        return f"{self.value} {self.unit.symbol}"

    def __format__(self, format_spec: str) -> str:
        """Format the numeric value with ``format_spec`` and append the
        unit symbol, so e.g. ``f"{q:.2f}"`` works."""
        return f"{format(self.value, format_spec)} {self.unit.symbol}"

    def __hash__(self) -> int:
        # Hash on the dimension only: __eq__ treats quantities of the same
        # dimension as comparable (e.g. 1 m == 100 cm), so anything more
        # specific could violate the hash/equality contract due to
        # floating-point conversion rounding.
        return hash(self.unit.dimension)

    def _ensure_compatible(self, other: "Quantity", message: str) -> None:
        """Raise :class:`IncompatibleUnitsError` with ``message`` if the
        dimensions of ``self`` and ``other`` differ. Centralizes the
        dimension-compatibility check shared by arithmetic and
        comparison operations."""
        if not self.unit.is_compatible_with(other.unit):
            raise IncompatibleUnitsError(message)

    # Arithmetic operations

    def __add__(self, other: Quantity) -> Quantity:
        """Add two quantities (must have compatible dimensions)."""
        if not isinstance(other, Quantity):
            raise TypeError(f"Cannot add Quantity and {type(other).__name__}")

        self._ensure_compatible(
            other,
            f"Cannot add {self.unit.symbol} and {other.unit.symbol}: "
            f"incompatible dimensions",
        )

        # Convert other to self's unit
        other_in_self_unit = other.to(self.unit)
        return Quantity(self.value + other_in_self_unit.value, self.unit)

    def __sub__(self, other: Quantity) -> Quantity:
        """Subtract two quantities (must have compatible dimensions)."""
        if not isinstance(other, Quantity):
            raise TypeError(f"Cannot subtract {type(other).__name__} from Quantity")

        self._ensure_compatible(
            other,
            f"Cannot subtract {other.unit.symbol} from {self.unit.symbol}: "
            f"incompatible dimensions",
        )

        other_in_self_unit = other.to(self.unit)
        return Quantity(self.value - other_in_self_unit.value, self.unit)

    def __mul__(self, other: Union[Quantity, int, float]) -> Quantity:
        """Multiply quantity by another quantity or scalar."""
        if isinstance(other, (int, float)):
            return Quantity(self.value * other, self.unit)

        if isinstance(other, Quantity):
            if self.unit.offset != 0 or other.unit.offset != 0:
                raise AffineUnitArithmeticError(
                    f"Cannot multiply affine units {self.unit.symbol} and "
                    f"{other.unit.symbol}; convert to an absolute (zero-offset) "
                    f"unit first."
                )
            # Convert both operands to base units, then combine
            self_base = self.unit.convert_to_base(self.value)
            other_base = other.unit.convert_to_base(other.value)
            base = self_base * other_base
            new_dimension = self.unit.dimension * other.unit.dimension

            # Use a canonical base unit for the resulting dimension if one exists
            canonical = _CANONICAL_UNITS.get(new_dimension)
            if canonical is not None:
                return Quantity(base / canonical.to_base, canonical)

            # Otherwise create a derived unit
            generated_to_base = self.unit.to_base * other.unit.to_base
            new_symbol = f"{self.unit.symbol}·{other.unit.symbol}"
            new_unit = Unit(
                name=f"{self.unit.name} {other.unit.name}",
                symbol=new_symbol,
                dimension=new_dimension,
                to_base=generated_to_base,
            )
            return Quantity(base / generated_to_base, new_unit)

        raise TypeError(f"Cannot multiply Quantity and {type(other).__name__}")

    def __rmul__(self, other: Union[int, float]) -> Quantity:
        """Right multiplication (scalar * quantity)."""
        return self.__mul__(other)

    def __truediv__(self, other: Union[Quantity, int, float]) -> Quantity:
        """Divide quantity by another quantity or scalar."""
        if isinstance(other, (int, float)):
            return Quantity(self.value / other, self.unit)

        if isinstance(other, Quantity):
            if self.unit.offset != 0 or other.unit.offset != 0:
                raise AffineUnitArithmeticError(
                    f"Cannot divide affine units {self.unit.symbol} and "
                    f"{other.unit.symbol}; convert to an absolute (zero-offset) "
                    f"unit first."
                )
            # Convert both operands to base units, then combine
            self_base = self.unit.convert_to_base(self.value)
            other_base = other.unit.convert_to_base(other.value)
            base = self_base / other_base
            new_dimension = self.unit.dimension / other.unit.dimension

            # Use a canonical base unit for the resulting dimension if one exists
            canonical = _CANONICAL_UNITS.get(new_dimension)
            if canonical is not None:
                return Quantity(base / canonical.to_base, canonical)

            # Otherwise create a derived unit
            generated_to_base = self.unit.to_base / other.unit.to_base
            new_symbol = f"{self.unit.symbol}/{other.unit.symbol}"
            new_unit = Unit(
                name=f"{self.unit.name} per {other.unit.name}",
                symbol=new_symbol,
                dimension=new_dimension,
                to_base=generated_to_base,
            )
            return Quantity(base / generated_to_base, new_unit)

        raise TypeError(f"Cannot divide Quantity by {type(other).__name__}")

    def __rtruediv__(self, other: Union[int, float]) -> Quantity:
        """Right division (scalar / quantity), producing the reciprocal
        dimension (e.g. ``1 / Quantity(2, units.meter)`` -> 0.5 m^-1)."""
        if isinstance(other, (int, float)):
            return Quantity(other, units.dimensionless) / self
        return NotImplemented

    def __pow__(self, exponent: Union[int, float]) -> Quantity:
        """Raise quantity to a power."""
        if not isinstance(exponent, (int, float)):
            raise TypeError(f"Exponent must be numeric, not {type(exponent).__name__}")

        if self.unit.offset != 0:
            raise AffineUnitArithmeticError(
                f"Cannot raise affine unit {self.unit.symbol} to a power; "
                f"convert to an absolute (zero-offset) unit first."
            )

        # For fractional exponents, we need to be careful
        # This is a simplified implementation
        if not isinstance(exponent, int):
            raise NotImplementedError("Fractional exponents not yet supported")

        # Convert to base units before raising to a power
        base = self.unit.convert_to_base(self.value) ** exponent
        new_dimension = self.unit.dimension**exponent

        # Use a canonical base unit for the resulting dimension if one exists
        canonical = _CANONICAL_UNITS.get(new_dimension)
        if canonical is not None:
            return Quantity(base / canonical.to_base, canonical)

        # Otherwise create a derived unit
        generated_to_base = self.unit.to_base**exponent
        new_symbol = f"{self.unit.symbol}^{exponent}"
        new_unit = Unit(
            name=f"{self.unit.name} to the power {exponent}",
            symbol=new_symbol,
            dimension=new_dimension,
            to_base=generated_to_base,
        )
        return Quantity(base / generated_to_base, new_unit)

    def __neg__(self) -> Quantity:
        """Negate quantity."""
        return Quantity(-self.value, self.unit)

    def __abs__(self) -> Quantity:
        """Absolute value."""
        return Quantity(abs(self.value), self.unit)

    # Comparison operations

    def __eq__(self, other: Any) -> bool:
        """Strict equality: exact value comparison after conversion to a
        common unit. Use :meth:`is_close` for tolerance-based comparison."""
        if not isinstance(other, Quantity):
            return False

        if not self.unit.is_compatible_with(other.unit):
            return False

        other_in_self_unit = other.to(self.unit)
        return self.value == other_in_self_unit.value

    def is_close(
        self,
        other: Quantity,
        rel_tol: float = 1e-9,
        abs_tol: float = 0.0,
    ) -> bool:
        """Check approximate equality with another quantity.

        Returns False for incompatible dimensions. Otherwise converts
        ``other`` to this quantity's unit and delegates to
        :func:`math.isclose` with the given relative and absolute
        tolerances.
        """
        if not isinstance(other, Quantity):
            return False

        if not self.unit.is_compatible_with(other.unit):
            return False

        other_in_self_unit = other.to(self.unit)
        return math.isclose(
            self.value,
            other_in_self_unit.value,
            rel_tol=rel_tol,
            abs_tol=abs_tol,
        )

    def __lt__(self, other: Quantity) -> bool:
        if not isinstance(other, Quantity):
            raise TypeError(f"Cannot compare Quantity and {type(other).__name__}")

        self._ensure_compatible(
            other,
            f"Cannot compare {self.unit.symbol} and {other.unit.symbol}",
        )

        other_in_self_unit = other.to(self.unit)
        return self.value < other_in_self_unit.value

    def __le__(self, other: Quantity) -> bool:
        return self == other or self < other

    def __gt__(self, other: Quantity) -> bool:
        return not self <= other

    def __ge__(self, other: Quantity) -> bool:
        return not self < other

    # Unit conversion

    def to(self, target_unit: Unit) -> Quantity:
        """Convert to another unit (must have compatible dimensions)."""
        if not self.unit.is_compatible_with(target_unit):
            raise IncompatibleUnitsError(
                f"Cannot convert {self.unit.symbol} to {target_unit.symbol}: "
                f"incompatible dimensions"
            )

        # Go through the base unit; convert_to_base/convert_from_base
        # handle both linear and affine (e.g., temperature) conversions.
        base_value = self.unit.convert_to_base(self.value)
        new_value = target_unit.convert_from_base(base_value)

        return Quantity(new_value, target_unit)

    def value_in(self, target_unit: Unit) -> float:
        """Convenience alias returning the numeric value in ``target_unit``
        (equivalent to ``self.to(target_unit).value``)."""
        return self.to(target_unit).value


# ============================================================================
# Unit Definitions
# ============================================================================


class units:
    """Namespace for predefined units."""

    # Dimensionless
    dimensionless = Unit("dimensionless", "", Dimension())

    # Angle (dimensionless in SI)
    radian = Unit("radian", "rad", Dimension())
    degree = Unit("degree", "°", Dimension(), to_base=math.pi / 180.0)
    arcminute = Unit("arcminute", "′", Dimension(), to_base=math.pi / 10800.0)
    arcsecond = Unit("arcsecond", "″", Dimension(), to_base=math.pi / 648000.0)

    # Length
    meter = Unit("meter", "m", Dimension(length=1))
    kilometer = Unit("kilometer", "km", Dimension(length=1), to_base=1000.0)
    centimeter = Unit("centimeter", "cm", Dimension(length=1), to_base=0.01)
    millimeter = Unit("millimeter", "mm", Dimension(length=1), to_base=0.001)

    # Imperial length
    inch = Unit("inch", "in", Dimension(length=1), to_base=0.0254)
    foot = Unit("foot", "ft", Dimension(length=1), to_base=0.3048)
    yard = Unit("yard", "yd", Dimension(length=1), to_base=0.9144)
    mile = Unit("mile", "mi", Dimension(length=1), to_base=1609.344)

    # Mass
    kilogram = Unit("kilogram", "kg", Dimension(mass=1))
    gram = Unit("gram", "g", Dimension(mass=1), to_base=0.001)
    milligram = Unit("milligram", "mg", Dimension(mass=1), to_base=1e-6)
    tonne = Unit("tonne", "t", Dimension(mass=1), to_base=1000.0)

    # Imperial mass
    pound = Unit("pound", "lb", Dimension(mass=1), to_base=0.453592)
    ounce = Unit("ounce", "oz", Dimension(mass=1), to_base=0.0283495)

    # Time
    second = Unit("second", "s", Dimension(time=1))
    minute = Unit("minute", "min", Dimension(time=1), to_base=60.0)
    hour = Unit("hour", "h", Dimension(time=1), to_base=3600.0)
    day = Unit("day", "d", Dimension(time=1), to_base=86400.0)

    # Temperature (absolute)
    kelvin = Unit("kelvin", "K", Dimension(temperature=1))
    celsius = Unit(
        "celsius", "°C", Dimension(temperature=1), to_base=1.0, offset=273.15
    )
    fahrenheit = Unit(
        "fahrenheit", "°F", Dimension(temperature=1), to_base=5 / 9, offset=459.67
    )

    # Current
    ampere = Unit("ampere", "A", Dimension(current=1))
    milliampere = Unit("milliampere", "mA", Dimension(current=1), to_base=0.001)

    # Chemistry
    mole = Unit("mole", "mol", Dimension(amount=1))

    # Derived units

    # Force (kg·m/s²)
    newton = Unit("newton", "N", Dimension(mass=1, length=1, time=-2))

    # Energy (kg·m²/s²)
    joule = Unit("joule", "J", Dimension(mass=1, length=2, time=-2))
    kilojoule = Unit(
        "kilojoule", "kJ", Dimension(mass=1, length=2, time=-2), to_base=1000.0
    )
    calorie = Unit(
        "calorie", "cal", Dimension(mass=1, length=2, time=-2), to_base=4.184
    )
    kilocalorie = Unit(
        "kilocalorie", "kcal", Dimension(mass=1, length=2, time=-2), to_base=4184.0
    )
    watt_hour = Unit(
        "watt hour", "Wh", Dimension(mass=1, length=2, time=-2), to_base=3600.0
    )

    # Power (kg·m²/s³)
    watt = Unit("watt", "W", Dimension(mass=1, length=2, time=-3))
    kilowatt = Unit(
        "kilowatt", "kW", Dimension(mass=1, length=2, time=-3), to_base=1000.0
    )

    # Electricity
    volt = Unit("volt", "V", Dimension(mass=1, length=2, time=-3, current=-1))
    ohm = Unit("ohm", "Ω", Dimension(mass=1, length=2, time=-3, current=-2))

    # Pressure (kg/(m·s²))
    pascal = Unit("pascal", "Pa", Dimension(mass=1, length=-1, time=-2))
    kilopascal = Unit(
        "kilopascal", "kPa", Dimension(mass=1, length=-1, time=-2), to_base=1000.0
    )

    # Velocity (m/s)
    meter_per_second = Unit("meter per second", "m/s", Dimension(length=1, time=-1))
    kilometer_per_hour = Unit(
        "kilometer per hour",
        "km/h",
        Dimension(length=1, time=-1),
        to_base=1000.0 / 3600.0,
    )
    mile_per_hour = Unit(
        "mile per hour",
        "mph",
        Dimension(length=1, time=-1),
        to_base=1609.344 / 3600.0,
    )

    # Acceleration (m/s²)
    meter_per_second_squared = Unit(
        "meter per second squared", "m/s²", Dimension(length=1, time=-2)
    )

    # Area and volume
    square_meter = Unit("square meter", "m²", Dimension(length=2))
    square_kilometer = Unit("square kilometer", "km²", Dimension(length=2), to_base=1e6)
    hectare = Unit("hectare", "ha", Dimension(length=2), to_base=1e4)
    cubic_meter = Unit("cubic meter", "m³", Dimension(length=3))
    liter = Unit("liter", "L", Dimension(length=3), to_base=0.001)
    milliliter = Unit("milliliter", "mL", Dimension(length=3), to_base=1e-6)

    # Frequency
    hertz = Unit("hertz", "Hz", Dimension(time=-1))
    kilohertz = Unit("kilohertz", "kHz", Dimension(time=-1), to_base=1000.0)
    megahertz = Unit("megahertz", "MHz", Dimension(time=-1), to_base=1e6)


# Canonical SI/base units keyed by dimension (built from every zero-offset,
# to_base==1.0 unit). The first matching unit in `units` wins, so the
# dimensionless unit stays canonical for Dimension() instead of being
# overwritten by e.g. radian (which is also dimensionless with to_base==1).
_CANONICAL_UNITS: Dict[Dimension, Unit] = {}
for _u in vars(units).values():
    if isinstance(_u, Unit) and _u.to_base == 1.0 and _u.offset == 0.0:
        _CANONICAL_UNITS.setdefault(_u.dimension, _u)
del _u
