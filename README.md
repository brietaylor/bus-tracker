## Building Protobuf Bindings

GTFS Protobuf file, downloaded from:  https://raw.githubusercontent.com/google/transit/refs/heads/master/gtfs-realtime/proto/gtfs-realtime.proto


Compile using protoc:

```sh
protoc --go_out=./ proto/gtfs-realtime.proto --go_opt=Mproto/gtfs-realtime.proto=./proto
```

This should generate a file in the `proto/` folder called `gtfs-realtime.pb.go`.

## Updating static GTFS

### From Translink

Go to this page to find the latest version: https://www.translink.ca/about-us/doing-business-with-translink/app-developer-resources/gtfs/gtfs-data

Extract in the `gtfs-static/translink-DATE` directory.