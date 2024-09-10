'use strict';

const autoprefixer = require('autoprefixer');
const cssdeclarationsorter = require('css-declaration-sorter');
const cssnano = require('cssnano');

exports.babel = () => {
  return {
    presets: ['@babel/preset-env']
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
    autoprefixer(), // Modern version doesn't need browser list, it's taken from .browserslistrc or package.json
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
    livereload: true,
    open: false,
    port: 8000,
  };
};
