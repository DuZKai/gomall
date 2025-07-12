package model

import "time"

type SeckillRequest struct {
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
	Captcha    string `json:"captcha"`
	Priority   int    `json:"priority"` // 0: 普通，1: 高优先级
}

type SeckillMessage struct {
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
	TS         string `json:"ts"`
}

type CheckoutRequest struct {
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
}

type Order struct {
	ID         int64      `db:"id"`
	UserID     string     `db:"user_id"`
	ActivityID string     `db:"activity_id"`
	Status     string     `db:"status"` // 可为：INIT / PAID / TIMEOUT
	CreateTime time.Time  `db:"create_time"`
	PayTime    *time.Time `db:"pay_time,omitempty"`
}

type TokenInfo struct {
	UserID       string `json:"user_id"`
	ActivityID   string `json:"activity_id"`
	CreateTime   int64  `json:"create_time"`    // 纳秒时间戳
	ExpireSecond int64  `json:"expire_seconds"` // 秒数
}

// 秒杀活动请求结构
type CreateActivityRequest struct {
	ActivityID string `json:"activity_id"`
	ProductID  string `json:"product_id"`
	Stock      int64  `json:"stock"`
	StartTime  string `json:"start_time"` // 字符串时间，如 "2025-07-12 14:30:00"
	EndTime    string `json:"end_time"`   // 字符串时间，如 "2025-07-12 14:30:00"
	Remark     string `json:"remark"`
}

// 活动信息结构体（数据库/Redis用）
type Activity struct {
	ID         string `gorm:"primaryKey;column:id"`
	ActivityID string `json:"activity_id"`
	ProductID  string `gorm:"column:product_id"`
	Stock      int64  `gorm:"column:stock"`
	StartTime  int64  `gorm:"column:start_time"`
	EndTime    int64  `gorm:"column:end_time"`
	Remark     string `gorm:"column:remark"`
	CreateAt   int64  `gorm:"column:create_at"`
}

type ActivityStock struct {
	ActivityID string `json:"activity_id"`
	Stock      int64  `json:"stock"`
}

type SeckillConfig struct {
	ActivityTimeout string `json:"activity_timeout"`
	RateLimitQPS    int    `json:"rate_limit_qps"`
}

type SeckillLimitConfig struct {
	TokenTTL            int `json:"token_ttl"`             // Token 有效时间（分钟）
	BlacklistTTL        int `json:"blacklist_ttl"`         // 黑名单拉黑时间（分钟）
	FreqLimitExpire     int `json:"freq_limit_expire"`     // 访问频率限制过期时间（秒）
	IdempotentKeyExpire int `json:"idempotent_key_expire"` // 幂等性校验键过期时间（分钟）
	BucketExpireSeconds int `json:"bucket_expire_seconds"` // Redis 令牌桶 key 过期时间（秒）
	CapacityFactor      int `json:"capacity_factor"`       // 动态桶容量调整因子
	RateFactor          int `json:"rate_factor"`           // 动态速率调整因子
	BaseTokenRate       int `json:"base_token_rate"`       // 默认令牌生成速率（每秒）
	TokenBucketFactor   int `json:"token_bucket_factor"`   // 限流桶倍率（扩展）
}
