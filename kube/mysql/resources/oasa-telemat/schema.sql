-- MySQL dump 10.13  Distrib 8.0.40, for Win64 (x86_64)
--
-- Host: localhost    Database: oasaTelemat
-- ------------------------------------------------------
-- Server version	8.0.40

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `line`
--
CREATE DATABASE IF NOT EXISTS `oasaTelemat`;
USE `oasaTelemat`;

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
) ENGINE=InnoDB AUTO_INCREMENT=456 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `route`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `route` (
  `id` int NOT NULL AUTO_INCREMENT,
  `route_code` int DEFAULT NULL,
  `ln_code` int NOT NULL,
  `route_descr` varchar(100) DEFAULT NULL,
  `route_descr_eng` varchar(100) DEFAULT NULL,
  `route_type` int DEFAULT NULL,
  `route_distance` decimal(7,2) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `route_code_un` (`route_code`),
  FOREIGN KEY (`ln_code`) REFERENCES `line` (`line_code`)
) ENGINE=InnoDB AUTO_INCREMENT=684 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `route01`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
-- Αυτός είναι ο πίνακας που αποθηκεύουμε τις λεπτομερειές για την κάθες διαδρομή
CREATE TABLE IF NOT EXISTS `route01` (
  `id` int NOT NULL AUTO_INCREMENT,
  `rt_code` int DEFAULT NULL,
  `routed_x` decimal(10,7) DEFAULT NULL,
  `routed_y` decimal(10,7) DEFAULT NULL,
  `routed_order` int NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `route01_code_un` (`rt_code`,`routed_order`),
  FOREIGN KEY (`rt_code`) REFERENCES `route` (`route_code`)
) ENGINE=InnoDB AUTO_INCREMENT=79464 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `route02`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
-- Αυτός είναι ο πίνακας που αποθηκεύουμε τις στάσεις ανα διαδρομή και σειρά Α/Α
CREATE TABLE IF NOT EXISTS `route02` (
  `rt_code` int NOT NULL,
  `stp_code` int NOT NULL,
  `senu` int NOT NULL,
  PRIMARY KEY (`rt_code`,`stp_code`,`senu`),
  FOREIGN KEY (`rt_code`) REFERENCES `route` (`route_code`),
  FOREIGN KEY (`stp_code`) REFERENCES `stop` (`stop_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `schedulemaster`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `schedulemaster` (
  `id` int NOT NULL AUTO_INCREMENT,
  `sdc_descr` varchar(50) DEFAULT NULL,
  `sdc_descr_eng` varchar(500) DEFAULT NULL,
  `sdc_code` int NOT NULL,
  `sdc_days` varchar(7) DEFAULT NULL,
  `sdc_months` varchar(12) DEFAULT NULL
  PRIMARY KEY (`id`),
  UNIQUE KEY `sdc_code_un` (`sdc_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `scheduletime`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `scheduletime` (
  `ln_code` int NOT NULL,
  `sdc_cd` int NOT NULL,
  `start_time` TIME NOT NULL,
  `end_time` TIME NOT NULL,
  `sort` int NOT NULL,
  `direction` int NOT NULL,
  PRIMARY KEY (`ln_code`,`sdc_cd`,`start_time`,`direction`),
  FOREIGN KEY (`sdc_cd`) REFERENCES `schedulemaster` (`sdc_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;


/*
ΑΥΤΟ ΘΑ ΠΡΕΠΕΙ ΝΑ ΦΎΓΕΙ ΔΕΝ ΧΡΕΙΑΖΕΤΑΙ 
  ΓΙΑΤΙ ΕΙΜΑΙ ΜΑΛΑΚΑΣ.
*/
/*CREATE TABLE IF NOT EXISTS `scheduleline` (
  `ln_code` int NOT NULL,
  `sdc_cd` int NOT NULL,
  PRIMARY KEY (`ln_code`,`sdc_cd`),
  FOREIGN KEY (`sdc_cd`) REFERENCES `schedulemaster` (`sdc_code`),
  FOREIGN KEY (`ln_code`) REFERENCES `line` (`line_code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
*/
--
-- Table structure for table `stop`
--

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
) ENGINE=InnoDB AUTO_INCREMENT=16105 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `syncversions`
--

/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE IF NOT EXISTS `syncversions` (
  `uv_descr` varchar(20) NOT NULL,
  `uv_lastupdatelong` int NOT NULL,
  PRIMARY KEY (`uv_descr`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;
