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

	fmt.Println("-----------------serialize json string to proto byte array-----------------------")
	resultByte := SerializeJsonToProto(path, name, messageName, jsonString)
	fmt.Println(resultByte)
	require.Equal(t, byteArray, resultByte)
}

func TestConvertProtoToJson(t *testing.T) {
	path := "./test/"
	name := "foo.proto"
	messageName := "Test"
	protobufByteArray := byteArray

	fmt.Println("-----------------serialize json string to proto byte array-----------------------")
	jsonString := ConvertProtoToJson(path, name, messageName, protobufByteArray)
	fmt.Println(json)
	jsonString = strings.TrimSuffix(jsonString, "\n")
	require.Equal(t, json, jsonString)
}
