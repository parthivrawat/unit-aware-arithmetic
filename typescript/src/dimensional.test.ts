/**
 * Tests for Unit-Aware Numeric Arithmetic Library
 */

import { describe, it, expect } from 'vitest';
import { Dimension, Unit, Quantity, units, IncompatibleUnitsError } from './dimensional';

describe('Dimension', () => {
  it('should create a dimension', () => {
    const d = new Dimension(1, 0, -2);
    expect(d.length).toBe(1);
    expect(d.time).toBe(-2);
    expect(d.mass).toBe(0);
  });

  it('should multiply dimensions', () => {
    const d1 = new Dimension(1); // L
    const d2 = new Dimension(0, 0, -1); // T^-1
    const result = d1.multiply(d2); // L·T^-1
    expect(result.length).toBe(1);
    expect(result.time).toBe(-1);
  });

  it('should divide dimensions', () => {
    const d1 = new Dimension(1, 0, -1); // L·T^-1
    const d2 = new Dimension(0, 0, 1); // T
    const result = d1.divide(d2); // L·T^-2
    expect(result.length).toBe(1);
    expect(result.time).toBe(-2);
  });

  it('should raise dimension to a power', () => {
    const d = new Dimension(1, 0, -1); // L·T^-1
    const result = d.power(2); // L²·T^-2
    expect(result.length).toBe(2);
    expect(result.time).toBe(-2);
  });

  it('should check if dimensionless', () => {
    expect(new Dimension().isDimensionless()).toBe(true);
    expect(new Dimension(1).isDimensionless()).toBe(false);
  });
});

describe('Unit', () => {
  it('should create a unit', () => {
    const u = new Unit('meter', 'm', new Dimension(1));
    expect(u.name).toBe('meter');
    expect(u.symbol).toBe('m');
    expect(u.dimension.length).toBe(1);
  });

  it('should check equality', () => {
    const u1 = new Unit('meter', 'm', new Dimension(1));
    const u2 = new Unit('meter', 'm', new Dimension(1));
    expect(u1.equals(u2)).toBe(true);
  });

  it('should have predefined units', () => {
    expect(units.meter.symbol).toBe('m');
    expect(units.kilogram.symbol).toBe('kg');
    expect(units.second.symbol).toBe('s');
  });
});

describe('Quantity - Basics', () => {
  it('should create a quantity', () => {
    const q = new Quantity(5.0, units.meter);
    expect(q.value).toBe(5.0);
    expect(q.unit).toBe(units.meter);
  });

  it('should convert to string', () => {
    const q = new Quantity(5.0, units.meter);
    expect(q.toString()).toBe('5 m');
  });
});

describe('Arithmetic', () => {
  it('should add quantities with same unit', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(3, units.meter);
    const result = q1.add(q2);
    expect(result.value).toBe(8);
    expect(result.unit).toBe(units.meter);
  });

  it('should add quantities with compatible units', () => {
    const q1 = new Quantity(1, units.meter);
    const q2 = new Quantity(100, units.centimeter);
    const result = q1.add(q2);
    expect(result.value).toBeCloseTo(2.0);
    expect(result.unit).toBe(units.meter);
  });

  it('should throw on adding incompatible units', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(3, units.second);
    expect(() => q1.add(q2)).toThrow(IncompatibleUnitsError);
  });

  it('should subtract quantities with same unit', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(3, units.meter);
    const result = q1.subtract(q2);
    expect(result.value).toBe(2);
    expect(result.unit).toBe(units.meter);
  });

  it('should subtract quantities with compatible units', () => {
    const q1 = new Quantity(1, units.kilometer);
    const q2 = new Quantity(500, units.meter);
    const result = q1.subtract(q2);
    expect(result.value).toBeCloseTo(0.5);
    expect(result.unit).toBe(units.kilometer);
  });

  it('should throw on subtracting incompatible units', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(3, units.kilogram);
    expect(() => q1.subtract(q2)).toThrow(IncompatibleUnitsError);
  });

  it('should multiply by scalar', () => {
    const q = new Quantity(5, units.meter);
    const result = q.multiply(3);
    expect(result.value).toBe(15);
    expect(result.unit).toBe(units.meter);
  });

  it('should multiply quantities', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(3, units.meter);
    const result = q1.multiply(q2);
    expect(result.value).toBe(15);
    expect(result.unit.dimension.length).toBe(2); // m²
  });

  it('should divide by scalar', () => {
    const q = new Quantity(10, units.meter);
    const result = q.divide(2);
    expect(result.value).toBe(5);
    expect(result.unit).toBe(units.meter);
  });

  it('should divide quantities', () => {
    const q1 = new Quantity(100, units.meter);
    const q2 = new Quantity(10, units.second);
    const result = q1.divide(q2);
    expect(result.value).toBe(10);
    expect(result.unit.dimension.length).toBe(1);
    expect(result.unit.dimension.time).toBe(-1);
  });

  it('should raise to integer power', () => {
    const q = new Quantity(3, units.meter);
    const result = q.power(2);
    expect(result.value).toBe(9);
    expect(result.unit.dimension.length).toBe(2);
  });

  it('should negate quantity', () => {
    const q = new Quantity(5, units.meter);
    const result = q.negate();
    expect(result.value).toBe(-5);
    expect(result.unit).toBe(units.meter);
  });

  it('should get absolute value', () => {
    const q = new Quantity(-5, units.meter);
    const result = q.abs();
    expect(result.value).toBe(5);
    expect(result.unit).toBe(units.meter);
  });
});

