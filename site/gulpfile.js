'use strict';

// Gulp
const gulp = require('gulp');

// Gulp plugins
const babel = require('gulp-babel');
const closureCompilerPackage = require('google-closure-compiler');
const closureCompiler = closureCompilerPackage.gulp();
const crisper = require('gulp-crisper');
const gulpif = require('gulp-if');
const htmlmin = require('gulp-htmlmin');
const merge = require('merge-stream');
const postcss = require('gulp-html-postcss');
const rename = require('gulp-rename');
const sass = require('gulp-sass');
const through = require('through2');
const useref = require('gulp-useref');
const vulcanize = require('gulp-vulcanize');
const watch = require('gulp-watch');
const webserver = require('gulp-webserver');

// Uglify ES6
const uglifyes = require('uglify-es');
const composer = require('gulp-uglify/composer');
const uglify = composer(uglifyes, console);

// Other helpers
const args = require('yargs').argv;
const childprocess = require('child_process');
const claat = require('./tasks/helpers/claat');
const del = require('del');
const fs = require('fs-extra');
const gcs = require('./tasks/helpers/gcs');
const glob = require('glob');
const opts = require('./tasks/helpers/opts');
const path = require('path');
const serveStatic = require('serve-static');
const spawn = childprocess.spawn;
const swig = require('swig-templates');
const url = require('url');

// DEFAULT_GA is the default Google Analytics tracker ID
const DEFAULT_GA = 'UA-49880327-14';

// DEFAULT_VIEW_META_PATH is the default path to view metadata.
const DEFAULT_VIEW_META_PATH = 'app/views/default/view.json';

// DEFAULT_VIEW_TMPL_PATH is the default path to view template.
const DEFAULT_VIEW_TMPL_PATH = 'app/views/default/index.html';

// DEFAULT_CATEGORY is the default name for uncategorized codelabs.
const DEFAULT_CATEGORY = 'Default';

// BASE_URL is the canonical base URL where the site will reside. This should
// always include the protocol (http:// or https://) and NOT including a
// trailing slash.
const BASE_URL = args.baseUrl || 'https://example.com';

// CODELABS_DIR is the directory where the actual codelabs exist on disk.
// Despite being a constant, this can be overridden with the --codelabs-dir
// flag.
const CODELABS_DIR = args.codelabsDir || '.';

// CODELABS_ENVIRONMENT is the environment for which to build codelabs.
const CODELABS_ENVIRONMENT = args.codelabsEnv || 'web';

// CODELABS_FILTER is an inclusion filter against codelab IDs.
const CODELABS_FILTER = args.codelabsFilter || '*';

// CODELABS_FORMAT is the output format for which to build codelabs.
const CODELABS_FORMAT = args.codelabsFormat || 'html';

// CODELABS_NAMESPACE is the content namespace.
const CODELABS_NAMESPACE = (args.codelabsNamespace || 'codelabs').replace(/^\/|\/$/g, '');

// DELETE_MISSING controls whether missing files at the destination are deleted.
// The default value is true.
const DELETE_MISSING = !!args.deleteMissing || false;

// DRY_RUN indicates if dry should be used.
const DRY_RUN = !!args.dry;

// PROD_BUCKET is the default bucket for prod.
const PROD_BUCKET = gcs.bucketName(args.prodBucket || 'DEFAULT_PROD_BUCKET');

// STAGING_BUCKET is the default bucket for staging.
const STAGING_BUCKET = gcs.bucketName(args.stagingBucket || 'DEFAULT_STAGING_BUCKET');

// VIEWS_FILTER is the filter to use for view inclusion.
const VIEWS_FILTER = args.viewsFilter || '*';

// clean:build removes the build directory
gulp.task('clean:build', (callback) => {
  return del('build')
});

// clean:dist removes the dist directory
gulp.task('clean:dist', (callback) => {
  return del('dist')
});

// clean:js removes the built javascript
// NOTE: this is not included in the default 'clean' task
gulp.task('clean:js', (callback) => {
  return del('app/js/bundle')
});

// clean removes all built files
gulp.task('clean', gulp.parallel(
  'clean:build',
  'clean:dist',
));

// build:codelabs copies the codelabs from the directory into build.
gulp.task('build:codelabs', (done) => {
  copyFilteredCodelabs('build');
  done();
});

