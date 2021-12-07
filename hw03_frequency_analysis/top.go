package hw03frequencyanalysis

import (
	"sort"
	"strings"
)

type wordCount struct {
	word  string
	count int
}

func Top10(s string) []string {
	pieces := strings.Split(fixString(s), " ")
	if len(pieces) == 1 {
		return []string{s}
	}
	data := map[string]*wordCount{}
	for _, word := range pieces {
		if word == "" {
			continue
		}
		if wc, ok := data[word]; ok {
			wc.count++
		} else {
			data[word] = &wordCount{word, 1}
		}
	}
	var wordCounts []*wordCount
	for _, wc := range data {
		wordCounts = append(wordCounts, wc)
	}
	sort.Slice(wordCounts, func(i, j int) bool {
		icount, jcount := wordCounts[i], wordCounts[j]
		if icount.count == jcount.count {
			return icount.word > jcount.word
		}
		return icount.count > jcount.count
	})
	var ret []string
	for i, wc := range wordCounts {
		ret = append(ret, wc.word)
		if i == 9 {
			break
		}
	}
	return ret
}

func fixString(s string) (ret string) {
	symbols := []string{"\n", "\t", "\r"}
	for _, sym := range symbols {
		ret = strings.ReplaceAll(s, sym, " ")
	}
	return
}
