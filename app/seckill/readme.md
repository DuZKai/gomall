# 秒杀滴答

## 实现功能
- 前置校验：开始时间、结束时间、库存是否有多等
- 限流策略：VIP进入Sentinel，非VIP进入Redis令牌桶限流
- Kafka消息队列：通过Kafka发送到消息队列，异步处理订单创建
- 库存安全：Lua脚本封装库存校验、库存预扣减与Token下发逻辑，防止库存超卖，支持一人一单校验，提升系统稳定性。
- 库存回滚：基于 Asynq 构建统一定时分布式任务轮询器，完成过期令牌库存回滚，替代单用户延时任务方案，大幅节省资源。

## 安装
在docker里面运行如下命令构建容器
```bash
docker-compose up -d
# 检验是否完全启动成功
docker ps
# 如果需要查看日志可以使用
docker-compose logs -f
```
PS: 如果有错误，可以先手动docker pull对应镜像再执行上面的命令。

需要停止可以使用
```bash
docker-compose down
```

mysql建表语句
```mysql
CREATE TABLE `orders` (
    `id` BIGINT NOT NULL AUTO_INCREMENT COMMENT '主键ID',
    `user_id` VARCHAR(64) NOT NULL COMMENT '用户ID',
    `activity_id` VARCHAR(64) NOT NULL COMMENT '活动ID',
    `status` VARCHAR(16) NOT NULL DEFAULT 'INIT' COMMENT '订单状态: INIT / PAID / TIMEOUT',
    `create_time` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `pay_time` DATETIME DEFAULT NULL COMMENT '支付时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_activity_id` (`activity_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀订单表';
CREATE TABLE `activities` (
    `id` VARCHAR(64) NOT NULL PRIMARY KEY COMMENT '雪花ID',
    `activity_id` VARCHAR(64) NOT NULL COMMENT '业务活动ID',
    `product_id` VARCHAR(64) NOT NULL COMMENT '关联商品ID',
    `stock` BIGINT NOT NULL COMMENT '库存数量',
    `start_time` BIGINT NOT NULL COMMENT '开始时间（时间戳）',
    `end_time` BIGINT NOT NULL COMMENT '结束时间（时间戳）',
    `remark` TEXT COMMENT '备注信息',
    `create_at` BIGINT NOT NULL COMMENT '创建时间（时间戳）'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='秒杀活动表';


```


安装完成后使用如下命令启动
```bash
go run .
```

