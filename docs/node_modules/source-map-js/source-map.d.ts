export interface StartOfSourceMap {
  file?: string;
  sourceRoot?: string;
}

/**
 * A source map as `SourceMapGenerator#toJSON` returns it. For a source map
 * passed to `SourceMapConsumer`, use `RawSourceMapInput`.
 */
export interface RawSourceMap extends StartOfSourceMap {
  /** The number `3` at runtime. Typed as a string for backward compatibility. */
  version: string;
  sources: string[];
  names: string[];
  /**
   * `null` for each source without content at runtime. Typed as `string[]`
   * for backward compatibility.
   */
  sourcesContent?: string[];
  mappings: string;
}

/**
 * A source map as `SourceMapConsumer` accepts it. Unlike `RawSourceMap`, the
 * `version` can be a number, `sourcesContent` can have `null` entries, and
 * `names` can be left out.
 */
export interface RawSourceMapInput extends StartOfSourceMap {
  version: number | string;
  sources: string[];
  names?: string[];
  sourcesContent?: (string | null)[];
  mappings: string;
}

export interface RawSection {
  /** Where the section starts in the generated code. Both are 0-based. */
  offset: Position;
  map: RawSourceMapInput | RawIndexMap;
}

/** A source map made of sections, each with its own source map. */
export interface RawIndexMap {
  version: number | string;
  file?: string;
  sections: RawSection[];
}

export interface Position {
  line: number;
  column: number;
}

export interface LineRange extends Position {
  lastColumn: number;
}

export interface FindPosition extends Position {
  // SourceMapConsumer.GREATEST_LOWER_BOUND or SourceMapConsumer.LEAST_UPPER_BOUND
  bias?: number;
}

export interface SourceFindPosition extends FindPosition {
  source: string;
}

export interface MappedPosition extends Position {
  source: string;
  name?: string;
}

export interface MappingItem {
  source: string | null;
  generatedLine: number;
  generatedColumn: number;
  originalLine: number | null;
  originalColumn: number | null;
  name: string | null;
}

export class SourceMapConsumer {
  static GENERATED_ORDER: number;
  static ORIGINAL_ORDER: number;

  static GREATEST_LOWER_BOUND: number;
  static LEAST_UPPER_BOUND: number;

  /**
   * @param rawSourceMap A source map or an index map, as an object or as a
   *   JSON string.
   * @param sourceMapUrl The URL of the source map. Relative `sources` are
   *   resolved against it.
   */
  constructor(
    rawSourceMap: RawSourceMapInput | RawIndexMap | string,
    sourceMapUrl?: string
  );
  static fromSourceMap(
    sourceMap: SourceMapGenerator,
    sourceMapUrl?: string
  ): SourceMapConsumer;
  readonly file: string | undefined | null;
  readonly sourceRoot: string | undefined | null;
  /** Can have `null` entries at runtime: one for each source without content. */
  readonly sourcesContent: readonly string[] | null | undefined;
  readonly sources: readonly string[];

  computeColumnSpans(): void;
  /**
   * At runtime, `source`, `line`, `column` and `name` are all `null` when no
   * mapping is found, and `name` is `null` when the mapping has no name.
   */
  originalPositionFor(generatedPosition: FindPosition): MappedPosition;
  /**
   * At runtime, `line`, `column` and `lastColumn` are all `null` when no
   * mapping is found, and `lastColumn` is `null` until `computeColumnSpans()`
   * is called.
   */
  generatedPositionFor(originalPosition: SourceFindPosition): LineRange;
  /**
   * Without a `column`, returns the mappings of the line, or of the next line
   * that has mappings. `lastColumn` is `null` until `computeColumnSpans()` is
   * called.
   */
  allGeneratedPositionsFor(
    originalPosition:
      MappedPosition | { source: string; line: number; column?: number }
  ): Array<Position & { lastColumn?: number | null }>;
  hasContentsOfAllSources(): boolean;
  sourceContentFor(
    source: string,
    returnNullOnMissing?: boolean
  ): string | null;
  eachMapping(
    callback: (mapping: MappingItem) => void,
    context?: any,
    order?: number
  ): void;
}

export interface Mapping {
  generated: Position;
  original?: Position | null;
  source?: string | null;
  name?: string | null;
}

export interface SourceMapGeneratorOptions extends StartOfSourceMap {
  /** Don't validate mappings in `addMapping`. */
  skipValidation?: boolean;
  /** Warn about invalid mappings and skip them, instead of throwing. */
  ignoreInvalidMapping?: boolean;
}

export class SourceMapGenerator {
  constructor(startOfSourceMap?: SourceMapGeneratorOptions);
  static fromSourceMap(
    sourceMapConsumer: SourceMapConsumer,
    startOfSourceMap?: SourceMapGeneratorOptions
  ): SourceMapGenerator;
  addMapping(mapping: Mapping): void;
  setSourceContent(
    sourceFile: string,
    sourceContent: string | null | undefined
  ): void;
  applySourceMap(
    sourceMapConsumer: SourceMapConsumer,
    sourceFile?: string,
    sourceMapPath?: string
  ): void;
  toString(): string;
  toJSON(): RawSourceMap;
}

export interface CodeWithSourceMap {
  code: string;
  map: SourceMapGenerator;
}

/** A string of generated code, a `SourceNode`, or an array of chunks. */
export type SourceNodeChunk = string | SourceNode | SourceNodeChunkArray;
export interface SourceNodeChunkArray extends Array<SourceNodeChunk> {}

export class SourceNode {
  constructor();
  constructor(
    line: number | null,
    column: number | null,
    source: string | null,
    chunks?: SourceNodeChunk,
    name?: string | null
  );
  static fromStringWithSourceMap(
    code: string,
    sourceMapConsumer: SourceMapConsumer,
    relativePath?: string
  ): SourceNode;
  /** Returns this node at runtime. */
  add(chunk: SourceNodeChunk): void;
  /** Returns this node at runtime. */
  prepend(chunk: SourceNodeChunk): void;
  setSourceContent(sourceFile: string, sourceContent: string): void;
  /**
   * At runtime, `source`, `line`, `column` and `name` are `null` for chunks
   * without an original position, and `name` is `null` when there's no name.
   */
  walk(fn: (chunk: string, mapping: MappedPosition) => void): void;
  walkSourceContents(fn: (file: string, content: string) => void): void;
  join(sep: string): SourceNode;
  replaceRight(
    pattern: string | RegExp,
    replacement: string | ((substring: string, ...args: any[]) => string)
  ): SourceNode;
  toString(): string;
  toStringWithSourceMap(
    startOfSourceMap?: SourceMapGeneratorOptions
  ): CodeWithSourceMap;
}
