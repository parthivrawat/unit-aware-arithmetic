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

use std::borrow::Cow;
use std::cmp::Ordering;
use std::collections::HashMap;
use std::fmt;
use std::ops::{Add, Div, Mul, Neg, Sub};
use std::sync::OnceLock;

/// Represents the dimensional formula of a unit (e.g., L^1 T^-2 for acceleration).
#[derive(Debug, Clone, Copy, PartialEq, Eq, Hash)]
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

    /// Check whether two dimensions are compatible, i.e. quantities with
    /// these dimensions can be converted between each other and combined
    /// by addition, subtraction, and comparison.
    pub fn is_compatible_with(&self, other: &Dimension) -> bool {
        self == other
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
#[derive(Debug, Clone, PartialEq)]
pub struct Unit {
    pub name: Cow<'static, str>,
    pub symbol: Cow<'static, str>,
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
            name: Cow::Borrowed(name),
            symbol: Cow::Borrowed(symbol),
            dimension,
            to_base,
            offset,
        }
    }

    /// Convert a value expressed in this unit to the corresponding value in
    /// base (SI) units, applying the affine transformation
    /// `(value + offset) * to_base`.
    ///
    /// For linear units (`offset == 0.0`) this is simply `value * to_base`.
    pub fn to_base_value(&self, value: f64) -> f64 {
        (value + self.offset) * self.to_base
    }

    /// Convert a value expressed in base (SI) units back to this unit,
    /// applying the inverse affine transformation
    /// `value / to_base - offset`.
    ///
    /// For linear units (`offset == 0.0`) this is simply `value / to_base`.
    pub fn from_base_value(&self, value: f64) -> f64 {
        value / self.to_base - self.offset
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

/// Error type for arithmetic that is undefined for affine (offset) units
/// such as Celsius or Fahrenheit.
#[derive(Debug, Clone)]
pub struct AffineUnitArithmeticError {
    pub message: String,
}

impl fmt::Display for AffineUnitArithmeticError {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{}", self.message)
    }
}

impl std::error::Error for AffineUnitArithmeticError {}

/// A numeric value with an associated unit.
///
/// The `value` and `unit` fields are public for direct read access
/// (a deliberate API choice kept for backwards compatibility); the
/// [`Quantity::value`] and [`Quantity::unit`] accessors are provided for
/// callers that prefer getter-style access. Construct values through
/// [`Quantity::new`] rather than mutating fields in place so that
/// value/unit pairs always stay consistent.
#[derive(Debug, Clone)]
pub struct Quantity {
    pub value: f64,
    pub unit: Unit,
}

impl Quantity {
    /// Create a new quantity.
    pub fn new(value: f64, unit: Unit) -> Self {
        Self { value, unit }
    }

    /// The numeric value expressed in [`Quantity::unit`].
    pub fn value(&self) -> f64 {
        self.value
    }

    /// The unit this quantity is expressed in.
    pub fn unit(&self) -> &Unit {
        &self.unit
    }

    /// Convert to another unit (must have compatible dimensions).
    pub fn to(&self, target_unit: Unit) -> Result<Self, IncompatibleUnitsError> {
        if !self
            .unit
            .dimension
            .is_compatible_with(&target_unit.dimension)
        {
            return Err(IncompatibleUnitsError {
                message: format!(
                    "Cannot convert {} to {}: incompatible dimensions",
                    self.unit.symbol, target_unit.symbol
                ),
            });
        }

        let new_value = if self.unit.offset != 0.0 || target_unit.offset != 0.0 {
            // Affine conversion (e.g., temperature)
            target_unit.from_base_value(self.unit.to_base_value(self.value))
        } else {
            // Linear conversion
            self.value * (self.unit.to_base / target_unit.to_base)
        };

        Ok(Self::new(new_value, target_unit))
    }

    /// Try to add another quantity, returning an error if dimensions differ.
    pub fn try_add(&self, other: &Quantity) -> Result<Quantity, IncompatibleUnitsError> {
        if !self
            .unit
            .dimension
            .is_compatible_with(&other.unit.dimension)
        {
            return Err(IncompatibleUnitsError {
                message: format!(
                    "Cannot add {} and {}: incompatible dimensions",
                    self.unit.symbol, other.unit.symbol
                ),
            });
        }
        let other_converted = other.to(self.unit.clone())?;
        Ok(Quantity::new(
            self.value + other_converted.value,
            self.unit.clone(),
        ))
    }

    /// Try to subtract another quantity, returning an error if dimensions differ.
    pub fn try_sub(&self, other: &Quantity) -> Result<Quantity, IncompatibleUnitsError> {
        if !self
            .unit
            .dimension
            .is_compatible_with(&other.unit.dimension)
        {
            return Err(IncompatibleUnitsError {
                message: format!(
                    "Cannot subtract {} from {}: incompatible dimensions",
                    other.unit.symbol, self.unit.symbol
                ),
            });
        }
        let other_converted = other.to(self.unit.clone())?;
        Ok(Quantity::new(
            self.value - other_converted.value,
            self.unit.clone(),
        ))
    }

