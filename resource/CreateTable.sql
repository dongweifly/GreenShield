CREATE TABLE `operate_record`
(
    `id`          bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `text`        varchar(128) DEFAULT NULL COMMENT '敏感词，最大128',
    `bizType`     varchar(32)  DEFAULT NULL COMMENT '业务类型',
    `sceneType`   varchar(32) NOT NULL COMMENT '场景编码，porn, politics, ad',
    `operator`    varchar(255) DEFAULT NULL COMMENT '扩展',
    `operateType` varchar(32) NOT NULL COMMENT '操作类型, remove, add',
    `createTime`  bigint(20) DEFAULT NULL COMMENT '创建时间戳',
    PRIMARY KEY (`id`),
    KEY           `idx_operator` (`operator`),
    KEY           `idx_text` (`text`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `words`
(
    `id`         bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `text`       varchar(128) DEFAULT NULL COMMENT '敏感词，最大128',
    `bizType`    varchar(32)  DEFAULT NULL COMMENT '业务类型',
    `sceneType`  varchar(32) NOT NULL COMMENT '场景编码，porn, politics, ad',
    `extent`     varchar(255) DEFAULT NULL COMMENT '扩展',
    `valid`      tinyint(1) DEFAULT '1' COMMENT '是否有效，true:有效， false:无效',
    `updateTime` bigint(20) DEFAULT NULL COMMENT '创建时间戳',
    `createTime` bigint(20) DEFAULT NULL COMMENT '创建时间戳',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_bizType_text` (`bizType`,`text`),
    KEY          `idx_bizType_sceneType` (`bizType`,`sceneType`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;

CREATE TABLE `words_info`
(
    `id`         bigint(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `name`       varchar(64) NOT NULL COMMENT '词库名称',
    `bizType`    varchar(32) NOT NULL COMMENT '业务类型',
    `sceneType`  varchar(32) NOT NULL COMMENT '场景编码，porn, politics, ad',
    `parentName` varchar(64) DEFAULT NULL,
    `updateTime` bigint(20) DEFAULT NULL COMMENT '更新时间戳',
    `createTime` bigint(20) DEFAULT NULL COMMENT '创建时间戳',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_biztype_scenetype` (`bizType`,`sceneType`),
    UNIQUE KEY `idx_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4;