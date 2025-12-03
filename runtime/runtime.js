if (typeof runtime === 'undefined') {
    var runtime = {}
}

Object.assign(runtime, {
    copy(dst, src) {
        if (!dst || !src || !dst.length || !src.length) {
            return 0;
        }
        let n = 0;
        const length = Math.min(dst.length, src.length);
        for (let i = 0; i < length; i++) {
            dst[i] = src[i];
            n++;
        }
        return n;
    },
    // Keep makeSlice/makeMap if they are needed by phase1, or remove if they conflict.
    // The previous runtime.js had them. I'll keep them but fix makeSlice zero value if needed.
    // But since I am overwriting, I should include copy.
    makeSlice(len, cap, zero) {
        len = len || 0;
        return new Array(len).fill(zero);
    },
    makeMap() {
        return {};
    }
});
