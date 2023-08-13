-- Adminer 4.8.1 MySQL 8.1.0 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;

SET NAMES utf8mb4;

CREATE DATABASE `haxball` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */ /*!80016 DEFAULT ENCRYPTION='N' */;
USE `haxball`;

CREATE TABLE `server` (
  `id` int NOT NULL AUTO_INCREMENT,
  `link` varchar(11) COLLATE utf8mb4_general_ci NOT NULL,
  `name` varchar(64) COLLATE utf8mb4_general_ci NOT NULL,
  `flag` varchar(2) COLLATE utf8mb4_general_ci NOT NULL,
  `private` bit(1) NOT NULL,
  `playersNow` int NOT NULL,
  `playersMax` int NOT NULL,
  `date` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- 2023-08-13 11:14:58
