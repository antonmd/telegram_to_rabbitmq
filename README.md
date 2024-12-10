# Telegram Bot with RabbitMQ Integration

This project is a Telegram bot written in Go. It receives incoming Telegram messages, processes them, and then publishes those messages to a RabbitMQ queue. This allows for decoupled, asynchronous processing of messages by other services.

## Features

*   **Telegram Integration**: Uses long polling to fetch incoming Telegram updates.
*   **RabbitMQ Integration**: Publishes incoming messages from Telegram chats (private, groups, channels) to a RabbitMQ queue.
*   **Dynamic Queue Creation (Optional)**: With the `mandatory` flag, queues are created on-demand if they don’t exist when a message arrives from a new chat.
*   **Graceful Shutdown**: Uses Go’s `context` and OS signals to handle graceful termination.
*   **Configuration via Environment Variables**: Easily set tokens, URLs, and connection strings without hard-coding values.

## Prerequisites

1.  **Go Environment**:  
    Make sure you have Go installed (e.g., go1.20+).  
      
    `go version`  
      
    
2.  **Telegram Bot Token**:  
    \- Create a new bot via [BotFather](https://t.me/BotFather).  
    \- Obtain the token in the form `123456789:ABC-Your_Token_Here`.  
      
    
3.  **RabbitMQ**:  
    Ensure RabbitMQ is running. For example, using Docker:  
      
    `docker run -d --hostname my-rabbit --name some-rabbit \ -e RABBITMQ_DEFAULT_USER=myuser \ -e RABBITMQ_DEFAULT_PASS=mypassword \ -p 5672:5672 \ -p 15672:15672 \ rabbitmq:3-management`  
      
    Once running, access the RabbitMQ management UI at [http://localhost:15672](http://localhost:15672) (user: `myuser`, password: `mypassword`).

## Configuration

Set the following environment variables before running the bot:

*   `TELEGRAM_BOT_TOKEN`: Your Telegram bot token from BotFather.
*   `TELEGRAM_API_URL`: (Optional) Defaults to `https://api.telegram.org` if not set.
*   `RABBITMQ_CONN`: The RabbitMQ connection string, e.g. `amqp://myuser:mypassword@localhost:5672/`.

**Example:**

export TELEGRAM\_BOT\_TOKEN="123456789:ABC-Your\_Token\_Here"
export TELEGRAM\_API\_URL="https://api.telegram.org"
export RABBITMQ\_CONN="amqp://myuser:mypassword@localhost:5672/"

## Running the Bot

1.  **Clone the Repository:**  
      
    `git clone https://github.com/yourusername/your-telegram-bot.git   cd your-telegram-bot`
  
3.  **Initialize and Tidy Modules:**  
      
    `go mod tidy`
  
5.  **Run the Bot:**  
      
    `go run cmd/bot/main.go`  
      
    You should see:  
    
        Configuration loaded successfully!
        Telegram client initialized!
        Connected to RabbitMQ
  
7.  **Interact with the Bot:**  
      
    \- Find your bot on Telegram by its username and send a message.  
    \- If you want to receive group messages, disable privacy mode using BotFather.  
    \- Check your terminal logs; you should see messages being logged and published to RabbitMQ.  
    \- In the RabbitMQ management console ([http://localhost:15672](http://localhost:15672)), observe that messages appear in your queues as configured.

## Privacy Mode in Groups

If the bot isn’t receiving group messages:

*   Talk to [BotFather](https://t.me/BotFather).
*   Select your bot, go to "Bot Settings" → "Group Privacy".
*   Disable privacy mode to allow the bot to receive all messages in the group.

## Advanced Configuration

*   **Queue Naming and Mandatory Publishing**: The code supports mandatory publishing. If a queue doesn’t exist for a chat, the message is returned, a queue is created, and the message is re-published. Adjust this logic in `queue/queue.go`.
*   **Multiple Channels or Connections**: Currently uses a single connection and channel. For more concurrency or parallelism, you can open multiple channels from the same connection.

## Contributing

Open issues or pull requests for bugs or feature requests. Run `go fmt` and lint your code before contributing.

## License

Distributed under the [MIT License](LICENSE).

tput. These options include:

*   headingStyle (setext or atx)
*   horizontalRule (\*, -, or \_)
*   bullet (\*, -, or )
*   codeBlockStyle (indented or fenced)
*   fence (\` or ~)
*   emDelimiter (\_ or \*)
*   strongDelimiter (\*\* or \_\_)
*   linkStyle (inlined or referenced)
*   linkReferenceStyle (full, collapsed, or shortcut)