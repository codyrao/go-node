const demoaddon = require('./output/hello.node')
var callbackCount = 0
const testArray13 = [1, 2, 3, 4, 5, "hello", "world"]

//  demoaddon.Hello6(testArray13, function(p1,res){
//     callbackCount++
    
//     console.log(`>>> 回调 #${callbackCount} 被调用!`)
//     console.log('>>> success:', res.success)
//     console.log('>>> error:', res.error)
//     console.log('>>> res:', res)

// })

 demoaddon.ProcessArray(testArray13)