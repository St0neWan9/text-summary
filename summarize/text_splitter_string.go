package summarize

import "unicode"

type StringTextSplitter struct {
	Punctuations []string
}

func (d StringTextSplitter) Sentences(text string) []string {
	buf := getBuffer()
	defer bufferPool.Put(buf)

	sentences := []string{}
	newSentence := true
	lastNonWhiteSpace := -1

	for _, r := range text {
		if d.OneOfPunct(string(r), d.Punctuations) {
			if buf.Len() > 0 {
				if lastNonWhiteSpace > 0 {
					buf.Truncate(lastNonWhiteSpace)
					buf.WriteRune(r)
					sentences = append(sentences, buf.String())
				}
				buf.Reset()
				newSentence = true
			}
		} else {
			isSpace := unicode.IsSpace(r)
			if newSentence && isSpace {
				continue
			}
			newSentence = false
			buf.WriteRune(r)
			if !isSpace {
				lastNonWhiteSpace = buf.Len()
			}
		}
	}

	if buf.Len() > 0 && lastNonWhiteSpace > 0 {
		buf.Truncate(lastNonWhiteSpace)
		sentences = append(sentences, buf.String())
		buf.Reset()
	}

	return sentences
}

func (d StringTextSplitter) Words(text string) []string {
	buf := getBuffer()
	defer bufferPool.Put(buf)

	words := []string{}

	for _, r := range text {
		if unicode.IsLetter(r) || unicode.IsNumber(r) {
			buf.WriteRune(r)
		} else if !unicode.IsOneOf([]*unicode.RangeTable{unicode.Hyphen}, r) {
			if buf.Len() > 0 {
				words = append(words, buf.String())
			}
			buf.Reset()
		}
	}

	if buf.Len() > 0 {
		words = append(words, buf.String())
	}

	return words
}

func (d StringTextSplitter) OneOfPunct(r string, punct []string) bool {
	for _, p := range punct {
		if p == r {
			return true
		}
	}
	return false
}

func (d StringTextSplitter) OneOfStartQuote(r string, quotes [][]string) int {
	for i, q := range quotes {
		if q[0] == r {
			return i
		}
	}
	return -1
}
