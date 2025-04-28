# discord-export

A very simple CLI to extract all messages sent in any Discord server's ***text*** channel.

> [!WARNING]
> Due to Discord's own limitations, source URL to files/attachments/images/videos in the message cannot be extracted.
> Text messages with markdown or embed will have strange formatting.

### Usage:

```console
discord-export <CHANNEL_ID>
```

### Requirements:

> auth.txt

Paste valid discord token inside `auth.txt` file before running. `auth.txt` file should ONLY include the token, nothing else. no extra characters, no newlines, etc.

### Output:

> logs/

Directory where any runtime errors will be put into.

> message-exports/

Directory where the JSON formatted exported messages will be put into.

### Exported JSON format:

```json
{
    "channel_id": "123",
    "messages": [
        {
            "message": "hello",
            "user_id": "456",
            "user": "USER1"
        },
        {
            "message": "hello1",
            "user_id": "789",
            "user": "USER2"
        }
    ]
}
```

This JSON format should be very easy to work with. The array contains messages sent latest-oldest when looping top-bottom.
