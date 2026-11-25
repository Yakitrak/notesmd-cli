# Copy anything 🎭

<a href="https://www.npmjs.com/package/copy-anything"><img src="https://img.shields.io/npm/v/copy-anything.svg" alt="Total Downloads"></a>
<a href="https://www.npmjs.com/package/copy-anything"><img src="https://img.shields.io/npm/dw/copy-anything.svg" alt="Latest Stable Version"></a>

```
npm i copy-anything
```

An optimised way to copy'ing (cloning) an object or array. A small and simple integration.

## Motivation

I created this package because I tried a lot of similar packages that do copy'ing/cloning. But all had its quirks, and _all of them break things they are not supposed to break_... 😞

I was looking for:

- a simple copy/clone function
- has to be fast!
- everything it clones must lose any reference to the original object
- works with arrays and objects in arrays!
- supports symbols
- can copy non-enumerable props as well
- **does not break special class instances**　‼️

This last one is crucial! So many libraries use custom classes that create objects with special prototypes, and such objects all break when trying to copy them improperly. So we gotta be careful!

copy-anything will copy objects and nested properties, but only as long as they're "plain objects". As soon as a sub-prop is not a "plain object" and has a special prototype, it will copy that instance over "as is". ♻️ The flip side is that such an instance stays shared with the original, along with anything nested inside it, so mutating `copy.someDate` is visible on `original.someDate`.

There is no depth limit — nesting is only bounded by available memory. Circular references are reproduced in the copy instead of being followed forever, and two props pointing at the same object still point at one shared copy afterwards.

## Usage

```js
import { copy } from 'copy-anything'

const original = { name: 'Ditto', type: { water: true } }
const copy = copy(original)

// now if we change a nested prop like the type
copy.type.water = false
// or add a new nested prop
copy.type.fire = true

// then the original object will still be the same:
original.type.water === true // true
original.type.fire === undefined // true
```

> Please note, by default copy-anything does not copy non-enumerable props. If you need to copy those, see the instructions further down below.

### Works with arrays

It will also clone arrays, **as well as objects inside arrays!** 😉

```js
const original = [{ name: 'Squirtle' }]
const copy = copy(original)

// now if we change a prop in the array
copy[0].name = 'Wartortle'
// or add a new item to the array
copy.push({ name: 'Charmander' })

// then the original array will still be the same:
original[0].name === 'Squirtle' // true
original[1] === undefined // true
```

### Non-enumerable

By default, copy-anything only copies enumerable properties. If you also want to copy non-enumerable properties you can do so by passing that as an option.

```js
const original = { name: 'Bulbasaur' }
// bulbasaur's ID is non-enumerable
Object.defineProperty(original, 'id', {
  value: '001',
  writable: true,
  enumerable: false,
  configurable: true,
})
const copy1 = copy(original)
copy1.id === undefined // true

const copy2 = copy(original, { nonenumerable: true })
copy2.id === '001' // true
```

### Limit to specific props

You can limit to specific props.

```js
const original = { name: 'Flareon', type: ['fire'], id: '136' }
const copy = copy(original, { props: ['name'] })

copy // will look like: `{ name: 'Flareon' }`
```

> Please note, if the props you have specified are non-enumerable, you will also need to pass `{nonenumerable: true}`.

## Benchmark

Regenerate this section any time with `npm run benchmark`. Every number and every ✅ below comes from actually running the libraries, not from reading their docs.

| clone function      | flat  | mixed  | 1000 objects | 1000 keys | 1000 levels |
| ------------------- | ----- | ------ | ------------ | --------- | ----------- |
| **copy-anything**   | 194ns | 693ns  | 94.4µs       | 79.5µs    | 79.1µs      |
| copy-anything@4.0.5 | 618ns | 1117ns | 176.1µs      | 235.5µs   | 102.3µs     |
| klona               | 279ns | 300ns  | 43.4µs       | 82.4µs    | 31.2µs      |
| lodash.cloneDeep    | 391ns | 1202ns | 223.9µs      | 143.8µs   | 177.6µs     |
| structuredClone     | 787ns | 1935ns | 207.3µs      | 35.4µs    | 223.6µs     |

- **flat** — one object, 10 primitive keys
- **mixed** — one object 4 levels deep, holding an array of strings and an array of objects
- **1000 objects** — an array of 1000 flat objects
- **1000 keys** — one object with 1000 primitive keys, to isolate the cost of enumerating them
- **1000 levels** — a single-child chain 1000 objects deep

