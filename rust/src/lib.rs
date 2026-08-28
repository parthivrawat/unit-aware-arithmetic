//! Unit-Aware Numeric Arithmetic Library
//!
//! A type-safe dimensional arithmetic library that tracks units at compile/runtime
//! and prevents invalid operations.
//!
//! # Examples
//!
//! ```
//! use unit_aware_arithmetic::{Quantity, units};
//!
//! let distance = Quantity::new(100.0, units::METER);
//! let time = Quantity::new(9.58, units::SECOND);
//! let velocity = distance / time;
//! println!("{}", velocity); // 10.438... m/s
//! ```

use std::fmt;
use std::ops::{Add, Sub, Mul, Div, Neg};

/// Represents the dimensional formula of a unit (e.g., L^1 T^-2 for acceleration).
#[derive(Debug, Clone, Copy, PartialEq, Eq)]
pub struct Dimension {
    pub length: i32,
    pub mass: i32,
    pub time: i32,
    pub current: i32,
    pub temperature: i32,
    pub amount: i32,
    pub luminosity: i32,
}

impl Dimension {
    /// Create a new dimension.
    pub const fn new(
        length: i32,
        mass: i32,
        time: i32,
        current: i32,
        temperature: i32,
        amount: i32,
        luminosity: i32,
    ) -> Self {
        Self {
            length,
            mass,
            time,
            current,
            temperature,
            amount,
            luminosity,
        }
    }

    /// Create a dimensionless quantity.
    pub const fn dimensionless() -> Self {
        Self::new(0, 0, 0, 0, 0, 0, 0)
    }

    /// Check if this is a dimensionless quantity.
    pub fn is_dimensionless(&self) -> bool {
        self.length == 0
            && self.mass == 0
            && self.time == 0
            && self.current == 0
            && self.temperature == 0
            && self.amount == 0
            && self.luminosity == 0
    }
}

impl Mul for Dimension {
    type Output = Self;

    fn mul(self, other: Self) -> Self {
        Self {
            length: self.length + other.length,
            mass: self.mass + other.mass,
            time: self.time + other.time,
            current: self.current + other.current,
            temperature: self.temperature + other.temperature,
            amount: self.amount + other.amount,
            luminosity: self.luminosity + other.luminosity,
        }
    }
}

impl Div for Dimension {
    type Output = Self;

    fn div(self, other: Self) -> Self {
        Self {
            length: self.length - other.length,
            mass: self.mass - other.mass,
            time: self.time - other.time,
            current: self.current - other.current,
            temperature: self.temperature - other.temperature,
            amount: self.amount - other.amount,
            luminosity: self.luminosity - other.luminosity,
        }
    }
}

/// Represents a unit of measurement with its dimension and conversion factor.
#[derive(Debug, Clone, Copy)]
pub struct Unit {
    pub name: &'static str,
    pub symbol: &'static str,
    pub dimension: Dimension,
    pub to_base: f64,
    pub offset: f64,
}

impl Unit {
    /// Create a new unit.
    pub const fn new(
        name: &'static str,
        symbol: &'static str,
        dimension: Dimension,
        to_base: f64,
        offset: f64,
    ) -> Self {
        Self {
            name,
            symbol,
            dimension,
            to_base,
            offset,
        }
    }
}

/// Error type for incompatible unit operations.
#[derive(Debug, Clone)]
pub struct IncompatibleUnitsError {
    pub message: String,
}

impl fmt::Display for IncompatibleUnitsError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.message)
    }
}

impl std::error::Error for IncompatibleUnitsError {}

/// A numeric value with an associated unit.
#[derive(Debug, Clone, Copy)]
pub struct Quantity {
    pub value: f64,
    pub unit: Unit,
}

impl Quantity {
    /// Create a new quantity.
    pub fn new(value: f64, unit: Unit) -> Self {
        Self { value, unit }
    }

