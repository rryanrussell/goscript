if (typeof runtime === 'undefined') {
    runtime = {}
}

runtime.time = {
    Now: () => new Date(),
    Since: (t) => (Date.now() - t.getTime()) * 1000000, // Duration in ns
    Sleep: (d) => new Promise(resolve => setTimeout(resolve, d / 1000000)),
    Second: 1000000000,
    Millisecond: 1000000,
    Microsecond: 1000,
    Nanosecond: 1,
    Unix: (sec, nsec) => new Date(sec * 1000 + nsec / 1000000),
    Date: (year, month, day, hour, min, sec, nsec, loc) => new Date(year, month - 1, day, hour, min, sec, nsec / 1000000),
};
