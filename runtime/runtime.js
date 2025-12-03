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
    }
});
