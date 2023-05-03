# gadget

Go utility library for rapid development of CLI programs and services.

## Repo Status

WIP

## Import

```go
// generally the modules within gadget are imported, e.g.:
import "github.com/uhdaptuhbl/gadget/exec"
import "github.com/uhdaptuhbl/gadget/harness"
import "github.com/uhdaptuhbl/gadget/logging"
```

## Documentation

WIP
TODO: add link to main docs

## Development

### Release Versioning

When determining the project structure, [a post](https://stackoverflow.com/a/64705638) led to discovering a better way than submodules, and possibly also [git subtree](https://www.atlassian.com/git/tutorials/git-subtree); Go will recognize version tags in multi-module repositories in the form of: [`<prefix>/<version>`](https://github.com/golang/go/wiki/Modules#publishing-a-release).

For example, creating the tag `teapot/v0.3.4` on the main branch will version the teapot module such that Go will only see that version tag specifically for the teapot module, allowing others to be versioned separately. This does mean that tagging releases does get slightly more complicated, but for now seems like a worthwhile tradeoff to make.

TODO: determine whether or not to adopt `git subtree` or multi-repo organization structure

### Contributing

TBD
