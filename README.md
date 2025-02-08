"# snap-rqt"

## Debugging

1. Start delve headless `dlv debug --headless --listen=:2345 --log --api-version=2 main.go`
2. Attach VSCode to delve. Add this launch.json

```
{
  "configurations": [
    {
      "name": "My remote debug",
      "type": "go",
      "request": "attach",
      "mode": "remote",
      "remotePath": "${workspaceFolder}",
      "port": 2345,
      "host": "127.0.0.1"
    }
  ]
}

```

```

```
