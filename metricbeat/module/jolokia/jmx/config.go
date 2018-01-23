package jmx

import (
	"encoding/json"
	"sort"
	"strings"
)

type JMXMapping struct {
	MBean      string
	Attributes []Attribute
	Target     Target
}

type Attribute struct {
	Attr  string
	Field string
}

type Target struct {
	Url      string
	User     string
	Password string
}

// RequestBlock is used to build the request blocks of the following format:
//
// [
//    {
//       "type":"read",
//       "mbean":"java.lang:type=Runtime",
//       "attribute":[
//          "Uptime"
//       ]
//    },
//    {
//       "type":"read",
//       "mbean":"java.lang:type=GarbageCollector,name=ConcurrentMarkSweep",
//       "attribute":[
//          "CollectionTime",
//          "CollectionCount"
//       ]
//    }
// ]
type RequestBlock struct {
	Type      string      `json:"type"`
	MBean     string      `json:"mbean"`
	Attribute []string    `json:"attribute"`
	Target    TargetBlock `json:"target"`
}

type TargetBlock struct {
	Url      string `json:"url"`
	User     string `json:"user"`
	Password string `json:"password"`
}

func buildRequestBodyAndMapping(mappings []JMXMapping) ([]byte, map[string]string, error) {
	responseMapping := map[string]string{}
	var blocks []RequestBlock

	for _, mapping := range mappings {
		rb := RequestBlock{
			Type:  "read",
			MBean: mapping.MBean,
			Target: TargetBlock{
				Url:      mapping.Target.Url,
				User:     mapping.Target.User,
				Password: mapping.Target.Password,
			},
		}

		for _, attribute := range mapping.Attributes {
			rb.Attribute = append(rb.Attribute, attribute.Attr)
			// MBean format (java.lang:type=GarbageCollector,name=PS MarkSweep)
			// Doamin : java.lang
			// Properties : type=GarbageCollector
			//              name=PS MarkSweep
			sortedMBeanKey := make([]byte, 0, len(mapping.MBean))
			mBeanDomainProp := strings.Split(mapping.MBean, ":")
			sortedMBeanKey = append(sortedMBeanKey, mBeanDomainProp[0]...)
			sortedMBeanKey = append(sortedMBeanKey, ':')
			mBeanProps := strings.Split(mBeanDomainProp[1], ",")
			sort.Strings(mBeanProps)
			for _, v := range mBeanProps {
				sortedMBeanKey = append(sortedMBeanKey, v...)
				sortedMBeanKey = append(sortedMBeanKey, ',')
			}
			responseMapping[string(sortedMBeanKey)+"_"+attribute.Attr] = attribute.Field
		}
		blocks = append(blocks, rb)
	}

	content, err := json.Marshal(blocks)
	return content, responseMapping, err
}
