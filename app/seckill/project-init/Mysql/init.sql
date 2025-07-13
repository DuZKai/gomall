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
CREATE TABLE activity_stocks (
                                 activity_id VARCHAR(64) PRIMARY KEY COMMENT "活动ID",
                                 stock INT NOT NULL COMMENT "活动库存"
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci COMMENT='秒杀活动库存表' COMMENT='秒杀活动库存表';
;