##Telegram bot implemented on go lang
A simple bot that sends the first 5 YouTube videos upon a user query.

##Instruction
1. [Go to.](https://console.developers.google.com/project/)
2. Click Create Project button.
3. [Go to](https://console.developers.google.com/apis/api/youtube.googleapis.com/credentials) and copy Api key.
4. Paste you Api key to `YoutubeToken` constant in `main.go`.
5. Send `/newbot` message in Telegram BotFather (`@BotFather`).
6. Choose a name and username for your bot.
7. Copy and paste token to `TelegramToken` in `main.go`.
8. In your project directory run `go build .` and `./golang_telebot`.
9. ![Bot is work!](https://ibb.co/cbkg2LS)

