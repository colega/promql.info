// Source: https://discuss.codemirror.net/t/codemirror-6-and-textareas/2731/3

import { EditorState } from "@codemirror/state";
import { EditorView, keymap, lineNumbers } from "@codemirror/view";
import { defaultKeymap } from "@codemirror/commands";

// Themes from https://thememirror.net/
import {espresso} from 'thememirror';

// Find the submit button
const submitButton = document.querySelector("#submit");

// Hide the textarea used as the editor's source
const editorSource = document.querySelector("#query");
editorSource.setAttribute("hidden", "true");

// Find the editor targtet
const editorTarget = document.querySelector("#editor");

// Set up Code mirror
let startState = EditorState.create({
    // Use the textarea's value as the initial text for CodeMirror
    doc: editorSource.value,
    extensions: [keymap.of(defaultKeymap), lineNumbers(), espresso],
});

let view = new EditorView({
    state: startState,
    parent: editorTarget,
});

// Sync the contents of CodeMirror with the textarea
const syncEditor = () => {
    // Is there a better way to get CodeMirror's contents?
    // Not sure, but `view.state.sliceDoc()` was the easiest way I found.
    editorSource.value = view.state.sliceDoc();
    console.log(editorSource.value);
};

submitButton.addEventListener("click", syncEditor);