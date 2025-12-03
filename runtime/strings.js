if (typeof strings === 'undefined') {
    strings = {}
}

Object.assign(strings, {
    Contains(s, substr) {
        return s.includes(substr);
    },
    HasPrefix(s, prefix) {
        return s.startsWith(prefix);
    },
    HasSuffix(s, suffix) {
        return s.endsWith(suffix);
    },
    Split(s, sep) {
        if (sep === "") {
            // Go strings.Split("", "") -> []
            if (s === "") {
                return [];
            }
            // Go strings.Split("abc", "") -> ["", "a", "b", "c", ""]
            // Use [...s] for unicode correctness
            return ["", ...s, ""];
        }
        return s.split(sep);
    },
    Join(arr, sep) {
        return arr.join(sep);
    },
    ToUpper(s) {
        return s.toUpperCase();
    },
    ToLower(s) {
        return s.toLowerCase();
    },
    Replace(s, old, new_val, n) {
        if (old === "") {
             if (n === 0) return s;

             let result = "";
             let count = 0;
             // Start
             if (n < 0 || count < n) {
                 result += new_val;
                 count++;
             }

             // Iterate over code points (runes) for correctness with unicode?
             // Go Replace operates on runes if old is empty?
             // "If old is empty, it matches at the beginning of the string and after each UTF-8 sequence"
             // JS string iteration is by code point or UTF-16 unit?
             // for ... of iterates code points.
             for (const char of s) {
                 result += char;
                 if (n < 0 || count < n) {
                     result += new_val;
                     count++;
                 }
             }
             return result;
        }

        if (n < 0) {
            return s.split(old).join(new_val);
        }

        let result = "";
        let i = 0;
        let count = 0;
        while (count < n) {
            const idx = s.indexOf(old, i);
            if (idx === -1) {
                break;
            }
            result += s.substring(i, idx) + new_val;
            i = idx + old.length;
            count++;
        }
        result += s.substring(i);
        return result;
    },
    ReplaceAll(s, old, new_val) {
        if (old === "") {
            return strings.Replace(s, old, new_val, -1);
        }
        return s.split(old).join(new_val);
    },
    TrimSpace(s) {
        return s.trim();
    },
    Trim(s, cutset) {
         if (cutset === "") return s;
         const escapedCutset = cutset.replace(/[.*+?^${}()|[\]\\-]/g, '\\$&');
         const pattern = new RegExp(`^[${escapedCutset}]+|[${escapedCutset}]+$`, 'g');
         return s.replace(pattern, '');
    }
});
