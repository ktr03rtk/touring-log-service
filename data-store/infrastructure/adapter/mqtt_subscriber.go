package adapter

import (
	"context"
	"crypto/tls"
	"log"
	"net/url"
	"time"

	"github.com/eclipse/paho.golang/autopaho"
	"github.com/eclipse/paho.golang/paho"
	"github.com/ktr03rtk/touring-log-service/data-store/domain/model"
	"github.com/ktr03rtk/touring-log-service/data-store/domain/repository"
	"github.com/pkg/errors"
)

const (
	connectionWaitTime = 1000 // milliseconds
	concurrency        = 100
)

type mqttAdapter struct {
	connectionManager *autopaho.ConnectionManager
	payloadCh         chan *paho.Publish
	topic             string
	qos               byte
}

func NewMqttAdapter(ctx context.Context) (repository.PayloadSubscribeRepository, error) {
	cfg, err := getConfig()
	if err != nil {
		return nil, err
	}

	tlsCfg, err := newTLSConfig()
	if err != nil {
		return nil, err
	}

	ch := make(chan *paho.Publish, concurrency)
	mqttCfg := getMqttConfig(cfg, tlsCfg, ch)

	cm, err := autopaho.NewConnection(ctx, mqttCfg)
	if err != nil {
		return nil, err
	}

	time.Sleep(connectionWaitTime * time.Millisecond)

	return &mqttAdapter{
		connectionManager: cm,
		payloadCh:         ch,
		topic:             cfg.topic,
		qos:               cfg.qos,
	}, nil
}

func (a *mqttAdapter) Subscribe(ctx context.Context, ch chan<- *model.Payload) error {
	defer close(ch)
	defer close(a.payloadCh)

	for {
		select {
		case <-ctx.Done():

			return nil
		case p := <-a.payloadCh:
			payload, err := model.NewPayload(p.Payload, p.Topic)
			if err != nil {
				return errors.Wrapf(err, "failed to convert mqtt data")
			}

			ch <- payload
		}
	}
}

func getMqttConfig(cfg config, tlsCfg *tls.Config, ch chan<- *paho.Publish) autopaho.ClientConfig {
	return autopaho.ClientConfig{
		BrokerUrls:        []*url.URL{cfg.serverURL},
		KeepAlive:         cfg.keepAlive,
		ConnectRetryDelay: cfg.connectRetryDelay,
		OnConnectionUp: func(cm *autopaho.ConnectionManager, connAck *paho.Connack) {
			log.Println("mqtt connection up")
			if _, err := cm.Subscribe(context.Background(), &paho.Subscribe{
				Subscriptions: map[string]paho.SubscribeOptions{
					cfg.topic: {QoS: cfg.qos},
				},
			}); err != nil {
				log.Printf("failed to subscribe (%s). This is likely to mean no messages will be received.", err)
				return
			}
			log.Println("mqtt subscription made")
		},
		OnConnectError: func(err error) { log.Printf("error whilst attempting connection: %s\n", err) },
		Debug:          paho.NOOPLogger{},
		TlsCfg:         tlsCfg,
		ClientConfig: paho.ClientConfig{
			ClientID: cfg.clientID,
			Router: paho.NewSingleHandlerRouter(func(m *paho.Publish) {
				ch <- m
			}),
			OnClientError: func(err error) { log.Printf("server requested disconnect: %s\n", err) },
			OnServerDisconnect: func(d *paho.Disconnect) {
				if d.Properties != nil {
					log.Printf("server requested disconnect: %s\n", d.Properties.ReasonString)
				} else {
					log.Printf("server requested disconnect; reason code: %d\n", d.ReasonCode)
				}
			},
		},
	}
}
