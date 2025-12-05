if (typeof runtime === 'undefined') {
    runtime = {}
}

runtime.strings = {
    Contains(s, substr) {
        return s.includes(substr);
    },
    HasPrefix(s, prefix) {
        return s.startsWith(prefix);
    },
    HasSuffix(s, suffix) {
        return s.endsWith(suffix);
    },
    Index(s, substr) {
        return s.indexOf(substr);
    },
    LastIndex(s, substr) {
        return s.lastIndexOf(substr);
    },
    Split(s, sep) {
        return s.split(sep);
    },
    Join(elems, sep) {
        return elems.join(sep);
    },
    ToLower(s) {
        return s.toLowerCase();
    },
    ToUpper(s) {
        return s.toUpperCase();
    },
    Replace(s, old, newS, n) {
        if (n === -1) {
            return s.split(old).join(newS);
        }
        if (old === "") {
            // Special case for empty old string: insert newS between every char
            if (n === 0) return s;
        }

        let res = "";
        let start = 0;
        for (let i = 0; i < n; i++) {
            const idx = s.indexOf(old, start);
            if (idx === -1) break;
            res += s.substring(start, idx) + newS;
            start = idx + old.length;
        }
        res += s.substring(start);
        return res;
    },
    ReplaceAll(s, old, newS) {
        return s.split(old).join(newS);
    },
    Trim(s, cutset) {
        // Simple implementation for common case where cutset is space
        // For full Go compatibility, need regex or loop
        if (cutset === " \t\n\r") {
            return s.trim();
        }
        // TODO: Implement full Trim
        let start = 0;
        while (start < s.length && cutset.includes(s[start])) start++;
        let end = s.length;
        while (end > start && cutset.includes(s[end - 1])) end--;
        return s.substring(start, end);
    },
    TrimSpace(s) {
        return s.trim();
    }
};
