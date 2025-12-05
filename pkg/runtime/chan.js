if (typeof runtime === 'undefined') {
    runtime = {}
}

class Channel {
    constructor(capacity = 0) {
        this.capacity = capacity;
        this.buffer = [];
        this.closed = false;
        this.sendq = []; // {resolve, val, select}
        this.recvq = []; // {resolve, select}
    }

    async send(val) {
        if (this.closed) throw new Error("send on closed channel");

        // Check for waiting receiver
        while (this.recvq.length > 0) {
            const waiter = this.recvq.shift();
            if (waiter.select && !waiter.select.tryLock()) continue;
            waiter.resolve(val);
            return;
        }

        // Check buffer
        if (this.buffer.length < this.capacity) {
            this.buffer.push(val);
            return;
        }

        // Block
        return new Promise((resolve) => {
            this.sendq.push({ resolve, val });
        });
    }

    async recv() {
        // Check buffer
        if (this.buffer.length > 0) {
            const val = this.buffer.shift();
            // If sender waiting, move to buffer
            while (this.sendq.length > 0) {
                const sender = this.sendq.shift();
                if (sender.select && !sender.select.tryLock()) continue;
                this.buffer.push(sender.val);
                sender.resolve();
                break;
            }
            return val;
        }

        if (this.closed) return null;

        // Check waiting sender (unbuffered)
        while (this.sendq.length > 0) {
            const sender = this.sendq.shift();
            if (sender.select && !sender.select.tryLock()) continue;
            const val = sender.val;
            sender.resolve();
            return val;
        }

        // Block
        return new Promise((resolve) => {
            this.recvq.push({ resolve });
        });
    }

    close() {
        if (this.closed) return;
        this.closed = true;
        // Release all receivers with null
        while (this.recvq.length > 0) {
            const waiter = this.recvq.shift();
            if (waiter.select && !waiter.select.tryLock()) continue;
            waiter.resolve(null);
        }
        // Reject all senders
        while (this.sendq.length > 0) {
            const sender = this.sendq.shift();
            if (sender.select && !sender.select.tryLock()) continue;
            // In Go, send on closed panics. Here we throw/reject.
        }
    }
}

class SelectGroup {
    constructor() {
        this.active = true;
    }
    tryLock() {
        if (!this.active) return false;
        this.active = false;
        return true;
    }
}

Object.assign(runtime, {
    Channel,

    makeChan(cap) {
        return new Channel(cap || 0);
    },

    go(fn) {
        // Fire and forget
        fn();
    },

    async select(cases) {
        // cases: [{op: 'recv', chan: c, handler: v => ...}, {op: 'send', chan: c, val: v, handler: () => ...}, {op: 'default', handler: () => ...}]

        const group = new SelectGroup();

        return new Promise((resolve) => {
            let hasDefault = false;
            let defaultHandler = null;

            // Register with all channels
            for (let c of cases) {
                if (c.op === 'default') {
                    hasDefault = true;
                    defaultHandler = c.handler;
                    continue;
                }

                if (c.op === 'recv') {
                    const ch = c.chan;

                    // Try immediate recv
                    if (ch.buffer.length > 0 || (!ch.closed && ch.sendq.length > 0)) {
                        if (group.tryLock()) {
                            ch.recv().then(v => resolve(c.handler(v)));
                            return;
                        }
                    }

                    if (ch.closed) {
                        if (group.tryLock()) {
                            resolve(c.handler(null));
                            return;
                        }
                    }

                    // Add to wait queue
                    ch.recvq.push({
                        resolve: (v) => resolve(c.handler(v)),
                        select: group
                    });
                } else if (c.op === 'send') {
                    const ch = c.chan;
                    if (ch.closed) throw new Error("send on closed channel");

                    // Try immediate send
                    if (ch.recvq.length > 0 || ch.buffer.length < ch.capacity) {
                        if (group.tryLock()) {
                            ch.send(c.val).then(() => resolve(c.handler()));
                            return;
                        }
                    }

                    ch.sendq.push({
                        resolve: () => resolve(c.handler()),
                        val: c.val,
                        select: group
                    });
                }
            }

            if (hasDefault && group.tryLock()) {
                resolve(defaultHandler());
            }
        });
    }
});
