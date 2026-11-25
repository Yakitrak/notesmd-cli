/*
  @license
	Rollup.js v4.68.8
	Sat, 21 Nov 2026 06:53:43 GMT - commit a845aca377e70a43bab0575931f2d403de1c35fa

	https://github.com/rollup/rollup

	Released under the MIT License.
*/
'use strict';

Object.defineProperty(exports, Symbol.toStringTag, { value: 'Module' });

const rollup = require('./rollup.js');
const path = require('node:path');
const process = require('node:process');
const index = require('./index.js');
const node_os = require('node:os');
const rollup_js = require('./node-entry.js');
require('./parseAst.js');
require('../native.js');
require('path');
require('node:perf_hooks');
require('node:fs/promises');
require('fs');
require('util');
require('stream');
require('os');
require('./fsevents-importer.js');
require('events');

class FileWatcher {
    constructor(task, chokidarOptions) {
        this.transformWatchers = new Map();
        this.chokidarOptions = chokidarOptions;
        this.task = task;
        this.watcher = this.createWatcher(null);
    }
    close() {
        this.watcher.close();
        for (const watcher of this.transformWatchers.values()) {
            watcher.close();
        }
    }
    unwatch(id) {
        this.watcher.unwatch(id);
        const transformWatcher = this.transformWatchers.get(id);
        if (transformWatcher) {
            this.transformWatchers.delete(id);
            transformWatcher.close();
        }
    }
    watch(id, isTransformDependency) {
        if (isTransformDependency) {
            const watcher = this.transformWatchers.get(id) ?? this.createWatcher(id);
            watcher.add(id);
            this.transformWatchers.set(id, watcher);
        }
        else {
            this.watcher.add(id);
        }
    }
    createWatcher(transformWatcherId) {
        const task = this.task;
        const isLinux = node_os.platform() === 'linux';
        const isFreeBSD = node_os.platform() === 'freebsd';
        const isTransformDependency = transformWatcherId !== null;
        const handleChange = (id, event) => {
            const changedId = transformWatcherId || id;
            if (isLinux || isFreeBSD) {
                // unwatching and watching fixes an issue with chokidar where on certain systems,
                // a file that was unlinked and immediately recreated would create a change event
                // but then no longer any further events
                watcher.unwatch(changedId);
                watcher.add(changedId);
            }
            task.invalidate(changedId, { event, isTransformDependency });
        };
        const watcher = index.chokidar
            .watch([], this.chokidarOptions)
            .on('add', id => handleChange(id, 'create'))
            .on('change', id => handleChange(id, 'update'))
            .on('unlink', id => handleChange(id, 'delete'));
        return watcher;
    }
}

