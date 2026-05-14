import { uploadImage } from './api.js';

let editorEl = null;

/**
 * Initialises the rich-text editor on the given contenteditable element.
 * Attaches drag-and-drop and paste handlers for inline image uploading,
 * then builds the formatting toolbar.
 * @param {HTMLElement} element - The contenteditable host element.
 * @param {HTMLElement} toolbarElement - The container element for toolbar buttons.
 * @returns {HTMLElement} The initialised editor element.
 */
export function initEditor(element, toolbarElement) {
    editorEl = element;
    editorEl.contentEditable = 'true';
    editorEl.setAttribute('role', 'textbox');
    editorEl.setAttribute('aria-multiline', 'true');
    editorEl.setAttribute('data-placeholder', 'Write your note here… (paste or drop images directly!)');
    editorEl.classList.add('editor-content');

    editorEl.addEventListener('drop', (event) => {
        const files = event.dataTransfer?.files;
        if (files && files.length > 0) {
            event.preventDefault();
            handleImageFiles(files);
        }
    });

    editorEl.addEventListener('paste', (event) => {
        const items = event.clipboardData?.items;
        if (!items) return;
        for (const item of items) {
            if (item.type.startsWith('image/')) {
                event.preventDefault();
                const file = item.getAsFile();
                if (file) handleImageFiles([file]);
                return;
            }
        }
    });

    buildToolbar(toolbarElement);
    return editorEl;
}

/**
 * Uploads a list of image files and inserts them as <img> tags at the cursor position.
 * @param {FileList|File[]} files - The image files to upload.
 */
async function handleImageFiles(files) {
    for (const file of files) {
        const result = await uploadImage(file);
        if (result.url) {
            insertHtmlAtCursor(`<img src="${result.url}" alt="image" />`);
        }
    }
}

/**
 * Executes a contenteditable formatting command via document.execCommand.
 * @param {string} command - The execCommand name (e.g. 'bold', 'formatBlock').
 * @param {string|null} [value=null] - Optional value for commands that require one.
 */
function exec(command, value = null) {
    if (!editorEl) return;
    editorEl.focus();
    document.execCommand(command, false, value);
}

/**
 * Inserts an arbitrary HTML string at the current cursor position in the editor.
 * @param {string} html - The HTML to insert.
 */
function insertHtmlAtCursor(html) {
    exec('insertHTML', html);
}

/**
 * Builds the formatting toolbar and appends buttons to the given container.
 * @param {HTMLElement} container - The element to render toolbar buttons into.
 */
function buildToolbar(container) {
    const buttons = [
        { label: 'B', command: () => exec('bold') },
        { label: 'I', command: () => exec('italic') },
        { label: 'H1', command: () => exec('formatBlock', 'h1') },
        { label: 'H2', command: () => exec('formatBlock', 'h2') },
        { label: 'H3', command: () => exec('formatBlock', 'h3') },
        { label: '•', command: () => exec('insertUnorderedList') },
        { label: '1.', command: () => exec('insertOrderedList') },
        { label: '""', command: () => exec('formatBlock', 'blockquote') },
        { label: '<>', command: () => exec('formatBlock', 'pre') },
        { label: '🔗', command: () => {
            const url = prompt('URL:');
            if (url) exec('createLink', url.trim());
        }},
        { label: '📷', command: () => {
            const input = document.createElement('input');
            input.type = 'file';
            input.accept = 'image/*';
            input.onchange = () => handleImageFiles(input.files);
            input.click();
        }},
    ];

    container.innerHTML = '';
    for (const btn of buttons) {
        const el = document.createElement('button');
        el.textContent = btn.label;
        el.type = 'button';
        el.addEventListener('click', btn.command);
        container.appendChild(el);
    }
}

/**
 * Returns the current HTML content of the editor.
 * @returns {string}
 */
export function getContent() {
    return editorEl ? editorEl.innerHTML : '';
}

/**
 * Sets the HTML content of the editor, replacing any existing content.
 * @param {string} html - The HTML string to load into the editor.
 */
export function setContent(html) {
    if (editorEl) editorEl.innerHTML = html || '';
}

/**
 * Destroys the editor instance by replacing the element with an inert clone,
 * removing all attached event listeners.
 */
export function destroyEditor() {
    if (editorEl) {
        editorEl.replaceWith(editorEl.cloneNode(true));
        editorEl = null;
    }
}
