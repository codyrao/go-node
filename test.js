const demoaddon = require('./output/hello.node')

console.log('=== Go2Node 多函数测试 ===\n')

// 测试1: 同步调用 Hello1(name string, value int) int
console.log('测试1: Hello1 - 同步调用 (name string, value int) -> int')
const result1 = demoaddon.Hello1("Test1", "10")
console.log('   结果:', result1)
console.log()

// 测试2: 同步调用 Hello2(name string, value float) float
console.log('测试2: Hello2 - 同步调用 (name string, value float) -> float')
const result2 = demoaddon.Hello2("Test2", "3.14")
console.log('   结果:', result2)
console.log()

// 测试3: 同步调用 Hello3(name bool, value float) bool
console.log('测试3: Hello3 - 同步调用 (name bool, value float) -> bool')
const result3 = demoaddon.Hello3("true", "5.5")
console.log('   结果:', result3)
console.log()

// 测试4: 同步调用-同步回调 Hello4(name, callback(test string, res string)) int
console.log('测试4: Hello4 - 同步调用-同步回调')
const params4 = {
    test: "Hello4 Test"
}
const result4 = demoaddon.Hello4(JSON.stringify(params4), function(callbackType, jsonData) {
    console.log('   回调 [Hello4]:', callbackType, '->', jsonData)
})
console.log('   结果:', result4)
console.log()

// 测试5: 同步调用-异步回调 Hello5(name, async_callback(test string, res string)) string
console.log('测试5: Hello5 - 同步调用-异步回调')
const params5 = {
    test: "Hello5 Test"
}
const result5 = demoaddon.Hello5(JSON.stringify(params5), function(callbackType, jsonData) {
    console.log('   回调 [Hello5]:', callbackType, '->', jsonData)
})
console.log('   结果:', result5)
console.log('   等待异步回调...\n')

// 测试6: 同步调用-异步无限次回调 Hello6(name, async_callback(test string, res string)) string
console.log('测试6: Hello6 - 同步调用-异步无限次回调 (已禁用)')
const params6 = {
    test: "Hello6 Test"
}
// const result6 = demoaddon.Hello6(JSON.stringify(params6), function(callbackType, jsonData) {
//     console.log('   回调 [Hello6]:', callbackType, '->', jsonData)
// })
// console.log('   结果:', result6)
// console.log('   等待无限异步回调...\n')

// 测试7: 返回字符串类型
console.log('测试7: ReturnString - 返回字符串类型')
const result7 = demoaddon.ReturnString("Test7", "World")
console.log('   结果:', result7)
console.log('   类型:', typeof result7)
console.log()

// 测试8: 返回整数类型
console.log('测试8: ReturnInt - 返回整数类型')
const result8 = demoaddon.ReturnInt("Test8", "10")
console.log('   结果:', result8)
console.log('   类型:', typeof result8)
console.log()

// 测试9: 返回浮点数类型
console.log('测试9: ReturnFloat - 返回浮点数类型')
const result9 = demoaddon.ReturnFloat("Test9", "3.14")
console.log('   结果:', result9)
console.log('   类型:', typeof result9)
console.log()

// 测试10: 返回布尔值类型
console.log('测试10: ReturnBool - 返回布尔值类型')
const result10 = demoaddon.ReturnBool("Test10", "5.5")
console.log('   结果:', result10)
console.log('   类型:', typeof result10)
console.log()

// 测试11: 返回对象类型
console.log('测试11: ReturnObject - 返回对象类型')
const result11 = demoaddon.ReturnObject("Test11", "30")
console.log('   结果:', result11)
console.log('   类型:', typeof result11)
console.log('   name:', result11.name)
console.log('   age:', result11.age)
console.log('   isActive:', result11.isActive)
console.log('   scores:', result11.scores)
console.log('   address:', result11.address)
console.log()

// 测试12: 返回嵌套对象类型
console.log('测试12: ReturnNestedObject - 返回嵌套对象类型')
const result12 = demoaddon.ReturnNestedObject("Test12", "100")
console.log('   结果:', result12)
console.log('   类型:', typeof result12)
console.log('   user.name:', result12.user.name)
console.log('   user.age:', result12.user.age)
console.log('   metadata:', result12.metadata)
console.log('   items:', result12.items)
console.log()

console.log('=== 基础测试完成 ===')
console.log('等待异步回调...')
