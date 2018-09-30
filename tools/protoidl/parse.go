package protoidl

import (
	"github.com/emicklei/proto"
	"strings"
)

// validate proto idl text
func ValidateProto(content string) (bool, error) {
	reader := strings.NewReader(content)
	parser := proto.NewParser(reader)
	_, err := parser.Parse()
	if err != nil {
		return false, err
	}
	return true, nil
}

// get all method from proto idl text
func GetMethods(content string) (methods []Method, err error) {
	reader := strings.NewReader(content)
	parser := proto.NewParser(reader)
	definition, err := parser.Parse()
	if err != nil {
		return methods, err
	}

	// iterate definition get all method
	var rs []*proto.RPC
	proto.Walk(definition,
		proto.WithRPC(func(r *proto.RPC) {
			rs = append(rs, r)
		}))

	// get method attribute from comment, each line comment only define one attribute
	for _, r := range rs {
		commentMap := make(map[string]string)
		if r.Comment != nil {
			commentMap = transferComment(r.Comment.Lines)
		}
		method := Method{
			r.Name,
			commentMap,
		}
		methods = append(methods, method)
	}
	return methods, nil
}

func transferComment(lines []string) map[string]string {
	commentMap := make(map[string]string)
	for _, line := range lines {
		line = strings.Replace(line, " ", "", -1)
		ss := strings.Split(line, ":")
		if len(ss) < 2 {
			continue
		}
		commentMap[ss[0]] = ss[1]
	}
	return commentMap
}
