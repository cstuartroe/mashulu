package ast

type Postverb interface {
	Post()
}

type Person string

const (
	FirstExclusive Person = "1excl"
	FirstInclusive Person = "1incl"
	Second         Person = "second"
	Proximate      Person = "proximate"
	Reflexive      Person = "reflexive"
	Resumptive     Person = "resumptive"
)

type NounClass string

const (
	Animate          NounClass = "animate"
	SemiAnimate      NounClass = "semi-animate"
	FlatLocation     NounClass = "flat/location"
	SimpleDiminutive NounClass = "simple/diminutive"
	Substance        NounClass = "substance"
	Abstract         NounClass = "abstract"
	Time             NounClass = "time"
	Pair             NounClass = "pair"
)

type NounHead interface {
	Head()
}

func (p Person) Head()     {}
func (nc NounClass) Head() {}

type Case string

const (
	Agent                         Case = "agent"
	Patient                       Case = "patient"
	CauseeInstrumental            Case = "causee/instrumental"
	GerundComparative             Case = "gerund/comparative"
	AllativeExperiencerBenefactor Case = "allative/experiencer/benefactor"
	Locative                      Case = "locative"
	AblativeContrastive           Case = "ablative/contrastive"
	Prolative                     Case = "prolative"
	About                         Case = "about"
	Comitative                    Case = "comitative"
)

type NounModifier interface {
	NM()
}

type RelativeClause struct {
	HeadCase  Case
	Verb      *Verb
	PostVerbs []Postverb
}

func (rrc *RelativeClause) AddPostverb(o Postverb) {
	rrc.PostVerbs = append(rrc.PostVerbs, o)
}

func (rrc RelativeClause) NM() {}

type NounPhrase struct {
	Head      NounHead
	Modifiers []NounModifier
}

type Argument struct {
	NounPhrase *NounPhrase
	Case       Case
}

func (o Argument) Post() {}
func (o Argument) NM()   {}

type Content struct {
	Ident string
}

type TameModifier string

const (
	Negative TameModifier = "negative"
)

type Verb struct {
	TameModifiers []TameModifier
	Stem          Content
}

type Attitudinal string

const (
	Assertive           Attitudinal = "assertive"
	SeekingConfirmation Attitudinal = "seeking confirmation"
	Interrogative       Attitudinal = "interrogative"
	Exhortative         Attitudinal = "exhortative"
)

type Sentence struct {
	Topic       *Argument
	Verb        *Verb
	Postverbs   []Postverb
	Attitudinal Attitudinal
}

func (fc Sentence) NM() {}
