#!/usr/bin/env node
/** TypeScript (compiled) dimensional arithmetic benchmark.
 *
 * Run from the `typescript/` directory after building:
 *   npm run build
 *   node ../benchmarks/bench_typescript.js
 *
 * This is a CommonJS wrapper that dynamically loads the compiled ESM module.
 */

const { performance } = require('node:perf_hooks');

const ITERATIONS = 200000;
let __benchSink;

function bench(label, rawFn, qtyFn) {
  let start = performance.now();
  for (let i = 0; i < ITERATIONS; i++) {
    __benchSink = rawFn();
  }
  let rawTotal = performance.now() - start;
  let rawNs = (rawTotal / ITERATIONS) * 1e6;

  start = performance.now();
  for (let i = 0; i < ITERATIONS; i++) {
    __benchSink = qtyFn();
  }
  let qtyTotal = performance.now() - start;
  let qtyNs = (qtyTotal / ITERATIONS) * 1e6;

  let overhead = rawNs > 1e-12 ? qtyNs / rawNs : Infinity;
  console.log(`${label.padEnd(20)}  raw ${rawNs.toFixed(1).padStart(10)} ns   qty ${qtyNs.toFixed(1).padStart(10)} ns   overhead ${overhead.toFixed(1).padStart(5)}x`);
}

import('../typescript/dist/dimensional.js').then(({ quantity, units }) => {
  console.log(`TypeScript unit-aware arithmetic benchmark (${ITERATIONS.toLocaleString()} iterations each)`);
  console.log('='.repeat(70));

  bench(
    'construction',
    () => 100.0,
    () => quantity(100, units.meter)
  );

  const q1_add = quantity(100, units.meter);
  const q2_add = quantity(50, units.centimeter);
  bench(
    'addition',
    () => 100.0 + 0.5,
    () => q1_add.add(q2_add)
  );

  const q1_mul = quantity(10, units.kilogram);
  const q2_mul = quantity(9.8, units.meterPerSecondSquared);
  bench(
    'multiplication',
    () => 10.0 * 9.8,
    () => q1_mul.multiply(q2_mul)
  );

  const q1_div = quantity(100, units.meter);
  const q2_div = quantity(10, units.second);
  bench(
    'division',
    () => 100.0 / 10.0,
    () => q1_div.divide(q2_div)
  );

  bench(
    'conversion',
    () => 100.0 * 0.001,
    () => q1_div.to(units.kilometer)
  );

  // Touch the sink so it is not dead-code eliminated.
  console.log(__benchSink === __benchSink ? '' : 'sink unused');
  console.log('Times are in nanoseconds per operation and depend on the host machine.');
}).catch((err) => {
  console.error('Failed to load compiled TypeScript module:', err);
  console.error('Make sure you have run `npm run build` in the typescript/ directory.');
  process.exit(1);
});
