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
    
    length: int = 0      # L
    mass: int = 0        # M
    time: int = 0        # T
    current: int = 0     # I
    temperature: int = 0 # Θ
    amount: int = 0      # N
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
        return all([
            self.length == 0,
            self.mass == 0,
            self.time == 0,
            self.current == 0,
            self.temperature == 0,
            self.amount == 0,
            self.luminosity == 0,
        ])


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
        self.offset = offset    # Offset for affine conversions (e.g., Celsius)
    
    def __repr__(self) -> str:
        return f"Unit({self.symbol})"
    
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


class Quantity:
    """A numeric value with an associated unit."""
    
    def __init__(self, value: Union[int, float], unit: Unit):
        self.value = float(value)
        self.unit = unit
    
    def __repr__(self) -> str:
        return f"{self.value} {self.unit.symbol}"
    
    def __str__(self) -> str:
        return self.__repr__()
    
    # Arithmetic operations
    
    def __add__(self, other: Quantity) -> Quantity:
        """Add two quantities (must have compatible dimensions)."""
        if not isinstance(other, Quantity):
            raise TypeError(f"Cannot add Quantity and {type(other).__name__}")
        
        if self.unit.dimension != other.unit.dimension:
            raise IncompatibleUnitsError(
                f"Cannot add {self.unit.symbol} and {other.unit.symbol}: "
                f"incompatible dimensions"
            )
        
        # Convert other to self's unit
        other_in_self_unit = other.to(self.unit)
        return Quantity(self.value + other_in_self_unit.value, self.unit)
    
    def __sub__(self, other: Quantity) -> Quantity:
        """Subtract two quantities (must have compatible dimensions)."""
        if not isinstance(other, Quantity):
            raise TypeError(f"Cannot subtract {type(other).__name__} from Quantity")
        
        if self.unit.dimension != other.unit.dimension:
            raise IncompatibleUnitsError(
                f"Cannot subtract {other.unit.symbol} from {self.unit.symbol}: "
                f"incompatible dimensions"
            )
        
        other_in_self_unit = other.to(self.unit)
        return Quantity(self.value - other_in_self_unit.value, self.unit)
    
    def __mul__(self, other: Union[Quantity, int, float]) -> Quantity:
        """Multiply quantity by another quantity or scalar."""
        if isinstance(other, (int, float)):
            return Quantity(self.value * other, self.unit)
        
        if isinstance(other, Quantity):
            # Multiply values and dimensions
            new_value = self.value * other.value
            new_dimension = self.unit.dimension * other.unit.dimension
            
            # Create a derived unit
            new_symbol = f"{self.unit.symbol}·{other.unit.symbol}"
            new_unit = Unit(
                name=f"{self.unit.name} {other.unit.name}",
                symbol=new_symbol,
                dimension=new_dimension,
                to_base=self.unit.to_base * other.unit.to_base,
            )
            return Quantity(new_value, new_unit)
        
        raise TypeError(f"Cannot multiply Quantity and {type(other).__name__}")
    
    def __rmul__(self, other: Union[int, float]) -> Quantity:
        """Right multiplication (scalar * quantity)."""
        return self.__mul__(other)
    
    def __truediv__(self, other: Union[Quantity, int, float]) -> Quantity:
        """Divide quantity by another quantity or scalar."""
        if isinstance(other, (int, float)):
            return Quantity(self.value / other, self.unit)
        
        if isinstance(other, Quantity):
            # Divide values and dimensions
            new_value = self.value / other.value
            new_dimension = self.unit.dimension / other.unit.dimension
            
            # Create a derived unit
            new_symbol = f"{self.unit.symbol}/{other.unit.symbol}"
            new_unit = Unit(
                name=f"{self.unit.name} per {other.unit.name}",
                symbol=new_symbol,
                dimension=new_dimension,
                to_base=self.unit.to_base / other.unit.to_base,
            )
            return Quantity(new_value, new_unit)
        
        raise TypeError(f"Cannot divide Quantity by {type(other).__name__}")
    
    def __pow__(self, exponent: Union[int, float]) -> Quantity:
        """Raise quantity to a power."""
        if not isinstance(exponent, (int, float)):
            raise TypeError(f"Exponent must be numeric, not {type(exponent).__name__}")
        
        new_value = self.value ** exponent
        
        # For integer exponents, we can compute the exact dimension
        if isinstance(exponent, int):
            new_dimension = self.unit.dimension ** exponent
            new_symbol = f"{self.unit.symbol}^{exponent}"
            new_unit = Unit(
                name=f"{self.unit.name} to the power {exponent}",
                symbol=new_symbol,
                dimension=new_dimension,
                to_base=self.unit.to_base ** exponent,
            )
            return Quantity(new_value, new_unit)
        
        # For fractional exponents, we need to be careful
        # This is a simplified implementation
        raise NotImplementedError("Fractional exponents not yet supported")
    
    def __neg__(self) -> Quantity:
        """Negate quantity."""
        return Quantity(-self.value, self.unit)
    
    def __abs__(self) -> Quantity:
        """Absolute value."""
        return Quantity(abs(self.value), self.unit)
    
    # Comparison operations
    
    def __eq__(self, other: Any) -> bool:
        if not isinstance(other, Quantity):
            return False
        
        if self.unit.dimension != other.unit.dimension:
            return False
        
        other_in_self_unit = other.to(self.unit)
        return math.isclose(self.value, other_in_self_unit.value)
    
    def __lt__(self, other: Quantity) -> bool:
        if not isinstance(other, Quantity):
            raise TypeError(f"Cannot compare Quantity and {type(other).__name__}")
        
        if self.unit.dimension != other.unit.dimension:
            raise IncompatibleUnitsError(
                f"Cannot compare {self.unit.symbol} and {other.unit.symbol}"
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
        if self.unit.dimension != target_unit.dimension:
            raise IncompatibleUnitsError(
                f"Cannot convert {self.unit.symbol} to {target_unit.symbol}: "
                f"incompatible dimensions"
            )
        
        # Handle affine conversions (e.g., temperature)
        if self.unit.offset != 0 or target_unit.offset != 0:
            # Convert to base unit first (remove offset)
            base_value = (self.value + self.unit.offset) * self.unit.to_base
            # Convert from base to target (apply offset)
            new_value = base_value / target_unit.to_base - target_unit.offset
        else:
            # Simple linear conversion
            new_value = self.value * (self.unit.to_base / target_unit.to_base)
        
        return Quantity(new_value, target_unit)


# ============================================================================
# Unit Definitions
# ============================================================================

class units:
    """Namespace for predefined units."""
    
    # Dimensionless
    dimensionless = Unit("dimensionless", "", Dimension())
    
    # Length
    meter = Unit("meter", "m", Dimension(length=1))
    kilometer = Unit("kilometer", "km", Dimension(length=1), to_base=1000.0)
    centimeter = Unit("centimeter", "cm", Dimension(length=1), to_base=0.01)
    millimeter = Unit("millimeter", "mm", Dimension(length=1), to_base=0.001)
    
    # Imperial length
    inch = Unit("inch", "in", Dimension(length=1), to_base=0.0254)
    foot = Unit("foot", "ft", Dimension(length=1), to_base=0.3048)
    yard = Unit("yard", "yd", Dimension(length=1), to_base=0.9144)
    mile = Unit("mile", "mi", Dimension(length=1), to_base=1609.34)
    
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
    celsius = Unit("celsius", "°C", Dimension(temperature=1), to_base=1.0, offset=273.15)
    fahrenheit = Unit("fahrenheit", "°F", Dimension(temperature=1), to_base=5/9, offset=459.67)
    
    # Current
    ampere = Unit("ampere", "A", Dimension(current=1))
    milliampere = Unit("milliampere", "mA", Dimension(current=1), to_base=0.001)
    
    # Derived units
    
    # Force (kg·m/s²)
    newton = Unit("newton", "N", Dimension(mass=1, length=1, time=-2))
    
    # Energy (kg·m²/s²)
    joule = Unit("joule", "J", Dimension(mass=1, length=2, time=-2))
    kilojoule = Unit("kilojoule", "kJ", Dimension(mass=1, length=2, time=-2), to_base=1000.0)
    
    # Power (kg·m²/s³)
    watt = Unit("watt", "W", Dimension(mass=1, length=2, time=-3))
    kilowatt = Unit("kilowatt", "kW", Dimension(mass=1, length=2, time=-3), to_base=1000.0)
    
    # Pressure (kg/(m·s²))
    pascal = Unit("pascal", "Pa", Dimension(mass=1, length=-1, time=-2))
    kilopascal = Unit("kilopascal", "kPa", Dimension(mass=1, length=-1, time=-2), to_base=1000.0)
    
    # Velocity (m/s)
    meter_per_second = Unit("meter per second", "m/s", Dimension(length=1, time=-1))
    kilometer_per_hour = Unit("kilometer per hour", "km/h", Dimension(length=1, time=-1), to_base=1000.0/3600.0)
    
    # Acceleration (m/s²)
    meter_per_second_squared = Unit("meter per second squared", "m/s²", Dimension(length=1, time=-2))
