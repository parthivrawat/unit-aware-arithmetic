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
 * A numeric value with an associated unit.
 */
export class Quantity {
  constructor(
    public readonly value: number,
    public readonly unit: Unit
  ) {}

  toString(): string {
    return `${this.value} ${this.unit.symbol}`;
  }

  // Arithmetic operations

  /**
   * Add two quantities (must have compatible dimensions).
   */
  add(other: Quantity): Quantity {
    if (!this.unit.dimension.equals(other.unit.dimension)) {
      throw new IncompatibleUnitsError(
        `Cannot add ${this.unit.symbol} and ${other.unit.symbol}: incompatible dimensions`
      );
    }

    const otherInSelfUnit = other.to(this.unit);
    return new Quantity(this.value + otherInSelfUnit.value, this.unit);
  }

  /**
   * Subtract two quantities (must have compatible dimensions).
   */
  subtract(other: Quantity): Quantity {
    if (!this.unit.dimension.equals(other.unit.dimension)) {
      throw new IncompatibleUnitsError(
        `Cannot subtract ${other.unit.symbol} from ${this.unit.symbol}: incompatible dimensions`
      );
    }

    const otherInSelfUnit = other.to(this.unit);
    return new Quantity(this.value - otherInSelfUnit.value, this.unit);
  }

  /**
   * Multiply quantity by another quantity or scalar.
   */
  multiply(other: Quantity | number): Quantity {
    if (typeof other === 'number') {
      return new Quantity(this.value * other, this.unit);
    }

    // Multiply values and dimensions
    const newValue = this.value * other.value;
    const newDimension = this.unit.dimension.multiply(other.unit.dimension);

    // Create a derived unit
    const newSymbol = `${this.unit.symbol}·${other.unit.symbol}`;
    const newUnit = new Unit(
      `${this.unit.name} ${other.unit.name}`,
      newSymbol,
      newDimension,
      this.unit.toBase * other.unit.toBase
    );

    return new Quantity(newValue, newUnit);
  }

  /**
   * Divide quantity by another quantity or scalar.
   */
  divide(other: Quantity | number): Quantity {
    if (typeof other === 'number') {
      return new Quantity(this.value / other, this.unit);
    }

    // Divide values and dimensions
    const newValue = this.value / other.value;
    const newDimension = this.unit.dimension.divide(other.unit.dimension);

    // Create a derived unit
    const newSymbol = `${this.unit.symbol}/${other.unit.symbol}`;
    const newUnit = new Unit(
      `${this.unit.name} per ${other.unit.name}`,
      newSymbol,
      newDimension,
      this.unit.toBase / other.unit.toBase
    );

    return new Quantity(newValue, newUnit);
  }

  /**
   * Raise quantity to a power.
   */
  power(exponent: number): Quantity {
    if (!Number.isInteger(exponent)) {
      throw new Error('Fractional exponents not yet supported');
    }

    const newValue = Math.pow(this.value, exponent);
    const newDimension = this.unit.dimension.power(exponent);
    const newSymbol = `${this.unit.symbol}^${exponent}`;
    const newUnit = new Unit(
      `${this.unit.name} to the power ${exponent}`,
      newSymbol,
      newDimension,
      Math.pow(this.unit.toBase, exponent)
    );

    return new Quantity(newValue, newUnit);
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
   * Check equality with another quantity.
   */
  equals(other: Quantity, tolerance: number = 1e-9): boolean {
    if (!this.unit.dimension.equals(other.unit.dimension)) {
      return false;
    }

    const otherInSelfUnit = other.to(this.unit);
    return Math.abs(this.value - otherInSelfUnit.value) < tolerance;
  }

  /**
   * Check if less than another quantity.
   */
  lessThan(other: Quantity): boolean {
    if (!this.unit.dimension.equals(other.unit.dimension)) {
      throw new IncompatibleUnitsError(
        `Cannot compare ${this.unit.symbol} and ${other.unit.symbol}`
      );
    }

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
   * Convert to another unit (must have compatible dimensions).
   */
  to(targetUnit: Unit): Quantity {
    if (!this.unit.dimension.equals(targetUnit.dimension)) {
      throw new IncompatibleUnitsError(
        `Cannot convert ${this.unit.symbol} to ${targetUnit.symbol}: incompatible dimensions`
      );
    }

    // Handle affine conversions (e.g., temperature)
    let newValue: number;
    if (this.unit.offset !== 0 || targetUnit.offset !== 0) {
      // Convert to base unit first (remove offset)
      const baseValue = (this.value + this.unit.offset) * this.unit.toBase;
      // Convert from base to target (apply offset)
      newValue = baseValue / targetUnit.toBase - targetUnit.offset;
    } else {
      // Simple linear conversion
      newValue = this.value * (this.unit.toBase / targetUnit.toBase);
    }

    return new Quantity(newValue, targetUnit);
  }
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

  // Length
  meter: new Unit('meter', 'm', new Dimension(1)),
  kilometer: new Unit('kilometer', 'km', new Dimension(1), 1000.0),
  centimeter: new Unit('centimeter', 'cm', new Dimension(1), 0.01),
  millimeter: new Unit('millimeter', 'mm', new Dimension(1), 0.001),

  // Imperial length
  inch: new Unit('inch', 'in', new Dimension(1), 0.0254),
  foot: new Unit('foot', 'ft', new Dimension(1), 0.3048),
  yard: new Unit('yard', 'yd', new Dimension(1), 0.9144),
  mile: new Unit('mile', 'mi', new Dimension(1), 1609.34),

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

  // Derived units

  // Force (kg·m/s²)
  newton: new Unit('newton', 'N', new Dimension(1, 1, -2)),

  // Energy (kg·m²/s²)
  joule: new Unit('joule', 'J', new Dimension(2, 1, -2)),
  kilojoule: new Unit('kilojoule', 'kJ', new Dimension(2, 1, -2), 1000.0),

  // Power (kg·m²/s³)
  watt: new Unit('watt', 'W', new Dimension(2, 1, -3)),
  kilowatt: new Unit('kilowatt', 'kW', new Dimension(2, 1, -3), 1000.0),

  // Pressure (kg/(m·s²))
  pascal: new Unit('pascal', 'Pa', new Dimension(-1, 1, -2)),
  kilopascal: new Unit('kilopascal', 'kPa', new Dimension(-1, 1, -2), 1000.0),

  // Velocity (m/s)
  meterPerSecond: new Unit('meter per second', 'm/s', new Dimension(1, 0, -1)),
  kilometerPerHour: new Unit('kilometer per hour', 'km/h', new Dimension(1, 0, -1), 1000.0 / 3600.0),

  // Acceleration (m/s²)
  meterPerSecondSquared: new Unit('meter per second squared', 'm/s²', new Dimension(1, 0, -2)),
} as const;
