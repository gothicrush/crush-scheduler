## 架构图

![](C:\Users\narli\Desktop\gopath\src\github.com\gothicrush\crush-scheduler\_image\0.PNG)



## web管理样式图

![](C:\Users\narli\Desktop\gopath\src\github.com\gothicrush\crush-scheduler\_image\1.PNG)

![](C:\Users\narli\Desktop\gopath\src\github.com\gothicrush\crush-scheduler\_image\2.PNG)

![](C:\Users\narli\Desktop\gopath\src\github.com\gothicrush\crush-scheduler\_image\3.PNG)



## Master

### Master

* 功能：初始化各个模块
* 组成：
  * 设置线程数目
  * 处理命令行参数
  * 加载配置
  * 初始化服务发现
  * 初始化日志管理器
  * 初始化任务管理器
  * 初始化 HTTP API 服务

### API Server

* 功能：提供HTTP API服务
* 组成：
  * 保存任务接口
  * 删除任务接口
  * 列出所有任务接口
  * 强制删除某个任务接口
  * 查询日志接口
  * 获取健康worker节点接口
  * 初始化API Server

### Job Manager

* 功能：对任务进行管理
* 组成：
  * 保存任务
  * 删除任务
  * 列举所有任务
  * 强杀任务
  * 初始化 Job Manager

### Log Manager

* 功能：对日志进行管理
* 组成：
  * 列举所有日志
  * 初始化 Log Manager

### Worker Manager

* 功能：用于服务发现
* 组成：
  * 列举所有的系统中所有的worker
  * 初始化 Worker Manager

### Config

* 功能：存储 master 的全局配置
* 组成：初始化 Config



## Worker

### Worker

### Executor

* 执行一个任务
* 初始化执行器

### Job Lock

* 功能：实现分布式锁
* 组成：
  * 锁结构定义
  * 创建锁
  * 上锁操作
    * 创建租约，自动续租
    * 使用事务抢锁
  * 释放锁操作

### Job Manager

* 初始化任务管理器
* 监听任务变化，派发任务给执行器
* 创建任务执行锁
* 监听强杀任务

### Log Sink

* 功能：日志记录
* 组成：
  * 日志存储
  * 批量写入日志
  * 发送日志
  * 初始化日志模块

### Register

* 功能：服务注册
* 组成：
  * 获取本机网卡IP地址，作为机器唯一标识
  * 服务注册
  * 初始化服务注册模块

### Scheduler

* 初始化调度器
* 处理任务事件
* 处理任务结果
* 调度协程
* 尝试执行任务
* 尝试执行，获取休眠时间
* 推送任务到执行器
* 初始化调度器
* 回传任务执行结果

### Config

- 功能：存储 master 的全局配置
- 组成：初始化 Config



## Common

* 任务 Job
* 反序列化 Job
* 任务事件
* 构建任务事件
* 任务执行状态
* 构建任务执行信息
* 任务调度计划
* 构建任务计划
* 任务执行结果
* HTTP接口应答
* 应答方法
* 任务执行日志
* 任务日志批次
* 任务日志排序规则
* 提取任务名



