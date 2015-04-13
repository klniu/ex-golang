// javac -encoding UTF-8 hash.java
// java -Xms256m -Xmx512m -XX:MaxPermSize=256m hash
// Author: liudng@gmail.com
// 2014-9-24

import java.util.*;

public class hash {
    public static void main(String[] args) {
        int totalStr = 200*10000;
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

        //测试开始
        long begin = java.util.Calendar.getInstance().getTimeInMillis();
        Dictionary<String, String> Dict = new Hashtable<String, String>();

        for (int i = 0; i < myArr.size();i++ ) {
            char[] charArr = myArr.get(i).toCharArray();
            Arrays.sort(charArr);
            String sKey = new String(charArr);
            if (Dict.get(sKey) == null) {
                Dict.put(sKey, myArr.get(i));
            }
        }

        long end = java.util.Calendar.getInstance().getTimeInMillis();

        System.out.println("共计" + totalStr / 10000 + "万条数据，数据总长度" + totalLength + "，其中" + Dict.size() + "条不重复数据");
        System.out.println("完成过滤共耗时" + (end - begin) /1000.000 + "s");
    }
}
