/* SELECT id, username, email, created_at
FROM users
WHERE id = $1; */

BEGIN;
-- name: InsertMonster :one
INSERT INTO monsters (name, level, focus_points, traits_rarity, traits_size, ac_value, ac_detail, hp_value, hp_detail)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id;
-- name: InsertAttributes :exec
INSERT INTO attributes (monster_id, str, dex, con, wis, int, cha)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: InsertSaves :exec
INSERT INTO saves (monster_id, fort, fort_detail, ref, ref_detail, will, will_detail, exception)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name:InsertImmunities 
INSERT INTO monster_immunities (monster_id, immunity)
VALUES ($1, $2), ($3, $4);
-- name: InsertDamageModifiers
INSERT INTO monster_damage_modifiers (monster_id, modifier_category, value, damage_type)
VALUES ($1, $2 $3, $4), ($5, $2, $6, $7);

INSERT INTO monster_actions (monster_id, action_type, name, text, actions, category, rarity)
VALUES (new_monster_id, 'action', 'Flaming Bite', 'The Fire Drake bites with intense heat.', '1', 'Offensive', 'Rare');

INSERT INTO monster_attacks (monster_id, attack_category, name, attack_type, to_hit_bonus)
VALUES (new_monster_id, 'melee', 'Fiery Claws', 'Slash', '+10');

INSERT INTO items (monster_id, id, name, category, description, level, rarity, bulk, quantity)
VALUES (new_monster_id, 'item_001', 'Draconic Gem', 'Treasure', 'A rare gemstone infused with dragon energy.', '5', 'Rare', 'Light', '1');

INSERT INTO innate_spell_casting (monster_id, dc, tradition, mod, name, description)
VALUES (new_monster_id, 18, 'Arcane', '+5', 'Dragon’s Breath', 'The Fire Drake expels a cone of flame.');

COMMIT; 