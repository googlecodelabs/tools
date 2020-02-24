workspace(name = "googlecodelabs_custom_elements")

maven_server(
    name = "default",
    url = "https://repo1.maven.org/maven2/"
)

maven_jar(
    name = "org_apache_httpcomponents_httpclient",
    artifact = "org.apache.httpcomponents:httpclient:4.5.5",
)

maven_jar(
    name = "org_apache_httpcomponents_httpmime",
    artifact = "org.apache.httpcomponents:httpmime:4.5.5",
)

maven_jar(
    name = "org_apache_httpcomponents_httpcore",
    artifact = "org.apache.httpcomponents:httpcore:4.4.9",
)

maven_jar(
    name = "org_apache_commons_exec",
    artifact = "org.apache.commons:commons-exec:1.3",
)

maven_jar(
    name = "org_seleniumhq_selenium_api",
    artifact = "org.seleniumhq:selenium-api:3.9.1",
)

maven_jar(
    name = "org_seleniumhq_selenium_remote_driver",
    artifact = "org.seleniumhq.selenium:selenium-remote-driver:3.8.1",
)

maven_jar(
    name = "net_java_dev_jna_platform",
    artifact = "net.java.dev:jna-client:4.5.1",
)

maven_jar(
    name = "net_java_dev_jna",
    artifact = "net.java.dev:jna:4.5.1",
)

maven_jar(
    name = "net_bytebuddy",
    artifact = "net.bytebuddy:byte-buddy:1.7.9",
)

maven_jar(
    name = "com_squareup_okio",
    artifact = "com.squareup:okio:1.14.0",
)

maven_jar(
    name = "com_squareup_okhttp3_okhttp",
    artifact = "com.squareup.okhttp3:okhttp:3.9.1",
)

maven_jar(
    name = "cglib_nodep",
    artifact = "cglib:cglib-nodep:3.2.6",
)

maven_jar(
    name = "junit",
    artifact = "junit:junit:4.12",
)

maven_jar(
    name = "commons_logging",
    artifact = "commons-logging:commons-logging:1.2",
)

maven_jar(
    name = "commons_codec",
    artifact = "commons-codec:commons-codec:1.11",
)

maven_jar(
    name = "org_hamcrest_core",
    artifact = "org.hamcrest:hamcrest-core:1.3",
)

new_http_archive(
    name = "com_google_javascript_closure_compiler",
    build_file = "third_party/BUILD.closure",
    url = "https://repo1.maven.org/maven2/com/google/javascript/closure-compiler-unshaded/v20180805/closure-compiler-unshaded-v20180805.jar",
)

# Required by io_bazel_rules_webtesting.
skylib_ver = "f9b0ff1dd3d119d19b9cacbbc425a9e61759f1f5"
http_archive(
    name = "bazel_skylib",
    sha256 = "ce27a2007deda8a1de65df9de3d4cd93a5360ead43c5ff3017ae6b3a2abe485e",
    strip_prefix = "bazel-skylib-{v}".format(v=skylib_ver),
    urls = [
        "https://github.com/bazelbuild/bazel-skylib/archive/{v}.tar.gz".format(v=skylib_ver),
    ],
)

rules_closure_ver = "0.9.0"
http_archive(
    name = "io_bazel_rules_closure",
    sha256 = "054717a2e6a415001bc4c608b208723526bdf6cace3592ca6efb3749ba18ce21",
    strip_prefix = "rules_closure-{v}".format(v=rules_closure_ver),
    url = "https://github.com/shawnbuso/rules_closure/archive/{v}.zip".format(v=rules_closure_ver),
)
load("@io_bazel_rules_closure//closure:defs.bzl", "closure_repositories")
closure_repositories()

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "53c8222c6eab05dd49c40184c361493705d4234e60c42c4cd13ab4898da4c6be",
    url = "https://github.com/bazelbuild/rules_go/releases/download/0.10.0/rules_go-0.10.0.tar.gz",
)
load("@io_bazel_rules_go//go:def.bzl", "go_rules_dependencies", "go_register_toolchains")
go_rules_dependencies()
go_register_toolchains()

rules_webtesting_ver = "936c760cff973a63031be0d0518b40a228e224e3"
http_archive(
    name = "io_bazel_rules_webtesting",
    sha256 = "797b75e792a34728a6a3846c7c3d3ad669f12cd8490b888cc969bad93d236b1b",
    strip_prefix = "rules_webtesting-{v}".format(v=rules_webtesting_ver),
    url = "https://github.com/bazelbuild/rules_webtesting/archive/{v}.zip".format(v=rules_webtesting_ver),
)
load(
    "@io_bazel_rules_webtesting//web:repositories.bzl",
    "browser_repositories",
    "web_test_repositories",
)
web_test_repositories()
browser_repositories(chromium = True)

prettify_ver = "2013-03-04"
new_http_archive(
    name = "prettify",
    build_file = "third_party/BUILD.prettify",
    strip_prefix = "code-prettify-{v}".format(v=prettify_ver),
    url = "https://github.com/google/code-prettify/archive/{v}.zip".format(v=prettify_ver),
)

new_http_archive(
    name = "polyfill",
    build_file = "third_party/BUILD.polyfill",
    sha256 = "9606cdeacbb67f21fb495a4b0a0e5ea6a137fc453945907822e1b930e77124d4",
    strip_prefix = "custom-elements-1.0.8",
    url = "https://github.com/webcomponents/custom-elements/archive/v1.0.8.zip",
)

git_repository(
    name = "io_bazel_rules_sass",
    remote = "https://github.com/bazelbuild/rules_sass.git",
    tag = "0.0.3",
)

load("@io_bazel_rules_sass//sass:sass.bzl", "sass_repositories")

sass_repositories()
