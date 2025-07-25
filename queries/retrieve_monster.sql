-- name: GetFullMonsterByID :one
SELECT row_to_json(monster_data)
FROM (
  SELECT m.*,
    (
      SELECT json_agg(mi)
      FROM monster_immunities mi
      WHERE mi.monster_id = m.id
    ) AS immunities,

    (
      SELECT json_agg(
        json_build_object(
          'id', md.id,
          'modifier_category', md.modifier_category,
          'value', md.value,
          'damage_type', md.damage_type,
          'exceptions', (
            SELECT json_agg(mme.exception)
            FROM monster_modifier_exceptions mme
            WHERE mme.modifier_id = md.id
          ),
          'doubles', (
            SELECT json_agg(mmd.double_value)
            FROM monster_modifier_doubles mmd
            WHERE mmd.modifier_id = md.id
          )
        )
      )
      FROM monster_damage_modifiers md
      WHERE md.monster_id = m.id
    ) AS damage_modifiers,

    (
      SELECT json_agg(ml.language)
      FROM monster_languages ml
      WHERE ml.monster_id = m.id
    ) AS languages,

    (
      SELECT json_agg(ms)
      FROM monster_senses ms
      WHERE ms.monster_id = m.id
    ) AS senses,

    (
      SELECT json_agg(
        json_build_object(
          'id', msk.id,
          'name', msk.name,
          'value', msk.value,
          'specials', (
            SELECT json_agg(mss)
            FROM monster_skill_specials mss
            WHERE mss.skill_id = msk.id
          )
        )
      )
      FROM monster_skills msk
      WHERE msk.monster_id = m.id
    ) AS skills,

    (
      SELECT json_agg(mm)
      FROM monster_movements mm
      WHERE mm.monster_id = m.id
    ) AS movements,

    (
      SELECT json_agg(
        json_build_object(
          'id', ma.id,
          'action_type', ma.action_type,
          'name', ma.name,
          'text', ma.text,
          'actions', ma.actions,
          'category', ma.category,
          'rarity', ma.rarity,
          'dc', ma.dc,
          'traits', (
            SELECT json_agg(mat.trait)
            FROM monster_action_traits mat
            WHERE mat.monster_action_id = ma.id
          )
        )
      )
      FROM monster_actions ma
      WHERE ma.monster_id = m.id
    ) AS actions,

    (
      SELECT json_agg(
        json_build_object(
          'id', ma2.id,
          'attack_category', ma2.attack_category,
          'name', ma2.name,
          'attack_type', ma2.attack_type,
          'to_hit_bonus', ma2.to_hit_bonus,
          'effects_custom_string', ma2.effects_custom_string,
          'effects_values', ma2.effects_values,
          'damage_blocks', (
            SELECT json_agg(adb)
            FROM attack_damage_blocks adb
            WHERE adb.attack_id = ma2.id
          )
        )
      )
      FROM monster_attacks ma2
      WHERE ma2.monster_id = m.id
    ) AS attacks,

    (
      SELECT json_agg(
        json_build_object(
          'id', fsc.id,
          'dc', fsc.dc,
          'mod', fsc.mod,
          'tradition', fsc.tradition,
          'spellcasting_id', fsc.spellcasting_id,
          'name', fsc.name,
          'description', fsc.description,
          'cast_level', fsc.cast_level,
          'spells', (
            SELECT json_agg(fss.spell_id)
            FROM focus_spell_casting_spells fss
            WHERE fss.focus_spell_casting_id = fsc.id
          )
        )
      )
      FROM focus_spell_casting fsc
      WHERE fsc.monster_id = m.id
    ) AS focus_spell_casting,

    (
      SELECT json_agg(
        json_build_object(
          'id', isc.id,
          'dc', isc.dc,
          'tradition', isc.tradition,
          'mod', isc.mod,
          'spellcasting_id', isc.spellcasting_id,
          'name', isc.name,
          'description', isc.description,
          'uses', (
            SELECT json_agg(iu)
            FROM innate_spell_uses iu
            WHERE iu.innate_spell_casting_id = isc.id
          )
        )
      )
      FROM innate_spell_casting isc
      WHERE isc.monster_id = m.id
    ) AS innate_spell_casting,

    (
      SELECT json_agg(
        json_build_object(
          'id', psc.id,
          'dc', psc.dc,
          'tradition', psc.tradition,
          'mod', psc.mod,
          'spellcasting_id', psc.spellcasting_id,
          'description', psc.description,
          'slots', (
            SELECT json_agg(
              json_build_object(
                'id', psl.id,
                'level', psl.level,
                'spell_id', psl.spell_id
              )
            )
            FROM prepared_slots psl
            WHERE psl.prepared_spell_casting_id = psc.id
          )
        )
      )
      FROM prepared_spell_casting psc
      WHERE psc.monster_id = m.id
    ) AS prepared_spell_casting,

    (
      SELECT json_agg(
        json_build_object(
          'id', ssc.id,
          'dc', ssc.dc,
          'id_string', ssc.id_string,
          'tradition', ssc.tradition,
          'mod', ssc.mod,
          'spontaneous_slots', (
            SELECT json_agg(
              json_build_object(
                'id', ssl.id,
                'level', ssl.level,
                'casts', ssl.casts
              )
            )
            FROM spontaneous_slots ssl
            WHERE ssl.spontaneous_spell_casting_id = ssc.id
          ),
          'spontaneous_spell_list', (
            SELECT json_agg(
              json_build_object(
                'id', ssl2.id,
                'spell_id', ssl2.spell_id
              )
            )
            FROM spontaneous_spell_list ssl2
            WHERE ssl2.spontaneous_spell_casting_id = ssc.id
          )
        )
      )
      FROM spontaneous_spell_casting ssc
      WHERE ssc.monster_id = m.id
    ) AS spontaneous_spell_casting,

    (
      SELECT json_agg(
        json_build_object(
          'id', i.id,
          'name', i.name,
          'category', i.category,
          'description', i.description,
          'level', i.level,
          'type', i.type,
          'rarity', i.rarity,
          'size', i.size,
          'range', i.range,
          'reload', i.reload,
          'bulk', i.bulk,
          'quantity', i.quantity,
          'price_per', i.price_per,
          'price_cp', i.price_cp,
          'price_sp', i.price_sp,
          'price_gp', i.price_gp,
          'price_pp', i.price_pp,
          'traits', (
            SELECT json_agg(it.trait)
            FROM item_traits it
            WHERE it.item_id = i.id
          )
        )
      )
      FROM items i
      WHERE i.monster_id = m.id
    ) AS items

  FROM monsters m
  WHERE m.id = $1
) monster_data;

