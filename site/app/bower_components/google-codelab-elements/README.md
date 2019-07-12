# Google Codelab Components

A set of Polymer elements, core of the Google Codelabs platform.

## Dev setup

`bower install` is the obvious first thing to do.

To play with the demos, an easy thing would be to use
[polyserve](https://github.com/PolymerLabs/polyserve):
`npm install -g polyserve`.

Fire up the server with `polyserve` command and point your browser to:

    http://localhost:8080/components/codelab-components/

Use [web-component-tester](https://github.com/Polymer/web-component-tester) to run tests.
Can be installed with `npm install -g web-component-tester`, or just fire up `polyserve`
and navigate to `/components/codelab_components/test/<test-file>`.
