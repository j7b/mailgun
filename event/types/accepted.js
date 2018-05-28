{
  "event": "accepted",
  "id": "ncV2XwymRUKbPek_MIM-Gw",
  "timestamp": 1377211256.096436,
  "tags": [],
  "envelope": {
    "sender": "sender@example.com"
  },
  "campaigns": [],
  "user-variables": {},
  "flags": {
    "is-authenticated": false,
    "is-test-mode": false
  },
  "routes": [
    {
      "priority": 1,
      "expression": "match_recipient(\".*@samples.mailgun.org\")",
      "description": "Sample route",
      "actions": [
        "stop()",
        "forward(\"http://host.com/messages\")"
      ]
    }
  ],
  "message": {
    "headers": {
      "to": "",
      "message-id": "77AF5C3CA1416D93FC47AF8AD42A60AD@example.com",
      "from": "John Doe <sender@example.com>",
      "subject": "Test Subject"
    },
    "attachments": [],
    "recipients": [
      "recipient@example.com"
    ],
    "size": 6021
  },
  "recipient": "recipient@example.com",
  "method": "smtp"
}