    /// Raise to an integer power.
    ///
    /// # Panics
    /// Panics with an `AffineUnitArithmeticError` message if the unit is affine
    /// (has a non-zero `offset`, e.g., Celsius or Fahrenheit).
    pub fn pow(&self, exponent: i32) -> Self {
        if self.unit.offset != 0.0 {
            panic!(
                "{}",
                AffineUnitArithmeticError {
                    message: format!(
                        "Cannot raise affine unit {} to a power; convert to an absolute (zero-offset) unit first.",
                        self.unit.symbol
                    ),
                }
            );
        }
        let base_value = self.unit.to_base_value(self.value).powi(exponent);
        let new_dimension = Dimension {
            length: self.unit.dimension.length * exponent,
            mass: self.unit.dimension.mass * exponent,
            time: self.unit.dimension.time * exponent,
            current: self.unit.dimension.current * exponent,
            temperature: self.unit.dimension.temperature * exponent,
            amount: self.unit.dimension.amount * exponent,
            luminosity: self.unit.dimension.luminosity * exponent,
        };
        if let Some(canon) = canonical_unit(new_dimension) {
            return Self::new(canon.from_base_value(base_value), canon);
        }
        let new_symbol = format!("{}^{}", self.unit.symbol, exponent);
        let new_unit = Unit {
            name: Cow::Owned(String::from("derived")),
            symbol: Cow::Owned(new_symbol),
            dimension: new_dimension,
            to_base: self.unit.to_base.powi(exponent),
            offset: 0.0,
        };
        Self::new(new_unit.from_base_value(base_value), new_unit)
    }

    /// Absolute value.
    pub fn abs(&self) -> Self {
        Self::new(self.value.abs(), self.unit.clone())
    }

    /// Check if two quantities are close within a relative and/or absolute tolerance.
    pub fn is_close(&self, other: &Self, rel_tol: f64, abs_tol: f64) -> bool {
        if !self
            .unit
            .dimension
            .is_compatible_with(&other.unit.dimension)
        {
            return false;
        }
        match other.to(self.unit.clone()) {
            Ok(converted) => {
                if self.value == converted.value {
                    return true;
                }
                if !self.value.is_finite() || !converted.value.is_finite() {
                    return false;
                }
                let diff = (self.value - converted.value).abs();
                let max = self.value.abs().max(converted.value.abs());
                diff <= abs_tol || diff <= rel_tol * max
            }
            Err(_) => false,
        }
    }

    /// Check equality with an absolute tolerance.
    pub fn approx_eq(&self, other: &Self, tolerance: f64) -> bool {
        self.is_close(other, 0.0, tolerance)
    }
}

impl PartialEq for Quantity {
    fn eq(&self, other: &Self) -> bool {
        if !self
            .unit
            .dimension
            .is_compatible_with(&other.unit.dimension)
        {
            return false;
        }
        match other.to(self.unit.clone()) {
            Ok(converted) => self.value == converted.value,
            Err(_) => false,
        }
    }
}

impl PartialOrd for Quantity {
    fn partial_cmp(&self, other: &Self) -> Option<Ordering> {
        if !self
            .unit
            .dimension
            .is_compatible_with(&other.unit.dimension)
        {
            return None;
        }
        match other.to(self.unit.clone()) {
            Ok(converted) => self.value.partial_cmp(&converted.value),
            Err(_) => None,
        }
    }
}

impl fmt::Display for Quantity {
    fn fmt(&self, f: &mut fmt::Formatter<'_>) -> fmt::Result {
        write!(f, "{} {}", self.value, self.unit.symbol)
    }
}

impl Add for Quantity {
    type Output = Self;

    fn add(self, other: Self) -> Self::Output {
        self.try_add(&other).unwrap_or_else(|e| panic!("{}", e))
    }
}

impl Sub for Quantity {
    type Output = Self;

    fn sub(self, other: Self) -> Self::Output {
        self.try_sub(&other).unwrap_or_else(|e| panic!("{}", e))
    }
}

impl Mul for Quantity {
    type Output = Self;

    /// # Panics
    /// Panics with an `AffineUnitArithmeticError` message if either unit is
    /// affine (has a non-zero `offset`, e.g., Celsius or Fahrenheit).
    fn mul(self, other: Self) -> Self {
        if self.unit.offset != 0.0 || other.unit.offset != 0.0 {
            panic!(
                "{}",
                AffineUnitArithmeticError {
                    message: format!(
                        "Cannot multiply affine units {} and {}; convert to an absolute (zero-offset) unit first.",
                        self.unit.symbol, other.unit.symbol
                    ),
                }
            );
        }
        let base_value =
            self.unit.to_base_value(self.value) * other.unit.to_base_value(other.value);
        let new_dimension = self.unit.dimension * other.unit.dimension;
        if let Some(canon) = canonical_unit(new_dimension) {
            return Self::new(canon.from_base_value(base_value), canon);
        }
        let new_symbol = format!("{}·{}", self.unit.symbol, other.unit.symbol);
        let new_unit = Unit {
            name: Cow::Owned(String::from("derived")),
            symbol: Cow::Owned(new_symbol),
            dimension: new_dimension,
            to_base: self.unit.to_base * other.unit.to_base,
            offset: 0.0,
        };
        Self::new(new_unit.from_base_value(base_value), new_unit)
    }
}

impl Mul<f64> for Quantity {
    type Output = Self;

    fn mul(self, scalar: f64) -> Self {
        Self::new(self.value * scalar, self.unit.clone())
    }
}

impl Mul<Quantity> for f64 {
    type Output = Quantity;

    fn mul(self, quantity: Quantity) -> Quantity {
        Quantity::new(self * quantity.value, quantity.unit.clone())
    }
}

impl Div for Quantity {
    type Output = Self;

