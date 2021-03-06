-- data stores generic data rows
CREATE TABLE IF NOT EXISTS `data` (
  `id` BIGINT(20) UNSIGNED NOT NULL AUTO_INCREMENT,
  `source` VARCHAR(255) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `type` VARCHAR(255) NOT NULL,
  `value` VARCHAR(255) NOT NULL,
  `created_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `source` (`source`),
  KEY `type` (`type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_520_ci ROW_FORMAT=dynamic;