// build:scss builds all the scss files into the dist dir
gulp.task('build:scss', () => {
  return gulp.src('app/**/*.scss')
    .pipe(sass(opts.sass()))
    .pipe(gulp.dest('build'));
});

// build:css builds all the css files into the dist dir
gulp.task('build:css', () => {
  const srcs = [
    'app/elements/codelab-elements/*.css',
  ];

  return gulp.src(srcs, { base: 'app/' })
    .pipe(gulp.dest('build'));
});

// build:html builds all the HTML files
gulp.task('build:html', () => {
  const streams = [];

  streams.push(gulp.src(`app/views/${VIEWS_FILTER}/view.json`, { base: 'app/' })
    .pipe(generateView())
    .pipe(useref({ searchPath: ['app'] }))
    .pipe(gulpif('*.js', babel(opts.babel())))
    .pipe(gulp.dest('build'))
    .pipe(gulpif(['*.html', '!index.html'], generateDirectoryIndex()))
  );

  streams.push(gulp.src(`app/views/${VIEWS_FILTER}/*.{css,gif,jpeg,jpg,png,svg,tff}`, { base: 'app/views' })
    .pipe(gulp.dest('build')));

  const otherSrcs = [
    'app/404.html',
    'app/browserconfig.xml',
    'app/robots.txt',
    'app/site.webmanifest',
  ]
  streams.push(gulp.src(otherSrcs, { base: 'app/' })
    .pipe(gulp.dest('build'))
  );

  return merge(...streams);
});

// build:images builds all the images into the build directory.
gulp.task('build:images', () => {
  const srcs = [
    'app/images/**/*',
    'app/favicon.ico',
  ];

  return gulp.src(srcs, { base: 'app/' })
    .pipe(gulp.dest('build'));
});

// build:js builds all the javascript into the dest dir
gulp.task('build:js', (callback) => {
  let streams = [];

  if (!fs.existsSync('app/js/bundle/cardsorter.js')) {
    // cardSorter is compiled into app/js, not build/scripts, because it is
    // vulcanized directly into the HTML.
    const cardSorterSrcs = [
      'app/js/claat/ui/cards/cardsorter.js',
      'app/js/claat/ui/cards/cardsorter_export.js',
    ];
    streams.push(gulp.src(cardSorterSrcs, { base: 'app/' })
      .pipe(closureCompiler(opts.closureCompiler(), { platform: 'javascript' }))
      .pipe(babel(opts.babel()))
      .pipe(gulp.dest('app/js/bundle'))
    );
  }

  const bowerSrcs = [
    'app/bower_components/webcomponentsjs/webcomponents-lite.min.js',
    // Needed for async loading - remove after polymer/polymer#2380
    'app/bower_components/google-codelab-elements/shared-style.html',
    'app/bower_components/google-prettify/src/prettify.js',
  ];
  streams.push(gulp.src(bowerSrcs, { base: 'app/' })
    .pipe(gulpif('*.js', babel(opts.babel())))
    .pipe(gulp.dest('build'))
  );

  return merge(...streams);
});

gulp.task('build:elements_js', () => {
  const srcs = [
    'app/elements/codelab-elements/*.js'
  ];

  return gulp.src(srcs, { base: 'app/' })
    .pipe(gulp.dest('build'));
})

// build:vulcanize vulcanizes html, js, and css
gulp.task('build:vulcanize', () => {
  const srcs = [
    'app/elements/codelab.html',
    'app/elements/elements.html',
  ];
  return gulp.src(srcs, { base: 'app/' })
    .pipe(vulcanize(opts.vulcanize()))
    .pipe(crisper(opts.crisper()))
    .pipe(gulp.dest('build'));
});

// build builds all the assets
gulp.task('build', gulp.series(
  'clean',
  'build:codelabs',
  'build:css',
  'build:scss',
  'build:html',
  'build:images',
  'build:js',
  'build:elements_js',
  'build:vulcanize',
));

// copy copies the built artifacts in build into dist/
gulp.task('copy', (callback) => {
  // Explicitly do not use gulp here. It's too slow and messes up the symlinks
  fs.rename('build', 'dist', callback);
});

// minify:css minifies the css
gulp.task('minify:css', () => {
  const srcs = [
    'dist/**/*.css',
    '!dist/codelabs/**/*',
    '!dist/elements/codelab-elements/*.css',
  ]
  return gulp.src(srcs, { base: 'dist/' })
    .pipe(postcss(opts.postcss()))
    .pipe(gulp.dest('dist'));
});