const eventsRewrites = {
    create: {
        create: 'buggy',
        delete: null, //delete file from map
        update: 'create'
    },
    delete: {
        create: 'update',
        delete: 'buggy',
        update: 'buggy'
    },
    update: {
        create: 'buggy',
        delete: 'delete',
        update: 'update'
    }
};
class Watcher {
    constructor(optionsList, emitter) {
        this.closed = false;
        this.buildDelay = 0;
        this.buildTimeout = null;
        this.invalidatedIds = new Map();
        this.rerun = false;
        // Starts held so that invalidations before the initial run only request a rerun.
        this.running = true;
        this.emitter = emitter;
        emitter.close = this.close.bind(this);
        this.tasks = optionsList.map(options => new Task(this, options));
        for (const { watch } of optionsList) {
            if (watch && typeof watch.buildDelay === 'number') {
                this.buildDelay = Math.max(this.buildDelay, watch.buildDelay);
            }
        }
        process.nextTick(() => this.runExclusively(() => this.run()));
    }
    async close() {
        if (this.closed)
            return;
        this.closed = true;
        if (this.buildTimeout)
            clearTimeout(this.buildTimeout);
        try {
            await rollup_js.waitForAllAndRethrowFirstFailure([
                ...this.tasks.map(task => task.close()),
                this.emitter.emit('close')
            ]);
        }
        finally {
            this.emitter.removeAllListeners();
        }
    }
    invalidate(file) {
        if (this.closed)
            return;
        if (file) {
            const previousEvent = this.invalidatedIds.get(file.id);
            const event = previousEvent ? eventsRewrites[previousEvent][file.event] : file.event;
            if (event === 'buggy') {
                //TODO: throws or warn? Currently just ignore, uses new event
                this.invalidatedIds.set(file.id, file.event);
            }
            else if (event === null) {
                this.invalidatedIds.delete(file.id);
            }
            else {
                this.invalidatedIds.set(file.id, event);
            }
        }
        if (this.running) {
            this.rerun = true;
            return;
        }
        if (this.buildTimeout)
            clearTimeout(this.buildTimeout);
        this.buildTimeout = setTimeout(() => {
            this.buildTimeout = null;
            this.runExclusively(() => this.announceChangesAndRebuild());
        }, this.buildDelay);
    }
    // The running flag spans the whole cycle, including async listeners, so
    // that invalidations meanwhile only request a rerun.
    async runExclusively(cycle) {
        this.running = true;
        try {
            await cycle();
        }
        catch (error) {
            await this.reportError(error);
        }
        finally {
            this.running = false;
        }
        if (this.rerun) {
            this.rerun = false;
            this.invalidate();
        }
    }
    async announceChangesAndRebuild() {
        await this.announcePendingChanges();
        if (this.closed || !this.hasInvalidatedTask())
            return;
        await this.emitter.emit('restart');
        await this.announcePendingChanges();
        if (this.closed)
            return;
        this.emitter.removeListenersForCurrentRun();
        await this.run();
    }
    // The upcoming run consumes announced changes via the task flags, making
    // their rerun requests obsolete. The reset must directly follow the last
    // emptiness check so that no change slips through unannounced.
    async announcePendingChanges() {
        while (this.invalidatedIds.size > 0 && !this.closed) {
            await this.announceChangeBatch();
        }
        this.rerun = false;
    }
    // Cleared before emitting so that changes arriving during slow hooks,
    // including follow-ups of the same file, form the next batch instead of
    // being dropped with this one.
    async announceChangeBatch() {
        const changes = [...this.invalidatedIds];
        this.invalidatedIds.clear();
        await rollup_js.waitForAllAndRethrowFirstFailure(changes.flatMap(([id, event]) => [
            ...this.tasks.map(task => task.announceChange(id, event)),
            this.emitter.emit('change', id, { event })
        ]));
    }
    hasInvalidatedTask() {
        return this.tasks.some(task => task.isInvalidated());
    }
    async reportError(error) {
        await this.emitter.emit('event', {
            code: 'ERROR',
            error,
            result: null
        });
        await this.emitter.emit('event', {
            code: 'END'
        });
    }
    async run() {
        await this.emitter.emit('event', {
            code: 'START'
        });
        for (const task of this.tasks) {
            await task.run();
        }
        if (this.closed)
            return;
        await this.emitter.emit('event', {
            code: 'END'
        });
    }
}
class Task {
    constructor(watcher, options) {
        this.cache = { modules: [] };
        this.watchFiles = [];
        this.invalidated = true;
        this.watched = new Set();
        this.watchHooks = null;
        this.watcher = watcher;
        this.options = options;
        this.skipWrite = Boolean(options.watch && options.watch.skipWrite);
        this.outputs = this.options.output;
        this.outputFiles = this.outputs.map(output => {
            if (output.file || output.dir)
                return path.resolve(output.file || output.dir);
            return undefined;
        });
        this.watchOptions = this.options.watch || {};
        this.filter = rollup.createFilter(this.watchOptions.include, this.watchOptions.exclude);
        this.fileWatcher = new FileWatcher(this, {
            ...this.watchOptions.chokidar,
            disableGlobbing: true,
            ignoreInitial: true
        });
    }
    async announceChange(id, event) {
        await this.watchHooks?.watchChange(id, { event });
    }
    async close() {
        this.fileWatcher.close();
        const watchHooks = this.watchHooks;
        this.watchHooks = null;
        await watchHooks?.closeWatcher();
    }
    invalidate(id, details) {
        this.invalidated = true;
        if (details.isTransformDependency) {
            for (const module of this.cache.modules) {
                if (!module.transformDependencies.includes(id))
                    continue;
                // effective invalidation
                module.originalCode = null;
            }
        }
        this.watcher.invalidate({ event: details.event, id });
        this.watchOptions.onInvalidate?.(id);
    }
    isInvalidated() {
        return this.invalidated;
    }
    async run() {
        if (!this.invalidated || this.watcher.closed)
            return;
        this.invalidated = false;
        const options = {
            ...this.options,
            cache: this.cache
        };
        const start = Date.now();
        await this.watcher.emitter.emit('event', {
            code: 'BUNDLE_START',
            input: this.options.input,
            output: this.outputFiles
        });
        if (this.watcher.closed) {
            return;
        }
        let result = null;
        try {
            result = await rollup.rollupInternal(options, watchHooks => {
                if (!this.watcher.closed)
                    this.watchHooks = watchHooks;
            });
            if (this.watcher.closed) {
                return;
            }
            this.updateWatchedFiles(result);
            if (!this.skipWrite) {
                await Promise.all(this.outputs.map(output => result.write(output)));
                if (this.watcher.closed) {
                    return;
                }
                this.updateWatchedFiles(result);
            }
            await this.watcher.emitter.emit('event', {
                code: 'BUNDLE_END',
                duration: Date.now() - start,
                input: this.options.input,
                output: this.outputFiles,
                result
            });
        }
        catch (error) {
            if (!this.watcher.closed) {
                if (Array.isArray(error.watchFiles)) {
                    for (const id of error.watchFiles) {
                        this.watchFile(id);
                    }
                }
                if (error.id) {
                    this.cache.modules = this.cache.modules.filter(module => module.id !== error.id);
                }
            }
            await this.watcher.emitter.emit('event', {
                code: 'ERROR',
                error,
                result
            });
        }
    }
    updateWatchedFiles(result) {
        const previouslyWatched = this.watched;
        this.watched = new Set();
        this.watchFiles = result.watchFiles;
        this.cache = result.cache;
        for (const id of this.watchFiles) {
            this.watchFile(id);
        }
        for (const module of this.cache.modules) {
            for (const dependencyId of module.transformDependencies) {
                this.watchFile(dependencyId, true);
            }
        }
        for (const id of previouslyWatched) {
            if (!this.watched.has(id)) {
                this.fileWatcher.unwatch(id);
            }
        }
    }
    watchFile(id, isTransformDependency = false) {
        if (!this.filter(id))
            return;
        this.watched.add(id);
        if (this.outputFiles.includes(id)) {
            throw new Error('Cannot import the generated bundle');
        }
        // this is necessary to ensure that any 'renamed' files
        // continue to be watched following an error
        this.fileWatcher.watch(id, isTransformDependency);
    }
}

exports.Task = Task;
exports.Watcher = Watcher;
//# sourceMappingURL=watch.js.map
