# AWS SES Exporter

Exporter for AWS SES send statistics. `awsses_exporter` will poll all AWS SES regions and report back the latest Complaints, Timestamp, DeliveryAttempts, Bounces, Rejects. The AWS SES API endpoint `GetSendStatisticsInput` is called once per region per call to the `/metrics` endpoint

## Usage
```
$ awsses_exporter --help
usage: collector [<flags>]

Flags:
  -h, --help              Show context-sensitive help (also try --help-long and --help-man).
      --web.listen-address=":9199"
                          Address to listen on for web interface and telemetry.
      --web.telemetry-path="/metrics"
                          Path under which to expose metrics.
      --log.level="info"  Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]
      --log.format="logger:stderr"
                          Set the log target and format. Example: "logger:syslog?appname=bob&local=7" or "logger:stdout?json=true"
      --version           Show application version.
```

AWS credentials are pulled via the rules in https://docs.aws.amazon.com/sdk-for-go/api/aws/#Config. Most notably you can use the environment variables `AWS_ACCESS_KEY` and `AWS_SECRET_ACCESS_KEY`.

## Build

This easiest way to build and run this utility is with the included docker image: `docker build .`

You can also install it locally with:
```
go get github.com/FlasHs1/awsses_exporter.git
go install github.com/FlasHs1/awsses_exporter.git
```


## Example Metrics
Below is an example of the output of the `/metrics` endpoint
```
# HELP awsses_exporter_bounces Bounces per region
# TYPE awsses_exporter_bounces gauge
awsses_exporter_bounces{aws_region="us-east-1"} 0.0
awsses_exporter_bounces{aws_region="us-west-2"} 0.0
# HELP awsses_exporter_complaints Complaints per region
# TYPE awsses_exporter_complaints gauge
awsses_exporter_complaints{aws_region="us-east-1"} 0.0
awsses_exporter_complaints{aws_region="us-west-2"} 0.0
# HELP awsses_exporter_deliveryAttempts Delivery attempts per region
# TYPE awsses_exporter_deliveryAttempts gauge
awsses_exporter_deliveryAttempts{aws_region="us-east-1"} 5.0
awsses_exporter_deliveryAttempts{aws_region="us-west-2"} 2.0
# HELP awsses_exporter_rejects Rejects per region
# TYPE awsses_exporter_rejects gauge
awsses_exporter_rejects{aws_region="us-east-1"} 0.0
awsses_exporter_rejects{aws_region="us-west-2"} 0.0
```