// minify:css minifies the html
gulp.task('minify:html', () => {
  const srcs = [
    'dist/**/*.html',
    '!dist/codelabs/**/*',
  ]
  return gulp.src(srcs, { base: 'dist/' })
    .pipe(postcss(opts.postcss()))
    .pipe(htmlmin(opts.htmlmin()))
    .pipe(gulp.dest('dist'));
});

// minify:js minifies the javascript
gulp.task('minify:js', () => {
  const srcs = [
    'dist/**/*.js',
    '!dist/codelabs/**/*',
    '!dist/elements/codelab-elements/*.js',
  ]
  return gulp.src(srcs, { base: 'dist/' })
    .pipe(uglify(opts.uglify()))
    .pipe(gulp.dest('dist'));
});

// minify minifies all minifiable things in dist
gulp.task('minify', gulp.parallel(
  'minify:css',
  'minify:html',
  'minify:js',
));

// dist packages the build for distribution, compressing and condensing where
// appropriate.
gulp.task('dist', gulp.series(
  'build',
  'copy',
  'minify',
));

// watch:css watches css files for changes and re-builds them
gulp.task('watch:css', () => {
  gulp.watch('app/**/*.scss', gulp.series('build:css'));
});

// watch:html watches html files for changes and re-builds them
gulp.task('watch:html', () => {
  const srcs = [
    'app/views/**/*',
    'app/*.html',
    'app/*.txt',
    'app/*.xml',
  ]
  gulp.watch(srcs, gulp.series('build:html'));
});

// watch:css watches image files for changes and updates them
gulp.task('watch:images', () => {
  gulp.watch('app/images/**/*', gulp.series('build:images'));
});

// watch:images watches js files for changes and re-builds them
gulp.task('watch:js', () => {
  const srcs = [
    'app/js/**/*',
    '!app/js/bundle/**/*',
    'app/scripts/**/*',
  ]
  gulp.watch(srcs, gulp.series('build:js', 'build:html'));
});

// watch starts all watchers
gulp.task('watch', gulp.parallel(
  'watch:css',
  'watch:html',
  'watch:images',
  'watch:js',
));

// serve builds the website, starts the webserver, and watches for changes.
gulp.task('serve', gulp.series(
  'build',
  gulp.parallel(
    'watch',
    () => {
      return gulp.src('build')
        .pipe(webserver(opts.webserver()));
    }
  )
));

// serve:dist serves the built and minified website from dist. It does not
// support live-reloading and should be used to verify final output before
// publishing.
gulp.task('serve:dist', gulp.series('dist', () => {
  return gulp.src('dist')
    .pipe(webserver(opts.webserver()));
}));

//
// Codelabs
//
// codelabs:export exports the codelabs
gulp.task('codelabs:export', (callback) => {
  const source = args.source;

  if (source !== undefined) {
    const sources = Array.isArray(source) ? source : [source];
    claat.run(CODELABS_DIR, 'export', CODELABS_ENVIRONMENT, CODELABS_FORMAT, DEFAULT_GA, sources, callback);
  } else {
    const codelabIds = collectCodelabs().map((c) => { return c.id });
    claat.run(CODELABS_DIR, 'update', CODELABS_ENVIRONMENT, CODELABS_FORMAT, DEFAULT_GA, codelabIds, callback);
  }
});



//
// Helpers
//

// gulpdebug is a helper for debugging gulp pipeline transforms. It logs the
// file path and then invokes the callback.
const gulpdebug = () => {
  return through.obj((file, enc, callback) => {
    console.log(file.path);
    callback(null, file)
  });
}

// run executes the given command with the specified arguments.
const run = (cmd, args, callback) => {
  const proc = spawn(cmd, args, { stdio: 'inherit' });
  proc.on('close', callback);
}

// parseViewMetadata parses the metadata of a single view and returns that metadata
// as JSON.
const parseViewMetadata = (filepath) => {
  if (path.basename(filepath) !== 'view.json') {
    throw new Error(`can only be called on view.json, got: ${filepath}`);
  }

  const dirname = path.dirname(filepath);
  const meta = JSON.parse(fs.readFileSync(filepath));

  meta.id = path.basename(dirname);
  meta.url = viewFilename(meta.id).replace(/\.html$/, '');

  if (meta.sort === undefined) {
    meta.sort = meta.id === 'default' ? 'mainCategory' : 'title';
  }

  if (fs.existsSync(path.resolve(dirname, 'style.css'))) {
    meta.customStyle = `/${meta.id}/style.css`;
  }

  const defMeta = defaultViewMetadata();
  for (let key in defMeta) {
    if (meta[key] === undefined) {
      meta[key] = defMeta[key];
    }
  }

  return meta;
}

