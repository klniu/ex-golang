CREATE TABLE `chinese_word` (
  `word_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(8) NOT NULL,
  `url` varchar(512) NOT NULL,
  PRIMARY KEY (`word_id`),
  KEY `title` (`title`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE `chinese_word_pinyin` (
  `word_pinyin_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `word_id` int(10) unsigned NOT NULL,
  `title` varchar(16) NOT NULL,
  `text` varchar(16) NOT NULL,
  `is_primary` tinyint(3) unsigned NOT NULL,
  PRIMARY KEY (`word_pinyin_id`),
  KEY `word_id` (`word_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ā', 'a');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'á', 'a');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ǎ', 'a');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'à', 'a');

UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ō', 'o');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ó', 'o');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ǒ', 'o');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ò', 'o');

UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ê', 'e');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ē', 'e');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'é', 'e');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ě', 'e');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'è', 'e');

UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ī', 'i');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'í', 'i');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ǐ', 'i');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ì', 'i');

UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ū', 'u');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ú', 'u');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ǔ', 'u');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ù', 'u');

UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ǖ', 'u');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ǘ', 'u');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ǚ', 'u');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ǜ', 'u');
UPDATE chinese_word_pinyin SET `text` = REPLACE(`text`, 'ü', 'u');

CREATE TABLE `pinyin` (
  `pinyin_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` varchar(8) NOT NULL,
  `pinyin` varchar(16) NOT NULL,
  PRIMARY KEY (`pinyin_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

INSERT INTO `pinyin`(`title`, `pinyin`)
SELECT a.title, (SELECT `text` FROM chinese_word_pinyin where word_id = a.word_id and is_primary = 1) 
FROM chinese_word a WHERE a.title <> '';
