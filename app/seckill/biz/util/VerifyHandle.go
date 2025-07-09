package util

import (
	"math/rand"
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