    /// # Panics
    /// Panics with an `AffineUnitArithmeticError` message if either unit is
    /// affine (has a non-zero `offset`, e.g., Celsius or Fahrenheit).
    fn div(self, other: Self) -> Self {
        if self.unit.offset != 0.0 || other.unit.offset != 0.0 {
            panic!(
                "{}",
                AffineUnitArithmeticError {
                    message: format!(
                        "Cannot divide affine units {} and {}; convert to an absolute (zero-offset) unit first.",
                        self.unit.symbol, other.unit.symbol
                    ),
                }
            );
        }
        let base_value =
            self.unit.to_base_value(self.value) / other.unit.to_base_value(other.value);
        let new_dimension = self.unit.dimension / other.unit.dimension;
        if let Some(canon) = canonical_unit(new_dimension) {
            return Self::new(canon.from_base_value(base_value), canon);
        }
        let new_symbol = format!("{}/{}", self.unit.symbol, other.unit.symbol);
        let new_unit = Unit {
            name: Cow::Owned(String::from("derived")),
            symbol: Cow::Owned(new_symbol),
            dimension: new_dimension,
            to_base: self.unit.to_base / other.unit.to_base,
            offset: 0.0,
        };
        Self::new(new_unit.from_base_value(base_value), new_unit)
    }
}

impl Div<f64> for Quantity {
    type Output = Self;

    fn div(self, scalar: f64) -> Self {
        Self::new(self.value / scalar, self.unit.clone())
    }
}

impl Neg for Quantity {
    type Output = Self;

    fn neg(self) -> Self {
        Self::new(-self.value, self.unit.clone())
    }
}

/// Predefined units.
pub mod units {
    use super::*;

    // Dimensionless
    pub const DIMENSIONLESS: Unit =
        Unit::new("dimensionless", "", Dimension::dimensionless(), 1.0, 0.0);

    // Angle (dimensionless)
    pub const RADIAN: Unit = Unit::new("radian", "rad", Dimension::dimensionless(), 1.0, 0.0);
    pub const DEGREE: Unit = Unit::new(
        "degree",
        "°",
        Dimension::dimensionless(),
        std::f64::consts::PI / 180.0,
        0.0,
    );
    pub const ARCMINUTE: Unit = Unit::new(
        "arcminute",
        "′",
        Dimension::dimensionless(),
        std::f64::consts::PI / 10800.0,
        0.0,
    );
    pub const ARCSECOND: Unit = Unit::new(
        "arcsecond",
        "″",
        Dimension::dimensionless(),
        std::f64::consts::PI / 648000.0,
        0.0,
    );

    // Length
    pub const METER: Unit = Unit::new("meter", "m", Dimension::new(1, 0, 0, 0, 0, 0, 0), 1.0, 0.0);
    pub const KILOMETER: Unit = Unit::new(
        "kilometer",
        "km",
        Dimension::new(1, 0, 0, 0, 0, 0, 0),
        1000.0,
        0.0,
    );
    pub const CENTIMETER: Unit = Unit::new(
        "centimeter",
        "cm",
        Dimension::new(1, 0, 0, 0, 0, 0, 0),
        0.01,
        0.0,
    );
    pub const MILLIMETER: Unit = Unit::new(
        "millimeter",
        "mm",
        Dimension::new(1, 0, 0, 0, 0, 0, 0),
        0.001,
        0.0,
    );

    // Imperial length
    pub const INCH: Unit = Unit::new(
        "inch",
        "in",
        Dimension::new(1, 0, 0, 0, 0, 0, 0),
        0.0254,
        0.0,
    );
    pub const FOOT: Unit = Unit::new(
        "foot",
        "ft",
        Dimension::new(1, 0, 0, 0, 0, 0, 0),
        0.3048,
        0.0,
    );
    pub const YARD: Unit = Unit::new(
        "yard",
        "yd",
        Dimension::new(1, 0, 0, 0, 0, 0, 0),
        0.9144,
        0.0,
    );
    pub const MILE: Unit = Unit::new(
        "mile",
        "mi",
        Dimension::new(1, 0, 0, 0, 0, 0, 0),
        1609.344,
        0.0,
    );

    // Mass
    pub const KILOGRAM: Unit = Unit::new(
        "kilogram",
        "kg",
        Dimension::new(0, 1, 0, 0, 0, 0, 0),
        1.0,
        0.0,
    );
    pub const GRAM: Unit = Unit::new("gram", "g", Dimension::new(0, 1, 0, 0, 0, 0, 0), 0.001, 0.0);
    pub const MILLIGRAM: Unit = Unit::new(
        "milligram",
        "mg",
        Dimension::new(0, 1, 0, 0, 0, 0, 0),
        1e-6,
        0.0,
    );
    pub const TONNE: Unit = Unit::new(
        "tonne",
        "t",
        Dimension::new(0, 1, 0, 0, 0, 0, 0),
        1000.0,
        0.0,
    );

    // Imperial mass
    pub const POUND: Unit = Unit::new(
        "pound",
        "lb",
        Dimension::new(0, 1, 0, 0, 0, 0, 0),
        0.453592,
        0.0,
    );
    pub const OUNCE: Unit = Unit::new(
        "ounce",
        "oz",
        Dimension::new(0, 1, 0, 0, 0, 0, 0),
        0.0283495,
        0.0,
    );

    // Time
    pub const SECOND: Unit =
        Unit::new("second", "s", Dimension::new(0, 0, 1, 0, 0, 0, 0), 1.0, 0.0);
    pub const MINUTE: Unit = Unit::new(
        "minute",
        "min",
        Dimension::new(0, 0, 1, 0, 0, 0, 0),
        60.0,
        0.0,
    );
    pub const HOUR: Unit = Unit::new(
        "hour",
        "h",
        Dimension::new(0, 0, 1, 0, 0, 0, 0),
        3600.0,
        0.0,
    );
    pub const DAY: Unit = Unit::new(
        "day",
        "d",
        Dimension::new(0, 0, 1, 0, 0, 0, 0),
        86400.0,
        0.0,
    );

