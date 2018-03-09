package iservice

import (
	"fmt"
	cmn "github.com/tendermint/tmlibs/common"
	"github.com/pkg/errors"
)

const (
	Tag_heigth = "height"
	Tag_Service_name = "service-name"
	Tag_chain_id = "chain_id"
	Tag_messaging_type = "messaging_type"
	Tag_key = "key"
	Tag_sender = "sender"
)

func IndexHeight(heigth int64) cmn.KVPair {
	return cmn.KVPair{[]byte(Tag_heigth), []byte(fmt.Sprintf("%d", heigth))}
}

func IndexServiceName(name string) cmn.KVPair {
	return cmn.KVPair{[]byte(Tag_Service_name), []byte(fmt.Sprintf("%s", name))}
}

func IndexChainId(chainId string) cmn.KVPair {
	return cmn.KVPair{[]byte(Tag_chain_id), []byte(fmt.Sprintf("%s", chainId))}
}

func IndexMessagingType(msgType string) cmn.KVPair {
	return cmn.KVPair{[]byte(Tag_messaging_type), []byte(fmt.Sprintf("%s", msgType))}
}

func IndexKVTag(key string,value string) cmn.KVPair{
	return cmn.KVPair{[]byte(key), []byte(value)}
}

func IndexKeyTag(key string) cmn.KVPair{
	return cmn.KVPair{[]byte(Tag_key), []byte(key)}
}

func IndexSender(address []byte) cmn.KVPair{
	return cmn.KVPair{[]byte(Tag_sender), address}
}

func Qheight(heigth int64) string{
	return FmtQuery(Tag_heigth,heigth)
}

func QServiceName(name string) string{
	return FmtQuery(Tag_Service_name,name)
}

func QChainId(chainId string) string{
	return FmtQuery(Tag_chain_id,chainId)
}

func QMessagingType(msgtype string) string{
	return FmtQuery(Tag_messaging_type,msgtype)
}

func QKVTag(key string,value string) string{
	return FmtQuery(key,value)
}

func QKeyTag(key string) string{
	return FmtQuery(Tag_key,key)
}

func QSender(sender string) string{
	return FmtQuery(Tag_sender,sender)
}

func FmtQuery(key string,value interface{}) string{
	switch value.(type) {
		case string:
			return fmt.Sprintf("%s='%s'", key,value)
		case int64:{
			return fmt.Sprintf("%s=%d", key,value)
		}
		case []byte:{
			return fmt.Sprintf("%s='%s'", key,value)
		}

	}
	panic(errors.New("unsupport tags type"))
}
