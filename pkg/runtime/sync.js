if (typeof runtime === 'undefined') {
    runtime = {}
}

runtime.sync = {
    Mutex: class {
        constructor() {
            this._queue = [];
            this._locked = false;
        }

        async Lock() {
            if (!this._locked) {
                this._locked = true;
                return;
            }
            return new Promise(resolve => this._queue.push(resolve));
        }

        Unlock() {
            if (this._queue.length > 0) {
                const resolve = this._queue.shift();
                resolve();
            } else {
                this._locked = false;
            }
        }
    },
    WaitGroup: class {
        constructor() {
            this._count = 0;
            this._waiters = [];
        }

        Add(delta) {
            this._count += delta;
            if (this._count < 0) {
                throw new Error("sync: negative WaitGroup counter");
            }
            if (this._count === 0) {
                this._notify();
            }
        }

        Done() {
            this.Add(-1);
        }

        async Wait() {
            if (this._count === 0) {
                return;
            }
            return new Promise(resolve => this._waiters.push(resolve));
        }

        _notify() {
            const waiters = this._waiters;
            this._waiters = [];
            for (const resolve of waiters) {
                resolve();
            }
        }
    },
    Once: class {
        constructor() {
            this._done = false;
            this._promise = null;
        }

        async Do(fn) {
            if (this._done) {
                return;
            }
            if (!this._promise) {
                this._promise = (async () => {
                    try {
                        await fn();
                    } finally {
                        this._done = true;
                    }
                })();
            }
            await this._promise;
        }
    }
};

// Expose sync globally for type references like 'new sync.Mutex()' if needed
// though the transpiler might need adjustment to use runtime.sync
if (typeof sync === 'undefined') {
    sync = runtime.sync;
}
