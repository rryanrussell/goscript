if (typeof runtime === 'undefined') {
    runtime = {}
}

Object.assign(runtime, {
    makeSlice(len, cap, zero) {
        len = len || 0;
        // In JS, arrays are dynamic, so we just return an array of length 'len'
        // filled with 'zero' value (or undefined if not provided)
        return new Array(len).fill(zero);
    },
    makeMap() {
        // Currently mapped to Object as per existing transpiler logic
        // If we switch to Map, we need to update IndexExpr handling too.
        return {};
    },
    copy(dst, src) {
        let n = 0;
        if (dst && src && dst.length && src.length) {
            n = Math.min(dst.length, src.length);
            for (let i = 0; i < n; i++) {
                dst[i] = src[i];
            }
        }
        return n;
    },
    // Control Flow
    __gs_panicValue: null,
    panic(v) {
        runtime.__gs_panicValue = v;
        throw new Error(v);
    },
    recover() {
        const v = runtime.__gs_panicValue;
        runtime.__gs_panicValue = null;
        return v;
    },
    runDefers(defers) {
        for (let i = defers.length - 1; i >= 0; i--) {
            try {
                defers[i]();
            } catch (e) {
                // If panic happens in defer, it replaces the previous panic
                // But in Go, it aborts the defer and continues panicking?
                // For now, let's just let it bubble up, but we need to handle it if we want proper semantics.
                // Simple version: just call it.
                console.error("panic in defer:", e);
            }
        }
    }
});
