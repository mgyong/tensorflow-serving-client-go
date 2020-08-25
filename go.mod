module github.com/mgyong/merlion/cmd/modelserving/

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	google.golang.org/grpc v1.30.0
	tensorflow v0.0.0-00010101000000-000000000000 // indirect
	tensorflow_serving v0.0.0-00010101000000-000000000000
)

replace tensorflow => ./tensorflow

replace tensorflow_serving => ./tensorflow_serving