    /// Convert to another unit (must have compatible dimensions).
    pub fn to(&self, target_unit: Unit) -> Result<Self, IncompatibleUnitsError> {
        if self.unit.dimension != target_unit.dimension {
            return Err(IncompatibleUnitsError {
                message: format!(
                    "Cannot convert {} to {}: incompatible dimensions",
                    self.unit.symbol, target_unit.symbol
                ),
            });
        }

        let new_value = if self.unit.offset != 0.0 || target_unit.offset != 0.0 {
            // Affine conversion (e.g., temperature)
            let base_value = (self.value + self.unit.offset) * self.unit.to_base;
            base_value / target_unit.to_base - target_unit.offset
        } else {
            // Linear conversion
            self.value * (self.unit.to_base / target_unit.to_base)
        };

        Ok(Self::new(new_value, target_unit))
    }

    /// Raise to an integer power.
    pub fn pow(&self, exponent: i32) -> Self {
        let new_value = self.value.powi(exponent);
        let new_dimension = Dimension {
            length: self.unit.dimension.length * exponent,
            mass: self.unit.dimension.mass * exponent,
            time: self.unit.dimension.time * exponent,
            current: self.unit.dimension.current * exponent,
            temperature: self.unit.dimension.temperature * exponent,
            amount: self.unit.dimension.amount * exponent,
            luminosity: self.unit.dimension.luminosity * exponent,
        };
        let new_symbol = format!("{}^{}", self.unit.symbol, exponent);
        let new_unit = Unit {
            name: "derived",
            symbol: Box::leak(new_symbol.into_boxed_str()),
            dimension: new_dimension,
            to_base: self.unit.to_base.powi(exponent),
            offset: 0.0,
        };
        Self::new(new_value, new_unit)
    }

    /// Absolute value.
    pub fn abs(&self) -> Self {
        Self::new(self.value.abs(), self.unit)
    }

    /// Check equality with tolerance.
    pub fn approx_eq(&self, other: &Self, tolerance: f64) -> bool {
        if self.unit.dimension != other.unit.dimension {
            return false;
        }
        match other.to(self.unit) {
            Ok(converted) => (self.value - converted.value).abs() < tolerance,
            Err(_) => false,
        }
    }
}

impl fmt::Display for Quantity {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{} {}", self.value, self.unit.symbol)
    }
}

impl Add for Quantity {
    type Output = Result<Self, IncompatibleUnitsError>;

    fn add(self, other: Self) -> Self::Output {
        if self.unit.dimension != other.unit.dimension {
            return Err(IncompatibleUnitsError {
                message: format!(
                    "Cannot add {} and {}: incompatible dimensions",
                    self.unit.symbol, other.unit.symbol
                ),
            });
        }
        let other_converted = other.to(self.unit)?;
        Ok(Self::new(self.value + other_converted.value, self.unit))
    }
}

impl Sub for Quantity {
    type Output = Result<Self, IncompatibleUnitsError>;

    fn sub(self, other: Self) -> Self::Output {
        if self.unit.dimension != other.unit.dimension {
            return Err(IncompatibleUnitsError {
                message: format!(
                    "Cannot subtract {} from {}: incompatible dimensions",
                    other.unit.symbol, self.unit.symbol
                ),
            });
        }
        let other_converted = other.to(self.unit)?;
        Ok(Self::new(self.value - other_converted.value, self.unit))
    }
}

impl Mul for Quantity {
    type Output = Self;

    fn mul(self, other: Self) -> Self {
        let new_value = self.value * other.value;
        let new_dimension = self.unit.dimension * other.unit.dimension;
        let new_symbol = format!("{}·{}", self.unit.symbol, other.unit.symbol);
        let new_unit = Unit {
            name: "derived",
            symbol: Box::leak(new_symbol.into_boxed_str()),
            dimension: new_dimension,
            to_base: self.unit.to_base * other.unit.to_base,
            offset: 0.0,
        };
        Self::new(new_value, new_unit)
    }
}

impl Mul<f64> for Quantity {
    type Output = Self;

    fn mul(self, scalar: f64) -> Self {
        Self::new(self.value * scalar, self.unit)
    }
}

impl Mul<Quantity> for f64 {
    type Output = Quantity;

    fn mul(self, quantity: Quantity) -> Quantity {
        Quantity::new(self * quantity.value, quantity.unit)
    }
}

impl Div for Quantity {
    type Output = Self;

    fn div(self, other: Self) -> Self {
        let new_value = self.value / other.value;
        let new_dimension = self.unit.dimension / other.unit.dimension;
        let new_symbol = format!("{}/{}", self.unit.symbol, other.unit.symbol);
        let new_unit = Unit {
            name: "derived",
            symbol: Box::leak(new_symbol.into_boxed_str()),
            dimension: new_dimension,
            to_base: self.unit.to_base / other.unit.to_base,
            offset: 0.0,
        };
        Self::new(new_value, new_unit)
    }
}

