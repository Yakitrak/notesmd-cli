type Headers = Record<string, string>;
type QueryParameters = Record<string, any>;
/**
 * The method of the request.
 */
type Method = 'DELETE' | 'GET' | 'PATCH' | 'POST' | 'PUT';
type Request = {
    method: Method;
    /**
     * The path of the REST API to send the request to.
     */
    path: string;
    queryParameters: QueryParameters;
    data?: Array<Record<string, any>> | Record<string, any> | undefined;
    headers: Headers;
    /**
     * If the given request should persist on the cache. Keep in mind,
     * that some methods may have this option enabled by default.
     */
    cacheable?: boolean | undefined;
    /**
     * Some POST methods in the Algolia REST API uses the `read` transporter.
     * This information is defined at the spec level.
     */
    useReadTransporter?: boolean | undefined;
};
type EndRequest = Pick<Request, 'headers' | 'method'> & {
    /**
     * The full URL of the REST API.
     */
    url: string;
    /**
     * The connection timeout, in milliseconds.
     */
    connectTimeout: number;
    /**
     * The response timeout, in milliseconds.
     */
    responseTimeout: number;
    data?: string | Uint8Array | undefined;
};
type Response = {
    /**
     * The body of the response.
     */
    content: string;
    /**
     * The headers of the response, with lower-cased names. Optional so that custom requesters that
     * don't capture headers remain valid, in which case the `WithHTTPInfo` methods expose no headers.
     */
    headers?: Headers | undefined;
    /**
     * Whether the API call is timed out or not.
     */
    isTimedOut: boolean;
    /**
     * The HTTP status code of the response.
     */
    status: number;
};
type Requester = {
    /**
     * Sends the given `request` to the server.
     */
    send: (request: EndRequest) => Promise<Response>;
    /**
     * Sends the given `request` and returns a raw byte stream for streaming responses (e.g. SSE).
     */
    sendStream?: (request: EndRequest) => Promise<ReadableStream<Uint8Array>>;
};

/**
 * The full HTTP response of an API call, as returned by the `WithHTTPInfo` variant of each method.
 */
type AlgoliaHttpResponse<TData> = {
    /**
     * The HTTP status code of the response.
     */
    status: number;
    /**
     * The headers of the response, with lower-cased names. Undefined when the requester does not capture headers.
     */
    headers?: Headers | undefined;
    /**
     * The raw body of the response.
     */
    content: string;
    /**
     * The deserialized body of the response.
     */
    data: TData;
};

type Cache = {
    /**
     * Gets the value of the given `key`.
     */
    get: <TValue>(key: Record<string, any> | string, defaultValue: () => Promise<TValue>, events?: CacheEvents<TValue> | undefined) => Promise<TValue>;
    /**
     * Sets the given value with the given `key`.
     */
    set: <TValue>(key: Record<string, any> | string, value: TValue) => Promise<TValue>;
    /**
     * Deletes the given `key`.
     */
    delete: (key: Record<string, any> | string) => Promise<void>;
    /**
     * Clears the cache.
     */
    clear: () => Promise<void>;
};
type CacheEvents<TValue> = {
    /**
     * The callback when the given `key` is missing from the cache.
     */
    miss: (value: TValue) => Promise<any>;
};
type MemoryCacheOptions = {
    /**
     * If keys and values should be serialized using `JSON.stringify`.
     */
    serializable?: boolean | undefined;
};
type BrowserLocalStorageOptions = {
    /**
     * The cache key.
     */
    key: string;
    /**
     * The time to live for each cached item in seconds.
     */
    timeToLive?: number | undefined;
    /**
     * The native local storage implementation.
     */
    localStorage?: Storage | undefined;
};
type BrowserLocalStorageCacheItem = {
    /**
     * The cache item creation timestamp.
     */
    timestamp: number;
    /**
     * The cache item value.
     */
    value: any;
};
type FallbackableCacheOptions = {
    /**
     * List of caches order by priority.
     */
    caches: Cache[];
};

/**
 * Shared configuration for chunked helpers that poll for task completion.
 */
