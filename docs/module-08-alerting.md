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

```go
type DiscordMessage struct {
    Content string `json:"content"`
    Embeds  []DiscordEmbed `json:"embeds"`
}

type DiscordEmbed struct {
    Title       string `json:"title"`
    Description string `json:"description"`
    Color       int    `json:"color"`  // red=16711680
}
```

```go
func SendAlert(webhookURL string, incident Incident, endpoint Endpoint) error {
    msg := DiscordMessage{
        Embeds: []DiscordEmbed{{
            Title:       "🔴 Incident Detected",
            Description: fmt.Sprintf("**%s** is DOWN\nIncident started: %s", endpoint.Name, incident.StartedAt),
            Color:       16711680,
        }},
    }
    body, _ := json.Marshal(msg)
    _, err := http.Post(webhookURL, "application/json", bytes.NewReader(body))
    return err
}
```

## Alert Rules (avoid spam)

- Alert once when incident OPENS
- Alert once when incident RESOLVES
- Do NOT alert on every failed check

```go
// In incident service:
if incidentJustCreated {
    go alerter.SendAlert(incident, endpoint)  // async, don't block
}
```

## Retry Logic

If Discord webhook fails, retry up to 3 times with backoff:

```go
func sendWithRetry(fn func() error, maxAttempts int) error {
    for i := 0; i < maxAttempts; i++ {
        err := fn()
        if err == nil {
            return nil
        }
        time.Sleep(time.Duration(i+1) * 2 * time.Second)  // 2s, 4s, 6s
    }
    return fmt.Errorf("alert failed after %d attempts", maxAttempts)
}
```

## Config

Store webhook URL in environment variable, not hardcoded.

```bash
DISCORD_WEBHOOK_URL=https://discord.com/api/webhooks/...
```

## File Structure Target

```
backend/
  internal/
    alerting/
      alerter.go     ← interface + implementations
      discord.go     ← Discord webhook impl
```

Define an interface so you can add Email later without changing incident service:

```go
type Alerter interface {
    SendIncidentAlert(incident Incident, endpoint Endpoint) error
    SendResolvedAlert(incident Incident, endpoint Endpoint) error
}
```

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
