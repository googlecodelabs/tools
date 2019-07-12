var
  dP = Object.defineProperty,
  gOPD = Object.getOwnPropertyDescriptor,
  createSearchParamsPollute = function (search) {
    /*jshint validthis:true */
    function append(name, value) {
      URLSearchParamsProto.append.call(this, name, value);
      name = this.toString();
      search.set.call(this._usp, name ? ('?' + name) : '');
    }
    function del(name) {
      URLSearchParamsProto.delete.call(this, name);
      name = this.toString();
      search.set.call(this._usp, name ? ('?' + name) : '');
    }
    function set(name, value) {
      URLSearchParamsProto.set.call(this, name, value);
      name = this.toString();
      search.set.call(this._usp, name ? ('?' + name) : '');
    }
    return function (sp, value) {
      sp.append = append;
      sp.delete = del;
      sp.set = set;
      return dP(sp, '_usp', {
        configurable: true,
        writable: true,
        value: value
      });
    };
  },
  createSearchParamsCreate = function (polluteSearchParams) {
    return function (obj, sp) {
      dP(
        obj, '_searchParams', {
          configurable: true,
          writable: true,
          value: polluteSearchParams(sp, obj)
        }
      );
      return sp;
    };
  },
  updateSearchParams = function (sp) {
    var append = sp.append;
    sp.append = URLSearchParamsProto.append;
    URLSearchParams.call(sp, sp._usp.search.slice(1));
    sp.append = append;
  },
  verifySearchParams = function (obj, Class) {
    if (!(obj instanceof Class)) throw new TypeError(
      "'searchParams' accessed on an object that " +
      "does not implement interface " + Class.name
    );
  },
  upgradeClass = function(Class){
    var
      ClassProto = Class.prototype,
      searchParams = gOPD(ClassProto, 'searchParams'),
      href = gOPD(ClassProto, 'href'),
      search = gOPD(ClassProto, 'search'),
      createSearchParams
    ;
    if (!searchParams && search && search.set) {
      createSearchParams = createSearchParamsCreate(
        createSearchParamsPollute(search)
      );
      Object.defineProperties(
        ClassProto,
        {
          href: {
            get: function () {
              return href.get.call(this);
            },
            set: function (value) {
              var sp = this._searchParams;
              href.set.call(this, value);
              if (sp) updateSearchParams(sp);
            }
          },
          search: {
            get: function () {
              return search.get.call(this);
            },
            set: function (value) {
              var sp = this._searchParams;
              search.set.call(this, value);
              if (sp) updateSearchParams(sp);
            }
          },
          searchParams: {
            get: function () {
              verifySearchParams(this, Class);
              return this._searchParams || createSearchParams(
                this,
                new URLSearchParams(this.search.slice(1))
              );
            },
            set: function (sp) {
              verifySearchParams(this, Class);
              createSearchParams(this, sp);
            }
          }
        }
      );
    }

  }
;