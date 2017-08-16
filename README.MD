###  Log forwarder for fluentd/fluentbit to graylog

input -> [fluentd/fluentbit] -> fluent2gelf (multiline concat) -> graylog

### Build and run

```
$ make build
$ ./dist/fluent2gelf -c sample.yml # run with sample config
```
