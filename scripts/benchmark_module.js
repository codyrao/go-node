'use strict';

const fs = require('fs');
const path = require('path');

function parseArgs(argv) {
  const args = {};
  for (const raw of argv) {
    if (!raw.startsWith('--')) {
      continue;
    }
    const body = raw.slice(2);
    const index = body.indexOf('=');
    if (index === -1) {
      args[body] = 'true';
      continue;
    }
    const key = body.slice(0, index);
    const value = body.slice(index + 1);
    args[key] = value;
  }
  return args;
}

function percentile(values, ratio) {
  if (values.length === 0) {
    return 0;
  }
  const sorted = [...values].sort((a, b) => a - b);
  const idx = Math.max(0, Math.min(sorted.length - 1, Math.ceil(sorted.length * ratio) - 1));
  return sorted[idx];
}

function runBenchmarkCase(name, fn, warmup, iterations) {
  for (let i = 0; i < warmup; i++) {
    fn(i);
  }

  if (global.gc) {
    global.gc();
  }

  const samples = new Array(iterations);
  let totalNs = 0;
  for (let i = 0; i < iterations; i++) {
    const start = process.hrtime.bigint();
    fn(i);
    const end = process.hrtime.bigint();
    const elapsed = Number(end - start);
    samples[i] = elapsed;
    totalNs += elapsed;
  }

  const avgNs = totalNs / iterations;
  const p50Ns = percentile(samples, 0.5);
  const p95Ns = percentile(samples, 0.95);
  const ops = avgNs > 0 ? 1e9 / avgNs : 0;

  return {
    name,
    iterations,
    warmup,
    avgNs,
    p50Ns,
    p95Ns,
    ops
  };
}

function ensureFunction(moduleObj, name) {
  if (typeof moduleObj[name] !== 'function') {
    throw new Error(`module function "${name}" not found`);
  }
}

function main() {
  const args = parseArgs(process.argv.slice(2));

  if (!args.module) {
    throw new Error('missing required argument: --module=/abs/path/to/module.node');
  }

  const modulePath = path.resolve(args.module);
  const iterations = Number(args.iterations || 5000);
  const warmup = Number(args.warmup || 500);
  const outputPath = args.output ? path.resolve(args.output) : '';

  if (!Number.isFinite(iterations) || iterations <= 0) {
    throw new Error('iterations must be a positive number');
  }
  if (!Number.isFinite(warmup) || warmup < 0) {
    throw new Error('warmup must be a non-negative number');
  }

  let addon;
  try {
    addon = require(modulePath);
  } catch (error) {
    const hint = [
      `failed to load module: ${modulePath}`,
      'this usually means Node ABI mismatch between current node runtime and compiled .node binary',
      `original error: ${error.message}`
    ].join('\n');
    throw new Error(hint);
  }
  ensureFunction(addon, 'hello1');
  ensureFunction(addon, 'processObject');

  const sanityHello = addon.hello1({ name: 'sanity', value: 21 });
  if (!sanityHello || typeof sanityHello.result === 'undefined') {
    throw new Error('sanity check failed for hello1');
  }
  const sanityProcess = addon.processObject({ name: 'sanity', age: 18, items: [1, 2, 3] });
  if (!sanityProcess || sanityProcess.processed !== true) {
    throw new Error('sanity check failed for processObject');
  }

  const results = [
    runBenchmarkCase(
      'hello1_object',
      (i) => addon.hello1({ name: 'bench', value: i }),
      warmup,
      iterations
    ),
    runBenchmarkCase(
      'process_object',
      (i) => addon.processObject({ name: 'bench', age: (i % 50) + 10, items: [1, 2, 3, i] }),
      warmup,
      iterations
    )
  ];

  const report = {
    timestamp: new Date().toISOString(),
    nodeVersion: process.version,
    modulePath,
    iterations,
    warmup,
    benchmarks: results
  };

  const text = JSON.stringify(report, null, 2);
  if (outputPath) {
    fs.mkdirSync(path.dirname(outputPath), { recursive: true });
    fs.writeFileSync(outputPath, text, 'utf8');
  }
  process.stdout.write(text + '\n');
  process.exit(0);
}

main();
