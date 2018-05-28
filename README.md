# A Go interface to the mailgun API.
This package contains a client for the mailgun REST API, all endpoints categories are implemented (including webhooks) and most operations and parameters have corresponding types, functions and methods.

## Prerequisite
Some knowledge of the mailgun API is helpful. Some of the package design and expression of endpoints is easier to understand with reference to the API.

## Install
`` go get github.com/j7b/mailgun ``

## Usage
See [godoc](https://godoc.org/github.com/j7b/mailgun). Most user code probably only needs the client Send method, that sends a simple HTML email. More sophisticated mailing, with attachments, custom headers, and variables is possible with the ``message`` package.

A http.Handler that dispatches webhooks and lower-level facilities are in the corresponding package.

## Versioning
Commits to master are releases. Compatability with previous releases will be in the spirit of [the Go 1 compatability document](https://golang.org/doc/go1compat).

## Contributing
Bug reports are welcome. Feature requests are likely to conflict with design decisions and other considerations or be out of the scope of an API client. Accepted PR's must follow the conventions of existing packages, must not patch generated code, any corresponding mailgun API facilities must be published in the API documentation and not noted as deprecated or obsolete, and the PR author must open an issue first.

## License
[MIT](./LICENSE)