type ChunkedHelperOptions = {
    /**
     * The maximum number of retries when polling for task completion. 100 by default.
     */
    maxRetries?: number | undefined;
};

/**
 * WHATWG-compliant Server-Sent Events parser.
 *
 * Three-layer architecture:
 *   1. iterLines()    — byte chunking → line decoding
 *   2. SSEDecoder     — line → SSE event decoding
 *   3. iterSSEEvents  — top-level composer (exported)
 *
 * @see https://html.spec.whatwg.org/multipage/server-sent-events.html#event-stream-interpretation
 */
type ServerSentEvent = {
    /** Concatenated data: field values, joined by '\n'. */
    data: string;
    /** Event type from the event: field. Defaults to "" (empty string). */
    event: string;
    /** Last event ID. Persists across dispatches until changed. */
    id: string | null;
    /** Reconnection time in ms. Persists across dispatches until changed. */
    retry: number | null;
};
/**
 * Wrapper for a parsed SSE event, yielded by the typed `*Stream` methods.
 *
 * - `data` is the JSON-parsed payload when parsing succeeds, `null` otherwise.
 * - `raw` is the original {@link ServerSentEvent} (always present).
 * - `error` is set when JSON parsing of `event.data` failed.
 */
type StreamEvent<T = Record<string, unknown>> = {
    /** Parsed data from the event, or `null` if parsing failed. */
    data: T | null;
    /** The original, unparsed SSE event. */
    raw: ServerSentEvent;
    /** The error that occurred while parsing `event.data`, if any. */
    error?: Error;
};
/**
 * Parses a byte stream as WHATWG Server-Sent Events.
 *
 * Accepts both ReadableStream<Uint8Array> (browser fetch) and
 * AsyncIterable<Uint8Array> (Node.js streams / Buffer chunks).
 *
 * @example
 * ```ts
 * const response = await fetch(url);
 * for await (const event of iterSSEEvents(response.body!)) {
 *   console.log(event.event, event.data);
 * }
 * ```
 */
declare function iterSSEEvents(stream: ReadableStream<Uint8Array> | AsyncIterable<Uint8Array>): AsyncGenerator<ServerSentEvent>;

type Host = {
    /**
     * The host URL.
     */
    url: string;
    /**
     * The accepted transporter.
     */
    accept: 'read' | 'readWrite' | 'write';
    /**
     * The protocol of the host URL.
     */
    protocol: 'http' | 'https';
    /**
     * The port of the host URL.
     */
    port?: number | undefined;
};
type StatefulHost = Host & {
    /**
     * The status of the host.
     */
    status: 'down' | 'timed out' | 'up';
    /**
     * The last update of the host status, used to compare with the expiration delay.
     */
    lastUpdate: number;
    /**
     * Returns whether the host is up or not.
     */
    isUp: () => boolean;
    /**
     * Returns whether the host is timed out or not.
     */
    isTimedOut: () => boolean;
};

declare const LogLevelEnum: Readonly<Record<string, LogLevelType>>;
type LogLevelType = 1 | 2 | 3;
type Logger = {
    /**
     * Logs debug messages.
     */
    debug: (message: string, args?: any | undefined) => Promise<void>;
    /**
     * Logs info messages.
     */
    info: (message: string, args?: any | undefined) => Promise<void>;
    /**
     * Logs warning messages. When not implemented, warnings fall back to `console.warn` so they stay visible with the default no-op logger.
     */
    warn?: (message: string, args?: any | undefined) => Promise<void>;
    /**
     * Logs error messages.
     */
    error: (message: string, args?: any | undefined) => Promise<void>;
};

