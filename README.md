# chiyou-code
code gen (默认 Spring 3.x Java 17)

**Spring 3.x (Java 17)**

> lombok

> springdoc-openapi-starter-webmvc-ui:2.3.0

> springboot:3.2.2


#### 安装教程

1.  初始化步骤

确保安装`cobra-cli`，通过`cobra-cli`进行命令行程序开发

```shell
go install github.com/spf13/cobra-cli@latest

go mod init chiyou.code/mmc

cobra-cli init --author "FF911" --license mit --viper

cobra-cli add java

cobra-cli add csharp
```

#### 表规范 (mysql 8.0)

```
-- sys_platform definition

CREATE TABLE `sys_platform` (
  `TransactionNumber` int NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `PlatformCode` varchar(50) NOT NULL DEFAULT '' COMMENT '平台编码',
  `PlatformName` varchar(128) NOT NULL DEFAULT '' COMMENT '平台名称',
  `PlatformUrl` varchar(255) NOT NULL COMMENT '平台URL',
  `PlatformIcon` varchar(50) NOT NULL DEFAULT '' COMMENT '平台图标',
  `Status` tinyint(3) NOT NULL DEFAULT '1' COMMENT '状态: 1-启用; 0-停用',
  `Description` varchar(255) NOT NULL DEFAULT '' COMMENT '备注',
  `Version` int NOT NULL DEFAULT '0' COMMENT '版本',
  `InUser` varchar(128) NOT NULL COMMENT '创建人',
  `InDate` bigint NOT NULL COMMENT '创建时间',
  `LastEditUser` varchar(128) DEFAULT NULL COMMENT '最后更新人',
  `LastEditDate` bigint DEFAULT NULL COMMENT '最后更新时间',
  `Deleted` tinyint(1) NOT NULL DEFAULT '0' COMMENT '是否删除: 1-已删除; 0-正常',
  PRIMARY KEY (`TransactionNumber`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci COMMENT='平台表';
```

#### 使用说明

1.  `mmc java` or `mmc_amd64.exe java --config .mmc.yml`