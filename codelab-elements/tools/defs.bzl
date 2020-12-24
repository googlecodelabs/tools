# Copyright 2018 Google Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#      http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

"""Repo's bazel rules and macros."""

load(
    "@io_bazel_rules_closure//closure:defs.bzl",
    _closure_js_binary_alias = "closure_js_binary",
)
load(
    "@io_bazel_rules_closure//closure:defs.bzl",
    _closure_js_library_alias = "closure_js_library",
)
load("@io_bazel_rules_webtesting//web:web.bzl", "web_test_suite")

def concat(ext):
    """Returns a genrule command to concat files with the extension ext."""
    return "ls $(SRCS) | grep -E '\\.{ext}$$' | xargs cat > $@".format(ext = ext)

def closure_js_library(**kwargs):
    """Invokes closure_js_library with non-test compilation defaults.

    Args:
      **kwargs: Additional arguments, passed to _closure_js_library_alias.
    """
    kwargs.setdefault("convention", "GOOGLE")
    suppress = kwargs.pop("suppress", [])
    suppress.append("JSC_UNKNOWN_EXPR_TYPE")
    kwargs.update(dict(suppress = suppress))
    _closure_js_library_alias(**kwargs)

def closure_js_binary(**kwargs):
    """Invokes closure_js_binary with non-test compilation defaults.

    Args:
      **kwargs: Additional arguments, passed to _closure_js_binary_alias.
    """
    kwargs.setdefault("compilation_level", "ADVANCED")
    kwargs.setdefault("dependency_mode", "STRICT")
    kwargs.setdefault("language", "ECMASCRIPT5_STRICT")
    kwargs.setdefault("defs", [
        "--assume_function_wrapper",
        "--rewrite_polyfills=false",
        "--new_type_inf",
        "--export_local_property_definitions",
        "--language_out=ES5_STRICT",
        "--isolation_mode=IIFE",
        "--generate_exports",
        "--jscomp_warning=newCheckTypes",
        "--jscomp_off=newCheckTypesExtraChecks",
        "--hide_warnings_for=closure/goog",
    ])
    _closure_js_binary_alias(**kwargs)

def _gen_test_html_impl(ctx):
    """Implementation of the gen_test_html rule."""
    ctx.actions.expand_template(
        template = ctx.file._template,
        output = ctx.outputs.html_file,
        substitutions = {
            "{{TEST_FILE_JS}}": ctx.attr.test_file_js,
        },
    )
    runfiles = ctx.runfiles(files = [ctx.outputs.html_file], collect_default = True)
    return [DefaultInfo(runfiles = runfiles)]

# A rule used by js_test to generate default test.html file
# suitable for running Closure-based JS tests.
# The test_file_js argument specifies the name of the JS file containing tests,
# typically created with closure_js_binary.
# The output is created from gen_test_html.template file.
gen_test_html = rule(
    implementation = _gen_test_html_impl,
    attrs = {
        "test_file_js": attr.string(mandatory = True),
        "_template": attr.label(
            default = Label("//codelab-elements/tools:gen_test_html.template"),
            allow_single_file = True,
        ),
    },
    outputs = {"html_file": "%{name}.html"},
)

def js_test(
        name,
        srcs,
        browsers,
        data = None,
        deps = None,
        compilation_level = None,
        css = None,
        entry_points = None,
        html = None,
        suppress = None,
        visibility = None,
        **kwargs):
    """A lower level macro which creates JS tests suite.

    It creates three targets: <name>_lib closure_js_library,
    <name>_bin closure_js_binary with the former as a dependencies,
    and <name> web_test_suite with the <name>_bin in its data dependencies.
    All targets have testonly attribute set to True.

    For more details about closure_js_library and closure_js_binary,
    see https://github.com/bazelbuild/rules_closure.

    Args:
      name: The name of the test target.
      srcs: A list of test source files with _test.js suffix.
      browsers: A list of browsers to run on. See rules_webtesting
        for list of supported browsers: https://goo.gl/VVH8tP.
      data: A list of data dependencies passed to closure_js_library.
      deps: list of code dependencies passed to closure_js_library.
      compilation_level: Closure compiler compilation level.
      css: A CSS class renaming target passed to closure_js_binary.
        It must point to a closure_css_binary rule.
      entry_points: List of unreferenced namespaces which should not
        be pruned by the compiler. See //demo:hello_test
        for a usage example.
      html: An HTML file which declares the generated closure_js_binary
        target it its <script src="..."> sources.
        If not specified, a default is generated using gen_test_html
        with gen_<name> target.
      suppress: List of codes the linter should ignore,
        passed to the generated closure_js_library target.
      visibility: Target visibility.
      **kwargs: Additional arguments, passed to the web_test_suite target.
    """
    _ignore = [compilation_level]
    if not srcs:
        fail("js_test rules can not have an empty 'srcs' list")
    for src in srcs:
        if not src.endswith("_test.js"):
            fail("js_test srcs must be files ending with _test.js")

    _closure_js_library_alias(
        name = "%s_lib" % name,
        srcs = srcs,
        data = data,
        deps = deps,
        suppress = suppress,
        visibility = visibility,
        testonly = True,
    )

    _closure_js_binary_alias(
        name = "%s_bin" % name,
        deps = [":%s_lib" % name],
        entry_points = entry_points,
        css = css,
        debug = True,
        defs = [
            "--rewrite_polyfills=false",
            "--jscomp_off=analyzerChecks",
            "--export_local_property_definitions",
        ],
        language = "ECMASCRIPT_2015",
        dependency_mode = "LOOSE",
        formatting = "PRETTY_PRINT",
        visibility = visibility,
        testonly = True,
    )

    if not html:
        gen_test_html(
            name = "gen_%s" % name,
            test_file_js = "%s_bin.js" % name,
        )
        html = "gen_%s" % name

    web_test_suite(
        name = name,
        data = [":%s_bin" % name, html],
        test = "//codelab-elements/tools:webtest",
        args = ["--test_url", "$(location %s)" % html],
        browsers = browsers,
        visibility = visibility,
        **kwargs
    )

def closure_js_test(**kwargs):
    """A handy higher level macro built on top of js_test.

    It sets useful defaults suitable for running all JS tests in this repo:
    - adds closure/library:testing to the deps
    - adds custom elements polyfill to the data
    - suppresses linter's known false positives
    - sets default list of browsers to run tests on

    Args:
      **kwargs: Additional arguments, passed to the js_test target.
    """
    deps = kwargs.pop("deps", [])
    deps.append("@io_bazel_rules_closure//closure/library:testing")

    data = kwargs.pop("data", [])
    data.append("@polyfill//:custom_elements")

    suppress = kwargs.pop("suppress", [])
    suppress.append("JSC_EXTRA_REQUIRE_WARNING")

    kwargs.update(dict(deps = deps, data = data, suppress = suppress))
    kwargs.setdefault("browsers", [
        # For experimental purposes only. Eventually you should
        # create your own browser definitions.
        # TODO: Add firefox.
        "@io_bazel_rules_webtesting//browsers:chromium-local",
    ])

    js_test(**kwargs)
