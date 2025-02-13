# httpserve

A little utility I wrote to myself to serve a folder as an HTTP static server. It is just tested under linux, but it should work on iOS and Windows too.

## Usage

You should run it and its default behaviour is exposing your `pwd` folder on an available port.
If you prefer to specify a port use the `--port <INT>` flag.
If you prefer to specify a folder to expose as root use `--path <STRING>` flag.
