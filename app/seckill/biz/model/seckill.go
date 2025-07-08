package model

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
