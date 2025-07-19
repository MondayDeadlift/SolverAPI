BEGIN;

ALTER TABLE katas ADD COLUMN url VARCHAR(255) NOT NULL DEFAULT '';
UPDATE katas SET url = 'https://www.codewars.com/kata/' || slug;

COMMIT;