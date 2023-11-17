-- Current sql file was generated after introspecting the database
-- If you want to run this migration please uncomment this code before executing migrations
/*
CREATE TABLE `officedrummer-frontend-new_account` (
	`userId` varchar(255) NOT NULL,
	`type` varchar(255) NOT NULL,
	`provider` varchar(255) NOT NULL,
	`providerAccountId` varchar(255) NOT NULL,
	`refresh_token` text,
	`access_token` text,
	`expires_at` int,
	`token_type` varchar(255),
	`scope` varchar(255),
	`id_token` text,
	`session_state` varchar(255),
	CONSTRAINT `officedrummer-frontend-new_account_provider_providerAccountId_pk` PRIMARY KEY(`provider`,`providerAccountId`)
);
--> statement-breakpoint
CREATE TABLE `officedrummer-frontend-new_requests` (
	`id` bigint AUTO_INCREMENT NOT NULL,
	`name` varchar(256) NOT NULL,
	`twitchId` varchar(256) NOT NULL,
	`requestText` varchar(256),
	`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updatedAt` timestamp ON UPDATE CURRENT_TIMESTAMP,
	`sliceSize` int,
	CONSTRAINT `officedrummer-frontend-new_requests_id_pk` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `officedrummer-frontend-new_session` (
	`sessionToken` varchar(255) NOT NULL,
	`userId` varchar(255) NOT NULL,
	`expires` timestamp NOT NULL,
	CONSTRAINT `officedrummer-frontend-new_session_sessionToken_pk` PRIMARY KEY(`sessionToken`)
);
--> statement-breakpoint
CREATE TABLE `officedrummer-frontend-new_user` (
	`id` varchar(255) NOT NULL,
	`name` varchar(255),
	`email` varchar(255) NOT NULL,
	`emailVerified` timestamp(3) DEFAULT CURRENT_TIMESTAMP(3),
	`image` varchar(255),
	CONSTRAINT `officedrummer-frontend-new_user_id_pk` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE TABLE `officedrummer-frontend-new_verificationToken` (
	`identifier` varchar(255) NOT NULL,
	`token` varchar(255) NOT NULL,
	`expires` timestamp NOT NULL,
	CONSTRAINT `officedrummer-frontend-new_verificationToken_identifier_token_pk` PRIMARY KEY(`identifier`,`token`)
);
--> statement-breakpoint
CREATE TABLE `officedrummer-frontend-new_wheelStatus` (
	`id` bigint AUTO_INCREMENT NOT NULL,
	`status` varchar(256),
	`created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
	`updatedAt` timestamp ON UPDATE CURRENT_TIMESTAMP,
	CONSTRAINT `officedrummer-frontend-new_wheelStatus_id_pk` PRIMARY KEY(`id`)
);
--> statement-breakpoint
CREATE INDEX `userId_idx` ON `officedrummer-frontend-new_account` (`userId`);--> statement-breakpoint
CREATE INDEX `twitchId_idx` ON `officedrummer-frontend-new_requests` (`twitchId`);--> statement-breakpoint
CREATE INDEX `userId_idx` ON `officedrummer-frontend-new_session` (`userId`);--> statement-breakpoint
CREATE INDEX `status_idx` ON `officedrummer-frontend-new_wheelStatus` (`status`);
*/