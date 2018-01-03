/*
Navicat MySQL Data Transfer

Source Server         : 内网
Source Server Version : 50173
Source Host           : 10.20.104.175:3306
Source Database       : kwebtool

Target Server Type    : MYSQL
Target Server Version : 50173
File Encoding         : 65001

Date: 2018-01-03 20:43:37
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for configtemplates
-- ----------------------------
DROP TABLE IF EXISTS `configtemplates`;
CREATE TABLE `configtemplates` (
  `TemplateName` varchar(256) NOT NULL,
  `TemplateType` int(11) DEFAULT NULL,
  `EditTime` bigint(20) DEFAULT NULL,
  `Content` longtext,
  `Remarks` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`TemplateName`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for filelist
-- ----------------------------
DROP TABLE IF EXISTS `filelist`;
CREATE TABLE `filelist` (
  `filename` varchar(256) DEFAULT NULL,
  `time` bigint(20) DEFAULT NULL,
  `content` longblob
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for nodelist
-- ----------------------------
DROP TABLE IF EXISTS `nodelist`;
CREATE TABLE `nodelist` (
  `InnerIp` varchar(64) NOT NULL,
  `OuterIp` varchar(64) NOT NULL,
  `lastUpdateTime` int(11) DEFAULT NULL,
  `remarks` varchar(1024) DEFAULT NULL,
  PRIMARY KEY (`InnerIp`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for processlist
-- ----------------------------
DROP TABLE IF EXISTS `processlist`;
CREATE TABLE `processlist` (
  `ServerName` varchar(64) NOT NULL,
  `Insid` varchar(8) NOT NULL,
  `Innerip` varchar(32) NOT NULL,
  `Outerip` varchar(32) NOT NULL,
  `ConfigContent` text NOT NULL,
  `Path` varchar(256) NOT NULL,
  `Port` int(11) NOT NULL,
  `Status` int(11) NOT NULL,
  `Lastupdatetime` int(64) NOT NULL,
  `Other` varchar(156) NOT NULL,
  PRIMARY KEY (`Insid`,`ServerName`)
) ENGINE=MyISAM DEFAULT CHARSET=utf8;

-- ----------------------------
-- Table structure for rpclist
-- ----------------------------
DROP TABLE IF EXISTS `rpclist`;
CREATE TABLE `rpclist` (
  `Module` int(11) NOT NULL,
  `Object` varchar(255) NOT NULL,
  `Function` varchar(255) NOT NULL,
  `Data` text NOT NULL
) ENGINE=MyISAM DEFAULT CHARSET=utf8;
