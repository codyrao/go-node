const demoaddon = require('./output/hello.node')

console.log('=== Go2Node 多函数测试 ===\n')

// 测试1: 同步调用 hello1(name string, value int) int
console.log('测试1: Hello1 - 同步调用 (name string, value int) -> int')
const result1 = demoaddon.Hello1("Test1", "10")
console.log('   结果:', result1)
console.log()

// 测试2: 同步调用 hello2(name string, value float) float
console.log('测试2: Hello2 - 同步调用 (name string, value float) -> float')
const result2 = demoaddon.Hello2("Test2", "3.14")
console.log('   结果:', result2)
console.log()

// 测试3: 同步调用 hello3(name bool, value float) bool
console.log('测试3: Hello3 - 同步调用 (name bool, value float) -> bool')
const result3 = demoaddon.Hello3("true", "5.5")
console.log('   结果:', result3)
console.log()

// 测试4: 同步调用-同步回调 hello4(name, callback(test string, res string)) int
console.log('测试4: Hello4 - 同步调用-同步回调')
const params4 = {
    test: "Hello4 Test"
}
const result4 = demoaddon.Hello4(JSON.stringify(params4), function(callbackType, jsonData) {
    console.log('   回调 [Hello4]:', callbackType, '->', jsonData)
})
console.log('   结果:', result4)
console.log()

// 测试5: 同步调用-异步回调 hello5(name, async_callback(test string, res string)) string
console.log('测试5: Hello5 - 同步调用-异步回调')
const params5 = {
    test: "Hello5 Test"
}
const result5 = demoaddon.Hello5(JSON.stringify(params5), function(callbackType, jsonData) {
    console.log('   回调 [Hello5]:', callbackType, '->', jsonData)
})
console.log('   结果:', result5)
console.log('   等待异步回调...\n')

// 测试6: 同步调用-异步无限次回调 hello6(name, async_callback(test string, res string)) string
console.log('测试6: Hello6 - 同步调用-异步无限次回调')
const params6 = {
    test: "Hello6 Test"
}
const result6 = demoaddon.Hello6(JSON.stringify(params6), function(callbackType, jsonData) {
    console.log('   回调 [Hello6]:', callbackType, '->', jsonData)
})
console.log('   结果:', result6)
console.log('   等待无限异步回调...\n')

// 保持程序运行
// let checkCount = 0
// const checkInterval = setInterval(() => {
//     checkCount++
//     console.log(`   等待中... (${checkCount}秒)`)
//     if (checkCount > 30) { // 30秒超时
//         console.log('   超时，退出')
//         clearInterval(checkInterval)
//         process.exit(0)
//     }
// }, 1000)

// console.log('=== 测试完成 ===')
