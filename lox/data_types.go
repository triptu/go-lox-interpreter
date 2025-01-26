package lox

import (
	"fmt"
	"strconv"
	"strings"
)

type dataType interface {
	getMethod(name token) callable
}

type loxList struct {
	elements []any
}

var _ dataType = &loxList{}

func getLoxList(elements []any) *loxList {
	return &loxList{elements: elements}
}

func (l *loxList) getMethod(name token) callable {
	arityCnt, method := l.getMethodAndArity(name)
	if method == nil {
		return nil
	}
	return nativeFunction{
		arityCnt: arityCnt,
		fn: func(i interpreter, a []any) (any, error) {
			return method(a), nil
		},
	}
}

func (l *loxList) getMethodAndArity(name token) (int, func(args []any) any) {
	switch name.lexeme {
	case "append":
		return 1, l.append
	case "extend":
		return 1, l.extend
	case "pop":
		return 0, l.pop
	case "remove":
		return 1, l.remove
	case "insert":
		return 2, l.insert
	case "concat":
		return 1, l.concat
	default:
		return 0, nil
	}
}

func (l *loxList) getAtIndex(index int) any {
	if index < 0 {
		index = len(l.elements) + index
	}
	if index >= len(l.elements) {
		logRuntimeError(token{}, "Index out of bounds")
		return nil
	}
	return l.elements[index]
}

func (l *loxList) append(args []any) any {
	l.elements = append(l.elements, args[0])
	return l
}

func (l *loxList) extend(args []any) any {
	l.elements = append(l.elements, args[0].(*loxList).elements...)
	return l
}

func (l *loxList) pop(args []any) any {
	last := len(l.elements) - 1
	val := l.elements[last]
	l.elements = l.elements[:last]
	return val
}

func (l *loxList) remove(args []any) any {
	index := int(args[0].(float64))
	l.elements = append(l.elements[:index], l.elements[index+1:]...)
	return l
}

func (l *loxList) insert(args []any) any {
	index := int(args[0].(float64))
	element := args[1]
	l.elements = append(l.elements, nil)
	copy(l.elements[index+1:], l.elements[index:])
	l.elements[index] = element
	return l
}

// creates a new list instead of modifying the existing one like extend
func (l *loxList) concat(args []any) any {
	other := args[0].(*loxList)
	newList := loxList{}
	newList.elements = append(newList.elements, l.elements...)
	newList.elements = append(newList.elements, other.elements...)
	return &newList
}

func (l *loxList) String() string {
	var sb strings.Builder
	sb.WriteString("[")
	isFirst := true
	for _, elem := range l.elements {
		if !isFirst {
			sb.WriteString(", ")
		}
		isFirst = false
		sb.WriteString(arrayElementToStr(elem))
	}
	sb.WriteString("]")
	return sb.String()
}

func arrayElementToStr(literal interface{}) string {
	if literal == nil {
		return "nil"
	}

	switch literal := literal.(type) {
	case float64:
		return strconv.FormatFloat(literal, 'f', -1, 64)
	case *loxList:
		return literal.String()
	default:
		return fmt.Sprintf("\"%s\"", getLiteralStr(literal))
	}
}
