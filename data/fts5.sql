-- create index table
CREATE VIRTUAL TABLE IF NOT EXISTS ead_AdditiveFTSI USING fts5(id, code, name, status, function, foods, notice, info)
;

-- populate index
INSERT INTO ead_AdditiveFTSI 
    SELECT a.id, a.code,
    (SELECT value_str FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'name') AS name,
    (SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'status') AS status,
    (SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'function') AS function,
    (SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'foods') AS foods,
    (SELECT value_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'notice') AS notice,
    (SELECT value_big_text FROM ead_AdditiveProps WHERE additive_id = a.id AND key_name = 'info') AS info
    FROM ead_Additive AS a
;
