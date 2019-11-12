CREATE TABLE `market` (
  `id` bigint(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(25) DEFAULT NULL,
  `condition` varchar(45) DEFAULT NULL,
  `price` int(10) unsigned DEFAULT NULL,
  `member` varchar(45) DEFAULT NULL,
  `updated` date DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `title_UNIQUE` (`title`),
  UNIQUE KEY `price_UNIQUE` (`price`),
  UNIQUE KEY `member_UNIQUE` (`member`),
  UNIQUE KEY `updated_UNIQUE` (`updated`),
  UNIQUE KEY `condition_UNIQUE` (`condition`)
) ENGINE=InnoDB AUTO_INCREMENT=36 DEFAULT CHARSET=utf8;