    // Temperature
    pub const KELVIN: Unit =
        Unit::new("kelvin", "K", Dimension::new(0, 0, 0, 0, 1, 0, 0), 1.0, 0.0);
    pub const CELSIUS: Unit = Unit::new(
        "celsius",
        "°C",
        Dimension::new(0, 0, 0, 0, 1, 0, 0),
        1.0,
        273.15,
    );
    pub const FAHRENHEIT: Unit = Unit::new(
        "fahrenheit",
        "°F",
        Dimension::new(0, 0, 0, 0, 1, 0, 0),
        5.0 / 9.0,
        459.67,
    );

    // Current
    pub const AMPERE: Unit =
        Unit::new("ampere", "A", Dimension::new(0, 0, 0, 1, 0, 0, 0), 1.0, 0.0);
    pub const MILLIAMPERE: Unit = Unit::new(
        "milliampere",
        "mA",
        Dimension::new(0, 0, 0, 1, 0, 0, 0),
        0.001,
        0.0,
    );

    // Amount of substance
    pub const MOLE: Unit = Unit::new("mole", "mol", Dimension::new(0, 0, 0, 0, 0, 1, 0), 1.0, 0.0);

    // Derived units
    pub const NEWTON: Unit = Unit::new(
        "newton",
        "N",
        Dimension::new(1, 1, -2, 0, 0, 0, 0),
        1.0,
        0.0,
    );
    pub const JOULE: Unit = Unit::new("joule", "J", Dimension::new(2, 1, -2, 0, 0, 0, 0), 1.0, 0.0);
    pub const WATT: Unit = Unit::new("watt", "W", Dimension::new(2, 1, -3, 0, 0, 0, 0), 1.0, 0.0);
    pub const PASCAL: Unit = Unit::new(
        "pascal",
        "Pa",
        Dimension::new(-1, 1, -2, 0, 0, 0, 0),
        1.0,
        0.0,
    );
    pub const KILOJOULE: Unit = Unit::new(
        "kilojoule",
        "kJ",
        Dimension::new(2, 1, -2, 0, 0, 0, 0),
        1000.0,
        0.0,
    );
    pub const KILOWATT: Unit = Unit::new(
        "kilowatt",
        "kW",
        Dimension::new(2, 1, -3, 0, 0, 0, 0),
        1000.0,
        0.0,
    );
    pub const KILOPASCAL: Unit = Unit::new(
        "kilopascal",
        "kPa",
        Dimension::new(-1, 1, -2, 0, 0, 0, 0),
        1000.0,
        0.0,
    );
    pub const METER_PER_SECOND: Unit = Unit::new(
        "meter per second",
        "m/s",
        Dimension::new(1, 0, -1, 0, 0, 0, 0),
        1.0,
        0.0,
    );
    pub const KILOMETER_PER_HOUR: Unit = Unit::new(
        "kilometer per hour",
        "km/h",
        Dimension::new(1, 0, -1, 0, 0, 0, 0),
        1000.0 / 3600.0,
        0.0,
    );
    pub const METER_PER_SECOND_SQUARED: Unit = Unit::new(
        "meter per second squared",
        "m/s²",
        Dimension::new(1, 0, -2, 0, 0, 0, 0),
        1.0,
        0.0,
    );

    // Area (m²)
    pub const SQUARE_METER: Unit = Unit::new(
        "square meter",
        "m²",
        Dimension::new(2, 0, 0, 0, 0, 0, 0),
        1.0,
        0.0,
    );

    // Volume (m³)
    pub const CUBIC_METER: Unit = Unit::new(
        "cubic meter",
        "m³",
        Dimension::new(3, 0, 0, 0, 0, 0, 0),
        1.0,
        0.0,
    );

    // Frequency (1/s)
    pub const HERTZ: Unit = Unit::new(
        "hertz",
        "Hz",
        Dimension::new(0, 0, -1, 0, 0, 0, 0),
        1.0,
        0.0,
    );
    pub const KILOHERTZ: Unit = Unit::new(
        "kilohertz",
        "kHz",
        Dimension::new(0, 0, -1, 0, 0, 0, 0),
        1e3,
        0.0,
    );
    pub const MEGAHERTZ: Unit = Unit::new(
        "megahertz",
        "MHz",
        Dimension::new(0, 0, -1, 0, 0, 0, 0),
        1e6,
        0.0,
    );

    // Area
    pub const SQUARE_KILOMETER: Unit = Unit::new(
        "square kilometer",
        "km²",
        Dimension::new(2, 0, 0, 0, 0, 0, 0),
        1e6,
        0.0,
    );
    pub const HECTARE: Unit = Unit::new(
        "hectare",
        "ha",
        Dimension::new(2, 0, 0, 0, 0, 0, 0),
        1e4,
        0.0,
    );

    // Volume
    pub const LITER: Unit = Unit::new("liter", "L", Dimension::new(3, 0, 0, 0, 0, 0, 0), 1e-3, 0.0);
    pub const MILLILITER: Unit = Unit::new(
        "milliliter",
        "mL",
        Dimension::new(3, 0, 0, 0, 0, 0, 0),
        1e-6,
        0.0,
    );

