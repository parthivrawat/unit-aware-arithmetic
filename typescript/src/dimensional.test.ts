/**
 * Tests for Unit-Aware Numeric Arithmetic Library
 */

import { describe, it, expect } from 'vitest';
import { Dimension, Unit, Quantity, quantity, units, IncompatibleUnitsError, AffineUnitArithmeticError } from './dimensional';

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

  it('should throw on multiplying affine units', () => {
    const q1 = new Quantity(2, units.celsius);
    const q2 = new Quantity(3, units.celsius);
    expect(() => q1.multiply(q2)).toThrow(AffineUnitArithmeticError);
  });

  it('should throw on multiplying fahrenheit quantities', () => {
    const q1 = new Quantity(2, units.fahrenheit);
    const q2 = new Quantity(1, units.fahrenheit);
    expect(() => q1.multiply(q2)).toThrow(AffineUnitArithmeticError);
  });

  it('should throw on dividing by an affine unit', () => {
    const q1 = new Quantity(100, units.meter);
    const q2 = new Quantity(2, units.celsius);
    expect(() => q1.divide(q2)).toThrow(AffineUnitArithmeticError);
  });

  it('should throw on raising an affine unit to a power', () => {
    const q = new Quantity(2, units.celsius);
    expect(() => q.power(2)).toThrow(AffineUnitArithmeticError);
  });

  it('should allow scalar multiplication of affine units', () => {
    const q = new Quantity(2, units.celsius);
    const result = q.multiply(3);
    expect(result.value).toBe(6);
    expect(result.unit).toBe(units.celsius);
  });

  it('should allow scalar division of affine units', () => {
    const q = new Quantity(10, units.celsius);
    const result = q.divide(2);
    expect(result.value).toBe(5);
    expect(result.unit).toBe(units.celsius);
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

describe('isClose', () => {
  it('should be close with relative tolerance', () => {
    const q1 = new Quantity(1, units.meter);
    const q2 = new Quantity(1.0000001, units.meter);
    // diff is 1e-7, so relTol must exceed ~1e-7 to be "close"
    expect(q1.isClose(q2, 1e-9)).toBe(false);
    expect(q1.isClose(q2, 1e-6)).toBe(true);
  });

  it('should be close with absolute tolerance', () => {
    const q1 = new Quantity(1e-9, units.meter);
    const q2 = new Quantity(2e-9, units.meter);
    expect(q1.isClose(q2, 0.0, 1.5e-9)).toBe(true);
  });

  it('should be close relatively', () => {
    const q1 = new Quantity(1e-9, units.meter);
    const q2 = new Quantity(2e-9, units.meter);
    // diff is 1e-9; relative tolerance must be >= 0.5 to cover it
    expect(q1.isClose(q2, 1e-9)).toBe(false);
    expect(q1.isClose(q2, 0.6)).toBe(true);
  });

  it('should not be close for different values', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(3, units.meter);
    expect(q1.isClose(q2)).toBe(false);
  });

  it('should not be close for incompatible dimensions', () => {
    const q1 = new Quantity(5, units.meter);
    const q2 = new Quantity(5, units.second);
    expect(q1.isClose(q2)).toBe(false);
  });

  it('should be close near zero with absolute tolerance', () => {
    const q1 = new Quantity(0, units.meter);
    const q2 = new Quantity(1e-12, units.meter);
    expect(q1.isClose(q2, 1e-9, 1e-9)).toBe(true);
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

describe('Canonical unit lookup', () => {
  it('should canonicalize area from length multiplication', () => {
    const result = new Quantity(5, units.meter).multiply(new Quantity(3, units.meter));
    expect(result.value).toBeCloseTo(15);
    expect(result.unit).toBe(units.squareMeter);
  });

  it('should canonicalize area from power', () => {
    const result = new Quantity(3, units.meter).power(2);
    expect(result.value).toBeCloseTo(9);
    expect(result.unit).toBe(units.squareMeter);
  });

  it('should canonicalize velocity', () => {
    const result = new Quantity(100, units.meter).divide(new Quantity(10, units.second));
    expect(result.value).toBeCloseTo(10);
    expect(result.unit).toBe(units.meterPerSecond);
  });

  it('should canonicalize force', () => {
    const result = new Quantity(10, units.kilogram).multiply(
      new Quantity(9.8, units.meterPerSecondSquared)
    );
    expect(result.value).toBeCloseTo(98);
    expect(result.unit).toBe(units.newton);
  });

  it('should canonicalize kinetic energy', () => {
    const mass = new Quantity(2, units.kilogram);
    const velocity = new Quantity(10, units.meterPerSecond);
    const result = mass.multiply(velocity.power(2)).multiply(0.5);
    expect(result.value).toBeCloseTo(100);
    expect(result.unit).toBe(units.joule);
  });

  it('should canonicalize power from kilojoule per second', () => {
    const result = new Quantity(1, units.kilojoule).divide(new Quantity(1, units.second));
    expect(result.value).toBeCloseTo(1000);
    expect(result.unit).toBe(units.watt);
  });

  it('should canonicalize force from joule per meter', () => {
    const result = new Quantity(1, units.joule).divide(new Quantity(1, units.meter));
    expect(result.value).toBeCloseTo(1);
    expect(result.unit).toBe(units.newton);
  });

  it('should canonicalize energy from force times length', () => {
    const result = new Quantity(1, units.newton).multiply(new Quantity(1, units.meter));
    expect(result.value).toBeCloseTo(1);
    expect(result.unit).toBe(units.joule);
  });

  it('should canonicalize hertz from dimensionless per second', () => {
    const result = new Quantity(1, units.dimensionless).divide(new Quantity(1, units.second));
    expect(result.value).toBeCloseTo(1);
    expect(result.unit).toBe(units.hertz);
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

describe('quantity() factory', () => {
  it('should create a quantity', () => {
    const q = quantity(5, units.meter);
    expect(q).toBeInstanceOf(Quantity);
    expect(q.value).toBe(5);
    expect(q.unit).toBe(units.meter);
  });
});

describe('isCompatibleWith / in', () => {
  it('should report compatible dimensions', () => {
    expect(quantity(1, units.meter).isCompatibleWith(quantity(1, units.kilometer))).toBe(true);
    expect(quantity(1, units.meter).isCompatibleWith(quantity(1, units.foot))).toBe(true);
  });

  it('should report incompatible dimensions', () => {
    expect(quantity(1, units.meter).isCompatibleWith(quantity(1, units.second))).toBe(false);
    expect(quantity(1, units.meter).isCompatibleWith(quantity(1, units.kilogram))).toBe(false);
  });

  it('should return numeric value in a target unit', () => {
    expect(quantity(1, units.kilometer).in(units.meter)).toBeCloseTo(1000);
    expect(quantity(1, units.hour).in(units.minute)).toBeCloseTo(60);
  });

  it('should throw on `in` with incompatible units', () => {
    expect(() => quantity(5, units.meter).in(units.second)).toThrow(IncompatibleUnitsError);
  });
});

describe('NaN and Infinity handling', () => {
  it('should allow constructing a NaN quantity', () => {
    const q = quantity(NaN, units.meter);
    expect(Number.isNaN(q.value)).toBe(true);
  });

  it('should allow constructing an Infinity quantity', () => {
    const q = quantity(Infinity, units.meter);
    expect(q.value).toBe(Infinity);
  });

  it('should propagate NaN through arithmetic', () => {
    const result = quantity(NaN, units.meter).add(quantity(1, units.meter));
    expect(Number.isNaN(result.value)).toBe(true);
  });

  it('should propagate Infinity through arithmetic', () => {
    const result = quantity(Infinity, units.meter).multiply(2);
    expect(result.value).toBe(Infinity);
  });

  it('should produce NaN for Infinity minus Infinity', () => {
    const result = quantity(Infinity, units.meter).subtract(quantity(Infinity, units.meter));
    expect(Number.isNaN(result.value)).toBe(true);
  });

  it('should produce NaN for zero times Infinity', () => {
    const result = quantity(Infinity, units.meter).multiply(0);
    expect(Number.isNaN(result.value)).toBe(true);
  });

  it('should produce Infinity when dividing by zero', () => {
    const result = quantity(1, units.meter).divide(0);
    expect(result.value).toBe(Infinity);
  });

  it('should not equal a NaN quantity', () => {
    const q = quantity(NaN, units.meter);
    expect(q.equals(quantity(NaN, units.meter))).toBe(false);
  });
});

describe('Fractional power errors', () => {
  it('should throw on a fractional exponent', () => {
    expect(() => quantity(4, units.meter).power(0.5)).toThrow('Fractional exponents');
  });

  it('should throw on a negative fractional exponent', () => {
    expect(() => quantity(4, units.meter).power(-1.5)).toThrow('Fractional exponents');
  });

  it('should allow negative integer exponents', () => {
    const result = quantity(2, units.meter).power(-1);
    expect(result.value).toBeCloseTo(0.5);
    expect(result.unit.dimension.length).toBe(-1);
  });
});

describe('Left-scalar division (scalar / quantity)', () => {
  it('should divide a scalar by a quantity via a dimensionless quantity', () => {
    const result = quantity(1, units.dimensionless).divide(quantity(2, units.second));
    expect(result.value).toBeCloseTo(0.5);
    expect(result.unit).toBe(units.hertz);
  });

  it('should produce inverse dimensions for scalar / quantity', () => {
    const result = quantity(10, units.dimensionless).divide(quantity(5, units.meterPerSecond));
    expect(result.value).toBeCloseTo(2);
    expect(result.unit.dimension.length).toBe(-1);
    expect(result.unit.dimension.time).toBe(1);
  });

  it('should yield dimensionless result for scalar / scalar-quantity', () => {
    const result = quantity(6, units.dimensionless).divide(quantity(3, units.dimensionless));
    expect(result.value).toBe(2);
    expect(result.unit.dimension.isDimensionless()).toBe(true);
  });
});

describe('Angle units', () => {
  it('should construct and convert angle units', () => {
    expect(quantity(1, units.radian).in(units.degree)).toBeCloseTo(180 / Math.PI);
    expect(quantity(360, units.degree).in(units.radian)).toBeCloseTo(2 * Math.PI);
    expect(quantity(1, units.arcminute).in(units.degree)).toBeCloseTo(1 / 60);
    expect(quantity(1, units.arcsecond).in(units.arcminute)).toBeCloseTo(1 / 60);
  });
});

describe('Frequency units', () => {
  it('should construct and convert frequency units', () => {
    expect(quantity(1, units.kilohertz).in(units.hertz)).toBeCloseTo(1000);
    expect(quantity(1, units.megahertz).in(units.hertz)).toBeCloseTo(1e6);
    expect(quantity(1, units.hertz).in(units.megahertz)).toBeCloseTo(1e-6);
  });
});

describe('Area and volume units', () => {
  it('should construct and convert area units', () => {
    expect(quantity(1, units.squareKilometer).in(units.squareMeter)).toBeCloseTo(1e6);
    expect(quantity(1, units.hectare).in(units.squareMeter)).toBeCloseTo(1e4);
    expect(quantity(1, units.squareMeter).in(units.hectare)).toBeCloseTo(1e-4);
  });

  it('should construct and convert volume units', () => {
    expect(quantity(1, units.liter).in(units.cubicMeter)).toBeCloseTo(1e-3);
    expect(quantity(1, units.milliliter).in(units.cubicMeter)).toBeCloseTo(1e-6);
    expect(quantity(1, units.cubicMeter).in(units.liter)).toBeCloseTo(1000);
  });

  it('should derive square meter from length multiplication', () => {
    const result = quantity(5, units.meter).multiply(quantity(3, units.meter));
    expect(result.value).toBeCloseTo(15);
    expect(result.unit).toBe(units.squareMeter);
  });
});

describe('Velocity units', () => {
  it('should construct and convert mile per hour', () => {
    expect(quantity(1, units.milePerHour).in(units.meterPerSecond)).toBeCloseTo(0.44704);
    expect(quantity(1, units.milePerHour).in(units.kilometerPerHour)).toBeCloseTo(1.609344);
  });
});

describe('Chemistry units', () => {
  it('should construct mole', () => {
    expect(units.mole.symbol).toBe('mol');
    expect(units.mole.dimension.amount).toBe(1);
    expect(quantity(2, units.mole).in(units.mole)).toBeCloseTo(2);
  });
});

describe('Energy units', () => {
  it('should construct and convert energy units', () => {
    expect(quantity(1, units.calorie).in(units.joule)).toBeCloseTo(4.184);
    expect(quantity(1, units.kilocalorie).in(units.joule)).toBeCloseTo(4184);
    expect(quantity(1, units.wattHour).in(units.joule)).toBeCloseTo(3600);
    expect(quantity(1, units.kilojoule).in(units.wattHour)).toBeCloseTo(1000 / 3600);
  });
});

describe('Electricity units', () => {
  it('should construct and convert volts and ohms', () => {
    expect(units.volt.symbol).toBe('V');
    expect(units.ohm.symbol).toBe('Ω');
    expect(units.volt.dimension).toEqual(new Dimension(2, 1, -3, -1));
    expect(units.ohm.dimension).toEqual(new Dimension(2, 1, -3, -2));
  });

  it('should canonicalize derived electric units', () => {
    const power = quantity(1, units.volt).multiply(quantity(1, units.ampere));
    expect(power.unit).toBe(units.watt);
    expect(power.value).toBeCloseTo(1);
    const current = quantity(1, units.volt).divide(quantity(1, units.ohm));
    expect(current.unit).toBe(units.ampere);
    expect(current.value).toBeCloseTo(1);
  });
});

describe('Derived canonicalization', () => {
  it('should canonicalize joule seconds', () => {
    const result = quantity(1, units.joule).multiply(quantity(1, units.second));
    expect(result.value).toBeCloseTo(1);
    expect(result.unit.dimension.length).toBe(2);
    expect(result.unit.dimension.mass).toBe(1);
    expect(result.unit.dimension.time).toBe(-1);
  });

  it('should canonicalize cubic meter from power', () => {
    const result = quantity(3, units.meter).power(3);
    expect(result.value).toBeCloseTo(27);
    expect(result.unit).toBe(units.cubicMeter);
  });

  it('should canonicalize hertz from inverse second', () => {
    const result = quantity(1, units.dimensionless).divide(quantity(1, units.second));
    expect(result.unit).toBe(units.hertz);
  });

  it('should canonicalize mole as amount of substance', () => {
    const result = quantity(1, units.mole);
    expect(result.unit).toBe(units.mole);
  });
});