type RequestOptions = Pick<Request, 'cacheable'> & {
    /**
     * Custom timeout for the request. Note that, in normal situations
     * the given timeout will be applied. But the transporter layer may
     * increase this timeout if there is need for it.
     */
    timeouts?: Partial<Timeouts> | undefined;
    /**
     * Custom headers for the request. This headers are
     * going to be merged the transporter headers.
     */
    headers?: Headers | undefined;
    /**
     * Custom query parameters for the request. This query parameters are
     * going to be merged the transporter query parameters.
     */
    queryParameters?: QueryParameters | undefined;
    /**
     * Custom data for the request. This data is
     * going to be merged the transporter data.
     */
    data?: Array<Record<string, any>> | Record<string, any> | undefined;
};
type StackFrame = {
    request: EndRequest;
    response: Response;
    host: Host;
    triesLeft: number;
};
type AlgoliaAgentOptions = {
    /**
     * The segment. Usually the integration name.
     */
    segment: string;
    /**
     * The version. Usually the integration version.
     */
    version?: string | undefined;
};
type AlgoliaAgent = {
    /**
     * The raw value of the user agent.
     */
    value: string;
    /**
     * Mutates the current user agent adding the given user agent options.
     */
    add: (options: AlgoliaAgentOptions) => AlgoliaAgent;
};
type Timeouts = {
    /**
     * Timeout in milliseconds before the connection is established.
     */
    connect: number;
    /**
     * Timeout in milliseconds before reading the response on a read request.
     */
    read: number;
    /**
     * Timeout in milliseconds before reading the response on a write request.
     */
    write: number;
};
type TransporterOptions = {
    /**
     * The cache of the hosts. Usually used to persist
     * the state of the host when its down.
     */
    hostsCache: Cache;
    /**
     * The logger instance to send events of the transporter.
     */
    logger: Logger;
    /**
     * The underlying requester used. Should differ
     * depending of the environment where the client
     * will be used.
     */
    requester: Requester;
    /**
     * Cache used to store in-flight requests.
     * When a request is marked as `cacheable`, its returned Promise
     * is stored in this cache so that identical requests can share
     * the same Promise before it resolves.
     *
     * @warning
     * The provided cache **must not** serialize stored values.
     *
     * Since in-flight requests are stored as Promises (which cannot be
     * serialized to JSON), using a serializing cache will cause failures.
     *
     * Make sure to use a non-serializing cache implementation, such as `createMemoryCache({ serializable: false })` or to disable request caching with `createNullCache()`
     */
    requestsCache: Cache;
    /**
     * The cache of the responses. When requests are
     * `cacheable`, the returned responses persists
     * in this cache to shared in similar requests.
     */
    responsesCache: Cache;
    /**
     * The timeouts used by the requester. The transporter
     * layer may increase this timeouts as defined on the
     * retry strategy.
     */
    timeouts: Timeouts;
    /**
     * How many times to wait and retry on the same host after HTTP 429.
     * Default is 3. `0` fails on the first 429 (no wait).
     * Wait time is `Retry-After` in whole seconds, or 1 second if the header is missing or invalid.
     */
    maxRateLimitRetries?: number | undefined;
    /**
     * The hosts used by the requester.
     */
    hosts: Host[];
    /**
     * The headers used by the requester. The transporter
     * layer may add some extra headers during the request
     * for the user agent, and others.
     */
    baseHeaders: Headers;
    /**
     * The query parameters used by the requester. The transporter
     * layer may add some extra headers during the request
     * for the user agent, and others.
     */
    baseQueryParameters: QueryParameters;
    /**
     * The user agent used. Sent on query parameters.
     */
    algoliaAgent: AlgoliaAgent;
    /**
     * An optional function to compress request bodies before sending.
     * When provided, POST/PUT bodies exceeding the compression threshold
     * will be compressed and `Content-Encoding: gzip` is added to the headers.
     * Node builds use node:zlib, browser/worker builds use CompressionStream when available.
     */
    compress?: (data: string) => Promise<Uint8Array>;
    compression?: 'gzip';
    /**
     * Where the generated Request-ID is sent: as the `Request-ID` header, or as the
     * `x-algolia-request-id` query parameter. When undefined, no Request-ID is sent.
     * A caller-supplied Request-ID is never overwritten.
     */
    requestIdChannel?: 'headers' | 'queryParameters' | undefined;
};
type Transporter = TransporterOptions & {
    /**
     * Performs a request.
     * The `baseRequest` and `baseRequestOptions` will be merged accordingly.
     */
    request: <TResponse>(baseRequest: Request, baseRequestOptions?: RequestOptions) => Promise<TResponse>;
    requestStream: (baseRequest: Request, baseRequestOptions?: RequestOptions) => AsyncGenerator<ServerSentEvent>;
};
/**
 * The transporter returned by `createTransporter`. Kept separate from `Transporter` so that
 * existing implementations of `Transporter` remain type-valid.
 */
