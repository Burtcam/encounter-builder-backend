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
	Perception   string
	Languages    []string
	Senses       []string
	Skils        []string
	Passives     []passive
	Movements    []movement
	Reactions    []reaction
	Melees       []attack
	Ranged       []attack
	SpellCasting []spellCasting
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
	Fort string
	Ref  string
	Will string
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
type spellCasting struct {
}
type innateSpellCasting struct {
	DC     int
	School string
}
