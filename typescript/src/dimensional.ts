/**
 * Unit-Aware Numeric Arithmetic Library
 *
 * A type-safe dimensional arithmetic library that tracks units at runtime
 * and prevents invalid operations.
 *
 * Zero dependencies, production-ready.
 */

/**
 * Represents the dimensional formula of a unit (e.g., L^1 T^-2 for acceleration).
 */
export class Dimension {
  constructor(
    public readonly length: number = 0,
    public readonly mass: number = 0,
    public readonly time: number = 0,
    public readonly current: number = 0,
    public readonly temperature: number = 0,
    public readonly amount: number = 0,
    public readonly luminosity: number = 0
  ) {}

  /**
   * Multiply dimensions (add exponents).
   */
  multiply(other: Dimension): Dimension {
    return new Dimension(
      this.length + other.length,
      this.mass + other.mass,
      this.time + other.time,
      this.current + other.current,
      this.temperature + other.temperature,
      this.amount + other.amount,
      this.luminosity + other.luminosity
    );
  }

  /**
   * Divide dimensions (subtract exponents).
   */
  divide(other: Dimension): Dimension {
    return new Dimension(
      this.length - other.length,
      this.mass - other.mass,
      this.time - other.time,
      this.current - other.current,
      this.temperature - other.temperature,
      this.amount - other.amount,
      this.luminosity - other.luminosity
    );
  }

  /**
   * Raise dimension to a power.
   */
  power(exponent: number): Dimension {
    return new Dimension(
      this.length * exponent,
      this.mass * exponent,
      this.time * exponent,
      this.current * exponent,
      this.temperature * exponent,
      this.amount * exponent,
      this.luminosity * exponent
    );
  }

  /**
   * Check if this is a dimensionless quantity.
   */
  isDimensionless(): boolean {
    return (
      this.length === 0 &&
      this.mass === 0 &&
      this.time === 0 &&
      this.current === 0 &&
      this.temperature === 0 &&
      this.amount === 0 &&
      this.luminosity === 0
    );
  }

  /**
   * Check equality with another dimension.
   */
  equals(other: Dimension): boolean {
    return (
      this.length === other.length &&
      this.mass === other.mass &&
      this.time === other.time &&
      this.current === other.current &&
      this.temperature === other.temperature &&
      this.amount === other.amount &&
      this.luminosity === other.luminosity
    );
  }

  /**
   * Check whether this dimension is compatible with another.
   */
  isCompatibleWith(other: Dimension): boolean {
    return this.equals(other);
  }

  /**
   * Return a string key for use in maps (comma-separated exponents).
   */
  toKey(): string {
    return [
      this.length,
      this.mass,
      this.time,
      this.current,
      this.temperature,
      this.amount,
      this.luminosity,
    ].join(',');
  }
}

/**
 * Represents a unit of measurement with its dimension and conversion factor to base units.
 */
export class Unit {
  constructor(
    public readonly name: string,
    public readonly symbol: string,
    public readonly dimension: Dimension,
    public readonly toBase: number = 1.0,
    public readonly offset: number = 0.0
  ) {}

  equals(other: Unit): boolean {
    return (
      this.symbol === other.symbol &&
      this.dimension.equals(other.dimension) &&
      this.toBase === other.toBase &&
      this.offset === other.offset
    );
  }

  /**
   * Convert a value in this unit to the base unit (handles affine units).
   */
  toBaseValue(value: number): number {
    return (value + this.offset) * this.toBase;
  }

  /**
   * Convert a value from the base unit to this unit.
   */
  fromBaseValue(value: number): number {
    return value / this.toBase - this.offset;
  }
}

/**
 * Error thrown when attempting incompatible unit operations.
 */
export class IncompatibleUnitsError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'IncompatibleUnitsError';
  }
}

/**
 * Error thrown when multiplying, dividing, or powering an affine unit (e.g., °C, °F).
 */