type TransporterWithHttpInfo = Transporter & {
    /**
     * Performs a request and returns the full HTTP response information — status code,
     * headers (when the requester captures them), raw body and deserialized data.
     * Both the requests and the responses caches are bypassed: it always hits the network.
     */
    requestWithHttpInfo: <TData>(baseRequest: Request, baseRequestOptions?: RequestOptions) => Promise<AlgoliaHttpResponse<TData>>;
};

type AuthMode = 'WithinHeaders' | 'WithinQueryParameters';
type OverriddenTransporterOptions = 'baseHeaders' | 'baseQueryParameters' | 'hosts';
type CreateClientOptions = Omit<TransporterOptions, OverriddenTransporterOptions | 'algoliaAgent'> & Partial<Pick<TransporterOptions, OverriddenTransporterOptions>> & {
    appId: string;
    apiKey: string;
    authMode?: AuthMode | undefined;
    algoliaAgents: AlgoliaAgentOptions[];
};
type ClientOptions = Partial<Omit<CreateClientOptions, 'apiKey' | 'appId' | 'compress'>>;

type IterableOptions<TResponse> = Partial<{
    /**
     * The function that runs right after the API call has been resolved, allows you to do anything with the response before `validate`.
     */
    aggregator: (response: TResponse) => unknown | PromiseLike<unknown>;
    /**
     * The `validate` condition to throw an error and its message.
     */
    error: {
        /**
         * The function to validate the error condition.
         */
        validate: (response: TResponse) => boolean | PromiseLike<boolean>;
        /**
         * The error message to throw.
         */
        message: (response: TResponse) => string | PromiseLike<string>;
    };
    /**
     * The function to decide how long to wait between iterations.
     */
    timeout: () => number | PromiseLike<number>;
}>;
type CreateIterablePromise<TResponse> = IterableOptions<TResponse> & {
    /**
     * The function to run, which returns a promise.
     *
     * The `previousResponse` parameter (`undefined` on the first call) allows you to build your request with incremental logic, to iterate on `page` or `cursor` for example.
     */
    func: (previousResponse?: TResponse | undefined) => Promise<TResponse>;
    /**
     * The validator function. It receive the resolved return of the API call.
     */
    validate: (response: TResponse) => boolean | PromiseLike<boolean>;
};

declare function createBrowserLocalStorageCache(options: BrowserLocalStorageOptions): Cache;

declare function createFallbackableCache(options: FallbackableCacheOptions): Cache;

declare function createMemoryCache(options?: MemoryCacheOptions): Cache;

declare function createNullCache(): Cache;

declare const DEFAULT_CONNECT_TIMEOUT_BROWSER = 1000;
declare const DEFAULT_READ_TIMEOUT_BROWSER = 2000;
declare const DEFAULT_WRITE_TIMEOUT_BROWSER = 30000;
declare const DEFAULT_CONNECT_TIMEOUT_NODE = 2000;
declare const DEFAULT_READ_TIMEOUT_NODE = 5000;
declare const DEFAULT_WRITE_TIMEOUT_NODE = 30000;
declare const DEFAULT_REPLACE_ALL_OBJECTS_MAX_RETRIES = 800;

declare function createAlgoliaAgent(version: string): AlgoliaAgent;

declare function createAuth(appId: string, apiKey: string, authMode?: AuthMode): {
    readonly headers: () => Headers;
    readonly queryParameters: () => QueryParameters;
};

/**
 * Helper: Returns the promise of a given `func` to iterate on, based on a given `validate` condition.
 *
 * @param createIterator - The createIterator options.
 * @param createIterator.func - The function to run, which returns a promise.
 * @param createIterator.validate - The validator function. It receives the resolved return of `func`.
 * @param createIterator.aggregator - The function that runs right after the `func` method has been executed, allows you to do anything with the response before `validate`.
 * @param createIterator.error - The `validate` condition to throw an error, and its message.
 * @param createIterator.timeout - The function to decide how long to wait between iterations.
 */
