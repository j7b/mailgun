{
  "event": "failed",
  "id": "pVqXGJWhTzysS9GpwF2hlQ",
  "timestamp": 1377198389.769129,
  "severity": "permanent",
  "tags": [],
  "envelope": {
    "transport": "smtp",
    "sender": "postmaster@samples.mailgun.org",
    "sending-ip": "184.173.153.199"
  },
  "delivery-status": {
    "message": "Relay Not Permitted",
    "code": 550,
    "description": null
  },
  "campaigns": [],
  "reason": "bounce",
  "user-variables": {},
  "flags": {
    "is-authenticated": true,
    "is-test-mode": false
  },
  "message": {
    "headers": {
      "to": "recipient@example.com",
      "message-id": "20130822185902.31528.73196@samples.mailgun.org",
      "from": "John Doe <sender@example.com>",
      "subject": "Test Subject"
    },
    "attachments": [],
    "recipients": [
      "recipient@example.com"
    ],
    "size": 557
  },
  "recipient": "recipient@example.com",
}

