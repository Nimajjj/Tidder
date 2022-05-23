-- MySQL Workbench Forward Engineering

SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0;
SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0;
SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='ONLY_FULL_GROUP_BY,STRICT_TRANS_TABLES,NO_ZERO_IN_DATE,NO_ZERO_DATE,ERROR_FOR_DIVISION_BY_ZERO,NO_ENGINE_SUBSTITUTION';

-- -----------------------------------------------------
-- Schema mydb
-- -----------------------------------------------------
-- -----------------------------------------------------
-- Schema tidder
-- -----------------------------------------------------

-- -----------------------------------------------------
-- Schema tidder
-- -----------------------------------------------------
CREATE SCHEMA IF NOT EXISTS `tidder` DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci ;
USE `tidder` ;

-- -----------------------------------------------------
-- Table `tidder`.`accounts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`accounts` (
  `id_account` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(25) NOT NULL,
  `email` VARCHAR(50) NOT NULL,
  `hashed_password` VARCHAR(255) NOT NULL,
  `birth_date` DATE NOT NULL,
  `creation_date` DATE NOT NULL,
  `karma` VARCHAR(50) NOT NULL DEFAULT 0,
  `profile_picture` LONGBLOB NULL DEFAULT NULL,
  `student_id` VARCHAR(14) NULL DEFAULT 'FuckTheSystem',
  PRIMARY KEY (`id_account`),
  UNIQUE INDEX `name` (`name` ASC) VISIBLE,
  UNIQUE INDEX `email` (`email` ASC) VISIBLE)
ENGINE = InnoDB
AUTO_INCREMENT = 35
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`subjects`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`subjects` (
  `id_subject` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(25) NOT NULL,
  `profile_picture` LONGBLOB NULL DEFAULT NULL,
  `id_owner` INT NOT NULL,
  `nsfw` TINYINT(1) NOT NULL DEFAULT false,
  `banner` LONGBLOB NULL DEFAULT NULL,
  `infos` VARCHAR(2048) NOT NULL DEFAULT 'No information has been set for this subtidder...',
  `can_create_post` TINYINT(1) NULL DEFAULT '1',
  PRIMARY KEY (`id_subject`),
  UNIQUE INDEX `name` (`name` ASC) VISIBLE,
  INDEX `id_owner` (`id_owner` ASC) VISIBLE,
  CONSTRAINT `subjects_ibfk_1`
    FOREIGN KEY (`id_owner`)
    REFERENCES `tidder`.`accounts` (`id_account`))
ENGINE = InnoDB
AUTO_INCREMENT = 41
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`posts`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`posts` (
  `id_post` INT NOT NULL AUTO_INCREMENT,
  `title` VARCHAR(255) NOT NULL,
  `media_url` LONGBLOB NULL DEFAULT NULL,
  `content` TEXT NOT NULL,
  `creation_date` DATETIME NOT NULL,
  `upvotes` INT NOT NULL DEFAULT 0,
  `downvotes` INT NOT NULL DEFAULT 0,
  `nsfw` TINYINT(1) NOT NULL DEFAULT false,
  `redacted` TINYINT(1) NOT NULL DEFAULT false,
  `pinned` TINYINT(1) NOT NULL DEFAULT false,
  `id_subject` INT NOT NULL,
  `id_author` INT NOT NULL,
  PRIMARY KEY (`id_post`),
  INDEX `id_subject` (`id_subject` ASC) VISIBLE,
  INDEX `id_author` (`id_author` ASC) VISIBLE,
  CONSTRAINT `posts_ibfk_1`
    FOREIGN KEY (`id_subject`)
    REFERENCES `tidder`.`subjects` (`id_subject`),
  CONSTRAINT `posts_ibfk_2`
    FOREIGN KEY (`id_author`)
    REFERENCES `tidder`.`accounts` (`id_account`))
ENGINE = InnoDB
AUTO_INCREMENT = 50
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`comments`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`comments` (
  `id_comment` INT NOT NULL AUTO_INCREMENT,
  `content` TEXT NOT NULL,
  `creation_date` DATETIME NOT NULL,
  `upvotes` INT NOT NULL DEFAULT 0,
  `downvotes` INT NOT NULL DEFAULT 0,
  `redacted` TINYINT(1) NOT NULL DEFAULT false,
  `id_author` INT NOT NULL,
  `response_to_id` INT NOT NULL DEFAULT -(1),
  `id_post` INT NOT NULL,
  PRIMARY KEY (`id_comment`),
  INDEX `id_author` (`id_author` ASC) VISIBLE,
  INDEX `response_to_id` (`response_to_id` ASC) VISIBLE,
  INDEX `id_post` (`id_post` ASC) VISIBLE,
  CONSTRAINT `comments_ibfk_1`
    FOREIGN KEY (`id_author`)
    REFERENCES `tidder`.`accounts` (`id_account`),
  CONSTRAINT `comments_ibfk_3`
    FOREIGN KEY (`id_post`)
    REFERENCES `tidder`.`posts` (`id_post`))
ENGINE = InnoDB
AUTO_INCREMENT = 77
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`global_access`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`global_access` (
  `id_global_access` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(25) NOT NULL,
  PRIMARY KEY (`id_global_access`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`global_roles`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`global_roles` (
  `id_global_role` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(25) NOT NULL,
  PRIMARY KEY (`id_global_role`),
  UNIQUE INDEX `name` (`name` ASC) VISIBLE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`global_roles_management`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`global_roles_management` (
  `id_account` INT NOT NULL,
  `id_global_role` INT NOT NULL,
  `id_global_access` INT NOT NULL,
  PRIMARY KEY (`id_account`, `id_global_role`, `id_global_access`),
  INDEX `id_global_role` (`id_global_role` ASC) VISIBLE,
  INDEX `id_global_access` (`id_global_access` ASC) VISIBLE,
  CONSTRAINT `global_roles_management_ibfk_1`
    FOREIGN KEY (`id_account`)
    REFERENCES `tidder`.`accounts` (`id_account`),
  CONSTRAINT `global_roles_management_ibfk_2`
    FOREIGN KEY (`id_global_role`)
    REFERENCES `tidder`.`global_roles` (`id_global_role`),
  CONSTRAINT `global_roles_management_ibfk_3`
    FOREIGN KEY (`id_global_access`)
    REFERENCES `tidder`.`global_access` (`id_global_access`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`subject_access`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`subject_access` (
  `id_subject_access` INT NOT NULL AUTO_INCREMENT,
  `create_post` TINYINT(1) NOT NULL DEFAULT '1',
  `pin_post` TINYINT(1) NOT NULL DEFAULT false,
  `manage_post` TINYINT(1) NOT NULL DEFAULT '0',
  `ban_user` TINYINT(1) NOT NULL DEFAULT false,
  `manage_role` TINYINT(1) NOT NULL DEFAULT '0',
  `manage_subtidder` TINYINT(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id_subject_access`))
