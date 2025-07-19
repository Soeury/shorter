create table `reflect_map` (
    `id` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'primary key',
    `create_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT 'create time',
    `create_by` VARCHAR(64) NOT NULL DEFAULT "rabbit" COMMENT 'create by',
    `is_del` TINYINT UNSIGNED NOT NULL DEFAULT '0' COMMENT 'soft delete or not ? 0-not delete 1-delete',
    
    `lurl` VARCHAR(2048) DEFAULT NULL COMMENT 'long url',
    `md5` CHAR(32) DEFAULT NULL COMMENT 'long url md5',
    `surl` VARCHAR(11) DEFAULT NULL COMMENT 'short url',
	 `expire_at` TIMESTAMP DEFAULT NULL COMMENT 'expire_time',
    
    PRIMARY KEY (`id`),
    INDEX(`is_del`),
    UNIQUE(`md5`),
    UNIQUE(`surl`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT = 'long-short url reflect map'; 