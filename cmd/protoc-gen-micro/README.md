# protoc-gen-micro

proto file generator for micro-lite

Test protoc-gen-micro only:

```
protoc --plugin ./protoc-gen-micro --micro_out=./ example/*.proto
```

Test full protoc:

(install protoc-gen-micro first)
```
go install github.com/ofavor/micro-let/cmd/protoc-gen-micro
```

```
protoc --go_out=./ --micro_out=./ example/proto/*.proto
```

## References

* [protoc-gen-goexample](https://github.com/drekle/protoc-gen-goexample)