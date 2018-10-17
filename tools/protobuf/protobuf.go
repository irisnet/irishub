package protobuf

/*
#cgo LDFLAGS: -L . -liservice -lstdc++ -lprotobuf -ljsoncpp
#cgo CFLAGS: -I ./
#include "irisnet.h"
*/
import "C"

import (
    "unsafe"
)

func init() {
}

func SerializeJsonToProto(protoPath, protoName, messageName string, jsonString string) []byte {
	path := C.CString(protoPath)
	name := C.CString(protoName)
	message := C.CString(messageName)
	json := C.CString(jsonString)

	defer func() {
		C.free(unsafe.Pointer(path))
		C.free(unsafe.Pointer(name))
		C.free(unsafe.Pointer(message))
		C.free(unsafe.Pointer(json))
	}()

	resultLength := C.int(0)
	resultByte := C.serializeJsonToProto(path, name, message, json, &resultLength)
	protoBufByte := C.GoBytes(unsafe.Pointer(resultByte), resultLength)

	return protoBufByte
}

func ConvertProtoToJson(protoPath, protoName, messageName string, protoBufByte []byte) string {
	path := C.CString(protoPath)
	name := C.CString(protoName)
	message := C.CString(messageName)

	defer func() {
		C.free(unsafe.Pointer(path))
		C.free(unsafe.Pointer(name))
		C.free(unsafe.Pointer(message))
	}()

	protoBufCByte := C.CBytes(protoBufByte)
	protoBufCByteLength := C.int(len(protoBufByte))
	resultLength := C.int(0)
	resultByte := C.convertProtoToJson(path, name, message, protoBufCByte, protoBufCByteLength, &resultLength)
	json := C.GoStringN(resultByte, resultLength)

	return json
}
