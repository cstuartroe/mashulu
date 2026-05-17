package dictionary

import (
	"fmt"
	"slices"

	"github.com/cstuartroe/mashulu/src/ast"
	"github.com/cstuartroe/mashulu/src/segment"
)

var _dictionary = []struct {
	ident string
	form  string
}{
	{"female", "li"},
	{"nonbinary", "pe"},
	{"male", "to"},

	{"exist", "ni"},
	{"copula", "ki"},
	{"go", "te"},
	{"see", "ji"},
	{"eat", "she"},
	{"speak", "lu"},
	{"build", "che"},
	{"grow", "ko"},
	{"hear", "mbo"},
	{"sleep", "ndu"},
	{"seek", "tsu"},
	{"live", "ju"},

	{"small", "ke"},
	{"large", "bo"},
	{"offspring", "de"},
	{"cold", "gu"},
	{"hot", "chi"},
	{"old", "mu"},
	{"new", "ku"},
	{"short", "mi"},
	{"long", "pu"},
	{"neat", "nji"},

	{"learn", "shu"},
	{"love", "mbi"},

	{"water", "yu"},
	{"time", "shi"},
	{"part", "tu"},
	{"head", "ndo"},

	{"apple", "balu"},
	{"mango", "mango"},
}

var Dictionary = map[string]ast.Content{}

// Maps ident to form
var ReverseDictionary = map[string]string{}

var TameModifiers = map[string]ast.TameModifier{
	"el": ast.Negative,
}

var TameModifierForms = map[ast.TameModifier]string{}

var Attitudinals = map[string]ast.Attitudinal{
	"yo": ast.Assertive,
	"ni": ast.SeekingConfirmation,
	"ki": ast.Interrogative,
	"te": ast.Exhortative,
}

var AttitudinalForms = map[ast.Attitudinal]string{}

func validateForm(form string) {
	segments := segment.Segment(form)
	if !(len(segments) == 2 || len(segments) == 4) {
		panic(fmt.Errorf("invalid dictionary entry: %q", form))
	}
	err := segment.Validate(segments[0:2])
	if err != nil {
		panic(err)
	}
	if segments[0] == "s" {
		panic(fmt.Errorf("content stem starts with /s/: %q", form))
	}
	if len(segments) == 4 {
		if segments[1] != "a" {
			panic(fmt.Errorf("Two-syllable content stem does not have /a/ as first vowel: %q", form))
		}
		err := segment.Validate(segments[2:4])
		if err != nil {
			panic(err)
		}
	} else {
		if segments[1] == "a" {
			panic(fmt.Errorf("One-syllable content stem has /a/ as first vowel: %q", form))
		}
	}
}

func validateTameModifier(form string) {
	segments := segment.Segment(form)
	if !slices.Contains(segment.Vowels, segments[0]) {
		panic(fmt.Errorf("Invalid TAME modifier: %q", form))
	}
	if !slices.Contains(segment.Consonants, segments[1]) {
		panic(fmt.Errorf("Invalid TAME modifier: %q", form))
	}
}

func init() {
	for _, entry := range _dictionary {
		validateForm(entry.form)
		if _, ok := Dictionary[entry.form]; ok {
			panic(fmt.Errorf("duplicate dictionary form: %q", entry.form))
		}
		Dictionary[entry.form] = ast.Content{
			Ident: entry.ident,
		}
		if _, ok := ReverseDictionary[entry.ident]; ok {
			panic(fmt.Errorf("duplicate dictionary ident: %q", entry.ident))
		}
		ReverseDictionary[entry.ident] = entry.form
	}

	for form, tm := range TameModifiers {
		validateTameModifier(form)
		if _, ok := TameModifierForms[tm]; ok {
			panic(fmt.Errorf("duplicate TAME modifier: %q", tm))
		}
		TameModifierForms[tm] = form
	}

	for form, att := range Attitudinals {
		if _, ok := AttitudinalForms[att]; ok {
			panic(fmt.Errorf("duplicate TAME modifier: %q", att))
		}
		AttitudinalForms[att] = form
	}
}
