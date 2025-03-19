# How to install

## Run this command

```bash
git clone https://github.com/Rahmat-Afriyadi/graphql-golang.git

```

```bash
cd graphql-golang
```

```bash
go mod tidy
```

## Set Enviroment Variable

```bash
MONGO_URI=
JWT_SECRET=
PORT=8080
```

## Run this command

```bash
go run main.go

```

## Open GraphQL in url

### Without JWT Middleware

```bash
http://localhost:8080/playground

```

### With JWT Middleware

```bash
http://localhost:8080/playground/auth

```
