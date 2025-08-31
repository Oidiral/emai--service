package helpers

import (
	"context"
	"fmt"
	"sync"

	"github.com/Oidiral/emai--service/internal/config"
	rmqc "github.com/Oidiral/rmq"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

type SubscribeDeliveryHandlerFn func(ctx context.Context, eventBody []byte, enentHeaders map[string]string) (ack bool, err error)

func SubscribeRMQ(ctx context.Context, consumer rmqc.RMQConsumer, name string, fn SubscribeDeliveryHandlerFn, cfg *config.Config) (err error) {

	loge := logrus.WithContext(ctx).WithField("consumer", name)
	var wg sync.WaitGroup
	jobs := make(chan amqp091.Delivery)

	for i := 0; i < cfg.MaxWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				case d, ok := <-jobs:
					if !ok {
						return
					}
					headers := DeliveryHeadersToMap(d)
					tctx, cancel := context.WithTimeout(ctx, cfg.WorkerTimeout)
					ack, err := fn(tctx, d.Body, headers)
					cancel()
					if err != nil {
						loge.WithError(err).Error("Failed to handle event")
					}
					if ack {
						d.Ack(false)
					} else {
						d.Reject(false)
					}
				}
			}
		}()
	}

	go func() {
		defer close(jobs)
		for {
			select {
			case <-ctx.Done():
				return
			case dlvr, ok := <-consumer.Deliveries():
				if !ok {
					loge.Warn("Consumer closed")
					return
				}
				jobs <- dlvr
			}
		}
	}()
	wg.Wait()
	return
}

func DeliveryHeadersToMap(d amqp091.Delivery) map[string]string {
	m := make(map[string]string, len(d.Headers))
	for k, v := range d.Headers {
		switch val := v.(type) {
		case string:
			m[k] = val
		case fmt.Stringer:
			m[k] = val.String()
		default:
			m[k] = fmt.Sprintf("%v", v) // Fallback.
		}
	}
	return m
}
