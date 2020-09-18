# Markdown Parser

The Markdown codelab parser converts a Markdown document into an HTML file and
codelab metadata file.

## Metadata

Metadata consists of key-value pairs of the form "key: value". Keys cannot
contain colons, and separate metadata fields must be separated by blank lines.
At present, values must all be on one line. All metadata must come before the
title. Any arbitrary keys and values may be used; however, only the following
will be understood by the renderer:

- Summary: A human-readable summary of the codelab. Defaults to blank.
- Id: An identifier composed of lowercase letters ideally describing the
  content of the codelab. This field should be unique among
  codelabs.
- Categories: A comma-separated list of the topics the codelab covers.
- Environments: A list of environments the codelab should be discoverable in.
  Codelabs marked "Web" will be visible at the codelabs index. Codelabs marked
  "Kiosk" will only be available at codelabs kiosks, which have special
  equipment attached.
- Status: The publication status of the codelab. Valid values are:
  - Draft: Codelab is not finished.
  - Published: Codelab is finished and visible.
  - Deprecated: Codelab is considered stale and should not be widely advertised.
  - Hidden: Codelab is not shown in index.
- Feedback Link: A link to send users to if they wish to leave feedback on the
  codelab.
- Analytics Account: A Google Analytics ID to include with all codelab pages.

## Title

The title of the codelab directly follows the metadata. The title is a Header 1.

```
# Title of codelab
```

## Steps

A step is declared by putting the step's title in a Header 2. All content
following a step title will be considered part of the step, until the next step
title, or the end of the document.

```
## Codelab Step
```

### Duration

Steps should be marked with the expected duration to complete them. To label a
step with a duration, put "Duration: TIME" by itself on the line directly
following the step title, where TIME is formatted like "hh:mm:ss" (or "mm:ss" if
only one `:` is provided).

```
## Codelab Step
Duration: 1:25
```

### Content

Codelab content may be written in standard Markdown. Some special constructs are
understood:

#### Fenced Code and Language Hints

Code blocks may be declared by placing them between two lines containing just
three backticks (fenced code blocks). The codelab renderer will attempt to
perform syntax highlighting on code blocks, but it is not always effective at
guessing the language to highlight in. Put the name of the code language after
the first fence to explicitly specify which highlighting plan to use.

    ```go
    This block will be highlighted as Go source code.
    ```

If you'd like to disable syntax highlighting, you can specify the language
hint to "console":

    ```console
    This block will not be syntax highlighted.
    ```

#### Info Boxes

Info boxes are colored callouts that enclose special information in codelabs.
Positive info boxes should contain positive information like best practices and
time-saving tips. Negative infoboxes should contain information like warnings
and API usage restriction. To create an infobox, put the type of infobox on a
line by itself, then begin the next line with a colon.

```
Positive
: This will appear in a positive info box.

Negative
: This will appear in a negative info box.
```

`<aside>` elements work as well:

```
<aside class="positive">
This will appear in a positive info box.
</aside>

<aside class="negative">
This will appear in a negative info box.
</aside>
```

#### Download Buttons

Codelabs sometimes contain links to SDKs or sample code. The codelab renderer
will apply special button-esque styling to any link that begins with the word
"Download".

```
<button>
  [Download SDK](https://www.google.com)
</button>
```

