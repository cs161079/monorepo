-- DELETE DATABASE IF EXISTS oasadb;
-- Create the database if it doesn't exist and not delete it if it exists and recreate it. Likewise with tables
-- CREATE DATABASE IF NOT EXISTS oasa-telemat;
USE oasa-telemat;

--
-- DROP AND CREATE LINE TABLE
--
CREATE TABLE IF NOT EXISTS `line` (
  `id` int NOT NULL AUTO_INCREMENT,
  `ml_code` int DEFAULT NULL,
  `sdc_code` int DEFAULT NULL,
  `line_code` int NOT NULL,
  `line_id` varchar(45) NOT NULL,
  `line_descr` varchar(100) DEFAULT NULL,
  `line_descr_eng` varchar(100) DEFAULT NULL,
  `mld_master` smallint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `LINE_CODE` (`line_code`)
);

--
-- DROP AND CREATE ROUTE TABLE
--
CREATE TABLE IF NOT EXISTS `route` (
  `id` int NOT NULL AUTO_INCREMENT,
  `route_code` int DEFAULT NULL,
  `line_code` int NOT NULL,
  `route_descr` varchar(100) DEFAULT NULL,
  `route_descr_eng` varchar(100) DEFAULT NULL,
  `route_type` int DEFAULT NULL,
  `route_distance` decimal(7,2) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `route_code_un` (`route_code`),
  KEY `line_code_indx` (`line_code`)
);

--
-- DROP AND CREATE STOP TABLE
--
CREATE TABLE IF NOT EXISTS `stop` (
  `id` int NOT NULL AUTO_INCREMENT,
  `stop_code` int NOT NULL,
  `stop_id` varchar(6) DEFAULT NULL,
  `stop_descr` varchar(100) DEFAULT NULL,
  `stop_descr_eng` varchar(100) DEFAULT NULL,
  `stop_street` varchar(100) DEFAULT NULL,
  `stop_street_eng` varchar(100) DEFAULT NULL,
  `stop_heading` int DEFAULT NULL,
  `stop_lat` decimal(10,7) DEFAULT NULL,
  `stop_lng` decimal(10,7) DEFAULT NULL,
  `stop_type` int DEFAULT NULL,
  `stop_amea` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `stop_code_un` (`stop_code`)
);

--
-- DROP AND CREATE ROUTESTOPS TABLE
--
CREATE TABLE IF NOT EXISTS `routestops` (
  `route_code` int NOT NULL,
  `stop_code` int NOT NULL,
  `senu` int NOT NULL,
  PRIMARY KEY (`route_code`,`stop_code`,`senu`)
);

CREATE TABLE IF NOT EXISTS `route01` (
  `id` int NOT NULL AUTO_INCREMENT,
  `route_code` int DEFAULT NULL,
  `routed_x` decimal(10,7) DEFAULT NULL,
  `routed_y` decimal(10,7) DEFAULT NULL,
  `routed_order` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `route01_code_un` (`route_code`, `routed_order`)
);

--
-- DROP AND CREATE SCHEDULEMASTER TABLE
--
CREATE TABLE IF NOT EXISTS `schedulemaster` (
  `id` int NOT NULL AUTO_INCREMENT,
  `sdc_descr` varchar(50) DEFAULT NULL,
  `sdc_descr_eng` varchar(500) DEFAULT NULL,
  `sdc_code` int NOT NULL,
  `line_code` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `LINE_SDC_CODE` (`line_code`,`sdc_code`)
);

--
-- DROP AND CREATE SCHEDULETIME TABLE
--
CREATE TABLE IF NOT EXISTS `scheduletime`(
  `line_code` int NOT NULL,
  `sdc_code` int NOT NULL,
  `start_time` varchar(5) NOT NULL,
  `type` int(1) NOT NULL,
  PRIMARY KEY (`line_code`, `sdc_code`, `start_time`, `type`)
);

CREATE TABLE IF NOT EXISTS `syncversions`(
    `uv_descr` varchar(20) NOT NULL,
    `uv_lastupdatelong` int NOT NULL,
    Primary Key(`uv_descr`)
);
