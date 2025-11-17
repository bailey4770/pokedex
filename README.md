# pokedex

A boot.dev guided project to build a pokedex in Go using the pokeAPI

## Change log

### Restructured command and repl structure

- Original boot.dev split repl and many commands into separate files.
- This meant the addition of new commands involved multiple files.
- By combining the repl and commands functions to one file, adding new commands is easier.

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

- store history of commands to cycle through with arrow keys.
- store pokedex to disk so it persists between sessions.
- level up catching ability to make it easier to catch harder pokemon.