declare function createIterablePromise<TResponse>({ func, validate, aggregator, error, timeout, }: CreateIterablePromise<TResponse>): Promise<TResponse>;

type GetAlgoliaAgent = {
    algoliaAgents: AlgoliaAgentOptions[];
    client: string;
    version: string;
};
declare function getAlgoliaAgent({ algoliaAgents, client, version }: GetAlgoliaAgent): AlgoliaAgent;

declare function createNullLogger(): Logger;

/**
 * Prefers the logger's optional `warn` method and falls back to `console.warn`, so warnings stay visible with the default no-op logger.
 */
declare function logWarning(logger: Logger, message: string): void;

declare const COMPRESSION_THRESHOLD = 750;

declare function createStatefulHost(host: Host, status?: StatefulHost['status']): StatefulHost;

declare function createTransporter({ hosts, hostsCache, baseHeaders, logger, baseQueryParameters, algoliaAgent, timeouts, requester, requestsCache, responsesCache, compress, compression, requestIdChannel, maxRateLimitRetries, }: TransporterOptions): TransporterWithHttpInfo;

declare class AlgoliaError extends Error {
    name: string;
    constructor(message: string, name: string);
}
declare class IndexNotFoundError extends AlgoliaError {
    constructor(indexName: string);
}
declare class IndicesInSameAppError extends AlgoliaError {
    constructor();
}
declare class IndexAlreadyExistsError extends AlgoliaError {
    constructor(indexName: string);
}
declare class ErrorWithStackTrace extends AlgoliaError {
    stackTrace: StackFrame[];
    correlationId?: string | undefined;
    constructor(message: string, stackTrace: StackFrame[], name: string, correlationId?: string | undefined);
}
declare class RetryError extends ErrorWithStackTrace {
    constructor(stackTrace: StackFrame[], correlationId?: string | undefined);
}
declare class ApiError extends ErrorWithStackTrace {
    status: number;
    constructor(message: string, status: number, stackTrace: StackFrame[], name?: string, correlationId?: string | undefined);
}
/**
 * Thrown by `sendStream` when the HTTP status is not 2xx, so `requestStream` can
 * wait and retry on 429 using `Retry-After` without parsing the error message.
 */
declare class StreamRequestError extends AlgoliaError {
    status: number;
    headers?: Headers | undefined;
    constructor(status: number, body: string, headers?: Headers | undefined);
}
declare class DeserializationError extends AlgoliaError {
    response: Response;
    correlationId?: string | undefined;
    constructor(message: string, response: Response, correlationId?: string | undefined);
}
type DetailedErrorWithMessage = {
    message: string;
    label: string;
};
type DetailedErrorWithTypeID = {
    id: string;
    type: string;
    name?: string | undefined;
};
type DetailedError = {
    code: string;
    details?: DetailedErrorWithMessage[] | DetailedErrorWithTypeID[] | undefined;
};
declare class DetailedApiError extends ApiError {
    error: DetailedError;
    constructor(message: string, status: number, error: DetailedError, stackTrace: StackFrame[], correlationId?: string | undefined);
}

declare function shuffle<TData>(array: TData[]): TData[];
declare function serializeUrl(host: Host, path: string, queryParameters: QueryParameters): string;
declare function serializeQueryParameters(parameters: QueryParameters): string;
declare function serializeData(request: Request, requestOptions: RequestOptions): string | undefined;
declare function serializeHeaders(baseHeaders: Headers, requestHeaders: Headers, requestOptionsHeaders?: Headers | undefined): Headers;
declare function deserializeSuccess<TObject>(response: Response): TObject;
declare function deserializeSuccessWithHttpInfo<TData>(response: Response): AlgoliaHttpResponse<TData>;
declare function getCorrelationId(headers: Headers | undefined): string | undefined;
declare function getLastCorrelationId(stackTrace: StackFrame[]): string | undefined;
declare function deserializeFailure(response: Response, stackFrame: StackFrame[]): Error;

