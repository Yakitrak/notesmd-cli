# mdast-util-to-hast

[![Build][badge-build-image]][badge-build-url]
[![Coverage][badge-coverage-image]][badge-coverage-url]
[![Downloads][badge-downloads-image]][badge-downloads-url]
[![Size][badge-size-image]][badge-size-url]

[mdast][github-mdast] utility to transform to [hast][github-hast].

## Contents

* [What is this?](#what-is-this)
* [When should I use this?](#when-should-i-use-this)
* [Install](#install)
* [Use](#use)
* [API](#api)
  * [`defaultFootnoteBackContent(referenceIndex, rereferenceIndex)`](#defaultfootnotebackcontentreferenceindex-rereferenceindex)
  * [`defaultFootnoteBackLabel(referenceIndex, rereferenceIndex)`](#defaultfootnotebacklabelreferenceindex-rereferenceindex)
  * [`defaultHandlers`](#defaulthandlers)
  * [`toHast(tree[, options])`](#tohasttree-options)
  * [`FootnoteBackContentTemplate`](#footnotebackcontenttemplate)
  * [`FootnoteBackLabelTemplate`](#footnotebacklabeltemplate)
  * [`Handler`](#handler)
  * [`Handlers`](#handlers)
  * [`Options`](#options)
  * [`Raw`](#raw)
  * [`State`](#state)
* [Examples](#examples)
  * [Example: supporting HTML in markdown naïvely](#example-supporting-html-in-markdown-naïvely)
  * [Example: supporting HTML in markdown properly](#example-supporting-html-in-markdown-properly)
  * [Example: footnotes in languages other than English](#example-footnotes-in-languages-other-than-english)
  * [Example: supporting custom nodes](#example-supporting-custom-nodes)
* [Algorithm](#algorithm)
  * [Default handling](#default-handling)
  * [Fields on nodes](#fields-on-nodes)
* [CSS](#css)
* [Syntax tree](#syntax-tree)
  * [Nodes](#nodes)
* [Types](#types)
* [Compatibility](#compatibility)
* [Security](#security)
* [Related](#related)
* [Contribute](#contribute)
* [License](#license)

## What is this?

This package is a utility that takes an [mdast][github-mdast] (markdown) syntax
tree as input and turns it into a [hast][github-hast] (HTML) syntax tree.

## When should I use this?

This project is useful when you want to deal with ASTs and turn markdown to
HTML.

The hast utility [`hast-util-to-mdast`][github-hast-util-to-mdast] does the
inverse of this utility.
It turns HTML into markdown.

The remark plugin [`remark-rehype`][github-remark-rehype] wraps this utility to
also turn markdown to HTML at a higher-level (easier) abstraction.

## Install

This package is [ESM only][github-gist-esm].
In Node.js (version 16+), install with [npm][npmjs-install]:

```sh
npm install mdast-util-to-hast
```

In browsers with [`esm.sh`][esmsh]:

```html
<script type="module">
  import {toHast} from 'https://esm.sh/mdast-util-to-hast@13?bundle'
</script>
```

## Use

Say we have the following `example.md`:

```markdown
## Hello **World**!
```

…and next to it a module `example.js`:

```js
import fs from 'node:fs/promises'
import {toHtml} from 'hast-util-to-html'
import {fromMarkdown} from 'mdast-util-from-markdown'
import {toHast} from 'mdast-util-to-hast'

const markdown = String(await fs.readFile('example.md'))
const mdast = fromMarkdown(markdown)
const hast = toHast(mdast)
const html = toHtml(hast)

console.log(html)
```

…now running `node example.js` yields:

```html
<h2>Hello <strong>World</strong>!</h2>
```

## API

This package exports the identifiers
[`defaultFootnoteBackContent`][api-default-footnote-back-content],
[`defaultFootnoteBackLabel`][api-default-footnote-back-label],
[`defaultHandlers`][api-default-handlers],
and [`toHast`][api-to-hast].
There is no default export.

### `defaultFootnoteBackContent(referenceIndex, rereferenceIndex)`

Generate the default content that GitHub uses on backreferences.

###### Parameters

* `referenceIndex` (`number`)
  — index of the definition in the order that they are first referenced,
  0-indexed
* `rereferenceIndex` (`number`)
  — index of calls to the same definition, 0-indexed

###### Returns

Content (`Array<ElementContent>`).

### `defaultFootnoteBackLabel(referenceIndex, rereferenceIndex)`

Generate the default label that GitHub uses on backreferences.

###### Parameters

* `referenceIndex` (`number`)
  — index of the definition in the order that they are first referenced,
  0-indexed
* `rereferenceIndex` (`number`)
  — index of calls to the same definition, 0-indexed

###### Returns

Label (`string`).

### `defaultHandlers`

Default handlers for nodes ([`Handlers`][api-handlers]).

### `toHast(tree[, options])`

Transform mdast to hast.

###### Parameters

* `tree` ([`MdastNode`][github-mdast-nodes])
  — mdast tree
* `options` ([`Options`][api-options], optional)
  — configuration

###### Returns

hast tree ([`HastNode`][github-hast-nodes]).

##### Notes

###### HTML

Raw HTML is available in mdast as [`html`][github-mdast-html] nodes and can be
embedded in hast as semistandard `raw` nodes.
Most utilities ignore `raw` nodes but two notable ones don’t:

* [`hast-util-to-html`][github-hast-util-to-html] also has an option
  `allowDangerousHtml` which will output the raw HTML.
  This is typically discouraged as noted by the option name but is useful if you
  completely trust authors
* [`hast-util-raw`][github-hast-util-raw] can handle the raw embedded HTML
  strings by parsing them into standard hast nodes (`element`, `text`, etc).
  This is a heavy task as it needs a full HTML parser,
  but it is the only way to support untrusted content

###### Footnotes

Many options supported here relate to footnotes.
Footnotes are not specified by CommonMark, which we follow by default.
They are supported by GitHub,
so footnotes can be enabled in markdown with
[`mdast-util-gfm`][github-mdast-util-gfm].

The options `footnoteBackLabel` and `footnoteLabel` define natural language that
explains footnotes,
which is hidden for sighted users but shown to assistive technology.
When your page is not in English, you must define translated values.

Back references use ARIA attributes,
but the section label itself uses a heading that is hidden with an `sr-only`
class.
To show it to sighted users,
define different attributes in `footnoteLabelProperties`.

###### Clobbering

Footnotes introduces a problem,
as it links footnote calls to footnote definitions on the page through `id`
attributes generated from user content,
which results in DOM clobbering.

DOM clobbering is this:

```html
<p id=x></p>
<script>alert(x) // `x` now refers to the DOM `p#x` element</script>
```

Elements by their ID are made available by browsers on the `window` object,
which is a security risk.
Using a prefix solves this problem.

More information on how to handle clobbering and the prefix is explained in
[Example: headings (DOM clobbering) in
`rehype-sanitize`][github-rehype-sanitize-example-clobbering].

###### Unknown nodes

Unknown nodes are nodes with a type that isn’t in `handlers` or `passThrough`.
The default behavior for unknown nodes is:

* when the node has a `value` (and doesn’t have `data.hName`,
  `data.hProperties`, or `data.hChildren`, see later),
  create a hast `text` node
* otherwise,
  create a `<div>` element (which could be changed with `data.hName`),
  with its children mapped from mdast to hast as well

This behavior can be changed by passing an `unknownHandler`.

### `FootnoteBackContentTemplate`

Generate content for the backreference dynamically.

For the following markdown:

```markdown
Alpha[^micromark], bravo[^micromark], and charlie[^remark].

[^remark]: things about remark
[^micromark]: things about micromark
```

This function will be called with:

* `0` and `0` for the backreference from `things about micromark` to `alpha`,
  as it is the first used definition,
  and the first call to it
* `0` and `1` for the backreference from `things about micromark` to `bravo`,
  as it is the first used definition,
  and the second call to it
* `1` and `0` for the backreference from `things about remark` to `charlie`,
  as it is the second used definition

###### Parameters

* `referenceIndex` (`number`)
  — index of the definition in the order that they are first referenced,
  0-indexed
* `rereferenceIndex` (`number`)
  — index of calls to the same definition, 0-indexed

###### Returns

Content for the backreference when linking back from definitions to their
reference (`Array<ElementContent>`, `ElementContent`, or `string`).

### `FootnoteBackLabelTemplate`

Generate a back label dynamically.

For the following markdown:

```markdown
Alpha[^micromark], bravo[^micromark], and charlie[^remark].

[^remark]: things about remark
[^micromark]: things about micromark
```

This function will be called with:

* `0` and `0` for the backreference from `things about micromark` to `alpha`,
  as it is the first used definition,
  and the first call to it
* `0` and `1` for the backreference from `things about micromark` to `bravo`,
  as it is the first used definition,
  and the second call to it
* `1` and `0` for the backreference from `things about remark` to `charlie`,
  as it is the second used definition

###### Parameters

* `referenceIndex` (`number`)
  — index of the definition in the order that they are first referenced,
  0-indexed
* `rereferenceIndex` (`number`)
  — index of calls to the same definition, 0-indexed

###### Returns

Back label to use when linking back from definitions to their reference
(`string`).

### `Handler`

Handle a node (TypeScript type).

###### Parameters

* `state` ([`State`][api-state])
  — info passed around
* `node` ([`MdastNode`][github-mdast-nodes])
  — node to handle
* `parent` ([`MdastNode | undefined`][github-mdast-nodes])
  — parent of `node`

###### Returns

Result ([`Array<HastNode> | HastNode | undefined`][github-hast-nodes]).

### `Handlers`

Handle nodes (TypeScript type).

###### Type

```ts
type Handlers = Partial<Record<Nodes['type'], Handler>>
```

### `Options`

Configuration (TypeScript type).

###### Fields

* `allowDangerousHtml` (`boolean`, default: `false`)
  — whether to persist raw HTML in markdown in the hast tree
* `clobberPrefix` (`string`, default: `'user-content-'`)
  — prefix to use before the `id` property on footnotes to prevent them from
  *clobbering*
* `file` ([`VFile`][github-vfile], optional)
  — corresponding virtual file representing the input document
* `footnoteBackContent`
  ([`FootnoteBackContentTemplate`][api-footnote-back-content-template] or
  `string`, default:
  [`defaultFootnoteBackContent`][api-default-footnote-back-content])
  — content of the backreference back to references
* `footnoteBackLabel`
  ([`FootnoteBackLabelTemplate`][api-footnote-back-label-template] or `string`,
  default: [`defaultFootnoteBackLabel`][api-default-footnote-back-label])
  — label to describe the backreference back to references
* `footnoteLabel` (`string`, default: `'Footnotes'`)
  — label to use for the footnotes section (affects screen readers)
* `footnoteLabelProperties`
  ([`Properties`][github-hast-properties], default: `{className: ['sr-only']}`)
  — properties to use on the footnote label
  (note that `id: 'footnote-label'` is always added as footnote calls use it
  with `aria-describedby` to provide an accessible label)
* `footnoteLabelTagName` (`string`, default: `h2`)
  — tag name to use for the footnote label
* `handlers` ([`Handlers`][api-handlers], optional)
  — extra handlers for nodes
* `passThrough` (`Array<Nodes['type']>`, optional)
  — list of custom mdast node types to pass through (keep) in hast (note that
  the node itself is passed, but eventual children are transformed)
* `unknownHandler` ([`Handler`][api-handler], optional)
  — handle all unknown nodes

### `Raw`

Raw string of HTML embedded into HTML AST (TypeScript type).

###### Type

```ts
import type {Data, Literal} from 'hast'

interface Raw extends Literal {
  type: 'raw'
  data?: RawData | undefined
}

interface RawData extends Data {}
```

### `State`

Info passed around about the current state (TypeScript type).

###### Fields

* `all` (`(node: MdastNode) => Array<HastNode>`)
  — transform the children of an mdast parent to hast
* `applyData`
  (`<Type extends HastNode>(from: MdastNode, to: Type) => Type | HastElement`)
  — honor the `data` of `from` and maybe generate an element instead of `to`
* `definitionById` (`Map<string, Definition>`)
  — definitions by their uppercased identifier
* `footnoteById` (`Map<string, FootnoteDefinition>`)
  — footnote definitions by their uppercased identifier
* `footnoteCounts` (`Map<string, number>`)
  — counts for how often the same footnote was called
* `footnoteOrder` (`Array<string>`)
  — identifiers of order when footnote calls first appear in tree order
* `handlers` ([`Handlers`][api-handlers])
  — applied node handlers
* `one`
  (`(node: MdastNode, parent: MdastNode | undefined) => HastNode | Array<HastNode> | undefined`)
  — transform an mdast node to hast
* `options` ([`Options`][api-options])
  — configuration
* `patch` (`(from: MdastNode, to: HastNode) => undefined`)
* `wrap`
  (`<Type extends HastNode>(nodes: Array<Type>, loose?: boolean) => Array<Type | HastText>`)
  — wrap `nodes` with line endings between each node,
  adds initial/final line endings when `loose`

## Examples

### Example: supporting HTML in markdown naïvely

If you completely trust authors (or plugins) and want to allow them to HTML *in*
markdown,
and the last utility has an `allowDangerousHtml` option as well (such as
`hast-util-to-html`) you can pass `allowDangerousHtml` to this utility
(`mdast-util-to-hast`):

```js
import {fromMarkdown} from 'mdast-util-from-markdown'
import {toHast} from 'mdast-util-to-hast'
import {toHtml} from 'hast-util-to-html'

const markdown = 'It <i>works</i>! <img onerror="alert(1)">'
const mdast = fromMarkdown(markdown)
const hast = toHast(mdast, {allowDangerousHtml: true})
const html = toHtml(hast, {allowDangerousHtml: true})

console.log(html)
```

…now running `node example.js` yields:

```html
<p>It <i>works</i>! <img onerror="alert(1)"></p>
```

> ⚠️ **Danger**: observe that the XSS attack through the `onerror` attribute is
> still present.

### Example: supporting HTML in markdown properly

If you do not trust the authors of the input markdown,
or if you want to make sure that further utilities can see HTML embedded in
markdown,
use [`hast-util-raw`][github-hast-util-raw].
The following example passes `allowDangerousHtml` to this utility
(`mdast-util-to-hast`), then turns the raw embedded HTML into proper HTML nodes
(`hast-util-raw`), and finally sanitizes the HTML by only allowing safe things
(`hast-util-sanitize`):

```js
import {raw} from 'hast-util-raw'
import {sanitize} from 'hast-util-sanitize'
import {toHtml} from 'hast-util-to-html'
import {fromMarkdown} from 'mdast-util-from-markdown'
import {toHast} from 'mdast-util-to-hast'

const markdown = 'It <i>works</i>! <img onerror="alert(1)">'
const mdast = fromMarkdown(markdown)
const hast = raw(toHast(mdast, {allowDangerousHtml: true}))
const safeHast = sanitize(hast)
const html = toHtml(safeHast)

console.log(html)
```

…now running `node example.js` yields:

```html
<p>It <i>works</i>! <img></p>
```

> 👉 **Note**: observe that the XSS attack through the `onerror` attribute is no
> longer present.

### Example: footnotes in languages other than English

If you know that the markdown is authored in a language other than English,
and you’re using `micromark-extension-gfm` and `mdast-util-gfm` to match how
GitHub renders markdown,
and you know that footnotes are (or can?) be used,
you should translate the labels associated with them.

Let’s first set the stage:

```js
import {toHtml} from 'hast-util-to-html'
import {gfm} from 'micromark-extension-gfm'
import {fromMarkdown} from 'mdast-util-from-markdown'
import {gfmFromMarkdown} from 'mdast-util-gfm'
import {toHast} from 'mdast-util-to-hast'

const markdown = 'Bonjour[^1]\n\n[^1]: Monde!'
const mdast = fromMarkdown(markdown, {
  extensions: [gfm()],
  mdastExtensions: [gfmFromMarkdown()]
})
const hast = toHast(mdast)
const html = toHtml(hast)

console.log(html)
```

…now running `node example.js` yields:

```html
<p>Bonjour<sup><a href="#user-content-fn-1" id="user-content-fnref-1" data-footnote-ref aria-describedby="footnote-label">1</a></sup></p>
<section data-footnotes class="footnotes"><h2 class="sr-only" id="footnote-label">Footnotes</h2>
<ol>
<li id="user-content-fn-1">
<p>Monde! <a href="#user-content-fnref-1" data-footnote-backref="" aria-label="Back to reference 1" class="data-footnote-backref">↩</a></p>
</li>
</ol>
</section>
```

This is a mix of English and French that screen readers can’t handle nicely.
Let’s say our program does know that the markdown is in French.
In that case,
it’s important to translate and define the labels relating to footnotes so that
screen reader users can properly pronounce the page:

```diff
@@ -9,7 +9,16 @@ const mdast = fromMarkdown(markdown, {
   extensions: [gfm()],
   mdastExtensions: [gfmFromMarkdown()]
 })
-const hast = toHast(mdast)
+const hast = toHast(mdast, {
+  footnoteLabel: 'Notes de bas de page',
+  footnoteBackLabel(referenceIndex, rereferenceIndex) {
+    return (
+      'Retour à la référence ' +
+      (referenceIndex + 1) +
+      (rereferenceIndex > 1 ? '-' + rereferenceIndex : '')
+    )
+  }
+})
 const html = toHtml(hast)

 console.log(html)
```

…now running `node example.js` with the above patch applied yields:

```diff
@@ -1,8 +1,8 @@
 <p>Bonjour<sup><a href="#user-content-fn-1" id="user-content-fnref-1" data-footnote-ref aria-describedby="footnote-label">1</a></sup></p>
-<section data-footnotes class="footnotes"><h2 class="sr-only" id="footnote-label">Footnotes</h2>
+<section data-footnotes class="footnotes"><h2 class="sr-only" id="footnote-label">Notes de bas de page</h2>
 <ol>
 <li id="user-content-fn-1">
-<p>Monde! <a href="#user-content-fnref-1" data-footnote-backref="" aria-label="Back to reference 1" class="data-footnote-backref">↩</a></p>
+<p>Monde! <a href="#user-content-fnref-1" data-footnote-backref="" aria-label="Retour à la référence 1" class="data-footnote-backref">↩</a></p>
 </li>
 </ol>
 </section>
```

### Example: supporting custom nodes

This project supports CommonMark and the GFM constructs (footnotes,
strikethrough, tables) and the frontmatter constructs YAML and TOML.
Support can be extended to other constructs in two ways: a) with handlers,
b) through fields on nodes.

For example,
when we represent a mark element in markdown and want to turn it into a `<mark>`
element in HTML,
we can use a handler:

```js
import {toHtml} from 'hast-util-to-html'
import {toHast} from 'mdast-util-to-hast'

const mdast = {
  type: 'paragraph',
  children: [{type: 'mark', children: [{type: 'text', value: 'x'}]}]
}

const hast = toHast(mdast, {
  handlers: {
    mark(state, node) {
      return {
        type: 'element',
        tagName: 'mark',
        properties: {},
        children: state.all(node)
      }
    }
  }
})

console.log(toHtml(hast))
```

We can do the same through certain fields on nodes:

```js
import {toHtml} from 'hast-util-to-html'
import {toHast} from 'mdast-util-to-hast'

const mdast = {
  type: 'paragraph',
  children: [
    {
      type: 'mark',
      children: [{type: 'text', value: 'x'}],
      data: {hName: 'mark'}
    }
  ]
}

console.log(toHtml(toHast(mdast)))
```

## Algorithm

This project by default handles CommonMark,
GFM (footnotes, strikethrough, tables) and common frontmatter (YAML, TOML).

Existing handlers can be overwritten and handlers for more nodes can be added.
It’s also possible to define how mdast is turned into hast through fields on
nodes.

### Default handling

The following table gives insight into what input turns into what output:

<table>
<thead>
<tr>
<th scope="col">mdast node</th>
<th scope="col">markdown example</th>
<th scope="col">hast node</th>
<th scope="col">html example</th>
</tr>
</thead>
<tbody>
<tr>
<th scope="row">

[`blockquote`](https://github.com/syntax-tree/mdast#blockquote)

</th>
<td>

```markdown
> A greater than…
```

</td>
<td>

[`element`][github-hast-element] (`blockquote`)

</td>
<td>

```html
<blockquote>
<p>A greater than…</p>
</blockquote>
```

</td>
</tr>
<tr>
<th scope="row">

[`break`](https://github.com/syntax-tree/mdast#break)

</th>
<td>

```markdown
A backslash\
before a line break…
```

</td>
<td>

[`element`][github-hast-element] (`br`)

</td>
<td>

```html
<p>A backslash<br>
before a line break…</p>
```

</td>
</tr>
<tr>
<th scope="row">

[`code`](https://github.com/syntax-tree/mdast#code)

</th>
<td>

````markdown
```js
backtick.fences('for blocks')
```
````

</td>
<td>

[`element`][github-hast-element] (`pre` and `code`)

</td>
<td>

```html
<pre><code className="language-js">backtick.fences('for blocks')
</code></pre>
```

</td>
</tr>
<tr>
<th scope="row">

[`delete`](https://github.com/syntax-tree/mdast#delete) (GFM)

</th>
<td>

```markdown
Two ~~tildes~~ for delete.
```

</td>
<td>

[`element`][github-hast-element] (`del`)

</td>
<td>

```html
<p>Two <del>tildes</del> for delete.</p>
```

</td>
</tr>
<tr>
<th scope="row">

[`emphasis`](https://github.com/syntax-tree/mdast#emphasis)

</th>
<td>

```markdown
Some *asterisks* for emphasis.
```

</td>
<td>

[`element`][github-hast-element] (`em`)

</td>
<td>

```html
<p>Some <em>asterisks</em> for emphasis.</p>
```

</td>
</tr>
<tr>
<th scope="row">

[`footnoteReference`](https://github.com/syntax-tree/mdast#footnotereference),
[`footnoteDefinition`](https://github.com/syntax-tree/mdast#footnotedefinition)
(GFM)

</th>
<td>

```markdown
With a [^caret].

[^caret]: Stuff
```

</td>
<td>

[`element`][github-hast-element] (`section`, `sup`, `a`)

</td>
<td>

```html
<p>With a <sup><a href="#fn-caret" …>1</a></sup>.</p>…
```

</td>
</tr>
<tr>
<th scope="row">

[`heading`](https://github.com/syntax-tree/mdast#heading)

</th>
<td>

```markdown
# One number sign…
###### Six number signs…
```

</td>
<td>

[`element`][github-hast-element] (`h1`…`h6`)

</td>
<td>

```html
<h1>One number sign…</h1>
<h6>Six number signs…</h6>
```

</td>
</tr>
<tr>
<th scope="row">

[`html`][github-mdast-html]

</th>
<td>

```html
<kbd>CMD+S</kbd>
```

</td>
<td>

Nothing (default), `raw` (when `allowDangerousHtml: true`)

</td>
<td>

n/a

</td>
</tr>
<tr>
<th scope="row">

[`image`](https://github.com/syntax-tree/mdast#image)

</th>
<td>

```markdown
![Alt text](/logo.png "title")
```

</td>
<td>

[`element`][github-hast-element] (`img`)

</td>
<td>

```html
<p><img src="/logo.png" alt="Alt text" title="title"></p>
```

</td>
</tr>
<tr>
<th scope="row">

[`imageReference`](https://github.com/syntax-tree/mdast#imagereference),
[`definition`][github-mdast-definition]

</th>
<td>

```markdown
![Alt text][logo]

[logo]: /logo.png "title"
```

</td>
<td>

[`element`][github-hast-element] (`img`)

</td>
<td>

```html
<p><img src="/logo.png" alt="Alt text" title="title"></p>
```

</td>
</tr>
<tr>
<th scope="row">

[`inlineCode`](https://github.com/syntax-tree/mdast#inlinecode)

</th>
<td>

```markdown
Some `backticks` for inline code.
```

</td>
<td>

[`element`][github-hast-element] (`code`)

</td>
<td>

```html
<p>Some <code>backticks</code> for inline code.</p>
```

</td>
</tr>
<tr>
<th scope="row">

[`link`](https://github.com/syntax-tree/mdast#link)

</th>
<td>

```markdown
[Example](https://example.com "title")
```

</td>
<td>

[`element`][github-hast-element] (`a`)

</td>
<td>

```html
<p><a href="https://example.com" title="title">Example</a></p>
```

</td>
</tr>
<tr>
<th scope="row">

[`linkReference`](https://github.com/syntax-tree/mdast#linkreference),
[`definition`][github-mdast-definition]

</th>
<td>

```markdown
[Example][]

[example]: https://example.com "title"
```

</td>
<td>

[`element`][github-hast-element] (`a`)

</td>
<td>

```html
<p><a href="https://example.com" title="title">Example</a></p>
```

</td>
</tr>
<tr>
<th scope="row">

[`list`](https://github.com/syntax-tree/mdast#list),
[`listItem`](https://github.com/syntax-tree/mdast#listitem)

</th>
<td>

```markdown
* asterisks for unordered items

1. decimals and a dot for ordered items
```

</td>
<td>

[`element`][github-hast-element] (`li` and `ol` or `ul`)

</td>
<td>

```html
<ul>
<li>asterisks for unordered items</li>
</ul>
<ol>
<li>decimals and a dot for ordered items</li>
</ol>
```

</td>
</tr>
<tr>
<th scope="row">

[`paragraph`](https://github.com/syntax-tree/mdast#paragraph)

</th>
<td>

```markdown
Just some text…
```

</td>
<td>

[`element`][github-hast-element] (`p`)

</td>
<td>

```html
<p>Just some text…</p>
```

</td>
</tr>
<tr>
<th scope="row">

[`root`](https://github.com/syntax-tree/mdast#root)

</th>
<td>

```markdown
Anything!
```

</td>
<td>

[`root`](https://github.com/syntax-tree/hast#root)

</td>
<td>

```html
<p>Anything!</p>
```

</td>
</tr>
<tr>
<th scope="row">

[`strong`](https://github.com/syntax-tree/mdast#strong)

</th>
<td>

```markdown
Two **asterisks** for strong.
```

</td>
<td>

[`element`][github-hast-element] (`strong`)

</td>
<td>

```html
<p>Two <strong>asterisks</strong> for strong.</p>
```

</td>
</tr>
<tr>
<th scope="row">

[`text`](https://github.com/syntax-tree/mdast#text)

</th>
<td>

```markdown
Anything!
```

</td>
<td>

[`text`](https://github.com/syntax-tree/hast#text)

</td>
<td>

```html
<p>Anything!</p>
```

</td>
</tr>
<tr>
<th scope="row">

[`table`](https://github.com/syntax-tree/mdast#table),
[`tableRow`](https://github.com/syntax-tree/mdast#tablerow),
[`tableCell`](https://github.com/syntax-tree/mdast#tablecell)

</th>
<td>

```markdown
| Pipes |
| ----- |
```

</td>
<td>

[`element`][github-hast-element] (`table`, `thead`, `tbody`, `tr`, `td`, `th`)

</td>
<td>

```html
<table>
<thead>
<tr>
<th>Pipes</th>
</tr>
</thead>
</table>
```

</td>
</tr>
<tr>
<th scope="row">

[`thematicBreak`](https://github.com/syntax-tree/mdast#thematicbreak)

</th>
<td>

```markdown
Three asterisks for a thematic break:

***
```

</td>
<td>

[`element`][github-hast-element] (`hr`)

</td>
<td>

```html
<p>Three asterisks for a thematic break:</p>
<hr>
```

</td>
</tr>
<tr>
<th scope="row">

`toml` (frontmatter)

</th>
<td>

```markdown
+++
fenced = true
+++
```

</td>
<td>

Nothing

</td>
<td>

n/a

</td>
</tr>
<tr>
<th scope="row">

[`yaml`](https://github.com/syntax-tree/mdast#yaml) (frontmatter)

</th>
<td>

```markdown
---
fenced: yes
---
```

</td>
<td>

Nothing

</td>
<td>

n/a

</td>
</tr>
</tbody>
</table>

> 👉 **Note**: GFM prescribes that the obsolete `align` attribute on `td` and
> `th` elements is used.
> To use `style` attributes instead of obsolete features,
> combine this utility with
> [`@mapbox/hast-util-table-cell-style`][github-hast-util-table-cell-style].

> 🧑‍🏫 **Info**: this project is concerned with turning one syntax tree into
> another.
> It does not deal with markdown syntax or HTML syntax.
> The preceding examples are illustrative rather than authoritative or
> exhaustive.

### Fields on nodes

A frequent problem arises when having to turn one syntax tree into another.
As the original tree (in this case, mdast for markdown) is in some cases limited
compared to the destination (in this case, hast for HTML) tree,
is it possible to provide more info in the original to define what the result
will be in the destination?
This is possible by defining data on mdast nodes,
which this utility will read as instructions on what hast nodes to create.

An example is math,
which is a nonstandard markdown extension,
that this utility doesn’t understand.
To solve this,
`mdast-util-math` defines instructions on mdast nodes that this plugin does
understand because they define a certain hast structure.

The following fields can be used:

* `node.data.hName` — define the element’s tag name
* `node.data.hProperties` — define extra properties to use
* `node.data.hChildren` — define hast children to use

###### `hName`

`node.data.hName` sets the tag name of an element.
The following [mdast][github-mdast]:

```js
{
  type: 'strong',
  data: {hName: 'b'},
  children: [{type: 'text', value: 'Alpha'}]
}
```

…yields ([hast][github-hast]):

```js
{
  type: 'element',
  tagName: 'b',
  properties: {},
  children: [{type: 'text', value: 'Alpha'}]
}
```

###### `hProperties`

`node.data.hProperties` sets the properties of an element.
The following [mdast][github-mdast]:

```js
{
  type: 'image',
  src: 'circle.svg',
  alt: 'Big red circle on a black background',
  data: {hProperties: {className: ['responsive']}}
}
```

…yields ([hast][github-hast]):

```js
{
  type: 'element',
  tagName: 'img',
  properties: {
    src: 'circle.svg',
    alt: 'Big red circle on a black background',
    className: ['responsive']
  },
  children: []
}
```

###### `hChildren`

`node.data.hChildren` sets the children of an element.
The following [mdast][github-mdast]:

```js
{
  type: 'code',
  lang: 'js',
  data: {
    hChildren: [
      {
        type: 'element',
        tagName: 'span',
        properties: {className: ['hljs-meta']},
        children: [{type: 'text', value: '"use strict"'}]
      },
      {type: 'text', value: ';'}
    ]
  },
  value: '"use strict";'
}
```

…yields ([hast][github-hast]):

```js
{
  type: 'element',
  tagName: 'pre',
  properties: {},
  children: [{
    type: 'element',
    tagName: 'code',
    properties: {className: ['language-js']},
    children: [
      {
        type: 'element',
        tagName: 'span',
        properties: {className: ['hljs-meta']},
        children: [{type: 'text', value: '"use strict"'}]
      },
      {type: 'text', value: ';'}
    ]
  }]
}
```

> 👉 **Note**: the `pre` and `language-js` class are normal `mdast-util-to-hast`
> functionality.

## CSS

Assuming you know how to use (semantic) HTML and CSS,
then it should generally be straightforward to style the HTML produced by this
plugin.
With CSS, you can get creative and style the results as you please.

Some semistandard features,
notably GFMs tasklists and footnotes,
generate HTML that be unintuitive,
as it matches exactly what GitHub produces for their website.
There is a project,
[`sindresorhus/github-markdown-css`][github-github-markdown-css],
that exposes the stylesheet that GitHub uses for rendered markdown,
which might either be inspirational for more complex features,
or can be used as-is to exactly match how GitHub styles rendered markdown.

The following CSS is needed to make footnotes look a bit like GitHub:

```css
/* Style the footnotes section. */
.footnotes {
  font-size: smaller;
  color: #8b949e;
  border-top: 1px solid #30363d;
}

/* Hide the section label for visual users. */
.sr-only {
  position: absolute;
  width: 1px;
  height: 1px;
  padding: 0;
  overflow: hidden;
  clip: rect(0, 0, 0, 0);
  word-wrap: normal;
  border: 0;
}

/* Place `[` and `]` around footnote calls. */
[data-footnote-ref]::before {
  content: '[';
}

[data-footnote-ref]::after {
  content: ']';
}
```

## Syntax tree

The following interfaces are added to **[hast][github-hast]** by this utility.

### Nodes

#### `Raw`

```idl
interface Raw <: Literal {
  type: 'raw'
}
```

**Raw** (**[Literal][github-hast-literal]**) represents a string if raw HTML
inside hast.
Raw nodes are typically ignored but are handled by
[`hast-util-to-html`][github-hast-util-to-html] and
[`hast-util-raw`][github-hast-util-raw].

## Types

This package is fully typed with [TypeScript][].
It exports the
[`FootnoteBackContentTemplate`][api-footnote-back-content-template],
[`FootnoteBackLabelTemplate`][api-footnote-back-label-template],
[`Handler`][api-handler],
[`Handlers`][api-handlers],
[`Options`][api-options],
[`Raw`][api-raw],
and [`State`][api-state] types.

It also registers the `Raw` node type with `@types/hast`.
If you’re working with the syntax tree (and you pass
`allowDangerousHtml: true`),
make sure to import this utility somewhere in your types,
as that registers the new node type in the tree.

```js
/**
 * @import {Root} from 'hast'
 * @import {} from 'mdast-util-to-hast'
 */

import {visit} from 'unist-util-visit'

/** @type {Root} */
const tree = { /* … */ }

visit(tree, function (node) {
  // `node` can now be `raw`.
})
```

Finally,
it also registers the `hChildren`,
`hName`,
and `hProperties` fields on `Data` of `@types/mdast`.
If you’re working with the syntax tree,
make sure to import this utility somewhere in your types,
as that registers the data fields in the tree.

```js
/**
 * @import {Root} from 'hast'
 * @import {} from 'mdast-util-to-hast'
 */

import {visit} from 'unist-util-visit'

/** @type {Root} */
const tree = { /* … */ }

console.log(tree.data?.hName) // Types as `string | undefined`.
```

## Compatibility

Projects maintained by the unified collective are compatible with maintained
versions of Node.js.

When we cut a new major release,
we drop support for unmaintained versions of Node.
This means we try to keep the current release line,
`mdast-util-to-hast@^13`,
compatible with Node.js 16.

## Security

Use of `mdast-util-to-hast` can open you up to a
[cross-site scripting (XSS)][wikipedia-xss] attack.
Embedded hast properties (`hName`, `hProperties`, `hChildren`),
custom handlers,
and the `allowDangerousHtml` option all provide openings.

The following example shows how a script is injected where a benign code block
is expected with embedded hast properties:

```js
const code = {type: 'code', value: 'alert(1)'}

code.data = {hName: 'script'}
```

Yields:

```html
<script>alert(1)</script>
```

The following example shows how an image is changed to fail loading and
therefore run code in a browser.

```js
const image = {type: 'image', url: 'existing.png'}

image.data = {hProperties: {src: 'missing', onError: 'alert(2)'}}
```

Yields:

```html
<img src="missing" onerror="alert(2)">
```

The following example shows the default handling of embedded HTML:

```markdown
# Hello

<script>alert(3)</script>
```

Yields:

```html
<h1>Hello</h1>
```

Passing `allowDangerousHtml: true` to `mdast-util-to-hast` is typically still
not enough to run unsafe code:

```html
<h1>Hello</h1>
&#x3C;script>alert(3)&#x3C;/script>
```

If `allowDangerousHtml: true` is also given to `hast-util-to-html` (or
`rehype-stringify`),
the unsafe code runs:

```html
<h1>Hello</h1>
<script>alert(3)</script>
```

Use [`hast-util-sanitize`][github-hast-util-sanitize] to make the hast tree
safe.

## Related

* [`hast-util-to-mdast`][github-hast-util-to-mdast]
  — transform hast to mdast
* [`hast-util-to-xast`](https://github.com/syntax-tree/hast-util-to-xast)
  — transform hast to xast
* [`hast-util-sanitize`][github-hast-util-sanitize]
  — sanitize hast nodes

## Contribute

See
[`contributing.md`][health-contributing]
in
[`syntax-tree/.github`][health]
for ways to get started.
See [`support.md`][health-support] for ways to get help.

This project has a [code of conduct][health-coc].
By interacting with this repository,
organization,
or community you agree to abide by its terms.

## License

[MIT][file-license] © [Titus Wormer][wooorm]

<!-- Definitions -->

[api-default-footnote-back-content]: #defaultfootnotebackcontentreferenceindex-rereferenceindex

[api-default-footnote-back-label]: #defaultfootnotebacklabelreferenceindex-rereferenceindex

[api-default-handlers]: #defaulthandlers

[api-footnote-back-content-template]: #footnotebackcontenttemplate

[api-footnote-back-label-template]: #footnotebacklabeltemplate

[api-handler]: #handler

[api-handlers]: #handlers

[api-options]: #options

[api-raw]: #raw

[api-state]: #state

[api-to-hast]: #tohasttree-options

[badge-build-image]: https://github.com/syntax-tree/mdast-util-to-hast/workflows/main/badge.svg

[badge-build-url]: https://github.com/syntax-tree/mdast-util-to-hast/actions

[badge-coverage-image]: https://img.shields.io/codecov/c/github/syntax-tree/mdast-util-to-hast.svg

[badge-coverage-url]: https://codecov.io/github/syntax-tree/mdast-util-to-hast

[badge-downloads-image]: https://img.shields.io/npm/dm/mdast-util-to-hast.svg

[badge-downloads-url]: https://www.npmjs.com/package/mdast-util-to-hast

[badge-size-image]: https://img.shields.io/bundlejs/size/mdast-util-to-hast

[badge-size-url]: https://bundlejs.com/?q=mdast-util-to-hast

[esmsh]: https://esm.sh

[file-license]: license

[github-gist-esm]: https://gist.github.com/sindresorhus/a39789f98801d908bbc7ff3ecc99d99c

[github-github-markdown-css]: https://github.com/sindresorhus/github-markdown-css

[github-hast]: https://github.com/syntax-tree/hast

[github-hast-element]: https://github.com/syntax-tree/hast#element

[github-hast-literal]: https://github.com/syntax-tree/hast#literal

[github-hast-nodes]: https://github.com/syntax-tree/hast#nodes

[github-hast-properties]: https://github.com/syntax-tree/hast#properties

[github-hast-util-raw]: https://github.com/syntax-tree/hast-util-raw

[github-hast-util-sanitize]: https://github.com/syntax-tree/hast-util-sanitize

[github-hast-util-table-cell-style]: https://github.com/mapbox/hast-util-table-cell-style

[github-hast-util-to-html]: https://github.com/syntax-tree/hast-util-to-html

[github-hast-util-to-mdast]: https://github.com/syntax-tree/hast-util-to-mdast

[github-mdast]: https://github.com/syntax-tree/mdast

[github-mdast-definition]: https://github.com/syntax-tree/mdast#definition

[github-mdast-html]: https://github.com/syntax-tree/mdast#html

[github-mdast-nodes]: https://github.com/syntax-tree/mdast#nodes

[github-mdast-util-gfm]: https://github.com/syntax-tree/mdast-util-gfm

[github-rehype-sanitize-example-clobbering]: https://github.com/rehypejs/rehype-sanitize#example-headings-dom-clobbering

[github-remark-rehype]: https://github.com/remarkjs/remark-rehype

[github-vfile]: https://github.com/vfile/vfile

[health]: https://github.com/syntax-tree/.github

[health-coc]: https://github.com/syntax-tree/.github/blob/main/code-of-conduct.md

[health-contributing]: https://github.com/syntax-tree/.github/blob/main/contributing.md

[health-support]: https://github.com/syntax-tree/.github/blob/main/support.md

[npmjs-install]: https://docs.npmjs.com/cli/install

[typescript]: https://www.typescriptlang.org

[wikipedia-xss]: https://en.wikipedia.org/wiki/Cross-site_scripting

[wooorm]: https://wooorm.com
