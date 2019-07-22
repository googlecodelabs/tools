'use strict';

const childprocess = require('child_process');
const spawn = childprocess.spawn;

// claat is a wrapper around the claat tool.
//
//   cwd - codelabs content dir
//   cmd - claat command, either 'update' or 'export'
//   fmt - output format, e.g. 'html'
//   ga - google analytics tracking code
//   args - an array of source doc IDs or codelab names (IDs)
//   callback - an async task callback function
//
exports.run = (cwd, cmd, env, fmt, ga, args, callback) => {
  args.unshift(cmd, '-e', env, '-f', fmt, '-ga', ga);
  const proc = spawn('claat', args, { stdio: 'inherit', cwd: cwd, env: process.env });

  proc.on('close', (e) => {
    if (e) {
      throw new Error(e);
    }
    callback();
  })
};
