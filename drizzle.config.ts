import { type Config } from "drizzle-kit";
import { env } from "./src/env";

export default {
  schema: "./src/db/schema/*",
  out: "./src/drizzle",
  driver: "mysql2",
  dbCredentials: {
    uri: env.DATABASE_URL,
  },
  tablesFilter: ["officedrummer_*"],
} satisfies Config;
