/**
 * Encode special characters in `value`.
 *
 * @param {string} value
 *   Value to encode.
 * @param {Options | null | undefined} [options]
 *   Configuration.
 * @returns {string}
 *   Encoded value.
 */
export function stringifyEntities(value: string, options?: Options | null | undefined): string;
/**
 * Encode special characters in `value` as hexadecimal character references.
 *
 * @param {string} value
 *   Value to encode.
 * @param {LightOptions | null | undefined} [options]
 *   Configuration.
 * @returns {string}
 *   Encoded value.
 */
export function stringifyEntitiesLight(value: string, options?: LightOptions | null | undefined): string;
/**
 * Configuration for `stringifyEntities`.
 *
 * Includes the fields of `LightOptions`.
 */
export type Options = LightOptions & FormatSmartOptions;
import type { LightOptions } from './core.js';
import type { FormatSmartOptions } from './util/format-smart.js';
//# sourceMappingURL=index.d.ts.map