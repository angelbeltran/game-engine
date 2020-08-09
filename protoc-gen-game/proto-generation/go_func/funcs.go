package go_func

import (
	"fmt"
	"math"
	"strings"
)

func BoolToBool_NOT(v bool) bool {
	return !v
}

func BoolToInt_CAST(v bool) int {
	if v {
		return 1
	}
	return 0
}

func BoolToFloat_CAST(v bool) float64 {
	if v {
		return 1
	}
	return 0
}

func BoolToString_CAST(v bool) string {
	if v {
		return "true"
	}
	return "false"
}

func IntToBool_CAST(v int) bool {
	return v != 0
}

func IntToInt_NEG(v int) int {
	return -v
}

func IntToInt_INC(v int) int {
	return v + 1
}

func IntToInt_DEC(v int) int {
	return v - 1
}

func IntToString_CAST(v int) string {
	return fmt.Sprint(v)
}

func StringToBool_CAST(v string) bool {
	return len(v) > 0
}

func StringToInt_LEN(v string) int {
	return len(v)
}

func StringToString_UPPER(v string) string {
	return strings.ToUpper(v)
}

func StringToString_LOWER(v string) string {
	return strings.ToLower(v)
}

func BoolAndBoolToBool_EQ(v, w bool) bool {
	return v == w
}

func BoolAndBoolToBool_NEQ(v, w bool) bool {
	return v != w
}

func BoolAndBoolToBool_GT(v, w bool) bool {
	return v && !w
}

func BoolAndBoolToBool_LT(v, w bool) bool {
	return !v && w
}

func BoolAndBoolToBool_GTE(v, w bool) bool {
	return v || !w
}

func BoolAndBoolToBool_LTE(v, w bool) bool {
	return !v || w
}

func BoolAndBoolToBool_AND(v, w bool) bool {
	return v && w
}

func BoolAndBoolToBool_OR(v, w bool) bool {
	return v || w
}

func IntAndIntToBool_EQ(v, w int) bool {
	return v == w
}

func IntAndIntToBool_NEQ(v, w int) bool {
	return v != w
}

func IntAndIntToBool_GT(v, w int) bool {
	return v > w
}

func IntAndIntToBool_LT(v, w int) bool {
	return v < w
}

func IntAndIntToBool_GTE(v, w int) bool {
	return v >= w
}

func IntAndIntToBool_LTE(v, w int) bool {
	return v <= w
}

func IntAndIntToInt_ADD(v, w int) int {
	return v + w
}

func IntAndIntToInt_SUB(v, w int) int {
	return v - w
}

func IntAndIntToInt_MULT(v, w int) int {
	return v * w
}

func IntAndIntToInt_DIV(v, w int) int {
	if v == 0 {
		return 0
	}
	if w == 0 {
		if w > 0 {
			return math.MaxInt64
		}
		return math.MinInt64
	}

	return v / w
}

func IntAndIntToInt_MOD(v, w int) int {
	if v == 0 || w == 0 {
		return 0
	}

	return v % w
}

func FloatAndFloatToBool_EQ(v, w float64) bool {
	return v == w
}

func FloatAndFloatToBool_NEQ(v, w float64) bool {
	return v != w
}

func FloatAndFloatToBool_GT(v, w float64) bool {
	return v > w
}

func FloatAndFloatToBool_LT(v, w float64) bool {
	return v < w
}

func FloatAndFloatToBool_GTE(v, w float64) bool {
	return v >= w
}

func FloatAndFloatToBool_LTE(v, w float64) bool {
	return v <= w
}

func FloatAndFloatToBool_ADD(v, w float64) float64 {
	return v + w
}

func FloatAndFloatToBool_SUB(v, w float64) float64 {
	return v - w
}

func FloatAndFloatToBool_MULT(v, w float64) float64 {
	return v * w
}

func FloatAndFloatToFloat_DIV(v, w float64) float64 {
	if v == 0 {
		return 0
	}
	if w == 0 {
		if w > 0 {
			return math.MaxFloat64
		}
		return math.SmallestNonzeroFloat64
	}

	return v / w
}

func StringAndStringToBool_EQ(v, w string) bool {
	return v == w
}

func StringAndStringToBool_NEQ(v, w string) bool {
	return v != w
}

func StringAndStringToBool_GT(v, w string) bool {
	return v > w
}

func StringAndStringToBool_LT(v, w string) bool {
	return v < w
}

func StringAndStringToBool_GTE(v, w string) bool {
	return v >= w
}

func StringAndStringToBool_LTE(v, w string) bool {
	return v <= w
}

func StringAndStringToBool_CONCAT(v, w string) string {
	return v + w
}

func BoolsToBool_EQ(v ...bool) bool {
	if len(v) == 0 {
		return false
	}

	for _, b := range v[1:] {
		if v[0] != b {
			return false
		}
	}

	return true
}

func BoolsToBool_NEQ(v ...bool) bool {
	switch len(v) {
	case 0:
		return false
	case 1:
		return true
	case 2:
		return v[0] != v[1]
	default:
		return false
	}
}

func BoolsToBool_AND(v ...bool) bool {
	if len(v) == 0 {
		return false
	}

	for _, b := range v {
		if !b {
			return false
		}
	}

	return true
}

func BoolsToBool_OR(v ...bool) bool {
	if len(v) == 0 {
		return false
	}

	for _, b := range v {
		if b {
			return true
		}
	}

	return false
}

func IntsToInt_ADD(v ...int) int {
	var res int

	for _, i := range v {
		res += i
	}

	return res
}

func IntsToInt_MULT(v ...int) int {
	var res int

	for _, i := range v {
		res *= i
	}

	return res
}

func StringsToString_CONCAT(v ...string) string {
	var res string

	for _, s := range v {
		res += s
	}

	return res
}
