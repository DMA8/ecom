package producer

import (
	"context"
	"ecom/loms/internal/model"
	"encoding/json"
	"fmt"
	"log"
)

const (
	topicStatusChanged string = "order_status_changed"
)

func (p *Producer) SendOrderStatusChanged(ctx context.Context, order model.OrderChangedMessage) {
	msg, err := json.Marshal(order)
	if err != nil {
		log.Println(err)
	}
	p.Send(ctx, topicStatusChanged, msg, fmt.Sprintf("%d", order.OrderID))
}
