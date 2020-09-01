package go_func

import (
	"fmt"
	"math"
	"strings"
)

// Chainable values

type BoolValue struct {
	value bool
	err   error
}

type IntValue struct {
	value int
	err   error
}

type FloatValue struct {
	value float64
	err   error
}

type StringValue struct {
	value string
	err   error
}

func Bool(v bool) BoolValue {
	return BoolValue{value: v}
}

func Int(v int) IntValue {
	return IntValue{value: v}
}

func Float(v float64) FloatValue {
	return FloatValue{value: v}
}

func String(v string) StringValue {
	return StringValue{value: v}
}

func (v BoolValue) Value() (bool, error) {
	return v.value, v.err
}

func (v IntValue) Value() (int, error) {
	return v.value, v.err
}

func (v FloatValue) Value() (float64, error) {
	return v.value, v.err
}

func (v StringValue) Value() (string, error) {
	return v.value, v.err
}

func (v BoolValue) Err() error {
	return v.err
}

func (v IntValue) Err() error {
	return v.err
}

func (v FloatValue) Err() error {
	return v.err
}

func (v StringValue) Err() error {
	return v.err
}

// Unary functions

// Bool -> .

func BoolToBool_NOT(v BoolValue) BoolValue {
	v.value = !v.value
	return v
}

func BoolToInt_CAST(v BoolValue) IntValue {
	i := IntValue{err: v.err}
	if v.value {
		i.value = 1
	}
	return i
}

func BoolToFloat_CAST(v BoolValue) FloatValue {
	f := FloatValue{err: v.err}
	if v.value {
		f.value = 1
	}
	return f
}

func BoolToString_CAST(v BoolValue) StringValue {
	s := StringValue{err: v.err}
	if v.value {
		s.value = "true"
	} else {
		s.value = "false"
	}
	return s
}

// Int -> .

func IntToBool_CAST(v IntValue) BoolValue {
	return BoolValue{value: v.value != 0, err: v.err}
}

func IntToInt_NEG(v IntValue) IntValue {
	v.value = -v.value
	return v
}

func IntToInt_INC(v IntValue) IntValue {
	v.value++
	return v
}

func IntToInt_DEC(v IntValue) IntValue {
	v.value--
	return v
}

func IntToString_CAST(v IntValue) StringValue {
	return StringValue{value: fmt.Sprint(v.value), err: v.err}
}

// String -> .

func StringToBool_CAST(v StringValue) BoolValue {
	return BoolValue{value: len(v.value) > 0, err: v.err}
}

func StringToInt_LEN(v StringValue) IntValue {
	return IntValue{value: len(v.value), err: v.err}
}

func StringToString_UPPER(v StringValue) StringValue {
	return StringValue{value: strings.ToUpper(v.value), err: v.err}
}

func StringToString_LOWER(v StringValue) StringValue {
	return StringValue{value: strings.ToLower(v.value), err: v.err}
}

// Binary functions

// (Bool, Bool) -> .

func BoolAndBoolToBool_EQ(v, w BoolValue) BoolValue {
	return BoolValue{value: v.value == w.value, err: firstError(v, w)}
}

func BoolAndBoolToBool_NEQ(v, w BoolValue) BoolValue {
	return BoolValue{value: v.value != w.value, err: firstError(v, w)}
}

func BoolAndBoolToBool_GT(v, w BoolValue) BoolValue {
	return BoolValue{value: v.value && !w.value, err: firstError(v, w)}
}

func BoolAndBoolToBool_LT(v, w BoolValue) BoolValue {
	return BoolValue{value: !v.value && w.value, err: firstError(v, w)}
}

func BoolAndBoolToBool_GTE(v, w BoolValue) BoolValue {
	return BoolValue{value: v.value || !w.value, err: firstError(v, w)}
}

func BoolAndBoolToBool_LTE(v, w BoolValue) BoolValue {
	return BoolValue{value: !v.value || w.value, err: firstError(v, w)}
}

