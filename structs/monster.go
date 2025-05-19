package structs

type Monster struct {
	Name         string
	Traits       Traits
	Attributes   Attributes
	Level        string
	Saves        Saves
	AClass       AC
	HP           HP
	Immunities   []string
	Weaknesses   []string
	Resistances  []string
	Perception   Perception
	Languages    []string
	Senses       []Sense
	Skills       []Skill
	Passives     []Passive //items.system.type =="actions" && items.
	Movements    []Movement
	Reactions    []Reaction
	Melees       []Attack
	Ranged       []Attack
	SpellCasting []SpellCasting
	Specials     []Special
}
type Perception struct {
	Mod    string
	Detail string
}
type HP struct {
	Detail string
	Value  int
}
type AC struct {
	Value  string
	Detail string
}
type Traits struct {
	Rarity    string
	Size      string
	TraitList []string
}
type Skill struct {
	Name     string
	Value    int
	Specials []SkillSpecial
}
type SkillSpecial struct {
	Value      int
	Label      string
	Predicates []string
}
type Attributes struct {
	Str string
	Dex string
	Con string
	Wis string
	Int string
	Cha string
}
type Saves struct {
	Fort       string
	FortDetail string
	Ref        string
	RefDetail  string
	Will       string
	WillDetail string // exceptions per type
	Exception  string // overall exceptions
}
type Passive struct {
	Name   string
	Text   string
	Traits []string
	Range  string
	Damage string
	DC     string
}
type Reaction struct {
	Name   string
	Text   string
	Traits []string
	Range  string
	Damage string
	DC     string
}
type Special struct {
	Name      string
	Text      string
	Traits    []string
	Range     string
	Damage    string
	Actions   string
	DC        string
	Frequency string
}
type Movement struct {
	Type  string
	Speed string
	Notes string
}
type Attack struct {
	AcountCount int
	ToHit       string
	Damage      string
	Type        string
}

type SpellCasting struct {
	InnateSpellCasting      []InnateSpellCasting
	PreparedSpellCasting    []PreparedSpellCasting
	SpontaneousSpellCasting []SpontaneousSpellCasting
}

// Arcane Innate Spells DC 30; 2nd darkness (at will)
// type: Innate, tradition: arcane, dc: 30, spelluses: [spellUse {name: Darkness, Level 2, description: xjklj, Targets: Nil, School, }]
type InnateSpellCasting struct {
	DC        int
	Tradition string
	SpellUses []SpellUse
	Mod       string
}
type PreparedSpellCasting struct {
	DC        int
	Tradition string
	Slots     []Spell
	Mod       string
}
type SpontaneousSpellCasting struct {
	DC        int
	Tradition string
	SpellList []Spell
	Slots     []slot
	Mod       string
}
type slot struct {
	Level  string
	Number string
}
type Spell struct {
	Name           string
	Level          string
	Description    string
	Range          string
	Area           SpellArea
	Duration       string
	Targets        string
	Traits         []string
	Defense        string
	CastTime       string
	CastComponents []string
	Heightened     string
}
type SpellUse struct {
	Spell Spell
	Uses  string
}
type SpellArea struct {
	Type  string
	Value string
}
type Sense struct {
	Name   string //darkvision, smell, etc
	Range  string // 60 feet
	Acuity string //precise or imprecise
	Detail string
}
