const FPID = (() => {
    'use strict';

    /**
     * Calculates the SHA-256 hash of a given message.
     * @param {string} message - The string to hash.
     * @returns {Promise<string>} The SHA-256 hash as a hex string.
     */
    const sha256 = async (message) => {
        try {
            const msgUint8 = new TextEncoder().encode(message);
            const hashBuffer = await crypto.subtle.digest('SHA-256', msgUint8);
            const hashArray = Array.from(new Uint8Array(hashBuffer));
            return hashArray.map(b => b.toString(16).padStart(2, '0')).join('');
        } catch (error) {
            return "hash_error";
        }
    };

    /**
     * Runs a function multiple times to check for consistent results, a common anti-fingerprinting countermeasure.
     * @param {Function} fn - The fingerprinting function to test.
     * @param {number} [attempts=3] - The number of times to run the check.
     * @returns {Promise<any|string>} The consistent result, or 'Blocked' if inconsistent.
     */
    const runConsistentCheck = async (fn, attempts = 3) => {
        try {
            const firstResult = await fn();
            // 'Blocked' is a definitive state from a sub-check, no need to re-run.
            if (firstResult === 'Blocked') return 'Blocked';
            for (let i = 1; i < attempts; i++) {
                await new Promise(resolve => setTimeout(resolve, 20));
                if (JSON.stringify(firstResult) !== JSON.stringify(await fn())) return "Blocked";
            }
            return firstResult;
        } catch (error) {
            return "Error";
        }
    };

    /**
     * A wrapper to run a fingerprinting function and format the output.
     * @param {Function} fn - The fingerprinting function.
     * @returns {Promise<{raw: any, hash: string, status: string}>} An object with the raw value, its hash, and the status.
     */
    const runTest = async (fn) => {
        const result = { raw: 'N/A', hash: 'N/A', status: 'Error' };
        try {
            const raw = await runConsistentCheck(fn);
            if (raw === 'Blocked' || raw === 'Error' || raw === 'NotSupported') {
                result.status = raw;
                result.raw = raw;
            } else {
                result.raw = raw;
                result.hash = await sha256(typeof raw === 'string' ? raw : JSON.stringify(raw));
                result.status = 'Success';
            }
        } catch (e) {
            result.raw = e.message;
        }
        return result;
    };


    // --- Private Fingerprinting Modules ---

    const getScreenFingerprint = () => JSON.stringify({ width: window.screen.width, height: window.screen.height, devicePixelRatio: window.devicePixelRatio, colorDepth: screen.colorDepth, pixelDepth: screen.pixelDepth, availWidth: screen.availWidth, availHeight: screen.availHeight, orientation: screen.orientation ? screen.orientation.type : 'N/A' });
    const getSensorFingerprint = () => { try { return JSON.stringify({ gyroscope: 'activated' in new Gyroscope(), accelerometer: 'activated' in new Accelerometer() }); } catch { return "Error"; } };
    const getFontFingerprint = () => { if (!document.fonts?.check) return "NotSupported"; const fonts = ["Andale Mono", "Arial Black", "Courier New", "Malgun Gothic", "Nanum Gothic", "Open Sans", "Noto Sans", "Noto Serif", "Adobe Arabic", "Acumin", "Sloop Script", "Cortado", "Ubuntu Mono", "Big Caslon", "Bodoni 72", "Yu Gothic", "Gulim", "Batang", "BatangChe", "monospace", "sans-serif"]; return fonts.map(font => document.fonts.check(`10px ${font}`)).join(','); };
    const getPluginFingerprint = () => { const plugins = navigator.plugins; if (!plugins || plugins.length === 0) return "NoPlugins"; for (const plugin of plugins) { if (plugin.name.toLowerCase().includes('brave')) return "Blocked"; } return Array.from(plugins).map(p => `${p.name}|${p.description}|${p.filename}`).join(';'); };
    const getBrowserApiFingerprint = () => ['MathMLElement', 'PointerEvent', 'mozInnerScreenX', 'u2f', 'WebGL2RenderingContext', 'SubtleCrypto', 'Text', 'Uint8Array', 'ArrayBuffer', 'ActiveXObject', 'Audio', 'AudioBuffer', 'AudioBufferSourceNode', 'Blob', 'Credential', 'Gamepad', 'Geolocation', 'openDatabase', 'open', 'alert', 'prompt', 'MouseEvent', 'RegExp', 'AuthenticatorResponse', 'AuthenticatorAttestationResponse', 'AuthenticatorAssertionResponse', 'applicationCache', 'Promise', 'indexedDB', 'Cache', 'CacheStorage', 'Clipboard'].map(api => typeof window[api]).join(',');
    const getMathFingerprint = () => [Math.PI, Math.E, Math.LN2, Math.LN10, Math.SQRT1_2, Math.SQRT2, Math.sin(10), Math.sinh(10), Math.cos(10), Math.cosh(10)].join(',');
    const getWebGLFingerprint = () => { const c = document.createElement("canvas"), gl = c.getContext("webgl") || c.getContext("experimental-webgl"); if (!gl) return "NotSupported"; try { const d = gl.getExtension('WEBGL_debug_renderer_info'); return JSON.stringify({ vendor: d ? gl.getParameter(d.UNMASKED_VENDOR_WEBGL) : 'N/A', renderer: d ? gl.getParameter(d.UNMASKED_RENDERER_WEBGL) : 'N/A', extensions: gl.getSupportedExtensions()?.join(',') }); } catch (e) { return "Error"; } finally { c.remove(); } };
    const getJsonOrderFingerprint = () => JSON.stringify({ z: 1, y: "test", x: true, a: [1, 2, 3], b: null });
    const getBatteryFingerprint = () => (typeof navigator.getBattery === 'function').toString();
    const getHardwareApisFingerprint = () => JSON.stringify({ hid: 'hid' in navigator, usb: 'usb' in navigator, serial: 'serial' in navigator });
    const getIntlFingerprint = () => { try { const o = Intl.DateTimeFormat().resolvedOptions(); return JSON.stringify({ locale: o.locale, timeZone: o.timeZone, numberingSystem: o.numberingSystem }); } catch (e) { return "Error"; } };
    
    const validateCanvasPixels = () => {
        try {
            const canvas = document.createElement('canvas'); canvas.width = 100; canvas.height = 50;
            const ctx = canvas.getContext('2d'); if (!ctx) return false;
            ctx.fillStyle = 'rgb(255, 255, 255)'; ctx.fillRect(0, 0, canvas.width, canvas.height);
            ctx.fillStyle = 'rgb(0, 0, 0)'; ctx.fillRect(10, 10, 80, 30);
            const pointsToCheck = { white: [[5, 5], [95, 5], [5, 45], [95, 45]], black: [[15, 15], [85, 15], [15, 35], [85, 35]] };
            for (const point of pointsToCheck.white) { const p = ctx.getImageData(point[0], point[1], 1, 1).data; if (p[0] !== 255 || p[1] !== 255 || p[2] !== 255 || p[3] !== 255) return false; }
            for (const point of pointsToCheck.black) { const p = ctx.getImageData(point[0], point[1], 1, 1).data; if (p[0] !== 0 || p[1] !== 0 || p[2] !== 0 || p[3] !== 255) return false; }
            return true;
        } catch (e) { return false; }
    };

    const getCanvasFingerprint = () => {
        if (!validateCanvasPixels()) return "Blocked";
        const c = document.createElement("canvas"); c.width = 300; c.height = 200; const x = c.getContext("2d"); if (!x) return "NotSupported"; try { const g = x.createLinearGradient(0, 0, 0, 150); g.addColorStop(0, "black"); g.addColorStop(1, "gray"); x.fillStyle = g; x.fillRect(0, 0, 300, 200); x.fillStyle = "white"; x.shadowBlur = 10; x.shadowColor = "yellow"; x.font = "16px 'Noto Sans'"; x.fillText("Hello World! 12345 &^%$#?\\/�🍳🍔🍟🍤😫🙄😑😐🤗😀ΑΕγξΕηθτξΞ", 10, 20); x.beginPath(); x.arc(75, 75, 50, 0, Math.PI * 2, true); x.moveTo(110, 75); x.arc(75, 75, 35, 0, Math.PI, false); x.moveTo(65, 65); x.arc(60, 65, 5, 0, Math.PI * 2, true); x.moveTo(95, 65); x.arc(90, 65, 5, 0, Math.PI * 2, true); x.strokeStyle = "rgba(0, 255, 0, 0.7)"; x.stroke(); return c.toDataURL("image/png"); } catch (e) { return "Error"; } finally { c.remove(); } 
    };
    
    const validateAudioContext = (buffer) => {
        if (!buffer) return false;
        let firstValue = buffer[0];
        for (let i = 1; i < buffer.length; i++) {
            if (buffer[i] !== firstValue) return true;
        }
        return false;
    };

    const getAudioContextFingerprint = () => new Promise(r => { 
        try {
            if (!!JSON.stringify(navigator.userAgentData).match(/Brave/)) return r("Blocked");
            const a = new (window.OfflineAudioContext || window.webkitOfflineAudioContext)(1, 44100, 44100); if (!a) return r("NotSupported"); 
            const o = a.createOscillator(); o.type = "triangle"; o.frequency.setValueAtTime(10000, a.currentTime); 
            const c = a.createDynamicsCompressor(); c.threshold.setValueAtTime(-50, a.currentTime); c.knee.setValueAtTime(40, a.currentTime); c.ratio.setValueAtTime(12, a.currentTime); c.attack.setValueAtTime(0, a.currentTime); c.release.setValueAtTime(0.25, a.currentTime); 
            o.connect(c); c.connect(a.destination); o.start(0); a.startRendering(); 
            a.oncomplete = e => {
                const buffer = e.renderedBuffer.getChannelData(0);
                if (!validateAudioContext(buffer)) return r("Blocked");
                r(buffer.slice(4500, 5000).reduce((s, v) => s + Math.abs(v), 0).toString());
            };
        } catch (e) { r("Error"); } 
    });

    // --- Public Method ---

    /**
     * Generates a comprehensive browser fingerprint.
     * @returns {Promise<Object>} A promise that resolves to an object containing all fingerprinting results and a final hash.
     */
    const generate = async () => {
        const results = {
            canvas: await runTest(getCanvasFingerprint),
            audio: await runTest(getAudioContextFingerprint),
            webgl: await runTest(getWebGLFingerprint),
            fonts: await runTest(getFontFingerprint),
            screen: await runTest(getScreenFingerprint),
            intl: await runTest(getIntlFingerprint),
            sensors: await runTest(getSensorFingerprint),
            plugins: await runTest(getPluginFingerprint),
            browserApis: await runTest(getBrowserApiFingerprint),
            hardwareApis: await runTest(getHardwareApisFingerprint),
            battery: await runTest(getBatteryFingerprint),
            math: await runTest(getMathFingerprint),
            jsonOrder: await runTest(getJsonOrderFingerprint),
        };
        const finalHashes = Object.values(results).map(v => v.hash).sort().join('');
        results.finalHash = await sha256(finalHashes);

        return results;
    };

    return {
        generate: generate
    };

})();

