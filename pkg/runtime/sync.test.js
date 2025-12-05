const assert = require('assert');
require('./runtime.js');
require('./sync.js');

async function testMutex() {
    const mu = new runtime.sync.Mutex();
    let count = 0;
    const p1 = (async () => {
        await mu.Lock();
        const c = count;
        await new Promise(r => setTimeout(r, 10));
        count = c + 1;
        mu.Unlock();
    })();
    const p2 = (async () => {
        await mu.Lock();
        const c = count;
        await new Promise(r => setTimeout(r, 10));
        count = c + 1;
        mu.Unlock();
    })();
    await Promise.all([p1, p2]);
    assert.strictEqual(count, 2, "Mutex should protect critical section");
    console.log("Mutex test passed");
}

async function testWaitGroup() {
    const wg = new runtime.sync.WaitGroup();
    let count = 0;
    wg.Add(2);
    setTimeout(() => {
        count++;
        wg.Done();
    }, 10);
    setTimeout(() => {
        count++;
        wg.Done();
    }, 20);
    await wg.Wait();
    assert.strictEqual(count, 2, "WaitGroup should wait for all tasks");
    console.log("WaitGroup test passed");
}

async function testOnce() {
    const once = new runtime.sync.Once();
    let count = 0;
    const fn = async () => {
        await new Promise(r => setTimeout(r, 10));
        count++;
    };
    await Promise.all([once.Do(fn), once.Do(fn), once.Do(fn)]);
    assert.strictEqual(count, 1, "Once should execute only once");
    console.log("Once test passed");
}

async function testTranspiledStructure() {
    // Mock runtime environment
    Object.assign(runtime, {
        go: (fn) => fn(),
        time: {
            Sleep: (d) => new Promise(resolve => setTimeout(resolve, d / 1000000)),
            Millisecond: 1000000
        },
        runDefers: (defers) => {
            for (let i = defers.length - 1; i >= 0; i--) {
                defers[i]();
            }
        }
    });
    // Use the global runtime.sync classes

    let mu = new runtime.sync.Mutex();
    let count = 0;
    let wg = new runtime.sync.WaitGroup();
    wg.Add(2);

    runtime.go(() => {
        (async function () {
            let __gs_deferred = [];
            try {
                await mu.Lock();
                __gs_deferred.push(() => { mu.Unlock(); });

                let c = count;
                await runtime.time.Sleep(10 * runtime.time.Millisecond);
                count = c + 1;
                wg.Done();
            } finally {
                runtime.runDefers(__gs_deferred);
            }
        })();
    });

    runtime.go(() => {
        (async function () {
            let __gs_deferred = [];
            try {
                await mu.Lock();
                __gs_deferred.push(() => { mu.Unlock(); });

                let c = count;
                await runtime.time.Sleep(10 * runtime.time.Millisecond);
                count = c + 1;
                wg.Done();
            } finally {
                runtime.runDefers(__gs_deferred);
            }
        })();
    });

    await wg.Wait();
    assert.strictEqual(count, 2, "Transpiled structure should work");
    console.log("Transpiled structure test passed");
}

(async () => {
    try {
        await testMutex();
        await testWaitGroup();
        await testOnce();
        await testTranspiledStructure();
        console.log("All sync tests passed");
    } catch (e) {
        console.error(e);
        process.exit(1);
    }
})();
