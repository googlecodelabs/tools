'use strict';

const autoprefixer = require('autoprefixer');
const closureCompilerPackage = require('google-closure-compiler');
const cssdeclarationsorter = require('css-declaration-sorter');
const cssnano = require('cssnano');

exports.babel = () => {
  return {
    presets: ['es2015'],
  };
};

exports.closureCompiler = () => {
  return {
    compilation_level: 'ADVANCED',
    warning_level: 'VERBOSE',
    language_out: 'ECMASCRIPT5_STRICT',
    generate_exports: true,
    export_local_property_definitions: true,
    output_wrapper: '(function(window, document){\n%output%\n})(window, document);',
    js_output_file: 'cardsorter.js',
  };
};

exports.crisper = () => {
  return {
    scriptInHead: false,
  };
};

exports.htmlmin = () => {
  return {
    collapseWhitespace: true,
    conservativeCollapse: true,
    preserveLineBreaks: true,
    removeComments: true,
    useShortDoctype: true,
  };
};

exports.postcss = () => {
  return [
    autoprefixer({
      browsers: [
        'ie >= 10',
        'ie_mob >= 10',
        'ff >= 30',
        'chrome >= 34',
        'safari >= 7',
        'opera >= 23',
        'ios >= 7.1',
        'android >= 4.4',
        'bb >= 10',
      ],
    }),
    cssdeclarationsorter({ order: 'alphabetically' }),
    cssnano(),
  ];
};

exports.sass = () => {
  return {
    outputStyle: 'expanded',
    precision: 5,
  };
};

exports.uglify = () => {
  return {
    compress: {
      drop_console: true,
      keep_infinity: true,
      passes: 5,
    },
    output: {
      beautify: false,
    },
    toplevel: false,
  };
};

exports.vulcanize = () => {
  return {
    excludes: ['prettify.js'], // prettify produces errors when inlined
    inlineCss: true,
    inlineScripts: true,
    stripComments: true,
    stripExcludes: ['iron-shadow-flex-layout.html'],
  };
};

exports.webserver = () => {
  return {
    livereload: false,
  };
};
