package validator

import (
	ut "github.com/go-playground/universal-translator"
	"reflect"
	"sync"
	"github.com/go-playground/validator/v10"
	"errors"
	"github.com/go-playground/locales/zh"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

// trans use single pattern
var trans = GetZhTranslator()

// GetZhTranslator get simple chinese translator
func GetZhTranslator() ut.Translator {
	//中文翻译器
	zh_ch := zh.New()
	uni := ut.New(zh_ch)
	trans, _ := uni.GetTranslator("zh")
	return trans
}

// Validator
type Validator struct {
	once     sync.Once
	validate *validator.Validate
}

// getKindOf return the kind of data
func getKindOf(data interface{}) reflect.Kind {
	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

// ValidateStruct validate struct
func (v *Validator) ValidateStruct(obj interface{}) error {

	if getKindOf(obj) == reflect.Struct {
		v.lazyInit()

		if err := v.validate.Struct(obj); err != nil {
			//验证器
			for _, err := range err.(validator.ValidationErrors) {
				return errors.New(err.Translate(trans))
			}

		}
	}

	return nil
}

// Engine
func (v *Validator) Engine() interface{} {
	v.lazyInit()
	return v.validate
}

// lazyInit
func (v *Validator) lazyInit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("binding")

		// define filed name
		v.validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
			return fld.Tag.Get("comment")
		})

		// use zh-CN
		_ = zh_translations.RegisterDefaultTranslations(v.validate, trans)
	})
}
