-- Adminer 4.8.1 MySQL 5.5.5-10.11.2-MariaDB-1:10.11.2+maria~ubu2204 dump

SET NAMES utf8;
SET time_zone = '+00:00';
SET foreign_key_checks = 0;
SET sql_mode = 'NO_AUTO_VALUE_ON_ZERO';

SET NAMES utf8mb4;

CREATE DATABASE `haxball` /*!40100 DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci */;
USE `haxball`;

CREATE TABLE `server` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `link` varchar(11) NOT NULL,
  `name` varchar(64) NOT NULL,
  `flag` varchar(2) NOT NULL,
  `private` bit(1) NOT NULL,
  `playersNow` int(11) NOT NULL,
  `playersMax` int(11) NOT NULL,
  PRIMARY KEY (`id`),
  KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_general_ci;


-- 2023-08-09 05:30:54
