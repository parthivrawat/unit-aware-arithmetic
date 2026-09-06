# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/), and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [2.0.0] - 2026-09-06

> **Note:** This release includes a **breaking Go API change**: `Multiply`, `Divide`, and `Power` now return `(Quantity, error)` instead of panicking.

### Added

- GitHub Actions CI workflows for all four language implementations (Python, TypeScript, Go, Rust).
- Language-specific benchmarks for all four implementations.
- `EXAMPLES.md` and README fixes across the repository.

#### Python
- `Quantity` is now a frozen dataclass and supports hashing.
- Split `__repr__` and `__str__` representations.
- Added `__rtruediv__`, `__format__`, and `value_in` methods.

#### TypeScript
- Added overloaded `multiply` and `divide` signatures.
- Added `isCompatibleWith` and `in` unit-conversion helpers.
- Added `quantity()` factory function.
- New edge-case tests.

#### Go
- Added `Value()` and `Unit()` getters for `Quantity`.
- Verified `LessThanOrEqual` and `GreaterThanOrEqual` are present.

#### Rust
- Added missing units: `KILOJOULE`, `KILOWATT`, `KILOPASCAL`, and `KILOMETER_PER_HOUR`.
- Implemented `PartialEq`/`PartialOrd` and an `approx_eq` helper.
- Added approximately 52 tests.
- Added clippy and rustfmt checks to CI.

### Changed

#### Go (breaking)
- `Multiply`, `Divide`, and `Power` now return `(Quantity, error)` instead of panicking. This is a breaking API change for Go consumers.

#### TypeScript
- Migrated ESLint configuration to the new flat config format.

### Removed

#### Python
- Removed `setup.py` in favor of `pyproject.toml` for package metadata.

### Fixed

#### Rust
- Fixed a `Box::leak` lifetime issue in the unit registry.
