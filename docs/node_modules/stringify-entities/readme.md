# stringify-entities

[![Build][badge-build-image]][badge-build-url]
[![Coverage][badge-coverage-image]][badge-coverage-url]
[![Downloads][badge-downloads-image]][badge-downloads-url]
[![Size][badge-size-image]][badge-size-url]

Serialize (encode) HTML character references.

## Contents

* [What is this?](#what-is-this)
* [When should I use this?](#when-should-i-use-this)
* [Install](#install)
* [Use](#use)
* [API](#api)
  * [`LightOptions`](#lightoptions)
  * [`Options`](#options)
  * [`stringifyEntitiesLight(value[, options])`](#stringifyentitieslightvalue-options)
  * [`stringifyEntities(value[, options])`](#stringifyentitiesvalue-options)
* [Algorithm](#algorithm)
* [Compatibility](#compatibility)
* [Security](#security)
* [Related](#related)
* [Contribute](#contribute)
* [License](#license)

## What is this?

This is a small and powerful encoder of HTML character references
(often called entities).
This one has either all the options you need for a minifier/formatter,
or a tiny size when using
[`stringifyEntitiesLight`][api-stringify-entities-light].

## When should I use this?

You can use this for spec-compliant encoding of character references.
It’s small and fast enough to do that well.
You can also use this when making an HTML formatter or minifier,
because there are different ways to produce pretty or tiny output.
This package is reliable:
``'`'`` characters are encoded to ensure no scripts run in
Internet Explorer 6 to 8.
Additionally, only named references recognized by HTML 4 are encoded,
meaning the infamous `&apos;`
(which people think is a [virus][telegraph-apos])
won’t show up.

## Install

This package is [ESM only][github-gist-esm].
In Node.js (version 14.14+),
install with [npm][npmjs-install]:

```sh
npm install stringify-entities
```

In browsers with [`esm.sh`][esmsh]:

```html
<script type="module">
  import {stringifyEntities} from 'https://esm.sh/stringify-entities@4?bundle'
</script>
```

## Use

```js
import {stringifyEntities} from 'stringify-entities'

stringifyEntities('alpha © bravo ≠ charlie 𝌆 delta')
// => 'alpha &#xA9; bravo &#x2260; charlie &#x1D306; delta'

stringifyEntities('alpha © bravo ≠ charlie 𝌆 delta', {useNamedReferences: true})
// => 'alpha &copy; bravo &ne; charlie &#x1D306; delta'
```

## API

This package exports the identifiers
[`stringifyEntities`][api-stringify-entities] and
[`stringifyEntitiesLight`][api-stringify-entities-light].
There is no default export.
It exports the additional [TypeScript][] types
[`LightOptions`][api-light-options] and
[`Options`][api-options].

### `LightOptions`

Configuration for [`stringifyEntitiesLight`][api-stringify-entities-light]
(TypeScript type).

###### Fields

* `escapeOnly` (`boolean`, default: `false`)
  — whether to only escape possibly dangerous characters;
  those characters are `"`, `&`, `'`, `<`, `>`, and `` ` ``
* `subset` (`ReadonlyArray<string>`, optional)
  — whether to only escape the given subset of characters;
  only BMP characters are supported here (so no emoji)

### `Options`

Configuration for [`stringifyEntities`][api-stringify-entities]
(TypeScript type).

Includes the fields of [`LightOptions`][api-light-options].
If you do not care about the following fields,
use [`stringifyEntitiesLight`][api-stringify-entities-light].

###### Fields

* `attribute` (`boolean`, default: `false`)
  — create character references which don’t fail in attributes;
  only applies when operating dangerously with
  `omitOptionalSemicolons: true`
* `omitOptionalSemicolons` (`boolean`, default: `false`)
  — whether to omit semicolons when possible;
  this creates what HTML calls “parse errors” but is otherwise still valid
  HTML,
  so don’t use this except when building a minifier;
  omitting semicolons is possible for certain named and numeric references in
  some cases
* `useNamedReferences` (`boolean`, default: `false`)
  — prefer named character references (`&amp;`) where possible
* `useShortestReferences` (`boolean`, default: `false`)
  — prefer the shortest possible reference,
  if that results in less bytes;
  `useNamedReferences` can be omitted when using `useShortestReferences`

### `stringifyEntitiesLight(value[, options])`

Encode special characters in `value` as hexadecimal character references.

###### Parameters

* `value` (`string`, required)
  — value to encode
* `options` ([`LightOptions`][api-light-options], optional)
  — configuration

###### Returns

Encoded value (`string`).

### `stringifyEntities(value[, options])`

Encode special characters in `value`.

###### Parameters

* `value` (`string`, required)
  — value to encode
* `options` ([`Options`][api-options], optional)
  — configuration

###### Returns

Encoded value (`string`).

## Algorithm

By default,
all dangerous, non-ASCII, and non-printable ASCII characters are encoded.
A `subset` of characters can be given to encode just those characters.
Alternatively,
pass `escapeOnly` to escape just the dangerous characters
(`"`, `'`, `<`, `>`, `&`, `` ` ``).
By default, hexadecimal character references are used.
Pass `useNamedReferences` to use named character references when possible,
or `useShortestReferences` to use whichever is shortest:
decimal, hexadecimal, or named.

There is also [`stringifyEntitiesLight`][api-stringify-entities-light],
which works just like [`stringifyEntities`][api-stringify-entities]
but without the formatting options:
it’s much smaller but always outputs hexadecimal character references.

## Compatibility

This package is compatible with maintained versions of Node.js.

When we cut a new major release,
we drop support for unmaintained versions of Node.
This means we try to keep the current release line,
`stringify-entities@4`,
compatible with Node.js 14.14.

## Security

By default,
`"`, `&`, `'`, `<`, `>`, and `` ` `` are always encoded,
so the result is safe to use in HTML text and in quoted attribute values.
Spaces and other ASCII characters are not encoded,
so the result is not safe in unquoted attribute values.

When you pass `subset`,
only those characters are encoded.
Leaving out characters such as `<` or `&` can open you up to a
[cross-site scripting (XSS)][wikipedia-xss] attack.

## Related

* [`wooorm/character-entities`][github-character-entities]
  — info on character references
* [`wooorm/character-entities-html4`][github-character-entities-html4]
  — info on HTML 4 character references
* [`wooorm/character-entities-legacy`][github-character-entities-legacy]
  — info on legacy character references
* [`wooorm/character-reference-invalid`][github-character-reference-invalid]
  — info on invalid numeric character references
* [`wooorm/parse-entities`][github-parse-entities]
  — parse (decode) HTML character references

## Contribute

Yes please!
See [How to Contribute to Open Source][opensource-guide-contribute].

## License

[MIT][file-license] © [Titus Wormer][wooorm]

<!-- Definitions -->

[api-light-options]: #lightoptions

[api-options]: #options

[api-stringify-entities]: #stringifyentitiesvalue-options

[api-stringify-entities-light]: #stringifyentitieslightvalue-options

[badge-build-image]: https://github.com/wooorm/stringify-entities/workflows/main/badge.svg

[badge-build-url]: https://github.com/wooorm/stringify-entities/actions

[badge-coverage-image]: https://img.shields.io/codecov/c/github/wooorm/stringify-entities.svg

[badge-coverage-url]: https://codecov.io/github/wooorm/stringify-entities

[badge-downloads-image]: https://img.shields.io/npm/dm/stringify-entities.svg

[badge-downloads-url]: https://www.npmjs.com/package/stringify-entities

[badge-size-image]: https://img.shields.io/bundlejs/size/stringify-entities

[badge-size-url]: https://bundlejs.com/?q=stringify-entities

[esmsh]: https://esm.sh

[file-license]: license

[github-character-entities]: https://github.com/wooorm/character-entities

[github-character-entities-html4]: https://github.com/wooorm/character-entities-html4

[github-character-entities-legacy]: https://github.com/wooorm/character-entities-legacy

[github-character-reference-invalid]: https://github.com/wooorm/character-reference-invalid

[github-gist-esm]: https://gist.github.com/sindresorhus/a39789f98801d908bbc7ff3ecc99d99c

[github-parse-entities]: https://github.com/wooorm/parse-entities

[npmjs-install]: https://docs.npmjs.com/cli/install

[opensource-guide-contribute]: https://opensource.guide/how-to-contribute/

[telegraph-apos]: https://www.telegraph.co.uk/technology/advice/10516839/Why-do-some-apostrophes-get-replaced-with-andapos.html

[typescript]: https://www.typescriptlang.org

[wikipedia-xss]: https://en.wikipedia.org/wiki/Cross-site_scripting

[wooorm]: https://wooorm.com
