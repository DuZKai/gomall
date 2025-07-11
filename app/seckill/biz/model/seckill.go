package model

import "time"

type SeckillRequest struct {
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
	Captcha    string `json:"captcha"`
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
	ActivityID int64  `json:"activity_id"`
	ProductID  int64  `json:"product_id"`
	Stock      int64  `json:"stock"`
	StartTime  int64  `json:"start_time"` // 时间戳，秒
	EndTime    int64  `json:"end_time"`   // 时间戳，秒
	Remark     string `json:"remark"`
}

// 活动信息结构体（数据库/Redis用）
type Activity struct {
	ID        int64  `gorm:"primaryKey;column:id"`
	ProductID int64  `gorm:"column:product_id"`
	Stock     int64  `gorm:"column:stock"`
	StartTime int64  `gorm:"column:start_time"`
	EndTime   int64  `gorm:"column:end_time"`
	Remark    string `gorm:"column:remark"`
	CreateAt  int64  `gorm:"column:create_at"`
}