// parseCodelabMetadata parses the codelab metadata at the given path.
const parseCodelabMetadata = (filepath) => {
  var meta = JSON.parse(fs.readFileSync(filepath));

  meta.category = meta.category || [];
  if (!Array.isArray(meta.category)) {
    meta.category = [meta.category];
  }

  meta.mainCategory = meta.category[0] || DEFAULT_CATEGORY;
  meta.categoryClass = categoryClass(meta);
  meta.url = path.join(CODELABS_NAMESPACE, meta.id, 'index.html');

  return meta;
}

// defaultViewMetadata returns the default view metadata. This is cached for
// performance.
let _defaultViewMetadata;
const defaultViewMetadata = () => {
  if (!_defaultViewMetadata) {
    _defaultViewMetadata = JSON.parse(fs.readFileSync(DEFAULT_VIEW_META_PATH));
    _defaultViewMetadata.id = 'default';
    _defaultViewMetadata.url = viewFilename('default');
    _defaultViewMetadata.catLevel = 0;
    delete _defaultViewMetadata.ga;
  }
  return _defaultViewMetadata;
}


// Parse view.json and codelab.json files in all directories. Value returned in
// the callback is an object:
//
//     {
//       views: Object.<String, ViewObj>,
//       codelabs: Array.<CodelabObj>,
//       categories: Array.<String>
//     }
//
// Both codelabs and categories are sorted alphabetically. Codelabs are sorted
// by main category.
let _allMetadata;
const collectMetadata = () => {
  if (_allMetadata === undefined) {
    let codelabs = [];
    let categories = {};
    let views = {};

    const viewFiles = glob.sync('app/views/*/view.json');
    for (let i = 0; i < viewFiles.length; i++) {
      const view = parseViewMetadata(viewFiles[i]);
      views[view.id] = view;
    }

    const codelabFiles = glob.sync(`${CODELABS_DIR}/*/codelab.json`);
    for (let i = 0; i < codelabFiles.length; i++) {
      const codelab = parseCodelabMetadata(codelabFiles[i]);
      codelabs.push(codelab);
      categories[codelab.mainCategory] = true;
    }

    _allMetadata = {
      categories: Object.keys(categories).sort(),
      codelabs: codelabs,
      views: views,
    }
  }

  return _allMetadata;
}

// generateDirectoryIndex accepts a stream of HTML files and converts them to
// directory indexes for use on a webserver. For example:
//
//    foo.html => foo/index.html
//
// Then the server can be configured to serve indexes and /foo will render
// properly.
const generateDirectoryIndex = () => {
  return through.obj((file, enc, callback) => {
    const srcPath = file.path;
    if (srcPath.match(/\.html$/) && fs.existsSync(srcPath)) {
      const dstPath = srcPath.replace(/\.html$/, '') + '/index.html';
      const dstPathDir = path.dirname(dstPath);
      const srcPathRel = path.relative(dstPathDir, srcPath);

      // Ensure directory exists
      fs.mkdirpSync(dstPathDir);

      // Change into the directory and create a relative symlink
      chdir(dstPathDir, () => {
        fs.ensureSymlinkSync(srcPathRel, 'index.html');
      });
    }
    callback(null, file);
  });
}

const chdir = (dir, callback) => {
  const cwd = process.cwd();
  process.chdir(dir);
  callback()
  process.chdir(cwd);
}

