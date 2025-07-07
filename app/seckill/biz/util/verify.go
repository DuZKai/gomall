package util

import (
	"github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
)

// 模拟验证码验证函数
func verifyCaptcha(code string) bool {
	// TODO模拟：随机返回是否正确（实际应调用图形验证码或短信验证码校验）
	// h := int64(0)
	// for _, c := range code {
	// 	h += int64(c)
	// }
	h := int64(0)
	r := rand.New(rand.NewSource(h)) // h 是 int64 类型的种子
	return r.Intn(2) == 1            // 50% 几率为真
}

type SeckillRequest struct {
	UserID     string `json:"user_id"`
	ActivityID string `json:"activity_id"`
	Captcha    string `json:"captcha"`
}

func SeckillRequestHandler(c *gin.Context) {
	// Step 1: 验证资格接口
	var req SeckillRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "参数错误: " + err.Error()})
		return
	}
	userID := req.UserID
	activityID := req.ActivityID
	captcha := req.ActivityID

	if userID == "" || activityID == "" || captcha == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Missing parameters"})
		return
	}

	// if !verifyCaptcha(captcha) {
	// 	c.JSON(http.StatusForbidden, gin.H{"error": "Invalid captcha"})
	// 	return
	// }
	// TODO 事实应该判断是否同一个人

	// 第二步：限流判断
	entry, blockErr := api.Entry("seckill_request", api.WithTrafficType(base.Inbound))
	if blockErr != nil {
		// 被限流
		c.JSON(http.StatusTooManyRequests, gin.H{"error": "Too many requests - rate limited"})
		return
	}
	defer entry.Exit()

	// 若通过限流，可继续后续逻辑（如 Kafka 入队）
	c.JSON(http.StatusOK, gin.H{"message": "Captcha valid, passed rate limit check"})

}
