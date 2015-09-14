// node hash.js
// Author: liudng@gmail.com
// 2014-9-25

var totalStr = 1000*10000
var myArr = []
var totalLength = 0
var input = "白日依山尽黄河入海流欲穷千里目更上一层楼危楼高百尺可以摘星辰不感高声语恐惊天上人".split("")
var i = 0

for (i = 0; i < totalStr; i++){
    var tempStr = ""
    var tempLength = parseInt(Math.random()*13+3)
    for (var j = 0; j < tempLength; j++){
        tempStr += input[parseInt(Math.random()*input.length)]
    }
    totalLength += tempLength
    myArr.push(tempStr)
}

//测试开始
var begin = new Date()
Dict = {}
var uniq = 0
for (i = 0; i < myArr.length; i++ ) {
    var charArr = myArr[i].split("")
    charArr.sort()
    var sKey = charArr.join("")
    if (typeof Dict[sKey] == "undefined") {
        Dict[sKey] = myArr[i]
        uniq++
    }
}

var end = new Date()

console.log("共计" + totalStr / 10000 + "万条数据，数据总长度" + totalLength + "，其中" + uniq + "条不重复数据")
console.log("完成过滤共耗时" + (end.getTime() - begin.getTime()) /1000.000 + "s")
