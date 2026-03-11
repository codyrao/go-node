const assert = require('assert');
const path = require('path');

const outputDir = process.env.GO_NODE_OUTPUT_DIR || 'output';
const testModule = require(path.join(__dirname, outputDir, 'test_callback_fix.node'));

async function main() {
  const results = {
    functionA: [],
    functionB: [],
    functionC: []
  };

  const finished = new Promise((resolve, reject) => {
    const startedAt = Date.now();

    const finishIfReady = () => {
      const total =
        results.functionA.length +
        results.functionB.length +
        results.functionC.length;

      if (total === 9) {
        resolve();
      } else if (Date.now() - startedAt > 4000) {
        reject(new Error(`Timed out waiting for callbacks, received ${total}/9`));
      } else {
        setTimeout(finishIfReady, 50);
      }
    };

    finishIfReady();
  });

  const resultA = testModule.functionA({ test: 'data' }, (data) => {
    results.functionA.push(data);
  });
  const resultB = testModule.functionB({ test: 'data' }, (data) => {
    results.functionB.push(data);
  });
  const resultC = testModule.functionC({ test: 'data' }, (data) => {
    results.functionC.push(data);
  });

  assert.strictEqual(resultA.function, 'FunctionA');
  assert.strictEqual(resultB.function, 'FunctionB');
  assert.strictEqual(resultC.function, 'FunctionC');

  await finished;

  assert.strictEqual(results.functionA.length, 3);
  assert.strictEqual(results.functionB.length, 3);
  assert.strictEqual(results.functionC.length, 3);

  assert(results.functionA.every((item) => item.function === 'FunctionA'));
  assert(results.functionB.every((item) => item.function === 'FunctionB'));
  assert(results.functionC.every((item) => item.function === 'FunctionC'));

  console.log('test/test_callback_fix.js passed');
  process.exit(0);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
