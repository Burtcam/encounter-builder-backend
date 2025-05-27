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
	Movements    []Movement
	Actions      []Action
	FreeActions  []FreeAction
	Reactions    []Reaction
	Passives     []Passive //items.system.type =="actions" && items.
	Melees       []Attack
	Ranged       []Attack
	SpellCasting SpellCasting
	FocusPoints  int
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
	Name     string
	Text     string
	Traits   []string
	DC       string
	Category string
	Rarity   string
}
type Reaction struct {
	Name     string
	Text     string
	Traits   []string
	Rarity   string
	Category string
}
type Action struct {
	Name     string
	Text     string
	Traits   []string
	Actions  string
	Category string
	Rarity   string
}
type FreeAction struct {
	Name     string
	Text     string
	Traits   []string
	Category string
	Rarity   string
}
type Movement struct {
	Type  string
	Speed string
	Notes string
}
type Attack struct {
	Name         string
	Type         string
	ToHitBonus   string
	DamageBlocks []DamageBlock
	Traits       []string
	Effects      DamageEffect
}
type DamageBlock struct {
	DamageRoll string
	DamageType string
}
type DamageEffect struct {
	CustomString string
	Value        []string
}

type SpellCasting struct {
	InnateSpellCasting      []InnateSpellCasting
	PreparedSpellCasting    []PreparedSpellCasting
	SpontaneousSpellCasting []SpontaneousSpellCasting
	FocusSpellCasting       []FocusSpellCasting
}

type FocusSpellCasting struct {
	DC             int
	Mod            string
	Tradition      string
	ID             string
	Name           string
	FocusSpellList []Spell
	Description    string
	CastLevel      string
}

// Arcane Innate Spells DC 30; 2nd darkness (at will)
// type: Innate, tradition: arcane, dc: 30, spelluses: [spellUse {name: Darkness, Level 2, description: xjklj, Targets: Nil, School, }]
type InnateSpellCasting struct {
	DC          int
	Tradition   string
	SpellUses   []SpellUse
	Mod         string
	ID          string
	Description string
}
type PreparedSpellCasting struct {
	DC          int
	Tradition   string
	Slots       []PreparedSlot
	Mod         string
	ID          string
	Description string
}
type SpontaneousSpellCasting struct {
	DC        int
	ID        string
	Tradition string
	SpellList []Spell
	Slots     []Slot
	Mod       string
}
type Slot struct {
	Level string
	Casts string
}
type PreparedSlot struct {
	Level   string
	SpellID string
	Spell   Spell
}
type Spell struct {
	ID                          string
	Name                        string
	CastLevel                   string
	SpellBaseLevel              string
	Description                 string
	Range                       string
	Area                        SpellArea
	Duration                    DurationBlock
	Targets                     string
	Traits                      []string
	Defense                     DefenseBlock
	CastTime                    string
	CastRequirements            string
	Rarity                      string
	AtWill                      bool
	SpellCastingBlockLocationID string
	Uses                        string
	Ritual                      bool
	RitualData                  RitualData
}
type RitualData struct {
	PrimaryCheck     string
	SecondaryCasters string
	SecondaryCheck   string
}

type DefenseBlock struct {
	Save  string
	Basic bool
}
type SpellUse struct {
	Spell Spell
	Level int
	Uses  string
}
type DurationBlock struct {
	Sustained bool
	Duration  string
}
type SpellArea struct {
	Type   string
	Value  string
	Detail string
}
type Sense struct {
	Name   string //darkvision, smell, etc
	Range  string // 60 feet
	Acuity string //precise or imprecise
	Detail string
}
