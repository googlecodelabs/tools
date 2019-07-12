
<!---

This README is automatically generated from the comments in these files:
iron-localstorage.html

Edit those files, and our readme bot will duplicate them over here!
Edit this file, and the bot will squash your changes :)

The bot does some handling of markdown. Please file a bug if it does the wrong
thing! https://github.com/PolymerLabs/tedium/issues

-->

[![Build status](https://travis-ci.org/PolymerElements/iron-localstorage.svg?branch=master)](https://travis-ci.org/PolymerElements/iron-localstorage)

_[Demo and API docs](https://elements.polymer-project.org/elements/iron-localstorage)_


##&lt;iron-localstorage&gt;

Element access to Web Storage API (window.localStorage).

Keeps `value` property in sync with localStorage.

Value is saved as json by default.

### Usage:

`ls-sample` will automatically save changes to its value.

```html
<dom-module id="ls-sample">
  <iron-localstorage name="my-app-storage"
    value="{{cartoon}}"
    on-iron-localstorage-load-empty="initializeDefaultCartoon"
  ></iron-localstorage>
</dom-module>

<script>
  Polymer({
    is: 'ls-sample',
    properties: {
      cartoon: {
        type: Object
      }
    },
    // initializes default if nothing has been stored
    initializeDefaultCartoon: function() {
      this.cartoon = {
        name: "Mickey",
        hasEars: true
      }
    },
    // use path set api to propagate changes to localstorage
    makeModifications: function() {
      this.set('cartoon.name', "Minions");
      this.set('cartoon.hasEars', false);
    }
  });
</script>
```

### Tech notes:

* `value.*` is observed, and saved on modifications. You must use
  path change notification methods such as `set()` to modify value
  for changes to be observed.


* Set `auto-save-disabled` to prevent automatic saving.


* Value is saved as JSON by default.


* To delete a key, set value to null



Element listens to StorageAPI `storage` event, and will reload upon receiving it.

__Warning__: do not bind value to sub-properties until Polymer
[bug 1550](https://github.com/Polymer/polymer/issues/1550)
is resolved. Local storage will be blown away.
`<iron-localstorage value="{{foo.bar}}"` will cause __data loss__.


