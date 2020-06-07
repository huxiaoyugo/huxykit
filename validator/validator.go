package validator

import (
	"fmt"
	"github.com/huxiaoyugo/huxykit/validator/rule"
	"reflect"
	"strings"
)

const (
	// 需要校验的tag
	validatorTag = "validator"
	// 规则分割符
	splitter = ","
)

type Validator struct {
	rules []rule.Rule
	// 是否需要检查所有的错误， 默认为true
	CheckAll bool
}

func NewValidator(rules ...rule.Rule) *Validator {
	validator := &Validator{
		rules: make([]rule.Rule, 0),
		CheckAll: true,
	}
	// register default rules
	validator.RegisterRules(rule.DefaultRules...)
	// register yourself rules
	validator.RegisterRules(rules...)
	return validator
}

/*
* IsValid check the target is valid
*
* PARAMS: t
* target: object need be checked, must be struct or struct pointer
*
* RETURN:
* nil: target is valid
* not nil: target is invalid and containing the error info
 */
func (v *Validator) IsValid(target interface{}) error{
	if target == nil {
		return fmt.Errorf("valid target can not be nil")
	}

	// 检查对象类型必须为struct或者struct pointer
	tType := reflect.TypeOf(target)
	if tType.Kind() == reflect.Ptr {
		tType = tType.Elem()
	}
	if tType.Kind() != reflect.Struct {
		return fmt.Errorf("valid target model must be struct or struct pointer")
	}
	tValue := reflect.Indirect(reflect.ValueOf(target))

	var resError = &ValidError{}
	for i := 0; i < tValue.NumField(); i++ {
		field := tType.Field(i)
		if tagStr, ok := field.Tag.Lookup(validatorTag); ok {
			if err := v.check(tagStr, &field, tValue.Field(i)); err != nil {
				resError.errors = append(resError.errors, err)
				if !v.CheckAll {
					return resError
				}
			}
		}
	}
	if len(resError.errors) == 0 {
		return nil
	}
	return resError
}

// RegisterRules: register check rules
func (v *Validator) RegisterRules(rules ...rule.Rule) {
	v.rules = append(v.rules, rules...)
}

func (v *Validator) check(ruleStr string, structField *reflect.StructField, val reflect.Value) error {
	var resError = &ValidError{}
	// 每个字段可能存在多条规则
	ruleExps := strings.Split(ruleStr, splitter)
	// 遍历每一条规则，进行检查
	for _, ruleExp := range ruleExps {
		ruleExp = strings.Trim(ruleExp, " ")
		if len(ruleExp) == 0 {
			continue
		}
		for _, r := range v.rules {
			// 不可导出，直接忽略
			if !val.CanInterface() {
				continue
			}
			ok, err := r.Check(val.Interface(), ruleExp)
			if !ok {
				// 不匹配该条规则，直接跳过
				continue
			}
			// 匹配该规则，但是不符合规则的要求，直接返回错误
			if err != nil {
				err = fmt.Errorf("[field: %s] %s", structField.Name, err.Error())
				resError.errors = append(resError.errors, err)
				if !v.CheckAll {
					return resError
				}
			}
		}
	}
	if len(resError.errors) == 0 {
		return nil
	}
	return resError
}

type ValidError struct {
	errors []error
}

func(ve *ValidError) Error()string{
	var b = strings.Builder{}
	if len(ve.errors) == 0 {
		return ""
	}
	if len(ve.errors) == 1 {
		return ve.errors[0].Error()
	}
	n := (len(ve.errors)-1) * len("\n")
	for _, e := range ve.errors {
		n += len(e.Error())
	}
	b.Grow(n)
	b.WriteString(ve.errors[0].Error())
	for _, e := range ve.errors[1:] {
		b.WriteString("\n")
		b.WriteString(e.Error())
	}
	return b.String()
}