# Step 06: Making things real(time)

---

在本节课中，将学到：
* 使用 goroutines
* 使用匿名函数(lambdas)
* 使用 channel
* 使用 select 语句同步读取 channel
* 使用 time 包


## 概述

---

当前有一个问题：只有当玩家移动时，敌人才会移动。

出现这个问题是因为读取输入是一个阻塞操作。我们需要以某种方式使其异步运行。


## Task 01: 重构输入代码

---