declare function generateRequestId(): string;
/**
 * Returns request options carrying a Request-ID on the transporter's channel, so that every
 * request a multi-request helper performs shares the same ID. Returns the options unchanged
 * when the channel is off or when a Request-ID is already supplied. Never apply ahead of a
 * cacheable operation: the ID enters the cache key and defeats caching and deduplication.
 */
declare function withRequestId(transporter: Pick<Transporter, 'requestIdChannel' | 'baseHeaders' | 'baseQueryParameters'>, requestOptions?: RequestOptions | undefined): RequestOptions | undefined;

declare function isNetworkError({ isTimedOut, status }: Omit<Response, 'content'>): boolean;
declare function isRetryable({ isTimedOut, status }: Omit<Response, 'content'>): boolean;
declare function isSuccess({ status }: Pick<Response, 'status'>): boolean;
declare function isRateLimited({ status }: Pick<Response, 'status'>): boolean;
/**
 * 429 from `sendStream`: a status-bearing error, or the `HTTP 429:` string requesters used to throw.
 */
declare function isRateLimitedError(error: unknown): boolean;
declare function headersFromError(error: unknown): Headers | undefined;
/**
 * `Retry-After` as a wait in milliseconds.
 * Only a positive whole-number-of-seconds string is honored; anything else (missing, `0`, HTTP-date, junk) waits 1s.
 */
declare function parseRetryAfterMs(headers: Headers | undefined): number;

declare function stackTraceWithoutCredentials(stackTrace: StackFrame[]): StackFrame[];
declare function stackFrameWithoutCredentials(stackFrame: StackFrame): StackFrame;

declare function validateRequired(field: string, method: string, value: unknown): void;

export { type AlgoliaAgent, type AlgoliaAgentOptions, AlgoliaError, type AlgoliaHttpResponse, ApiError, type AuthMode, type BrowserLocalStorageCacheItem, type BrowserLocalStorageOptions, COMPRESSION_THRESHOLD, type Cache, type CacheEvents, type ChunkedHelperOptions, type ClientOptions, type CreateClientOptions, type CreateIterablePromise, DEFAULT_CONNECT_TIMEOUT_BROWSER, DEFAULT_CONNECT_TIMEOUT_NODE, DEFAULT_READ_TIMEOUT_BROWSER, DEFAULT_READ_TIMEOUT_NODE, DEFAULT_REPLACE_ALL_OBJECTS_MAX_RETRIES, DEFAULT_WRITE_TIMEOUT_BROWSER, DEFAULT_WRITE_TIMEOUT_NODE, DeserializationError, DetailedApiError, type DetailedError, type DetailedErrorWithMessage, type DetailedErrorWithTypeID, type EndRequest, ErrorWithStackTrace, type FallbackableCacheOptions, type GetAlgoliaAgent, type Headers, type Host, IndexAlreadyExistsError, IndexNotFoundError, IndicesInSameAppError, type IterableOptions, LogLevelEnum, type LogLevelType, type Logger, type MemoryCacheOptions, type Method, type QueryParameters, type Request, type RequestOptions, type Requester, type Response, RetryError, type ServerSentEvent, type StackFrame, type StatefulHost, type StreamEvent, StreamRequestError, type Timeouts, type Transporter, type TransporterOptions, type TransporterWithHttpInfo, createAlgoliaAgent, createAuth, createBrowserLocalStorageCache, createFallbackableCache, createIterablePromise, createMemoryCache, createNullCache, createNullLogger, createStatefulHost, createTransporter, deserializeFailure, deserializeSuccess, deserializeSuccessWithHttpInfo, generateRequestId, getAlgoliaAgent, getCorrelationId, getLastCorrelationId, headersFromError, isNetworkError, isRateLimited, isRateLimitedError, isRetryable, isSuccess, iterSSEEvents, logWarning, parseRetryAfterMs, serializeData, serializeHeaders, serializeQueryParameters, serializeUrl, shuffle, stackFrameWithoutCredentials, stackTraceWithoutCredentials, validateRequired, withRequestId };
