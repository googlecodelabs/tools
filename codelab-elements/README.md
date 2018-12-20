# Codelab Custom Elements

The next generation of the codelab elements without any framework or library
dependencies, only the [Custom Elements](https://html.spec.whatwg.org/multipage/custom-elements.html)
standard spec.

If this is a release bundle, produced with a `bazel build :bundle` command,
you should see `codelab-elements.js`, `codelab-elements.css` and other files,
ready to be added to an HTML page like the following. Only relevant parts are shown:

```html
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, minimum-scale=1.0, initial-scale=1.0, user-scalable=yes">
  <title>A codelab demo</title>
  <link rel="stylesheet" href="//fonts.googleapis.com/css?family=Source+Code+Pro:400|Roboto:400,300,400italic,500,700|Roboto+Mono">
  <link rel="stylesheet" href="//fonts.googleapis.com/icon?family=Material+Icons">
  <link rel="stylesheet" href="codelab-elements.css">
</head>
<body>
  <google-codelab-analytics gaid="UA-123"></google-codelab-analytics>
  <google-codelab codelab-gaid="UA-345" id="codelab-demo" title="A codelab demo">
    <google-codelab-step label="Overview" duration="1">
      Contents of the first step.
    </google-codelab-step>
    <google-codelab-step label="Second" duration="10">
      Contents of the second step.
    </google-codelab-step>
  </google-codelab>
  <script src="native-shim.js"></script>
  <script src="custom-elements.min.js"></script>
  <script src="prettify.js"></script>
  <script src="codelab-elements.js"></script>
</body>
</html>
```

You can download the latest version
from https://github.com/googlecodelabs/codelab-elements.

## Dev environment

All you need is [bazel](https://docs.bazel.build/versions/master/install.html).

After bazel is installed, try executing the following:

    bazel test --test_output=all //demo:hello_test

It will take some time at the first run because bazel will download and compile
all dependencies needed to work with the code and run tests. This includes
Google Closure library and compiler, Go language and browsers to run local JS
tests on.

### Building

Check out a demo HelloElement target. To build the element, execute the following:

    bazel build //demo:hello_bin

It should output something like this:

    INFO: Analysed target //demo:hello_bin (0 packages loaded).
    INFO: Found 1 target...
    Target //demo:hello_bin up-to-date:
      bazel-bin/demo/hello_bin.js
      bazel-bin/demo/hello_bin.js.map
    INFO: Elapsed time: 0.716s, Critical Path: 0.03s
    INFO: Build completed successfully, 1 total action

### Testing

All elements should have their test targets.
As a starting point, check out HelloElement tests:

    bazel test --test_output=errors //demo:hello_test

You should see something like this:

    INFO: Elapsed time: 5.394s, Critical Path: 4.60s
    INFO: Build completed successfully, 2 total actions
    //demo:hello_test_chromium-local                      PASSED in 4.6s

When things go wrong, it is usually easier to inspect and analyze output
with debug enabled:

    bazel test -s --verbose_failures --test_output=all --test_arg=-debug demo/hello_test

### Manual inspection from a browser

To browse things around manually with a real browser, execute the following:

    bazel run //tools:server

and navigate to http://localhost:8080.

## Notes

This is not an official Google product.
