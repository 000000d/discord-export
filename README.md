# discord-export

A very simple CLI to extract all text messages sent in any Discord server's ***text*** channels.

> [!WARNING]
> Due to Discord's own limitations, source URL to files/attachments/images/videos in the message cannot be extracted.
> Text messages with markdown or embed will have strange formatting or may be omitted entirely.

### Usage:

```console
discord-export <CHANNEL_ID> ...
```

#### Example:

```console
discord-export 123 456 789
```

### Requirements:

`auth.txt`: Paste valid discord token inside `auth.txt` file before running. The `auth.txt` file should ONLY include the token, nothing else. no extra characters, no extra words, etc.

### Output:

`logs/`: Directory where any runtime errors will be put into.

`message-exports/`: Directory where the JSON formatted exported messages will be put into.

### Exported JSON format:

```json
{
    "channel_id": "123",
    "messages": [
        {
            "message_id": "456",
            "message": "hello",
            "user_id": "789",
            "user": "USER1",
            "bot": false
        },
        {
            "message_id": "012",
            "message": "hello1",
            "user_id": "345",
            "user": "USER2",
            "bot": true
        }
    ]
}
```

The `messages` array contains messages sent latest-oldest when looping top-bottom. Sometimes the `message` field may be empty or omitted even if `message_id` exists, eg. bot sends embed or other improper format of message that is beyond the scope of this software.
