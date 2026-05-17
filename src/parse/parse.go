package parse

import (
	"fmt"

	"github.com/cstuartroe/mashulu/src/ast"
	"github.com/cstuartroe/mashulu/src/dictionary"
)

var CaseConsonants = map[string]ast.Case{
	"n":  ast.Agent,
	"k":  ast.Patient,
	"l":  ast.CauseeInstrumental,
	"m":  ast.GerundComparative,
	"g":  ast.AllativeExperiencerBenefactor,
	"ch": ast.Locative,
	"p":  ast.AblativeContrastive,
	"sh": ast.Prolative,
	"b":  ast.About,
	"j":  ast.Comitative,
}

func getDefaultCase(head ast.NounHead) (ast.Case, error) {
	switch p := head.(type) {
	case ast.Person:
		switch p {
		case ast.FirstExclusive, ast.FirstInclusive, ast.Second:
			return ast.Agent, nil
		default:
			return "", fmt.Errorf("%s pronoun does not have default case", p)
		}
	case ast.NounClass:
		switch p {
		case ast.Animate:
			return ast.Agent, nil
		case ast.SemiAnimate:
			return ast.CauseeInstrumental, nil
		case ast.FlatLocation:
			return ast.Locative, nil
		case ast.Abstract:
			return ast.GerundComparative, nil
		case ast.Time:
			return ast.Prolative, nil
		default:
			return ast.Patient, nil
		}
	default:
		return "", fmt.Errorf("Unknown noun head type")
	}
}

var PersonConsonants = map[string]ast.Person{
	"g":  ast.FirstExclusive,
	"gw": ast.FirstInclusive,
	"w":  ast.Second,
	"t":  ast.Proximate,
	"d":  ast.Resumptive,
	"ny": ast.Reflexive,
}

var NounClassConsonants = map[string]ast.NounClass{
	"n":  ast.Animate,
	"l":  ast.SemiAnimate,
	"ch": ast.FlatLocation,
	"k":  ast.SimpleDiminutive,
	"b":  ast.Substance,
	"m":  ast.Abstract,
	"sh": ast.Time,
}

type parser struct {
	segments []string
	i        int
}

func New(segments []string) parser {
	return parser{
		segments: segments,
		i:        0,
	}
}

func (p *parser) hasMore() bool {
	return p.i < len(p.segments)
}

func (p *parser) current() (string, string) {
	if !p.hasMore() {
		return "", ""
	}
	return p.segments[p.i], p.segments[p.i+1]
}

func (p *parser) advance() {
	p.i += 2
}

func (p *parser) Parse() ([]*ast.Sentence, error) {
	out := []*ast.Sentence{}

	topic, err := p.grabSubject()
	if err != nil {
		return nil, err
	}

	for {
		verb, postverbs, err := p.grabClause(true)
		if err != nil {
			return nil, err
		}

		c, v := p.current()
		if c == "s" {
			out = append(out, &ast.Sentence{
				Topic:     topic,
				Verb:      verb,
				Postverbs: postverbs[0 : len(postverbs)-1],
			})
			switch p := postverbs[len(postverbs)-1].(type) {
			case *ast.Argument:
				topic = p
			default:
				return nil, fmt.Errorf("Top-level verb is not preceded by a noun phrase")
			}
		} else if att, ok := dictionary.Attitudinals[c+v]; ok {
			out = append(out, &ast.Sentence{
				Topic:       topic,
				Verb:        verb,
				Postverbs:   postverbs,
				Attitudinal: att,
			})
			p.advance()

			if p.hasMore() {
				topic, err = p.grabSubject()
				if err != nil {
					return nil, err
				}
			} else {
				return out, nil
			}
		} else {
			out = append(out, &ast.Sentence{
				Topic:     topic,
				Verb:      verb,
				Postverbs: postverbs,
			})
			if p.hasMore() {
				topic, err = p.grabSubject()
				if err != nil {
					return nil, err
				}
			} else {
				return out, nil
			}
		}
	}
}

func (p *parser) grabSubject() (*ast.Argument, error) {
	np, err := p.grabNounPhrase()
	if err != nil {
		return nil, err
	}

	var ncase ast.Case
	c, v := p.current()
	if c == "s" {
		ncase, err = getDefaultCase(np.Head)
		if err != nil {
			return nil, err
		}
	} else if v == "i" {
		var ok bool
		ncase, ok = CaseConsonants[c]
		if !ok {
			return nil, fmt.Errorf("Unknown case consonant: %q", c)
		}
		p.advance()
	} else {
		return nil, fmt.Errorf("Unexpected continuation after subject noun phrase: %q", c+v)
	}

	return &ast.Argument{
		NounPhrase: np,
		Case:       ncase,
	}, nil
}

