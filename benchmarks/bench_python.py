#!/usr/bin/env python3
"""Standalone benchmark for the Python dimensional arithmetic library."""

import os
import sys
import timeit

# Make the local `python/` package importable without installation.
REPO_ROOT = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
sys.path.insert(0, os.path.join(REPO_ROOT, "python"))

from dimensional import Quantity, units

REPEATS = 200_000


def time_op(stmt, setup="pass", globs=None):
    """Return the average nanoseconds per operation for `stmt`."""
    globs = globs or {}
    total = timeit.timeit(stmt, setup=setup, number=REPEATS, globals=globs)
    return (total / REPEATS) * 1e9


def print_row(label, raw_ns, qty_ns):
    overhead = qty_ns / raw_ns if raw_ns > 1e-12 else float("inf")
    print(f"{label:20}  raw {raw_ns:>10.1f} ns   qty {qty_ns:>10.1f} ns   overhead {overhead:>5.1f}x")


def main():
    print(f"Python unit-aware arithmetic benchmark ({REPEATS:,} iterations each)")
    print("=" * 70)

    # Construction
    raw_ns = time_op("x = float(100.0)")
    globs = {"Quantity": Quantity, "units": units}
    qty_ns = time_op("x = Quantity(100, units.meter)", globs=globs)
    print_row("construction", raw_ns, qty_ns)

    # Addition (same dimension, different units, converts q2)
    setup = """
q1 = Quantity(100, units.meter)
q2 = Quantity(50, units.centimeter)
"""
    raw_ns = time_op("x = 100.0 + 0.5", setup=setup, globs=globs)
    qty_ns = time_op("x = q1 + q2", setup=setup, globs={"Quantity": Quantity, "units": units})
    print_row("addition", raw_ns, qty_ns)

    # Multiplication (canonical derived unit: Newton)
    setup = """
q1 = Quantity(10, units.kilogram)
q2 = Quantity(9.8, units.meter_per_second_squared)
"""
    raw_ns = time_op("x = 10.0 * 9.8", setup=setup, globs=globs)
    qty_ns = time_op("x = q1 * q2", setup=setup, globs={"Quantity": Quantity, "units": units})
    print_row("multiplication", raw_ns, qty_ns)

    # Division (canonical derived unit: m/s)
    setup = """
q1 = Quantity(100, units.meter)
q2 = Quantity(10, units.second)
"""
    raw_ns = time_op("x = 100.0 / 10.0", setup=setup, globs=globs)
    qty_ns = time_op("x = q1 / q2", setup=setup, globs={"Quantity": Quantity, "units": units})
    print_row("division", raw_ns, qty_ns)

    # Unit conversion
    setup = """
q1 = Quantity(100, units.meter)
"""
    raw_ns = time_op("x = 100.0 * 0.001", setup=setup, globs=globs)
    qty_ns = time_op("x = q1.to(units.kilometer)", setup=setup, globs={"Quantity": Quantity, "units": units})
    print_row("conversion", raw_ns, qty_ns)

    print()
    print("Times are in nanoseconds per operation and depend on the host machine.")


if __name__ == "__main__":
    main()
