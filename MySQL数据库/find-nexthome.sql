/*
 Navicat Premium Dump SQL

 Source Server         : 10001
 Source Server Type    : MySQL
 Source Server Version : 90400 (9.4.0)
 Source Host           : localhost:3306
 Source Schema         : find-nexthome

 Target Server Type    : MySQL
 Target Server Version : 90400 (9.4.0)
 File Encoding         : 65001

 Date: 02/04/2026 12:45:38
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for accounts
-- ----------------------------
DROP TABLE IF EXISTS `accounts`;
CREATE TABLE `accounts`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `nickname` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '昵称',
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '账号=手机号/自设账号',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '密码',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 10 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of accounts
-- ----------------------------
INSERT INTO `accounts` VALUES (1, '关圣帝君', '261001', '123456');
INSERT INTO `accounts` VALUES (2, '瑶池第二吃货', '262002', '123456');
INSERT INTO `accounts` VALUES (3, '爱上福建女', '262003', '123456');
INSERT INTO `accounts` VALUES (4, '不想学习GO', '261004', '123456');
INSERT INTO `accounts` VALUES (5, 'xx', '242005', '123456');
INSERT INTO `accounts` VALUES (6, '测试', '242006', '123456');
INSERT INTO `accounts` VALUES (7, 'qwe444', 'qwe444', '123456');
INSERT INTO `accounts` VALUES (9, 'God', 'admin', '123456');

-- ----------------------------
-- Table structure for address
-- ----------------------------
DROP TABLE IF EXISTS `address`;
CREATE TABLE `address`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '账号',
  `rel_way` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '联系方式',
  `location` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL COMMENT '联系账号',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 5 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of address
-- ----------------------------
INSERT INTO `address` VALUES (1, '261001', '电话', '19722695149');
INSERT INTO `address` VALUES (2, '262002', '陌陌', '17912698457');
INSERT INTO `address` VALUES (3, 'qwe444', 'qq', '15115151515');
INSERT INTO `address` VALUES (4, 'admin', 'QQ', '123456789');

-- ----------------------------
-- Table structure for goods
-- ----------------------------
DROP TABLE IF EXISTS `goods`;
CREATE TABLE `goods`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `name` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `description` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `price` decimal(10, 2) NULL DEFAULT NULL,
  `type` int NULL DEFAULT NULL,
  `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 14 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of goods
-- ----------------------------
INSERT INTO `goods` VALUES (9, '242005', '普中开发板', '89c51开发板，2025年9月入，至今只使用两次', 50.00, 2, '1775047682609_242005_good');
INSERT INTO `goods` VALUES (10, '242005', '防风火机', '没啥用，夏天有文字点蚊香用两次，买回去点根烟', 3.00, 2, '1775047748166_242005_good');
INSERT INTO `goods` VALUES (11, '262002', '半卷卫生纸', '吃完饭擦一擦', 1.00, 2, '1775047857668_262002_good');
INSERT INTO `goods` VALUES (12, '262002', '钥匙串', '配有水果刀和指甲剪、小剪子，实用', 9.00, 2, '1775047904137_262002_good');
INSERT INTO `goods` VALUES (13, 'admin', '拖鞋', '主播的拖鞋', 10.00, 2, '1775048858559_admin_good');

-- ----------------------------
-- Table structure for need
-- ----------------------------
DROP TABLE IF EXISTS `need`;
CREATE TABLE `need`  (
  `id` int NOT NULL AUTO_INCREMENT,
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `description` varchar(1000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 19 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of need
-- ----------------------------
INSERT INTO `need` VALUES (16, '242005', '想要一个能用来赚钱的锦囊', '1775047602048_242005_need');
INSERT INTO `need` VALUES (17, '262002', '找不到图片了，凑合着用把', '1775047948906_262002_need');
INSERT INTO `need` VALUES (18, '261001', '求一个上课防睡着的神器', '1775048027685_261001_need');

-- ----------------------------
-- Table structure for self_made
-- ----------------------------
DROP TABLE IF EXISTS `self_made`;
CREATE TABLE `self_made`  (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '自增主键',
  `account` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `work_name` varchar(5000) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `introduce` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  `price` decimal(10, 2) NULL DEFAULT NULL,
  `type` int NULL DEFAULT NULL,
  `image` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 11 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of self_made
-- ----------------------------
INSERT INTO `self_made` VALUES (8, 'qwe444', '纯手工指甲盖', '无添加纯天然手指甲盖，可入药', 3.00, 1, '1775047357202_qwe444_selfmade');
INSERT INTO `self_made` VALUES (9, '242005', '纸飞机', '手叠纸飞机，尊重艺术', 66.00, 1, '1775047465420_242005_selfmade');
INSERT INTO `self_made` VALUES (10, '242005', '眼镜水', '自来水制作眼镜水，纯自制', 6.00, 1, '1775047523362_242005_selfmade');

SET FOREIGN_KEY_CHECKS = 1;
