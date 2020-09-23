
package render

var entities = map[string]string{
	"&amp;lt;": "&lt;",
	"&amp;apos;": "&apos;",
	"&amp;gt;": "&gt;",
	"&amp;nbsp;": "&nbsp;",
	"&amp;iexcl;": "&iexcl;",
	"&amp;cent;": "&cent;",
	"&amp;pound;": "&pound;",
	"&amp;curren;": "&curren;",
	"&amp;yen;": "&yen;",
	"&amp;brvbar;": "&brvbar;",
	"&amp;sect;": "&sect;",
	"&amp;uml;": "&uml;",
	"&amp;copy;": "&copy;",
	"&amp;ordf;": "&ordf;",
	"&amp;laquo;": "&laquo;",
	"&amp;not;": "&not;",
	"&amp;reg;": "&reg;",
	"&amp;macr;": "&macr;",
	"&amp;deg;": "&deg;",
	"&amp;plusmn;": "&plusmn;",
	"&amp;sup2;": "&sup2;",
	"&amp;sup3;": "&sup3;",
	"&amp;acute;": "&acute;",
	"&amp;micro;": "&micro;",
	"&amp;para;": "&para;",
	"&amp;middot;": "&middot;",
	"&amp;cedil;": "&cedil;",
	"&amp;sup1;": "&sup1;",
	"&amp;ordm;": "&ordm;",
	"&amp;raquo;": "&raquo;",
	"&amp;frac14;": "&frac14;",
	"&amp;frac12;": "&frac12;",
	"&amp;frac34;": "&frac34;",
	"&amp;iquest;": "&iquest;",
	"&amp;Agrave;": "&Agrave;",
	"&amp;Aacute;": "&Aacute;",
	"&amp;Acirc;": "&Acirc;",
	"&amp;Atilde;": "&Atilde;",
	"&amp;Auml;": "&Auml;",
	"&amp;Aring;": "&Aring;",
	"&amp;AElig;": "&AElig;",
	"&amp;Ccedil;": "&Ccedil;",
	"&amp;Egrave;": "&Egrave;",
	"&amp;Eacute;": "&Eacute;",
	"&amp;Ecirc;": "&Ecirc;",
	"&amp;Euml;": "&Euml;",
	"&amp;Igrave;": "&Igrave;",
	"&amp;Iacute;": "&Iacute;",
	"&amp;Icirc;": "&Icirc;",
	"&amp;Iuml;": "&Iuml;",
	"&amp;ETH;": "&ETH;",
	"&amp;Ntilde;": "&Ntilde;",
	"&amp;Ograve;": "&Ograve;",
	"&amp;Oacute;": "&Oacute;",
	"&amp;Ocirc;": "&Ocirc;",
	"&amp;Otilde;": "&Otilde;",
	"&amp;Ouml;": "&Ouml;",
	"&amp;times;": "&times;",
	"&amp;Oslash;": "&Oslash;",
	"&amp;Ugrave;": "&Ugrave;",
	"&amp;Uacute;": "&Uacute;",
	"&amp;Ucirc;": "&Ucirc;",
	"&amp;Uuml;": "&Uuml;",
	"&amp;Yacute;": "&Yacute;",
	"&amp;THORN;": "&THORN;",
	"&amp;szlig;": "&szlig;",
	"&amp;agrave;": "&agrave;",
	"&amp;aacute;": "&aacute;",
	"&amp;acirc;": "&acirc;",
	"&amp;atilde;": "&atilde;",
	"&amp;auml;": "&auml;",
	"&amp;aring;": "&aring;",
	"&amp;aelig;": "&aelig;",
	"&amp;ccedil;": "&ccedil;",
	"&amp;egrave;": "&egrave;",
	"&amp;eacute;": "&eacute;",
	"&amp;ecirc;": "&ecirc;",
	"&amp;euml;": "&euml;",
	"&amp;igrave;": "&igrave;",
	"&amp;iacute;": "&iacute;",
	"&amp;icirc;": "&icirc;",
	"&amp;iuml;": "&iuml;",
	"&amp;eth;": "&eth;",
	"&amp;ntilde;": "&ntilde;",
	"&amp;ograve;": "&ograve;",
	"&amp;oacute;": "&oacute;",
	"&amp;ocirc;": "&ocirc;",
	"&amp;otilde;": "&otilde;",
	"&amp;ouml;": "&ouml;",
	"&amp;divide;": "&divide;",
	"&amp;oslash;": "&oslash;",
	"&amp;ugrave;": "&ugrave;",
	"&amp;uacute;": "&uacute;",
	"&amp;ucirc;": "&ucirc;",
	"&amp;uuml;": "&uuml;",
	"&amp;yacute;": "&yacute;",
	"&amp;thorn;": "&thorn;",
	"&amp;yuml;": "&yuml;",
	"&amp;OElig;": "&OElig;",
	"&amp;oelig;": "&oelig;",
	"&amp;Scaron;": "&Scaron;",
	"&amp;scaron;": "&scaron;",
	"&amp;Yuml;": "&Yuml;",
	"&amp;fnof;": "&fnof;",
	"&amp;circ;": "&circ;",
	"&amp;tilde;": "&tilde;",
	"&amp;Alpha;": "&Alpha;",
	"&amp;Beta;": "&Beta;",
	"&amp;Gamma;": "&Gamma;",
	"&amp;Delta;": "&Delta;",
	"&amp;Epsilon;": "&Epsilon;",
	"&amp;Zeta;": "&Zeta;",
	"&amp;Eta;": "&Eta;",
	"&amp;Theta;": "&Theta;",
	"&amp;Iota;": "&Iota;",
	"&amp;Kappa;": "&Kappa;",
	"&amp;Lambda;": "&Lambda;",
	"&amp;Mu;": "&Mu;",
	"&amp;Nu;": "&Nu;",
	"&amp;Xi;": "&Xi;",
	"&amp;Omicron;": "&Omicron;",
	"&amp;Pi;": "&Pi;",
	"&amp;Rho;": "&Rho;",
	"&amp;Sigma;": "&Sigma;",
	"&amp;Tau;": "&Tau;",
	"&amp;Upsilon;": "&Upsilon;",
	"&amp;Phi;": "&Phi;",
	"&amp;Chi;": "&Chi;",
	"&amp;Psi;": "&Psi;",
	"&amp;Omega;": "&Omega;",
	"&amp;alpha;": "&alpha;",
	"&amp;beta;": "&beta;",
	"&amp;gamma;": "&gamma;",
	"&amp;delta;": "&delta;",
	"&amp;epsilon;": "&epsilon;",
	"&amp;zeta;": "&zeta;",
	"&amp;eta;": "&eta;",
	"&amp;theta;": "&theta;",
	"&amp;iota;": "&iota;",
	"&amp;kappa;": "&kappa;",
	"&amp;lambda;": "&lambda;",
	"&amp;mu;": "&mu;",
	"&amp;nu;": "&nu;",
	"&amp;xi;": "&xi;",
	"&amp;omicron;": "&omicron;",
	"&amp;pi;": "&pi;",
	"&amp;rho;": "&rho;",
	"&amp;sigmaf;": "&sigmaf;",
	"&amp;sigma;": "&sigma;",
	"&amp;tau;": "&tau;",
	"&amp;upsilon;": "&upsilon;",
	"&amp;phi;": "&phi;",
	"&amp;chi;": "&chi;",
	"&amp;psi;": "&psi;",
	"&amp;omega;": "&omega;",
	"&amp;thetasym;": "&thetasym;",
	"&amp;Upsih;": "&Upsih;",
	"&amp;piv;": "&piv;",
	"&amp;ndash;": "&ndash;",
	"&amp;mdash;": "&mdash;",
	"&amp;lsquo;": "&lsquo;",
	"&amp;rsquo;": "&rsquo;",
	"&amp;sbquo;": "&sbquo;",
	"&amp;ldquo;": "&ldquo;",
	"&amp;rdquo;": "&rdquo;",
	"&amp;bdquo;": "&bdquo;",
	"&amp;dagger;": "&dagger;",
	"&amp;Dagger;": "&Dagger;",
	"&amp;bull;": "&bull;",
	"&amp;hellip;": "&hellip;",
	"&amp;permil;": "&permil;",
	"&amp;prime;": "&prime;",
	"&amp;Prime;": "&Prime;",
	"&amp;lsaquo;": "&lsaquo;",
	"&amp;rsaquo;": "&rsaquo;",
	"&amp;oline;": "&oline;",
	"&amp;frasl;": "&frasl;",
	"&amp;euro;": "&euro;",
	"&amp;image;": "&image;",
	"&amp;weierp;": "&weierp;",
	"&amp;real;": "&real;",
	"&amp;trade;": "&trade;",
	"&amp;alefsym;": "&alefsym;",
	"&amp;larr;": "&larr;",
	"&amp;uarr;": "&uarr;",
	"&amp;rarr;": "&rarr;",
	"&amp;darr;": "&darr;",
	"&amp;harr;": "&harr;",
	"&amp;crarr;": "&crarr;",
	"&amp;lArr;": "&lArr;",
	"&amp;UArr;": "&UArr;",
	"&amp;rArr;": "&rArr;",
	"&amp;dArr;": "&dArr;",
	"&amp;hArr;": "&hArr;",
	"&amp;forall;": "&forall;",
	"&amp;part;": "&part;",
	"&amp;exist;": "&exist;",
	"&amp;empty;": "&empty;",
	"&amp;nabla;": "&nabla;",
	"&amp;isin;": "&isin;",
	"&amp;notin;": "&notin;",
	"&amp;ni;": "&ni;",
	"&amp;prod;": "&prod;",
	"&amp;sum;": "&sum;",
	"&amp;minus;": "&minus;",
	"&amp;lowast;": "&lowast;",
	"&amp;radic;": "&radic;",
	"&amp;prop;": "&prop;",
	"&amp;infin;": "&infin;",
	"&amp;ang;": "&ang;",
	"&amp;and;": "&and;",
	"&amp;or;": "&or;",
	"&amp;cap;": "&cap;",
	"&amp;cup;": "&cup;",
	"&amp;int;": "&int;",
	"&amp;there4;": "&there4;",
	"&amp;sim;": "&sim;",
	"&amp;cong;": "&cong;",
	"&amp;asymp;": "&asymp;",
	"&amp;ne;": "&ne;",
	"&amp;equiv;": "&equiv;",
	"&amp;le;": "&le;",
	"&amp;ge;": "&ge;",
	"&amp;sub;": "&sub;",
	"&amp;sup;": "&sup;",
	"&amp;nsub;": "&nsub;",
	"&amp;sube;": "&sube;",
	"&amp;supe;": "&supe;",
	"&amp;oplus;": "&oplus;",
	"&amp;otimes;": "&otimes;",
	"&amp;perp;": "&perp;",
	"&amp;sdot;": "&sdot;",
	"&amp;lceil;": "&lceil;",
	"&amp;rceil;": "&rceil;",
	"&amp;lfloor;": "&lfloor;",
	"&amp;rfloor;": "&rfloor;",
	"&amp;lang;": "&lang;",
	"&amp;rang;": "&rang;",
	"&amp;loz;": "&loz;",
	"&amp;spades;": "&spades;",
	"&amp;clubs;": "&clubs;",
	"&amp;hearts;": "&hearts;",
	"&amp;diams;": "&diams;",
}