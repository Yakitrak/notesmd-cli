/**
 * Encode certain characters in `value`.
 *
 * @param {string} value
 *   Value to encode.
 * @param {CoreWithFormatOptions} options
 *   Configuration.
 * @returns {string}
 *   Encoded value.
 */
export function core(value: string, options: CoreWithFormatOptions): string;
/**
 * Configuration for `stringifyEntitiesLight`.
 */
export type LightOptions = {
    /**
     * Whether to only escape possibly dangerous characters (default: `false`);
     * those characters are `"`, `&`, `'`, `<`, `>`, and `` ` ``.
     */
    escapeOnly?: boolean | null | undefined;
    /**
     * Whether to only escape the given subset of characters (optional);
     * only BMP characters are supported here (so no emoji).
     */
    subset?: ReadonlyArray<string> | null | undefined;
};
/**
 * Internal format options.
 */
export type FormatOptions = {
    /**
     *   Format strategy.
     */
    format: (code: number, next: number, options: CoreWithFormatOptions) => string;
};
export type CoreWithFormatOptions = LightOptions & FormatOptions & FormatSmartOptions;
import type { FormatSmartOptions } from './util/format-smart.js';
//# sourceMappingURL=core.d.ts.map