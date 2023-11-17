// src/env.mjs
import { createEnv } from "@t3-oss/env-core";
import { z } from "zod";

export const env = createEnv({
  /*
   * Serverside Environment variables, not available on the client.
   * Will throw if you access these variables on the client.
   */
  server: {
    MONGO_URI: z.string().url(),
    DATABASE_URL: z
      .string()
      .url()
      .refine(
        (str) => !str.includes("YOUR_MYSQL_URL_HERE"),
        "You forgot to change the default URL"
      ),
    TWITCH_CLIENT_ID: z.string().min(1),
    TWITCH_CLIENT_SECRET: z.string().min(1),
    TWITCH_PASS: z.string().min(1),
    TWITCH_REDIRECT_URI: z.string().min(1),
    SESSION_SECRET: z.string().min(1),
    WEBHOOK_SECRET: z.string().min(1),
    BASE_URL: z.string().url(),
    OFD_USER_ID: z.string().min(1),
  },
  runtimeEnv: process.env,
});
