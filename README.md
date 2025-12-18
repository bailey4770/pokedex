# Pokedex

A boot.dev guided project to build a pokedex in Go using the pokeAPI

## How to run

Simply download the latest release and run the binary.
You may have to `chmod +x pokedex` in order to run it (if using Linux).

Alternatively, if you prefer to build from source, or are not running on Linux:

- clone source code into a local directory.
- ensure go tool chain is downloaded.
- run `go build`. The command should build the appropriate binary for your OS and architecture.
- run the resulting binary.

If you find any bugs, please email me <bailey4770@outlook.com>

## Change log

### Store pokedex to disk between sessions

- Program stores pokedex to file using gob (encodes go structs) when exit command is executed.
- Loads pokedex file at start up.
- New reset command allows user to clear pokedex.

### Added command history and tab-completion

- Replaced the bufio.Scanner with a chyzer/readline.
- Library comes with very easily configured history, with navigation by arrow keys, and tab completion

### Restructured command and repl structure

- Originally, commands were added to the callback map in repl.go but declared in commands.go.
- Addition of a getCommands() func cleans up repl.go file.
- Now maintainers can add new commands from just the commands.go file. No annoying file switching is necessary.

### Refactored the http request function

- Original function was called by both map and mapb commands.
- OG function handled making request, parsing into JSON, and returning a specific struct type for map commands.
- now there is a GetData entry point to my pokeAPI package.
- this GetData function orchestrates helper functions:
  - first checks cache to see if we can avoid http request.
  - method on client makes the http request.
  - caches http response.
  - unmarshals json into a Generic struct.
- The benefits of these changes are:
  - increased separation of concerns. Allows for easier unit testing and debugging.
  - Abstractions now make more logical sense. Call stack is easy to follow.
  - http request logic is hidden from GetData func. Handled by the client.
  - Method on client now only deals with http request. nothing else.
  - Entry-point function takes a generic type as specified by the caller. This allows for greater code reuse.
  - Generics are defined in Response interface.

## To-do

- level up catching ability to make it easier to catch harder pokemon.
