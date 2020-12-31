# protoc-gen-go-micro

proto file generator for micro-lite

Test protoc-gen-go-micro only:

```
protoc --plugin ./protoc-gen-go-micro --go-micro_out=./ example/*.proto
```

Test full protoc:

(install protoc-gen-go-micro first)
```
go install github.com/ofavor/micro-let/cmd/protoc-gen-go-micro
```

```
protoc --go_out=./ --micro_out=./ example/proto/*.proto
```

## References

* [protoc-gen-goexample](https://github.com/drekle/protoc-gen-goexample)