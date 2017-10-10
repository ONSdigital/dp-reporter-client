package reporter

import (
	"errors"
	"testing"

	"fmt"

	"github.com/ONSdigital/dp-reporter-client/model"
	"github.com/ONSdigital/dp-reporter-client/reporter/reportertest"
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
	Convey("Given ImportErrorReporter has been configured correctly", t, func() {
		marshalParams := make([]interface{}, 0)
		output, kafkaProducer, marshalFunc := setup(&marshalParams, schema.ReportEventSchema.Marshal)
		target := ImportErrorReporter{kafkaProducer: kafkaProducer, serviceName: "Bob", marshal: marshalFunc}

		Convey("When Notify is called with valid parameters", func() {
			err := target.Notify(testInstanceID, errContext, cause)

			avroBytes := <-output

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And Marshall is called once with the expected parameters", func() {
				So(len(marshalParams), ShouldEqual, 1)
				So(marshalParams[0], ShouldResemble, expectedReportEvent)
			})

			Convey("And kafkaProducer.Output is called once with the expected parameters", func() {
				var actual model.ReportEvent
				schema.ReportEventSchema.Unmarshal(avroBytes, &actual)

				So(kafkaProducer.OutputCalls(), ShouldEqual, 1)
				So(&actual, ShouldResemble, expectedReportEvent)
			})

		})
	})
}

func TestHandler_HandleMarshalError(t *testing.T) {
	Convey("Given ImportErrorReporter has been configured correctly", t, func() {

		marshalParams := make([]interface{}, 0)
		_, kafkaProducer, marshalFunc := setup(&marshalParams, func(s interface{}) ([]byte, error) {
			return nil, errors.New("Bork!")
		})
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

			Convey("And kafkaProducer.Output is never called", func() {
				So(kafkaProducer.OutputCalls(), ShouldEqual, 0)
			})
		})
	})
}

func TestHandler_HandleInstanceIDEmpty(t *testing.T) {
	Convey("Given ImportErrorReporter has been configured correctly", t, func() {
		marshalParams := make([]interface{}, 0)
		_, kafkaProducer, marshalFunc := setup(&marshalParams, schema.ReportEventSchema.Marshal)
		target := &ImportErrorReporter{kafkaProducer: kafkaProducer, serviceName: "Bob", marshal: marshalFunc}

		Convey("When Notify is called with an empty instanceID", func() {
			err := target.Notify("", errContext, cause)

			Convey("Then the expected error is returned", func() {
				So(err, ShouldResemble, errors.New(idEmpty))
			})

			Convey("And marshal is never called", func() {
				So(len(marshalParams), ShouldEqual, 0)
			})

			Convey("And kafkaProducer.Output is never called", func() {
				So(kafkaProducer.OutputCalls(), ShouldEqual, 0)
			})
		})
	})
}

func TestHandler_HandleErrContextEmpty(t *testing.T) {
	Convey("Given ImportErrorReporter has been configured correctly", t, func() {
		marshalParams := make([]interface{}, 0)
		_, kafkaProducer, marshalFunc := setup(&marshalParams, schema.ReportEventSchema.Marshal)
		target := &ImportErrorReporter{kafkaProducer: kafkaProducer, serviceName: "Bob", marshal: marshalFunc}

		Convey("When Notify is called with an empty errContext", func() {
			err := target.Notify(testInstanceID, "", cause)

			Convey("Then the expected error is returned", func() {
				So(err, ShouldResemble, errors.New(contextEmpty))
			})

			Convey("And marshal is never called", func() {
				So(len(marshalParams), ShouldEqual, 0)
			})

			Convey("And kafkaProducer.Output is never called", func() {
				So(kafkaProducer.OutputCalls(), ShouldEqual, 0)
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
		kafkaProducer := &reportertest.KafkaProducerMock{}

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
		kafkaProducer := &reportertest.KafkaProducerMock{}

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

func setup(marshalParams *[]interface{}, marshal func(s interface{}) ([]byte, error)) (chan []byte, *reportertest.KafkaProducerMock, marshalFunc) {
	output := make(chan []byte, 1)
	producer := reportertest.NewKafkaProducerMock(output)

	marshalFunc := func(s interface{}) ([]byte, error) {
		*marshalParams = append(*marshalParams, s)
		return marshal(s)
	}
	return output, producer, marshalFunc
}
