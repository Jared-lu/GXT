package saramax

import (
	"encoding/json"
	"github.com/IBM/sarama"
)

type Handler[T any] struct {
	// 如何提供一个通用的日志输出，让调用方传入给我
	//l  func(msg string, args ...interface{})
	fn func(msg *sarama.ConsumerMessage, t T) error
}

func (h *Handler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *Handler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		var t T
		err := json.Unmarshal(msg.Value, &t)
		if err != nil {
			continue
		}
		for i := 0; i < 3; i++ {
			err = h.fn(msg, t)
			if err == nil {
				break
			}
			// 记录日志

		}
		if err != nil {
			// 重试次数达到上限
		} else {
			session.MarkMessage(msg, "")
		}
	}
	return nil
}
