import express from "express";
import "dotenv/config";
import { env } from "./env";
import expressWs from "express-ws";

import tmi from 'tmi.js';
import { db } from "./db";
import { officedrummerRequests } from "./drizzle/schema";
import { eq } from "drizzle-orm";

const twitchClient = new tmi.Client({
	options: { debug: true },
	identity: {
		username: 'pepega_bot21',
		password: env.TWITCH_PASS
	},
	channels: [ 'opti_21' ]
});

twitchClient.connect().catch(console.error);

twitchClient.on('message', async (channel, tags, message, self) => {
  if (self) return;
  const splitMessage = message.split(' ');
  const command = splitMessage[0]?.toLowerCase();

	if(command === '!add') {
    // check if the user is subscribed to the channel
    if (tags.subscriber) {
      const userTwitchId = tags["user-id"];
      const twitchUsername = tags.username;

      if (!userTwitchId) {
        console.log("no twitch id");
        return;
      }

      if (!twitchUsername) {
        console.log("no twitch username");
        return;
      }

      const requestText = message.split(' ').slice(1).join(' ');

      if (!requestText) {
        twitchClient.say(channel, `@${tags.username} you forgot your request! Type !add <request> to add your request`);
        return;
      }

      if (requestText.length > 100) {
        twitchClient.say(channel, `@${tags.username} your request is too long! Please keep it under 100 characters`);
        return;
      }

      // sanitize the request text and make sure it does not have any html or special characters expect for spaces and dashes and colons and periods and commas and apostrophes
      // replace all the special characters with nothing

      const sanitizedRequestText = requestText.replace(/[^a-zA-Z0-9\s\-\:\.\,\']/g, '');

      const existingRequest = await db.query.officedrummerRequests.findFirst({
        where: (requests, { eq }) => eq(requests.twitchId, userTwitchId),
      })

      if (existingRequest) {
        twitchClient.say(channel, `@${tags.username} you already have a request in the queue!`);
        return;
      }

      await db.insert(officedrummerRequests).values({
        twitchId: userTwitchId,
        twitchUser: tags.username as string,
        requestText: sanitizedRequestText,
      })
      .catch((err) => {
        console.log(err);
      })

      twitchClient.say(channel, `@${tags.username} your request has been added to the queue!`);

    }
	}
  if (command === '!remove') {
    // check if the user is subscribed to the channel
    if (tags.subscriber) {
      const userTwitchId = tags["user-id"];
      if (!userTwitchId) {
        return;
      }

      const existingRequest = await db.query.officedrummerRequests.findFirst({
        where: (requests, { eq }) => eq(requests.twitchId, userTwitchId),
      })

      if (!existingRequest) {
        twitchClient.say(channel, `@${tags.username} you don't have a request in the queue!`);
        return;
      }

      await db.delete(officedrummerRequests).where(eq(officedrummerRequests.id, existingRequest.id))

    }
  }

  if(command === '!modadd') {
    // check if the user is a mod
    if (tags.mod) {
      const userTwitchId = tags["user-id"];
      const twitchUsername = tags.username;

      if (!userTwitchId || !twitchUsername) {
        return;
      }

      const usernameToUse = message.split(' ')[1];

      if (!usernameToUse) {
        twitchClient.say(channel, `@${tags.username} you forgot to add the user! Type !modadd <user> <request> to add your request`);
        return;
      }

      const requestText = message.split(' ').slice(2).join(' ');

      if (!requestText) {
        twitchClient.say(channel, `@${tags.username} you forgot your request! Type !modadd <request> to add your request`);
        return;
      }

      const existingRequest = await db.query.officedrummerRequests.findFirst({
        where: (requests, { eq }) => eq(requests.twitchId, userTwitchId),
      })

      if (existingRequest) {
        twitchClient.say(channel, `@${tags.username} you already have a request in the queue!`);
        return;
      }

      await db.insert(officedrummerRequests).values({
        twitchId: "0",
        twitchUser: usernameToUse as string,
        requestText: requestText,
      })
      .catch((err) => {
        console.log(err);
      })

      twitchClient.say(channel, `@${tags.username} your request has been added to the queue!`);

    }
  }

  if (command === '!modremove') {
    // check if the user is a mod
    if (tags.mod) {
      const usernameToRemove = message.split(' ')[1];

      if (!usernameToRemove) {
        twitchClient.say(channel, `@${tags.username} you forgot to add the user! Type !modremove <user> to remove their request`);
        return;
      }

      const existingRequest = await db.query.officedrummerRequests.findFirst({
        where: (requests, { eq }) => eq(requests.twitchUser, usernameToRemove),
      })

      if (!existingRequest) {
        twitchClient.say(channel, `@${tags.username} that user doesn't have a request in the queue!`);
        return;
      }

      await db.delete(officedrummerRequests).where(eq(officedrummerRequests.id, existingRequest.id))

    }
  }

});

export default twitchClient;

const app = express();

app.use(express.json());

app.get("/", async (req, res) => {
  return res.send("Yep he's bald");
});

app.listen(process.env.PORT ?? 3030, async () => {

  console.log("Server is running on port 3030");
    
});