export class AffineUnitArithmeticError extends Error {
  constructor(message: string) {
    super(message);
    this.name = 'AffineUnitArithmeticError';
  }
}

/**
 * A numeric value with an associated unit.
 */
export class Quantity {
  constructor(
    public readonly value: number,
    public readonly unit: Unit
  ) {}

  /**
   * Throw IncompatibleUnitsError if the target does not have a compatible dimension.
   */
  private ensureCompatible(target: Unit | Quantity, message: string): void {
    const dimension = target instanceof Quantity ? target.unit.dimension : target.dimension;
    if (!this.unit.dimension.isCompatibleWith(dimension)) {
      throw new IncompatibleUnitsError(message);
    }
  }

  toString(): string {
    return `${this.value} ${this.unit.symbol}`;
  }

  // Arithmetic operations

  /**
   * Add two quantities (must have compatible dimensions).
   */
  add(other: Quantity): Quantity {
    this.ensureCompatible(
      other,
      `Cannot add ${this.unit.symbol} and ${other.unit.symbol}: incompatible dimensions`
    );

    const otherInSelfUnit = other.to(this.unit);
    return new Quantity(this.value + otherInSelfUnit.value, this.unit);
  }

  /**
   * Subtract two quantities (must have compatible dimensions).
   */
  subtract(other: Quantity): Quantity {
    this.ensureCompatible(
      other,
      `Cannot subtract ${other.unit.symbol} from ${this.unit.symbol}: incompatible dimensions`
    );

    const otherInSelfUnit = other.to(this.unit);
    return new Quantity(this.value - otherInSelfUnit.value, this.unit);
  }

  /**
   * Multiply quantity by another quantity or scalar.
   */
  multiply(other: number): Quantity;
  multiply(other: Quantity): Quantity;
  multiply(other: Quantity | number): Quantity {
    if (typeof other === 'number') {
      return new Quantity(this.value * other, this.unit);
    }

    if (this.unit.offset !== 0 || other.unit.offset !== 0) {
      throw new AffineUnitArithmeticError(
        `Cannot multiply affine units ${this.unit.symbol} and ${other.unit.symbol}; ` +
        `convert to an absolute (zero-offset) unit first.`
      );
    }

    // Convert both operands to base units, then combine
    const base = this.unit.toBaseValue(this.value) * other.unit.toBaseValue(other.value);
    const newDimension = this.unit.dimension.multiply(other.unit.dimension);

    const canonical = CANONICAL_UNITS[newDimension.toKey()];
    if (canonical !== undefined) {
      return new Quantity(canonical.fromBaseValue(base), canonical);
    }

    // Create a derived unit
    const generatedToBase = this.unit.toBase * other.unit.toBase;
    const newSymbol = `${this.unit.symbol}·${other.unit.symbol}`;
    const newUnit = new Unit(
      `${this.unit.name} ${other.unit.name}`,
      newSymbol,
      newDimension,
      generatedToBase
    );

    return new Quantity(newUnit.fromBaseValue(base), newUnit);
  }

  /**
   * Divide quantity by another quantity or scalar.
   */
  divide(other: number): Quantity;
  divide(other: Quantity): Quantity;
  divide(other: Quantity | number): Quantity {
    if (typeof other === 'number') {
      return new Quantity(this.value / other, this.unit);
    }

    if (this.unit.offset !== 0 || other.unit.offset !== 0) {
      throw new AffineUnitArithmeticError(
        `Cannot divide affine units ${this.unit.symbol} and ${other.unit.symbol}; ` +
        `convert to an absolute (zero-offset) unit first.`
      );
    }

    // Convert both operands to base units, then combine
    const base = this.unit.toBaseValue(this.value) / other.unit.toBaseValue(other.value);
    const newDimension = this.unit.dimension.divide(other.unit.dimension);

    const canonical = CANONICAL_UNITS[newDimension.toKey()];
    if (canonical !== undefined) {
      return new Quantity(canonical.fromBaseValue(base), canonical);
    }

    // Create a derived unit
    const generatedToBase = this.unit.toBase / other.unit.toBase;
    const newSymbol = `${this.unit.symbol}/${other.unit.symbol}`;
    const newUnit = new Unit(
      `${this.unit.name} per ${other.unit.name}`,
      newSymbol,
      newDimension,
      generatedToBase
    );

    return new Quantity(newUnit.fromBaseValue(base), newUnit);
  }

