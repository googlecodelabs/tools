## Developer setup instructions

1. `git clone git@github.com:web-animations/web-animations-next.git`
1. `git submodule update --init --recursive` (Necessary for running tests.)
1. Install [node](https://nodejs.org/en/) and make sure `npm` is in your $PATH
1. Run `npm install` in the respository to pull in development dependencies.
1. Run `npm install -g grunt grunt-cli` to get the build tools for the command line.
1. Run `grunt` to build the polyfill.
1. Run `grunt test` to run polyfill and web-platform-tests tests.


## Debugging tests

You can run the tests in an interactive mode with `grunt debug`. This starts the
Karma server once for each polyfill target for each test framework.
Navigate to `http://localhost:9876/debug.html` to open the test runner in your
browser of choice, all test results appear in the Javascript console.
Test failures can be accessed via `window.failures` and `window.formattedFailures`
once the tests have completed.

The polyfill target and tests can be specified as arguments to the `debug` task.  
Example: `grunt debug:web-animations-next:test/web-platform-tests/web-animations/animation/pause.html`  
Multiple test files may be listed with comma separation. Specifying files will output their URL in the command line.  
Example: `http://localhost:9876/base/test/web-platform-tests/web-animations/animation/pause.html`


## Design notes

[Design diagrams](https://drive.google.com/folderview?id=0B9rpPoIDv3vTNlZxOVp6a2tNa1E&usp=sharing)


## Publishing a release

1.  Determine the version number for the release

    * Increment the first number and reset others to 0 when there are large breaking changes
    * Increment the second number and reset the third to 0 when there are significant new, but backwards compatible features
    * Otherwise, increment the third number

1.  Add versioned release notes to `History.md`, for example:

        ### 3.13.37 — *November 1, 2001*

          * Fixed a bug where nothing worked

    Use the following to generate a summary of commits, but edit the list to contain only
    relevant information.

        git log --first-parent `git describe --tags --abbrev=0 web-animations-js/master`..web-animations-next/master --pretty=format:"  * %s"

1.  Specify the new version inside `package.json` (for NPM), for example:

    ```js
      "version": "3.13.37",
    ```

1.  Build the polyfill with `npm install && grunt` then update `README.md`'s Build Target Comparison with the current gzipped sizes.

1.  Submit both changes to web-animations-next then follow the procedure to push from web-animations-next to web-animations-js.

1.  Draft a [new release](https://github.com/web-animations/web-animations-js/releases) at the
    commit pushed to web-animations-js in step #4. Copy the release notes from `History.md`
    added in step #2.

1. Once you've pushed to web-animations-js, run `npm publish` from that checked-out folder

   To do this, you'll need to be a collaborator [on the NPM project](https://www.npmjs.com/package/web-animations-js), or have a collaborator help you.

1. If there are any breaking changes to the API in this release you must notify web-animations-changes@googlegroups.com.

   Only owners of the group may post to it so you may need to request ownership or ask someone to post it for you.

## Pushing from web-animations-next to web-animations-js

    git fetch web-animations-next
    git fetch web-animations-js
    git checkout web-animations-js/master
    git merge web-animations-next/master --no-edit --quiet
    npm install
    grunt
    # Optional "grunt test" to make sure everything still passes.
    git add -f *.min.js*
    git rm .gitignore
    git commit -m 'Add build artifacts from '`cat .git/refs/remotes/web-animations-next/master`
    git push web-animations-js HEAD:refs/heads/master

## Testing architecture

This is an overview of what happens when `grunt test` is run.

1. Polyfill tests written in mocha and chai are run.
    1. grunt creates a karma config with mocha and chai adapters.
    1. grunt adds the test/js files as includes to the karma config.
    1. grunt starts the karma server with the config and waits for the result.
    1. The mocha adaptor runs the included tests and reports the results to karma.
    1. karma outputs results to the console and returns the final pass/fail result to grunt.
1. web-platform-tests/web-animations tests written in testtharness.js are run.
    1. grunt creates a karma config with karma-testharness-adaptor.js included.
    1. grunt adds the web-platform-tests/web-animations files to the custom testharnessTests config in the karma config.
    1. grunt adds failure expectations to the custom testharnessTests config in the karma config.
    1. grunt starts the karma server with the config and waits for the result.
    1. The testharness.js adaptor runs the included tests (ignoring expected failures) and reports the results to karma.
    1. karma outputs results to the console and returns the final pass/fail result to grunt.
1. grunt exits successfully if both test runs passed.

