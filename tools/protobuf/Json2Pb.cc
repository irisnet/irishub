#include <string>
#include <json/json.h>

#include <google/protobuf/message.h>
#include <google/protobuf/descriptor.h>

#include "Json2Pb.h"

namespace {
#include "Bin2Ascii.h"
}

using google::protobuf::Message;
using google::protobuf::MessageFactory;
using google::protobuf::Descriptor;
using google::protobuf::FieldDescriptor;
using google::protobuf::EnumDescriptor;
using google::protobuf::EnumValueDescriptor;
using google::protobuf::Reflection;

static int32_t _Pb2Json(const Message& stMsg, Json::Value& stRoot);
static int32_t _Field2Json(const Message& stMsg, const FieldDescriptor* field, size_t index, Json::Value& stRoot)
{
    const Reflection *ref = stMsg.GetReflection();
    const bool repeated = field->is_repeated();
    switch (field->cpp_type())
    {
#define _CONVERT(type, ctype, sfunc, afunc)		\
        case FieldDescriptor::type: {			\
            ctype value = (repeated)?		\
                ref->afunc(stMsg, field, index):	\
                ref->sfunc(stMsg, field);		\
            stRoot = Json::Value(value);			\
            break;					\
        }

        _CONVERT(CPPTYPE_DOUBLE, double, GetDouble, GetRepeatedDouble);
        _CONVERT(CPPTYPE_FLOAT, float, GetFloat, GetRepeatedFloat);
        _CONVERT(CPPTYPE_INT64, long long int, GetInt64, GetRepeatedInt64);
        _CONVERT(CPPTYPE_UINT64, unsigned long long int, GetUInt64, GetRepeatedUInt64);
        _CONVERT(CPPTYPE_INT32, int32_t, GetInt32, GetRepeatedInt32);
        _CONVERT(CPPTYPE_UINT32, uint32_t, GetUInt32, GetRepeatedUInt32);
        _CONVERT(CPPTYPE_BOOL, bool, GetBool, GetRepeatedBool);
#undef _CONVERT
        case FieldDescriptor::CPPTYPE_STRING: {
            std::string scratch;
            const std::string &value = (repeated)?
                ref->GetRepeatedStringReference(stMsg, field, index, &scratch):
                ref->GetStringReference(stMsg, field, &scratch);
            if (field->type() == FieldDescriptor::TYPE_BYTES)
                stRoot = Json::Value(B64Encode(value));
            else
                stRoot = Json::Value(value);
            break;
        }
        case FieldDescriptor::CPPTYPE_MESSAGE: {
            const Message& mf = (repeated)?
                ref->GetRepeatedMessage(stMsg, field, index):
                ref->GetMessage(stMsg, field);
            if (_Pb2Json(mf, stRoot) != 0) return -1;
            break;
        }
        case FieldDescriptor::CPPTYPE_ENUM: {
            const EnumValueDescriptor* ef = (repeated)?
                ref->GetRepeatedEnum(stMsg, field, index):
                ref->GetEnum(stMsg, field);

            stRoot = Json::Value(ef->number());
            break;
        }
        default:
            break;
    }
    if (stRoot.empty()) {
        printf("%s: Fail to convert to json", field->name().c_str());
        return -1;
    }
    return 0;
}

static int32_t _Pb2Json(const Message& stMsg, Json::Value& stRoot)
{
    const Descriptor *d = stMsg.GetDescriptor();
    const Reflection *ref = stMsg.GetReflection();
    if (!d || !ref) return -1;

    std::vector<const FieldDescriptor *> fields;
    ref->ListFields(stMsg, &fields);

    for (size_t i = 0; i != fields.size(); i++)
    {
        const FieldDescriptor *field = fields[i];

        Json::Value jf;
        if(field->is_repeated()) {
            int count = ref->FieldSize(stMsg, field);
            if (count == 0) continue;
            Json::Value tmp(Json::arrayValue);
            for (int j = 0; j < count; j++)
                if (_Field2Json(stMsg, field, j, tmp[j]) != 0) return -1;
            jf = tmp;
        } else if (ref->HasField(stMsg, field)) {
            if (_Field2Json(stMsg, field, 0, jf)) return -1;
        }
        else
            continue;

        const std::string &name = (field->is_extension())?field->full_name():field->name();
        stRoot[name] = jf;
    }
    return 0;
}