  /**
   * Raise quantity to a power.
   */
  power(exponent: number): Quantity {
    if (this.unit.offset !== 0) {
      throw new AffineUnitArithmeticError(
        `Cannot raise affine unit ${this.unit.symbol} to a power; ` +
        `convert to an absolute (zero-offset) unit first.`
      );
    }

    if (!Number.isInteger(exponent)) {
      throw new Error('Fractional exponents not yet supported');
    }

    // Convert to base units before raising to a power
    const base = Math.pow(this.unit.toBaseValue(this.value), exponent);
    const newDimension = this.unit.dimension.power(exponent);

    const canonical = CANONICAL_UNITS[newDimension.toKey()];
    if (canonical !== undefined) {
      return new Quantity(canonical.fromBaseValue(base), canonical);
    }

    const generatedToBase = Math.pow(this.unit.toBase, exponent);
    const newSymbol = `${this.unit.symbol}^${exponent}`;
    const newUnit = new Unit(
      `${this.unit.name} to the power ${exponent}`,
      newSymbol,
      newDimension,
      generatedToBase
    );

    return new Quantity(newUnit.fromBaseValue(base), newUnit);
  }

  /**
   * Negate quantity.
   */
  negate(): Quantity {
    return new Quantity(-this.value, this.unit);
  }

  /**
   * Absolute value.
   */
  abs(): Quantity {
    return new Quantity(Math.abs(this.value), this.unit);
  }

  // Comparison operations

  /**
   * Check strict equality with another quantity: exact value comparison
   * after conversion to a common unit. Use `isClose` for tolerance-based
   * comparison.
   */
  equals(other: Quantity): boolean {
    if (!this.unit.dimension.equals(other.unit.dimension)) {
      return false;
    }

    const otherInSelfUnit = other.to(this.unit);
    return this.value === otherInSelfUnit.value;
  }

  /**
   * Check approximate equality with another quantity.
   *
   * Returns false for incompatible dimensions. Otherwise converts `other`
   * to this quantity's unit and applies relative/absolute tolerances
   * (same semantics as Python's `math.isclose`):
   * `diff <= absTol || diff <= relTol * max(|a|, |b|)`.
   */
  isClose(other: Quantity, relTol: number = 1e-9, absTol: number = 0.0): boolean {
    if (!this.unit.dimension.equals(other.unit.dimension)) {
      return false;
    }

    const otherInSelfUnit = other.to(this.unit);
    if (this.value === otherInSelfUnit.value) {
      return true;
    }
    const diff = Math.abs(this.value - otherInSelfUnit.value);
    return (
      diff <= absTol ||
      diff <= relTol * Math.max(Math.abs(this.value), Math.abs(otherInSelfUnit.value))
    );
  }

  /**
   * Check if less than another quantity.
   */
  lessThan(other: Quantity): boolean {
    this.ensureCompatible(
      other,
      `Cannot compare ${this.unit.symbol} and ${other.unit.symbol}`
    );

    const otherInSelfUnit = other.to(this.unit);
    return this.value < otherInSelfUnit.value;
  }

  /**
   * Check if less than or equal to another quantity.
   */
  lessThanOrEqual(other: Quantity): boolean {
    return this.equals(other) || this.lessThan(other);
  }

  /**
   * Check if greater than another quantity.
   */
  greaterThan(other: Quantity): boolean {
    return !this.lessThanOrEqual(other);
  }

  /**
   * Check if greater than or equal to another quantity.
   */
  greaterThanOrEqual(other: Quantity): boolean {
    return !this.lessThan(other);
  }

