if (typeof runtime === 'undefined') {
    runtime = {}
}

runtime.sort = {
    Ints: (a) => a.sort((x, y) => x - y),
    Strings: (a) => a.sort(),
    Float64s: (a) => a.sort((x, y) => x - y),
    Slice: (slice, less) => {
        // sort.Slice(slice, func(i, j int) bool)
        // We need to sort 'slice' using 'less' which takes indices.
        // Create an array of indices
        const indices = Array.from({ length: slice.length }, (_, i) => i);
        indices.sort((i, j) => less(i, j) ? -1 : (less(j, i) ? 1 : 0));

        // Reorder slice based on sorted indices
        const sorted = indices.map(i => slice[i]);
        for (let i = 0; i < slice.length; i++) {
            slice[i] = sorted[i];
        }
    },
    SearchInts: (a, x) => {
        // Binary search
        let i = 0, j = a.length;
        while (i < j) {
            const h = Math.floor((i + j) / 2);
            if (a[h] < x) {
                i = h + 1;
            } else {
                j = h;
            }
        }
        return i;
    }
};
