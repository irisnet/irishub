package protoidl

import (
	"github.com/emicklei/proto"
	"strings"
	"fmt"
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

// validate proto idl text and get all method
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
		attributes := make(map[string]string)
		if r.Comment != nil {
			attributes, err = transferComment(r.Comment.Lines)
			if err != nil {
				return methods, err
			}
		}
		method := Method{
			r.Name,
			attributes,
		}
		methods = append(methods, method)
	}
	return methods, nil
}

func transferComment(lines []string) (map[string]string, error) {
	attributes := make(map[string]string)
	for _, line := range lines {
		index := strings.Index(line, "@Attribute")
		if index == -1 {
			continue
		}
		ss := strings.SplitN(line[index+10:], ":", 2)
		if len(ss) < 2 {
			return attributes, fmt.Errorf("invalid attribute at %s", line)
		}
		key := strings.Replace(ss[0], " ", "", -1)
		if key == "" {
			return attributes, fmt.Errorf("attribute has empty key at %s", line)
		}
		value := strings.TrimSpace(ss[1])
		if key == "" {
			return attributes, fmt.Errorf("attribute has empty value at %s", line)
		}
		attributes[key] = value
	}
	return attributes, nil
}
