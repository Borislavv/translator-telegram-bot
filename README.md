# Translator-telegram-bot

1. Migrations: 
    UP  : migrate -database "mysql://username:userpassword@tcp(localhost:3306)/translator_telegram_bot" -path migrations up
    Down: migrate -database "mysql://username:userpassword@tcp(localhost:3306)/translator_telegram_bot" -path migrations down