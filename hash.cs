// csc hash.cs
// origin: http://bbs.csdn.net/topics/310249367
// 2014-9-23

using System;
using System.Collections.Generic;
namespace tree
{
    class Program
    {
        static void Main(string[] args)
        {
            //模拟输入条件,随机生成200万个不同长度的字符串,长度在3到6之间
            //有谁还可以写出更快的吗？可以拿这个输入条件做测试。
            //200万个字符串，需要280毫秒，不包括显示在内。

            int totalStr = 200*10000;
            List<string> myArr = new List<string>();
            int totalLength = 0;
            char[] input = "白日依山尽黄河入海流欲穷千里目更上一层楼危楼高百尺可以摘星辰不感高声语恐惊天上人".ToCharArray();
            Random rand = new Random();
            for (int i = 0; i < totalStr; i++) {//一千个数据
                string tempStr = "";
                int tempLength = rand.Next(3, 16);
                for (int j = 0; j < tempLength; j++) {
                    tempStr += input[rand.Next(0, input.Length)];//用a-z英文测试
                }
                totalLength += tempLength;
                myArr.Add(tempStr);
            }

            //测试开始
            long begin = System.DateTime.Now.Ticks;
            Dictionary<string, string> Dict = new Dictionary<string, string>();

            for (int i = 0; i < myArr.Count;i++ ) {
                char[] charArr = myArr[i].ToCharArray();
                Array.Sort(charArr);
                string sKey = new string(charArr);
                if (!Dict.ContainsKey(sKey)) {
                    Dict.Add(sKey, myArr[i]);
                }
            }

            long end = System.DateTime.Now.Ticks;

            Console.WriteLine("共计" + totalStr/10000 + "万条数据，数据总长度" + totalLength + "，其中" + Dict.Count + "条不重复数据");
            Console.WriteLine("完成过滤共耗时" + System.TimeSpan.FromTicks(end - begin).Milliseconds /1000.000 + "s");
        }
    }
}
