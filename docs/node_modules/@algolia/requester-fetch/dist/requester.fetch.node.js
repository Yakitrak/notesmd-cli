// src/createFetchRequester.ts
import { StreamRequestError } from "@algolia/client-common";
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
      throw new StreamRequestError(fetchRes.status, text, toResponseHeaders(fetchRes.headers));
    }
    if (fetchRes.body === null) {
      throw new Error("Response body is null \u2014 streaming not supported");
    }
    return fetchRes.body;
  }
  return { send, sendStream };
}
export {
  createFetchRequester
};
//# sourceMappingURL=requester.fetch.node.js.map