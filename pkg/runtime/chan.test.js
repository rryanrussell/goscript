import { test, expect } from "bun:test";
import { readFileSync } from "fs";

// Setup global runtime
global.runtime = {};

// Load chan.js
// Note: chan.js is a raw script intended for concatenation, not a module.
// We use eval to load it into the global scope for testing.
const chanCode = readFileSync("chan.js", "utf8");
eval(chanCode);

test("Channel send/recv", async () => {
    const ch = runtime.makeChan();

    runtime.go(async () => {
        await ch.send(42);
    });

    const val = await ch.recv();
    expect(val).toBe(42);
});

test("Buffered Channel", async () => {
    const ch = runtime.makeChan(2);
    await ch.send(1);
    await ch.send(2);

    expect(await ch.recv()).toBe(1);
    expect(await ch.recv()).toBe(2);
});

test("Select Recv", async () => {
    const c1 = runtime.makeChan();
    const c2 = runtime.makeChan();

    runtime.go(async () => {
        await c1.send(10);
    });

    const res = await runtime.select([
        { op: 'recv', chan: c1, handler: (v) => v },
        { op: 'recv', chan: c2, handler: (v) => v }
    ]);

    expect(res).toBe(10);
});

test("Select Send", async () => {
    const c1 = runtime.makeChan(1);

    const res = await runtime.select([
        { op: 'send', chan: c1, val: 20, handler: () => 'sent' },
        { op: 'default', handler: () => 'default' }
    ]);

    expect(res).toBe('sent');
    expect(await c1.recv()).toBe(20);
});

test("Select Default", async () => {
    const c1 = runtime.makeChan();

    const res = await runtime.select([
        { op: 'recv', chan: c1, handler: (v) => v },
        { op: 'default', handler: () => 'default' }
    ]);

    expect(res).toBe('default');
});