// generateView creates an index for the given view and codelabs data.
const generateView = () => {
  return through.obj((file, enc, callback) => {
    // viewId is the basename of the dirname of the event (e.g.
    // app/views/my-event/view.json -> my-event).
    const viewId = path.basename(path.dirname(file.path));

    // Read the HTML template for this view (or the default template if no
    // view-specific one exists).
    let templatePath = `app/views/${viewId}/index.html`;
    if (!fs.existsSync(templatePath)) {
      templatePath = DEFAULT_VIEW_TMPL_PATH;
    }
    const template = fs.readFileSync(templatePath);

    // Get the metadata about this view.
    const view = parseViewMetadata(file.path);

    // Aanalytics information.
    const ga = DEFAULT_GA;

    // Full list of views
    const all = collectMetadata();

    // Calculate URL parameters to append.
    let codelabUrlParams = 'index=' + encodeURIComponent('../..' + view.url);
    if (view.ga || args.indexGa) {
      let viewGa = args.indexGa ? args.indexGa : view.ga;
      codelabUrlParams += '&viewga=' + viewGa;
    }

    // Get the list of codelabs and categories for this view
    const filtered = filterCodelabs(view, all.codelabs);
    const codelabs = filtered.codelabs;
    const categories = filtered.categories;

    let locals = {
      baseUrl: BASE_URL,
      categories: categories,
      codelabs: codelabs,
      ga: ga,
      showcats: categories.length > 1,
      view: view,
      views: all.views,

      canonicalViewUrl: viewFuncs.canonicalViewUrl(),
      categoryClass: viewFuncs.categoryClass(),
      categoryHasShowableCodelabs: viewFuncs.categoryHasShowableCodelabs(viewId),
      codePrettyDate: viewFuncs.codelabPrettyDate(),
      codelabPin: viewFuncs.codelabPin(),
      codelabUrl: viewFuncs.codelabUrl(codelabUrlParams),
      levelledCategory: viewFuncs.levelledCategory(),
    };

    const html = swig.render(template.toString(), { locals: locals });
    file.contents = new Buffer.from(html);
    file.base = __dirname;
    file.path = viewFilename(viewId);

    callback(null, file);
  });
}

// viewFuncs is the list of functions shared by views (usually for rendering).
// This is a constant for performance instead of defining the functions inline.
const viewFuncs = {
  // categoryClass returns the top-level categoryClass function.
  categoryClass: () => {
    return categoryClass;
  },

  // canonicalViewUrl returns the canonical URL to the given view.
  canonicalViewUrl: () => {
    return (view) => {
      if (view.id === 'default' || view.id === '') {
        return `${BASE_URL}/`;
      }
      return `${BASE_URL}/${view.id}/`;
    };
  },

  // categoryHasShowableCodelabs returns true if the given category is present
  // in any of the given codelabs, or false otherwise.
  categoryHasShowableCodelabs: (viewId) => {
    return (category, codelabs) => {
      const codelabsUsingCategory = codelabs
        .filter((codelab) => {
          return codelab.category.indexOf(category) !== -1;
        })
        .filter((codelab) => {
          // Rilter hidden codelabs from the default view. All other views are
          // explictly opt-in via metadata.
          return viewId !== 'default' || codelab.status.indexOf('hidden') === -1;
        });

      return codelabsUsingCategory.length > 0;
    }
  },

  // codelabPin returns the index of the pinned item, if it exists.
  codelabPin: () => {
    return (view, codelab) => {
      const i = view.pins.indexOf(codelab.id);
      if (i >= 0) {
        return i + 1;
      }
      return '';
    }
  },

  // codelabUrl returns the url for this codelab, combined with any additional
  // parameters if given.
  codelabUrl: (params) => {
    return (view, codelab) => {
      let url = codelab.url;
      if (params !== undefined) {
        url = `${url}?${params}`;
      }
      if (url.length > 0 && url[0] !== '/') {
        url = `/${url}`;
      }
      return url;
    }
  },

  // codelabPrettyDate returns a pretty-formatted date for the codelab view
  // page.
  codelabPrettyDate: () => {
    return (ts) => {
      const monthNames = ['Jan', 'Feb', 'Mar', 'Apr', 'May', 'Jun', 'Jul', 'Aug', 'Sep', 'Oct', 'Nov', 'Dec'];
      const d = new Date(ts);
      const month = monthNames[d.getMonth()];
      const date = d.getUTCDate();
      const year = d.getFullYear();

      return `${month} ${date}, ${year}`;
    }
  },

  // levelledCategory returns the levelledCategory top-level function.
  levelledCategory: () => {
    return levelledCategory;
  },
};

// viewFilename returns the view index file name based on provided view.
const viewFilename = (view) => {
  if (view === 'default') {
    return 'index.html';
  }
  return `${view}.html`;
}

