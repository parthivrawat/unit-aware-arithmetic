//! Standalone Rust benchmark for unit-aware arithmetic.
//!
//! Run with: cargo run --example bench

use std::hint::black_box;
use std::time::Instant;
use unit_aware_arithmetic::{units, Quantity};

const ITERS: u64 = 200_000;

fn print_row(label: &str, raw_ns: f64, qty_ns: f64) {
    let overhead = if raw_ns > 1e-12 {
        qty_ns / raw_ns
    } else {
        f64::INFINITY
    };
    println!(
        "{:20}  raw {:>10.1} ns   qty {:>10.1} ns   overhead {:>5.1}x",
        label, raw_ns, qty_ns, overhead
    );
}

fn bench_construction() {
    let iters = ITERS;

    let start = Instant::now();
    let mut s = 0.0_f64;
    for _ in 0..iters {
        s = black_box(100.0_f64);
    }
    let _ = black_box(s);
    let raw = (start.elapsed().as_nanos() as f64) / iters as f64;

    let start = Instant::now();
    let mut q = Quantity::new(100.0, units::METER);
    for _ in 0..iters {
        q = black_box(Quantity::new(100.0, units::METER));
    }
    let _ = black_box(q);
    let qty = (start.elapsed().as_nanos() as f64) / iters as f64;

    print_row("construction", raw, qty);
}

fn bench_addition() {
    let iters = ITERS;
    let a = 100.0_f64;
    let b = 0.5_f64;

    let start = Instant::now();
    let mut s = 0.0_f64;
    for _ in 0..iters {
        s = black_box(a) + black_box(b);
    }
    let _ = black_box(s);
    let raw = (start.elapsed().as_nanos() as f64) / iters as f64;

    let q1 = black_box(Quantity::new(100.0, units::METER));
    let q2 = black_box(Quantity::new(0.5, units::METER));
    let start = Instant::now();
    let mut q = q1.clone();
    for _ in 0..iters {
        q = q1.clone() + q2.clone();
    }
    let _ = black_box(q);
    let qty = (start.elapsed().as_nanos() as f64) / iters as f64;

    print_row("addition", raw, qty);
}

fn bench_multiplication() {
    let iters = ITERS;
    let a = 10.0_f64;
    let b = 9.8_f64;

    let start = Instant::now();
    let mut s = 0.0_f64;
    for _ in 0..iters {
        s = black_box(a) * black_box(b);
    }
    let _ = black_box(s);
    let raw = (start.elapsed().as_nanos() as f64) / iters as f64;

    let q1 = black_box(Quantity::new(10.0, units::KILOGRAM));
    let q2 = black_box(Quantity::new(9.8, units::METER_PER_SECOND_SQUARED));
    let start = Instant::now();
    let mut q = q1.clone();
    for _ in 0..iters {
        q = q1.clone() * q2.clone();
    }
    let _ = black_box(q);
    let qty = (start.elapsed().as_nanos() as f64) / iters as f64;

    print_row("multiplication", raw, qty);
}

fn bench_division() {
    let iters = ITERS;
    let a = 100.0_f64;
    let b = 10.0_f64;

    let start = Instant::now();
    let mut s = 0.0_f64;
    for _ in 0..iters {
        s = black_box(a) / black_box(b);
    }
    let _ = black_box(s);
    let raw = (start.elapsed().as_nanos() as f64) / iters as f64;

    let q1 = black_box(Quantity::new(100.0, units::METER));
    let q2 = black_box(Quantity::new(10.0, units::SECOND));
    let start = Instant::now();
    let mut q = q1.clone();
    for _ in 0..iters {
        q = q1.clone() / q2.clone();
    }
    let _ = black_box(q);
    let qty = (start.elapsed().as_nanos() as f64) / iters as f64;

    print_row("division", raw, qty);
}

fn bench_conversion() {
    let iters = ITERS;
    let a = 100.0_f64;

    let start = Instant::now();
    let mut s = 0.0_f64;
    for _ in 0..iters {
        s = black_box(a) * black_box(0.001);
    }
    let _ = black_box(s);
    let raw = (start.elapsed().as_nanos() as f64) / iters as f64;

    let q1 = black_box(Quantity::new(100.0, units::METER));
    let start = Instant::now();
    let mut q = q1.clone();
    for _ in 0..iters {
        q = q1.to(units::KILOMETER).unwrap();
    }
    let _ = black_box(q);
    let qty = (start.elapsed().as_nanos() as f64) / iters as f64;

    print_row("conversion", raw, qty);
}

fn main() {
    println!(
        "Rust unit-aware arithmetic benchmark ({} iterations each)",
        ITERS
    );
    println!("{}", "=".repeat(70));
    bench_construction();
    bench_addition();
    bench_multiplication();
    bench_division();
    bench_conversion();
    println!();
    println!("Times are in nanoseconds per operation and depend on the host machine.");
}
