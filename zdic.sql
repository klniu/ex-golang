CREATE TABLE `chinese_word` (
  `word_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `title` char(1) NOT NULL,
  PRIMARY KEY (`word_id`),
  UNIQUE KEY `title` (`title`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `chinese_word_pinyin` (
  `word_pinyin_id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `word_id` int(10) unsigned NOT NULL,
  `pinyin` varchar(16) NOT NULL,
  `is_primary` tinyint(3) unsigned NOT NULL,
  PRIMARY KEY (`word_pinyin_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
