FROM golang

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Context of current work dir.
WORKDIR /app

# Copy and download dependency using go mod
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the code into the container
COPY . .

# Install necessery depsendencies
RUN go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest && go get github.com/go-sql-driver/mysql && go get github.com/mgutz/ansi

# Export necessary port
EXPOSE 8000

# Execute migrations
RUN migrate -database "mysql://root:colahonda@tcp(db:3306)/translatortelegrambot" -path migrations up; exit 0;

# Run the app
CMD [ "go", "run", "cmd/app/main.go"]
