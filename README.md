dp-reporter-client
================

A client for sending report events to the dp-import-reporter

### Getting started
##### Reporting an import error
Report an import pipeline error via the `ImportErrorReporter`. To create a new ImportErrorReporter you need to provide a
 KafkaProducer [see go-ns kafka.Producer](https://github.com/ONSdigital/go-ns/blob/master/kafka/producer.go) and service
  name. The producer should already be configured to talk to the desired instance of the reporter and the service name 
  should be the name of your service - as this is where the error event has occurred. __NOTE:__ It is the responsibility
   of the caller to gracefully close the KafkaProducer and handle any error it returns.
```go
	var p KafkaProducer
	errReporter, err := reporter.NewImportErrorReporter(p, "MyService")
	if err != nil {
		//...handle error
	}
``` 
##### Report an error event
To report an error event you __MUST__ provide an __instanceID__ and an __error context__ the reporter requires 
both parameters in order to update the instance. If neither are provided the client will not attempt to report the 
event and will return an error. Additionally you can provide the error causing the event the and any additional 
logging data.
```go
    if err := errReporter.Notfiy("instance-123", "unexpected http response status", err); err != nil {
        //...handle error
    }
```
#### Testing
The `reporter.ErrorReporter` interface can be used in place of the `ImportErrorReporter` struct in cases where it is a
member of another struct or function parameter. Accepting this interface instead of the concrete struct allows the 
reporter to be easily tested. For example given:
```go
    type Foo struct {
        errReporter reporter.ErrorReporter
    }
    
    func (f Foo) Bar(id string, context string, cause error) error {
        return f.errReporter.Notify(id, context, cause)
    }
``` 

Testing the above:

```go
	Convey("Test Foo.Bar", t, func() {
		// The expected parameters to the Notify func
		expectedParams := reportertest.NotfiyParams{
			ID:         "666",
			ErrContext: "666",
			Err:        errors.New("bar"),
		}
 
		// The error to return when Notify is called.
		expectedErr := errors.New("expected error")
  
		// Create the mock & configure to return the expected error
		reporterMock := reportertest.NewImportErrorReporterMock(expectedErr)

		foo := Foo{
			errReporter: reporterMock,
		}
		
		// run the test
		err := foo.Bar("666", "666", errors.New("bar"))
 
		// Check the return val, number of invocations and the parameters supplied to Notify
		So(err, ShouldResemble, expectedErr)
		So(len(reporterMock.NotifyCalls()), ShouldEqual, 1)
		So(reporterMock.NotifyCalls()[0], ShouldResemble, expectedParams)
	})
```
The `reporter.ErrorReporter` package also provides a mock of the `KafkaProducer`
```
    output := make(chan []byte, 1)
    p := reportertest.NewKafkaProducerMock(output)
    ...
    So(len(p.OutputCalls()), ShouldEqual, 1)
```
### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2016-2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
