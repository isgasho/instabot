# Instabot

Unofficial Instagram's private API written in Golang, inspired by Python **[instabot](https://github.com/instagrambot/instabot)**.

## Features

- Account login/logout, two-factor authentication supported.
- Post/delete a photo or an album to/from your feed.
- Like/unlike posts.
- Get post comments.
- Post/reply/delete comments.
- Like/unlike comments.
- Follow/unfollow users.
- Easily extend APIs by yourself.
- And more...

## Install

```sh
go get -u github.com/winterssy/instabot
```

## Usage

```go
import "github.com/winterssy/instabot"
```

## QuickStart

```go
bot := instabot.New("YOUR_USERNAME", "YOUR_PASSWORD")
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

_, err := bot.Login(ctx, true)
if err != nil {
    log.Fatal(err)
}

data, err := bot.GetMe(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Println(data)
```

## License

**[MIT](LICENSE)**