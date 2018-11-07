package protobuf

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

var json = `{"age":20,"gender":"female","key":"1234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890","value":10}`
var byteArray = []byte{10, 100, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 49, 50, 51, 52, 53, 54, 55, 56, 57, 48, 16, 10, 24, 20, 34, 6, 102, 101, 109, 97, 108, 101}

func TestSerializeJsonToProto(t *testing.T) {
	path := "./test/"
	name := "foo.proto"
	messageName := "Test"
	jsonString := json

	resultByte, err := SerializeJsonToProto(path, name, messageName, jsonString)
	require.Nil(t, err)
	fmt.Println(resultByte)
	require.Equal(t, byteArray, resultByte)

	// Check error handler
	path = "./test1/"
	_, err = SerializeJsonToProto(path, name, messageName, jsonString)
	fmt.Println(err.Error())
	require.EqualErrorf(t, err, "proto file doesn't exist", "Error message doesn't match")

	path = "./test/"
	name = "foo_invalid.proto"
	_, err = SerializeJsonToProto(path, name, messageName, jsonString)
	require.EqualErrorf(t, err, "failed to serialize json to protobuf", "Error message doesn't match")

	path = "./test/"
	name = "foo.proto"
	messageName = "Test1"
	_, err = SerializeJsonToProto(path, name, messageName, jsonString)
	require.EqualErrorf(t, err, "failed to serialize json to protobuf", "Error message doesn't match")
}

func TestConvertProtoToJson(t *testing.T) {
	path := "./test/"
	name := "foo.proto"
	messageName := "Test"
	protobufByteArray := byteArray

	jsonString, err := ConvertProtoToJson(path, name, messageName, protobufByteArray)
	require.Nil(t, err)
	fmt.Println(json)
	jsonString = strings.TrimSuffix(jsonString, "\n")
	require.Equal(t, json, jsonString)

	// Check error handler
	path = "./test1/"
	_, err = ConvertProtoToJson(path, name, messageName, protobufByteArray)
	fmt.Println(err.Error())
	require.EqualErrorf(t, err, "proto file doesn't exist", "Error message doesn't match")

	path = "./test/"
	name = "foo_invalid.proto"
	_, err = ConvertProtoToJson(path, name, messageName, protobufByteArray)
	require.EqualErrorf(t, err, "failed to convert protobuf to json", "Error message doesn't match")

	path = "./test/"
	name = "foo.proto"
	messageName = "Test1"
	_, err = ConvertProtoToJson(path, name, messageName, protobufByteArray)
	require.EqualErrorf(t, err, "failed to convert protobuf to json", "Error message doesn't match")
}
