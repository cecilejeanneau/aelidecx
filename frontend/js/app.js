import { fetchNotes, fetchNote, createNote, updateNote, deleteNote, searchNotes, setApiKey, hasApiKey, clearApiKey } from './api.js';
import { initEditor, getContent, setContent } from './editor.js';

let currentNoteId = null;
let currentFolder = '';
let searchTimeout = null;

// DOM elements
const notesList = document.getElementById('notes-list');
const titleInput = document.getElementById('note-title');
const folderSelect = document.getElementById('note-folder');
const saveBtn = document.getElementById('save-btn');
const deleteBtn = document.getElementById('delete-btn');
const newNoteBtn = document.getElementById('new-note-btn');
const searchInput = document.getElementById('search-input');
const editorEl = document.getElementById('editor');
const toolbarEl = document.getElementById('toolbar');
const folderBtns = document.querySelectorAll('.folder-btn');

// Init editor
initEditor(editorEl, toolbarEl);

// Check for existing API key and show auth modal or load notes
if (!hasApiKey()) {
    showAuthModal();
} else {
    loadNotesList();
}

// Event listeners
newNoteBtn.addEventListener('click', handleNewNote);
saveBtn.addEventListener('click', handleSave);
deleteBtn.addEventListener('click', handleDelete);
searchInput.addEventListener('input', handleSearch);

folderBtns.forEach(btn => {
    btn.addEventListener('click', () => {
        folderBtns.forEach(b => b.classList.remove('active'));
        btn.classList.add('active');
        currentFolder = btn.dataset.folder;
        loadNotesList();
    });
});

/**
 * Fetches the note list (filtered by the active folder) and renders it in the sidebar.
 */
async function loadNotesList() {
    try {
        const notes = await fetchNotes(currentFolder);
        renderNotesList(notes);
    } catch (err) {
        handleApiError(err);
    }
}

/**
 * Renders a flat list of notes into the sidebar.
 * @param {Note[]} notes - The notes to display.
 */
function renderNotesList(notes) {
    notesList.innerHTML = '';
    for (const note of notes) {
        const li = document.createElement('li');
        li.dataset.id = note.id;
        if (note.id === currentNoteId) li.classList.add('active');

        const title = document.createElement('div');
        title.className = 'note-item-title';
        title.textContent = note.title || 'Untitled';

        const date = document.createElement('div');
        date.className = 'note-item-date';
        date.textContent = formatDate(note.updated_at);

        li.appendChild(title);
        li.appendChild(date);
        li.addEventListener('click', () => openNote(note.id));
        notesList.appendChild(li);
    }
}

/**
 * Fetches a note by ID and populates the editor with its content.
 * @param {number} id - The note ID to open.
 */
async function openNote(id) {
    try {
        const note = await fetchNote(id);
        currentNoteId = note.id;
        titleInput.value = note.title;
        folderSelect.value = note.folder;
        setContent(note.content);

        document.querySelectorAll('#notes-list li').forEach(li => {
            li.classList.toggle('active', parseInt(li.dataset.id) === id);
        });
    } catch (err) {
        handleApiError(err);
    }
}

/**
 * Resets the editor to a blank state for composing a new note.
 */
function handleNewNote() {
    currentNoteId = null;
    titleInput.value = '';
    folderSelect.value = currentFolder || 'inbox';
    setContent('');
    titleInput.focus();
}

/**
 * Saves the current note (creates a new one or updates the existing one).
 * Trims the title and folder before sending to the API.
 */
async function handleSave() {
    const cleanTitle = titleInput.value.trim();
    const cleanFolder = String(folderSelect.value || '').trim();

    const data = {
        title: cleanTitle || 'Untitled',
        content: getContent(),
        folder: cleanFolder,
    };

    try {
        if (currentNoteId) {
            await updateNote(currentNoteId, data);
        } else {
            const result = await createNote(data);
            currentNoteId = result.id;
        }

        await loadNotesList();
    } catch (err) {
        handleApiError(err);
    }
}

/**
 * Deletes the currently open note after user confirmation.
 */
async function handleDelete() {
    if (!currentNoteId) return;
    if (!confirm('Delete this note?')) return;

    try {
        await deleteNote(currentNoteId);
        currentNoteId = null;
        titleInput.value = '';
        setContent('');
        await loadNotesList();
    } catch (err) {
        handleApiError(err);
    }
}

/**
 * Handles the search input with 300 ms debounce.
 * Clears search and reloads the full list when the input is empty.
 */
function handleSearch() {
    clearTimeout(searchTimeout);
    const q = searchInput.value.trim();

    if (!q) {
        loadNotesList();
        return;
    }

    searchTimeout = setTimeout(async () => {
        try {
            const results = await searchNotes(q);
            renderNotesList(results);
        } catch (err) {
            handleApiError(err);
        }
    }, 300);
}

/**
 * Formats an ISO date string as a short human-readable date/time.
 * @param {string} dateStr - UTC date string from the API.
 * @returns {string}
 */
function formatDate(dateStr) {
    if (!dateStr) return '';
    const d = new Date(dateStr + 'Z');
    return d.toLocaleDateString('en-GB', { day: '2-digit', month: 'short', hour: '2-digit', minute: '2-digit' });
}

// ========== AUTHENTICATION ==========

/**
 * Displays the authentication modal and handles the login flow.
 * Validates the entered API key by making a test request to the API.
 */
function showAuthModal() {
    const modal = document.getElementById('auth-modal');
    const input = document.getElementById('auth-input');
    const btn = document.getElementById('auth-btn');
    const error = document.getElementById('auth-error');

    modal.classList.remove('hidden');
    input.focus();

    const handleAuth = async () => {
        const key = input.value.trim();
        if (!key) {
            error.textContent = 'Please enter an API key';
            error.classList.remove('hidden');
            return;
        }

        // Try loading notes with the provided key
        setApiKey(key);
        try {
            await fetchNotes();
            // Success
            hideAuthModal();
            loadNotesList();
        } catch (err) {
            // Authentication error
            error.textContent = 'Invalid API key or access denied';
            error.classList.remove('hidden');
            clearApiKey();
        }
    };

    btn.onclick = handleAuth;
    input.onkeypress = (e) => {
        if (e.key === 'Enter') handleAuth();
    };
}

/**
 * Hides the authentication modal.
 */
function hideAuthModal() {
    const modal = document.getElementById('auth-modal');
    modal.classList.add('hidden');
}

/**
 * Handles API errors globally.
 * Redirects to the auth modal on 401/403; shows a generic alert for other errors.
 * @param {Error & { status?: number }} err - The error thrown by the API layer.
 */
function handleApiError(err) {
    if (err && (err.status === 401 || err.status === 403)) {
        clearApiKey();
        showAuthModal();
        return;
    }
    alert('An error occurred. Please try again.');
}