//@@START_CONFIG@@
const TELEMETRY_FP_VERSION = 1;
const TELEMETRY_BASEURL = "https://telemetry.ex.gosuda.org";
const CLIENT_VERSION = "20250807-V1ALPHA1";
//@@END_CONFIG@@

async function checkClientStatus() {
    // check if client is registered
    let clientID = localStorage.getItem("telemetry_client_id");
    let clientToken = localStorage.getItem("telemetry_client_token");

    if (!clientID || !clientToken) {
        return false
    }

    const resp = await fetch(TELEMETRY_BASEURL + "/client/status", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            id: clientID,
            token: clientToken,
        }),
    });

    if (resp.status == 200) {
        return true;
    }

    return false;
}

async function registerClient() {
    // register client
    const resp = await fetch(TELEMETRY_BASEURL + "/client/register", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
    });
    if (resp.status !== 201) {
        throw new Error("Failed to register client");
    }

    const clientIdentity = await resp.json();
    localStorage.setItem("telemetry_client_id", clientIdentity.id);
    localStorage.setItem("telemetry_client_token", clientIdentity.token);

    return clientIdentity;
}

async function registerFingerprint(fingerprint) {
    let clientID = localStorage.getItem("telemetry_client_id");
    let clientToken = localStorage.getItem("telemetry_client_token");

    // register fingerprint
    const resp = await fetch(TELEMETRY_BASEURL + "/client/checkin", {
        method: "POST",
        headers: {
            "Content-Type": "application/json",
        },
        body: JSON.stringify({
            client_id: clientID,
            client_token: clientToken,
            version: CLIENT_VERSION,
            fpv: TELEMETRY_FP_VERSION,
            fp: fingerprint,
            ua: navigator.userAgent,
            uad: JSON.stringify(navigator.userAgentData),
        }),
    });
    if (resp.status !== 200) {
        throw new Error("Failed to register fingerprint");
    }
}

async function telemetry() {
    // check if client is registered
    let clientFingerprint = localStorage.getItem("telemetry_client_fingerprint");

    if (!await checkClientStatus()) {
        await registerClient();

        const ok = await checkClientStatus();
        if (!ok) {
            throw new Error("Failed to register client");
        }
    }

    const fp = await FPID.generate();
    console.log("fpid:", fp.finalHash);
    if (fp.finalHash !== clientFingerprint) {
        await registerFingerprint(fp.finalHash);
    }
}

telemetry();
