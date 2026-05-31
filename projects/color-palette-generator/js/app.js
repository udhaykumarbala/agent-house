/**
 * Main application logic
 * State management, DOM rendering, event handling
 */
(function() {
  'use strict';

  // App state — single source of truth
  var state = {
    palette: [],
    harmonyMode: 'random',
    colorFormat: 'hex'
  };

  // DOM references
  var swatches = document.querySelectorAll('.swatch');
  var generateBtn = document.getElementById('generate-btn');
  var formatBtns = document.querySelectorAll('.format-btn');
  var harmonySelect = document.getElementById('harmony-select');
  var lockBtns = document.querySelectorAll('.lock-btn');
  var copyAllBtn = document.getElementById('copy-all-btn');
  var srAnnounce = document.getElementById('sr-announce');

  /**
   * Announce a message to screen readers via the aria-live region
   */
  function announce(message) {
    if (!srAnnounce) return;
    srAnnounce.textContent = '';
    requestAnimationFrame(function() {
      srAnnounce.textContent = message;
    });
  }

  /**
   * Render palette state to DOM
   * All updates use textContent — never innerHTML
   */
  function render() {
    state.palette.forEach(function(color, i) {
      var swatch = swatches[i];
      if (!swatch) return;

      var codeEl = swatch.querySelector('.swatch-color-code');
      swatch.style.backgroundColor = color.hex;

      var contrastColor = getContrastColor(color.hex);
      codeEl.style.color = contrastColor;
      codeEl.textContent = formatColor(color, state.colorFormat);

      // Set dark/light swatch class for contrast-aware elements
      var rgb = hexToRgb(color.hex);
      var luminance = 0.299 * rgb.r + 0.587 * rgb.g + 0.114 * rgb.b;
      if (luminance > 150) {
        swatch.classList.add('light-swatch');
        swatch.classList.remove('dark-swatch');
      } else {
        swatch.classList.add('dark-swatch');
        swatch.classList.remove('light-swatch');
      }

      swatch.setAttribute('aria-label',
        'Color ' + (i + 1) + ': ' + color.hex + '. ' + (color.locked ? 'Locked' : 'Unlocked')
      );

      // Update lock button state
      var lockBtn = swatch.querySelector('.lock-btn');
      if (lockBtn) {
        if (color.locked) {
          swatch.classList.add('locked');
          lockBtn.style.color = contrastColor;
          lockBtn.setAttribute('aria-label', 'Unlock color ' + (i + 1) + ' ' + color.hex);
          lockBtn.setAttribute('aria-pressed', 'true');
        } else {
          swatch.classList.remove('locked');
          lockBtn.style.color = contrastColor;
          lockBtn.setAttribute('aria-label', 'Lock color ' + (i + 1) + ' ' + color.hex);
          lockBtn.setAttribute('aria-pressed', 'false');
        }
      }
    });

    // Update generate button muted state when all locked
    var allLocked = state.palette.length > 0 && state.palette.every(function(c) { return c.locked; });
    generateBtn.classList.toggle('muted', allLocked);
  }

  /**
   * Generate a new palette and render
   */
  function generate() {
    var allLocked = state.palette.length > 0 && state.palette.every(function(c) { return c.locked; });
    if (allLocked) {
      showToast('Unlock a color to generate new ones');
      return;
    }
    state.palette = generatePalette(state.harmonyMode, state.palette);
    render();
    announce('New palette generated');
  }

  /**
   * Handle swatch click — copy color to clipboard from data model
   */
  function handleSwatchClick(e) {
    // Don't copy when clicking the lock button
    if (e.target.closest('.lock-btn')) return;

    var swatch = e.target.closest('.swatch');
    if (!swatch) return;

    var index = parseInt(swatch.dataset.index, 10);
    var color = state.palette[index];
    if (!color) return;

    var displayValue = formatColor(color, state.colorFormat);
    copyToClipboard(displayValue);
    showToast('Copied ' + displayValue + '!');
    announce('Copied ' + displayValue);
  }

  /**
   * Handle lock button click — toggle locked state
   */
  function handleLockClick(e) {
    e.stopPropagation();
    var btn = e.currentTarget;
    var index = parseInt(btn.dataset.index, 10);
    var color = state.palette[index];
    if (!color) return;

    color.locked = !color.locked;
    render();
    announce('Color ' + (index + 1) + (color.locked ? ' locked' : ' unlocked'));
  }

  /**
   * Handle harmony mode change
   */
  function handleHarmonyChange(e) {
    state.harmonyMode = e.target.value;
    generate();
  }

  /**
   * Get the currently focused or hovered swatch index
   * Returns -1 if no swatch is focused/hovered
   */
  function getActiveSwatch() {
    // Check focused swatch first
    var focused = document.activeElement;
    if (focused && focused.classList.contains('swatch')) {
      return parseInt(focused.dataset.index, 10);
    }
    // Check hovered swatch
    var hovered = document.querySelector('.swatch:hover');
    if (hovered) {
      return parseInt(hovered.dataset.index, 10);
    }
    return -1;
  }

  /**
   * Handle keyboard events
   */
  function handleKeydown(e) {
    var tag = e.target.tagName;
    if (tag === 'INPUT' || tag === 'TEXTAREA' || tag === 'SELECT') return;

    // Don't fire shortcuts when interacting with toolbar buttons
    var toolbar = document.querySelector('.toolbar');
    if (toolbar && toolbar.contains(e.target) && tag === 'BUTTON') return;

    if (e.code === 'Space') {
      e.preventDefault();
      generate();
      return;
    }

    // L key — toggle lock on focused/hovered swatch
    if (e.code === 'KeyL') {
      var lockIndex = getActiveSwatch();
      if (lockIndex >= 0 && state.palette[lockIndex]) {
        state.palette[lockIndex].locked = !state.palette[lockIndex].locked;
        render();
        announce('Color ' + (lockIndex + 1) + (state.palette[lockIndex].locked ? ' locked' : ' unlocked'));
      }
      return;
    }

    // C key — copy focused/hovered swatch color
    if (e.code === 'KeyC') {
      var copyIndex = getActiveSwatch();
      if (copyIndex >= 0 && state.palette[copyIndex]) {
        var displayValue = formatColor(state.palette[copyIndex], state.colorFormat);
        copyToClipboard(displayValue);
        showToast('Copied ' + displayValue + '!');
        announce('Copied ' + displayValue);
      }
      return;
    }

    // Keys 1-5 — focus corresponding swatch
    if (e.code >= 'Digit1' && e.code <= 'Digit5') {
      var swatchIndex = parseInt(e.code.charAt(5), 10) - 1;
      if (swatches[swatchIndex]) {
        swatches[swatchIndex].focus();
      }
      return;
    }
  }

  /**
   * Handle format toggle click — switch color display format
   */
  function handleFormatToggle(e) {
    var btn = e.target.closest('.format-btn');
    if (!btn) return;

    var format = btn.dataset.format;
    if (format === state.colorFormat) return;

    state.colorFormat = format;

    formatBtns.forEach(function(b) {
      var isActive = b.dataset.format === format;
      b.classList.toggle('active', isActive);
      b.setAttribute('aria-checked', isActive ? 'true' : 'false');
    });

    render();
  }

  /**
   * Handle Copy All — copy entire palette to clipboard in current format
   */
  function handleCopyAll() {
    if (state.palette.length === 0) return;
    var values = state.palette.map(function(c) {
      return formatColor(c, state.colorFormat);
    });
    copyToClipboard(values.join(', '));
    showToast('Palette copied!');
    announce('Palette copied to clipboard');
  }

  // Event listeners
  generateBtn.addEventListener('click', generate);

  swatches.forEach(function(swatch) {
    swatch.addEventListener('click', handleSwatchClick);
  });

  lockBtns.forEach(function(btn) {
    btn.addEventListener('click', handleLockClick);
  });

  harmonySelect.addEventListener('change', handleHarmonyChange);

  if (copyAllBtn) {
    copyAllBtn.addEventListener('click', handleCopyAll);
  }

  formatBtns.forEach(function(btn) {
    btn.addEventListener('click', handleFormatToggle);
  });

  document.addEventListener('keydown', handleKeydown);

  // Generate initial palette on page load
  generate();
})();
