# go-docker

简单来说 docker 本质其实是一个**特殊的进程**

这个进程特殊在它被`Namespace`和`Cgroup`技术做了装饰

- Namespace 

  将该进程与 Linux 系统进行隔离开来，让该进程处于一个虚拟的沙盒中

- Cgroup 

  对该进程做了一系列的资源限制，两者配合模拟出来一个沙盒的环境

## Namespace

> Linux 对线程提供了六种隔离机制，分别为：`uts` `pid` `user` `mount` `network` `ipc` ，它们的作用如下：

- `uts`: 用来隔离主机名
- `pid`：用来隔离进程 PID 号的
- `user`: 用来隔离用户的
- `mount`：用来隔离各个进程看到的挂载点视图
- `network`: 用来隔离网络
- `ipc`：用来隔离 System V IPC 和 POSIX message queues

## Cgroup

**Linux Cgroup** 提供了对一组进程及子进程的资源限制，控制和统计的能力

- 这些资源包括 CPU，内存，存储，网络等

> 通过 Cgroup可以方便的吸纳之某个进程的资源占用，并且可以实时监控进程和统计信息。

`Cgroup`完成资源限制主要通过下面三个组件

- `cgroup`: 是对**进程分组管理**的一种机制
- `subsystem`: 是一组**资源控制**的模块
- `hierarchy`: 把一组 cgroup 串成一个树状结构 (可让其实现继承)

### 使用

在Linux中使用`Cgroup`可以通过创建文件使其挂载`hierarchy`

一旦挂载后就会生成以下一些默认文件：

- `cgroup.clone_children`：cpuset 的 subsystem 会读取该文件，如果该文件里面的值为 1 的话，那么子 cgroup 将会继承父 cgroup 的 cpuset 配置
- `cgroup.procs`：记录了树中当前节点 cgroup 中的进程组 ID
- `task`: 标识该 cgroup 下的进程 ID，如果将**某个进程的 ID** 写到该文件中，那么便会将该进程**加入到当前的 cgroup** 中。

而新建其子`Cgroup`则在其父`Cgroup`文件夹下新建一个文件则该文件就会自动被`kernel`标记为该`Cgroup`为子`Cgroup`,且同样其子`Cgroup`下会自动生成默认文件，且会继承其父的配置。

我们提到的**限制进程的资源**则需要关联到`subsystem`，而系统默认已经为每个`subsystem`创建了一个默认的`hierarchy`，它在Linux的`/sys/fs/cgroup`路径下。其路径下有各种资源文件。

![](http://img.zhengyua.cn/img/20200831172200.png)

而怎么来限制呢？我们只需要在其对应资源的文件下创建一个文件夹，`kernel`会自动将该文件夹标记为一个`Cgroup`。

比如我们想要限制某个进程的内存，我们就需要在`/sys/fs/cgroup/memory`下创建一个限制memory的`cgroup`，创建`cgroup-demo-memory`文件夹后系统会自动帮助我们生成默认文件(限制资源文件)。
![](http://img.zhengyua.cn/img/20200831172445.png)

我们只需要将**进程ID**写入`tasks`文件中，然后**修改**其`memory.limit_in_bytes`的文件，就能够实现限制该进程的内存使用。











## 参考

https://learnku.com/articles/42072