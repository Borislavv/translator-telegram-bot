# Translator-telegram-bot

**Author**: Glazunov Borislav

**Github**: github.com/Borislavv **|**
**Gitlab**: gitlab.com/Zendden

**Email**: glazunov2142@gmail.com **|**
**Phone**: +7(904)-939-89-83 (Telegram)

# Installation and running

1. Choose the target configuration file and fill it. Config files path: `root_app_dir/config`.

2. Install [golang](https://go.dev/doc/install) (optional, you can run you app with binary file), but you will need to run [migrations](https://github.com/golang-migrate/migrate), then install [it](https://github.com/golang-migrate/migrate/blob/master/cmd/migrate/README.md).

3. Migrations (database in: `utf8_general_ci`):<br />
&nbsp;&nbsp;&nbsp;&nbsp;UP  : `migrate -database "mysql://user:pass@tcp(localhost:3306)/translatortelegrambot" -path migrations up`<br />
&nbsp;&nbsp;&nbsp;&nbsp;Down: `migrate -database "mysql://user:pass@tcp(localhost:3306)/translatortelegrambot" -path migrations down`

4. Compile executable file run it or just run the code (in the second case, binary file will be compiling and remove after all):<br />
&nbsp;&nbsp;&nbsp;&nbsp;4.1 Compile and run:<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Compiling:<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- `cd cmd/app`<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- `go build`<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Runnig:<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- `./main --config-path ./../../config/.env.dev.toml`<br />
&nbsp;&nbsp;&nbsp;&nbsp;4.2. Run without compiling (actualy 'with', but you will not see it):<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;Running:<br />
&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;&nbsp;- `go run cmd/app/main.go --config-path config/.env.dev.toml`<br />
