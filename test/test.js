const hello = require('./output/hello.node');

console.log('=== Testing go-node with new structure ===\n');

console.log('1. Test Hello1 - basic object parameter:');
const result1 = hello.hello1({ name: 'Alice', value: 21 });
console.log('Result:', result1);
console.log();

console.log('2. Test Hello2 - float calculation:');
const result2 = hello.hello2({ name: 'Bob', value: 10.5 });
console.log('Result:', result2);
console.log();

console.log('3. Test ProcessObject - object processing:');
const result3 = hello.processObject({
	name: 'Charlie',
	age: 25,
	items: ['item1', 'item2', 'item3']
});
console.log('Result:', result3);
console.log();

console.log('4. Test Calculate - different operations:');
const result4a = hello.calculate({ a: 10, b: 5, operation: 'add' });
console.log('Add (10 + 5):', result4a);

const result4b = hello.calculate({ a: 10, b: 5, operation: 'subtract' });
console.log('Subtract (10 - 5):', result4b);

const result4c = hello.calculate({ a: 10, b: 5, operation: 'multiply' });
console.log('Multiply (10 * 5):', result4c);

const result4d = hello.calculate({ a: 10, b: 5, operation: 'divide' });
console.log('Divide (10 / 5):', result4d);
console.log();

console.log('5. Test HelloWithCallback - synchronous callback:');
const result5 = hello.helloWithCallback({ test: 'Hello from Node' }, (data) => {
	console.log('Callback received:', data);
});
console.log('Result:', result5);
console.log();

console.log('6. Test AsyncHello - asynchronous callback:');
const result6 = hello.asyncHello({ test: 'Async test' }, (data) => {
	console.log('Async callback received:', data);
});
console.log('Result:', result6);
console.log();

console.log('7. Test NoReturn - no return value:');
const result7 = hello.noReturn({});
console.log('Result:', result7);
console.log();

console.log('=== All tests completed ===');
