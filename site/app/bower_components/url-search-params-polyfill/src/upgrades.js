
upgradeClass(HTMLAnchorElement);
if (/^function|object$/.test(typeof URL)) upgradeClass(URL);

/*

function spUpdate(a) {
  var
    search = this.toString(),
    password = a.password,
    username = a.username
  ;
  a.href = ''.concat(
    a.protocol, '//',
    username,
    password ? (':' + password) : '',
    username ? '@' : '',
    a.host,
    a.pathname,
    search ? ('?' + search) : '',
    a.hash
  );
}

*/
