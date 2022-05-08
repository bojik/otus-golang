package logger

import "fmt"

type Parameter interface {
	GetKeyValue() string
}

type StringParam struct {
	k string
	v string
}

func (s StringParam) GetKeyValue() string {
	return fmt.Sprintf("%s=%s", s.k, s.v)
}

var _ Parameter = (*StringParam)(nil)

func NewStringParam(key, value string) StringParam {
	return StringParam{k: key, v: value}
}

type IntParam struct {
	k string
	v int
}

func (i IntParam) GetKeyValue() string {
	return fmt.Sprintf("%s=%d", i.k, i.v)
}

var _ Parameter = (*IntParam)(nil)

func NewIntParam(key string, val int) IntParam {
	return IntParam{k: key, v: val}
}

type Int32Param struct {
	k string
	v int32
}

func (i Int32Param) GetKeyValue() string {
	return fmt.Sprintf("%s=%d", i.k, i.v)
}

var _ Parameter = (*Int32Param)(nil)

func NewInt32Param(key string, val int32) Int32Param {
	return Int32Param{k: key, v: val}
}

type Int64Param struct {
	k string
	v int64
}

func (i Int64Param) GetKeyValue() string {
	return fmt.Sprintf("%s=%d", i.k, i.v)
}

var _ Parameter = (*Int64Param)(nil)

func NewInt64Param(key string, val int64) Int64Param {
	return Int64Param{k: key, v: val}
}

type Float32Param struct {
	k string
	v float32
}

func (i Float32Param) GetKeyValue() string {
	return fmt.Sprintf("%s=%f", i.k, i.v)
}

var _ Parameter = (*Float32Param)(nil)

func NewFloat32Param(key string, val float32) Float32Param {
	return Float32Param{k: key, v: val}
}

type Float64Param struct {
	k string
	v float64
}

func (i Float64Param) GetKeyValue() string {
	return fmt.Sprintf("%s=%f", i.k, i.v)
}

var _ Parameter = (*Float64Param)(nil)

func NewFloat64Param(key string, val float64) Float64Param {
	return Float64Param{k: key, v: val}
}
