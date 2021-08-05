package middleware

import "github.com/bcowtech/structproto"

var _ structproto.TagResolver = NoneTagResolver

func NoneTagResolver(fieldname, token string) (*structproto.Tag, error) {
	var tag = &structproto.Tag{
		Name: fieldname,
	}
	return tag, nil
}
