# Unit-Aware Arithmetic

A type-safe dimensional arithmetic library that tracks units at compile/runtime and prevents invalid operations.

**Status**: ✅ Production-ready across Python, TypeScript, Go, and Rust  
**Languages**: Python, TypeScript, Go, Rust  
**License**: MIT

## Overview

This library solves a critical problem in scientific and engineering software: preventing unit-conversion bugs through dimensional analysis. It tracks physical dimensions (length, mass, time, etc.) and prevents incompatible operations like adding meters to seconds.

### The Problem

```python
# Without unit tracking - compiles but wrong!
distance = 100  # meters
time = 9.58     # seconds
result = distance + time  # 109.58 what? This is nonsense!
```

### The Solution

```python
from dimensional import Quantity, units

distance = Quantity(100, units.meter)
time = Quantity(9.58, units.second)
result = distance + time  # IncompatibleUnitsError: Cannot add m and s
```

## Features

- ✅ **Type-safe dimensional analysis**: Prevents incompatible unit operations
- ✅ **Zero dependencies**: Core functionality has no external dependencies
- ✅ **Comprehensive unit coverage**: SI, imperial, and derived units
- ✅ **Arithmetic operations**: Add, subtract, multiply, divide with automatic unit tracking
- ✅ **Unit conversion**: Automatic and explicit conversion between compatible units
- ✅ **Clear error messages**: Helpful errors for incompatible operations
- ✅ **Production-ready**: Comprehensive test coverage (>95%)
- ✅ **Cross-language consistency**: Similar APIs across Python, TypeScript, Go, and Rust

## Quick Start

### Python

```python
from dimensional import Quantity, units

# Create quantities
distance = Quantity(100, units.meter)
time = Quantity(9.58, units.second)

# Automatic unit tracking
velocity = distance / time  # 10.438 m/s

# Unit conversion
distance_km = distance.to(units.kilometer)  # 0.1 km

# Type safety
distance + time  # IncompatibleUnitsError!
```

### TypeScript

```typescript
import { Quantity, units } from 'unit-aware-arithmetic';

// Create quantities
const distance = new Quantity(100, units.meter);
const time = new Quantity(9.58, units.second);

// Automatic unit tracking
const velocity = distance.divide(time); // 10.438 m/s

// Unit conversion
const distanceKm = distance.to(units.kilometer); // 0.1 km

// Type safety
distance.add(time); // IncompatibleUnitsError!
```

## Installation

### Python

```bash
cd python
pip install -e .
```

### TypeScript

```bash
cd typescript
npm install
```

### Go

```bash
cd go
go test -v
```

### Rust

```bash
cd rust
cargo test
```

## Package Names

The packages have different distribution names versus the in-code module/crate names:

- **Python**: package `unit-aware-arithmetic` on PyPI; importable module is `dimensional` (`from dimensional import Quantity, units`).
- **TypeScript**: npm package `unit-aware-arithmetic`.
- **Go**: module path `github.com/parthivrawat/unit-aware-arithmetic/go/v2` (v2 semantic import versioning); package name `dimensional`.
- **Rust**: crate `unit-aware-arithmetic` on crates.io.

## Cross-Language API

A quick reference for the same operation in each implementation. All examples use meters/seconds where applicable.

| Operation | Python | TypeScript | Go | Rust |
|---|---|---|---|---|
| Construction | `Quantity(value, unit)` | `quantity(value, unit)` or `new Quantity(value, unit)` | `dimensional.NewQuantity(value, unit)` | `Quantity::new(value, units::METER)` |
| Add | `a + b` | `a.add(b)` | `a.Add(b)` returns `(Quantity, error)` | `a + b` (panics if incompatible; also `a.try_add(&b) -> Result`) |
| Subtract | `a - b` | `a.subtract(b)` | `a.Subtract(b)` returns `(Quantity, error)` | `a - b` (panics if incompatible; also `a.try_sub(&b) -> Result`) |
| Multiply | `a * b` / `a * scalar` | `a.multiply(b)` / `a.multiply(scalar)` | `a.Multiply(b)` returns `(Quantity, error)`; `a.MultiplyScalar(f)` | `a * b` / `a * scalar` |
| Divide | `a / b` / `a / scalar` | `a.divide(b)` / `a.divide(scalar)` | `a.Divide(b)` returns `(Quantity, error)`; `a.DivideScalar(f)` | `a / b` / `a / scalar` |
| Power | `a ** exponent` | `a.power(exponent)` | `a.Power(exponent)` returns `(Quantity, error)` | `a.pow(exponent)` |
| Convert | `a.to(unit)` (Quantity); `a.value_in(unit)` (float) | `a.to(unit)` (Quantity); `a.in(unit)` (number) | `a.To(unit)` returns `(Quantity, error)` | `a.to(unit)` returns `Result<Quantity, IncompatibleUnitsError>` |
| Equality / Compare | `a == b` (or `a.is_close(b)`) | `a.equals(b)` (or `a.isClose(b, relTol, absTol)`) | `a.Equals(b, tol)`; `a.IsClose(b, relTol, absTol)`; `a.LessThan(b)` / `a.GreaterThan(b)` / `a.LessThanOrEqual(b, tol)` / `a.GreaterThanOrEqual(b, tol)` | `a == b` / `a.partial_cmp(&b)`; `a.is_close(&b, relTol, absTol)` / `a.approx_eq(&b, tol)` |

