import { drizzle } from "drizzle-orm/planetscale-serverless";
import { connect } from "@planetscale/database";
import { env } from "../env";
import * as schema from "../drizzle/schema"

// create the connection
const connection = connect({
    url: env.DATABASE_URL,
});
 
export const db = drizzle(connection, {
    schema
});