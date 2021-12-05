package hw02unpackstring

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(str string) (string, error) {
	if str == "" {
		return "", nil
	}
	if unicode.IsDigit(rune(str[0])) {
		return "", ErrInvalidString
	}
	builder := &strings.Builder{}
	temp := str
	const searching = `0123456789\`
	for {
		index := strings.IndexAny(temp, searching)
		if index == -1 { // ничего не нашли, выходим из цикла
			builder.WriteString(temp)
			break
		}
		if index > 1 { // собираем все символы до целевого
			builder.WriteString(temp[:index-1])
		}
		var err error
		index, err = processString(temp, index, builder)
		if err != nil {
			return "", err
		}
		temp = temp[index+1:]
	}
	return builder.String(), nil
}

func processString(temp string, index int, builder *strings.Builder) (int, error) {
	lastIdx := len(temp) - 1

	if unicode.IsDigit(rune(temp[index])) {
		if index < lastIdx && unicode.IsDigit(rune(temp[index+1])) {
			return 0, ErrInvalidString
		}
		char := temp[index-1]
		number, _ := strconv.Atoi(string(temp[index]))
		builder.WriteString(strings.Repeat(string(char), number))
	}

	if !unicode.IsDigit(rune(temp[index])) { // slash
		if index == lastIdx { // слэш в конце, считаем ошибкой
			return 0, ErrInvalidString
		}
		if index > 0 { // записываем символ перед найденным
			builder.WriteRune(rune(temp[index-1]))
		}
		nextChar := temp[index+1]
		if lastIdx >= index+2 && unicode.IsDigit(rune(temp[index+2])) {
			number, _ := strconv.Atoi(string(temp[index+2]))
			builder.WriteString(strings.Repeat(getRepString(rune(nextChar)), number))
			index += 2
		} else {
			builder.WriteRune(rune(nextChar))
			index++
		}
	}
	return index, nil
}

func getRepString(sym rune) string {
	rep := string(sym)
	if sym == 'n' {
		rep = `\n`
	}
	return rep
}
