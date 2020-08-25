package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"sync"
	pb "tensorflow_serving/apis"
	"golang.org/x/net/context"
	//tf_core_framework "tensorflow/core/framework"
	tf_core_framework "tensorflow/core/framework"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
	//google_protobuf "github.com/golang/protobuf/ptypes/wrappers"
)

var (
	serverAddr         = flag.String("server_addr", "127.0.0.1:9000", "The server address in the format of host:port")
	modelName          = flag.String("model_name", "cancer", "TensorFlow model name")
	imagePath       = flag.String("image_file", "./img.png", "Input image file")
	modelVersion       = flag.Int64("model_version", 1, "TensorFlow model version")
	tls                = flag.Bool("tls", false, "Connection uses TLS if true, else plain TCP")
	caFile             = flag.String("ca_file", "testdata/ca.pem", "The file containning the CA root cert file")
	serverHostOverride = flag.String("server_host_override", "x.test.youtube.com", "The server name use to verify the hostname returned by TLS handshake")
)

type PredictionClient struct {
	mu      sync.RWMutex
	rpcConn *grpc.ClientConn
	psvcConn pb.PredictionServiceClient
}

func NewClient()(*PredictionClient, error){
	var opts []grpc.DialOption
	if *tls {
		var sn string
		if *serverHostOverride != "" {
			sn = *serverHostOverride
		}
		var creds credentials.TransportCredentials
		if *caFile != "" {
			var err error
			creds, err = credentials.NewClientTLSFromFile(*caFile, sn)
			if err != nil {
				grpclog.Fatalf("Failed to create TLS credentials %v", err)
				return nil, err
			}
		} else {
			creds = credentials.NewClientTLSFromCert(nil, sn)
		}
		opts = append(opts, grpc.WithTransportCredentials(creds))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}
	conn, err := grpc.Dial(*serverAddr, opts...)
	if err != nil {
		grpclog.Fatalf("fail to dial: %v", err)
		return nil, err
	}
	return &PredictionClient{rpcConn: conn, psvcConn: pb.NewPredictionServiceClient(conn)}, nil
}

func constructImageRequest(imageBytes []byte  ) (*pb.PredictRequest){

	request := &pb.PredictRequest{
		ModelSpec: &pb.ModelSpec{
			Name:          "inception",
			SignatureName: "predict_images",
		},
		Inputs: map[string]*tf_core_framework.TensorProto{
			"images": &tf_core_framework.TensorProto{
				Dtype: tf_core_framework.DataType_DT_STRING,
				TensorShape: &tf_core_framework.TensorShapeProto{
					Dim: []*tf_core_framework.TensorShapeProto_Dim{
						&tf_core_framework.TensorShapeProto_Dim{
							Size: int64(1),
						},
					},
				},
				StringVal: [][]byte{imageBytes},
			},
		},
	}
	return request

}

func (c *PredictionClient) PredictImage(modelName string, request *pb.PredictRequest) (*pb.PredictResponse, error) {
	resp, err := c.psvcConn.Predict(context.Background(), request)
	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return resp, nil
}

func (c *PredictionClient) Close() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.psvcConn = nil
	return c.rpcConn.Close()
}

func main() {
	flag.Parse()

	imgPath, err := filepath.Abs(*imagePath)
	if err != nil {
		log.Fatalln(err)
	}

	//imageBytes, err := ioutil.ReadFile(imgPath)
	imageBytes, err := ioutil.ReadFile(imgPath)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := NewClient()
	if err !=nil {
		log.Fatalln(err)
	}
	//request := constructImageRequest(imageBytes)
	fmt.Println(client)
	fmt.Println(imageBytes)
	//fmt.Println(request)

}
