dp-reporter-client
================

A client for sending report events to the dp-import-reporter

### Getting started
##### Create a new client
```
	p := &KafkaProducer{
	    //...Implementation of interface
	}
	cli, err := client.NewReporterClient(p, "MyService")
	if err != nil {
		//...handle error
	}
```
##### Report an error event
```
cli.ReportError("instance-123", "unexpected http response status", err, log.Data{"key": "value")
```

##### Close the client as part of your application shutdown
```
	ctx, _ := context.WithTimeout(context.Background(), time.Second * 5)
	cli.Close(ctx)
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
