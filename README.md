# glui
go list ui

![ezgif-4-7e4e643bc6](https://user-images.githubusercontent.com/2561547/191633768-3a462b1b-e36b-48e0-a486-5525fab4185d.gif)

comes from the need to easily go through package dependencies, provides basic information and the ability to edit them using the default editor.

## notes:
- the golang binary must be installed since it uses the command "go list"
- by default the calls go list with the arguments: '-e --json ...', beware that in some cases this might be not ideal, I welcome any recommendations.
- uses the EDITOR environment variable.

## installation:
```
go install github.com/alvarolm/glui@latest
```
