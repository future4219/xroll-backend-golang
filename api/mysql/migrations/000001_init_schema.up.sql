-- MySQL dump 10.13  Distrib 8.0.28, for Linux (x86_64)
--
-- Host: localhost    Database: e-privado
-- ------------------------------------------------------
-- Server version       8.0.28

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
-- Table structure for table `archive_categories`
--

DROP TABLE IF EXISTS `archive_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `archive_categories` (
  `archive_category_id` varchar(26) NOT NULL,
  `name` varchar(255) NOT NULL,
  `color` varchar(255) NOT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`archive_category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `archive_categories`
--

LOCK TABLES `archive_categories` WRITE;
/*!40000 ALTER TABLE `archive_categories` DISABLE KEYS */;
/*!40000 ALTER TABLE `archive_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `archive_document_files`
--

DROP TABLE IF EXISTS `archive_document_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `archive_document_files` (
  `archive_id` varchar(26) NOT NULL,
  `file_id` varchar(26) NOT NULL,
  PRIMARY KEY (`archive_id`,`file_id`),
  KEY `fk_archive_document_files_file` (`file_id`),
  CONSTRAINT `fk_archive_document_files_archive` FOREIGN KEY (`archive_id`) REFERENCES `archives` (`archive_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_archive_document_files_file` FOREIGN KEY (`file_id`) REFERENCES `files` (`file_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `archive_document_files`
--

LOCK TABLES `archive_document_files` WRITE;
/*!40000 ALTER TABLE `archive_document_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `archive_document_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `archive_files`
--

DROP TABLE IF EXISTS `archive_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `archive_files` (
  `archive_id` varchar(26) NOT NULL,
  `file_id` varchar(26) NOT NULL,
  PRIMARY KEY (`archive_id`,`file_id`),
  KEY `fk_archive_files_file` (`file_id`),
  CONSTRAINT `fk_archive_files_archive` FOREIGN KEY (`archive_id`) REFERENCES `archives` (`archive_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_archive_files_file` FOREIGN KEY (`file_id`) REFERENCES `files` (`file_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `archive_files`
--

LOCK TABLES `archive_files` WRITE;
/*!40000 ALTER TABLE `archive_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `archive_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `archive_tags`
--

DROP TABLE IF EXISTS `archive_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `archive_tags` (
  `archive_id` varchar(26) NOT NULL,
  `tag_id` varchar(26) NOT NULL,
  PRIMARY KEY (`archive_id`,`tag_id`),
  KEY `fk_archive_tags_tag` (`tag_id`),
  CONSTRAINT `fk_archive_tags_archive` FOREIGN KEY (`archive_id`) REFERENCES `archives` (`archive_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_archive_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`tag_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `archive_tags`
--

LOCK TABLES `archive_tags` WRITE;
/*!40000 ALTER TABLE `archive_tags` DISABLE KEYS */;
/*!40000 ALTER TABLE `archive_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `archive_to_archive_categories`
--

DROP TABLE IF EXISTS `archive_to_archive_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `archive_to_archive_categories` (
  `archive_id` varchar(26) NOT NULL,
  `archive_category_id` varchar(26) NOT NULL,
  PRIMARY KEY (`archive_id`,`archive_category_id`),
  KEY `fk_archive_to_archive_categories_archive_category` (`archive_category_id`),
  CONSTRAINT `fk_archive_to_archive_categories_archive` FOREIGN KEY (`archive_id`) REFERENCES `archives` (`archive_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_archive_to_archive_categories_archive_category` FOREIGN KEY (`archive_category_id`) REFERENCES `archive_categories` (`archive_category_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `archive_to_archive_categories`
--

LOCK TABLES `archive_to_archive_categories` WRITE;
/*!40000 ALTER TABLE `archive_to_archive_categories` DISABLE KEYS */;
/*!40000 ALTER TABLE `archive_to_archive_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `archives`
--

DROP TABLE IF EXISTS `archives`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `archives` (
  `archive_id` varchar(26) NOT NULL,
  `start_at` datetime(3) DEFAULT NULL,
  `end_at` datetime(3) DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  `message_json` longtext NOT NULL,
  `source` varchar(255) NOT NULL,
  `created_user_id` varchar(26) DEFAULT NULL,
  `updated_user_id` varchar(26) DEFAULT NULL,
  `thumbnail_file_id` varchar(26) DEFAULT NULL,
  `movie_file_id` varchar(26) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`archive_id`),
  KEY `fk_archives_created_user` (`created_user_id`),
  KEY `fk_archives_updated_user` (`updated_user_id`),
  KEY `fk_archives_thumbnail` (`thumbnail_file_id`),
  KEY `fk_archives_movie` (`movie_file_id`),
  CONSTRAINT `fk_archives_created_user` FOREIGN KEY (`created_user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_archives_movie` FOREIGN KEY (`movie_file_id`) REFERENCES `files` (`file_id`) ON DELETE SET NULL,
  CONSTRAINT `fk_archives_thumbnail` FOREIGN KEY (`thumbnail_file_id`) REFERENCES `files` (`file_id`) ON DELETE SET NULL,
  CONSTRAINT `fk_archives_updated_user` FOREIGN KEY (`updated_user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `archives`
--

LOCK TABLES `archives` WRITE;
/*!40000 ALTER TABLE `archives` DISABLE KEYS */;
/*!40000 ALTER TABLE `archives` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `clients`
--

DROP TABLE IF EXISTS `clients`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `clients` (
  `client_id` varchar(26) NOT NULL,
  `name` varchar(255) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`client_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `clients`
--

LOCK TABLES `clients` WRITE;
/*!40000 ALTER TABLE `clients` DISABLE KEYS */;
/*!40000 ALTER TABLE `clients` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `election_answer_question_choices`
--

DROP TABLE IF EXISTS `election_answer_question_choices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `election_answer_question_choices` (
  `election_answer_question_election_answer_question_id` varchar(26) NOT NULL,
  `election_question_choice_election_question_choice_id` varchar(26) NOT NULL,
  PRIMARY KEY (`election_answer_question_election_answer_question_id`,`election_question_choice_election_question_choice_id`),
  KEY `fk_election_answer_question_choices_election_question_choice` (`election_question_choice_election_question_choice_id`),
  CONSTRAINT `fk_election_answer_question_choices_election_answer_question` FOREIGN KEY (`election_answer_question_election_answer_question_id`) REFERENCES `election_answer_questions` (`election_answer_question_id`),
  CONSTRAINT `fk_election_answer_question_choices_election_question_choice` FOREIGN KEY (`election_question_choice_election_question_choice_id`) REFERENCES `election_question_choices` (`election_question_choice_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `election_answer_question_choices`
--

LOCK TABLES `election_answer_question_choices` WRITE;
/*!40000 ALTER TABLE `election_answer_question_choices` DISABLE KEYS */;
/*!40000 ALTER TABLE `election_answer_question_choices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `election_answer_questions`
--

DROP TABLE IF EXISTS `election_answer_questions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `election_answer_questions` (
  `election_answer_question_id` varchar(26) NOT NULL,
  `election_answer_id` varchar(26) DEFAULT NULL,
  `election_question_id` varchar(26) DEFAULT NULL,
  `election_question_text` longtext,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`election_answer_question_id`),
  KEY `fk_election_answers_election_answers` (`election_answer_id`),
  KEY `fk_election_answer_questions_election_question` (`election_question_id`),
  CONSTRAINT `fk_election_answer_questions_election_question` FOREIGN KEY (`election_question_id`) REFERENCES `election_questions` (`election_question_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_election_answers_election_answers` FOREIGN KEY (`election_answer_id`) REFERENCES `election_answers` (`election_answer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `election_answer_questions`
--

LOCK TABLES `election_answer_questions` WRITE;
/*!40000 ALTER TABLE `election_answer_questions` DISABLE KEYS */;
/*!40000 ALTER TABLE `election_answer_questions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `election_answers`
--

DROP TABLE IF EXISTS `election_answers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `election_answers` (
  `election_answer_id` varchar(26) NOT NULL,
  `election_id` varchar(26) DEFAULT NULL,
  `user_id` varchar(26) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`election_answer_id`),
  KEY `fk_election_answers_election` (`election_id`),
  KEY `fk_election_answers_user` (`user_id`),
  CONSTRAINT `fk_election_answers_election` FOREIGN KEY (`election_id`) REFERENCES `elections` (`election_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_election_answers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `election_answers`
--

LOCK TABLES `election_answers` WRITE;
/*!40000 ALTER TABLE `election_answers` DISABLE KEYS */;
/*!40000 ALTER TABLE `election_answers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `election_files`
--

DROP TABLE IF EXISTS `election_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `election_files` (
  `election_id` varchar(26) NOT NULL,
  `file_id` varchar(191) NOT NULL,
  PRIMARY KEY (`election_id`,`file_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `election_files`
--

LOCK TABLES `election_files` WRITE;
/*!40000 ALTER TABLE `election_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `election_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `election_question_choices`
--

DROP TABLE IF EXISTS `election_question_choices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `election_question_choices` (
  `election_question_choice_id` varchar(26) NOT NULL,
  `election_question_id` varchar(26) DEFAULT NULL,
  `value` varchar(255) NOT NULL,
  `sort_number` bigint DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`election_question_choice_id`),
  KEY `fk_election_questions_election_choices` (`election_question_id`),
  CONSTRAINT `fk_election_questions_election_choices` FOREIGN KEY (`election_question_id`) REFERENCES `election_questions` (`election_question_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `election_question_choices`
--

LOCK TABLES `election_question_choices` WRITE;
/*!40000 ALTER TABLE `election_question_choices` DISABLE KEYS */;
/*!40000 ALTER TABLE `election_question_choices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `election_questions`
--

DROP TABLE IF EXISTS `election_questions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `election_questions` (
  `election_question_id` varchar(26) NOT NULL,
  `election_id` varchar(26) DEFAULT NULL,
  `main_description` varchar(255) NOT NULL,
  `sub_description` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  `sort_number` bigint DEFAULT NULL,
  `max_choices` bigint DEFAULT NULL,
  `is_required` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`election_question_id`),
  KEY `fk_elections_election_questions` (`election_id`),
  CONSTRAINT `fk_elections_election_questions` FOREIGN KEY (`election_id`) REFERENCES `elections` (`election_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `election_questions`
--

LOCK TABLES `election_questions` WRITE;
/*!40000 ALTER TABLE `election_questions` DISABLE KEYS */;
/*!40000 ALTER TABLE `election_questions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `election_to_visibility_prefectures`
--

DROP TABLE IF EXISTS `election_to_visibility_prefectures`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `election_to_visibility_prefectures` (
  `election_id` varchar(26) NOT NULL,
  `prefecture_id` varchar(26) NOT NULL,
  PRIMARY KEY (`election_id`,`prefecture_id`),
  KEY `fk_election_to_visibility_prefectures_prefecture` (`prefecture_id`),
  CONSTRAINT `fk_election_to_visibility_prefectures_election` FOREIGN KEY (`election_id`) REFERENCES `elections` (`election_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_election_to_visibility_prefectures_prefecture` FOREIGN KEY (`prefecture_id`) REFERENCES `prefectures` (`prefecture_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `election_to_visibility_prefectures`
--

LOCK TABLES `election_to_visibility_prefectures` WRITE;
/*!40000 ALTER TABLE `election_to_visibility_prefectures` DISABLE KEYS */;
/*!40000 ALTER TABLE `election_to_visibility_prefectures` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `elections`
--

DROP TABLE IF EXISTS `elections`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `elections` (
  `election_id` varchar(26) NOT NULL,
  `title` varchar(255) NOT NULL,
  `message_json` longtext NOT NULL,
  `type` varchar(255) NOT NULL,
  `start_at` datetime(3) DEFAULT NULL,
  `end_at` datetime(3) DEFAULT NULL,
  `visibility` varchar(255) NOT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`election_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `elections`
--

LOCK TABLES `elections` WRITE;
/*!40000 ALTER TABLE `elections` DISABLE KEYS */;
/*!40000 ALTER TABLE `elections` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `files`
--

DROP TABLE IF EXISTS `files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `files` (
  `file_id` varchar(26) NOT NULL,
  `file_kind` varchar(255) NOT NULL,
  `file_name` varchar(1024) NOT NULL,
  `content_type` varchar(255) NOT NULL,
  `file_size` int NOT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `created_user_id` varchar(26) DEFAULT NULL,
  `updated_user_id` varchar(26) DEFAULT NULL,
  PRIMARY KEY (`file_id`),
  KEY `fk_files_created_user` (`created_user_id`),
  KEY `fk_files_updated_user` (`updated_user_id`),
  CONSTRAINT `fk_files_created_user` FOREIGN KEY (`created_user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_files_updated_user` FOREIGN KEY (`updated_user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `files`
--

LOCK TABLES `files` WRITE;
/*!40000 ALTER TABLE `files` DISABLE KEYS */;
/*!40000 ALTER TABLE `files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `information`
--

DROP TABLE IF EXISTS `information`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `information` (
  `information_id` varchar(26) NOT NULL,
  `start_at` datetime(3) DEFAULT NULL,
  `end_at` datetime(3) DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  `message_json` longtext NOT NULL,
  `created_user_id` varchar(26) DEFAULT NULL,
  `updated_user_id` varchar(26) DEFAULT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  PRIMARY KEY (`information_id`),
  KEY `fk_information_created_user` (`created_user_id`),
  KEY `fk_information_updated_user` (`updated_user_id`),
  CONSTRAINT `fk_information_created_user` FOREIGN KEY (`created_user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_information_updated_user` FOREIGN KEY (`updated_user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `information`
--

LOCK TABLES `information` WRITE;
/*!40000 ALTER TABLE `information` DISABLE KEYS */;
/*!40000 ALTER TABLE `information` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `information_categories`
--

DROP TABLE IF EXISTS `information_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `information_categories` (
  `information_category_id` varchar(26) NOT NULL,
  `name` varchar(255) NOT NULL,
  `color` varchar(255) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`information_category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `information_categories`
--

LOCK TABLES `information_categories` WRITE;
/*!40000 ALTER TABLE `information_categories` DISABLE KEYS */;
/*!40000 ALTER TABLE `information_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `information_clients`
--

DROP TABLE IF EXISTS `information_clients`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `information_clients` (
  `information_id` varchar(26) NOT NULL,
  `client_id` varchar(26) NOT NULL,
  PRIMARY KEY (`information_id`,`client_id`),
  KEY `fk_information_clients_client` (`client_id`),
  CONSTRAINT `fk_information_clients_client` FOREIGN KEY (`client_id`) REFERENCES `clients` (`client_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_information_clients_information` FOREIGN KEY (`information_id`) REFERENCES `information` (`information_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `information_clients`
--

LOCK TABLES `information_clients` WRITE;
/*!40000 ALTER TABLE `information_clients` DISABLE KEYS */;
/*!40000 ALTER TABLE `information_clients` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `information_files`
--

DROP TABLE IF EXISTS `information_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `information_files` (
  `information_id` varchar(26) NOT NULL,
  `file_id` varchar(26) NOT NULL,
  PRIMARY KEY (`information_id`,`file_id`),
  KEY `fk_information_files_file` (`file_id`),
  CONSTRAINT `fk_information_files_file` FOREIGN KEY (`file_id`) REFERENCES `files` (`file_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_information_files_information` FOREIGN KEY (`information_id`) REFERENCES `information` (`information_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `information_files`
--

LOCK TABLES `information_files` WRITE;
/*!40000 ALTER TABLE `information_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `information_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `information_tags`
--

DROP TABLE IF EXISTS `information_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `information_tags` (
  `information_id` varchar(26) NOT NULL,
  `tag_id` varchar(26) NOT NULL,
  PRIMARY KEY (`information_id`,`tag_id`),
  KEY `fk_information_tags_tag` (`tag_id`),
  CONSTRAINT `fk_information_tags_information` FOREIGN KEY (`information_id`) REFERENCES `information` (`information_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_information_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`tag_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `information_tags`
--

LOCK TABLES `information_tags` WRITE;
/*!40000 ALTER TABLE `information_tags` DISABLE KEYS */;
/*!40000 ALTER TABLE `information_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `information_to_information_categories`
--

DROP TABLE IF EXISTS `information_to_information_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `information_to_information_categories` (
  `information_id` varchar(26) NOT NULL,
  `information_category_id` varchar(26) NOT NULL,
  PRIMARY KEY (`information_id`,`information_category_id`),
  KEY `fk_information_to_information_categories_information_category` (`information_category_id`),
  CONSTRAINT `fk_information_to_information_categories_information` FOREIGN KEY (`information_id`) REFERENCES `information` (`information_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_information_to_information_categories_information_category` FOREIGN KEY (`information_category_id`) REFERENCES `information_categories` (`information_category_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `information_to_information_categories`
--

LOCK TABLES `information_to_information_categories` WRITE;
/*!40000 ALTER TABLE `information_to_information_categories` DISABLE KEYS */;
/*!40000 ALTER TABLE `information_to_information_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `member_occupations`
--

DROP TABLE IF EXISTS `member_occupations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `member_occupations` (
  `member_id` varchar(26) NOT NULL,
  `occupation_id` varchar(26) NOT NULL,
  PRIMARY KEY (`member_id`,`occupation_id`),
  KEY `fk_member_occupations_occupation` (`occupation_id`),
  CONSTRAINT `fk_member_occupations_member` FOREIGN KEY (`member_id`) REFERENCES `members` (`member_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_member_occupations_occupation` FOREIGN KEY (`occupation_id`) REFERENCES `occupations` (`occupation_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `member_occupations`
--

LOCK TABLES `member_occupations` WRITE;
/*!40000 ALTER TABLE `member_occupations` DISABLE KEYS */;
/*!40000 ALTER TABLE `member_occupations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `members`
--

DROP TABLE IF EXISTS `members`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `members` (
  `member_id` varchar(26) NOT NULL,
  `user_id` varchar(26) DEFAULT NULL,
  `first_name` varchar(255) NOT NULL,
  `last_name` varchar(255) NOT NULL,
  `first_name_kana` varchar(255) NOT NULL,
  `last_name_kana` varchar(255) NOT NULL,
  `gender` varchar(255) NOT NULL,
  `birthday` datetime(3) DEFAULT NULL,
  `license_number` varchar(255) NOT NULL,
  `external_organization_id` varchar(26) DEFAULT NULL,
  `send_to` varchar(255) NOT NULL,
  `date_agree` datetime(3) DEFAULT NULL,
  `post_code` varchar(255) NOT NULL,
  `prefecture_id` varchar(26) DEFAULT NULL,
  `city` varchar(255) NOT NULL,
  `address` varchar(255) NOT NULL,
  `building` varchar(255) NOT NULL,
  `phone` varchar(255) NOT NULL,
  `office_name` varchar(255) NOT NULL,
  `office_branch` varchar(255) NOT NULL,
  `office_department` varchar(255) NOT NULL,
  `position` varchar(255) NOT NULL,
  `office_post_code` varchar(255) NOT NULL,
  `office_prefecture_id` varchar(26) DEFAULT NULL,
  `office_city` varchar(255) NOT NULL,
  `office_address` varchar(255) NOT NULL,
  `office_building` varchar(255) NOT NULL,
  `office_phone` varchar(255) NOT NULL,
  `status` varchar(255) NOT NULL,
  `remarks` varchar(255) NOT NULL,
  `reason_for_withdrawal` varchar(255) NOT NULL,
  `withdrawal_date` datetime(3) DEFAULT NULL,
  `is_activated` tinyint(1) NOT NULL DEFAULT '0',
  `has_agreed_to_privacy_policy` tinyint(1) NOT NULL DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`member_id`),
  KEY `fk_members_external_organization` (`external_organization_id`),
  KEY `fk_members_prefecture` (`prefecture_id`),
  KEY `fk_members_office_prefecture` (`office_prefecture_id`),
  KEY `fk_users_member` (`user_id`),
  CONSTRAINT `fk_members_external_organization` FOREIGN KEY (`external_organization_id`) REFERENCES `prefectures` (`prefecture_id`),
  CONSTRAINT `fk_members_office_prefecture` FOREIGN KEY (`office_prefecture_id`) REFERENCES `prefectures` (`prefecture_id`),
  CONSTRAINT `fk_members_prefecture` FOREIGN KEY (`prefecture_id`) REFERENCES `prefectures` (`prefecture_id`),
  CONSTRAINT `fk_users_member` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `members`
--

LOCK TABLES `members` WRITE;
/*!40000 ALTER TABLE `members` DISABLE KEYS */;
/*!40000 ALTER TABLE `members` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `occupations`
--

DROP TABLE IF EXISTS `occupations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `occupations` (
  `occupation_id` varchar(26) NOT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`occupation_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `occupations`
--

LOCK TABLES `occupations` WRITE;
/*!40000 ALTER TABLE `occupations` DISABLE KEYS */;
/*!40000 ALTER TABLE `occupations` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_categories`
--

DROP TABLE IF EXISTS `package_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_categories` (
  `package_category_id` varchar(26) NOT NULL,
  `course_category` varchar(255) NOT NULL,
  `name` varchar(255) NOT NULL,
  `color` varchar(255) NOT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_category_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_categories`
--

LOCK TABLES `package_categories` WRITE;
/*!40000 ALTER TABLE `package_categories` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_completions`
--

DROP TABLE IF EXISTS `package_completions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_completions` (
  `package_completion_id` varchar(26) NOT NULL,
  `package_id` varchar(26) DEFAULT NULL,
  `package_plan_id` varchar(26) DEFAULT NULL,
  `user_id` varchar(26) DEFAULT NULL,
  `video_viewing_completed_at` datetime(3) DEFAULT NULL,
  `exam_passed_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_completion_id`),
  KEY `fk_package_completions_package` (`package_id`),
  KEY `fk_package_completions_package_plan` (`package_plan_id`),
  KEY `fk_package_completions_user` (`user_id`),
  CONSTRAINT `fk_package_completions_package` FOREIGN KEY (`package_id`) REFERENCES `packages` (`package_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_completions_package_plan` FOREIGN KEY (`package_plan_id`) REFERENCES `package_plans` (`package_plan_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_completions_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_completions`
--

LOCK TABLES `package_completions` WRITE;
/*!40000 ALTER TABLE `package_completions` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_completions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_content_to_files`
--

DROP TABLE IF EXISTS `package_content_to_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_content_to_files` (
  `package_content_id` varchar(26) NOT NULL,
  `file_id` varchar(26) NOT NULL,
  PRIMARY KEY (`package_content_id`,`file_id`),
  KEY `fk_package_content_to_files_file` (`file_id`),
  CONSTRAINT `fk_package_content_to_files_file` FOREIGN KEY (`file_id`) REFERENCES `files` (`file_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_content_to_files_package_content` FOREIGN KEY (`package_content_id`) REFERENCES `package_contents` (`package_content_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_content_to_files`
--

LOCK TABLES `package_content_to_files` WRITE;
/*!40000 ALTER TABLE `package_content_to_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_content_to_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_contents`
--

DROP TABLE IF EXISTS `package_contents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_contents` (
  `package_content_id` varchar(26) NOT NULL,
  `package_id` varchar(26) DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  `message_json` longtext NOT NULL,
  `sort_number` int NOT NULL,
  `type` varchar(255) NOT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_content_id`),
  KEY `fk_packages_package_contents` (`package_id`),
  CONSTRAINT `fk_packages_package_contents` FOREIGN KEY (`package_id`) REFERENCES `packages` (`package_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_contents`
--

LOCK TABLES `package_contents` WRITE;
/*!40000 ALTER TABLE `package_contents` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_contents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_document_to_files`
--

DROP TABLE IF EXISTS `package_document_to_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_document_to_files` (
  `package_video_document_content_id` varchar(26) NOT NULL,
  `file_id` varchar(26) NOT NULL,
  PRIMARY KEY (`package_video_document_content_id`,`file_id`),
  KEY `fk_package_document_to_files_file` (`file_id`),
  CONSTRAINT `fk_package_document_to_files_file` FOREIGN KEY (`file_id`) REFERENCES `files` (`file_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_document_to_files_package_video_document_content` FOREIGN KEY (`package_video_document_content_id`) REFERENCES `package_video_document_contents` (`package_video_document_content_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_document_to_files`
--

LOCK TABLES `package_document_to_files` WRITE;
/*!40000 ALTER TABLE `package_document_to_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_document_to_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_exam_answer_choices`
--

DROP TABLE IF EXISTS `package_exam_answer_choices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_exam_answer_choices` (
  `package_exam_content_answer_question_id` varchar(26) NOT NULL,
  `package_exam_content_choice_id` varchar(26) NOT NULL,
  PRIMARY KEY (`package_exam_content_answer_question_id`,`package_exam_content_choice_id`),
  KEY `fk_package_exam_answer_choices_package_exam_content_choice` (`package_exam_content_choice_id`),
  CONSTRAINT `fk_package_exam_answer_choices_package_exam_content_answc902e2ec` FOREIGN KEY (`package_exam_content_answer_question_id`) REFERENCES `package_exam_content_answer_questions` (`package_exam_content_answer_question_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_exam_answer_choices_package_exam_content_choice` FOREIGN KEY (`package_exam_content_choice_id`) REFERENCES `package_exam_content_choices` (`package_exam_content_choice_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_exam_answer_choices`
--

LOCK TABLES `package_exam_answer_choices` WRITE;
/*!40000 ALTER TABLE `package_exam_answer_choices` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_exam_answer_choices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_exam_content_answer_questions`
--

DROP TABLE IF EXISTS `package_exam_content_answer_questions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_exam_content_answer_questions` (
  `package_exam_content_answer_question_id` varchar(26) NOT NULL,
  `package_exam_content_answer_id` varchar(26) DEFAULT NULL,
  `package_exam_content_question_id` varchar(26) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_exam_content_answer_question_id`),
  KEY `fk_package_exam_content_answer_questions_package_exam_co52dd3d3e` (`package_exam_content_question_id`),
  KEY `fk_package_exam_content_answers_package_exam_content_ans0a959491` (`package_exam_content_answer_id`),
  CONSTRAINT `fk_package_exam_content_answer_questions_package_exam_co52dd3d3e` FOREIGN KEY (`package_exam_content_question_id`) REFERENCES `package_exam_content_questions` (`package_exam_content_question_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_exam_content_answers_package_exam_content_ans0a959491` FOREIGN KEY (`package_exam_content_answer_id`) REFERENCES `package_exam_content_answers` (`package_exam_content_answer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_exam_content_answer_questions`
--

LOCK TABLES `package_exam_content_answer_questions` WRITE;
/*!40000 ALTER TABLE `package_exam_content_answer_questions` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_exam_content_answer_questions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_exam_content_answers`
--

DROP TABLE IF EXISTS `package_exam_content_answers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_exam_content_answers` (
  `package_exam_content_answer_id` varchar(26) NOT NULL,
  `package_exam_content_id` varchar(26) DEFAULT NULL,
  `user_id` varchar(26) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_exam_content_answer_id`),
  KEY `fk_package_exam_content_answers_package_exam_content` (`package_exam_content_id`),
  KEY `fk_package_exam_content_answers_user` (`user_id`),
  CONSTRAINT `fk_package_exam_content_answers_package_exam_content` FOREIGN KEY (`package_exam_content_id`) REFERENCES `package_exam_contents` (`package_exam_content_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_exam_content_answers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_exam_content_answers`
--

LOCK TABLES `package_exam_content_answers` WRITE;
/*!40000 ALTER TABLE `package_exam_content_answers` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_exam_content_answers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_exam_content_choices`
--

DROP TABLE IF EXISTS `package_exam_content_choices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_exam_content_choices` (
  `package_exam_content_choice_id` varchar(26) NOT NULL,
  `package_exam_content_question_id` varchar(26) DEFAULT NULL,
  `value` varchar(255) NOT NULL,
  `is_correct` tinyint(1) NOT NULL DEFAULT '0',
  `sort_number` int DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_exam_content_choice_id`),
  KEY `fk_package_exam_content_questions_choices` (`package_exam_content_question_id`),
  CONSTRAINT `fk_package_exam_content_questions_choices` FOREIGN KEY (`package_exam_content_question_id`) REFERENCES `package_exam_content_questions` (`package_exam_content_question_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_exam_content_choices`
--

LOCK TABLES `package_exam_content_choices` WRITE;
/*!40000 ALTER TABLE `package_exam_content_choices` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_exam_content_choices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_exam_content_questions`
--

DROP TABLE IF EXISTS `package_exam_content_questions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_exam_content_questions` (
  `package_exam_content_question_id` varchar(26) NOT NULL,
  `package_exam_content_id` varchar(26) DEFAULT NULL,
  `main_description` varchar(255) NOT NULL,
  `sub_description` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  `sort_number` int DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_exam_content_question_id`),
  KEY `fk_package_exam_contents_questions` (`package_exam_content_id`),
  CONSTRAINT `fk_package_exam_contents_questions` FOREIGN KEY (`package_exam_content_id`) REFERENCES `package_exam_contents` (`package_exam_content_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_exam_content_questions`
--

LOCK TABLES `package_exam_content_questions` WRITE;
/*!40000 ALTER TABLE `package_exam_content_questions` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_exam_content_questions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_exam_contents`
--

DROP TABLE IF EXISTS `package_exam_contents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_exam_contents` (
  `package_exam_content_id` varchar(26) NOT NULL,
  `package_content_id` varchar(26) DEFAULT NULL,
  `passing_score` int NOT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_exam_content_id`),
  KEY `fk_package_contents_package_exam_content` (`package_content_id`),
  CONSTRAINT `fk_package_contents_package_exam_content` FOREIGN KEY (`package_content_id`) REFERENCES `package_contents` (`package_content_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_exam_contents`
--

LOCK TABLES `package_exam_contents` WRITE;
/*!40000 ALTER TABLE `package_exam_contents` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_exam_contents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_exam_to_completions`
--

DROP TABLE IF EXISTS `package_exam_to_completions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_exam_to_completions` (
  `package_completion_id` varchar(26) NOT NULL,
  `package_exam_content_id` varchar(26) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_completion_id`,`package_exam_content_id`),
  KEY `fk_package_exam_to_completions_package_exam_content` (`package_exam_content_id`),
  CONSTRAINT `fk_package_exam_to_completions_package_completion` FOREIGN KEY (`package_completion_id`) REFERENCES `package_completions` (`package_completion_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_exam_to_completions_package_exam_content` FOREIGN KEY (`package_exam_content_id`) REFERENCES `package_exam_contents` (`package_exam_content_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_exam_to_completions`
--

LOCK TABLES `package_exam_to_completions` WRITE;
/*!40000 ALTER TABLE `package_exam_to_completions` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_exam_to_completions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_files`
--

DROP TABLE IF EXISTS `package_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_files` (
  `package_id` varchar(26) NOT NULL,
  `file_id` varchar(26) NOT NULL,
  PRIMARY KEY (`package_id`,`file_id`),
  KEY `fk_package_files_file` (`file_id`),
  CONSTRAINT `fk_package_files_file` FOREIGN KEY (`file_id`) REFERENCES `files` (`file_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_files_package` FOREIGN KEY (`package_id`) REFERENCES `packages` (`package_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_files`
--

LOCK TABLES `package_files` WRITE;
/*!40000 ALTER TABLE `package_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_likes`
--

DROP TABLE IF EXISTS `package_likes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_likes` (
  `package_id` varchar(26) NOT NULL,
  `user_id` varchar(26) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_id`,`user_id`),
  KEY `fk_package_likes_user` (`user_id`),
  CONSTRAINT `fk_package_likes_package` FOREIGN KEY (`package_id`) REFERENCES `packages` (`package_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_likes_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_likes`
--

LOCK TABLES `package_likes` WRITE;
/*!40000 ALTER TABLE `package_likes` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_likes` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_plan_payments`
--

DROP TABLE IF EXISTS `package_plan_payments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_plan_payments` (
  `package_plan_payment_id` varchar(26) NOT NULL,
  `package_plan_id` varchar(26) DEFAULT NULL,
  `user_id` varchar(26) DEFAULT NULL,
  `session_id` longtext NOT NULL,
  `session_url` longtext NOT NULL,
  `status` longtext NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_plan_payment_id`),
  KEY `fk_package_plan_payments_package_plan` (`package_plan_id`),
  KEY `fk_package_plan_payments_user` (`user_id`),
  CONSTRAINT `fk_package_plan_payments_package_plan` FOREIGN KEY (`package_plan_id`) REFERENCES `package_plans` (`package_plan_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_plan_payments_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_plan_payments`
--

LOCK TABLES `package_plan_payments` WRITE;
/*!40000 ALTER TABLE `package_plan_payments` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_plan_payments` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_plan_to_users`
--

DROP TABLE IF EXISTS `package_plan_to_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_plan_to_users` (
  `package_plan_id` varchar(26) NOT NULL,
  `user_id` varchar(26) NOT NULL,
  PRIMARY KEY (`package_plan_id`,`user_id`),
  KEY `fk_package_plan_to_users_user` (`user_id`),
  CONSTRAINT `fk_package_plan_to_users_package_plan` FOREIGN KEY (`package_plan_id`) REFERENCES `package_plans` (`package_plan_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_plan_to_users_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_plan_to_users`
--

LOCK TABLES `package_plan_to_users` WRITE;
/*!40000 ALTER TABLE `package_plan_to_users` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_plan_to_users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_plans`
--

DROP TABLE IF EXISTS `package_plans`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_plans` (
  `package_plan_id` varchar(26) NOT NULL,
  `package_id` varchar(26) DEFAULT NULL,
  `name` varchar(255) NOT NULL,
  `price` int NOT NULL,
  `provider_price_id` longtext NOT NULL,
  `capacity` int NOT NULL,
  `is_member` tinyint(1) DEFAULT '0',
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_plan_id`),
  KEY `fk_packages_package_plans` (`package_id`),
  CONSTRAINT `fk_packages_package_plans` FOREIGN KEY (`package_id`) REFERENCES `packages` (`package_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_plans`
--

LOCK TABLES `package_plans` WRITE;
/*!40000 ALTER TABLE `package_plans` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_plans` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_questionnaire_answer_choices`
--

DROP TABLE IF EXISTS `package_questionnaire_answer_choices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_questionnaire_answer_choices` (
  `package_questionnaire_content_answer_question_id` varchar(26) NOT NULL,
  `package_questionnaire_choice_id` varchar(26) NOT NULL,
  PRIMARY KEY (`package_questionnaire_content_answer_question_id`,`package_questionnaire_choice_id`),
  KEY `fk_package_questionnaire_answer_choices_package_question43a6eef1` (`package_questionnaire_choice_id`),
  CONSTRAINT `fk_package_questionnaire_answer_choices_package_question43a6eef1` FOREIGN KEY (`package_questionnaire_choice_id`) REFERENCES `package_questionnaire_choices` (`package_questionnaire_choice_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_questionnaire_answer_choices_package_questiondeb9727a` FOREIGN KEY (`package_questionnaire_content_answer_question_id`) REFERENCES `package_questionnaire_content_answer_questions` (`package_questionnaire_content_answer_question_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_questionnaire_answer_choices`
--

LOCK TABLES `package_questionnaire_answer_choices` WRITE;
/*!40000 ALTER TABLE `package_questionnaire_answer_choices` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_questionnaire_answer_choices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_questionnaire_choices`
--

DROP TABLE IF EXISTS `package_questionnaire_choices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_questionnaire_choices` (
  `package_questionnaire_choice_id` varchar(26) NOT NULL,
  `package_questionnaire_content_id` varchar(26) DEFAULT NULL,
  `value` varchar(255) NOT NULL,
  `sort_number` int DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_questionnaire_choice_id`),
  KEY `fk_package_questionnaire_contents_package_questionnaire_choices` (`package_questionnaire_content_id`),
  CONSTRAINT `fk_package_questionnaire_contents_package_questionnaire_choices` FOREIGN KEY (`package_questionnaire_content_id`) REFERENCES `package_questionnaire_contents` (`package_questionnaire_content_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_questionnaire_choices`
--

LOCK TABLES `package_questionnaire_choices` WRITE;
/*!40000 ALTER TABLE `package_questionnaire_choices` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_questionnaire_choices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_questionnaire_content_answer_questions`
--

DROP TABLE IF EXISTS `package_questionnaire_content_answer_questions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_questionnaire_content_answer_questions` (
  `package_questionnaire_content_answer_question_id` varchar(26) NOT NULL,
  `package_questionnaire_content_answer_id` varchar(26) DEFAULT NULL,
  `package_questionnaire_content_id` varchar(26) DEFAULT NULL,
  `text_answer` varchar(255) NOT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_questionnaire_content_answer_question_id`),
  KEY `fk_package_questionnaire_content_answer_questions_packag4442001c` (`package_questionnaire_content_id`),
  KEY `fk_package_questionnaire_content_answers_package_questio0594d44c` (`package_questionnaire_content_answer_id`),
  CONSTRAINT `fk_package_questionnaire_content_answer_questions_packag4442001c` FOREIGN KEY (`package_questionnaire_content_id`) REFERENCES `package_questionnaire_contents` (`package_questionnaire_content_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_questionnaire_content_answers_package_questio0594d44c` FOREIGN KEY (`package_questionnaire_content_answer_id`) REFERENCES `package_questionnaire_content_answers` (`package_questionnaire_content_answer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_questionnaire_content_answer_questions`
--

LOCK TABLES `package_questionnaire_content_answer_questions` WRITE;
/*!40000 ALTER TABLE `package_questionnaire_content_answer_questions` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_questionnaire_content_answer_questions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_questionnaire_content_answers`
--

DROP TABLE IF EXISTS `package_questionnaire_content_answers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_questionnaire_content_answers` (
  `package_questionnaire_content_answer_id` varchar(26) NOT NULL,
  `package_content_id` varchar(26) DEFAULT NULL,
  `user_id` varchar(26) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_questionnaire_content_answer_id`),
  KEY `fk_package_questionnaire_content_answers_package_content` (`package_content_id`),
  KEY `fk_package_questionnaire_content_answers_user` (`user_id`),
  CONSTRAINT `fk_package_questionnaire_content_answers_package_content` FOREIGN KEY (`package_content_id`) REFERENCES `package_contents` (`package_content_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_questionnaire_content_answers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_questionnaire_content_answers`
--

LOCK TABLES `package_questionnaire_content_answers` WRITE;
/*!40000 ALTER TABLE `package_questionnaire_content_answers` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_questionnaire_content_answers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_questionnaire_contents`
--

DROP TABLE IF EXISTS `package_questionnaire_contents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_questionnaire_contents` (
  `package_questionnaire_content_id` varchar(26) NOT NULL,
  `package_content_id` varchar(26) DEFAULT NULL,
  `main_description` varchar(255) NOT NULL,
  `sub_description` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  `sort_number` int NOT NULL,
  `is_required` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_questionnaire_content_id`),
  KEY `fk_package_contents_package_questionnaire_contents` (`package_content_id`),
  CONSTRAINT `fk_package_contents_package_questionnaire_contents` FOREIGN KEY (`package_content_id`) REFERENCES `package_contents` (`package_content_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_questionnaire_contents`
--

LOCK TABLES `package_questionnaire_contents` WRITE;
/*!40000 ALTER TABLE `package_questionnaire_contents` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_questionnaire_contents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_tags`
--

DROP TABLE IF EXISTS `package_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_tags` (
  `package_id` varchar(26) NOT NULL,
  `tag_id` varchar(26) NOT NULL,
  PRIMARY KEY (`package_id`,`tag_id`),
  KEY `fk_package_tags_tag` (`tag_id`),
  CONSTRAINT `fk_package_tags_package` FOREIGN KEY (`package_id`) REFERENCES `packages` (`package_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`tag_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_tags`
--

LOCK TABLES `package_tags` WRITE;
/*!40000 ALTER TABLE `package_tags` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_to_categories`
--

DROP TABLE IF EXISTS `package_to_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_to_categories` (
  `package_id` varchar(26) NOT NULL,
  `package_category_id` varchar(26) NOT NULL,
  PRIMARY KEY (`package_id`,`package_category_id`),
  KEY `fk_package_to_categories_package_category` (`package_category_id`),
  CONSTRAINT `fk_package_to_categories_package` FOREIGN KEY (`package_id`) REFERENCES `packages` (`package_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_to_categories_package_category` FOREIGN KEY (`package_category_id`) REFERENCES `package_categories` (`package_category_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_to_categories`
--

LOCK TABLES `package_to_categories` WRITE;
/*!40000 ALTER TABLE `package_to_categories` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_to_categories` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_video_document_contents`
--

DROP TABLE IF EXISTS `package_video_document_contents`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_video_document_contents` (
  `package_video_document_content_id` varchar(26) NOT NULL,
  `package_content_id` varchar(26) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_video_document_content_id`),
  KEY `fk_package_contents_package_video_document_content` (`package_content_id`),
  CONSTRAINT `fk_package_contents_package_video_document_content` FOREIGN KEY (`package_content_id`) REFERENCES `package_contents` (`package_content_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_video_document_contents`
--

LOCK TABLES `package_video_document_contents` WRITE;
/*!40000 ALTER TABLE `package_video_document_contents` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_video_document_contents` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_video_to_files`
--

DROP TABLE IF EXISTS `package_video_to_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_video_to_files` (
  `package_video_document_content_id` varchar(26) NOT NULL,
  `file_id` varchar(26) NOT NULL,
  PRIMARY KEY (`package_video_document_content_id`,`file_id`),
  KEY `fk_package_video_to_files_file` (`file_id`),
  CONSTRAINT `fk_package_video_to_files_file` FOREIGN KEY (`file_id`) REFERENCES `files` (`file_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_video_to_files_package_video_document_content` FOREIGN KEY (`package_video_document_content_id`) REFERENCES `package_video_document_contents` (`package_video_document_content_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_video_to_files`
--

LOCK TABLES `package_video_to_files` WRITE;
/*!40000 ALTER TABLE `package_video_to_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_video_to_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `package_video_viewed_users`
--

DROP TABLE IF EXISTS `package_video_viewed_users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `package_video_viewed_users` (
  `package_completion_id` varchar(26) NOT NULL,
  `file_id` varchar(26) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_completion_id`,`file_id`),
  KEY `fk_package_video_viewed_users_file` (`file_id`),
  CONSTRAINT `fk_package_video_viewed_users_file` FOREIGN KEY (`file_id`) REFERENCES `files` (`file_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_package_video_viewed_users_package_completion` FOREIGN KEY (`package_completion_id`) REFERENCES `package_completions` (`package_completion_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `package_video_viewed_users`
--

LOCK TABLES `package_video_viewed_users` WRITE;
/*!40000 ALTER TABLE `package_video_viewed_users` DISABLE KEYS */;
/*!40000 ALTER TABLE `package_video_viewed_users` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `packages`
--

DROP TABLE IF EXISTS `packages`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `packages` (
  `package_id` varchar(26) NOT NULL,
  `title` varchar(255) NOT NULL,
  `thumbnail_id` varchar(26) DEFAULT NULL,
  `difficulty` varchar(255) NOT NULL,
  `start_at` datetime(3) DEFAULT NULL,
  `end_at` datetime(3) DEFAULT NULL,
  `start_duration_at` datetime(3) DEFAULT NULL,
  `end_duration_at` datetime(3) DEFAULT NULL,
  `message_json` longtext NOT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`package_id`),
  KEY `fk_packages_thumbnail` (`thumbnail_id`),
  CONSTRAINT `fk_packages_thumbnail` FOREIGN KEY (`thumbnail_id`) REFERENCES `files` (`file_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `packages`
--

LOCK TABLES `packages` WRITE;
/*!40000 ALTER TABLE `packages` DISABLE KEYS */;
/*!40000 ALTER TABLE `packages` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `prefectures`
--

DROP TABLE IF EXISTS `prefectures`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `prefectures` (
  `prefecture_id` varchar(26) NOT NULL,
  `name` varchar(255) NOT NULL,
  PRIMARY KEY (`prefecture_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `prefectures`
--

LOCK TABLES `prefectures` WRITE;
/*!40000 ALTER TABLE `prefectures` DISABLE KEYS */;
/*!40000 ALTER TABLE `prefectures` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `questionnaire_answer_details`
--

DROP TABLE IF EXISTS `questionnaire_answer_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `questionnaire_answer_details` (
  `questionnaire_answer_detail_id` varchar(26) NOT NULL,
  `questionnaire_answer_id` varchar(26) DEFAULT NULL,
  `questionnaire_question_id` varchar(26) DEFAULT NULL,
  `text_answer` varchar(255) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`questionnaire_answer_detail_id`),
  KEY `fk_questionnaire_answers_questionnaire_answer_details` (`questionnaire_answer_id`),
  KEY `fk_questionnaire_answer_details_questionnaire_question` (`questionnaire_question_id`),
  CONSTRAINT `fk_questionnaire_answer_details_questionnaire_question` FOREIGN KEY (`questionnaire_question_id`) REFERENCES `questionnaire_questions` (`questionnaire_question_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_questionnaire_answers_questionnaire_answer_details` FOREIGN KEY (`questionnaire_answer_id`) REFERENCES `questionnaire_answers` (`questionnaire_answer_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `questionnaire_answer_details`
--

LOCK TABLES `questionnaire_answer_details` WRITE;
/*!40000 ALTER TABLE `questionnaire_answer_details` DISABLE KEYS */;
/*!40000 ALTER TABLE `questionnaire_answer_details` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `questionnaire_answers`
--

DROP TABLE IF EXISTS `questionnaire_answers`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `questionnaire_answers` (
  `questionnaire_answer_id` varchar(26) NOT NULL,
  `questionnaire_id` varchar(26) DEFAULT NULL,
  `user_id` varchar(26) DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`questionnaire_answer_id`),
  KEY `fk_questionnaire_answers_user` (`user_id`),
  KEY `fk_questionnaire_answers_questionnaire` (`questionnaire_id`),
  CONSTRAINT `fk_questionnaire_answers_questionnaire` FOREIGN KEY (`questionnaire_id`) REFERENCES `questionnaires` (`questionnaire_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_questionnaire_answers_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `questionnaire_answers`
--

LOCK TABLES `questionnaire_answers` WRITE;
/*!40000 ALTER TABLE `questionnaire_answers` DISABLE KEYS */;
/*!40000 ALTER TABLE `questionnaire_answers` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `questionnaire_files`
--

DROP TABLE IF EXISTS `questionnaire_files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `questionnaire_files` (
  `questionnaire_id` varchar(26) NOT NULL,
  `file_id` varchar(26) NOT NULL,
  PRIMARY KEY (`questionnaire_id`,`file_id`),
  KEY `fk_questionnaire_files_file` (`file_id`),
  CONSTRAINT `fk_questionnaire_files_file` FOREIGN KEY (`file_id`) REFERENCES `files` (`file_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_questionnaire_files_questionnaire` FOREIGN KEY (`questionnaire_id`) REFERENCES `questionnaires` (`questionnaire_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `questionnaire_files`
--

LOCK TABLES `questionnaire_files` WRITE;
/*!40000 ALTER TABLE `questionnaire_files` DISABLE KEYS */;
/*!40000 ALTER TABLE `questionnaire_files` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `questionnaire_question_choices`
--

DROP TABLE IF EXISTS `questionnaire_question_choices`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `questionnaire_question_choices` (
  `questionnaire_question_choice_id` varchar(26) NOT NULL,
  `questionnaire_question_id` varchar(26) DEFAULT NULL,
  `value` varchar(255) NOT NULL,
  `sort_number` bigint DEFAULT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`questionnaire_question_choice_id`),
  KEY `fk_questionnaire_questions_questionnaire_question_choices` (`questionnaire_question_id`),
  CONSTRAINT `fk_questionnaire_questions_questionnaire_question_choices` FOREIGN KEY (`questionnaire_question_id`) REFERENCES `questionnaire_questions` (`questionnaire_question_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `questionnaire_question_choices`
--

LOCK TABLES `questionnaire_question_choices` WRITE;
/*!40000 ALTER TABLE `questionnaire_question_choices` DISABLE KEYS */;
/*!40000 ALTER TABLE `questionnaire_question_choices` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `questionnaire_question_choices_answer`
--

DROP TABLE IF EXISTS `questionnaire_question_choices_answer`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `questionnaire_question_choices_answer` (
  `questionnaire_answer_detail_id` varchar(26) NOT NULL,
  `questionnaire_question_choice_id` varchar(26) NOT NULL,
  PRIMARY KEY (`questionnaire_answer_detail_id`,`questionnaire_question_choice_id`),
  KEY `fk_questionnaire_question_choices_answer_questionnaire_qe269624b` (`questionnaire_question_choice_id`),
  CONSTRAINT `fk_questionnaire_question_choices_answer_questionnaire_a426143ee` FOREIGN KEY (`questionnaire_answer_detail_id`) REFERENCES `questionnaire_answer_details` (`questionnaire_answer_detail_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_questionnaire_question_choices_answer_questionnaire_qe269624b` FOREIGN KEY (`questionnaire_question_choice_id`) REFERENCES `questionnaire_question_choices` (`questionnaire_question_choice_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `questionnaire_question_choices_answer`
--

LOCK TABLES `questionnaire_question_choices_answer` WRITE;
/*!40000 ALTER TABLE `questionnaire_question_choices_answer` DISABLE KEYS */;
/*!40000 ALTER TABLE `questionnaire_question_choices_answer` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `questionnaire_questions`
--

DROP TABLE IF EXISTS `questionnaire_questions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `questionnaire_questions` (
  `questionnaire_question_id` varchar(26) NOT NULL,
  `questionnaire_id` varchar(26) DEFAULT NULL,
  `main_description` varchar(255) NOT NULL,
  `sub_description` varchar(255) NOT NULL,
  `type` varchar(255) NOT NULL,
  `sort_number` bigint DEFAULT NULL,
  `is_required` tinyint(1) NOT NULL DEFAULT '0',
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`questionnaire_question_id`),
  KEY `fk_questionnaires_questionnaire_questions` (`questionnaire_id`),
  CONSTRAINT `fk_questionnaires_questionnaire_questions` FOREIGN KEY (`questionnaire_id`) REFERENCES `questionnaires` (`questionnaire_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `questionnaire_questions`
--

LOCK TABLES `questionnaire_questions` WRITE;
/*!40000 ALTER TABLE `questionnaire_questions` DISABLE KEYS */;
/*!40000 ALTER TABLE `questionnaire_questions` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `questionnaire_tags`
--

DROP TABLE IF EXISTS `questionnaire_tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `questionnaire_tags` (
  `questionnaire_id` varchar(26) NOT NULL,
  `tag_id` varchar(26) NOT NULL,
  PRIMARY KEY (`questionnaire_id`,`tag_id`),
  KEY `fk_questionnaire_tags_tag` (`tag_id`),
  CONSTRAINT `fk_questionnaire_tags_questionnaire` FOREIGN KEY (`questionnaire_id`) REFERENCES `questionnaires` (`questionnaire_id`) ON DELETE CASCADE,
  CONSTRAINT `fk_questionnaire_tags_tag` FOREIGN KEY (`tag_id`) REFERENCES `tags` (`tag_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `questionnaire_tags`
--

LOCK TABLES `questionnaire_tags` WRITE;
/*!40000 ALTER TABLE `questionnaire_tags` DISABLE KEYS */;
/*!40000 ALTER TABLE `questionnaire_tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `questionnaires`
--

DROP TABLE IF EXISTS `questionnaires`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `questionnaires` (
  `questionnaire_id` varchar(26) NOT NULL,
  `title` varchar(255) NOT NULL,
  `message_json` longtext NOT NULL,
  `start_at` datetime(3) DEFAULT NULL,
  `end_at` datetime(3) DEFAULT NULL,
  `visibility` varchar(255) NOT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`questionnaire_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `questionnaires`
--

LOCK TABLES `questionnaires` WRITE;
/*!40000 ALTER TABLE `questionnaires` DISABLE KEYS */;
/*!40000 ALTER TABLE `questionnaires` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `tags`
--

DROP TABLE IF EXISTS `tags`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `tags` (
  `tag_id` varchar(26) NOT NULL,
  `name` varchar(255) NOT NULL,
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`tag_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `tags`
--

LOCK TABLES `tags` WRITE;
/*!40000 ALTER TABLE `tags` DISABLE KEYS */;
/*!40000 ALTER TABLE `tags` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `user_id` varchar(26) NOT NULL,
  `login_id` varchar(191) NOT NULL,
  `email` varchar(255) NOT NULL,
  `email_to_update` varchar(255) DEFAULT NULL,
  `hashed_password` varchar(255) NOT NULL,
  `user_type` varchar(255) NOT NULL,
  `is_deleted` tinyint(1) DEFAULT '0',
  `created_at` datetime(3) DEFAULT NULL,
  `updated_at` datetime(3) DEFAULT NULL,
  PRIMARY KEY (`user_id`),
  UNIQUE KEY `login_id` (`login_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `users`
--

LOCK TABLES `users` WRITE;
/*!40000 ALTER TABLE `users` DISABLE KEYS */;
/*!40000 ALTER TABLE `users` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2024-03-12 12:20:39