// levelledCategory returns the codelab category level and name. If no category
// exists at that level, the nearest is returned. If the codelab has no
// categories, the default category and level 0 are returned.
const levelledCategory = (codelab, level) => {
  let i;
  for (i = level; i >= 0; i--) {
    var name = codelab.category[i];
    if (name) {
      break;
    }
  }
  if (!name) {
    name = DEFAULT_CATEGORY;
  }
  return { name: name, level: i };
}

// categoryClass converts the codelab to its corresponding CSS class.
const categoryClass = (codelab, level) => {
  var name = codelab.mainCategory;
  if (level > 0) {
    name += ' ' + codelab.category[level];
  }
  return name.toLowerCase().replace(/\s/g, '-');
}

// Filters out codelabs which do not match view spec. Currently, the matching is
// done by intersecting view.tags with each codelab.tags, and view.categories
// with codelab.category. All codelabs are matched if the view.tags and
// view.categories are empty.
//
// Additionally, each codelab is tested against view.exclude exclusion filters,
// matching codelab ID, status, tags and categories.
const filterCodelabs = (view, codelabs) => {
  var vtags = cleanStringList(view.tags);
  var vcats = cleanStringList(view.categories);

  // Filter out codelabs with tags that don't intersect with view.tags
  // unless view.tags is empty - equivalent to any tag.
  // Same for categories.
  codelabs = codelabs.filter(function(codelab) {
    // Matches by default if both tags and cats are empty.
    var match = !vtags.length && !vcats.length;
    var ctags = cleanStringList(codelab.tags);
    var ccats = cleanStringList(codelab.category);
    // Match by tags.
    if (vtags.length) {
      for (var i = 0; i < ctags.length; i++) {
        if (vtags.indexOf(ctags[i]) > -1) {
          match = true;
          break;
        }
      }
    }
    // Match by category.
    if (!match && vcats.length) {
      for (var i = 0; i < ccats.length; i++) {
        if (vcats.indexOf(ccats[i]) > -1) {
          match = true;
          break;
        }
      }
    }
    if (!match || !view.exclude.length) {
      return match;
    }
    // Tag or category matches. Test the exclusion filter.
    var cstatus = cleanStringList(codelab.status);
    for (var i = 0; i < view.exclude.length; i++) {
      var rex = view.exclude[i];
      if (codelab.id.match(rex)) {
        return false;
      }
      for (var j = 0; j < cstatus.length; j++) {
        if (cstatus[j].match(rex)) {
          return false;
        }
      }
      for (var j = 0; j < ctags.length; j++) {
        if (ctags[j].match(rex)) {
          return false;
        }
      }
      for (var j = 0; j < ccats.length; j++) {
        if (ccats[j].match(rex)) {
          return false;
        }
      }
    }
    return true;
  });

  // Compute distinct categories.
  var categories = {};
  for (var i in codelabs) {
    var cat = levelledCategory(codelabs[i], view.catLevel);
    categories[cat.name] = true;
  }

  // sort the codelabs.
  sortCodelabs(codelabs, view);

  return {
    codelabs: codelabs,
    categories: Object.keys(categories).sort(),
  };
}


