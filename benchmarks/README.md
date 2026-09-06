# Unit-Aware Arithmetic Benchmarks

This directory contains reproducible, standalone benchmarks that compare the
overhead of the unit-aware `Quantity` arithmetic against raw numeric arithmetic.

The target implementations cover the four languages in this repo: Python,
TypeScript, Go, and Rust.

## What is measured

Each benchmark measures the same typical operations:

| Operation | Python | TypeScript | Go | Rust |
|-----------|--------|------------|----|----|
| Quantity construction | `Quantity(...)` | `quantity(...)` | `NewQuantity(...)` | `Quantity::new(...)` |
| Addition | `q1 + q2` | `q1.add(q2)` | `q1.Add(q2)` | `q1 + q2` |
| Multiplication | `q1 * q2` | `q1.multiply(q2)` | `q1.Multiply(q2)` | `q1 * q2` |
| Division | `q1 / q2` | `q1.divide(q2)` | `q1.Divide(q2)` | `q1 / q2` |
| Unit conversion | `q1.to(units.kilometer)` | `q1.to(units.kilometer)` | `q1.To(units.Kilometer)` | `q1.to(units::KILOMETER).unwrap()` |

The raw numeric baseline uses the same values without any unit tracking.

## Running the benchmarks

> **Note:** All timings are machine-dependent. Run them on your own hardware
> before drawing conclusions.

### Python

```bash
python benchmarks/bench_python.py
```

`bench_python.py` appends the `python/` directory to `sys.path` so it can
import the local `dimensional` module without installing it.

### TypeScript

The TypeScript implementation must be compiled first, then the plain JavaScript
benchmark is executed with Node:

```bash
cd typescript
npm install
npm run build
npm run bench
```

`npm run bench` builds the library and runs `node ../benchmarks/bench_typescript.js`.

### Go

```bash
cd go
go test -bench=. -benchmem
```

The benchmark file lives at `go/bench_test.go` inside the Go module so it can
access the `dimensional` package directly.

### Rust

No external benchmark dependencies are used. The benchmark is a plain example
binary that can be run with:

```bash
cd rust
cargo run --example bench
```

## Interpreting results

The library is expected to add a small, constant overhead (often in the range
of ~2–5x) over raw numeric operations. The exact factor depends on:

- CPU, memory, and runtime version
- Whether the operation allocates a new object/struct
- Whether a dimension lookup or unit conversion is required

For performance-critical paths, the recommended approach is to perform unit
validation at the edges of your code and extract raw numeric values for tight
loops, exactly as shown in the main `README.md`.
