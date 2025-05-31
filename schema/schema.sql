CREATE TABLE monsters (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    level VARCHAR(50),
    focus_points INTEGER,
    -- Traits (one-to-one)
    traits_rarity VARCHAR(50),
    traits_size VARCHAR(50),
    -- Attributes (one-to-one)
    attr_str VARCHAR(10),
    attr_dex VARCHAR(10),
    attr_con VARCHAR(10),
    attr_wis VARCHAR(10),
    attr_int VARCHAR(10),
    attr_cha VARCHAR(10),
    -- Saves (one-to-one)
    saves_fort VARCHAR(20),
    saves_fort_detail TEXT,
    saves_ref VARCHAR(20),
    saves_ref_detail TEXT,
    saves_will VARCHAR(20),
    saves_will_detail TEXT,
    saves_exception TEXT,
    -- AC and HP
    ac_value VARCHAR(50),
    ac_detail TEXT,
    hp_detail TEXT,
    hp_value INTEGER,
    -- Perception
    perception_mod VARCHAR(50),
    perception_detail TEXT
);
CREATE TABLE monster_immunities (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    immunity VARCHAR(100)
);

CREATE TABLE monster_languages (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    language VARCHAR(50)
);
CREATE TABLE monster_damage_modifiers (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    modifier_category VARCHAR(20) CHECK (modifier_category IN ('weakness', 'resistance')),
    value INTEGER,
    damage_type VARCHAR(50)
);

CREATE TABLE monster_modifier_exceptions (
    id SERIAL PRIMARY KEY,
    modifier_id INTEGER REFERENCES monster_damage_modifiers(id) ON DELETE CASCADE,
    exception VARCHAR(100)
);

CREATE TABLE monster_modifier_doubles (
    id SERIAL PRIMARY KEY,
    modifier_id INTEGER REFERENCES monster_damage_modifiers(id) ON DELETE CASCADE,
    double_value VARCHAR(100)
);
CREATE TABLE monster_senses (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    name VARCHAR(50),
    range VARCHAR(50),
    acuity VARCHAR(50),
    detail TEXT
);
CREATE TABLE monster_skills (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    name VARCHAR(100),
    value INTEGER
);

CREATE TABLE monster_skill_specials (
    id SERIAL PRIMARY KEY,
    skill_id INTEGER REFERENCES monster_skills(id) ON DELETE CASCADE,
    value INTEGER,
    label VARCHAR(100),
    predicates TEXT[]  -- storing an array of strings
);
CREATE TABLE monster_movements (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    movement_type VARCHAR(50),
    speed VARCHAR(50),
    notes TEXT
);
CREATE TABLE monster_actions (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    action_type VARCHAR(20) CHECK (action_type IN ('action', 'free_action', 'reaction', 'passive')),
    name VARCHAR(100),
    text TEXT,
    actions VARCHAR(100),  -- used for standard "actions"; leave NULL if not applicable
    category VARCHAR(50),
    rarity VARCHAR(50),
    dc VARCHAR(50)  -- used for passives
);

CREATE TABLE monster_action_traits (
    id SERIAL PRIMARY KEY,
    monster_action_id INTEGER REFERENCES monster_actions(id) ON DELETE CASCADE,
    trait VARCHAR(50)
);
CREATE TABLE monster_attacks (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    attack_category VARCHAR(20) CHECK (attack_category IN ('melee', 'ranged')),
    name VARCHAR(100),
    attack_type VARCHAR(50),
    to_hit_bonus VARCHAR(50),
    effects_custom_string TEXT,
    effects_values TEXT[]  -- array of strings for the DamageEffect.Value field
);

CREATE TABLE attack_damage_blocks (
    id SERIAL PRIMARY KEY,
    attack_id INTEGER REFERENCES monster_attacks(id) ON DELETE CASCADE,
    damage_roll VARCHAR(50),
    damage_type VARCHAR(50)
);
CREATE TABLE spells (
    id VARCHAR(50) PRIMARY KEY,
    name VARCHAR(100),
    cast_level VARCHAR(50),
    spell_base_level VARCHAR(50),
    description TEXT,
    range VARCHAR(100),
    cast_time VARCHAR(50),
    cast_requirements TEXT,
    rarity VARCHAR(50),
    at_will BOOLEAN,
    spell_casting_block_location_id VARCHAR(50),
    uses VARCHAR(50),
    ritual BOOLEAN,
    targets TEXT
);

