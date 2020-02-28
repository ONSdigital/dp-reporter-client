package reporter

import (
	"errors"
	"testing"

	"fmt"

	kafka "github.com/ONSdigital/dp-kafka"
	"github.com/ONSdigital/dp-kafka/kafkatest"
	"github.com/ONSdigital/dp-reporter-client/model"
	"github.com/ONSdigital/dp-reporter-client/schema"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	testInstanceID = "666"
	cause          = errors.New("flubba wubba dub dub")
	errContext     = "Ricky Ticky Tic Tac"

	expectedReportEvent = &model.ReportEvent{
		InstanceID:  testInstanceID,
		EventType:   "error",
		ServiceName: "Bob",
		EventMsg:    fmt.Sprintf(eventMessageFMT, errContext, cause.Error()),
	}
)

func TestHandler_Handle(t *testing.T) {
	Convey("Given ImportErrorReporter has been configured correctly", t, func(c C) {
		marshalParams := make([]interface{}, 0)
		kafkaProducer, marshalFunc, pChannels := setup(&marshalParams, schema.ReportEventSchema.Marshal, true)
		target := ImportErrorReporter{kafkaProducer: kafkaProducer, serviceName: "Bob", marshal: marshalFunc}

		Convey("When Notify is called with valid parameters, marshall is called as expected and data is sent to output channel exactly once", func(c C) {

			// Call notify in a separate go routine
			go func() {
				err := target.Notify(testInstanceID, errContext, cause)
				c.Convey("Then no error is returned", func() {
					So(err, ShouldBeNil)
				})
			}()

			// Wait for output chan and close it, so that test fails if data is sent to channel more than once.
			avroBytes := <-pChannels.Output
			close(pChannels.Output)

			// Validate marshall
			So(len(marshalParams), ShouldEqual, 1)
			So(marshalParams[0], ShouldResemble, expectedReportEvent)

			// Validate kafka received data
			var actual model.ReportEvent
			schema.ReportEventSchema.Unmarshal(avroBytes, &actual)
			So(&actual, ShouldResemble, expectedReportEvent)
		})
	})
}

func TestHandler_HandleMarshalError(t *testing.T) {
	Convey("Given ImportErrorReporter has been configured correctly", t, func() {

		marshalParams := make([]interface{}, 0)
		kafkaProducer, marshalFunc, _ := setup(&marshalParams, func(s interface{}) ([]byte, error) {
			return nil, errors.New("Bork!")
		}, false)
		target := &ImportErrorReporter{kafkaProducer: kafkaProducer, serviceName: "Bob", marshal: marshalFunc}

		Convey("When marshal returns an error", func() {
			err := target.Notify(testInstanceID, errContext, cause)

			Convey("Then the error returned to the caller", func() {
				So(err, ShouldResemble, errors.New("Bork!"))
			})

			Convey("And Marshall is called once with the expected parameters", func() {
				So(len(marshalParams), ShouldEqual, 1)
				So(marshalParams[0], ShouldResemble, expectedReportEvent)
			})
		})
	})
}

func TestHandler_HandleInstanceIDEmpty(t *testing.T) {
	Convey("Given ImportErrorReporter has been configured correctly", t, func() {
		marshalParams := make([]interface{}, 0)
		kafkaProducer, marshalFunc, _ := setup(&marshalParams, schema.ReportEventSchema.Marshal, false)
		target := &ImportErrorReporter{kafkaProducer: kafkaProducer, serviceName: "Bob", marshal: marshalFunc}

		Convey("When Notify is called with an empty instanceID", func() {
			err := target.Notify("", errContext, cause)

			Convey("Then the expected error is returned", func() {
				So(err, ShouldResemble, errors.New(idEmpty))
			})

			Convey("And marshal is never called", func() {
				So(len(marshalParams), ShouldEqual, 0)
			})
		})
	})
}

func TestHandler_HandleErrContextEmpty(t *testing.T) {
	Convey("Given ImportErrorReporter has been configured correctly", t, func() {
		marshalParams := make([]interface{}, 0)
		kafkaProducer, marshalFunc, _ := setup(&marshalParams, schema.ReportEventSchema.Marshal, false)
		target := &ImportErrorReporter{kafkaProducer: kafkaProducer, serviceName: "Bob", marshal: marshalFunc}

		Convey("When Notify is called with an empty errContext", func() {
			err := target.Notify(testInstanceID, "", cause)

			Convey("Then the expected error is returned", func() {
				So(err, ShouldResemble, errors.New(contextEmpty))
			})

			Convey("And marshal is never called", func() {
				So(len(marshalParams), ShouldEqual, 0)
			})
		})
	})
}

func TestNewReporterClient(t *testing.T) {
	Convey("Given kafkaProducer is nil", t, func() {
		serviceName := "awesome-service"

		Convey("When NewReportClient is called", func() {
			cli, err := NewImportErrorReporter(nil, serviceName)

			Convey("Then cli is nil and the expected error is returned", func() {
				So(cli, ShouldResemble, ImportErrorReporter{})
				So(err, ShouldResemble, errors.New(kafkaProducerNil))
			})
		})
	})

	Convey("Given serviceName is empty", t, func() {
		serviceName := ""
		kafkaProducer := &kafkatest.MessageProducer{}

		Convey("When NewReportClient is called", func() {
			cli, err := NewImportErrorReporter(kafkaProducer, serviceName)

			Convey("Then cli is nil and the expected error is returned", func() {
				So(cli, ShouldResemble, ImportErrorReporter{})
				So(err, ShouldResemble, errors.New(serviceNameEmpty))
			})
		})
	})

	Convey("Given a valid kafkaProducer and serviceName", t, func() {
		serviceName := "awesome-service"
		kafkaProducer := &kafkatest.MessageProducer{}

		Convey("When NewReportClient is called", func() {
			cli, err := NewImportErrorReporter(kafkaProducer, serviceName)

			Convey("Then cli is configured as expected", func() {
				So(cli.serviceName, ShouldEqual, serviceName)
				So(cli.kafkaProducer, ShouldEqual, kafkaProducer)
				So(cli.marshal, ShouldEqual, schema.ReportEventSchema.Marshal)
				So(err, ShouldBeNil)
			})

			Convey("And no error is returned", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}

// setup creates a testing Kakfa producer and marshal func. If your test doesn't expect to use any kafka channel,
// please pass createChannels=false, so that the test fails if it does.
func setup(marshalParams *[]interface{}, marshal func(s interface{}) ([]byte, error), createChannels bool) (*kafkatest.MessageProducer, marshalFunc, kafka.ProducerChannels) {
	pChannels := kafka.ProducerChannels{}
	if createChannels {
		pChannels = kafka.CreateProducerChannels()
	}
	producer := kafkatest.NewMessageProducerWithChannels(pChannels)

	marshalFunc := func(s interface{}) ([]byte, error) {
		*marshalParams = append(*marshalParams, s)
		return marshal(s)
	}
	return producer, marshalFunc, pChannels
}
