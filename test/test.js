const assert = require('assert');
const path = require('path');

const outputDir = process.env.GO_NODE_OUTPUT_DIR || 'output';
const hello = require(path.join(__dirname, outputDir, 'hello.node'));

async function main() {
  const result1 = hello.hello1({ name: 'Alice', value: 21 });
  assert.deepStrictEqual(result1, {
    name: 'Alice',
    value: 21,
    result: 42
  });

  const result2 = hello.processObject({
    name: 'Charlie',
    age: 25,
    items: ['item1', 'item2', 'item3']
  });
  assert.strictEqual(result2.processed, true);
  assert.strictEqual(result2.name, 'Charlie');
  assert.strictEqual(result2.itemCount, 3);
  assert.strictEqual(result2.isAdult, true);

  const syncCallbacks = await new Promise((resolve, reject) => {
    const events = [];
    const result3 = hello.helloWithCallback({ test: 'sync-check' }, (data) => {
      events.push(data);
      if (events.length === 3) {
        resolve({ result: result3, events });
      }
    });

    setTimeout(() => reject(new Error('Timed out waiting for helloWithCallback events')), 4000);
  });
  assert.strictEqual(syncCallbacks.result.status, 'success');
  assert.strictEqual(syncCallbacks.events.length, 3);
  assert(syncCallbacks.events.every((item, index) => item.result === `Callback ${index + 1}`));
  assert(syncCallbacks.events.every((item) => item.test === 'sync-check'));

  const asyncCallbacks = await new Promise((resolve, reject) => {
    const events = [];

    const result4 = hello.asyncHello({ test: 'async-check' }, (data) => {
      events.push(data);
    });

    assert.strictEqual(result4.status, 'success');
    assert.strictEqual(result4.result, 'Async started');

    setTimeout(() => resolve({ result: result4, events }), 1500);
  });

  assert.strictEqual(asyncCallbacks.events.length, 0);

  console.log('test/test.js passed');
  process.exit(0);
}

main().catch((error) => {
  console.error(error);
  process.exit(1);
});
