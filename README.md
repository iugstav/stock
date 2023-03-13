<h1 align="center">stock 📈</h1>

## Idea

Stock is an application calculates your directory size with its number of items or a file size. If the given path refers to a directory, it will call a function that creates a goroutine and traverse the dir until we find another subdirectory, which will also be traversed recursively. If the path refers to a file, it will simply return its name and file size.

## How to run

```shell
# see documentation
> go doc -all

# run application
> go run main.go
```

## Contributing

Wanna give suggestions? Write an issue.