  // Unit conversion

  /**
   * Check whether another quantity has a compatible dimension
   * (i.e., can be added, subtracted, compared, or converted).
   */
  isCompatibleWith(other: Quantity): boolean {
    return this.unit.dimension.isCompatibleWith(other.unit.dimension);
  }

  /**
   * Get this quantity's numeric value expressed in a target unit.
   * Throws `IncompatibleUnitsError` for incompatible dimensions.
   */
  in(targetUnit: Unit): number {
    return this.to(targetUnit).value;
  }

  /**
   * Convert to another unit (must have compatible dimensions).
   */
  to(targetUnit: Unit): Quantity {
    this.ensureCompatible(
      targetUnit,
      `Cannot convert ${this.unit.symbol} to ${targetUnit.symbol}: incompatible dimensions`
    );

    const baseValue = this.unit.toBaseValue(this.value);
    const newValue = targetUnit.fromBaseValue(baseValue);

    return new Quantity(newValue, targetUnit);
  }
}

/**
 * Convenience factory for creating a `Quantity`.
 *
 * Example: `const q = quantity(5, units.meter);`
 */
export function quantity(value: number, unit: Unit): Quantity {
  return new Quantity(value, unit);
}

// ============================================================================
// Unit Definitions
// ============================================================================

/**
 * Namespace for predefined units.
 */
