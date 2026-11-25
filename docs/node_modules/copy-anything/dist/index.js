/**
 * Assigns a prop the copy is keeping, preserving whether it was enumerable. Callers decide
 * _whether_ to keep a non-enumerable prop before getting here, so by this point the answer is
 * always yes.
 */
function assignProp(carry, key, newVal, originalObject) {
    if (Object.prototype.propertyIsEnumerable.call(originalObject, key)) {
        carry[key] = newVal;
        return;
    }
    Object.defineProperty(carry, key, {
        value: newVal,
        enumerable: false,
        writable: true,
        configurable: true,
    });
}
/**
 * Whether `value` is an object literal, as opposed to a class instance, a built-in like `Date`, an
 * `arguments` object, or an object with a null prototype.
 *
 * The prototype check comes first because it is cheap and rejects almost everything; the tag check
 * is only reached by values that are about to be cloned anyway.
 */
function isPlainObject(value) {
    return (Object.getPrototypeOf(value) === Object.prototype &&
        Object.prototype.toString.call(value) === '[object Object]');
}
/**
 * Returns the clone that `value` should be replaced with, and queues it to be filled in later.
 * Anything that isn't an array or a plain object is returned untouched.
 *
 * Reusing the clone already registered in `clones` is what makes circular and shared references
 * work: the placeholder is registered before its contents exist, so a reference back to an ancestor
 * resolves to it rather than walking forever.
 */
function cloneRef(value, clones, sources, dests) {
    if (typeof value !== 'object' || value === null)
        return value;
    const array = Array.isArray(value);
    if (!array && !isPlainObject(value))
        return value;
    const existing = clones.get(value);
    if (existing !== undefined)
        return existing;
    const clone = array ? new Array(value.length) : {};
    clones.set(value, clone);
    sources.push(value);
    dests.push(clone);
    return clone;
}
/**
 * Copy (clone) an object and all its props recursively to get rid of any prop referenced of the
 * original object. Arrays and the objects inside them are cloned as well.
 *
 * Nesting depth is limited only by available memory, and circular references are reproduced in the
 * copy rather than followed forever. Two props pointing at the same object still point at one
 * shared copy afterwards.
 *
 * @param target Target can be anything
 * @param [options={}] See type {@link Options} for more details.
 *
 *   - `{ props: ['key1'] }` will only copy the `key1` property. When using this you will need to cast
 *       the return type manually (in order to keep the TS implementation in here simple I didn't
 *       built a complex auto resolved type for those few cases people want to use this option)
 *   - `{ nonenumerable: true }` will copy all non-enumerable properties. Default is `{}`
 *
 * @returns The target with replaced values
 */
export function copy(target, options = {}) {
    if (typeof target !== 'object' || target === null)
        return target;
    const clones = new Map();
    // An explicit work list rather than recursion, so nesting depth is bounded by the heap instead
    // of the call stack. Deeply nested input used to throw a RangeError here.
    const sources = [];
    const dests = [];
    const result = cloneRef(target, clones, sources, dests);
    const onlyProps = Array.isArray(options.props) ? options.props : undefined;
    const nonenumerable = options.nonenumerable === true;
    while (sources.length) {
        const source = sources.pop();
        const dest = dests.pop();
        if (Array.isArray(source)) {
            for (let i = 0; i < source.length; i++) {
                if (i in source)
                    dest[i] = cloneRef(source[i], clones, sources, dests);
            }
            continue;
        }
        if (nonenumerable) {
            for (const key of Object.getOwnPropertyNames(source)) {
                // Skip __proto__ properties to prevent prototype pollution
                if (key === '__proto__')
                    continue;
                if (onlyProps && !onlyProps.includes(key))
                    continue;
                const newVal = cloneRef(source[key], clones, sources, dests);
                assignProp(dest, key, newVal, source);
            }
        }
        else {
            // `for in` skips non-enumerable keys and allocates no key array, which is why it beats
            // getOwnPropertyNames here. The own check is needed because it also walks the prototype.
            for (const key in source) {
                if (!Object.prototype.hasOwnProperty.call(source, key))
                    continue;
                // Skip __proto__ properties to prevent prototype pollution
                if (key === '__proto__')
                    continue;
                if (onlyProps && !onlyProps.includes(key))
                    continue;
                dest[key] = cloneRef(source[key], clones, sources, dests);
            }
        }
        for (const key of Object.getOwnPropertySymbols(source)) {
            if (onlyProps && !onlyProps.includes(key))
                continue;
            if (!nonenumerable && !Object.prototype.propertyIsEnumerable.call(source, key))
                continue;
            const newVal = cloneRef(source[key], clones, sources, dests);
            assignProp(dest, key, newVal, source);
        }
    }
    return result;
}