-- name: SearchMonsterByName :many
SELECT row_to_json(monster_data)
FROM (
  SELECT m.*
  FROM monsters m
  WHERE m.name ILIKE '%' || $1 || '%'
) monster_data;

-- name: GetMonstersByTrait :many
SELECT m.*
FROM monsters m
JOIN monster_traits mt ON m.id = mt.monster_id
WHERE mt.trait = $1;


-- name: GetMonstersByLevelRange :many
SELECT row_to_json(monster_data)
FROM (
  SELECT m.*,
    (
      SELECT json_agg(mi)
      FROM monster_immunities mi
      WHERE mi.monster_id = m.id
    ) AS immunities,

    (
      SELECT json_agg(
        json_build_object(
          'id', md.id,
          'modifier_category', md.modifier_category,
          'value', md.value,
          'damage_type', md.damage_type,
          'exceptions', (
            SELECT json_agg(mme.exception)
            FROM monster_modifier_exceptions mme
            WHERE mme.modifier_id = md.id
          ),
          'doubles', (
            SELECT json_agg(mmd.double_value)
            FROM monster_modifier_doubles mmd
            WHERE mmd.modifier_id = md.id
          )
        )
      )
      FROM monster_damage_modifiers md
      WHERE md.monster_id = m.id
    ) AS damage_modifiers,

    (
      SELECT json_agg(ml.language)
      FROM monster_languages ml
      WHERE ml.monster_id = m.id
    ) AS languages,

    (
      SELECT json_agg(ms)
      FROM monster_senses ms
      WHERE ms.monster_id = m.id
    ) AS senses,

    (
      SELECT json_agg(
        json_build_object(
          'id', msk.id,
          'name', msk.name,
          'value', msk.value,
          'specials', (
            SELECT json_agg(mss)
            FROM monster_skill_specials mss
            WHERE mss.skill_id = msk.id
          )
        )
      )
      FROM monster_skills msk
      WHERE msk.monster_id = m.id
    ) AS skills,

    (
      SELECT json_agg(mm)
      FROM monster_movements mm
      WHERE mm.monster_id = m.id
    ) AS movements,

    (
      SELECT json_agg(
        json_build_object(
          'id', ma.id,
          'action_type', ma.action_type,
          'name', ma.name,
          'text', ma.text,
          'actions', ma.actions,
          'category', ma.category,
          'rarity', ma.rarity,
          'dc', ma.dc,
          'traits', (
            SELECT json_agg(mat.trait)
            FROM monster_action_traits mat
            WHERE mat.monster_action_id = ma.id
          )
        )
      )
      FROM monster_actions ma
      WHERE ma.monster_id = m.id
    ) AS actions,

    (
      SELECT json_agg(
        json_build_object(
          'id', ma2.id,
          'attack_category', ma2.attack_category,
          'name', ma2.name,
          'attack_type', ma2.attack_type,
          'to_hit_bonus', ma2.to_hit_bonus,
          'effects_custom_string', ma2.effects_custom_string,
          'effects_values', ma2.effects_values,
          'damage_blocks', (
            SELECT json_agg(adb)
            FROM attack_damage_blocks adb
            WHERE adb.attack_id = ma2.id
          )
        )
      )
      FROM monster_attacks ma2
      WHERE ma2.monster_id = m.id
    ) AS attacks,

    (
      SELECT json_agg(
        json_build_object(
          'id', fsc.id,
          'dc', fsc.dc,
          'mod', fsc.mod,
          'tradition', fsc.tradition,
          'spellcasting_id', fsc.spellcasting_id,
          'name', fsc.name,
          'description', fsc.description,
          'cast_level', fsc.cast_level,
          'spells', (
            SELECT json_agg(fss.spell_id)
            FROM focus_spell_casting_spells fss
            WHERE fss.focus_spell_casting_id = fsc.id
          )
        )
      )
      FROM focus_spell_casting fsc
      WHERE fsc.monster_id = m.id
    ) AS focus_spell_casting,

    (
      SELECT json_agg(
        json_build_object(
          'id', isc.id,
          'dc', isc.dc,
          'tradition', isc.tradition,
          'mod', isc.mod,
          'spellcasting_id', isc.spellcasting_id,
          'name', isc.name,
          'description', isc.description,
          'uses', (
            SELECT json_agg(iu)
            FROM innate_spell_uses iu
            WHERE iu.innate_spell_casting_id = isc.id
          )
        )
      )
      FROM innate_spell_casting isc
      WHERE isc.monster_id = m.id
    ) AS innate_spell_casting,

    (
      SELECT json_agg(
        json_build_object(
          'id', psc.id,
          'dc', psc.dc,
          'tradition', psc.tradition,
          'mod', psc.mod,
          'spellcasting_id', psc.spellcasting_id,
          'description', psc.description,
          'slots', (
            SELECT json_agg(
              json_build_object(
                'id', psl.id,
                'level', psl.level,
                'spell_id', psl.spell_id
              )
            )
            FROM prepared_slots psl
            WHERE psl.prepared_spell_casting_id = psc.id
          )
        )
      )
      FROM prepared_spell_casting psc
      WHERE psc.monster_id = m.id
    ) AS prepared_spell_casting,

    (
      SELECT json_agg(
        json_build_object(
          'id', ssc.id,
          'dc', ssc.dc,
          'id_string', ssc.id_string,
          'tradition', ssc.tradition,
          'mod', ssc.mod,
          'spontaneous_slots', (
            SELECT json_agg(
              json_build_object(
                'id', ssl.id,
                'level', ssl.level,
                'casts', ssl.casts
              )
            )
            FROM spontaneous_slots ssl
            WHERE ssl.spontaneous_spell_casting_id = ssc.id
          ),
          'spontaneous_spell_list', (
            SELECT json_agg(
              json_build_object(
                'id', ssl2.id,
                'spell_id', ssl2.spell_id
              )
            )
            FROM spontaneous_spell_list ssl2
            WHERE ssl2.spontaneous_spell_casting_id = ssc.id
          )
        )
      )
      FROM spontaneous_spell_casting ssc
      WHERE ssc.monster_id = m.id
    ) AS spontaneous_spell_casting,

    (
      SELECT json_agg(
        json_build_object(
          'id', i.id,
          'name', i.name,
          'category', i.category,
          'description', i.description,
          'level', i.level,
          'type', i.type,
          'rarity', i.rarity,
          'size', i.size,
          'range', i.range,
          'reload', i.reload,
          'bulk', i.bulk,
          'quantity', i.quantity,
          'price_per', i.price_per,
          'price_cp', i.price_cp,
          'price_sp', i.price_sp,
          'price_gp', i.price_gp,
          'price_pp', i.price_pp,
          'traits', (
            SELECT json_agg(it.trait)
            FROM item_traits it
            WHERE it.item_id = i.id
          )
        )
      )
      FROM items i
      WHERE i.monster_id = m.id
    ) AS items
  FROM monsters m
  WHERE m.level BETWEEN $1 AND $2
) monster_data;