static int32_t _Json2Pb(Message& stMsg, const Json::Value& stRoot);
static int32_t _Json2Field(Message& stMsg, const FieldDescriptor *field, const Json::Value& stRoot)
{
    const Reflection *ref = stMsg.GetReflection();
    const bool repeated = field->is_repeated();
    
    switch (field->cpp_type())
    {
#define _SET_OR_ADD(sfunc, afunc, value)			\
        do {						\
            if (repeated)				\
                ref->afunc(&stMsg, field, value);	\
            else					\
                ref->sfunc(&stMsg, field, value);	\
        } while (0)

#define _CONVERT(type, ctype, func, sfunc, afunc) 		\
        case FieldDescriptor::type: {			\
            ctype value = stRoot.func();				\
            _SET_OR_ADD(sfunc, afunc, value);	\
            break;					\
        }

        _CONVERT(CPPTYPE_DOUBLE, double, asDouble, SetDouble, AddDouble);
        _CONVERT(CPPTYPE_FLOAT, float, asFloat, SetFloat, AddFloat);
        _CONVERT(CPPTYPE_INT64, int64_t, asInt64, SetInt64, AddInt64);
        _CONVERT(CPPTYPE_UINT64, uint64_t, asUInt64, SetUInt64, AddUInt64);
        _CONVERT(CPPTYPE_INT32, int32_t, asInt, SetInt32, AddInt32);
        _CONVERT(CPPTYPE_UINT32, uint32_t, asUInt, SetUInt32, AddUInt32);
        _CONVERT(CPPTYPE_BOOL, bool, asBool, SetBool, AddBool);

        case FieldDescriptor::CPPTYPE_STRING: {
            if (!stRoot.isString()) {
                printf("%s: Not a string", field->name().c_str());
                return -1;
            }
            const char * value = stRoot.asCString();
            
            if(field->type() == FieldDescriptor::TYPE_BYTES)
                _SET_OR_ADD(SetString, AddString, B64Decode(value));
            else
                _SET_OR_ADD(SetString, AddString, value);
            break;
        }
        case FieldDescriptor::CPPTYPE_MESSAGE: {
            Message *mf = (repeated)?
                ref->AddMessage(&stMsg, field):
                ref->MutableMessage(&stMsg, field);
            if (_Json2Pb(*mf, stRoot) != 0) return -1;
            break;
        }
        case FieldDescriptor::CPPTYPE_ENUM: {
            const EnumDescriptor *ed = field->enum_type();
            const EnumValueDescriptor *ev = 0;
            if (stRoot.isIntegral()) {
                ev = ed->FindValueByNumber(stRoot.asInt());
            } else if (stRoot.isString()) {
                ev = ed->FindValueByName(stRoot.asCString());
            } else {
                printf("%s: Not an integer or string", field->name().c_str());
                return -1;
            }
            if (!ev) {
                printf("%s: Enum value not found", field->name().c_str());
                return -1;
            }
            _SET_OR_ADD(SetEnum, AddEnum, ev);
            break;
        }
        default:
            break;
    }
    return 0;
}

static int32_t _Json2Pb(Message& stMsg, const Json::Value& stRoot)
{
    const Descriptor *d = stMsg.GetDescriptor();
    const Reflection *ref = stMsg.GetReflection();
    if (!d || !ref) {
        printf("No descriptor or reflection");
        return -1;
    }
    
    for (Json::Value::const_iterator  it = stRoot.begin(); it != stRoot.end(); ++it)
    {
        const FieldDescriptor *field = d->FindFieldByName(it.name().data());
        if (!field)
            field = ref->FindKnownExtensionByName(it.name().data());

        if (!field) {
            printf("%s: %s",it.name().data(), "Unknown field");
            return -1;
        }
        if (field->is_repeated()) {
            if (!it->isArray()) {
                printf("%s: Not array", field->name().c_str());
                return -1;
            }
            for (int j = 0; j < it->size(); j++) {
                if (_Json2Field(stMsg, field, (*it)[j]) != 0) return -1;
            }
        } else
            if (_Json2Field(stMsg, field, *it) != 0) return -1;
    }
    
    return 0;
}

int32_t Json2Pb(Message& stMsg, const std::string& sJson)
{
    Json::Value root;
    Json::Reader reader;
    if(!reader.parse(sJson, root))
    {
        printf("parse json data failed data:%s", sJson.c_str());
        return -1;
    }

    if (!root.isObject()) {
        printf("Malformed JSON: not an object");
        return -1;
    }
    
    return _Json2Pb(stMsg, root);
}


int32_t Json2Pb(Message& stMsg, const char* sJson)
{
    std::string json(sJson);
    return Json2Pb(stMsg, json);
}

int32_t Pb2Json(const Message& stMsg, std::string& sJson)
{
    Json::Value root;
    Json::FastWriter writer;
    if (_Pb2Json(stMsg, root) != 0) return -1;
    sJson = writer.write(root);
    return 0;
}
