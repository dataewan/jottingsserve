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
The rendering is pretty basic.
If you want a more functional interface, look at [jotting-frontend](https://github.com/dataewan/jotting-frontend).

### Command line options

```
NAME:
   jottingsserve - A new cli application

USAGE:
   jottingsserve [global options] command [command options] [arguments...]

VERSION:
   0.0.1

DESCRIPTION:
   Tools for working with markdown linked notes

COMMANDS:
   checklinks  Check for missing links
   help, h     Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --port value, -p value       (default: "8080")
   --directory value, -d value  (default: ".")
   --help, -h                   show help (default: false)
   --version, -v                print the version (default: false)
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


## CLI usage

This is also a command line application (CLI).


### checklinks

Check if there are any links that point to other markdown notes in your notes,
but the note being pointed to doesn't exist.
You might have deleted that note or renamed it.

```
jottingsserve checklinks
```
