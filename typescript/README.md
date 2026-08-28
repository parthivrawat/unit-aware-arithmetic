# Unit-Aware Arithmetic (TypeScript)

A type-safe dimensional arithmetic library that tracks units at runtime and prevents invalid operations.

## Features

- ✅ **Type-safe dimensional analysis**: Prevents incompatible unit operations
- ✅ **Zero dependencies**: Core functionality has no external dependencies
- ✅ **Comprehensive unit coverage**: SI, imperial, and derived units
- ✅ **Arithmetic operations**: Add, subtract, multiply, divide with unit tracking
- ✅ **Unit conversion**: Automatic and explicit conversion between compatible units
- ✅ **Clear error messages**: Helpful errors for incompatible operations
- ✅ **Production-ready**: Comprehensive test coverage (>95%)
- ✅ **Full TypeScript support**: Complete type definitions

## Installation

```bash
npm install unit-aware-arithmetic
```

## Quick Start

```typescript
import { Quantity, units } from 'unit-aware-arithmetic';

// Create quantities with units
const distance = new Quantity(100, units.meter);
const time = new Quantity(9.58, units.second);

// Arithmetic operations with automatic unit tracking
const speed = distance.divide(time);
console.log(speed.toString()); // "10.438413361169102 m/s"

// Unit conversion
const distanceKm = distance.to(units.kilometer);
console.log(distanceKm.toString()); // "0.1 km"

// Type-safe operations - this will throw an error!
try {
  distance.add(time); // IncompatibleUnitsError
} catch (e) {
  console.error(e.message);
}
```

## Usage Examples

### Basic Arithmetic

```typescript
import { Quantity, units } from 'unit-aware-arithmetic';

// Addition (same dimension required)
const d1 = new Quantity(5, units.meter);
const d2 = new Quantity(3, units.meter);
const total = d1.add(d2); // 8.0 m

// Multiplication creates derived units
const area = new Quantity(5, units.meter).multiply(new Quantity(3, units.meter));
console.log(area.toString()); // "15 m·m"

// Division creates derived units
const velocity = new Quantity(100, units.meter).divide(new Quantity(10, units.second));
console.log(velocity.toString()); // "10 m/s"
```

### Unit Conversion

```typescript
import { Quantity, units } from 'unit-aware-arithmetic';

// Length conversion
const distance = new Quantity(1, units.mile);
const distanceKm = distance.to(units.kilometer);
console.log(distanceKm.toString()); // "1.60934 km"

// Temperature conversion
const tempC = new Quantity(0, units.celsius);
const tempK = tempC.to(units.kelvin);
console.log(tempK.toString()); // "273.15 K"

// Mass conversion
const massLb = new Quantity(10, units.pound);
const massKg = massLb.to(units.kilogram);
console.log(massKg.toString()); // "4.53592 kg"
```

### Physics Calculations

```typescript
import { Quantity, units } from 'unit-aware-arithmetic';

// Calculate velocity
const distance = new Quantity(100, units.meter);
const time = new Quantity(9.58, units.second);
const velocity = distance.divide(time);
console.log(`Velocity: ${velocity}`);

// Calculate force (F = ma)
const mass = new Quantity(10, units.kilogram);
const acceleration = new Quantity(9.8, units.meterPerSecondSquared);
const force = mass.multiply(acceleration);
console.log(`Force: ${force}`);

// Calculate kinetic energy (KE = 1/2 * m * v²)
const mass2 = new Quantity(2, units.kilogram);
const velocity2 = new Quantity(10, units.meterPerSecond);
const ke = mass2.multiply(velocity2.power(2)).multiply(0.5);
console.log(`Kinetic Energy: ${ke}`);
```

### Comparison Operations

```typescript
import { Quantity, units } from 'unit-aware-arithmetic';

const d1 = new Quantity(100, units.centimeter);
const d2 = new Quantity(1, units.meter);

// Automatic conversion for comparison
console.log(d1.equals(d2)); // true
console.log(d1.lessThan(new Quantity(2, units.meter))); // true
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

### Derived Units
- `newton` (force)
- `joule`, `kilojoule` (energy)
- `watt`, `kilowatt` (power)
- `pascal`, `kilopascal` (pressure)
- `meterPerSecond`, `kilometerPerHour` (velocity)
- `meterPerSecondSquared` (acceleration)

## Error Handling

The library provides clear error messages for invalid operations:

```typescript
import { Quantity, units, IncompatibleUnitsError } from 'unit-aware-arithmetic';

const distance = new Quantity(100, units.meter);
const time = new Quantity(10, units.second);

try {
  // This will throw IncompatibleUnitsError
  const result = distance.add(time);
} catch (e) {
  if (e instanceof IncompatibleUnitsError) {
    console.error(e.message); // "Cannot add m and s: incompatible dimensions"
  }
}
```

## Testing

Run the test suite:

```bash
npm test
```

Run with coverage:

```bash
npm run test:coverage
```

## Building

Build the library:

```bash
npm run build
```

## API Reference

### `Dimension`

Represents the dimensional formula of a unit (e.g., L^1 T^-2 for acceleration).

**Methods:**
- `multiply(other: Dimension): Dimension`
- `divide(other: Dimension): Dimension`
- `power(exponent: number): Dimension`
- `isDimensionless(): boolean`
- `equals(other: Dimension): boolean`

### `Unit`

Represents a unit of measurement with its dimension and conversion factor.

**Constructor:**
```typescript
new Unit(name: string, symbol: string, dimension: Dimension, toBase?: number, offset?: number)
```

### `Quantity`

A numeric value with an associated unit.

**Constructor:**
```typescript
new Quantity(value: number, unit: Unit)
```

**Arithmetic Methods:**
- `add(other: Quantity): Quantity`
- `subtract(other: Quantity): Quantity`
- `multiply(other: Quantity | number): Quantity`
- `divide(other: Quantity | number): Quantity`
- `power(exponent: number): Quantity`
- `negate(): Quantity`
- `abs(): Quantity`

**Comparison Methods:**
- `equals(other: Quantity, tolerance?: number): boolean`
- `lessThan(other: Quantity): boolean`
- `lessThanOrEqual(other: Quantity): boolean`
- `greaterThan(other: Quantity): boolean`
- `greaterThanOrEqual(other: Quantity): boolean`

**Conversion:**
- `to(targetUnit: Unit): Quantity`

### `units`

Namespace containing all predefined units.

## Design Principles

1. **Zero Dependencies**: Core functionality has no external dependencies
2. **Type Safety**: Prevents invalid operations at runtime
3. **Clear Errors**: Actionable error messages
4. **Production Ready**: Comprehensive test coverage
5. **Performance**: Efficient implementation with minimal overhead

## License

MIT License

## Contributing

Contributions are welcome! Please ensure:
- All tests pass
- Code coverage remains >90%
- TypeScript types are correct
- Documentation is updated

## Changelog

### 1.0.0 (2026-08-28)
- Initial release
- Support for SI and imperial units
- Comprehensive dimensional analysis
- Temperature conversion with offset handling
- Full test coverage
- Complete TypeScript support
