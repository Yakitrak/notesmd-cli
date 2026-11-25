"use strict";
var __defProp = Object.defineProperty;
var __getOwnPropDesc = Object.getOwnPropertyDescriptor;
var __getOwnPropNames = Object.getOwnPropertyNames;
var __hasOwnProp = Object.prototype.hasOwnProperty;
var __export = (target, all) => {
  for (var name in all)
    __defProp(target, name, { get: all[name], enumerable: true });
};
var __copyProps = (to, from, except, desc) => {
  if (from && typeof from === "object" || typeof from === "function") {
    for (let key of __getOwnPropNames(from))
      if (!__hasOwnProp.call(to, key) && key !== except)
        __defProp(to, key, { get: () => from[key], enumerable: !(desc = __getOwnPropDesc(from, key)) || desc.enumerable });
  }
  return to;
};
var __toCommonJS = (mod) => __copyProps(__defProp({}, "__esModule", { value: true }), mod);

// src/index.ts
var src_exports = {};
__export(src_exports, {
  createFetchRequester: () => createFetchRequester
});
module.exports = __toCommonJS(src_exports);

// src/createFetchRequester.ts
var import_client_common = require("@algolia/client-common");
function isAbortError(error) {
  return error instanceof Error && error.name === "AbortError";
}
function getErrorMessage(error, abortContent) {
  if (isAbortError(error)) {
    return abortContent;
  }
  return error instanceof Error ? error.message : "Network request failed";
}
function toResponseHeaders(fetchHeaders) {
  const headers = {};
  fetchHeaders.forEach((value, name) => {
    headers[name] = value;
  });
  return headers;
}
function createFetchRequester({ requesterOptions = {} } = {}) {
  async function send(request) {
    const abortController = new AbortController();
    const signal = abortController.signal;
    const createTimeout = (timeout) => {
      return setTimeout(() => {
        abortController.abort();
      }, timeout);
    };
    const connectTimeout = createTimeout(request.connectTimeout);
    let fetchRes;
    try {
      fetchRes = await fetch(request.url, {
        method: request.method,
        body: request.data || null,
        redirect: "manual",
        signal,
        ...requesterOptions,
        headers: {
          ...requesterOptions.headers,
          ...request.headers
        }
      });
    } catch (error) {
      return {
        status: 0,
        content: getErrorMessage(error, "Connection timeout"),
        isTimedOut: isAbortError(error)
      };
    }
    clearTimeout(connectTimeout);
    createTimeout(request.responseTimeout);
    try {
      const content = await fetchRes.text();
      return {
        content,
        headers: toResponseHeaders(fetchRes.headers),
        isTimedOut: false,
        status: fetchRes.status
      };
    } catch (error) {
      return {
        status: 0,
        content: getErrorMessage(error, "Socket timeout"),
        isTimedOut: isAbortError(error)
      };
    }
  }
  async function sendStream(request) {
    const fetchRes = await fetch(request.url, {
      method: request.method,
      body: request.data || null,
      redirect: "manual",
      ...requesterOptions,
      headers: {
        ...requesterOptions.headers,
        ...request.headers
      }
    });
    if (!fetchRes.ok) {
      const text = await fetchRes.text();
      throw new import_client_common.StreamRequestError(fetchRes.status, text, toResponseHeaders(fetchRes.headers));
    }
    if (fetchRes.body === null) {
      throw new Error("Response body is null \u2014 streaming not supported");
    }
    return fetchRes.body;
  }
  return { send, sendStream };
}
// Annotate the CommonJS export names for ESM import in node:
0 && (module.exports = {
  createFetchRequester
});
//# sourceMappingURL=requester.fetch.node.cjs.map