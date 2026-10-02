# Module 8 — Alerting

## Status: NOT STARTED

## What You're Building

When incident is created, notify someone. Start with one channel: Discord or Email. Pick one, do it well.

## Concepts to Learn First

- **Asynchronous notifications** — send alert without blocking the check loop
- **Retries** — what if notification fails?
- **Failure handling** — alert system must not crash the monitoring system
- **Notification policies** — when to alert, when not to (no spam)

## Pick One Channel First

### Option A: Discord Webhook
- Create Discord server → Settings → Integrations → Webhooks
- Get webhook URL
- POST JSON to webhook URL
- No external library needed

### Option B: Email (SMTP)
- Use Gmail SMTP or Mailgun free tier
- Go standard library + `net/smtp` or `gomail` package

**Recommendation:** Discord webhook. Simpler. No email credentials. Instant visible result.

## Discord Webhook Payload

**Exercise: the payload structs.** Discord expects a JSON body with a `content` string and an `embeds` array. Each embed has a `title`, `description`, and `color` (an int — red is `16711680`). Define two Go structs (`DiscordMessage` and `DiscordEmbed`) with the right fields and `json:"..."` tags so they marshal to that shape. Look at how you tagged structs in earlier modules.

**Exercise: `SendAlert`.** Write a function that takes a webhook URL, an incident, and an endpoint, and returns an `error`. In English:

1. Build a `DiscordMessage` with one embed — a title like "Incident Detected", a description naming the endpoint and start time (which package formats a string with values?), and the red color int.
2. Marshal it to JSON bytes (`encoding/json`).
3. POST those bytes to the webhook URL with content type `application/json` (`http.Post` wants an `io.Reader` — how do you turn a byte slice into one?).
4. Return the error.

## Alert Rules (avoid spam)

- Alert once when incident OPENS
- Alert once when incident RESOLVES
- Do NOT alert on every failed check

**Exercise: fire the alert without blocking.** In the incident service, when an incident is *just created*, call the alerter — but in a way that does not block the check loop while the HTTP POST to Discord happens. Which keyword runs a call in the background?

## Retry Logic

If Discord webhook fails, retry up to 3 times with backoff.

**Exercise: `sendWithRetry`.** Write a helper that takes a function (`func() error`) and a max-attempts count, and returns an `error`. In English:

1. Loop up to `maxAttempts` times.
2. Call the function each time; if it returns no error, return `nil` immediately (success).
3. If it failed, sleep for a growing delay before the next try (2s, then 4s, then 6s — how do you compute that from the loop index? remember `time.Sleep` wants a `time.Duration`).
4. If all attempts fail, return an error saying so.

This is a **higher-order function** — it takes another function as an argument. That's the pattern that makes retry reusable for any operation.

## Config

Store webhook URL in an environment variable, not hardcoded (e.g. a `DISCORD_WEBHOOK_URL` var). Read it the same way your config already reads other env values.

## File Structure Target

```
backend/
  internal/
    alerting/
      alerter.go     ← interface + implementations
      discord.go     ← Discord webhook impl
```

**Exercise: the `Alerter` interface.** Define an interface so you can add Email later without changing the incident service. It should declare two methods — one for sending an incident (opened) alert and one for a resolved alert — each taking an incident and an endpoint and returning an `error`. Recall Go interface syntax: `type Name interface { MethodName(args) returnType }`. The incident service depends on this interface, not on Discord directly — that's how you swap implementations later.

## Definition of Done

- [ ] Discord alert sent when incident opens
- [ ] Discord alert sent when incident resolves
- [ ] No duplicate alerts for same incident
- [ ] Alert failure does not crash worker
- [ ] Webhook URL stored in env var
- [ ] Retry logic on failure
- [ ] Can explain why alert is sent async

## Questions Before Moving On

1. Why send alert in a goroutine (`go alerter.Send(...)`) instead of directly?
2. What is exponential backoff? When would you use it?
3. Why define an `Alerter` interface instead of calling Discord directly?
4. How do you prevent alert spam during a long outage?

## Next: Module 9 — Redis & Queues
