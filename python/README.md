# Unit-Aware Arithmetic (Python)

A type-safe dimensional arithmetic library that tracks units at runtime and prevents invalid operations.

## Features

- ✅ **Type-safe dimensional analysis**: Prevents incompatible unit operations
- ✅ **Zero dependencies**: Core functionality has no external dependencies
- ✅ **Comprehensive unit coverage**: SI, imperial, and derived units
- ✅ **Arithmetic operations**: Add, subtract, multiply, divide with unit tracking
- ✅ **Unit conversion**: Automatic and explicit conversion between compatible units
- ✅ **Clear error messages**: Helpful errors for incompatible operations
- ✅ **Production-ready**: Comprehensive test coverage (>95%)

## Installation

```bash
pip install -e .
```

## Quick Start

```python
from dimensional import Quantity, units

# Create quantities with units
distance = Quantity(100, units.meter)
time = Quantity(9.58, units.second)

# Arithmetic operations with automatic unit tracking
speed = distance / time
print(speed)  # 10.438413361169102 m/s

# Unit conversion
distance_km = distance.to(units.kilometer)
print(distance_km)  # 0.1 km

# Type-safe operations - this will raise an error!
try:
    distance + time  # IncompatibleUnitsError
except Exception as e:
    print(f"Error: {e}")
```

## Usage Examples

### Basic Arithmetic

```python
from dimensional import Quantity, units

# Addition (same dimension required)
d1 = Quantity(5, units.meter)
d2 = Quantity(3, units.meter)
total = d1 + d2  # 8.0 m

# Multiplication creates derived units
area = Quantity(5, units.meter) * Quantity(3, units.meter)
print(area)  # 15.0 m·m

# Division creates derived units
velocity = Quantity(100, units.meter) / Quantity(10, units.second)
print(velocity)  # 10.0 m/s
```

### Unit Conversion

```python
from dimensional import Quantity, units

# Length conversion
distance = Quantity(1, units.mile)
distance_km = distance.to(units.kilometer)
print(distance_km)  # 1.60934 km

# Temperature conversion
temp_c = Quantity(0, units.celsius)
temp_k = temp_c.to(units.kelvin)
print(temp_k)  # 273.15 K

# Mass conversion
mass_lb = Quantity(10, units.pound)
mass_kg = mass_lb.to(units.kilogram)
print(mass_kg)  # 4.53592 kg
```

### Physics Calculations

```python
from dimensional import Quantity, units

# Calculate velocity
distance = Quantity(100, units.meter)
time = Quantity(9.58, units.second)
velocity = distance / time
print(f"Velocity: {velocity}")

# Calculate force (F = ma)
mass = Quantity(10, units.kilogram)
acceleration = Quantity(9.8, units.meter_per_second_squared)
force = mass * acceleration
print(f"Force: {force}")

# Calculate kinetic energy (KE = 1/2 * m * v²)
mass = Quantity(2, units.kilogram)
velocity = Quantity(10, units.meter_per_second)
ke = 0.5 * mass * (velocity ** 2)
print(f"Kinetic Energy: {ke}")
```

### Comparison Operations

```python
from dimensional import Quantity, units

d1 = Quantity(100, units.centimeter)
d2 = Quantity(1, units.meter)

# Automatic conversion for comparison
print(d1 == d2)  # True
print(d1 < Quantity(2, units.meter))  # True

# `==` is strict exact equality (after conversion to a common unit).
# For approximate comparison within a tolerance, use `is_close`:
print(Quantity(1, units.meter).is_close(Quantity(1.0000001, units.meter), rel_tol=1e-6))  # True
print(Quantity(0, units.meter).is_close(Quantity(1e-12, units.meter), abs_tol=1e-9))       # True
```

## Available Units

### Length
- `meter`, `kilometer`, `centimeter`, `millimeter`
- `inch`, `foot`, `yard`, `mile`

### Mass
- `kilogram`, `gram`, `milligram`, `tonne`
- `pound`, `ounce`

### Time
- `second`, `minute`, `hour`, `day`

### Temperature
- `kelvin`, `celsius`, `fahrenheit`

### Current
- `ampere`, `milliampere`

### Chemistry
- `mole`

### Derived Units
- `newton` (force)
- `joule`, `kilojoule` (energy)
- `calorie`, `kilocalorie`, `watt_hour` (energy)
- `watt`, `kilowatt` (power)
- `pascal`, `kilopascal` (pressure)
- `meter_per_second`, `kilometer_per_hour`, `mile_per_hour` (velocity)
- `meter_per_second_squared` (acceleration)
- `square_meter`, `square_kilometer`, `hectare` (area)
- `cubic_meter`, `liter`, `milliliter` (volume)
- `hertz`, `kilohertz`, `megahertz` (frequency)
- `radian`, `degree`, `arcminute`, `arcsecond` (angle)
- `volt`, `ohm` (electricity)

## Error Handling

The library provides clear error messages for invalid operations:

```python
from dimensional import Quantity, units, IncompatibleUnitsError, AffineUnitArithmeticError

distance = Quantity(100, units.meter)
time = Quantity(10, units.second)

try:
    # This will raise IncompatibleUnitsError
    result = distance + time
except IncompatibleUnitsError as e:
    print(e)  # "Cannot add m and s: incompatible dimensions"

try:
    # This will raise AffineUnitArithmeticError
    result = Quantity(2, units.celsius) * Quantity(3, units.celsius)
except AffineUnitArithmeticError as e:
    print(e)  # "Cannot multiply affine units °C and °C; ..."
```

## Testing

Run the test suite:

```bash
pytest test_dimensional.py -v
```

Run with coverage:

```bash
pytest test_dimensional.py --cov=dimensional --cov-report=html
```

## Type Checking

This library includes type hints. Run type checking with:

```bash
mypy dimensional.py
```

## Design Principles

1. **Zero Dependencies**: Core functionality has no external dependencies
2. **Type Safety**: Prevents invalid operations at runtime
3. **Clear Errors**: Actionable error messages
4. **Production Ready**: Comprehensive test coverage
5. **Performance**: Efficient implementation with minimal overhead

## API Reference

### `Dimension`

Represents the dimensional formula of a unit (e.g., L^1 T^-2 for acceleration).

### `Unit`

Represents a unit of measurement with its dimension and conversion factor.

### `Quantity`

A numeric value with an associated unit. Supports:
- Arithmetic: `+`, `-`, `*`, `/`, `**`, `-` (negation), `abs()`
- Comparison: `==`, `!=`, `<`, `<=`, `>`, `>=` (`==` is strict exact equality after conversion to a common unit)
- Approximate comparison: `.is_close(other, rel_tol=1e-9, abs_tol=0.0)`
- Conversion: `.to(target_unit)`

### `units`

Namespace containing all predefined units.

## License

MIT License

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- Code coverage remains >90%
- Type hints are included
- Documentation is updated

## Changelog

### 2.0.0 (Unreleased)
- Added angle units: `radian`, `degree`, `arcminute`, `arcsecond`
- Added frequency units: `kilohertz`, `megahertz`
- Added area units: `square_kilometer`, `hectare`
- Added volume units: `liter`, `milliliter`
- Added velocity unit: `mile_per_hour`
- Added chemistry unit: `mole`
- Added energy units: `calorie`, `kilocalorie`, `watt_hour`
- Added electricity units: `volt`, `ohm`
- Updated canonical-unit lookup table to canonicalize Hz, m², m³, V, and Ω

### 1.0.0 (2026-08-28)
- Initial release
- Support for SI and imperial units
- Comprehensive dimensional analysis
- Temperature conversion with offset handling
- Full test coverage