func BoolAndBoolToBool_AND(v, w BoolValue) BoolValue {
	return BoolValue{value: v.value && w.value, err: firstError(v, w)}
}

func BoolAndBoolToBool_OR(v, w BoolValue) BoolValue {
	return BoolValue{value: v.value || w.value, err: firstError(v, w)}
}

// (Int, Int) -> .

func IntAndIntToBool_EQ(v, w IntValue) BoolValue {
	return BoolValue{value: v.value == w.value, err: firstError(v, w)}
}

func IntAndIntToBool_NEQ(v, w IntValue) BoolValue {
	return BoolValue{value: v.value != w.value, err: firstError(v, w)}
}

func IntAndIntToBool_GT(v, w IntValue) BoolValue {
	return BoolValue{value: v.value > w.value, err: firstError(v, w)}
}

func IntAndIntToBool_LT(v, w IntValue) BoolValue {
	return BoolValue{value: v.value < w.value, err: firstError(v, w)}
}

func IntAndIntToBool_GTE(v, w IntValue) BoolValue {
	return BoolValue{value: v.value >= w.value, err: firstError(v, w)}
}

func IntAndIntToBool_LTE(v, w IntValue) BoolValue {
	return BoolValue{value: v.value <= w.value, err: firstError(v, w)}
}

func IntAndIntToInt_ADD(v, w IntValue) IntValue {
	return IntValue{value: v.value + w.value, err: firstError(v, w)}
}

func IntAndIntToInt_SUB(v, w IntValue) IntValue {
	return IntValue{value: v.value - w.value, err: firstError(v, w)}
}

func IntAndIntToInt_MULT(v, w IntValue) IntValue {
	return IntValue{value: v.value * w.value, err: firstError(v, w)}
}

func IntAndIntToInt_DIV(v, w IntValue) IntValue {
	i := IntValue{err: firstError(v, w)}

	if v.value == 0 {
	} else if w.value != 0 {
		i.value = v.value / w.value
	} else if v.value > 0 {
		i.value = math.MaxInt64
	} else {
		i.value = math.MinInt64
	}

	return i
}

func IntAndIntToInt_MOD(v, w IntValue) IntValue {
	i := IntValue{err: firstError(v, w)}

	if v.value != 0 && w.value != 0 {
		i.value = v.value % w.value
	}

	return i
}

// (Float, Float) -> .

func FloatAndFloatToBool_EQ(v, w FloatValue) BoolValue {
	return BoolValue{value: v.value == w.value, err: firstError(v, w)}
}

func FloatAndFloatToBool_NEQ(v, w FloatValue) BoolValue {
	return BoolValue{value: v.value != w.value, err: firstError(v, w)}
}

func FloatAndFloatToBool_GT(v, w FloatValue) BoolValue {
	return BoolValue{value: v.value > w.value, err: firstError(v, w)}
}

func FloatAndFloatToBool_LT(v, w FloatValue) BoolValue {
	return BoolValue{value: v.value < w.value, err: firstError(v, w)}
}

func FloatAndFloatToBool_GTE(v, w FloatValue) BoolValue {
	return BoolValue{value: v.value >= w.value, err: firstError(v, w)}
}

func FloatAndFloatToBool_LTE(v, w FloatValue) BoolValue {
	return BoolValue{value: v.value <= w.value, err: firstError(v, w)}
}

func FloatAndFloatToBool_ADD(v, w FloatValue) FloatValue {
	return FloatValue{value: v.value + w.value, err: firstError(v, w)}
}

func FloatAndFloatToBool_SUB(v, w FloatValue) FloatValue {
	return FloatValue{value: v.value - w.value, err: firstError(v, w)}
}

func FloatAndFloatToBool_MULT(v, w FloatValue) FloatValue {
	return FloatValue{value: v.value * w.value, err: firstError(v, w)}
}

func FloatAndFloatToFloat_DIV(v, w FloatValue) FloatValue {
	f := FloatValue{err: firstError(v, w)}

	if v.value == 0 {
	} else if w.value != 0 {
		f.value = v.value / w.value
	} else if v.value > 0 {
		f.value = math.MaxFloat64
	} else {
		f.value = -math.MaxFloat64
	}

	return f
}