impl Div<f64> for Quantity {
    type Output = Self;

    fn div(self, scalar: f64) -> Self {
        Self::new(self.value / scalar, self.unit)
    }
}

impl Neg for Quantity {
    type Output = Self;

    fn neg(self) -> Self {
        Self::new(-self.value, self.unit)
    }
}

/// Predefined units.
pub mod units {
    use super::*;

    // Dimensionless
    pub const DIMENSIONLESS: Unit = Unit::new("dimensionless", "", Dimension::dimensionless(), 1.0, 0.0);

    // Length
    pub const METER: Unit = Unit::new("meter", "m", Dimension::new(1, 0, 0, 0, 0, 0, 0), 1.0, 0.0);
    pub const KILOMETER: Unit = Unit::new("kilometer", "km", Dimension::new(1, 0, 0, 0, 0, 0, 0), 1000.0, 0.0);
    pub const CENTIMETER: Unit = Unit::new("centimeter", "cm", Dimension::new(1, 0, 0, 0, 0, 0, 0), 0.01, 0.0);
    pub const MILLIMETER: Unit = Unit::new("millimeter", "mm", Dimension::new(1, 0, 0, 0, 0, 0, 0), 0.001, 0.0);

    // Imperial length
    pub const INCH: Unit = Unit::new("inch", "in", Dimension::new(1, 0, 0, 0, 0, 0, 0), 0.0254, 0.0);
    pub const FOOT: Unit = Unit::new("foot", "ft", Dimension::new(1, 0, 0, 0, 0, 0, 0), 0.3048, 0.0);
    pub const YARD: Unit = Unit::new("yard", "yd", Dimension::new(1, 0, 0, 0, 0, 0, 0), 0.9144, 0.0);
    pub const MILE: Unit = Unit::new("mile", "mi", Dimension::new(1, 0, 0, 0, 0, 0, 0), 1609.34, 0.0);

    // Mass
    pub const KILOGRAM: Unit = Unit::new("kilogram", "kg", Dimension::new(0, 1, 0, 0, 0, 0, 0), 1.0, 0.0);
    pub const GRAM: Unit = Unit::new("gram", "g", Dimension::new(0, 1, 0, 0, 0, 0, 0), 0.001, 0.0);
    pub const MILLIGRAM: Unit = Unit::new("milligram", "mg", Dimension::new(0, 1, 0, 0, 0, 0, 0), 1e-6, 0.0);
    pub const TONNE: Unit = Unit::new("tonne", "t", Dimension::new(0, 1, 0, 0, 0, 0, 0), 1000.0, 0.0);

    // Imperial mass
    pub const POUND: Unit = Unit::new("pound", "lb", Dimension::new(0, 1, 0, 0, 0, 0, 0), 0.453592, 0.0);
    pub const OUNCE: Unit = Unit::new("ounce", "oz", Dimension::new(0, 1, 0, 0, 0, 0, 0), 0.0283495, 0.0);

    // Time
    pub const SECOND: Unit = Unit::new("second", "s", Dimension::new(0, 0, 1, 0, 0, 0, 0), 1.0, 0.0);
    pub const MINUTE: Unit = Unit::new("minute", "min", Dimension::new(0, 0, 1, 0, 0, 0, 0), 60.0, 0.0);
    pub const HOUR: Unit = Unit::new("hour", "h", Dimension::new(0, 0, 1, 0, 0, 0, 0), 3600.0, 0.0);
    pub const DAY: Unit = Unit::new("day", "d", Dimension::new(0, 0, 1, 0, 0, 0, 0), 86400.0, 0.0);

    // Temperature
    pub const KELVIN: Unit = Unit::new("kelvin", "K", Dimension::new(0, 0, 0, 0, 1, 0, 0), 1.0, 0.0);
    pub const CELSIUS: Unit = Unit::new("celsius", "°C", Dimension::new(0, 0, 0, 0, 1, 0, 0), 1.0, 273.15);
    pub const FAHRENHEIT: Unit = Unit::new("fahrenheit", "°F", Dimension::new(0, 0, 0, 0, 1, 0, 0), 5.0 / 9.0, 459.67);

