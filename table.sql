CREATE TABLE `market` (
  `id` bigint(11) NOT NULL AUTO_INCREMENT,
  `title` varchar(250) DEFAULT NULL,
  `size` smallint DEFAULT NULL,
  `condition` varchar(45) DEFAULT NULL,
  `price` int(10) unsigned DEFAULT NULL,
  `member` varchar(45) DEFAULT NULL,
  `isNew` boolean DEFAULT 1,
  `updated` date DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`),
  UNIQUE KEY `title_UNIQUE` (`title`),
  KEY `isNew_Index` (`isNew`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8;

CREATE TABLE `conditions` (
  `id` bigint(11) NOT NULL AUTO_INCREMENT,
  `size` smallint DEFAULT NULL,
  `keyword` varchar(45) DEFAULT NULL,
  `created` date DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `id_UNIQUE` (`id`)
) ENGINE = InnoDB AUTO_INCREMENT = 1 DEFAULT CHARSET = utf8;

INSERT INTO `conditions` (`size`, `keyword`, `created`) VALUES (280, "고어", now());
INSERT INTO `conditions` (`size`, `keyword`, `created`) VALUES (285, "고어", now());
