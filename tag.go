package di

import (
	"reflect"
	"strings"
)

type Tag struct {
	tag        string
	name       string
	allowEmpty bool
}

func (tag *Tag) IsSkip() bool {
	return len(tag.tag) == 0 || tag.name == "-"
}

func (tag *Tag) AllowEmpty() bool {
	return tag.allowEmpty
}

func (tag *Tag) GetName() string {
	return tag.name
}

func (tag *Tag) ParseTag(field reflect.StructField, tagStr string) *Tag {
	result := tag
	result.tag = tagStr
	if len(tagStr) == 0 {
		return result
	}

	parts := strings.Split(tagStr, ",")
	result.name = parts[0]
	if len(result.name) == 0 {
		result.name = field.Name
	}

	for i := 1; i < len(parts); i++ {
		switch parts[i] {
		case "allowEmpty":
			result.allowEmpty = true
		}
	}

	return result
}
