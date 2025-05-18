package structs

type Monster struct {
	Name         string
	Traits       []string
	Attributes   attributes
	Size         string
	Level        string
	Saves        saves
	AClass       string
	HP           int
	Immunities   []string
	Weaknesses   []string
	Resistances  []string
	Perception   string
	Languages    []string
	Senses       []sense
	Skils        []string
	Passives     []passive
	Movements    []movement
	Reactions    []reaction
	Melees       []attack
	Ranged       []attack
	SpellCasting []spellCasting
	Specials     []special
}
type attributes struct {
	Str string
	Dex string
	Con string
	Wis string
	Int string
	Cha string
}
type saves struct {
	Fort       string
	FortDetail string
	Ref        string
	RefDetail  string
	Will       string
	WillDetail string // exceptions per type
	Exception  string // overall exceptions
}
type passive struct {
	Name   string
	Text   string
	Traits []string
	Range  string
	Damage string
	DC     string
}
type reaction struct {
	Name   string
	Text   string
	Traits []string
	Range  string
	Damage string
	DC     string
}
type special struct {
	Name      string
	Text      string
	Traits    []string
	Range     string
	Damage    string
	Actions   string
	DC        string
	Frequency string
}
type movement struct {
	Type  string
	Speed string
	Notes string
}
type attack struct {
	AcountCount int
	ToHit       string
	Damage      string
	Type        string
}

// Arcane Innate Spells DC 30; 2nd darkness (at will)
// type: Innate, tradition: arcane, dc: 30, spelluses: [spellUse {name: Darkness, Level 2, description: xjklj, Targets: Nil, School, }]
type spellCasting struct {
	DC        int
	Tradition string
	SpellUses []spellUse
	Type      string
}
type spell struct {
	Name           string
	Level          string
	Description    string
	Range          string
	Area           spellArea
	Duration       string
	Targets        string
	Traits         []string
	Defense        string
	CastTime       string
	CastComponents []string
	Heightened     string
}
type spellUse struct {
	Spell spell
	Uses  string
}
type spellArea struct {
	Type  string
	Value string
}
type sense struct {
	Name   string //darkvision, smell, etc
	Range  string // 60 feet
	Acuity string //precise or imprecise

}
