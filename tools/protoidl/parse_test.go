package protoidl

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var (
	tests = []struct {
		content      string
		methodNumber int
		expectPass   bool
	}{{`// Copyright 2015 gRPC authors.
	//
	// Licensed under the Apache License, Version 2.0 (the "License");
	// you may not use this file except in compliance with the License.
	// You may obtain a copy of the License at
	//
	//     http://www.apache.org/licenses/LICENSE-2.0
	//
	// Unless required by applicable law or agreed to in writing, software
	// distributed under the License is distributed on an "AS IS" BASIS,
	// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	// See the License for the specific language governing permissions and
	// limitations under the License.

	syntax = "proto3";

	package helloworld;

	// The greeting service definition.
	service Greeter {
		//@Attribute description:sayHello
		//@Attribute output_privacy:NoPrivacy
		//@Attribute output_cached:NoCached
		rpc SayHello (HelloRequest) returns (HelloReply) {}
	}

	// The request message containing the user's name.
	message HelloRequest {
		string name = 1;
	}

	// The response message containing the greetings
	message HelloReply {
		string message = 1;
	}`, 1, true}, {`{}`, 0, false}}
)

func TestValidateProto(t *testing.T) {
	for _, tc := range tests {
		validate, _ := ValidateProto(tc.content)
		require.Equal(t, validate, tc.expectPass)
	}
}

func TestGetMethods(t *testing.T) {
	for _, tc := range tests {
		methods, err := GetMethods(tc.content)
		if !tc.expectPass {
			assert.Error(t, err)
			return
		}
		assert.NoError(t, err)
		require.Len(t, methods, tc.methodNumber)
		require.Equal(t, methods[0].Name, "SayHello")
		require.Equal(t, methods[0].Attributes, map[string]string{"description": "sayHello", "output_cached": "NoCached", "output_privacy": "NoPrivacy"})
	}
}
