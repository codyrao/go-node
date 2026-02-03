const koffi = require('koffi');
const path = require('path');

// DLL 路径
const dllPath = path.join(__dirname, 'output', 'hello.dll');
console.log('Loading DLL from:', dllPath);

try {
    // 加载 DLL
    const hello = koffi.load(dllPath);

    console.log('\n=== Testing DLL Functions ===\n');

    // 声明函数 - 所有函数都接受两个字符串参数，返回字符串
    const Hello1 = hello.func('Hello1', 'string', ['string', 'string']);
    const Hello2 = hello.func('Hello2', 'string', ['string', 'string']);
    const Hello3 = hello.func('Hello3', 'string', ['string', 'string']);
    const ReturnString = hello.func('ReturnString', 'string', ['string', 'string']);
    const ReturnInt = hello.func('ReturnInt', 'string', ['string', 'string']);
    const ReturnFloat = hello.func('ReturnFloat', 'string', ['string', 'string']);
    const ReturnBool = hello.func('ReturnBool', 'string', ['string', 'string']);
    const ReturnObject = hello.func('ReturnObject', 'string', ['string', 'string']);
    const ProcessArray = hello.func('ProcessArray', 'string', ['string']);
    const ProcessObject = hello.func('ProcessObject', 'string', ['string']);

    // 测试 Hello1
    console.log('Hello1("", ""):', Hello1("", ""));

    // 测试 Hello2
    console.log('Hello2("World", "123"):', Hello2("World", "123"));

    // 测试 Hello3
    console.log('Hello3("Hello", "DLL"):', Hello3("Hello", "DLL"));

    // 测试 ReturnString
    console.log('ReturnString("User", "Test"):', ReturnString("User", "Test"));

    // 测试 ReturnInt - 传入数字字符串
    console.log('ReturnInt("", "42"):', ReturnInt("", "42"));

    // 测试 ReturnFloat - 传入浮点数字符串
    console.log('ReturnFloat("", "3.14"):', ReturnFloat("", "3.14"));

    // 测试 ReturnBool - 传入 bool 字符串
    console.log('ReturnBool("", "true"):', ReturnBool("", "true"));

    // 测试 ReturnObject
    console.log('ReturnObject("Test", "Value"):', ReturnObject("Test", "Value"));

    // 测试 ProcessArray
    const testArray = JSON.stringify([1, 2, 3, 4, 5]);
    console.log('ProcessArray([1,2,3,4,5]):', ProcessArray(testArray));

    // 测试 ProcessObject
    const testObj = JSON.stringify({ name: 'test', value: 123 });
    console.log('ProcessObject({name:"test",value:123}):', ProcessObject(testObj));

    console.log('\n=== All tests passed! ===');

} catch (error) {
    console.error('Error loading or calling DLL:', error.message);
    console.error('\nStack:', error.stack);
    process.exit(1);
}
