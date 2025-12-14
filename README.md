<!-- gomarkdoc:embed:start -->


<!-- gomarkdoc:embed:end -->

## Dependencies

1. postgres
1. clang-21

## Helpful Developer Cmds

```sh
go run ./bs build.bs                # Initial command to build the build system (only ever needs to be run once)
./sbbs install.bashAutocomplete     # Optional: installs autocomplete for the build system (local to current user)
./sbbs --help                       # Prints all targets
./sbbs <target>                     # Runs the selected target
```
