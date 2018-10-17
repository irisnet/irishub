#pragma once

#include <stdio.h>
#include <stdlib.h>

#ifdef __cplusplus
extern "C" {
#endif
char* serializeJsonToProto(char* protoFilePath, char* protoFileName, char* messageName, char* jsonString, int* length);
char* convertProtoToJson(char* protoFilePath, char* protoFileName, char* messageName, void* protoByteArray, int byteArrayLength, int* length);
#ifdef __cplusplus
}
#endif
