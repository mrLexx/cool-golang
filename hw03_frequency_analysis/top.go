package hw03frequencyanalysis

import (
	"errors"
	"fmt"
	"sort"
	"strings"
)

var ErrEmptyString = errors.New("empty string")

type word struct {
	word string
	cnt  int
}

func Top10(s string) []string {
	return Top(s, 10)
}

func Top(s string, cnt int) []string {
	if s == "" {
		return nil
	}

	wds := make(map[string]int)

	// collect all words in map with quantity

	for _, v := range strings.Fields(s) {
		if cl, err := wordCleaning(v); err == nil {
			wds[cl]++
		}
	}

	// move to slice for sort

	sl := make([]word, 0, len(wds))

	for k, v := range wds {
		sl = append(sl, word{word: k, cnt: v})
	}

	sort.Slice(sl, func(i, j int) bool {
		if sl[i].cnt != sl[j].cnt {
			return sl[i].cnt > sl[j].cnt
		}
		return sl[i].word < sl[j].word
	})

	// get top10

	top := make([]string, 0, cnt)

	for i, v := range sl {
		if i >= cnt {
			break
		}
		top = append(top, v.word)
	}

	return top
}

// Word cleaning

func wordCleaning(s string) (string, error) {
	if s == "" {
		return "", ErrEmptyString
	}

	// clearing a word of punctuation marks
	cleaned := strings.Trim(strings.ToLower(s), ",.!-?\"'`")

	// single punctuation mark, len(punctuation_mark)==1!
	if cleaned == "" && len(s) == 1 {
		return "", fmt.Errorf("punctuation mark `%s`: %w", s, ErrEmptyString)
	}

	// word from punctuation mark, its good
	if cleaned == "" {
		return s, nil
	}

	return cleaned, nil
}
