if (typeof runtime === 'undefined') {
    runtime = {}
}

runtime.errors = {
    New: (text) => new Error(text),
    Is: (err, target) => {
        if (err === target) return true;
        if (err && err.message === target.message) return true;
        return false;
    },
    As: (err, target) => {
        // Basic implementation: check if err is instance of target's constructor if possible
        // But target is usually a pointer to a type.
        // For now, return false as placeholder.
        return false;
    },
    Unwrap: (err) => {
        return err.cause || null;
    }
};
