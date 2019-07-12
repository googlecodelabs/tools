/*
Note: <body is="body-bind"> is an experiment taken from
https://github.com/PolymerLabs/polymer-experiments/tree/master/body-bind.
It allows bindings in <body>, and other Polymer features which are lazily
upgraded. This approach does not block rendering!.
*/

(function() {

  var origTryReady = Polymer.Base._tryReady;

  Polymer.Base._tryReady = function() {
    // Configure properties bound from host lazily at startup
    // This is opt-in because it is expensive; it either needs to be in
    // a permanent top-level host with `lazyChildren` property set, or else
    // each instance needs to set the `lazy-bind` attribute
    if (this.is && this.hasAttribute('lazy-bind') ||
        (this.dataHost && this.dataHost.lazyChildren)) {
      for (var p in this._propertyEffects) {
        if (this.hasOwnProperty(p)) {
          var v = this[p];
          delete this[p];
          this._config[p] = v;
        }
      }
    }
    origTryReady.apply(this, arguments);
  };

  Polymer.LightDomBindingBehavior = {

    lazyChildren: true,

    _initFeatures: function() {
      // By default, LightDomBindingBehavior becomes the permanent top-level
      // host, so that any lazy-loaded children are automatically hosted by
      // this element, and by virtue of `dataHost.lazyChildren: true` will
      // lazily configure any properties bound from the host.
      this._beginHosting();
      // Instance time binding of light children
      this.root = this;
      this._template = this;
      this._content = this;
      this._notes = null;
      this._prepAnnotations();
      this._prepEffects();
      this._prepBindings();
      this._setupConfigure();
      this._marshalAnnotationReferences();
      this._marshalInstanceEffects();
      this._tryReady();
    },
  };

}());

// Similar to `<template is="dom-bind">`, but uses the body's light DOM
// as the template.  If more advanced features such as event listeners,
// computed functions, etc. are needed, just make an `extends: 'body'`
// element that uses the `Polymer.LightDomBindingBehavior`, and use
// normal Polymer idioms.

Polymer({

  extends: 'body',

  is: 'body-bind',

  behaviors: [Polymer.LightDomBindingBehavior]

});
