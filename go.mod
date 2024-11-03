module github.com/brietaylor/online-bus-tracker

go 1.22.8

replace github.com/google/transit/proto => ./proto

require (
	github.com/gocarina/gocsv v0.0.0-20240520201108-78e41c74b4b1
	google.golang.org/protobuf v1.35.1
)