## Real-World Examples

### Physics Calculations

```python
# Calculate force (F = ma)
mass = Quantity(10, units.kilogram)
acceleration = Quantity(9.8, units.meter_per_second_squared)
force = mass * acceleration  # 98 kg·m/s² (Newtons)

# Calculate kinetic energy (KE = 1/2 * m * v²)
mass = Quantity(2, units.kilogram)
velocity = Quantity(10, units.meter_per_second)
ke = 0.5 * mass * (velocity ** 2)  # 100 kg·m²/s² (Joules)
```

### Engineering Calculations

```python
# Convert imperial to metric
height_ft = Quantity(6, units.foot)
height_m = height_ft.to(units.meter)  # 1.8288 m

# Temperature conversion
temp_f = Quantity(72, units.fahrenheit)
temp_c = temp_f.to(units.celsius)  # 22.22 °C
```

### Financial Calculations

```python
# Prevent mixing currencies (future feature)
# price_usd = Quantity(100, units.usd)
# price_eur = Quantity(90, units.eur)
# total = price_usd + price_eur  # IncompatibleUnitsError!
```

## Available Units

### Base Units
- **Length**: meter, kilometer, centimeter, millimeter, inch, foot, yard, mile
- **Mass**: kilogram, gram, milligram, tonne, pound, ounce
- **Time**: second, minute, hour, day
- **Temperature**: kelvin, celsius, fahrenheit
- **Current**: ampere, milliampere
- **Amount of substance**: mole
- **Angle**: radian, degree, arcminute, arcsecond

### Derived Units
- **Force**: newton (kg·m/s²)
- **Energy**: joule, kilojoule, calorie, kilocalorie, watt_hour (kg·m²/s²)
- **Power**: watt, kilowatt (kg·m²/s³)
- **Pressure**: pascal, kilopascal (kg/(m·s²))
- **Frequency**: hertz, kilohertz, megahertz (1/s)
- **Area**: square_meter, square_kilometer, hectare
- **Volume**: cubic_meter, liter, milliliter
- **Velocity**: meter_per_second, kilometer_per_hour, mile_per_hour
- **Acceleration**: meter_per_second_squared
- **Electricity**: volt (kg·m²/(s³·A)), ohm (kg·m²/(s³·A²))

## Implementation Status

| Language   | Status | Tests | Coverage | Notes |
|------------|--------|-------|----------|-------|
| Python     | ✅ Production-ready | 104 tests | >95% | Stable and fully tested |
| TypeScript | ✅ Production-ready | 100 tests | >95% | Stable and fully tested |
| Go         | ✅ Production-ready | 80 tests | >95% | Typed APIs, error returns, immutable `Quantity` |
| Rust       | ✅ Production-ready | 65 tests + 1 doc-test | >95% | No leaks, `PartialEq`/`PartialOrd`, clippy-clean |

## Design Principles

1. **Zero Dependencies**: Core functionality has no external dependencies
2. **Type Safety**: Prevents invalid operations at compile/runtime
3. **Clear Errors**: Actionable error messages with context
4. **Production Ready**: Comprehensive test coverage and real-world validation
5. **Performance**: Efficient implementation with minimal overhead
6. **Cross-Language Consistency**: Similar APIs across all implementations

## Architecture

### Core Concepts

1. **Dimension**: Represents the dimensional formula (e.g., L^1 T^-2 for acceleration)
2. **Unit**: A specific measurement unit with its dimension and conversion factor
3. **Quantity**: A numeric value with an associated unit

### Dimensional Analysis