|                                  | copy-anything | copy-anything@4.0.5 | klona | lodash.cloneDeep | structuredClone |
| -------------------------------- | ------------- | ------------------- | ----- | ---------------- | --------------- |
| unlimited nesting depth          | ✅            | ❌                  | ❌    | ❌               | ❌              |
| circular references              | ✅            | ❌                  | ❌    | ✅               | ✅              |
| shared references stay shared    | ✅            | ❌                  | ❌    | ✅               | ✅              |
| symbol keys                      | ✅            | ✅                  | ❌    | ✅               | ❌              |
| non-enumerable props             | ✅            | ✅                  | ❌    | ❌               | ❌              |
| class instances left alone       | ✅            | ✅                  | ❌    | ❌               | ❌              |
| clones Date / Map / Set / RegExp | ❌            | ❌                  | ✅    | ✅               | ✅              |
| survives functions in the input  | ✅            | ✅                  | ✅    | ✅               | ❌              |

_Measured on Apple M3 Pro, Node v24.19.0. Lower is better._

### Why klona is faster in some metrics

klona is the fastest of the bunch and that is not an accident — it is faster **because** it does less. It recurses with a plain `for in` loop and keeps no record of what it has already visited, which is exactly what buys it the speed. Those same two choices are what the ❌ column above is made of.

If your objects are plain data, shallow, and acyclic, klona is a great pick and you should use it. Here is what changes when they are not.

**Deeply nested data.** Any clone function that recurses is bounded by the JavaScript call stack, so it dies on deep input. copy-anything walks an explicit work list on the heap instead, so it has no such limit.

```js
let deep = { value: 'leaf' }
for (let i = 0; i < 10000; i++) deep = { nested: deep }

klona(deep) // 💥 RangeError: Maximum call stack size exceeded
copy(deep) // ✅ fine — depth is only limited by memory
```

**Circular references.** A cycle is just infinite depth, so an unvisited-set-free clone follows it until the stack runs out. copy-anything registers each clone before filling it in, so a reference back up the tree resolves to the copy already being built.

```js
const user = { name: 'Luca' }
user.self = user

klona(user) // 💥 RangeError: Maximum call stack size exceeded

const copied = copy(user)
copied.self === copied // ✅ true — the cycle is reproduced in the copy
copied !== user // ✅ true — and it is still a real copy
```

**Two props pointing at one object.** Without a record of what it has cloned, klona cannot know the two props were the same object, so the copy ends up with two separate ones.

```js
const shared = { count: 0 }
const original = { a: shared, b: shared }

const klonad = klona(original)
klonad.a === klonad.b // ❌ false — silently split into two separate objects

const copied = copy(original)
copied.a === copied.b // ✅ true — still one shared object
copied.a !== shared // ✅ true — and detached from the original
```

**Class instances.** klona rebuilds them with `new x.constructor()`, which throws outright if the constructor requires arguments, and silently skips any prop the constructor sets to the same value. copy-anything copies such instances over as is.

```js
class User {
  constructor(id) {
    if (id === undefined) throw new Error('id is required')
    this.id = id
  }
}
const original = { user: new User(1) }

klona(original) // 💥 Error: id is required
copy(original).user === original.user // ✅ true — copied over as is
```

**Symbols and non-enumerable props.** klona's `for in` loop cannot see either one, so both are dropped from the copy. copy-anything copies symbol keys by default and non-enumerable props with `{ nonenumerable: true }`.

The trade goes the other way in one place, and it is worth knowing: klona clones `Date`, `Map`, `Set`, `RegExp` and typed arrays, where copy-anything deliberately copies those over as is. If you need those cloned, klona (or `structuredClone`) is the better tool.

## Meet the family (more tiny utils with TS support)

- [is-what 🙉](https://github.com/mesqueeb/is-what)
- [is-where 🙈](https://github.com/mesqueeb/is-where)
- [merge-anything 🥡](https://github.com/mesqueeb/merge-anything)
- [check-anything 👁](https://github.com/mesqueeb/check-anything)
- [remove-anything ✂️](https://github.com/mesqueeb/remove-anything)
- [getorset-anything 🐊](https://github.com/mesqueeb/getorset-anything)
- [map-anything 🗺](https://github.com/mesqueeb/map-anything)
- [filter-anything ⚔️](https://github.com/mesqueeb/filter-anything)
- [copy-anything 🎭](https://github.com/mesqueeb/copy-anything)
- [case-anything 🐫](https://github.com/mesqueeb/case-anything)
- [flatten-anything 🏏](https://github.com/mesqueeb/flatten-anything)
- [nestify-anything 🧅](https://github.com/mesqueeb/nestify-anything)

## Source code

copy-anything has zero dependencies and the whole implementation is one short file you can read in a minute: [src/index.ts](src/index.ts).

It walks the object with an explicit work list rather than by recursing, which is what lets it handle any nesting depth — a recursive clone throws `RangeError: Maximum call stack size exceeded` somewhere around 2500 levels, and so does every other clone function that recurses, including the browser's own `structuredClone`.

Each clone is registered before it is filled in, so a reference pointing back up the tree resolves to the clone that is already being built instead of being followed forever. That is what makes circular references work.
