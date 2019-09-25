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

**Currently, the build process requires bazel version 0.18.1.**

After bazel is installed, try executing the following:

    bazel test --test_output=all //codelab-elements/demo:hello_test

It will take some time at the first run because bazel will download and compile
all dependencies needed to work with the code and run tests. This includes
Google Closure library and compiler, Go language and browsers to run local JS
tests on.

### Building

Check out a demo HelloElement target. To build the element, execute the following:

    bazel build //codelab-elements/demo:hello_bin

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

    bazel test --test_output=errors //codelab-elements/demo:hello_test

You should see something like this:

    INFO: Elapsed time: 5.394s, Critical Path: 4.60s
    INFO: Build completed successfully, 2 total actions
    //codelab-elements/demo:hello_test_chromium-local                      PASSED in 4.6s

When things go wrong, it is usually easier to inspect and analyze output
with debug enabled:

    bazel test -s --verbose_failures --test_output=all --test_arg=-debug codelab-elements/demo/hello_test

### Manual inspection from a browser

To browse things around manually with a real browser, execute the following:

    bazel run //codelab-elements/tools:server

and navigate to http://localhost:8080.

# Deploying the built artifacts to a CDN

We now store the built artifacts in a world-readable Google Cloud Storage bucket (gs://codelab-elements). The list of artifacts we need to serve from this bucket are as follows:

- codelab-elements.css
- native-shim.js
- custom-elements.min.js
- prettify.js
- codelab-elements.js

The manual process for deploying these artifacts is as follows:

```
cd googlecodelabs/tools # start at repo level
bazel build //... # build everything
mkdir pkg
cd pkg
unzip ../bazel-genfiles/bundle.zip
export SRC="codelab-elements.css native-shim.js custom-elements.min.js prettify.js codelab-elements.js codelab-index.css codelab-index.js"
gsutil -m cp -a public-read $SRC  gs://codelab-elements
```
Then you need to include these artifacts using their full path, which is:

`https://storage.googleapis.com/codelab-elements/FILE`

where 'FILE' is replaced by one of the source files defined in the SRC variable above.

Here's an index.html obtaining the codelab-elements from Google Cloud Storage:

```
<!doctype html>
<html>
<head>
  <meta charset="utf-8">
  <meta name="viewport" content="width=device-width, minimum-scale=1.0,initial-scale=1.0, user-scalable=yes">
  <title>A codelab demo</title>
  <link rel="stylesheet" href="https://fonts.googleapis.com/css?family=Source+Code+Pro:400|Roboto:400,300,400italic,500,700|Roboto+Mono">
  <link rel="stylesheet" href="https://fonts.googleapis.com/icon?family=Material+Icons">
  <link rel="stylesheet" href="https://storage.googleapis.com/codelab-elements/codelab-elements.css">
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
  <script src="https://storage.googleapis.com/codelab-elements/native-shim.js"></script>
  <script src="https://storage.googleapis.com/codelab-elements/custom-elements.min.js"></script>
  <script src="https://storage.googleapis.com/codelab-elements/prettify.js"></script>
  <script src="https://storage.googleapis.com/codelab-elements/codelab-elements.js"></script>
</body>
</html>
```

## NPM

The library is published on NPM.js at: https://www.npmjs.com/package/codelab-elements

We also support a build workflow using NPM. To build the library run from the root of the repo:

```bash
# Install dependencies. This takes care of installing the right version of Bazel.
npm install

# Build the library
npm run build

# The output is a zip file under bazel-genfiles
ls bazel-genfiles npm_dist.zip

# Publish a new version of the library to NPM
npm run dist
```

## Notes

This is not an official Google product.
