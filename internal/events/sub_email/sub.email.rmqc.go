package sub_email

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Oidiral/emai--service/internal/config"
	"github.com/Oidiral/emai--service/internal/events/helpers"
	rmqc "github.com/Oidiral/rmq"
	"github.com/sirupsen/logrus"
)

type EmailSubscribeRMQ struct {
	wg      sync.WaitGroup
	send    rmqc.RMQConsumer
	handler EmailHandlerInterface
	cfg     *config.Config
}

func NewEmailSubscribeRMQ(handler EmailHandlerInterface, cfg *config.Config) *EmailSubscribeRMQ {
	return &EmailSubscribeRMQ{handler: handler, cfg: cfg}
}

func (s *EmailSubscribeRMQ) SetSend(consumer rmqc.RMQConsumer) {
	s.send = consumer
}

func (s *EmailSubscribeRMQ) Subscribe(ctx context.Context) {
	loge := logrus.WithContext(ctx).WithField("consumer", "send")

	if s.send == nil {
		loge.Error("send consumer not set")
		return
	}

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		if err := helpers.SubscribeRMQ(ctx, s.send, "send", s.handleMessage, s.cfg); err != nil {
			loge.WithError(err).Error("subscribe failed")
		}
	}()

	s.wg.Wait()
}

func (s *EmailSubscribeRMQ) handleMessage(ctx context.Context, eventBody []byte, eventHeaders map[string]string) (ack bool, err error) {
	loge := logrus.WithContext(ctx)
	bodyLog := string(eventBody)
	headersLog := eventHeaders

	if xDeathRaw, ok := eventHeaders["x-death"]; ok {
		var deaths []map[string]interface{}
		if err := json.Unmarshal([]byte(xDeathRaw), &deaths); err != nil {
			loge.WithError(err).Warn("failed to parse x-death")
			return true, nil
		}
		if len(deaths) > 0 {
			if countRaw, ok := deaths[0]["count"]; ok {
				count, ok := countRaw.(float64)
				if ok && int(count) > s.cfg.MaxRetries {
					loge.WithFields(logrus.Fields{
						"retry_count":  count,
						"eventBody":    bodyLog,
						"eventHeaders": headersLog,
					}).Error("max retries exceeded, discarding message")
					// TODO: Optional - publish to parking queue.
					return true, nil
				}
				loge.WithField("retry_count", count).Info("processing retry from DLQ")
			}
		}
	}

	e := &SendEvent{}
	if err := json.Unmarshal(eventBody, e); err != nil {
		loge.WithError(err).Error("unmarshal failed")
		return false, fmt.Errorf("unmarshal: %w", err)
	}

	if e.Email == "" || e.Template == "" || e.Locale == "" {
		loge.WithField("event", e).Error("invalid event: missing required fields")
		return false, fmt.Errorf("invalid event")
	}

	if err := s.handler.Send(ctx, e); err != nil {
		loge.WithError(err).Error("handler failed")
		return false, err
	}

	return true, nil
}
