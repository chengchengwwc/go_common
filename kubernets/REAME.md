## kubernetes


#### 云计算


#### kubernetes
##### 声明试系统
kubernets的所有管理能力构建在对象抽象的基础上，核心对象包括
- node: 计算节点的抽象，用来描述计算节点的资源抽象，健康状态
- namespace: 资源隔离的基本单位，可以简单理解为文件系统中的目录结构
- pod:用来描述应用实例，包括镜像地址，资源需求等。kubernetes中最核心的对象，也是打通应用和基础架构的秘密武器
- service:服务如何将应用发布成服务，本质上是负载均衡和域名服务的声明

##### 主节点
- api server: 这是控制板面中唯一带有用户访问许可的API以及用户交互组件，API服务器会暴漏一个restful的kubernets api并使用JSON格式的清单文件
- cluster data store: kubernetes 使用etcd。这是一个强大的稳定的，高可用的键值存储，被kubernetes用于长久存储所有的API对象
- controller manager: kube-controller manager 它运行着所有的处理集群日常任务的控制器，包括节点控制器，副本控制器，端点控制器以及服务账户
- scheduler: 调度器会监控新建的pods(一个组件或是一个容器)并将其分配给其他节点。

##### etcd
- 基于key-value
- 监听机制
- key的过期以及续约机制，用于监听和服务发现
- 原子

##### API server
kube-APIServer是kubernetes最重要的核心组件之一，主要提供以下功能
- 提供集群管理的Rest Api接口，包括
   - 认证
   - 授权
   - 准入
- 提供其他模块之间的数据交互和通信的枢纽（其他模块通过API Server查询或是修改数据，只有API server可以直接操作etcd）
- APIServer提供etcd数据缓存，减少集群对etcd的访问


##### controller manager
- controller manager 是集群的大脑，是确保整个集群动起来的关键
- 确保kubernetes是遵循声明系统规范，确保系统真实状态和用户定义的期望状态一致
- 是多个控制器的组合，每个controller其实是一个control loop,负责侦听其管控的对象，当对象发送变更的时候完成配置
- controller配置失败，通常会触发自动重试，整个集群会在控制器不断重试的机制下保持一致


