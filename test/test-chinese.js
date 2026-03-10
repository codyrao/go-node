const hello = require('./output/hello.node')

console.log('=== 测试中文编码 ===\n')

console.log('1. 测试 Hello1 - 基本对象参数（中文）：')
const result1 = hello.hello1({ name: '守望先锋', value: 21 })
console.log('Result:', result1)
console.log('name:', result1.name)
console.log('value:', result1.value)
console.log('result:', result1.result)
console.log()

console.log('2. 测试 ProcessObject - 对象处理（中文）：')
const result2 = hello.processObject({
    name: '守望先锋',
    chineseName: '守望先锋',
    age: 25,
    items: ['游戏1', '游戏2', '游戏3']
})
console.log('Result:', result2)
console.log('name:', result2.name)
console.log('chineseName:', result2.chineseName)
console.log('nameLength:', result2.nameLength)
console.log()

console.log('3. 测试 HelloWithCallback - 同步回调（中文）：')
const result3 = hello.helloWithCallback({ test: '你好，世界' }, (data) => {
    console.log('Callback received:', data)
    console.log('test:', data.test)
    console.log('result:', data.result)
})
console.log('Result:', result3)
console.log()

console.log('=== 所有测试完成 ===')
