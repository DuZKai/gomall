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
