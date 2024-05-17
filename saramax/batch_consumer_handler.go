package saramax

import (
	"context"
	"encoding/json"
	"github.com/IBM/sarama"
	"time"
)

// BatchHandler 批量消费
type BatchHandler[T any] struct {
	fn func(msgs []*sarama.ConsumerMessage, ts []T) error
	// 用 option 模式来设置这个 batchSize 和 duration
	batchSize     int
	batchDuration time.Duration
}

func (b *BatchHandler[T]) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchHandler[T]) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	msgsCh := claim.Messages()
	batchSize := b.batchSize
	for {
		ctx, cancel := context.WithTimeout(context.Background(), b.batchDuration)
		done := false
		msgs := make([]*sarama.ConsumerMessage, 0, batchSize)
		ts := make([]T, 0, batchSize)
		for i := 0; i < batchSize && !done; i++ {
			select {
			case <-ctx.Done():
				done = true
			case msg, ok := <-msgsCh:
				// 再按照 key 或者业务 ID 转发到不同的 channel
				if !ok {
					cancel()
					// 代表消费者被关闭了
					return nil
				}
				var t T
				err := json.Unmarshal(msg.Value, &t)
				if err != nil {

					continue
				}
				msgs = append(msgs, msg)
				ts = append(ts, t)
			}
		}
		cancel()
		if len(msgs) == 0 {
			continue
		}
		err := b.fn(msgs, ts)
		if err != nil {
			// 你这里整个批次都要记下来

			// 还要继续往前消费
		}
		for _, msg := range msgs {
			// 这样，万无一失
			session.MarkMessage(msg, "")
		}
	}
}
