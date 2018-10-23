1. 命令本质是一个程序
2. 命令有两种调用方法
    * 交互式调用
    * 程序调用
3. /bin/bash命令(bash程序)可以执行其他命令(程序)
4. Windows下cygwin可以代替/bin/bash
5. // 睡眠100毫秒
    select {
    case <- time.NewTimer(100 * time.Millisecond).C
    }
6. time.AfterFunc(5 * time.Second, func() { ... })