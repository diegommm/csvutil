package csvutil

import (
	"reflect"
	"strings"
)

type tag struct {
	name      string
	prefix    string
	empty     bool
	omitEmpty bool
	ignore    bool
	inline    bool
}

func parseTags(tagname string, field reflect.StructField) []tag {
	var t tag
	tags := strings.Split(field.Tag.Get(tagname), ",")
	if len(tags) == 1 && tags[0] == "" {
		t.name = field.Name
		t.empty = true
		return []tag{t}
	}

	switch tags[0] {
	case "-":
		t.ignore = true
		return []tag{t}
	case "":
		t.name = field.Name
	default:
		t.name = tags[0]
	}

	var multiTags []string
	for _, tagOpt := range tags[1:] {
		switch tagOpt {
		default:
			if tagOpt[:6] == "multi=" {
				multiTags = strings.Split(tagOpt[6:], " ")
			}
		case "omitempty":
			t.omitEmpty = true
		case "inline":
			if walkType(field.Type).Kind() == reflect.Struct {
				t.inline = true
				t.prefix = tags[0]
			}
		}
	}

	if len(multiTags) == 0 {
		return []tag{t}
	}

	t.inline = false
	ret := make([]tag, len(multiTags))
	for i, _ := range multiTags {
		t.name = multiTags[i]
		ret[i] = t
	}
	return ret
}