ENGINE = InnoDB
AUTO_INCREMENT = 42
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`subject_roles`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`subject_roles` (
  `id_subject_role` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(25) NOT NULL,
  `id_subject` INT NOT NULL,
  `id_subject_access` INT NOT NULL,
  PRIMARY KEY (`id_subject_role`),
  UNIQUE INDEX `id_subject_role_UNIQUE` (`id_subject_role` ASC) VISIBLE,
  INDEX `id_subject` (`id_subject` ASC) VISIBLE,
  INDEX `id_subject_access` (`id_subject_access` ASC) VISIBLE,
  CONSTRAINT `subject_roles_ibfk_1`
    FOREIGN KEY (`id_subject`)
    REFERENCES `tidder`.`subjects` (`id_subject`),
  CONSTRAINT `subject_roles_ibfk_2`
    FOREIGN KEY (`id_subject_access`)
    REFERENCES `tidder`.`subject_access` (`id_subject_access`))
ENGINE = InnoDB
AUTO_INCREMENT = 39
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`has_subject_role`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`has_subject_role` (
  `id_account` INT NOT NULL,
  `id_subject` INT NOT NULL,
  `id_subject_role` INT NOT NULL,
  PRIMARY KEY (`id_account`, `id_subject`, `id_subject_role`),
  INDEX `id_subject` (`id_subject` ASC) VISIBLE,
  INDEX `id_subject_role` (`id_subject_role` ASC) VISIBLE,
  CONSTRAINT `has_subject_role_ibfk_1`
    FOREIGN KEY (`id_account`)
    REFERENCES `tidder`.`accounts` (`id_account`),
  CONSTRAINT `has_subject_role_ibfk_2`
    FOREIGN KEY (`id_subject`)
    REFERENCES `tidder`.`subjects` (`id_subject`),
  CONSTRAINT `has_subject_role_ibfk_3`
    FOREIGN KEY (`id_subject_role`)
    REFERENCES `tidder`.`subject_roles` (`id_subject_role`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`is_ban`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`is_ban` (
  `id_account` INT NOT NULL,
  `id_subject` INT NOT NULL,
  PRIMARY KEY (`id_account`, `id_subject`),
  INDEX `id_subject` (`id_subject` ASC) VISIBLE,
  CONSTRAINT `is_ban_ibfk_1`
    FOREIGN KEY (`id_account`)
    REFERENCES `tidder`.`accounts` (`id_account`),
  CONSTRAINT `is_ban_ibfk_2`
    FOREIGN KEY (`id_subject`)
    REFERENCES `tidder`.`subjects` (`id_subject`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`sessions`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`sessions` (
  `id_session` VARCHAR(45) NOT NULL,
  `id_account` INT NOT NULL,
  `creation_date` DATETIME(1) NOT NULL,
  PRIMARY KEY (`id_session`),
  UNIQUE INDEX `session_id_UNIQUE` (`id_session` ASC) VISIBLE)
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`subscribe_to_subject`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`subscribe_to_subject` (
  `id_account` INT NOT NULL,
  `id_subject` INT NOT NULL,
  PRIMARY KEY (`id_account`, `id_subject`),
  INDEX `id_subject` (`id_subject` ASC) VISIBLE,
  CONSTRAINT `subscribe_to_subject_ibfk_1`
    FOREIGN KEY (`id_account`)
    REFERENCES `tidder`.`accounts` (`id_account`),
  CONSTRAINT `subscribe_to_subject_ibfk_2`
    FOREIGN KEY (`id_subject`)
    REFERENCES `tidder`.`subjects` (`id_subject`))
ENGINE = InnoDB
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`vote_comment_to`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`vote_comment_to` (
  `idvote_comment_to` INT NOT NULL AUTO_INCREMENT,
  `id_account` INT NULL DEFAULT NULL,
  `id_comment` INT NULL DEFAULT NULL,
  `vote` INT NULL DEFAULT NULL,
  PRIMARY KEY (`idvote_comment_to`))
ENGINE = InnoDB
AUTO_INCREMENT = 23
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


-- -----------------------------------------------------
-- Table `tidder`.`vote_to`
-- -----------------------------------------------------
CREATE TABLE IF NOT EXISTS `tidder`.`vote_to` (
  `id_vote_to` INT NOT NULL AUTO_INCREMENT,
  `id_account` INT NULL DEFAULT NULL,
  `id_post` INT NULL DEFAULT NULL,
  `vote` INT NULL DEFAULT NULL,
  PRIMARY KEY (`id_vote_to`),
  UNIQUE INDEX `id_vote_to_UNIQUE` (`id_vote_to` ASC) VISIBLE)
ENGINE = InnoDB
AUTO_INCREMENT = 63
DEFAULT CHARACTER SET = utf8mb4
COLLATE = utf8mb4_0900_ai_ci;


SET SQL_MODE=@OLD_SQL_MODE;
SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS;
SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS;