    // Current
    pub const AMPERE: Unit = Unit::new("ampere", "A", Dimension::new(0, 0, 0, 1, 0, 0, 0), 1.0, 0.0);
    pub const MILLIAMPERE: Unit = Unit::new("milliampere", "mA", Dimension::new(0, 0, 0, 1, 0, 0, 0), 0.001, 0.0);

    // Derived units
    pub const NEWTON: Unit = Unit::new("newton", "N", Dimension::new(1, 1, -2, 0, 0, 0, 0), 1.0, 0.0);
    pub const JOULE: Unit = Unit::new("joule", "J", Dimension::new(2, 1, -2, 0, 0, 0, 0), 1.0, 0.0);
    pub const WATT: Unit = Unit::new("watt", "W", Dimension::new(2, 1, -3, 0, 0, 0, 0), 1.0, 0.0);
    pub const PASCAL: Unit = Unit::new("pascal", "Pa", Dimension::new(-1, 1, -2, 0, 0, 0, 0), 1.0, 0.0);
    pub const METER_PER_SECOND: Unit = Unit::new("meter per second", "m/s", Dimension::new(1, 0, -1, 0, 0, 0, 0), 1.0, 0.0);
    pub const METER_PER_SECOND_SQUARED: Unit = Unit::new("meter per second squared", "m/s²", Dimension::new(1, 0, -2, 0, 0, 0, 0), 1.0, 0.0);
}

#[cfg(test)]
mod tests {
    use super::*;

    const TOLERANCE: f64 = 1e-9;

    #[test]
    fn test_dimension_creation() {
        let d = Dimension::new(1, 0, -2, 0, 0, 0, 0);
        assert_eq!(d.length, 1);
        assert_eq!(d.time, -2);
    }

    #[test]
    fn test_dimension_multiply() {
        let d1 = Dimension::new(1, 0, 0, 0, 0, 0, 0);
        let d2 = Dimension::new(0, 0, -1, 0, 0, 0, 0);
        let result = d1 * d2;
        assert_eq!(result.length, 1);
        assert_eq!(result.time, -1);
    }

    #[test]
    fn test_quantity_creation() {
        let q = Quantity::new(5.0, units::METER);
        assert_eq!(q.value, 5.0);
    }

    #[test]
    fn test_add_same_unit() {
        let q1 = Quantity::new(5.0, units::METER);
        let q2 = Quantity::new(3.0, units::METER);
        let result = (q1 + q2).unwrap();
        assert_eq!(result.value, 8.0);
    }

    #[test]
    fn test_add_incompatible() {
        let q1 = Quantity::new(5.0, units::METER);
        let q2 = Quantity::new(3.0, units::SECOND);
        assert!((q1 + q2).is_err());
    }

    #[test]
    fn test_multiply_quantities() {
        let q1 = Quantity::new(5.0, units::METER);
        let q2 = Quantity::new(3.0, units::METER);
        let result = q1 * q2;
        assert_eq!(result.value, 15.0);
        assert_eq!(result.unit.dimension.length, 2);
    }

    #[test]
    fn test_divide_quantities() {
        let q1 = Quantity::new(100.0, units::METER);
        let q2 = Quantity::new(10.0, units::SECOND);
        let result = q1 / q2;
        assert_eq!(result.value, 10.0);
        assert_eq!(result.unit.dimension.length, 1);
        assert_eq!(result.unit.dimension.time, -1);
    }

    #[test]
    fn test_conversion() {
        let q = Quantity::new(1.0, units::METER);
        let result = q.to(units::CENTIMETER).unwrap();
        assert!((result.value - 100.0).abs() < TOLERANCE);
    }

    #[test]
    fn test_temperature_conversion() {
        let q = Quantity::new(0.0, units::CELSIUS);
        let result = q.to(units::KELVIN).unwrap();
        assert!((result.value - 273.15).abs() < TOLERANCE);
    }

    #[test]
    fn test_physics_velocity() {
        let distance = Quantity::new(100.0, units::METER);
        let time = Quantity::new(9.58, units::SECOND);
        let velocity = distance / time;
        assert!((velocity.value - 10.438).abs() < 0.001);
    }
}
