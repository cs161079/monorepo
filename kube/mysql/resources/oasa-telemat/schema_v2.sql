-- MySQL dump 10.13  Distrib 8.0.40, for Win64 (x86_64)
--
-- Host: localhost    Database: oasatelemat

--
-- Table structure for table `line`
--

CREATE DATABASE IF NOT EXISTS `oasa-telemat`;
USE `oasa-telemat`;

-- DROP TABLE IF EXISTS `line`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `line` (
  `id` int NOT NULL AUTO_INCREMENT,
  `ml_code` int DEFAULT NULL,
  `sdc_cd` int DEFAULT NULL,
  `line_code` int NOT NULL,
  `line_id` varchar(45) NOT NULL,
  `line_descr` varchar(100) DEFAULT NULL,
  `line_descr_eng` varchar(100) DEFAULT NULL,
  `mld_master` smallint DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `LINE_CODE` (`line_code`)
) ENGINE=InnoDB AUTO_INCREMENT=12637 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `route`
--

-- DROP TABLE IF EXISTS `route`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `route` (
  `id` int NOT NULL AUTO_INCREMENT,
  `route_code` int NOT NULL,
  `ln_code` int NOT NULL,
  `route_descr` varchar(150) DEFAULT NULL,
  `route_descr_eng` varchar(150) DEFAULT NULL,
  `route_type` int DEFAULT NULL,
  `route_distance` decimal(7,2) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `route_code_un` (`route_code`),
  KEY `ln_code` (`ln_code`),
  CONSTRAINT `route_ibfk_1` FOREIGN KEY (`ln_code`) REFERENCES `line` (`line_code`)
) ENGINE=InnoDB AUTO_INCREMENT=18320 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;\

--
-- Table structure for table `stop`
--

-- DROP TABLE IF EXISTS `stop`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `stop` (
  `id` int NOT NULL AUTO_INCREMENT,
  `stop_code` int NOT NULL,
  `stop_id` varchar(10) DEFAULT NULL,
  `stop_descr` varchar(100) DEFAULT NULL,
  `stop_descr_eng` varchar(100) DEFAULT NULL,
  `stop_street` varchar(100) DEFAULT NULL,
  `stop_street_eng` varchar(100) DEFAULT NULL,
  `stop_heading` int DEFAULT NULL,
  `stop_lat` decimal(10,7) DEFAULT NULL,
  `stop_lng` decimal(10,7) DEFAULT NULL,
  `stop_type` int DEFAULT NULL,
  `stop_amea` int DEFAULT NULL,
  `destinations` varchar(1000) DEFAULT NULL,
  `destinations_eng` varchar(1000) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `stop_code_un` (`stop_code`)
) ENGINE=InnoDB AUTO_INCREMENT=207552 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `route01`
--

-- DROP TABLE IF EXISTS `route01`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `route01` (
  `id` int NOT NULL AUTO_INCREMENT,
  `rt_code` int DEFAULT NULL,
  `routed_x` decimal(10,7) DEFAULT NULL,
  `routed_y` decimal(10,7) DEFAULT NULL,
  `routed_order` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `route01_code_un` (`rt_code`,`routed_order`),
  CONSTRAINT `route01_ibfk_1` FOREIGN KEY (`rt_code`) REFERENCES `route` (`route_code`)
) ENGINE=InnoDB AUTO_INCREMENT=1738840 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `route02`
--

-- DROP TABLE IF EXISTS `route02`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `route02` (
  `rt_code` int NOT NULL,
  `stp_code` int NOT NULL,
  `senu` int NOT NULL,
  PRIMARY KEY (`rt_code`,`stp_code`,`senu`),
  KEY `stp_code` (`stp_code`),
  CONSTRAINT `route02_ibfk_1` FOREIGN KEY (`rt_code`) REFERENCES `route` (`route_code`),
  CONSTRAINT `route02_ibfk_2` FOREIGN KEY (`stp_code`) REFERENCES `stop` (`stop_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `schedulemaster`
--

-- DROP TABLE IF EXISTS `schedulemaster`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `schedulemaster` (
  `id` int NOT NULL AUTO_INCREMENT,
  `sdc_descr` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci DEFAULT NULL,
  `sdc_descr_eng` varchar(500) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `sdc_code` int NOT NULL,
  `sdc_months` varchar(12) COLLATE utf8mb4_general_ci DEFAULT NULL,
  `sdc_days` varchar(7) COLLATE utf8mb4_general_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `sdc_code_un` (`sdc_code`)
) ENGINE=InnoDB AUTO_INCREMENT=2350 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `scheduletime`
--

-- DROP TABLE IF EXISTS `scheduletime`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `scheduletime` (
  `ln_code` int NOT NULL,
  `sdc_cd` int NOT NULL,
  `start_time` time NOT NULL,
  `end_time` time NOT NULL,
  `sort` int NOT NULL,
  `direction` int NOT NULL,
  PRIMARY KEY (`ln_code`,`sdc_cd`,`start_time`,`direction`),
  KEY `sdc_cd` (`sdc_cd`),
  CONSTRAINT `scheduletime_ibfk_1` FOREIGN KEY (`sdc_cd`) REFERENCES `schedulemaster` (`sdc_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `syncversions`
--

-- DROP TABLE IF EXISTS `syncversions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `syncversions` (
  `uv_descr` varchar(20) NOT NULL,
  `uv_lastupdatelong` int NOT NULL,
  PRIMARY KEY (`uv_descr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

