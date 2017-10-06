dp-reporter-client
================

A client for sending report events to the dp-import-reporter

### Getting started
##### Create a new client
To create a new client you need to provide a KafkaProducer [see go-ns kafka.Producer](https://github.com/ONSdigital/go-ns/blob/master/kafka/producer.go) 
and service name. The producer should already be configured to talk to the desired instance of the reporter and the 
service name should be the name of your service - as this is where the error event has occurred. __NOTE:__ It is the 
responsibility of the caller to gracefully close the KafkaProducer and reactor to any error it might return.
```go
	var p KafkaProducer
	...

	cli, err := client.NewReporterClient(p, "MyService")
	if err != nil {
 
		//...handle error
 
	}
```
##### Report an error event
To report an error event you __MUST__ provide an __instanceID__ and an __error context__ the reporter requires 
both parameters in order to update the instance. If neither are provided the client will not attempt to report the 
event and will return an error. Additionally you can provide the error causing the event the and any additional 
logging data.
```
cli.ReportError("instance-123", "unexpected http response status", err, log.Data{"key": "value")
```

### Configuration

An overview of the configuration options available, either as a table of
environment variables, or with a link to a configuration guide.

| Environment variable | Default | Description
| -------------------- | ------- | -----------
| BIND_ADDR            | :8080   | The host and port to bind to

### Contributing

See [CONTRIBUTING](CONTRIBUTING.md) for details.

### License

Copyright © 2016-2017, Office for National Statistics (https://www.ons.gov.uk)

Released under MIT license, see [LICENSE](LICENSE.md) for details.