// (String, String) -> .

func StringAndStringToBool_EQ(v, w StringValue) BoolValue {
	return BoolValue{value: v.value == w.value, err: firstError(v, w)}
}

func StringAndStringToBool_NEQ(v, w StringValue) BoolValue {
	return BoolValue{value: v.value != w.value, err: firstError(v, w)}
}

func StringAndStringToBool_GT(v, w StringValue) BoolValue {
	return BoolValue{value: v.value > w.value, err: firstError(v, w)}
}

func StringAndStringToBool_LT(v, w StringValue) BoolValue {
	return BoolValue{value: v.value < w.value, err: firstError(v, w)}
}

func StringAndStringToBool_GTE(v, w StringValue) BoolValue {
	return BoolValue{value: v.value >= w.value, err: firstError(v, w)}
}

func StringAndStringToBool_LTE(v, w StringValue) BoolValue {
	return BoolValue{value: v.value <= w.value, err: firstError(v, w)}
}

func StringAndStringToBool_CONCAT(v, w StringValue) StringValue {
	return StringValue{value: v.value + w.value, err: firstError(v, w)}
}

// N-ary functions

// (Bool..) -> .

func BoolsToBool_EQ(v ...BoolValue) BoolValue {
	b := BoolValue{err: firstError(boolValueErrors(v...)...)}

	if len(v) == 0 {
		return b
	}

	for _, val := range v[1:] {
		if v[0].value != val.value {
			return b
		}
	}

	b.value = true
	return b
}

func BoolsToBool_NEQ(v ...BoolValue) BoolValue {
	return BoolValue{
		value: (len(v) == 1) || (len(v) == 2 && v[0].value != v[1].value),
		err:   firstError(boolValueErrors(v...)...),
	}
}

func BoolsToBool_AND(v ...BoolValue) BoolValue {
	b := BoolValue{err: firstError(boolValueErrors(v...)...)}

	if len(v) == 0 {
		return b
	}

	for _, val := range v {
		if !val.value {
			return b
		}
	}

	b.value = true
	return b
}

func BoolsToBool_OR(v ...BoolValue) BoolValue {
	b := BoolValue{err: firstError(boolValueErrors(v...)...)}

	for _, val := range v {
		if val.value {
			b.value = true
			break
		}
	}

	return b
}

// (Int..) -> .

func IntsToInt_ADD(v ...IntValue) IntValue {
	i := IntValue{err: firstError(intValueErrors(v...)...)}

	for _, val := range v {
		i.value += val.value
	}

	return i
}

func IntsToInt_MULT(v ...IntValue) IntValue {
	i := IntValue{err: firstError(intValueErrors(v...)...)}

	for _, val := range v {
		i.value *= val.value
	}

	return i
}

// (String..) -> .

func StringsToString_CONCAT(v ...StringValue) StringValue {
	s := StringValue{err: firstError(stringValueErrors(v...)...)}

	for _, val := range v {
		s.value += val.value
	}

	return s
}

// Error handling

type errHolder interface {
	Err() error
}

func boolValueErrors(vs ...BoolValue) []errHolder {
	a := make([]errHolder, len(vs))
	for i, v := range vs {
		a[i] = v
	}
	return a
}

func intValueErrors(vs ...IntValue) []errHolder {
	a := make([]errHolder, len(vs))
	for i, v := range vs {
		a[i] = v
	}
	return a
}

func floatValueErrors(vs ...FloatValue) []errHolder {
	a := make([]errHolder, len(vs))
	for i, v := range vs {
		a[i] = v
	}
	return a
}

func stringValueErrors(vs ...StringValue) []errHolder {
	a := make([]errHolder, len(vs))
	for i, v := range vs {
		a[i] = v
	}
	return a
}

func firstError(values ...errHolder) error {
	for _, v := range values {
		if err := v.Err(); err != nil {
			return err
		}
	}

	return nil
}
