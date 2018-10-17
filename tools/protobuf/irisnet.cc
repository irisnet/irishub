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
/*
constexpr char hexmap[] = {'0', '1', '2', '3', '4', '5', '6', '7',
                           '8', '9', 'a', 'b', 'c', 'd', 'e', 'f'};

string hexStr(const char *data, int len)
{
    string s(len * 2, ' ');
    for (int i = 0; i < len; ++i) {
        s[2 * i]     = hexmap[(data[i] & 0xF0) >> 4];
        s[2 * i + 1] = hexmap[data[i] & 0x0F];
    }
    return s;
}
*/

char* serializeJsonToProto(char* protoFilePath, char* protoFileName, char* messageName, char* jsonString, int* length)
{
    DiskSourceTree sourceTree;
    //look up .proto file in current directory
    string filePath(protoFilePath);
    sourceTree.MapPath("", filePath);
    Importer importer(&sourceTree, NULL);
    //runtime compile foo.proto
    string fileName(protoFileName);
    importer.Import(fileName);

    string protoMessageName(messageName);
    const Descriptor *descriptor = importer.pool()->FindMessageTypeByName(protoMessageName);

    // build a dynamic message by "Pair" proto
    DynamicMessageFactory factory;
    const Message *message = factory.GetPrototype(descriptor);
    // create a real instance of "Pair"
    Message *pair = message->New();

    string json(jsonString);
    if (Json2Pb(*pair, json) < 0 ) {
        std::cout << "ERROR" << std::endl;
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
    importer.Import(fileName);

    string protoMessageName(messageName);
    const Descriptor *descriptor = importer.pool()->FindMessageTypeByName(protoMessageName);

    // build a dynamic message by "Pair" proto
    DynamicMessageFactory factory;
    const Message *message = factory.GetPrototype(descriptor);
    // create a real instance of "Pair"
    Message *pair = message->New();

    stringstream ss;
    ss << ((char *)protoByteArray);
    pair->ParseFromIstream(&ss);

    string resultString;
    Pb2Json(*pair, resultString);
    *length = resultString.size();
    delete pair;

    char* resultByte = new char[resultString.size()]();
    for (int i=0; i< resultString.size(); i++){
        resultByte[i]=resultString[i];
    }
    return resultByte;
}
