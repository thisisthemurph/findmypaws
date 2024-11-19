# Find my paws server

## Clerk webhooks

Clerk provides webhooks for listening to events such as user created, updated, and deleted.

For local development, Clerk uses Ngrok to forward events to the local server. 

The following need to be done before this sync will work locally:
- Set up Ngrok
  - Steps 1 & 2 are required from the [Ngrok quickstart guide](https://ngrok.com/docs/getting-started/#step-1-install)
- Set up a webhook endpoint in Clerk
- Add the Clerk Signing Secret to the .env file

Run Ngrok forwarding to the local server, the port should be whatever port the server is listening on:

```
ngrok http --url=well-plainly-midge.ngrok-free.app 42069
```


**Links**
- [Clerk data sync documentation](https://clerk.com/docs/integrations/webhooks/sync-data)
  - Note that this documentation is for Next.js, but is a good guide to get started.
- [Manage Clerk webhooks](https://dashboard.clerk.com/apps/app_2ns6pcXvTCGrSf3kD5nMSAYGu7X/instances/ins_2ns6pa6yUAQJrW1eBi5iqTBhO0f/webhooks)
- [Ngrok forwarding](https://dashboard.ngrok.com/domains/rd_2p3YGPtYOyZQsojdxZSh1kKunVp)
