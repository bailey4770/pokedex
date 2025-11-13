# pokedex

A boot.dev guided project to build a pokedex in Go using the pokeAPI

# Change log

## Refactored the http request function

- Original function was called by both map and mapb commands.
- OG function handled making request, parsing into JSON, and returning a specific struct type for map commands.
- now there is a GetData entry point to my pokeAPI package.
- this GetData function orchestrates helper functions:
  - calls separate function to check cache first.
  - method on client makes the http request.
  - caches http response.
  - calls separate function to parse JSON into a Generic Type
  - returns parsed struct
- The benefits of these changes are:
  - increased separation of concerns. Allows for easier unit testing and debugging.
  - Abstractions now make more logical sense. Call stack is easy to follow.
  - HTTP request logic, checking cache, and adding to cache logic are hidden from generic GetData and caller.
  - Method on client now only deals with http request. nothing else.
  - Entry-point function takes a generic type as specified by the caller. This allows for greater code reuse.
  - Generics are defined in Resposne interface.

## To-do

- create tests for map and mapb commands
  - it was unclear that these were broken until too late
  - generate a list of 200 indexed locations
  - then test map and mapb against expected output.
- create new tests for client API
