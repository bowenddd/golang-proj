package session

import (
	"reflect"
)

const (
	BeforeQuery  = "BeforeQuery"
	AfterQuery   = "AfterQuery"
	BeforeUpdate = "BeforeUpdate"
	AfterUpdate  = "AfterUpdate"
	BeforeDelete = "BeforeDelete"
	AfterDelete  = "AfterDelete"
	BeforeInsert = "BeforeInsert"
	AfterInsert  = "AfterInsert"
)

type IBeforeQuery interface {
	BeforeQuery(s *Session) error
}

type IAfterQuery interface {
	AfterQuery(s *Session) error
}

type IBeforeUpdate interface {
	BeforeUpdate(s *Session) error
}

type IAfterUpdate interface {
	AfterUpdate(s *Session) error
}

type IBeforeDelete interface {
	BeforeDelete(s *Session) error
}

type IAfterDelete interface {
	AfterDelete(s *Session) error
}

type IBeforeInsert interface {
	BeforeInsert(s *Session) error
}

type IAfterInsert interface {
	AfterInsert(s *Session) error
}

func (s *Session) CallMethod(name string, value interface{}) {
	param := reflect.ValueOf(s.RefTable().Model).Interface()
	if value != nil {
		param = reflect.ValueOf(value).Interface()
	}
	switch name {
	case BeforeQuery:
		if call, ok := param.(IBeforeQuery); ok {
			_ = call.BeforeQuery(s)
		}
	case AfterQuery:
		if call, ok := param.(IAfterQuery); ok {
			_ = call.AfterQuery(s)
		}
	case BeforeUpdate:
		if call, ok := param.(IBeforeUpdate); ok {
			_ = call.BeforeUpdate(s)
		}
	case AfterUpdate:
		if call, ok := param.(IAfterUpdate); ok {
			_ = call.AfterUpdate(s)
		}
	case BeforeDelete:
		if call, ok := param.(IBeforeDelete); ok {
			_ = call.BeforeDelete(s)
		}
	case AfterDelete:
		if call, ok := param.(IAfterDelete); ok {
			_ = call.AfterDelete(s)
		}
	case BeforeInsert:
		if call, ok := param.(IBeforeInsert); ok {
			_ = call.BeforeInsert(s)
		}
	case AfterInsert:
		if call, ok := param.(IAfterInsert); ok {
			_ = call.AfterInsert(s)
		}
	default:
		panic("unsupported hook method")
	}
	return
}
