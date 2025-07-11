# 秒杀

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
    `id` BIGINT NOT NULL PRIMARY KEY COMMENT '活动ID',
    `product_id` BIGINT NOT NULL COMMENT '商品ID',
    `stock` INT NOT NULL COMMENT '初始库存数量',
    `start_time` BIGINT NOT NULL COMMENT '活动开始时间（Unix秒）',
    `end_time` BIGINT NOT NULL COMMENT '活动结束时间（Unix秒）',
    `remark` VARCHAR(255) DEFAULT NULL COMMENT '活动备注',
    `create_at` BIGINT NOT NULL COMMENT '创建时间（Unix秒）'
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='秒杀活动表';

```


安装完成后使用如下命令启动
```bash
go run .
```