describe('Comparison', () => {
  it('should check equality with same unit', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(5, units.meter);
    expect(q1.equals(q2)).toBe(true);
  });

  it('should check equality with different compatible units', () => {
    const q1 = new Quantity(1, units.meter);
    const q2 = new Quantity(100, units.centimeter);
    expect(q1.equals(q2)).toBe(true);
  });

  it('should check inequality', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(3, units.meter);
    expect(q1.equals(q2)).toBe(false);
  });

  it('should check less than', () => {
    const q1 = new Quantity(3, units.meter);
    const q2 = new Quantity(5, units.meter);
    expect(q1.lessThan(q2)).toBe(true);
  });

  it('should check less than or equal', () => {
    const q1 = new Quantity(3, units.meter);
    const q2 = new Quantity(5, units.meter);
    const q3 = new Quantity(3, units.meter);
    expect(q1.lessThanOrEqual(q2)).toBe(true);
    expect(q1.lessThanOrEqual(q3)).toBe(true);
  });

  it('should check greater than', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(3, units.meter);
    expect(q1.greaterThan(q2)).toBe(true);
  });

  it('should check greater than or equal', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(3, units.meter);
    const q3 = new Quantity(5, units.meter);
    expect(q1.greaterThanOrEqual(q2)).toBe(true);
    expect(q1.greaterThanOrEqual(q3)).toBe(true);
  });

  it('should throw on comparing incompatible units', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(3, units.second);
    expect(() => q1.lessThan(q2)).toThrow(IncompatibleUnitsError);
  });
});

describe('Conversion', () => {
  it('should convert length', () => {
    const q = new Quantity(1, units.meter);
    const result = q.to(units.centimeter);
    expect(result.value).toBeCloseTo(100);
    expect(result.unit).toBe(units.centimeter);
  });

  it('should convert mass', () => {
    const q = new Quantity(1, units.kilogram);
    const result = q.to(units.gram);
    expect(result.value).toBeCloseTo(1000);
    expect(result.unit).toBe(units.gram);
  });

  it('should convert time', () => {
    const q = new Quantity(1, units.hour);
    const result = q.to(units.second);
    expect(result.value).toBeCloseTo(3600);
    expect(result.unit).toBe(units.second);
  });

  it('should convert imperial to metric', () => {
    const q = new Quantity(1, units.mile);
    const result = q.to(units.kilometer);
    expect(result.value).toBeCloseTo(1.60934, 5);
    expect(result.unit).toBe(units.kilometer);
  });

  it('should throw on converting incompatible units', () => {
    const q = new Quantity(5, units.meter);
    expect(() => q.to(units.second)).toThrow(IncompatibleUnitsError);
  });

  it('should convert temperature celsius to kelvin', () => {
    const q = new Quantity(0, units.celsius);
    const result = q.to(units.kelvin);
    expect(result.value).toBeCloseTo(273.15);
    expect(result.unit).toBe(units.kelvin);
  });

  it('should convert temperature fahrenheit to celsius', () => {
    const q = new Quantity(32, units.fahrenheit);
    const result = q.to(units.celsius);
    expect(result.value).toBeCloseTo(0, 1);
    expect(result.unit).toBe(units.celsius);
  });
});

describe('Physics Examples', () => {
  it('should calculate velocity', () => {
    const distance = new Quantity(100, units.meter);
    const time = new Quantity(9.58, units.second);
    const velocity = distance.divide(time);
    expect(velocity.value).toBeCloseTo(10.438, 3);
    expect(velocity.unit.dimension.length).toBe(1);
    expect(velocity.unit.dimension.time).toBe(-1);
  });

  it('should calculate acceleration', () => {
    const velocity = new Quantity(10, units.meterPerSecond);
    const time = new Quantity(2, units.second);
    const acceleration = velocity.divide(time);
    expect(acceleration.value).toBe(5);
    expect(acceleration.unit.dimension.length).toBe(1);
    expect(acceleration.unit.dimension.time).toBe(-2);
  });

  it('should calculate force', () => {
    // F = ma
    const mass = new Quantity(10, units.kilogram);
    const acceleration = new Quantity(9.8, units.meterPerSecondSquared);
    const force = mass.multiply(acceleration);
    expect(force.value).toBeCloseTo(98);
    expect(force.unit.dimension.mass).toBe(1);
    expect(force.unit.dimension.length).toBe(1);
    expect(force.unit.dimension.time).toBe(-2);
  });

  it('should calculate kinetic energy', () => {
    // KE = 1/2 * m * v²
    const mass = new Quantity(2, units.kilogram);
    const velocity = new Quantity(10, units.meterPerSecond);
    const ke = mass.multiply(velocity.power(2)).multiply(0.5);
    expect(ke.value).toBe(100);
    expect(ke.unit.dimension.mass).toBe(1);
    expect(ke.unit.dimension.length).toBe(2);
    expect(ke.unit.dimension.time).toBe(-2);
  });
});

describe('Edge Cases', () => {
  it('should handle zero quantity', () => {
    const q = new Quantity(0, units.meter);
    expect(q.value).toBe(0);
  });

  it('should handle negative quantity', () => {
    const q = new Quantity(-5, units.meter);
    expect(q.value).toBe(-5);
  });

  it('should handle very large quantity', () => {
    const q = new Quantity(1e100, units.meter);
    expect(q.value).toBe(1e100);
  });

  it('should handle very small quantity', () => {
    const q = new Quantity(1e-100, units.meter);
    expect(q.value).toBe(1e-100);
  });
});