    // Velocity
    pub const MILE_PER_HOUR: Unit = Unit::new(
        "mile per hour",
        "mph",
        Dimension::new(1, 0, -1, 0, 0, 0, 0),
        0.44704,
        0.0,
    );

    // Energy
    pub const CALORIE: Unit = Unit::new(
        "calorie",
        "cal",
        Dimension::new(2, 1, -2, 0, 0, 0, 0),
        4.184,
        0.0,
    );
    pub const KILOCALORIE: Unit = Unit::new(
        "kilocalorie",
        "kcal",
        Dimension::new(2, 1, -2, 0, 0, 0, 0),
        4184.0,
        0.0,
    );
    pub const WATT_HOUR: Unit = Unit::new(
        "watt-hour",
        "Wh",
        Dimension::new(2, 1, -2, 0, 0, 0, 0),
        3600.0,
        0.0,
    );

    // Electricity
    pub const VOLT: Unit = Unit::new("volt", "V", Dimension::new(2, 1, -3, -1, 0, 0, 0), 1.0, 0.0);
    pub const OHM: Unit = Unit::new("ohm", "Ω", Dimension::new(2, 1, -3, -2, 0, 0, 0), 1.0, 0.0);
    pub const COULOMB: Unit = Unit::new(
        "coulomb",
        "C",
        Dimension::new(0, 0, 1, 1, 0, 0, 0),
        1.0,
        0.0,
    );
}

