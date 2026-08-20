// Mirrors icon-composer's compose/palette.go Gradients array — just for
// rendering swatches; the server is the source of truth for what a given
// gradient/seed actually composes into.
const GRADIENTS = [
  ["#ff5a3c", "#ffb627"],
  ["#ffb627", "#ff5a3c"],
  ["#7a5cff", "#ff5a3c"],
  ["#1fb37a", "#ffb627"],
  ["#ff3d81", "#ffb627"],
  ["#ff7a5c", "#7a5cff"],
];

const nameInput = document.getElementById("name");
const lettersInput = document.getElementById("letters");
const iconSearch = document.getElementById("icon-search");
const iconGrid = document.getElementById("icon-grid");
const swatchesEl = document.getElementById("swatches");
const previewBox = document.getElementById("preview-box");
const copyBtn = document.getElementById("copy-btn");
const copyStatus = document.getElementById("copy-status");

const state = {
  icon: "",
  lettersDirty: false,
  gradient: null, // null = auto (derived from name on the server)
  allIcons: [],
  currentSvg: "",
};

function debounce(fn, ms) {
  let t;
  return (...args) => {
    clearTimeout(t);
    t = setTimeout(() => fn(...args), ms);
  };
}

async function loadIcons() {
  const res = await fetch("/api/icons");
  const data = await res.json();
  state.allIcons = data.icons ?? [];
  renderIconGrid(state.allIcons);
}

function renderIconGrid(names) {
  iconGrid.innerHTML = "";

  const noneBtn = document.createElement("button");
  noneBtn.type = "button";
  noneBtn.className = "icon-btn none" + (state.icon === "" ? " selected" : "");
  noneBtn.textContent = "None";
  noneBtn.addEventListener("click", () => selectIcon(""));
  iconGrid.appendChild(noneBtn);

  for (const name of names) {
    const btn = document.createElement("button");
    btn.type = "button";
    btn.className = "icon-btn" + (state.icon === name ? " selected" : "");
    btn.title = name;
    const img = document.createElement("img");
    img.src = "/api/icon-glyph?name=" + encodeURIComponent(name);
    img.alt = name;
    btn.appendChild(img);
    btn.addEventListener("click", () => selectIcon(name));
    iconGrid.appendChild(btn);
  }
}

function selectIcon(name) {
  state.icon = name;
  for (const btn of iconGrid.querySelectorAll(".icon-btn")) {
    btn.classList.toggle(
      "selected",
      (btn.classList.contains("none") && name === "") || btn.title === name
    );
  }
  recompose();
}

function renderSwatches() {
  swatchesEl.innerHTML = "";

  const auto = document.createElement("button");
  auto.type = "button";
  auto.className = "swatch auto selected";
  auto.textContent = "auto";
  auto.title = "Derive from app name";
  auto.addEventListener("click", () => selectGradient(null, auto));
  swatchesEl.appendChild(auto);

  for (const [from, to] of GRADIENTS) {
    const el = document.createElement("button");
    el.type = "button";
    el.className = "swatch";
    el.style.background = `linear-gradient(135deg, ${from}, ${to})`;
    el.addEventListener("click", () => selectGradient([from, to], el));
    swatchesEl.appendChild(el);
  }
}

function selectGradient(gradient, el) {
  state.gradient = gradient;
  for (const s of swatchesEl.querySelectorAll(".swatch")) {
    s.classList.remove("selected");
  }
  el.classList.add("selected");
  recompose();
}

async function suggestLetters(name) {
  if (!name) return;
  const res = await fetch("/api/suggest?name=" + encodeURIComponent(name));
  const data = await res.json();
  if (!state.lettersDirty) {
    lettersInput.value = (data.letters ?? "").toUpperCase();
  }
  recompose();
}
const suggestLettersDebounced = debounce(suggestLetters, 200);

async function recompose() {
  const body = {
    icon: state.icon,
    letters: lettersInput.value,
    name: nameInput.value,
  };
  if (state.gradient) {
    body.from = state.gradient[0];
    body.to = state.gradient[1];
  }

  const res = await fetch("/api/compose", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify(body),
  });
  const data = await res.json();
  if (data.svg) {
    state.currentSvg = data.svg;
    previewBox.innerHTML = data.svg;
  }
}
const recomposeDebounced = debounce(recompose, 150);

nameInput.addEventListener("input", () => {
  suggestLettersDebounced(nameInput.value);
  recomposeDebounced();
});

lettersInput.addEventListener("input", () => {
  state.lettersDirty = lettersInput.value.trim() !== "";
  recomposeDebounced();
});

iconSearch.addEventListener("input", () => {
  const q = iconSearch.value.trim().toLowerCase();
  const filtered = q ? state.allIcons.filter((n) => n.includes(q)) : state.allIcons;
  renderIconGrid(filtered);
});

copyBtn.addEventListener("click", async () => {
  try {
    await navigator.clipboard.writeText(state.currentSvg);
    copyStatus.textContent = "Copied!";
  } catch {
    copyStatus.textContent = "Couldn't copy — select and copy manually.";
  }
  setTimeout(() => (copyStatus.textContent = ""), 2000);
});

renderSwatches();
loadIcons().then(recompose);
