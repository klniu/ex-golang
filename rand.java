// javac -encoding UTF-8 rand.java
// java -Xms256m -Xmx1300m rand
// Author: liudng@gmail.com
// 2015-9-22

import java.util.*;

public class rand {
    public static void main(String[] args) {
        //Begin
        long begin = java.util.Calendar.getInstance().getTimeInMillis();
        
        int totalStr = 600*10000;
        List<String> myArr = new ArrayList<String>();
        int totalLength = 0;
        char[] input = "白日依山尽黄河入海流欲穷千里目更上一层楼危楼高百尺可以摘星辰不感高声语恐惊天上人".toCharArray();
        Random rand = new Random();
        for (int i = 0; i < totalStr; i++){//一千个数据
            String tempStr = "";
            int tempLength = rand.nextInt(13) + 3;
            for (int j = 0; j < tempLength; j++) {
                tempStr += input[rand.nextInt(input.length)];//用a-z英文测试
            }
            totalLength += tempLength;
            myArr.add(tempStr);
        }

        long end = java.util.Calendar.getInstance().getTimeInMillis();

        System.out.println("共计" + totalStr / 10000 + "万条数据，数据总长度" + totalLength);
        System.out.println("完成过滤共耗时" + (end - begin) /1000.000 + "s");
    }
}
