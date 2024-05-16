# discord-export

very simple cli to extract all messages sent in any Discord server's ***text*** channel.

** due to Discord's own limitations, source URL to files/attachments/images/videos in the message cannot be extracted. text messages with markdown or embed will have strange formatting.

### usage: 

> discord-export <CHANNEL_ID>

### requirements: 

> auth.txt

paste valid discord token inside `auth.txt` file before running. `auth.txt` file should ONLY include the token, nothing else. no extra characters, no newlines, etc.

### output:

> logs/

directory where any runtime errors will be placed inside.

> message-exports/

directory where the json formatted exported messages will be placed inside.

### exported json format:

```
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

this JSON format should be very easy to work with. the array contains messages sent latest-oldest when looping top-bottom.