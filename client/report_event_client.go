package client

import (
	"github.com/ONSdigital/go-ns/log"
	"fmt"
	"github.com/ONSdigital/dp-reporter-client/model"
	"errors"
	"context"
	"github.com/ONSdigital/dp-reporter-client/schema"
	"time"
)

const (
	instanceEmpty    = "cannot send reportEvent as instanceID is empty"
	sendingEvent     = "sending reportEvent for application error"
	failedToMarshal  = "failed to marshal reportEvent to avro"
	eventMessageFMT  = "%s: %s"
	serviceNameEmpty = "cannot create new reporter client as serviceName is empty"
	kafkaProducerNil = "cannot create new reporter client as kafkaProducer is nil"
	eventTypeErr     = "error"
	reportEventKey   = "reportEvent"
	defaultTimeout   = 10
)

// KafkaProducer interface of the go-ns kafka.Producer
type KafkaProducer interface {
	Output() chan []byte
	Close(ctx context.Context) (err error)
}

type marshalFunc func(s interface{}) ([]byte, error)

//ReporterClient a client for sending error reports to the import-reporte
type ReporterClient struct {
	kafkaProducer KafkaProducer
	marshal       marshalFunc
	serviceName   string
}

// NewReporterClientMock create a new ReporterClient to send error reports to the import-reporter
func NewReporterClient(kafkaProducer KafkaProducer, serviceName string) (*ReporterClient, error) {
	if kafkaProducer == nil {
		return nil, errors.New(kafkaProducerNil)
	}
	if len(serviceName) == 0 {
		return nil, errors.New(serviceNameEmpty)
	}
	cli := &ReporterClient{
		serviceName:   serviceName,
		kafkaProducer: kafkaProducer,
		marshal:       schema.ReportEventSchema.Marshal,
	}
	return cli, nil
}

// ReportError send an error report to the import-reporter
func (c *ReporterClient) ReportError(instanceID string, eventContext string, err error, data log.Data) error {
	log.ErrorC(eventContext, err, data)

	if len(instanceID) == 0 {
		log.Info(instanceEmpty, nil)
		return nil
	}

	reportEvent := &model.ReportEvent{
		InstanceID:  instanceID,
		EventMsg:    eventMsg(eventContext, err),
		ServiceName: c.serviceName,
		EventType:   eventTypeErr,
	}

	reportEventData := log.Data{reportEventKey: reportEvent}
	log.Info(sendingEvent, reportEventData)

	avroBytes, err := c.marshal(reportEvent)
	if err != nil {
		log.ErrorC(failedToMarshal, err, reportEventData)
		return err
	}

	c.kafkaProducer.Output() <- avroBytes
	return nil
}

// Close properly closes the ReporterClient
func (c *ReporterClient) Close(ctx context.Context) error {
	if ctx == nil {
		ctx, _ = context.WithTimeout(context.Background(), time.Second*defaultTimeout)
	}
	if _, ok := ctx.Deadline(); !ok {
		ctx, _ = context.WithTimeout(ctx, time.Second*defaultTimeout)
	}
	return c.kafkaProducer.Close(ctx)
}

func eventMsg(prefix string, err error) string {
	return fmt.Sprintf(eventMessageFMT, prefix, err.Error())
}
