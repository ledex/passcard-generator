# PassCard
PassCard is a small tool that generates password-cards from a given identifier.

## Installing/Running
Installing all needed dependencies:
```bash
go install github.com/a-h/templ/cmd/templ@latest
go get -v
```
To run PassCard:
```bash
templ generate
go run .
```
To run PassCard in dev-mode:
```bash
templ generate --watch --proxy="http://localhost:8080" --cmd="go run ."
```