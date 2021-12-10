package hw03frequencyanalysis

import (
	"regexp"
	"sort"
	"strings"
)

// MaxResultCount кол-во элементов, которое возвращать в результате работы Top10
const MaxResultCount = 10

type wordCount struct {
	word  string
	count int
}

//Top10 функция, принимающая на вход строку с текстом и
//возвращающую слайс с 10-ю наиболее часто встречаемыми в тексте словами
func Top10(s string) []string {
	if s == "" {
		return []string{}
	}
	pieces := strings.Split(fixString(s), " ")
	if len(pieces) == 1 {
		return []string{s}
	}
	data := map[string]*wordCount{}
	for _, word := range pieces {
		tw := strings.ToLower(trimWordRight(trimWordLeft(word)))
		if tw == "" {
			continue
		}
		if wc, ok := data[tw]; ok {
			wc.count++
		} else {
			data[tw] = &wordCount{tw, 1}
		}
	}

	wordCounts := []*wordCount{}
	for _, wc := range data {
		wordCounts = append(wordCounts, wc)
	}
	sortSlice(wordCounts)
	return cutSlice(wordCounts)
}

func cutSlice(wordCounts []*wordCount) []string {
	ret := []string{}
	for i, wc := range wordCounts {
		ret = append(ret, wc.word)
		if i == MaxResultCount-1 {
			break
		}
	}
	return ret
}

func sortSlice(wordCounts []*wordCount) {
	sort.Slice(wordCounts, func(i, j int) bool {
		ic, jc := wordCounts[i], wordCounts[j]
		if ic.count == jc.count {
			return ic.word < jc.word
		}
		return ic.count > jc.count
	})
}

func fixString(s string) string {
	symbols := []string{"\n", "\t", "\r"}
	ret := s
	for _, sym := range symbols {
		ret = strings.ReplaceAll(ret, sym, " ")
	}
	return ret
}

func trimWordLeft(s string) string {
	re := regexp.MustCompile(`^[,.!\-]+`)
	return string(re.ReplaceAll([]byte(s), []byte("")))
}

func trimWordRight(s string) string {
	re := regexp.MustCompile(`[,.!\-]+$`)
	return string(re.ReplaceAll([]byte(s), []byte("")))
}