func (p *parser) grabNounPhrase() (*ast.NounPhrase, error) {
	var head ast.NounHead
	modifiers := []ast.NounModifier{}

	c, v := p.current()
	if c == "p" {
		// TODO: relational nouns
	} else if person, ok := PersonConsonants[c]; ok {
		if v == "a" {
			head = person
			p.advance()
		}
	} else if nc, ok := NounClassConsonants[c]; ok {
		if v == "a" {
			head = nc
			p.advance()
		} else if v == "e" {
			head = nc
			p.advance()

			verb, postverbs, err := p.grabClause(false)
			if err != nil {
				return nil, err
			}
			modifiers = append(modifiers, &ast.RelativeClause{
				HeadCase:  ast.Patient,
				Verb:      verb,
				PostVerbs: postverbs,
			})
		} else if v == "o" {
			ncase, err := getDefaultCase(nc)
			if err != nil && ncase != ast.Patient {
				head = nc
				p.advance()

				verb, postverbs, err := p.grabClause(false)
				if err != nil {
					return nil, err
				}
				modifiers = append(modifiers, &ast.RelativeClause{
					HeadCase:  ncase,
					Verb:      verb,
					PostVerbs: postverbs,
				})
			}
		}
	}
	if head == nil {
		return nil, fmt.Errorf("Invalid start to noun phrase: %q", c+v)
	}

	modifier, err := p.grabNounModifier()
	if err != nil {
		return nil, err
	}
	for modifier != nil {
		modifiers = append(modifiers, modifier)
		modifier, err = p.grabNounModifier()
		if err != nil {
			return nil, err
		}
	}

	return &ast.NounPhrase{
		Head:      head,
		Modifiers: modifiers,
	}, nil
}

func (p *parser) grabClause(topLevel bool) (*ast.Verb, []ast.Postverb, error) {
	verb := ast.Verb{}
	var err error

	c, _ := p.current()
	if c == "s" {
		verb.TameModifiers, err = p.grabTameModifiers()
		if err != nil {
			return nil, nil, err
		}
	} else if topLevel {
		return nil, nil, fmt.Errorf("Top-level verb requires TAME complex")
	}
	verb.Stem, err = p.grabContent()
	if err != nil {
		return nil, nil, err
	}

	postverbs := []ast.Postverb{}
	postverb, err := p.grabPostverb()
	if err != nil {
		return nil, nil, err
	}
	for postverb != nil {
		postverbs = append(postverbs, postverb)
		postverb, err = p.grabPostverb()
		if err != nil {
			return nil, nil, err
		}
	}

	return &verb, postverbs, nil
}

func (p *parser) grabTameModifiers() ([]ast.TameModifier, error) {
	var mods []ast.TameModifier

	c, v := p.current()
	p.advance()
	if c != "s" {
		return nil, fmt.Errorf("Invalid start to TAME modifiers: %q", c+v)
	}

	var lastV string
	for v != "a" {
		lastV = v

		c, v = p.current()
		p.advance()

		form := lastV + c
		tm, ok := dictionary.TameModifiers[form]
		if !ok {
			return nil, fmt.Errorf("Unknown TAME modifier: %q", form)
		}
		mods = append(mods, tm)
	}

	return mods, nil
}

func (p *parser) grabContent() (ast.Content, error) {
	c, v := p.current()
	p.advance()

	form := c + v
	if v == "a" {
		c2, v2 := p.current()
		p.advance()
		form += c2 + v2
	}

	content, ok := dictionary.Dictionary[form]
	if !ok {
		return ast.Content{}, fmt.Errorf("Unknown content stem %q", form)
	}
	return content, nil
}

// unlike most "grab" methods, gracefully returns nil, nil if a noun modifier is not next
func (p *parser) grabNounModifier() (ast.NounModifier, error) {
	c, v := p.current()
	ncase, ok := CaseConsonants[c]
	if v == "u" && ok {
		p.advance()

		verb, postverbs, err := p.grabClause(false)
		if err != nil {
			return nil, err
		}

		return &ast.RelativeClause{
			HeadCase:  ncase,
			Verb:      verb,
			PostVerbs: postverbs,
		}, nil
	}

	// TODO: noun-noun modifiers
	// TODO: possessive pronouns

	return nil, nil
}

// unlike most "grab" methods, gracefully returns nil, nil if a noun modifier is not next
func (p *parser) grabPostverb() (ast.Postverb, error) {
	i := p.i

	np, err := p.grabNounPhrase()
	if err == nil {
		c, v := p.current()
		p.advance()

		ncase, ok := CaseConsonants[c]
		if v != "i" || !ok {
			p.i = i
			return nil, nil
		}

		return &ast.Argument{
			NounPhrase: np,
			Case:       ncase,
		}, nil
	}

	p.i = i
	return nil, nil
}
