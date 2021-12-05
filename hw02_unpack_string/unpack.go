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
	builder := strings.Builder{}
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
		lastIdx := len(temp) - 1
		if unicode.IsDigit(rune(temp[index])) {
			if index < lastIdx && unicode.IsDigit(rune(temp[index+1])) {
				return "", ErrInvalidString
			}
			char := temp[index-1]
			number, _ := strconv.Atoi(string(temp[index]))
			builder.WriteString(strings.Repeat(string(char), number))
		} else { // slash
			if index == lastIdx { // слэш в конце, считаем ошибкой
				return "", ErrInvalidString
			}
			if index > 0 { // записываем символ перед найденным
				builder.WriteRune(rune(temp[index-1]))
			}
			nextChar := temp[index+1]
			if lastIdx >= index+2 && unicode.IsDigit(rune(temp[index+2])) {
				number, _ := strconv.Atoi(string(temp[index+2]))
				rep := string(nextChar)
				if nextChar == 'n' {
					rep = `\n`
				}
				builder.WriteString(strings.Repeat(rep, number))
				index += 2
			} else {
				builder.WriteRune(rune(nextChar))
				index++
			}
		}
		temp = temp[index+1:]
	}
	return builder.String(), nil
}
