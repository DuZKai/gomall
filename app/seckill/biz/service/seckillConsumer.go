package service

import (
	"encoding/json"
	"github.com/IBM/sarama"
	"gomall/app/seckill/biz/model"
	"log"
)

// 实现 sarama.ConsumerGroupHandler 接口
type SeckillConsumer struct{}

func (h *SeckillConsumer) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (h *SeckillConsumer) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (h *SeckillConsumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var req model.SeckillMessage
		err := json.Unmarshal(msg.Value, &req)
		if err != nil {
			panic(err)
		}

		log.Printf("[Consumer] Processing user %s in activity %s at %s", req.UserID, req.ActivityID, req.TS)
		// TODO: 第四步：这里执行实际的业务逻辑（黑名单检查、库存判断、生成token等）

		sess.MarkMessage(msg, "")
	}
	return nil
}
