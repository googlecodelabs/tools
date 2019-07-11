
<!---

This README is automatically generated from the comments in these files:
iron-jsonp-library.html

Edit those files, and our readme bot will duplicate them over here!
Edit this file, and the bot will squash your changes :)

-->

[![Build Status](https://travis-ci.org/PolymerElements/iron-jsonp-library.svg?branch=master)](https://travis-ci.org/PolymerElements/iron-jsonp-library)

_[Demo and API Docs](https://elements.polymer-project.org/elements/iron-jsonp-library)_


##&lt;iron-jsonp-library&gt;


Loads specified jsonp library.

Example:

    <iron-jsonp-library
      library-url="https://apis.google.com/js/plusone.js?onload=%%callback%%"
      notify-event="api-load"
      library-loaded="{{loaded}}"></iron-jsonp-library>

Will emit 'api-load' event when loaded, and set 'loaded' to true

Implemented by  Polymer.IronJsonpLibraryBehavior. Use it
to create specific library loader elements.



##Polymer.IronJsonpLibraryBehavior


    `Polymer.IronJsonpLibraryBehavior` loads a jsonp library.
    Multiple components can request same library, only one copy will load.

    Some libraries require a specific global function be defined.
    If this is the case, specify the `callbackName` property.

    You should use an HTML Import to load library dependencies
    when possible instead of using this element.

    