// cleanStringList removes empty elements and converts the rest to lower case.
// The a argument is an array of strings or null/undefined. Alternatively, the
// argument can be a string of JSON-serialized array form.
const cleanStringList = (a) => {
  if (typeof a === 'string') {
    // Legacy codelab.json have array field as a string, for instance:
    // "status": "[u'ready for review']".
    a = JSON.parse(a.replace(/u'/, "'").replace(/'/g, '"'));
  }
  a = (a || []).filter((v) => { return (v || '').length > 0; });
  return a.map((v) => { return v.toLowerCase(); });
}

// sortCodelabs reorders codelabs array in-place for the given view. Pinned
// codelabs always come first, sorted by their pin index.
const sortCodelabs = (codelabs, view) => {
  let attr = function(codelab) {
    return levelledCategory(codelab, view.catLevel).name;
  }

  if (view.sort !== "mainCategory") {
    attr = function(codelab) { return codelab[view.sort] };
  }

  codelabs.sort(function(a, b) {
    // Move pinned codelabs to the top.
    var ia = view.pins.indexOf(a.id);
    var ib = view.pins.indexOf(b.id);
    if (ia >= 0 && ib < 0) {
      return -1;
    }
    if (ia < 0 && ib >= 0) {
      return 1;
    }
    if (ia >= 0 && ib >= 0) {
      return ia - ib;
    }
    // Regular sorting.
    var aa = attr(a);
    var ba = attr(b);
    if (aa < ba) { return -1; }
    if (aa > ba) { return 1; }
    return 0;
  });
}

// copyFilteredCodelabs copies all codelabs that match the given id regular
// expression or view regular expression into the build/ folder. If no filters
// are specified (i.e. if codelabRe and viewRe are both undefined), then this
// function returns all codelabs in the codelabs directory.
const copyFilteredCodelabs = (dest) =>  {
  // No filters were specified, symlink the codelabs folder directly and save
  // processing.
  if (CODELABS_FILTER === '*' && VIEWS_FILTER === '*') {
    const source = path.join(CODELABS_DIR);
    const target = path.join(dest, CODELABS_NAMESPACE);
    fs.ensureSymlinkSync(source, target, 'dir');
    return
  }

  const codelabs = collectCodelabs();

  for(let i = 0; i < codelabs.length; i++) {
    const codelab = codelabs[i];
    const source = path.join(CODELABS_DIR, codelab.id);
    const target = path.join(dest, CODELABS_NAMESPACE, codelab.id);
    fs.ensureSymlinkSync(source, target, 'dir');
  }
};

// collectCodelabs collects the list of codelabs that match the given view or
// codelab filter.
const collectCodelabs = () => {
  const meta = collectMetadata();
  let codelabs = meta.codelabs;

  // Only select codelabs that match the given codelab ID.
  if (CODELABS_FILTER !== '*') {
    codelabs = meta.codelabs.filter((codelab) => {
      return codelab.id.match(CODELABS_FILTER);
    });

    if (codelabs.length === 0) {
      throw new Error(`no codelabs matched: ${CODELABS_FILTER}`);
    }
  }

  // Only select codelabs that match the given view ID.
  if (VIEWS_FILTER !== '*') {
    let views = [];
    Object.keys(meta.views).forEach((key) => {
      if (key.match(VIEWS_FILTER)) {
        views.push(meta.views[key]);
      }
    });

    if (views.length === 0) {
      throw new Error(`no views matched: ${VIEWS_FILTER}`);
    }

    // Iterate over each view and include codelabs for that view.
    let s = new Set();
    for(let i = 0; i < views.length; i++) {
      let filtered = filterCodelabs(views[i], codelabs).codelabs;
      s.add(...filtered);
    }
    codelabs = s.values();
  }

  // Check if we have any codelabs
  if (codelabs.length === 0) {
    throw new Error('no codelabs matched given filters');
  }

  return codelabs;
}

// publish:staging:codelabs uploads the dist folder codelabs to a staging
// bucket. This only uploads the codelabs, the views remain unchanged.
gulp.task('publish:staging:codelabs', (callback) => {
  const opts = { dry: DRY_RUN, deleteMissing: DELETE_MISSING };
  const src = path.join('dist', CODELABS_NAMESPACE, '/');
  const dest = gcs.bucketFolderPath(STAGING_BUCKET, CODELABS_NAMESPACE);
  gcs.rsync(src, dest, opts, callback);
});

// publish:staging:views uploads the dist folder views and associated assets to
// a staging bucket. This does not upload any of the codelabs.
gulp.task('publish:staging:views', (callback) => {
  const opts = { exclude: CODELABS_NAMESPACE, dry: DRY_RUN, deleteMissing: DELETE_MISSING };
  gcs.rsync('dist', STAGING_BUCKET, opts, callback);
});

// publish:prod:codelabs syncs codelabs from the staging to the production
// bucket.
gulp.task('publish:prod:codelabs', (callback) => {
  const opts = { dry: DRY_RUN, deleteMissing: DELETE_MISSING };
  const src = gcs.bucketFolderPath(STAGING_BUCKET, CODELABS_NAMESPACE);
  const dest = gcs.bucketFolderPath(PROD_BUCKET, CODELABS_NAMESPACE);
  gcs.rsync(src, dest, opts, callback);
});

// publish:prod:views syncs views and associated assets from the staging to the
// production bucket.
gulp.task('publish:prod:views', (callback) => {
  const opts = { exclude: CODELABS_NAMESPACE, dry: DRY_RUN, deleteMissing: DELETE_MISSING };
  gcs.rsync(STAGING_BUCKET, PROD_BUCKET, opts, callback);
});
