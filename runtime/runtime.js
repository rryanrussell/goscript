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
    }
});
