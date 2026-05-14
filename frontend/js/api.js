const BASE = '';

/**
 * Retrieves the Authorization header for authenticated API requests.
 * Reads the Bearer token stored in sessionStorage.
 * @returns {{ Authorization: string } | {}}
 */
function getAuthHeader() {
    const apiKey = sessionStorage.getItem('apiKey');
    return apiKey ? { 'Authorization': `Bearer ${apiKey}` } : {};
}

/**
 * Parses a fetch Response, throwing an enriched Error with a status property on failure.
 * @param {Response} res - The fetch Response object.
 * @returns {Promise<any>} Resolves with the parsed JSON body.
 * @throws {Error} When the HTTP response status is not OK (error.status is set).
 */
async function parseResponse(res) {
    if (!res.ok) {
        const err = new Error(`API error: ${res.status}`);
        err.status = res.status;
        throw err;
    }
    return res.json();
}

/**
 * Stores the API key in sessionStorage (cleared automatically when the tab closes).
 * @param {string} key - The Bearer token to store.
 */
export function setApiKey(key) {
    sessionStorage.setItem('apiKey', key);
}

/**
 * Removes the API key from sessionStorage, effectively logging out the current session.
 */
export function clearApiKey() {
    sessionStorage.removeItem('apiKey');
}

/**
 * Returns true when an API key is currently stored in sessionStorage.
 * @returns {boolean}
 */
export function hasApiKey() {
    return !!sessionStorage.getItem('apiKey');
}

/**
 * Fetches the list of notes, optionally filtered by folder.
 * @param {string} [folder=''] - The folder slug to filter by (empty string = all folders).
 * @returns {Promise<Note[]>}
 */
export async function fetchNotes(folder = '') {
    const params = folder ? `?folder=${encodeURIComponent(folder)}` : '';
    const res = await fetch(`${BASE}/api/notes${params}`, {
        headers: getAuthHeader(),
    });
    return parseResponse(res);
}

/**
 * Fetches a single note by its database ID.
 * @param {number} id - The note ID.
 * @returns {Promise<Note>}
 */
export async function fetchNote(id) {
    const res = await fetch(`${BASE}/api/notes/${id}`, {
        headers: getAuthHeader(),
    });
    return parseResponse(res);
}

/**
 * Creates a new note.
 * @param {{ title: string, content: string, folder: string }} data - Note payload.
 * @returns {Promise<{ id: number }>}
 */
export async function createNote(data) {
    const res = await fetch(`${BASE}/api/notes`, {
        method: 'POST',
        headers: { 
            'Content-Type': 'application/json',
            ...getAuthHeader()
        },
        body: JSON.stringify(data),
    });
    return parseResponse(res);
}

/**
 * Updates an existing note.
 * @param {number} id - The ID of the note to update.
 * @param {{ title: string, content: string, folder: string }} data - Updated note payload.
 * @returns {Promise<{ ok: boolean }>}
 */
export async function updateNote(id, data) {
    const res = await fetch(`${BASE}/api/notes/${id}`, {
        method: 'PUT',
        headers: { 
            'Content-Type': 'application/json',
            ...getAuthHeader()
        },
        body: JSON.stringify(data),
    });
    return parseResponse(res);
}

/**
 * Deletes a note by its ID.
 * @param {number} id - The ID of the note to delete.
 * @returns {Promise<{ ok: boolean }>}
 */
export async function deleteNote(id) {
    const res = await fetch(`${BASE}/api/notes/${id}`, {
        method: 'DELETE',
        headers: getAuthHeader(),
    });
    return parseResponse(res);
}

/**
 * Searches notes via FTS5 full-text search.
 * @param {string} query - The search query (max 128 chars, alphanumeric and common punctuation).
 * @returns {Promise<Note[]>}
 */
export async function searchNotes(query) {
    const res = await fetch(`${BASE}/api/search?q=${encodeURIComponent(query)}`, {
        headers: getAuthHeader(),
    });
    return parseResponse(res);
}

/**
 * Uploads an image file via multipart/form-data.
 * @param {File} file - The image file to upload.
 * @returns {Promise<{ url: string }>} The server-relative URL of the uploaded image.
 */
export async function uploadImage(file) {
    const form = new FormData();
    form.append('image', file);
    const res = await fetch(`${BASE}/api/upload`, {
        method: 'POST',
        headers: getAuthHeader(),
        body: form,
    });
    return parseResponse(res);
}
