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
    }
});
