module github.com/ONSdigital/dp-reporter-client

go 1.23.1

//to avoid  [CVE-2022-29153] CWE-918: Server-Side Request Forgery (SSRF)
exclude github.com/hashicorp/consul/api v1.1.0

// To fix: [CVE-2024-24786] CWE-835: Loop with Unreachable Exit Condition ('Infinite Loop')
replace google.golang.org/protobuf => google.golang.org/protobuf v1.35.2

require (
	github.com/ONSdigital/dp-kafka/v2 v2.8.0
	github.com/ONSdigital/go-ns v0.0.0-20191104121206-f144c4ec2e58
	github.com/ONSdigital/log.go/v2 v2.4.1
	github.com/smartystreets/goconvey v1.8.0
)

require (
	github.com/ONSdigital/dp-api-clients-go/v2 v2.252.0 // indirect
	github.com/ONSdigital/dp-healthcheck v1.6.1 // indirect
	github.com/ONSdigital/dp-net/v2 v2.9.1 // indirect
	github.com/Shopify/sarama v1.30.1 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/eapache/go-resiliency v1.2.0 // indirect
	github.com/eapache/go-xerial-snappy v0.0.0-20180814174437-776d5712da21 // indirect
	github.com/eapache/queue v1.1.0 // indirect
	github.com/fatih/color v1.15.0 // indirect
	github.com/go-avro/avro v0.0.0-20171219232920-444163702c11 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/gopherjs/gopherjs v1.17.2 // indirect
	github.com/hashicorp/go-uuid v1.0.3 // indirect
	github.com/hokaccha/go-prettyjson v0.0.0-20211117102719-0474bc63780f // indirect
	github.com/jcmturner/aescts/v2 v2.0.0 // indirect
	github.com/jcmturner/dnsutils/v2 v2.0.0 // indirect
	github.com/jcmturner/gofork v1.7.6 // indirect
	github.com/jcmturner/gokrb5/v8 v8.4.4 // indirect
	github.com/jcmturner/rpc/v2 v2.0.3 // indirect
	github.com/jtolds/gls v4.20.0+incompatible // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.18 // indirect
	github.com/pierrec/lz4 v2.6.1+incompatible // indirect
	github.com/rcrowley/go-metrics v0.0.0-20201227073835-cf1acfcdf475 // indirect
	github.com/smartystreets/assertions v1.13.1 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/net v0.21.0 // indirect
	golang.org/x/sys v0.27.0 // indirect
)