CREATE TABLE spell_areas (
    id SERIAL PRIMARY KEY,
    spell_id VARCHAR(50) REFERENCES spells(id) ON DELETE CASCADE,
    area_type VARCHAR(50),
    value VARCHAR(50),
    detail TEXT
);

CREATE TABLE spell_durations (
    id SERIAL PRIMARY KEY,
    spell_id VARCHAR(50) REFERENCES spells(id) ON DELETE CASCADE,
    sustained BOOLEAN,
    duration VARCHAR(50)
);

CREATE TABLE spell_defenses (
    id SERIAL PRIMARY KEY,
    spell_id VARCHAR(50) REFERENCES spells(id) ON DELETE CASCADE,
    save VARCHAR(50),
    basic BOOLEAN
);

CREATE TABLE ritual_data (
    id SERIAL PRIMARY KEY,
    spell_id VARCHAR(50) REFERENCES spells(id) ON DELETE CASCADE,
    primary_check VARCHAR(50),
    secondary_casters VARCHAR(50),
    secondary_check VARCHAR(50)
);

CREATE TABLE spell_traits (
    id SERIAL PRIMARY KEY,
    spell_id VARCHAR(50) REFERENCES spells(id) ON DELETE CASCADE,
    trait VARCHAR(50)
);
CREATE TABLE focus_spell_casting (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    dc INTEGER,
    mod VARCHAR(50),
    tradition VARCHAR(50),
    spellcasting_id VARCHAR(50),
    name VARCHAR(100),
    description TEXT,
    cast_level VARCHAR(50)
);

CREATE TABLE focus_spell_casting_spells (
    id SERIAL PRIMARY KEY,
    focus_spell_casting_id INTEGER REFERENCES focus_spell_casting(id) ON DELETE CASCADE,
    spell_id VARCHAR(50) REFERENCES spells(id)
);
CREATE TABLE innate_spell_casting (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    dc INTEGER,
    tradition VARCHAR(50),
    mod VARCHAR(50),
    spellcasting_id VARCHAR(50),
    description TEXT,
    name VARCHAR(100)
);

CREATE TABLE innate_spell_uses (
    id SERIAL PRIMARY KEY,
    innate_spell_casting_id INTEGER REFERENCES innate_spell_casting(id) ON DELETE CASCADE,
    spell_id VARCHAR(50) REFERENCES spells(id),
    level INTEGER,
    uses VARCHAR(50)
);
CREATE TABLE prepared_spell_casting (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    dc INTEGER,
    tradition VARCHAR(50),
    mod VARCHAR(50),
    spellcasting_id VARCHAR(50),
    description TEXT
);

CREATE TABLE prepared_slots (
    id SERIAL PRIMARY KEY,
    prepared_spell_casting_id INTEGER REFERENCES prepared_spell_casting(id) ON DELETE CASCADE,
    level VARCHAR(50),
    spell_id VARCHAR(50) REFERENCES spells(id)
);
CREATE TABLE spontaneous_spell_casting (
    id SERIAL PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    dc INTEGER,
    id_string VARCHAR(50),
    tradition VARCHAR(50),
    mod VARCHAR(50)
);

CREATE TABLE spontaneous_slots (
    id SERIAL PRIMARY KEY,
    spontaneous_spell_casting_id INTEGER REFERENCES spontaneous_spell_casting(id) ON DELETE CASCADE,
    level VARCHAR(50),
    casts VARCHAR(50)
);

CREATE TABLE spontaneous_spell_list (
    id SERIAL PRIMARY KEY,
    spontaneous_spell_casting_id INTEGER REFERENCES spontaneous_spell_casting(id) ON DELETE CASCADE,
    spell_id VARCHAR(50) REFERENCES spells(id)
);
CREATE TABLE items (
    id VARCHAR(50) PRIMARY KEY,
    monster_id INTEGER REFERENCES monsters(id) ON DELETE CASCADE,
    name VARCHAR(100),
    category VARCHAR(50),
    description TEXT,
    level VARCHAR(50),
    type VARCHAR(50),
    rarity VARCHAR(50),
    size VARCHAR(50),
    range VARCHAR(50),
    reload VARCHAR(50),
    bulk VARCHAR(50),
    quantity VARCHAR(50),
    price_per INTEGER,
    price_cp INTEGER,
    price_sp INTEGER,
    price_gp INTEGER,
    price_pp INTEGER
);

CREATE TABLE item_traits (
    id SERIAL PRIMARY KEY,
    item_id VARCHAR(50) REFERENCES items(id) ON DELETE CASCADE,
    trait VARCHAR(50)
);
