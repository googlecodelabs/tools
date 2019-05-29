'use strict';

const childprocess = require('child_process');
const path = require('path');
const spawn = childprocess.spawn;

// GZIP_TYPES is the list of file types that can be gzipped.
const GZIP_TYPES = 'css,html,ico,js,json,svg,xml';

// ALWAYS_EXCLUDE is the list of filters to always use for exclusion.
const ALWAYS_EXCLUDE = '\\.git.*|\\.DS_Store';

// HEADER_CACHE_CONTROL is the default cache control header.
const HEADER_CACHE_CONTROL = 'Cache-Control:public,max-age=3600';

// rsync syncs files from the given src to dest
exports.rsync = (src, dest, opts = {}, callback) => {
  const args = [];

  if (opts && opts.exclude) {
    args.push('-x', `${opts.exclude}|${ALWAYS_EXCLUDE}`);
  } else {
    args.push('-x', ALWAYS_EXCLUDE);
  }

  if (opts && opts.dry) {
    args.push('-n');
  }

  if (opts && opts.deleteMissing) {
    args.push('-d');
  }

  const proc = spawn('gsutil', [
    '-o', 'GSUtil:parallel_process_count=8',
    '-o', 'GSUtil:parallel_thread_count=1',
    '-h', HEADER_CACHE_CONTROL,
    '-m',
    'rsync',
      '-c', // compare checksums instead of mtime
      '-C', // continue on error (but exit non-zero at end)
      '-j', GZIP_TYPES, // gzip
      '-r', // recurse
      ...args,
      src, dest
  ], { stdio: 'inherit' });

  proc.on('close', (err) => {
    if (err) {
      throw new Error(err);
    }
    callback();
  });
};

// bucketName returns the normalized bucket arg or defaultValue if bucket is
// blank:
//
// - prepend gs:// if missing
// - remove '/' suffix
exports.bucketName = (bucket, defaultValue) => {
  bucket = bucket || defaultValue;

  if (bucket.substring(0, 5) !== 'gs://') {
    bucket = 'gs://' + bucket;
  }

  bucket = removeTrailingSlash(bucket);

  return bucket;
};

// bucketFilePath is the path in the bucket to the file specified at the given
// path. The path is normalized.
//
// This function exists because path.join() removes the protocol from gs:// and
// doesn't give much control over the trailing slash.
exports.bucketFilePath = (bucket, ...paths) => {
  bucket = removeTrailingSlash(bucket);
  paths = paths.map((p) => {
    return removeTrailingSlash(removeLeadingSlash(p));
  });
  return bucket + '/' + path.join(...paths);
};

// bucketFolderPath is the path in the bucket to the folder specified by the
// given paths.
exports.bucketFolderPath = (bucket, ...paths) => {
  return exports.bucketFilePath(bucket, ...paths) + '/';
};

// removeLeadingSlash removes the trailing slash from a string.
const removeLeadingSlash = (s) => {
  if (s[0] === '/') {
    return s.substring(1, s.length);
  }
  return s;
};

// removeTrailingSlash removes the trailing slash from a string.
const removeTrailingSlash = (s) => {
  if (s[s.length - 1] === '/') {
    return s.substring(0, s.length - 1);
  }
  return s;
};
