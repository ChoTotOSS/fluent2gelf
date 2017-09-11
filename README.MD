###  Log forwarder for fluentd/fluentbit to graylog

input -> [fluentd/fluentbit] -> fluent2gelf (multiline concat) -> graylog

### Build and run

```
$ make build
$ ./dist/fluent2gelf -c sample.yml # run with sample config
```

### Sample config
```
- match: kube.log.*_default_shop-profile*.log
  host: 10.60.6.48
  port: 12203
  multiline:
    regexp: ^201

- match: kube.log.*_default_inbox*.log
  host: 10.60.6.48
  port: 12204
  multiline:
    begin: 201

- match: kube.log.*_default_auth*.log
  host: 10.60.6.48
  port: 12205
  multiline:
    begin: 201
    index_lt: 5
```
