🚧 Ongoing Project 🚧
This is a work-in-progress rewrite in **Go** of [switcher-backend](https://github.com/IngSoft1-Capybaras/switcher-backend), a university group project for the subject **Software Engineering 1** .

# El Switcher - Backend

Switcher is a real-time multiplayer board game where players have to form figures moving the pieces, according to their movement cards. Each player will have a deck with figures to discard as they form on the board, and the first to finish will be the winner. It is the web version of a real board game called [**El Switcher**](https://maldon.com.ar/blog/projects/el-switcher/).

## Installation

1. Clone the project:

```sh
git clone https://github.com/NachoGz/switcher-backend-go.git
cd switcher-backend-go
```

2. Intall dependancies

```sh
go mod download
```

3. Create a database named **switcher**

```sh
psql -U <user> -h localhost -W
CREATE DATABASE switcher;
```

4. Create a `.env` file:

```txt
PORT=8080
DB_URL=postgres://user:password@localhost:5432/switcher?sslmode=disable
```

4. Run the application

```sh
go run cmd/switcher/main.go
```

To run the front-end, you can go to the [switcher-frontend](https://github.com/IngSoft1-Capybaras/switcher-frontend) repository and follow the instructions.

## Testing

Run all tests:

```sh
go test ./...
```

Run tests with coverage:

```sh
go test -cover ./...
```

## License

[MIT](https://choosealicense.com/licenses/mit/)
