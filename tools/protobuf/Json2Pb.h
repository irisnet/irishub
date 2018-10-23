#pragma once

#ifndef __JSON2PB_H__
#define __JSON2PB_H__

#include <string>

namespace google {
namespace protobuf {
class Message;
}
}

int32_t Json2Pb(google::protobuf::Message& stMsg, const std::string& sJson);
int32_t Json2Pb(google::protobuf::Message& stMsg, const char* sJson);
int32_t Pb2Json(const google::protobuf::Message& stMsg, std::string& sJson);

#endif//__JSON2PB_H__
