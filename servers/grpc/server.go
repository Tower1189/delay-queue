package grpc

import (
	"golang.org/x/net/context"
	"strings"
	"github.com/ouqiang/delay-queue/delayqueue"
	"net"
	"log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"errors"
)

type server struct{}

func (s *server) Pop(ctx context.Context, in *PopRequest) (*DataResponse, error) {
	topic := strings.TrimSpace(in.Topic)
	if topic == "" {
		return nil, errors.New("topic 不能为空")
	}
	topics := strings.Split(topic, ",")
	job, err := delayqueue.Pop(topics)
	if err != nil {
		return  nil, err
	}
	if job == nil {
		return nil, nil
	}
	return &DataResponse{
		Data: &Data{Id: job.Id, Body: job.Body},
	}, nil
}



//
//func generateSuccessBody(msg string, data *Data) *DataResponse {
//	return generateResponseBody(0, msg, data)
//}
//
//func generateFailureBody(msg string) *DataResponse {
//	return generateResponseBody(1, msg, nil)
//}
//
//func generateResponseBody(code int32, msg string, data *Data) *DataResponse {
//	body := &DataResponse{}
//	body.Code = code
//	body.Message = msg
//	body.Data = data
//	return body
//}

func Run(port string) {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	RegisterDelayQueueServer(s, &server{})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}