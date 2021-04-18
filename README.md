# Jottings serve

Render markdown pages as html and serve them.

## installation

```
go install github.com/dataewan/jottingsserve
```

## usage

```
jottingsserve
```

Accessing `localhost:8080` will then display your rendered markdown pages.

### Command line options

```
USAGE:
   jottingsserve [global options] command [command options] [arguments...]

COMMANDS:
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value       (default: "8080")
   --directory value, -d value  (default: ".")
   --help, -h                   show help (default: false)
```


## API usage

I've dot an number of API endpoints that return JSON about the notes.


| Endpoint                    | Contents                                                                                                                 |
| --------------------------- | ------------------------------------------------------------------------------------------------------------------------ |
| /api/files                  | List all files                                                                                                           |
| /api/links                  | List all links between files and to external places like wikipedia                                                       |
| /api/links/{title}          | Get all links that link to a specific file file                                                                          |
| /api/files/{title}          | Get information about the file, including the filename and title                                                         |
| /api/files/{title}/contents | Break the markdown file into sections, for each section convert into html and provide both the html and the raw markdown |