export const units = {
  // Dimensionless
  dimensionless: new Unit('dimensionless', '', new Dimension()),

  // Angle (dimensionless)
  radian: new Unit('radian', 'rad', new Dimension()),
  degree: new Unit('degree', '°', new Dimension(), Math.PI / 180.0),
  arcminute: new Unit('arcminute', '′', new Dimension(), Math.PI / 10800.0),
  arcsecond: new Unit('arcsecond', '″', new Dimension(), Math.PI / 648000.0),

  // Length
  meter: new Unit('meter', 'm', new Dimension(1)),
  kilometer: new Unit('kilometer', 'km', new Dimension(1), 1000.0),
  centimeter: new Unit('centimeter', 'cm', new Dimension(1), 0.01),
  millimeter: new Unit('millimeter', 'mm', new Dimension(1), 0.001),

  // Imperial length
  inch: new Unit('inch', 'in', new Dimension(1), 0.0254),
  foot: new Unit('foot', 'ft', new Dimension(1), 0.3048),
  yard: new Unit('yard', 'yd', new Dimension(1), 0.9144),
  mile: new Unit('mile', 'mi', new Dimension(1), 1609.344),

  // Mass
  kilogram: new Unit('kilogram', 'kg', new Dimension(0, 1)),
  gram: new Unit('gram', 'g', new Dimension(0, 1), 0.001),
  milligram: new Unit('milligram', 'mg', new Dimension(0, 1), 1e-6),
  tonne: new Unit('tonne', 't', new Dimension(0, 1), 1000.0),

  // Imperial mass
  pound: new Unit('pound', 'lb', new Dimension(0, 1), 0.453592),
  ounce: new Unit('ounce', 'oz', new Dimension(0, 1), 0.0283495),

  // Time
  second: new Unit('second', 's', new Dimension(0, 0, 1)),
  minute: new Unit('minute', 'min', new Dimension(0, 0, 1), 60.0),
  hour: new Unit('hour', 'h', new Dimension(0, 0, 1), 3600.0),
  day: new Unit('day', 'd', new Dimension(0, 0, 1), 86400.0),

  // Temperature (absolute)
  kelvin: new Unit('kelvin', 'K', new Dimension(0, 0, 0, 0, 1)),
  celsius: new Unit('celsius', '°C', new Dimension(0, 0, 0, 0, 1), 1.0, 273.15),
  fahrenheit: new Unit('fahrenheit', '°F', new Dimension(0, 0, 0, 0, 1), 5 / 9, 459.67),

  // Current
  ampere: new Unit('ampere', 'A', new Dimension(0, 0, 0, 1)),
  milliampere: new Unit('milliampere', 'mA', new Dimension(0, 0, 0, 1), 0.001),

  // Amount of substance
  mole: new Unit('mole', 'mol', new Dimension(0, 0, 0, 0, 0, 1)),

  // Luminous intensity
  candela: new Unit('candela', 'cd', new Dimension(0, 0, 0, 0, 0, 0, 1)),

  // Derived units

  // Force (kg·m/s²)
  newton: new Unit('newton', 'N', new Dimension(1, 1, -2)),

  // Energy (kg·m²/s²)
  joule: new Unit('joule', 'J', new Dimension(2, 1, -2)),
  kilojoule: new Unit('kilojoule', 'kJ', new Dimension(2, 1, -2), 1000.0),
  calorie: new Unit('calorie', 'cal', new Dimension(2, 1, -2), 4.184),
  kilocalorie: new Unit('kilocalorie', 'kcal', new Dimension(2, 1, -2), 4184.0),
  wattHour: new Unit('watt hour', 'Wh', new Dimension(2, 1, -2), 3600.0),

  // Power (kg·m²/s³)
  watt: new Unit('watt', 'W', new Dimension(2, 1, -3)),
  kilowatt: new Unit('kilowatt', 'kW', new Dimension(2, 1, -3), 1000.0),

  // Pressure (kg/(m·s²))
  pascal: new Unit('pascal', 'Pa', new Dimension(-1, 1, -2)),
  kilopascal: new Unit('kilopascal', 'kPa', new Dimension(-1, 1, -2), 1000.0),

  // Electricity
  volt: new Unit('volt', 'V', new Dimension(2, 1, -3, -1)),
  ohm: new Unit('ohm', 'Ω', new Dimension(2, 1, -3, -2)),
  coulomb: new Unit('coulomb', 'C', new Dimension(0, 0, 1, 1)),

  // Velocity (m/s)
  meterPerSecond: new Unit('meter per second', 'm/s', new Dimension(1, 0, -1)),
  kilometerPerHour: new Unit('kilometer per hour', 'km/h', new Dimension(1, 0, -1), 1000.0 / 3600.0),
  milePerHour: new Unit('mile per hour', 'mph', new Dimension(1, 0, -1), 1609.344 / 3600.0),

  // Acceleration (m/s²)
  meterPerSecondSquared: new Unit('meter per second squared', 'm/s²', new Dimension(1, 0, -2)),

  // Area and volume
  squareMeter: new Unit('square meter', 'm²', new Dimension(2)),
  squareKilometer: new Unit('square kilometer', 'km²', new Dimension(2), 1e6),
  hectare: new Unit('hectare', 'ha', new Dimension(2), 1e4),
  cubicMeter: new Unit('cubic meter', 'm³', new Dimension(3)),
  liter: new Unit('liter', 'L', new Dimension(3), 1e-3),
  milliliter: new Unit('milliliter', 'mL', new Dimension(3), 1e-6),

  // Frequency
  hertz: new Unit('hertz', 'Hz', new Dimension(0, 0, -1)),
  kilohertz: new Unit('kilohertz', 'kHz', new Dimension(0, 0, -1), 1e3),
  megahertz: new Unit('megahertz', 'MHz', new Dimension(0, 0, -1), 1e6),
} as const;

// Canonical SI/base units keyed by dimension (built from every zero-offset, toBase==1.0 unit)
const CANONICAL_UNITS: Record<string, Unit> = Object.values(units)
  .filter((u): u is Unit => u instanceof Unit && u.toBase === 1.0 && u.offset === 0.0)
  .reduce((acc, u) => {
    acc[u.dimension.toKey()] = u;
    return acc;
  }, {} as Record<string, Unit>);

// Preserve dimensionless as the canonical dimensionless unit (radian also has toBase 1)
CANONICAL_UNITS[new Dimension().toKey()] = units.dimensionless;
