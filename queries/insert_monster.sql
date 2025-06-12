-- name: InsertMonster :one
INSERT INTO monsters (name, 
                        level, 
                        focus_points, 
                        traits_rarity, 
                        traits_size, 
                        attr_str, 
                        attr_dex, 
                        attr_con, 
                        attr_wis,
                        attr_int, 
                        attr_cha, 
                        saves_fort, 
                        saves_fort_detail,
                        saves_ref, 
                        saves_ref_detail, 
                        saves_will, 
                        saves_will_detail,
                        saves_exception, 
                        ac_value, 
                        ac_detail,
                        hp_value, 
                        hp_detail, 
                        perception_mod, 
                        perception_detail)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19, $20, $21, $22, $23, $24)
RETURNING id;

-- name: InsertMonsterImmunities :exec
INSERT INTO monster_immunities (monster_id, immunity)
VALUES ($1, $2);

-- name: InsertMonsterDamageModifier :one
INSERT INTO monster_damage_modifiers (monster_id, modifier_category, value, damage_type)
VALUES ($1, $2, $3, $4)
RETURNING id;

-- name: InsertMonsterModifierExceptions :exec
INSERT INTO monster_modifier_exceptions (modifier_id, exception)
VALUES ($1, $2);

-- name: InsertMonsterModifierDoubles :exec
INSERT INTO monster_modifier_doubles (modifier_id, double_value)
VALUES ($1, $2);

-- name: InsertMonsterLanguages :exec
INSERT INTO monster_languages (monster_id, language)
values ($1, $2);

-- name: InsertMonsterSenses :exec
INSERT INTO monster_senses (monster_id, name, range, acuity, detail)
VALUES ($1, $2, $3, $4, $5); 

-- name: InsertMonsterSkills :one
INSERT INTO monster_skills (monster_id, name, value)
VALUES($1, $2, $3)
RETURNING id; 

-- name: InsertMonsterSkillSpecials :exec
INSERT INTO monster_skill_specials(skill_id, value, label, predicates)
VALUES ($1, $2, $3, $4);

-- name: InsertMonsterMovements :exec
INSERT INTO monster_movements(monster_id, movement_type, speed, notes)
VALUES ($1, $2, $3, $4);

-- name: InsertMonsterAction :one
INSERT INTO monster_actions (monster_id, action_type, name, text, actions, category, rarity, dc)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id; 

-- name: InsertMonsterActionTraits :exec
INSERT INTO monster_action_traits (monster_action_id, trait)
VALUES($1, $2);

-- name: InsertMonsterAttacks :one
INSERT INTO monster_attacks (monster_id, attack_category, name, attack_type, to_hit_bonus, effects_custom_string, effects_values)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id; 

-- name: InsertMonsterAttackDamageBlock :exec
INSERT INTO attack_damage_blocks (attack_id, damage_roll, damage_type)
VALUES ($1, $2, $3);

-- name: InsertSpell :one
INSERT INTO spells (name, cast_level, spell_base_level, description, range, cast_time, cast_requirements, rarity, at_will, spell_casting_block_location_id, uses, ritual, targets)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
RETURNING id; 

-- name: InsertSpellArea :exec
INSERT INTO spell_areas (spell_id, area_type, value, detail)
VALUES($1, $2, $3, $4);

-- name: InsertSpellDuration :exec 
INSERT INTO spell_durations (spell_id, sustained, duration)
VALUES ($1, $2, $3);

-- name: InsertSpellDefences :exec
INSERT INTO spell_defenses (spell_id, save, basic)
VALUES ($1, $2, $3);

-- name: InsertRitualData :exec
INSERT INTO ritual_data (spell_id, primary_check, secondary_casters, secondary_check)
VALUES ($1, $2, $3, $4); 

-- name: InsertSpellTraits :exec
INSERT INTO spell_traits (spell_id, trait)
VALUES ($1, $2); 

-- name: InsertFocusSpellCasting :one
INSERT INTO focus_spell_casting (monster_id, dc, mod, tradition, spellcasting_id, name, description, cast_level)
Values($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING id; 

-- name: InsertFocusSpellsCasts :exec
INSERT INTO focus_spell_casting_spells (focus_spell_casting_id, spell_id)
VALUES ($1, $2); 

-- name: InsertInnateSpellCasting :one 
INSERT INTO innate_spell_casting (monster_id, dc, tradition, mod, spellcasting_id, description, name)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id; 

-- name: InsertInnateSpellUse :exec
INSERT INTO innate_spell_uses (innate_spell_casting_id, spell_id, level, uses)
VALUES ($1, $2, $3, $4);

-- name: InsertPreparedSpellCasting :one
INSERT INTO prepared_spell_casting (monster_id, dc, tradition, mod, spellcasting_id, description)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id; 

-- name: InsertPreparedSlots :exec
INSERT INTO prepared_slots (prepared_spell_casting_id, level, spell_id)
VALUES ($1, $2, $3);

-- name: InsertSpontaneousSpells :one
INSERT INTO spontaneous_spell_casting (monster_id, dc, id_string, tradition, mod)
VALUES ($1, $2, $3, $4, $5)
RETURNING id; 

-- name: InsertSpontaneousSpellSlots :exec
INSERT INTO spontaneous_slots (spontaneous_spell_casting_id, level, casts)
VALUES ($1, $2, $3); 

-- name: InsertSpontaneousSpellList :exec
INSERT INTO spontaneous_spell_list (spontaneous_spell_casting_id, spell_id)
VALUES ($1, $2); 

-- name: InsertItems :one
INSERT INTO items (id, monster_id, name, category, description, level, rarity, bulk, quantity, price_per, price_cp, price_gp, price_sp, price_pp)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14)
RETURNING id; 

-- name: InsertItemTraits :exec
INSERT INTO item_traits (item_id, trait)
VALUES ($1, $2);
