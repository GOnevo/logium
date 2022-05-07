# logium

**logium** is a tool for log monitoring in your browser.

## Installation

```shell
go install github.com/gonevo/logium@latest
```

## Usage

To start the **logium**, you only need to configure your own server and sources.

You can run **logium** with `-init` flag to create default config.json file

```shell
logium -init
```

## Config options

- **server.port**: port of **logium**, default: `3333`
- **server.last**: how many last lines from log **logium** will show for the new clients, default: `5`
- **server.auth**: map of users (login, password) for basic auth, default: `nil`
- **sources**: list of the sources. Each source has the `name` for UI and real `path` to the file

## Config example

```json
{
  "server": {
    "port": 3333,
    "last": 5,
    "auth": {
      "user1": "password1",
      "user2": "password2"
    }
  },
  "sources": [
    {
      "name": "errors",
      "path": "/var/error.log"
    }
  ]
}
```

## License

The logium is open-sourced software licensed under the [MIT license](https://opensource.org/licenses/MIT).
