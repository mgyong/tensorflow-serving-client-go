module main

go 1.14

replace tensorflow_serving => ./tensorflow_serving/

//replace tensorflow => ./tensorflow

require (
	google.golang.org/grpc v1.30.0
	//	tensorflow v0.0.0
	tensorflow_serving v0.0.0 // indirect
)
