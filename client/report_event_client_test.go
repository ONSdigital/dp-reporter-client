package client

import (
	"testing"
	. "github.com/smartystreets/goconvey/convey"
	"errors"
	"github.com/ONSdigital/dp-reporter-client/model"
	"github.com/ONSdigital/dp-reporter-client/schema"
	"github.com/ONSdigital/dp-reporter-client/client/clienttest"
	"github.com/ONSdigital/go-ns/log"
)

var (
	testInstanceID = "666"
	cause          = errors.New("Flubba Wubba Dub Dub")
	errContext     = "Ricky Ticky Tic Tac"

	expectedReportEvent = &model.ReportEvent{
		InstanceID:  testInstanceID,
		EventType:   "error",
		ServiceName: "Bob",
		EventMsg:    eventMsg(errContext, cause),
	}
)

func TestHandler_Handle(t *testing.T) {
	Convey("Given ReporterClient has been configured correctly", t, func() {
		marshalParams := make([]interface{}, 0)
		output, kafkaProducer, marshalFunc := setup(&marshalParams, schema.ReportEventSchema.Marshal)
		target := &ReporterClient{kafkaProducer: kafkaProducer, serviceName: "Bob", marshal: marshalFunc}

		Convey("When ReportError is called with valid parameters", func() {
			err := target.ReportError(testInstanceID, errContext, cause, log.Data{"KEY": "value"})

			avroBytes := <-output

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
			})

			Convey("And Marshall is called once with the expectedReportEvent parameters", func() {
				So(len(marshalParams), ShouldEqual, 1)
				So(marshalParams[0], ShouldResemble, expectedReportEvent)
			})

			Convey("And kafkaProducer.Output is called once with the expectedReportEvent parameters", func() {
				var actual model.ReportEvent
				schema.ReportEventSchema.Unmarshal(avroBytes, &actual)

				So(kafkaProducer.OutputCalls(), ShouldEqual, 1)
				So(&actual, ShouldResemble, expectedReportEvent)
			})

		})
	})
}

func TestHandler_HandleMarshalError(t *testing.T) {
	Convey("Given ReporterClient has been configured correctly", t, func() {

		marshalParams := make([]interface{}, 0)
		_, kafkaProducer, marshalFunc := setup(&marshalParams, func(s interface{}) ([]byte, error) {
			return nil, errors.New("Bork!")
		})
		target := &ReporterClient{kafkaProducer: kafkaProducer, serviceName: "Bob", marshal: marshalFunc}

		Convey("When marshal returns an error", func() {
			err := target.ReportError(testInstanceID, errContext, cause, nil)

			Convey("Then the error returned to the caller", func() {
				So(err, ShouldResemble, errors.New("Bork!"))
			})

			Convey("And Marshall is called once with the expectedReportEvent parameters", func() {
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
	Convey("Given ReporterClient has been configured correctly", t, func() {
		marshalParams := make([]interface{}, 0)
		_, kafkaProducer, marshalFunc := setup(&marshalParams, schema.ReportEventSchema.Marshal)
		target := &ReporterClient{kafkaProducer: kafkaProducer, serviceName: "Bob", marshal: marshalFunc}

		Convey("When ReportError is called with an empty instanceID", func() {
			err := target.ReportError("", errContext, cause, nil)

			Convey("Then no error is returned", func() {
				So(err, ShouldBeNil)
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
			cli, err := NewReporterClient(nil, serviceName)

			Convey("Then cli is nil and the expected error is returned", func() {
				So(cli, ShouldBeNil)
				So(err, ShouldResemble, errors.New(kafkaProducerNil))
			})
		})
	})

	Convey("Given serviceName is empty", t, func() {
		serviceName := ""
		kafkaProducer := &clienttest.KafkaProducerMock{}

		Convey("When NewReportClient is called", func() {
			cli, err := NewReporterClient(kafkaProducer, serviceName)

			Convey("Then cli is nil and the expected error is returned", func() {
				So(cli, ShouldBeNil)
				So(err, ShouldResemble, errors.New(serviceNameEmpty))
			})
		})
	})

	Convey("Given a valid kafkaProducer and serviceName", t, func() {
		serviceName := "awesome-service"
		kafkaProducer := &clienttest.KafkaProducerMock{}

		Convey("When NewReportClient is called", func() {
			cli, err := NewReporterClient(kafkaProducer, serviceName)

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

func setup(marshalParams *[]interface{}, marshal func(s interface{}) ([]byte, error)) (chan []byte, *clienttest.KafkaProducerMock, marshalFunc) {
	output := make(chan []byte, 1)
	producer := clienttest.NewKafkaProducerMock(output, nil)

	marshalFunc := func(s interface{}) ([]byte, error) {
		*marshalParams = append(*marshalParams, s)
		return marshal(s)
	}
	return output, producer, marshalFunc
}
