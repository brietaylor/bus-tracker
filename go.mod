module github.com/brietaylor/online-bus-tracker

go 1.22.8

replace github.com/google/transit/proto => ./proto

require (
	github.com/alecthomas/kingpin/v2 v2.4.0
	github.com/gocarina/gocsv v0.0.0-20240520201108-78e41c74b4b1
	google.golang.org/protobuf v1.35.1
)

require (
	github.com/alecthomas/units v0.0.0-20240927000941-0f3dac36c52b // indirect
	github.com/xhit/go-str2duration/v2 v2.1.0 // indirect
)
