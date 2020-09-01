# Tools for authoring and serving codelabs

[![Demo](https://storage.googleapis.com/claat/demo.png)](https://storage.googleapis.com/claat/demo.mp4)

Codelabs are interactive instructional tutorials, which can be authored in Google Docs
using some simple formatting conventions. You can also author codelabs using markdown syntax.
This repo contains all the tools and documentation you’ll need for building and publishing
your own codelabs.

If you're interested in authoring codelabs, create a doc following the [Codelab Formatting Guide](FORMAT-GUIDE.md). 
and see the [claat](claat) directory for a detailed description of the `claat` command line tool.

Also, consider joining the [codelab-authors Google Group](https://groups.google.com/forum/#!forum/codelab-authors),
which connects you with other authors and provides updates on new releases. 

## What is this?

For the past 3+ years, the CLaaT (Codelabs as a Thing) project has given developers around the
world a hands-on experience with Google products and tools.  We’ve accumulated over 500 high quality
codelabs, served millions of web visitors, and supported over 100 events, from local meetups
all the way up to Google I/O.

This project has been implemented as a volunteer project by a small group of dedicated Googlers
who care deeply about this kind of “learning by doing” approach to education.

## What's special about this tool?

* Powerful and flexible authoring flow via Google Docs
* Optional support for authoring in Markdown text
* Ability to produce interactive web or markdown tutorials without writing any code
* Easy interactive previewing
* Usage monitoring via Google Analytics
* Support for multiple target environments (kiosk, web, markdown, offline, etc.)
* Support for anonymous use - ideal for public computers at developer events
* Looks great, with a responsive web implementation
* Remembers where the student left off when returning to a codelab
* Mobile friendly user experience

## Can I use this to create my own codelabs and serve my own codelabs online?

Yes, the claat tool and the serving mechanism can be used by anyone to author their
own codelabs and to serve up their own codelabs on the web.

You can also use this tool to create a nice looking summary page like the one you see on the official [Google Codelabs site](https://g.co/codelabs).

If you're interested in authoring codelabs, join [codelab-authors group](https://groups.google.com/forum/#!forum/codelab-authors),
which connects you with other authors and provides access to the
[Codelab Formatting Guide](FORMAT-GUIDE.md).

## Ok, how do I use it?

Check out this [excellent tutorial](https://medium.com/@zarinlo/publish-technical-tutorials-in-google-codelab-format-b07ef76972cd).

1. Create a doc following the syntax conventions described in the [Codelab Formatting Guide](FORMAT-GUIDE.md). Here’s an [example doc](https://docs.google.com/document/d/1E6XMcdTexh5O8JwGy42SY3Ehzi8gOfUGiqTiUX6N04o/). Feel free to copy that doc as a starter template. Once you have your own source doc, note its DocId, which is the long string near the end of the URL (right after docs.google.com/document/d/).

1. Make one or more changes and preview your codelab, using the preview app provided by Google. To preview a codelab, install the [Preview Codelab Chrome extension](https://chrome.google.com/webstore/detail/preview-codelab/lhojjnijnkiglhkggagbapfonpdlinji) in your browser. Now you can preview a codelab directly from the Google Doc view by clicking the Chrome extension’s button, which will open a new tab to display the preview. Alternatively, navigate manually to https://codelabs-preview.appspot.com/?file_id=<google-doc-id>

1. Install the claat command -- see the [README in the claat directory](https://github.com/googlecodelabs/tools/blob/master/claat/README.md) of this repo for instructions..

1. Run the claat command to transform the doc contents into one of the supported output formats. The default supported formats are html and markdown but the claat tool supports adding additional formats by specifying a path to a Go template. For example, using the example document above:

        $ claat export 1rpHleSSeY-MJZ8JvncvYA8CFqlnlcrW8-a4uEaqizPY  
        ok      your-first-pwapp

    You can also specify a markdown document (.md file) as input. It has to adhere to the syntax conventions described [here](https://github.com/googlecodelabs/tools/tree/master/claat/parser/md)

        $ claat export document.md
        ok      your-md-based-codelab

1. Run the claat serve command.

        $ claat serve

This will start a local web server and open a browser tab to the local server. Click on the
hyperlink represent your codelab of interest to experience a fully rendered version.

## How do I generate my own landing page?

See instructions in the [site directory's readme](site/README.md).

## How do I generate a custom view?

Copy the [sample view](site/app/views/vslive), customize it to your liking,
tag and rebuild the codelabs you want included, and then generate your view.

## How do I publish my codelabs?

The output generated by `claat` is a purely static set of HTML or Markdown code. As such,
it can be served by any web serving mechanism, including any of the following options:

* Github Pages (`*.github.io`)
* [Google App Engine](https://cloud.google.com/appengine)
* [Firebase Static Serving](https://firebase.google.com/products/hosting)
* [Google Cloud Storage](https://cloud.google.com/storage)
* Amazon Web Services S3
* Netlify
* Any open source web server (Nginx, Apache)
* `python -m SimpleHTTPServer` (Python 2)
* `python3 -m http.server` (Python 3)

Simply commit the artifacts generated by the claat command into your preferred serving vehicle
and you should be ready to go.

The [site directory](site) contains tools for building your own custom landing page(s) and publishing both landing
pages and codelabs to Google Cloud Storage.

## Why bother with this approach when I can write tutorials directly in Markdown?

Some people like the Google Docs authoring flow, others prefer to specify their codelabs
directly in Markdown. Using the Docs approach, one source format can be used to generate
numerous output formats. Also, you can use a doc for the initial formulation stage, where
WYSIWYG and easy collaboration are extremely useful. Once the content stabilizes, typically
after the first launch, you are free to make the generated markdown your source of truth
and discard the Google Doc as a controlling source. This is desirable because it gives you
the ability to manage the content as code in a source control system, but it comes at the
cost of having to commit to one specific output format, or having to maintain multiple
sources of truth.

This tool and corresponding authoring approach are agnostic with respect to whether (and when)
you choose to manage your source as a Google Doc or as Markdown text checked into a repo.
The only hard and fast rule is that, at any one point in time, you should choose one or the
other. Trying to simultaneously maintain a doc and a corresponding repository is a recipe
for disaster.

## What are the supported input formats?

* Google Docs (following FORMAT-GUIDE.md formatting conventions)
* Markdown

## What are the supported output formats?

* Google Codelabs - HTML and Markdown
* Qwiklabs - Markdown
* Cloud Shell Tutorials - Markdown with special directives
* Jupyter, Kaggle Kernels, Colaboratory, et. al. - Markdown with format specific cells

There’s no one “best” publication format. Each format has its own advantages,
disadvantages, community, and application domain. For example, Jupyter has a very strong
following in the data science and Python communities.

This variety of formats is healthy because we’re seeing new innovative approaches all the
time (for example, see observablehq.com, which recently launched their Beta release).

While this evolving format ecosystem is generally a good thing, having to maintain tutorials in
multiple formats, or switch from one format to another can be painful. The Codelabs doc format
(as specified in FORMAT-GUIDE.md) can provide a high level specification for maintaining
a single source of truth, programmatically translated into one or more tutorial specific formats.

## Can I contribute?

Yes, by all means. Have feature ideas? Send us a pull request or file a bug.

## Where did this come from?

For several years, Googlers would rush to build new tutorials and related assets for our
annual developer event, Google I/O. But every year the authoring platform and distribution
mechanism changed. As a result, there was little reuse of content and serving infrastructure,
And every year we essentially kept reinventing the same wheel.

For Google I/O 2014, Shawn Simister wrote a Python program which retrieved
specially formatted documents from Google Drive, parsed them, and generated
a nice interactive web-based user experience. This allowed authors to design their
codelabs using Google Docs, with its great interactivity and collaboration features,
and automatically convert those documents into beautiful web based tutorials,
without needing to write a single line of code.

Later, Ewa Gasperowicz wrote a site generator, supporting the ability to
publish custom landing pages, with associated branding and an inventory of codelabs
specially curated for a given event.

Alex Vaghin later rewrote Shawn's Python program as a statically linked Go program (the claat command in this repo), eliminating many runtime dependencies, improving translation
performance. Alex also added, among many other enhancements, a proper abstract syntax
tree (to facilitate translation to different output formats), an app engine based previewer, an extensible rendering engine, support for generating markdown output. Alex also wrote the web serving infrastructure, the build tooling (based on gulp), and, with the author, the ability to self-publish codelabs directly from the preview app.

Clare Bayley has been the guru of onsite codelab experiences, running events large and small, while Sam Thorogood and Chris Broadfoot made major contributions to the onsite kiosks you may have seen at Google I/O.

Eric Bidelman redesigned the codelab user interface using Polymer components and built the g.co/codelabs landing page, to provide a beautiful user experience that looks great and works equally well on desktop and mobile devices.

Lots of other contributions have been made over the years and I’m sure that I’m neglecting some important advances but for the sake of brevity, I’ll leave it at that.

## Acknowledgements

Google Codelabs exists thanks to the talents and efforts of many fine volunteers, including:
Alex Vaghin, Marc Cohen, Shawn Simister, Ewa Gasperowicz, Eric Bidelman, Robert Kubis, Clare Bayley, Cassie Recher, Chris Broadfoot, Sam Thorogood, Ryan Seys, and the many codelab authors, inside and outside of Google, who have generated a veritable [treasure trove of content](https://g.co/codelabs).

## Notes

This is not an official Google product.
