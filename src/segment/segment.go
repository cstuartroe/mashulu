package segment

import (
	"fmt"
	"slices"
)

var Vowels = []string{
	"a",
	"e",
	"o",
	"i",
	"u",
	"ai",
	"au",
}

var Consonants = []string{
	"m",
	"n",
	"ny",
	"p",
	"t",
	"k",
	"kw",
	"b",
	"d",
	"g",
	"gw",
	"ts",
	"ch",
	"j",
	"s",
	"sh",
	"l",
	"y",
	"w",

	"mp",
	"nt",
	"nk",
	"nkw",
	"mb",
	"nd",
	"ng",
	"ngw",
	"nts",
	"nch",
	"nj",
}

var roundedVowels = []string{"o", "u"}
var noIConsonants = []byte{'t', 'd', 's', 'y'}

var segments = slices.Concat(Vowels, Consonants)

const longestSegment = 3

func findSegment(s string, i int) (string, int) {
	for l := longestSegment; l > 0; l-- {
		if i+l > len(s) {
			continue
		}

		if slices.Contains(segments, s[i:i+l]) {
			return s[i : i+l], i + l
		}
	}
	return "", i + 1
}

func Validate(segments []string) error {
	if len(segments)%2 == 1 {
		return fmt.Errorf("odd number of segments")
	}

	for i := 0; i < len(segments); i += 2 {
		cons, vow := segments[i], segments[i+1]

		if !slices.Contains(Consonants, cons) {
			return fmt.Errorf("expected consonant: %q", cons+vow)
		}
		if !slices.Contains(Vowels, vow) {
			return fmt.Errorf("expected vowel: %q", cons+vow)
		}

		illegalSequence := false
		lastConsonantLetter := cons[len(cons)-1]
		if lastConsonantLetter == 'w' && slices.Contains(roundedVowels, vow) {
			illegalSequence = true
		}
		if slices.Contains(noIConsonants, lastConsonantLetter) && vow == "i" {
			illegalSequence = true
		}
		if illegalSequence {
			return fmt.Errorf("illegal sequence: %q", cons+vow)
		}
	}

	return nil
}

func Segment(s string) []string {
	var out []string

	i := 0
	for i < len(s) {
		var seg string
		seg, i = findSegment(s, i)
		if seg != "" {
			out = append(out, seg)
		}
	}

	return out
}

func SegmentAndValidate(s string) ([]string, error) {
	segments := Segment(s)

	return segments, Validate(segments)
}
