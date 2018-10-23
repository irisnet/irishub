#include "Json2Pb.h"
#include "irisnet.h"
#include <stdio.h>
#include <iostream>
#include <sstream>
#include <fstream>
#include <google/protobuf/descriptor.h>
#include <google/protobuf/descriptor.pb.h>
#include <google/protobuf/dynamic_message.h>
#include <google/protobuf/compiler/importer.h>

using namespace std;
using namespace google::protobuf;
using namespace google::protobuf::compiler;
using google::protobuf::Message;

char* serializeJsonToProto(char* protoFilePath, char* protoFileName, char* messageName, char* jsonString, int* length)
{
    DiskSourceTree sourceTree;
    //look up .proto file in current directory
    string filePath(protoFilePath);
    sourceTree.MapPath("", filePath);
    Importer importer(&sourceTree, NULL);
    //runtime compile foo.proto
    string fileName(protoFileName);
    if (importer.Import(fileName) == NULL) {
        std::cout << "ERROR: failed to import proto file, please verifiy your proto file path and its validity" << std::endl;
        return NULL;
    }


    string protoMessageName(messageName);
    const Descriptor *descriptor = importer.pool()->FindMessageTypeByName(protoMessageName);
    if (descriptor == NULL) {
        std::cout << "ERROR: failed to get message descriptor, please check your message name" << std::endl;
        return NULL;
    }

    // build a dynamic message by "Pair" proto
    DynamicMessageFactory factory;
    const Message *message = factory.GetPrototype(descriptor);
    // create a real instance of "Pair"
    Message *pair = message->New();

    string json(jsonString);
    if (Json2Pb(*pair, json) < 0 ) {
        std::cout << "ERROR: failed to covert json to protobuf message" << std::endl;
        return NULL;
    }

    ostringstream oss;
    pair->SerializeToOstream(&oss);
    string resultString = oss.str();
    *length =resultString.size();

    delete pair;
    char* resultByte = new char[resultString.size()]();
    for (int i=0; i< resultString.size(); i++){
        resultByte[i]=resultString[i];
    }
    return resultByte;
}

char* convertProtoToJson(char* protoFilePath, char* protoFileName, char* messageName, void* protoByteArray, int byteArrayLength, int* length)
{
    DiskSourceTree sourceTree;
    //look up .proto file in current directory
    string filePath(protoFilePath);
    sourceTree.MapPath("", filePath);
    Importer importer(&sourceTree, NULL);
    //runtime compile foo.proto
    string fileName(protoFileName);
    if (importer.Import(fileName) == NULL) {
            std::cout << "ERROR: failed to import proto file, please verifiy your proto file path and its validity" << std::endl;
            return NULL;
        }

    string protoMessageName(messageName);
    const Descriptor *descriptor = importer.pool()->FindMessageTypeByName(protoMessageName);
    if (descriptor == NULL) {
            std::cout << "ERROR: failed to get message descriptor, please check your message name" << std::endl;
            return NULL;
        }

    // build a dynamic message by "Pair" proto
    DynamicMessageFactory factory;
    const Message *message = factory.GetPrototype(descriptor);
    // create a real instance of "Pair"
    Message *pair = message->New();

    stringstream ss;
    ss << ((char *)protoByteArray);
    pair->ParseFromIstream(&ss);

    string resultString;
    if (Pb2Json(*pair, resultString)< 0 ) {
        std::cout << "ERROR: failed to covert protobuf message to json" << std::endl;
    }
    *length = resultString.size();
    delete pair;

    char* resultByte = new char[resultString.size()]();
    for (int i=0; i< resultString.size(); i++){
        resultByte[i]=resultString[i];
    }
    return resultByte;
}
