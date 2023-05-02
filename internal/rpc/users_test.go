package rpc

import (
	"context"
	"testing"

	pb "github.com/dom/user/api/dom/user/v1"
	"github.com/stretchr/testify/assert"
)

func Test_SayHello(t *testing.T) {

	service := &userSvc{}
	helloRequest := &pb.HelloRequest{Name: "Peter"}

	reply, err := service.SayHello(context.Background(), helloRequest)

	assert.Nil(t, err)
	assert.Equal(t, "Hello "+helloRequest.Name, reply.Message)
	assert.True(t, true, "should be true")

}