The library uses the SI base dimensions:
- Length (L)
- Mass (M)
- Time (T)
- Electric Current (I)
- Temperature (Θ)
- Amount of Substance (N)
- Luminous Intensity (J)

All derived units are expressed as combinations of these base dimensions.

## Testing

### Python

```bash
cd python
pytest test_dimensional.py -v
pytest test_dimensional.py --cov=dimensional --cov-report=html
```

### TypeScript

```bash
cd typescript
npm test
npm run test:coverage
```

### Go

```bash
cd go
go test -v ./...
go test -bench=. -benchmem   # benchmarks
```

### Rust

```bash
cd rust
cargo test
cargo clippy --all-targets -- -D warnings
```

## Performance

Every operation performs a dimension check and (for add/subtract/compare) a
unit conversion, so `Quantity` arithmetic is meaningfully slower than raw
numeric operations. Measure with the reproducible benchmark suite in
[`benchmarks/`](./benchmarks/) on your own hardware — results below are from a
single run on an Intel i5-8250U and are **machine-dependent**:

| Operation | Python | TypeScript | Go | Rust |
|---|---|---|---|---|
| Construction | ~8x | ~5x | ~17x | ~10x |
| Addition | ~90x | ~3x | ~70x | ~95x |
| Multiplication | ~420x | ~25x | ~80x | ~95x |
| Division | ~160x | ~40x | ~110x | ~90x |
| Conversion | ~60x | ~2x | ~35x | ~20x |

(overhead vs. the equivalent raw numeric operation; absolute times are tens to
hundreds of nanoseconds per operation — see `benchmarks/` for details)

- **Memory**: Negligible overhead (one object per quantity)

For performance-critical code, convert to raw numbers after validation:

```python
# Validate with units
distance = Quantity(100, units.meter)
time = Quantity(10, units.second)

# Extract raw values for tight loops
d_val = distance.value
t_val = time.value
for i in range(1000000):
    v = d_val / t_val  # Fast!
```

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- Code coverage remains >90%
- Type hints/annotations are included
- Documentation is updated
- Examples are provided

## Roadmap

### Version 1.0 (Current)
- ✅ Python implementation
- ✅ TypeScript implementation
- ✅ Go implementation
- ✅ Rust implementation
- ✅ Core dimensional analysis
- ✅ SI and imperial units
- ✅ Temperature conversion with offsets
- ✅ Cross-language consistency

### Version 1.1 (Planned)
- 🚧 Additional units (angle, frequency, etc.)
- 🚧 Custom unit definitions

### Version 2.0 (Future)
- Currency units with exchange rates
- Uncertainty propagation
- Vector quantities
- Complex numbers with units
- Performance optimizations

## FAQ

### Q: Why not use existing libraries like Pint (Python) or mathjs (JS)?

**A**: Existing libraries are excellent but have limitations:
- **Pint**: Heavy dependencies, complex API, Python-only
- **mathjs**: Large bundle size, not type-safe, JavaScript-only
- **Our library**: Zero dependencies, cross-language, production-ready

### Q: How does this compare to static type systems?

**A**: Static type systems (like F#'s units of measure) are ideal but not available in Python/TypeScript/Go. This library provides runtime checking with minimal overhead.

### Q: Can I define custom units?

**A**: Yes! Create a new `Unit` instance:

```python
from dimensional import Unit, Dimension

# Define a custom unit
furlong = Unit("furlong", "fur", Dimension(length=1), to_base=201.168)
```

### Q: What about performance?

**A**: `Quantity` operations cost tens to hundreds of nanoseconds — roughly 5–400x a raw numeric operation depending on language and operation, because each call performs dimensional checks and unit conversion. See [`benchmarks/`](./benchmarks/) for reproducible numbers. For performance-critical code, extract raw values after validation.

## License

MIT License - see LICENSE file for details.

## Authors

- Parthiv Rawat <parthiv05022000@gmail.com>

## Acknowledgments

- Inspired by F#'s units of measure
- Based on SI dimensional analysis principles
- Thanks to the scientific Python and TypeScript communities

## Links

- [Python Implementation](./python/)
- [TypeScript Implementation](./typescript/)
- [Domain-Specific Proposals](../../../36-custom-library-proposals/DOMAIN_LIBRARY_PROPOSALS.md)
- [Main Implementations](../../README.md)

---

**Last Updated**: 2026-09-06  
**Version**: 2.0.0  
**Status**: Production-ready across Python, TypeScript, Go, and Rust
