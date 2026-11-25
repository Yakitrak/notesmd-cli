/*
  @license
	Rollup.js v4.68.8
	Sat, 21 Nov 2026 06:53:43 GMT - commit a845aca377e70a43bab0575931f2d403de1c35fa

	https://github.com/rollup/rollup

	Released under the MIT License.
*/
'use strict';

const rollup = require('./rollup.js');
const parseAst_js = require('./parseAst.js');
const fseventsImporter = require('./fsevents-importer.js');

function waitForAllAndRethrowFirstFailure(promises) {
    return Promise.allSettled(promises).then(results => {
        for (const result of results) {
            if (result.status === 'rejected') {
                throw result.reason;
            }
        }
    });
}
class WatchEmitter {
    constructor() {
        this.currentHandlers = Object.create(null);
        this.persistentHandlers = Object.create(null);
    }
    // Will be overwritten by Rollup
    async close() { }
    emit(event, ...parameters) {
        // The async wrapper turns a synchronous throw into a rejection so that no sibling listener is skipped.
        return waitForAllAndRethrowFirstFailure([...this.getCurrentHandlers(event), ...this.getPersistentHandlers(event)].map(async (handler) => handler(...parameters)));
    }
    off(event, listener) {
        const listeners = this.persistentHandlers[event];
        if (listeners) {
            // A hack stolen from "mitt": ">>> 0" does not change numbers >= 0, but -1
            // (which would remove the last array element if used unchanged) is turned
            // into max_int, which is outside the array and does not change anything.
            listeners.splice(listeners.indexOf(listener) >>> 0, 1);
        }
        return this;
    }
    on(event, listener) {
        this.getPersistentHandlers(event).push(listener);
        return this;
    }
    onCurrentRun(event, listener) {
        this.getCurrentHandlers(event).push(listener);
        return this;
    }
    once(event, listener) {
        const selfRemovingListener = (...parameters) => {
            this.off(event, selfRemovingListener);
            return listener(...parameters);
        };
        this.on(event, selfRemovingListener);
        return this;
    }
    removeAllListeners() {
        this.removeListenersForCurrentRun();
        this.persistentHandlers = Object.create(null);
        return this;
    }
    removeListenersForCurrentRun() {
        this.currentHandlers = Object.create(null);
        return this;
    }
    getCurrentHandlers(event) {
        return this.currentHandlers[event] || (this.currentHandlers[event] = []);
    }
    getPersistentHandlers(event) {
        return this.persistentHandlers[event] || (this.persistentHandlers[event] = []);
    }
}

function watch(configs) {
    const emitter = new WatchEmitter();
    let isClosed = false;
    // Until the watcher takes over closing, a close prevents its construction
    // and emits close itself, as no watcher will.
    emitter.close = async () => {
        if (isClosed)
            return;
        isClosed = true;
        try {
            await emitter.emit('close');
        }
        finally {
            emitter.removeAllListeners();
        }
    };
    watchInternal(configs, emitter, () => isClosed).catch(error => {
        rollup.handleError(error);
    });
    return emitter;
}
function ensureTrailingSlash(path) {
    if (path[path.length - 1] !== '/') {
        return `${path}/`;
    }
    return path;
}
function checkWatchConfig(config) {
    for (const item of config) {
        if (typeof item.watch !== 'boolean' && item.watch?.allowInputInsideOutputPath) {
            break;
        }
        if (item.input && item.output) {
            const input = typeof item.input === 'string' ? rollup.ensureArray(item.input) : item.input;
            const outputs = rollup.ensureArray(item.output);
            for (const index in input) {
                const inputPath = input[index];
                if (typeof inputPath !== 'string') {
                    continue;
                }
                const outputWithInputAsSubPath = outputs.find(({ dir }) => dir && ensureTrailingSlash(inputPath).startsWith(ensureTrailingSlash(dir)));
                if (outputWithInputAsSubPath) {
                    parseAst_js.error(parseAst_js.logInvalidOption('watch', parseAst_js.URL_WATCH, `the input "${inputPath}" is a subpath of the output "${outputWithInputAsSubPath.dir}"`));
                }
            }
        }
    }
}
async function watchInternal(configs, emitter, wasClosed) {
    const optionsList = await Promise.all(rollup.ensureArray(configs).map(config => rollup.mergeOptions(config, true)));
    const watchOptionsList = optionsList.filter(config => config.watch !== false);
    if (watchOptionsList.length === 0) {
        return parseAst_js.error(parseAst_js.logInvalidOption('watch', parseAst_js.URL_WATCH, 'there must be at least one config where "watch" is not set to "false"'));
    }
    checkWatchConfig(watchOptionsList);
    await fseventsImporter.loadFsEvents();
    const { Watcher } = await Promise.resolve().then(() => require('./watch.js'));
    if (wasClosed()) {
        return;
    }
    new Watcher(watchOptionsList, emitter);
}

const VERSION = rollup.package_.version;

exports.VERSION = VERSION;
exports.waitForAllAndRethrowFirstFailure = waitForAllAndRethrowFirstFailure;
exports.watch = watch;
//# sourceMappingURL=node-entry.js.map
