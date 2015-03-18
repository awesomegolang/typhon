// this file tests error passing between services

package test

import (
	"testing"

	"github.com/b2aio/typhon/client"
	"github.com/b2aio/typhon/errors"
	"github.com/b2aio/typhon/example/proto/callhello"
	"github.com/b2aio/typhon/server"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestErrorPropagation(t *testing.T) {

	s := InitServer(t, "test")
	defer s.Close()

	var (
		errorMessage = "failure"
		errorCode    = 1234
	)

	// Register test endpoints
	s.RegisterEndpoint(&server.Endpoint{
		Name: "callerror",
		Handler: func(req server.Request) (proto.Message, error) {
			// simulate some failure
			return nil, errors.New(errorCode, errorMessage, map[string]string{
				"public key": "public value",
			}, map[string]string{
				"private key": "private value",
			})
		},

		// for convienience use example request & response
		Request:  &callhello.Request{},
		Response: &callhello.Response{},
	})

	// call the service
	resp := &callhello.Response{}
	err := client.Request(
		nil,                                // context
		"test",                             // service
		"callerror",                        // service endpoint to call
		&callhello.Request{Value: "Bunny"}, // request
		resp, // response
	)

	// Assert our error matches
	require.NotNil(t, err)
	typhonErr := err.(*errors.Error)
	assert.Equal(t, errorCode, typhonErr.Code)
	assert.Equal(t, errorMessage, typhonErr.Error())
	assert.Equal(t, errorMessage, typhonErr.Message)
	assert.Equal(t, map[string]string{
		"public key": "public value",
	}, typhonErr.PublicContext)
	assert.Equal(t, map[string]string{
		"private key": "private value",
		"service":     "test",
		"endpoint":    "callerror",
	}, typhonErr.PrivateContext)
}
