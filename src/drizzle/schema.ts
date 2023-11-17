import { mysqlTable, mysqlSchema, AnyMySqlColumn, index, primaryKey, varchar, text, int, bigint, timestamp } from "drizzle-orm/mysql-core"
import { sql } from "drizzle-orm"


export const officedrummerAccount = mysqlTable("officedrummer_account", {
	userId: varchar("userId", { length: 255 }).notNull(),
	type: varchar("type", { length: 255 }).notNull(),
	provider: varchar("provider", { length: 255 }).notNull(),
	providerAccountId: varchar("providerAccountId", { length: 255 }).notNull(),
	refreshToken: text("refresh_token"),
	accessToken: text("access_token"),
	expiresAt: int("expires_at"),
	tokenType: varchar("token_type", { length: 255 }),
	scope: varchar("scope", { length: 255 }),
	idToken: text("id_token"),
	sessionState: varchar("session_state", { length: 255 }),
},
(table) => {
	return {
		userIdIdx: index("userId_idx").on(table.userId),
		officedrummerAccountProviderProviderAccountIdPk: primaryKey({ columns: [table.provider, table.providerAccountId], name: "officedrummer_account_provider_providerAccountId_pk"}),
	}
});

export const officedrummerRequests = mysqlTable("officedrummer_requests", {
	id: bigint("id", { mode: "number" }).autoincrement().notNull(),
	twitchUser: varchar("twitchUser", { length: 256 }).notNull(),
	twitchId: varchar("twitchId", { length: 256 }).notNull(),
	requestText: varchar("requestText", { length: 256 }),
	createdAt: timestamp("created_at", { mode: 'string' }).default(sql`CURRENT_TIMESTAMP`).notNull(),
	updatedAt: timestamp("updatedAt", { mode: 'string' }).onUpdateNow(),
	sliceSize: int("sliceSize"),
},
(table) => {
	return {
		twitchIdIdx: index("twitchId_idx").on(table.twitchId),
		officedrummerRequestsIdPk: primaryKey({ columns: [table.id], name: "officedrummer_requests_id_pk"}),
	}
});

export const officedrummerSession = mysqlTable("officedrummer_session", {
	sessionToken: varchar("sessionToken", { length: 255 }).notNull(),
	userId: varchar("userId", { length: 255 }).notNull(),
	expires: timestamp("expires", { mode: 'string' }).notNull(),
},
(table) => {
	return {
		userIdIdx: index("userId_idx").on(table.userId),
		officedrummerSessionSessionTokenPk: primaryKey({ columns: [table.sessionToken], name: "officedrummer_session_sessionToken_pk"}),
	}
});

export const officedrummerUser = mysqlTable("officedrummer_user", {
	id: varchar("id", { length: 255 }).notNull(),
	name: varchar("name", { length: 255 }),
	email: varchar("email", { length: 255 }).notNull(),
	emailVerified: timestamp("emailVerified", { fsp: 3, mode: 'string' }).default(sql`CURRENT_TIMESTAMP(3)`),
	image: varchar("image", { length: 255 }),
},
(table) => {
	return {
		officedrummerUserIdPk: primaryKey({ columns: [table.id], name: "officedrummer_user_id_pk"}),
	}
});

export const officedrummerVerificationToken = mysqlTable("officedrummer_verificationToken", {
	identifier: varchar("identifier", { length: 255 }).notNull(),
	token: varchar("token", { length: 255 }).notNull(),
	expires: timestamp("expires", { mode: 'string' }).notNull(),
},
(table) => {
	return {
		officedrummerVerificationTokenIdentifierTokenPk: primaryKey({ columns: [table.identifier, table.token], name: "officedrummer_verificationToken_identifier_token_pk"}),
	}
});

export const officedrummerWheelStatus = mysqlTable("officedrummer_wheelStatus", {
	id: bigint("id", { mode: "number" }).autoincrement().notNull(),
	status: varchar("status", { length: 256 }),
	createdAt: timestamp("created_at", { mode: 'string' }).default(sql`CURRENT_TIMESTAMP`).notNull(),
	updatedAt: timestamp("updatedAt", { mode: 'string' }).onUpdateNow(),
},
(table) => {
	return {
		statusIdx: index("status_idx").on(table.status),
		officedrummerWheelStatusIdPk: primaryKey({ columns: [table.id], name: "officedrummer_wheelStatus_id_pk"}),
	}
});