/// Return the canonical SI/base unit for a dimension, if one is known.
fn canonical_unit(dimension: Dimension) -> Option<Unit> {
    use units::*;
    static CANONICAL: OnceLock<HashMap<Dimension, Unit>> = OnceLock::new();
    let map = CANONICAL.get_or_init(|| {
        [
            DIMENSIONLESS,
            METER,
            KILOGRAM,
            SECOND,
            KELVIN,
            AMPERE,
            MOLE,
            NEWTON,
            JOULE,
            WATT,
            PASCAL,
            VOLT,
            OHM,
            COULOMB,
            METER_PER_SECOND,
            METER_PER_SECOND_SQUARED,
            SQUARE_METER,
            CUBIC_METER,
            HERTZ,
        ]
        .into_iter()
        .map(|unit| (unit.dimension, unit))
        .collect()
    });
    map.get(&dimension).cloned()
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
        let result = q1 + q2;
        assert_eq!(result.value, 8.0);
    }

    #[test]
    fn test_add_incompatible() {
        let q1 = Quantity::new(5.0, units::METER);
        let q2 = Quantity::new(3.0, units::SECOND);
        assert!(q1.try_add(&q2).is_err());
    }

    #[test]
    #[should_panic(expected = "Cannot add")]
    fn test_add_operator_panics_incompatible() {
        let q1 = Quantity::new(5.0, units::METER);
        let q2 = Quantity::new(3.0, units::SECOND);
        let _ = q1 + q2;
    }

    #[test]
    fn test_try_add_error_message() {
        let q1 = Quantity::new(5.0, units::METER);
        let q2 = Quantity::new(3.0, units::SECOND);
        let err = q1.try_add(&q2).unwrap_err();
        assert!(err.message.contains("Cannot add"));
    }

    #[test]
    fn test_try_sub_error_message() {
        let q1 = Quantity::new(5.0, units::METER);
        let q2 = Quantity::new(3.0, units::SECOND);
        let err = q1.try_sub(&q2).unwrap_err();
        assert!(err.message.contains("Cannot subtract"));
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
    fn test_kilojoule_to_joule() {
        let q = Quantity::new(2.5, units::KILOJOULE);
        let result = q.to(units::JOULE).unwrap();
        assert!((result.value - 2500.0).abs() < TOLERANCE);
    }

    #[test]
    fn test_kilowatt_to_watt() {
        let q = Quantity::new(1.2, units::KILOWATT);
        let result = q.to(units::WATT).unwrap();
        assert!((result.value - 1200.0).abs() < TOLERANCE);
    }

    #[test]
    fn test_kilopascal_to_pascal() {
        let q = Quantity::new(100.0, units::KILOPASCAL);
        let result = q.to(units::PASCAL).unwrap();
        assert!((result.value - 100_000.0).abs() < TOLERANCE);
    }

    #[test]
    fn test_kilometer_per_hour_to_meter_per_second() {
        let q = Quantity::new(36.0, units::KILOMETER_PER_HOUR);
        let result = q.to(units::METER_PER_SECOND).unwrap();
        assert!((result.value - 10.0).abs() < TOLERANCE);
    }

    #[test]
    fn test_temperature_conversion() {
        let q = Quantity::new(0.0, units::CELSIUS);
        let result = q.to(units::KELVIN).unwrap();
        assert!((result.value - 273.15).abs() < TOLERANCE);

        let f = Quantity::new(98.6, units::FAHRENHEIT);
        let c = f.to(units::CELSIUS).unwrap();
        assert!((c.value - 37.0).abs() < TOLERANCE);
    }

    #[test]
    fn test_partial_eq_compatible() {
        let q1 = Quantity::new(1.0, units::METER);
        let q2 = Quantity::new(100.0, units::CENTIMETER);
        assert_eq!(q1, q2);
    }

    #[test]
    fn test_partial_eq_incompatible() {
        let q1 = Quantity::new(1.0, units::METER);
        let q2 = Quantity::new(1.0, units::SECOND);
        assert_ne!(q1, q2);
    }

    #[test]
    fn test_is_close() {
        let q1 = Quantity::new(1.0, units::METER);
        let q2 = Quantity::new(1.0 + 1e-7, units::METER);
        assert!(q1.is_close(&q2, 1e-6, 0.0));
        assert!(q1.approx_eq(&q2, 1e-6));

        let q3 = Quantity::new(2.0, units::METER);
        assert!(!q1.is_close(&q3, 1e-9, 0.0));

        let q4 = Quantity::new(1e-12, units::METER);
        let q5 = Quantity::new(2e-12, units::METER);
        assert!(q4.is_close(&q5, 1e-9, 1e-9));

        let q6 = Quantity::new(1.0, units::SECOND);
        assert!(!q1.is_close(&q6, 1e-9, 1e-9));

        let q7 = Quantity::new(1.0, units::METER);
        assert!(q1.is_close(&q7, 0.0, 0.0));

        let exact = Quantity::new(1.0, units::METER);
        let exact_cm = Quantity::new(100.0, units::CENTIMETER);
        assert_eq!(exact, exact_cm);
    }

    #[test]
    fn test_partial_ord() {
        let q1 = Quantity::new(1.0, units::METER);
        let q2 = Quantity::new(50.0, units::CENTIMETER);
        assert!(q1 > q2);
        assert!(q2 < q1);
    }

    #[test]
    fn test_derived_force() {
        let m = Quantity::new(10.0, units::KILOGRAM);
        let a = Quantity::new(5.0, units::METER_PER_SECOND_SQUARED);
        let f = m * a;
        assert!((f.value - 50.0).abs() < TOLERANCE);
        assert_eq!(f.unit.dimension, units::NEWTON.dimension);
    }

    #[test]
    fn test_derived_kinetic_energy() {
        let m = Quantity::new(2.0, units::KILOGRAM);
        let v = Quantity::new(3.0, units::METER_PER_SECOND);
        let ke = 0.5 * m * v.clone() * v;
        assert!((ke.value - 9.0).abs() < TOLERANCE);
        assert_eq!(ke.unit.dimension, units::JOULE.dimension);
    }

    #[test]
    fn test_physics_velocity() {
        let distance = Quantity::new(100.0, units::METER);
        let time = Quantity::new(9.58, units::SECOND);
        let velocity = distance / time;
        assert!((velocity.value - 10.438).abs() < 0.001);
    }

    #[test]
    #[should_panic(expected = "Cannot multiply affine units")]
    fn test_multiply_celsius_panics() {
        let _ = Quantity::new(2.0, units::CELSIUS) * Quantity::new(3.0, units::CELSIUS);
    }

    #[test]
    #[should_panic(expected = "Cannot divide affine units")]
    fn test_divide_by_celsius_panics() {
        let _ = Quantity::new(100.0, units::METER) / Quantity::new(2.0, units::CELSIUS);
    }

    #[test]
    #[should_panic(expected = "Cannot divide affine units")]
    fn test_divide_celsius_panics() {
        let _ = Quantity::new(2.0, units::CELSIUS) / Quantity::new(100.0, units::METER);
    }

    #[test]
    #[should_panic(expected = "Cannot raise affine unit")]
    fn test_pow_celsius_panics() {
        let _ = Quantity::new(2.0, units::CELSIUS).pow(2);
    }

    #[test]
    #[should_panic(expected = "Cannot multiply affine units")]
    fn test_multiply_fahrenheit_panics() {
        let _ = Quantity::new(2.0, units::FAHRENHEIT) * Quantity::new(1.0, units::FAHRENHEIT);
    }

    #[test]
    fn test_scalar_ops_on_affine_units_allowed() {
        let q = Quantity::new(2.0, units::CELSIUS);
        assert_eq!((q.clone() * 3.0).value, 6.0);
        assert_eq!((q / 2.0).value, 1.0);
    }

    #[test]
    fn test_kelvin_multiply_allowed() {
        let result = Quantity::new(2.0, units::KELVIN) * Quantity::new(3.0, units::KELVIN);
        assert_eq!(result.value, 6.0);
    }

    #[test]
    fn test_affine_error_is_error() {
        let err = AffineUnitArithmeticError {
            message: String::from("test"),
        };
        let _: &dyn std::error::Error = &err;
        assert_eq!(format!("{}", err), "test");
    }

    #[test]
    fn test_canonical_square_meter_multiply() {
        let result = Quantity::new(5.0, units::METER) * Quantity::new(3.0, units::METER);
        assert_eq!(result.value, 15.0);
        assert_eq!(result.unit.symbol, "m²");
        assert_eq!(result.unit.dimension, units::SQUARE_METER.dimension);
    }

    #[test]
    fn test_canonical_square_meter_pow() {
        let result = Quantity::new(3.0, units::METER).pow(2);
        assert_eq!(result.value, 9.0);
        assert_eq!(result.unit.symbol, "m²");
    }

    #[test]
    fn test_canonical_meter_per_second() {
        let result = Quantity::new(100.0, units::METER) / Quantity::new(10.0, units::SECOND);
        assert_eq!(result.value, 10.0);
        assert_eq!(result.unit.symbol, "m/s");
    }

    #[test]
    fn test_canonical_newton() {
        let force = Quantity::new(10.0, units::KILOGRAM)
            * Quantity::new(9.8, units::METER_PER_SECOND_SQUARED);
        assert!((force.value - 98.0).abs() < TOLERANCE);
        assert_eq!(force.unit.symbol, "N");
    }

    #[test]
    fn test_canonical_kinetic_energy_joule() {
        let m = Quantity::new(2.0, units::KILOGRAM);
        let v = Quantity::new(10.0, units::METER_PER_SECOND);
        let energy = m * v.pow(2);
        assert!((energy.value - 200.0).abs() < TOLERANCE);
        assert_eq!(energy.unit.symbol, "J");
        let ke = energy * 0.5;
        assert!((ke.value - 100.0).abs() < TOLERANCE);
        assert_eq!(ke.unit.symbol, "J");
    }

    #[test]
    fn test_canonical_watt() {
        let result = Quantity::new(1.0, units::KILOJOULE) / Quantity::new(1.0, units::SECOND);
        assert!((result.value - 1000.0).abs() < TOLERANCE);
        assert_eq!(result.unit.symbol, "W");
    }

    #[test]
    fn test_canonical_joule_per_meter_newton() {
        let result = Quantity::new(1.0, units::JOULE) / Quantity::new(1.0, units::METER);
        assert!((result.value - 1.0).abs() < TOLERANCE);
        assert_eq!(result.unit.symbol, "N");
    }

    #[test]
    fn test_canonical_newton_times_meter_joule() {
        let result = Quantity::new(1.0, units::NEWTON) * Quantity::new(1.0, units::METER);
        assert!((result.value - 1.0).abs() < TOLERANCE);
        assert_eq!(result.unit.symbol, "J");
    }

    #[test]
    fn test_canonical_hertz() {
        let result = Quantity::new(1.0, units::DIMENSIONLESS) / Quantity::new(1.0, units::SECOND);
        assert!((result.value - 1.0).abs() < TOLERANCE);
        assert_eq!(result.unit.symbol, "Hz");
    }

    #[test]
    fn test_unknown_dimension_keeps_generated_symbol() {
        let result = Quantity::new(1.0, units::METER) * Quantity::new(1.0, units::SECOND);
        assert_eq!(result.unit.symbol, "m·s");
    }

    // --- Edge cases: non-finite values, zero/negative, powers ---

    #[test]
    fn test_nan_value_never_equals() {
        let nan = Quantity::new(f64::NAN, units::METER);
        let one = Quantity::new(1.0, units::METER);
        assert_ne!(nan, one);
        assert!(!nan.is_close(&one, f64::INFINITY, f64::INFINITY));
        assert!(nan.partial_cmp(&one).is_none());
    }

    #[test]
    fn test_infinity_arithmetic() {
        let inf = Quantity::new(f64::INFINITY, units::METER);
        let one = Quantity::new(1.0, units::METER);
        let sum = inf + one;
        assert_eq!(sum.value, f64::INFINITY);
        assert_eq!(
            Quantity::new(f64::INFINITY, units::METER),
            Quantity::new(f64::INFINITY, units::METER)
        );
    }

    #[test]
    fn test_zero_quantity() {
        let zero = Quantity::new(0.0, units::METER);
        let five = Quantity::new(5.0, units::METER);
        assert_eq!((zero.clone() + five.clone()).value, 5.0);
        assert_eq!(zero, Quantity::new(0.0, units::KILOMETER));
    }

    #[test]
    fn test_negative_values_and_neg() {
        let q = Quantity::new(-5.0, units::METER);
        assert_eq!((-q.clone()).value, 5.0);
        assert_eq!(q.abs().value, 5.0);
        assert!(q < Quantity::new(0.0, units::METER));
    }

    #[test]
    fn test_pow_zero_gives_dimensionless() {
        let result = Quantity::new(42.0, units::METER).pow(0);
        assert_eq!(result.value, 1.0);
        assert!(result.unit.dimension.is_dimensionless());
    }

    #[test]
    fn test_pow_negative_exponent() {
        // pow only supports integer exponents; fractional powers are not
        // part of the API. A negative exponent inverts the dimension.
        let result = Quantity::new(2.0, units::SECOND).pow(-1);
        assert!((result.value - 0.5).abs() < TOLERANCE);
        assert_eq!(result.unit.dimension, units::HERTZ.dimension);
    }

    #[test]
    fn test_scalar_and_reverse_scalar_mul() {
        let q = Quantity::new(3.0, units::KILOGRAM);
        assert_eq!((q.clone() * 2.0).value, 6.0);
        assert_eq!((2.0 * q.clone()).value, 6.0);
        assert_eq!((q / 3.0).value, 1.0);
    }

    #[test]
    fn test_comparison_incompatible_is_none() {
        let m = Quantity::new(1.0, units::METER);
        let s = Quantity::new(1.0, units::SECOND);
        assert!(m.partial_cmp(&s).is_none());
    }

    #[test]
    fn test_approx_eq_across_units() {
        let m = Quantity::new(1.0, units::METER);
        let cm = Quantity::new(100.000001, units::CENTIMETER);
        assert!(m.approx_eq(&cm, 1e-4));
        assert!(!m.approx_eq(&cm, 1e-12));
    }

    #[test]
    fn test_display_format() {
        let q = Quantity::new(9.81, units::METER_PER_SECOND_SQUARED);
        assert_eq!(format!("{}", q), "9.81 m/s²");
    }

    #[test]
    fn test_to_incompatible_error_message() {
        let err = Quantity::new(1.0, units::METER)
            .to(units::SECOND)
            .unwrap_err();
        assert!(err.message.contains("Cannot convert"));
        let _: &dyn std::error::Error = &err;
    }

    #[test]
    fn test_angle_conversions() {
        let rad = Quantity::new(1.0, units::RADIAN).to(units::DEGREE).unwrap();
        assert!((rad.value - 180.0 / std::f64::consts::PI).abs() < TOLERANCE);

        let one_deg = Quantity::new(1.0, units::DEGREE).to(units::RADIAN).unwrap();
        assert!((one_deg.value - std::f64::consts::PI / 180.0).abs() < TOLERANCE);

        let one_deg_arcsec = Quantity::new(1.0, units::DEGREE)
            .to(units::ARCSECOND)
            .unwrap();
        assert!((one_deg_arcsec.value - 3600.0).abs() < TOLERANCE);
    }

    #[test]
    fn test_frequency_conversions() {
        let khz = Quantity::new(1.0, units::KILOHERTZ)
            .to(units::HERTZ)
            .unwrap();
        assert!((khz.value - 1000.0).abs() < TOLERANCE);

        let mhz = Quantity::new(1.0, units::MEGAHERTZ)
            .to(units::HERTZ)
            .unwrap();
        assert!((mhz.value - 1_000_000.0).abs() < TOLERANCE);
    }

    #[test]
    fn test_area_conversions() {
        let ha = Quantity::new(1.0, units::HECTARE)
            .to(units::SQUARE_METER)
            .unwrap();
        assert!((ha.value - 10000.0).abs() < TOLERANCE);

        let sq_km = Quantity::new(1.0, units::SQUARE_KILOMETER)
            .to(units::HECTARE)
            .unwrap();
        assert!((sq_km.value - 100.0).abs() < TOLERANCE);
    }

    #[test]
    fn test_volume_conversions() {
        let l = Quantity::new(1.0, units::LITER)
            .to(units::CUBIC_METER)
            .unwrap();
        assert!((l.value - 0.001).abs() < TOLERANCE);

        let ml = Quantity::new(1.0, units::MILLILITER)
            .to(units::LITER)
            .unwrap();
        assert!((ml.value - 0.001).abs() < TOLERANCE);
    }

    #[test]
    fn test_mile_per_hour_to_meter_per_second() {
        let mph = Quantity::new(60.0, units::MILE_PER_HOUR)
            .to(units::METER_PER_SECOND)
            .unwrap();
        assert!((mph.value - 26.8224).abs() < TOLERANCE);
    }

    #[test]
    fn test_mole_dimension() {
        let q = Quantity::new(1.0, units::MOLE);
        assert_eq!(q.unit.dimension.amount, 1);
    }

    #[test]
    fn test_energy_conversions() {
        let cal = Quantity::new(1.0, units::CALORIE).to(units::JOULE).unwrap();
        assert!((cal.value - 4.184).abs() < TOLERANCE);

        let kcal = Quantity::new(1.0, units::KILOCALORIE)
            .to(units::JOULE)
            .unwrap();
        assert!((kcal.value - 4184.0).abs() < TOLERANCE);

        let wh = Quantity::new(1.0, units::WATT_HOUR)
            .to(units::JOULE)
            .unwrap();
        assert!((wh.value - 3600.0).abs() < TOLERANCE);
    }

    #[test]
    fn test_electricity_canonical() {
        let volt = Quantity::new(1.0, units::WATT) / Quantity::new(1.0, units::AMPERE);
        assert!((volt.value - 1.0).abs() < TOLERANCE);
        assert_eq!(volt.unit.symbol, "V");
        assert_eq!(volt.unit.dimension, units::VOLT.dimension);

        let ohm = Quantity::new(1.0, units::VOLT) / Quantity::new(1.0, units::AMPERE);
        assert!((ohm.value - 1.0).abs() < TOLERANCE);
        assert_eq!(ohm.unit.symbol, "Ω");
        assert_eq!(ohm.unit.dimension, units::OHM.dimension);
    }

    #[test]
    fn test_canonical_cubic_meter_pow() {
        let result = Quantity::new(2.0, units::METER).pow(3);
        assert_eq!(result.value, 8.0);
        assert_eq!(result.unit.symbol, "m³");
        assert_eq!(result.unit.dimension, units::CUBIC_METER.dimension);
    }

    #[test]
    fn test_canonical_coulomb() {
        let charge = Quantity::new(2.0, units::AMPERE) * Quantity::new(3.0, units::SECOND);
        assert!((charge.value - 6.0).abs() < TOLERANCE);
        assert_eq!(charge.unit.symbol, "C");
        assert_eq!(charge.unit.dimension, units::COULOMB.dimension);
    }

    #[test]
    fn test_dimension_is_compatible_with() {
        assert!(units::METER
            .dimension
            .is_compatible_with(&units::KILOMETER.dimension));
        assert!(!units::METER
            .dimension
            .is_compatible_with(&units::SECOND.dimension));
    }

    #[test]
    fn test_unit_base_value_roundtrip() {
        // Linear unit
        let km = units::KILOMETER;
        assert!((km.to_base_value(2.0) - 2000.0).abs() < TOLERANCE);
        assert!((km.from_base_value(2000.0) - 2.0).abs() < TOLERANCE);

        // Affine unit (temperature)
        let c = units::CELSIUS;
        assert!((c.to_base_value(100.0) - 373.15).abs() < TOLERANCE);
        assert!((c.from_base_value(273.15)).abs() < TOLERANCE);
        let f = units::FAHRENHEIT;
        assert!((f.to_base_value(32.0) - 273.15).abs() < TOLERANCE);
        assert!((f.from_base_value(273.15) - 32.0).abs() < TOLERANCE);
    }

    #[test]
    fn test_quantity_accessors() {
        let q = Quantity::new(5.0, units::METER);
        assert_eq!(q.value(), 5.0);
        assert_eq!(q.unit().symbol, "m");
    }
}
