package lox

import (
	"fmt"
	"strconv"
	"strings"
)

type LoxList struct {
	elements []any
}

func getLoxList(elements []any) *LoxList {
	return &LoxList{elements: elements}
}

// implement python like methods - append, pop, remove, insert, concat

func (l *LoxList) append(element any) *LoxList {
	l.elements = append(l.elements, element)
	return l
}

func (l *LoxList) extend(other *LoxList) *LoxList {
	l.elements = append(l.elements, other.elements...)
	return l
}

func (l *LoxList) pop() any {
	last := len(l.elements) - 1
	val := l.elements[last]
	l.elements = l.elements[:last]
	return val
}

func (l *LoxList) remove(index int) *LoxList {
	l.elements = append(l.elements[:index], l.elements[index+1:]...)
	return l
}

func (l *LoxList) insert(index int, element any) *LoxList {
	l.elements = append(l.elements, nil)
	copy(l.elements[index+1:], l.elements[index:])
	l.elements[index] = element
	return l
}

// creates a new list instead of modifying the existing one like extend
func (l *LoxList) concat(other *LoxList) *LoxList {
	newList := LoxList{}
	newList.elements = append(newList.elements, l.elements...)
	newList.elements = append(newList.elements, other.elements...)
	return &newList
}

func (l *LoxList) String() string {
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
	case *LoxList:
		return literal.String()
	default:
		return fmt.Sprintf("\"%s\"", getLiteralStr(literal))
	}
}
