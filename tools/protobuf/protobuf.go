package protobuf

/*
#cgo LDFLAGS: -L . -liservice -lstdc++ -lprotobuf -ljsoncpp
#cgo CFLAGS: -I ./
#include "irisnet.h"
*/
import "C"

import (
	"fmt"
	"os"
	"unsafe"
)

func checkFileExistence(protoPath, protoName string) error {
	if _, err := os.Stat(protoPath + "/" + protoName); os.IsNotExist(err) {
		return fmt.Errorf("ERROR: proto file doesn't exist")
	}
	return nil
}

func SerializeJsonToProto(protoPath, protoName, messageName string, jsonString string) ([]byte, error) {
	if err := checkFileExistence(protoPath, protoName); err != nil {
		return nil, err
	}
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
	if resultLength == 0 {
		return nil, fmt.Errorf("failed to serialize json to protobuf")
	}
	protoBufByte := C.GoBytes(unsafe.Pointer(resultByte), resultLength)
	C.free(unsafe.Pointer(resultByte))
	return protoBufByte, nil
}

func ConvertProtoToJson(protoPath, protoName, messageName string, protoBufByte []byte) (string, error) {
	if err := checkFileExistence(protoPath, protoName); err != nil {
		return "", err
	}

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
	if resultLength == 0 {
		return "", fmt.Errorf("failed to convert protobuf to json")
	}
	json := C.GoStringN(resultByte, resultLength)
	C.free(unsafe.Pointer(resultByte))
	return json, nil
}
