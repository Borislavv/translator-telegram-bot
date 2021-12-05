# Translator-telegram-bot

**Author**: Glazunov Borislav

**Github**: github.com/Borislavv **|**
**Gitlab**: gitlab.com/Zendden

**Email**: glazunov2142@gmail.com **|**
**Phone**: +7(904)-939-89-83 (Telegram)

# Installation and running

1. Choose the target configuration file and fill it. Config files path: `root_app_dir/config`.

2. Migrations: 
    UP  : `migrate -database "mysql://username:userpassword@tcp(localhost:3306)/translator_telegram_bot" -path migrations up`
    Down: `migrate -database "mysql://username:userpassword@tcp(localhost:3306)/translator_telegram_bot" -path migrations down`

3. Compile executable file run it or just run the code (in the second case, binary file will be compiling and remove after all):
    2.1 Compile and run:
            Compiling: 
                - `cd cmd/app`
                - `go build`
            Runnig:
                - `./main --config-path ./../../config/.env.dev.toml`
    2.2. Run without compiling (actualy 'with', but you will not see it):
            Running: 
                - `go run cmd/app/main.go --config-path config/.env.dev.toml`
