package hw03frequencyanalysis

import (
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
)

var ErrEmptyString = errors.New("empty string")

type Top struct {
	Word string
	Cnt  int
}

var re = regexp.MustCompile(`^[",.!-]+|[",.!-]+$`)

func Top10(s string) []string {
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

	sl := make([]Top, 0, len(wds))

	for k, v := range wds {
		sl = append(sl, Top{Word: k, Cnt: v})
	}

	sort.Slice(sl, func(i, j int) bool {
		if sl[i].Cnt != sl[j].Cnt {
			return sl[i].Cnt > sl[j].Cnt
		}
		return sl[i].Word < sl[j].Word
	})

	// get top10

	top10 := make([]string, 0, 10)

	for i, v := range sl {
		if i > 9 {
			break
		}
		top10 = append(top10, v.Word)
	}

	return top10
}

// Word cleaning

func wordCleaning(s string) (string, error) {
	if s == "" {
		return "", ErrEmptyString
	}

	// clearing a word of punctuation marks
	cleaned := re.ReplaceAllString(strings.ToLower(s), "")

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
