package reporter

import (
	"errors"
	"fmt"

	kafka "github.com/ONSdigital/dp-kafka/v2"
	"github.com/ONSdigital/dp-reporter-client/model"
	"github.com/ONSdigital/dp-reporter-client/schema"
	"github.com/ONSdigital/log.go/log"
)

const (
	idEmpty          = "cannot Notify, ID is a required field but was empty"
	contextEmpty     = "cannot Notify, errContext is a required field but was empty"
	sendingEvent     = "sending reportEvent for application error"
	failedToMarshal  = "failed to marshal reportEvent to avro"
	eventMessageFMT  = "%s: %s"
	serviceNameEmpty = "cannot create new import error reporter as serviceName is empty"
	kafkaProducerNil = "cannot create new import error reporter as kafkaProducer is nil"
	eventTypeErr     = "error"
	reportEventKey   = "reportEvent"
)

// KafkaProducer interface of the dp-kafka kafka.Producer
type KafkaProducer interface {
	Channels() *kafka.ProducerChannels
}

type marshalFunc func(s interface{}) ([]byte, error)

// ErrorReporter is the interface that wraps the Notify method.
// ID and errContent are required parameters. If neither is provided or there is any error while attempting to
// report the event then an error is returned which the caller can handle as they see fit.
type ErrorReporter interface {
	Notify(id string, errContext string, err error) error
}

// ImportErrorReporter a reporter for sending error reports to the import-reporter
type ImportErrorReporter struct {
	kafkaProducer KafkaProducer
	marshal       marshalFunc
	serviceName   string
}

// NewImportErrorReporter create a new ImportErrorReporter to send error reports to the import-reporter
func NewImportErrorReporter(kafkaProducer KafkaProducer, serviceName string) (ImportErrorReporter, error) {
	if kafkaProducer == nil {
		return ImportErrorReporter{}, errors.New(kafkaProducerNil)
	}
	if len(serviceName) == 0 {
		return ImportErrorReporter{}, errors.New(serviceNameEmpty)
	}
	return ImportErrorReporter{
		serviceName:   serviceName,
		kafkaProducer: kafkaProducer,
		marshal:       schema.ReportEventSchema.Marshal,
	}, nil
}

// Notify send an error report to the import-reporter
// ID and errContent are required parameters. If neither is provided or there is any error while attempting to
// report the event then an error is returned which the caller can handle as they see fit.
func (c ImportErrorReporter) Notify(id string, errContext string, err error) error {
	if len(id) == 0 {
		log.Event(nil, idEmpty)
		return errors.New(idEmpty)
	}
	if len(errContext) == 0 {
		log.Event(nil, contextEmpty)
		return errors.New(contextEmpty)
	}

	reportEvent := &model.ReportEvent{
		InstanceID:  id,
		EventMsg:    fmt.Sprintf(eventMessageFMT, errContext, err.Error()),
		ServiceName: c.serviceName,
		EventType:   eventTypeErr,
	}

	reportEventData := log.Data{reportEventKey: reportEvent}
	log.Event(nil, sendingEvent, reportEventData)

	avroBytes, err := c.marshal(reportEvent)
	if err != nil {
		log.Event(nil, failedToMarshal, log.Error(err), reportEventData)
		return err
	}
	fmt.Printf("Sending bytes to output channel")

	c.kafkaProducer.Channels().Output <- avroBytes
	return nil
}
