package utils

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
)

// FormatValidationError 将 validator 校验错误格式化为中文可读字符串
// 针对注册、登录、地址等常用字段提供友好提示
func FormatValidationError(err error) string {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return err.Error()
	}
	var msgs []string
	for _, fe := range validationErrors {
		msgs = append(msgs, fieldErrorMessage(fe))
	}
	return strings.Join(msgs, "; ")
}

// fieldErrorMessage 根据字段名和校验规则返回中文错误提示
func fieldErrorMessage(fe validator.FieldError) string {
	field := fe.Field()
	tag := fe.Tag()
	param := fe.Param()
	switch field {
	case "Phone":
		switch tag {
		case "required":
			return "手机号不能为空"
		case "len":
			return "手机号必须为" + param + "位"
		case "numeric":
			return "手机号必须为数字"
		}
	case "Password":
		switch tag {
		case "required":
			return "密码不能为空"
		case "min":
			return "密码长度不能小于" + param + "位"
		case "max":
			return "密码长度不能大于" + param + "位"
		}
	case "Name":
		switch tag {
		case "required":
			return "姓名不能为空"
		case "max":
			return "姓名不能超过" + param + "个字符"
		}
	case "Address":
		switch tag {
		case "required":
			return "地址不能为空"
		}
	case "Quantity":
		switch tag {
		case "required":
			return "数量不能为空"
		case "min":
			return "数量不能小于" + param
		}
	case "FoodID":
		return "商品不能为空"
	case "MerchantID":
		return "商家不能为空"
	case "AddressID":
		return "地址不能为空"
	case "DeliveryFee":
		return "配送费不能小于" + param
	case "Remark":
		return "备注不能超过" + param + "个字符"
	}
	return fe.Error()
}
