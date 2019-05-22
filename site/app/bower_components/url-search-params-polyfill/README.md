url-search-params
=================

[![build status](https://secure.travis-ci.org/WebReflection/url-search-params.svg)](http://travis-ci.org/WebReflection/url-search-params)

This is a polyfill for the [URLSearchParams API](https://developer.mozilla.org/en-US/docs/Web/API/URLSearchParams).

It is possible to simply include [build/url-search-params.js](build/url-search-params.js) or grab it via npm.

```
npm install url-search-params
```

The function is exported directly.
```js
var URLSearchParams = require('url-search-params');
```

MIT Style License

#### About HTMLAnchorElement.prototype.searchParams
This property is already implemented in Firefox and polyfilled here only for browsers that exposes getters and setters
through the `HTMLAnchorElement.prototype`.

In order to know if such property is supported, you **must** do the check as such:
```
if ('searchParams' in HTMLAnchorElement.prototype) {
  // polyfill for <a> links supported
}
```
If you do this check instead:
```js
if (HTMLAnchorElement.prototype.searchParams) {
  // throws a TypeError
}
```
this polyfill will reflect native behavior, throwing a type error due access to a property in a non instance of `HTMLAnchorElement`.

Nothing new to learn here, [just a reminder](http://webreflection.blogspot.co.uk/2011/08/please-stop-reassigning-for-no-reason.html).
