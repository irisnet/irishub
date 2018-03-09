package iservice

import (
	"fmt"
	cmn "github.com/tendermint/tmlibs/common"
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
	return fmt.Sprintf("%s=%d", Tag_heigth,heigth)
}

func QServiceName(name string) string{
	return fmt.Sprintf("%s='%s'", Tag_Service_name,name)
}

func QChainId(chainId string) string{
	return fmt.Sprintf("%s='%s'", Tag_chain_id,chainId)
}

func QMessagingType(msgtype string) string{
	return fmt.Sprintf("%s='%s'", Tag_messaging_type,msgtype)
}

func QKVTag(key string,value string) string{
	return fmt.Sprintf("%s='%s'", key,value)
}

func QKeyTag(key string) string{
	return fmt.Sprintf("%s='%s'", Tag_key,key)
}

func QSender(sender string) string{
	return fmt.Sprintf("%s='%s'", Tag_sender,sender)
}
