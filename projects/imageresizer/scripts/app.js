/**
 * Image Resizer — Main Application Logic
 * Handles: upload, state management, controls, preview, download.
 * All DOM rendering uses textContent (never innerHTML) per security spec.
 */

(function () {
  'use strict';

  // ---- State ----
  const state = {
    originalFile: null,
    originalImage: null,  // HTMLImageElement
    originalWidth: 0,
    originalHeight: 0,
    originalSize: 0,
    originalFormat: 'jpeg',
    width: 0,
    height: 0,
    aspectRatio: 1,
    aspectLocked: true,
    format: 'jpeg',
    quality: 0.85,
    resizedBlob: null,
    resizedBlobUrl: null,
    isProcessing: false,
    activePreset: null,  // currently selected preset or scale button
  };

  // ---- DOM Elements ----
  const dom = {
    uploadView: document.getElementById('uploadView'),
    editorView: document.getElementById('editorView'),
    dropzone: document.getElementById('dropzone'),
    dropzoneHeading: document.getElementById('dropzoneHeading'),
    dropzoneSubtext: document.getElementById('dropzoneSubtext'),
    dropzoneError: document.getElementById('dropzoneError'),
    fileInput: document.getElementById('fileInput'),
    newImageBtn: document.getElementById('newImageBtn'),
    previewImage: document.getElementById('previewImage'),
    previewLoading: document.getElementById('previewLoading'),
    progressFill: document.getElementById('progressFill'),
    previewFilename: document.getElementById('previewFilename'),
    previewDimensions: document.getElementById('previewDimensions'),
    previewSize: document.getElementById('previewSize'),
    widthInput: document.getElementById('widthInput'),
    heightInput: document.getElementById('heightInput'),
    aspectLockBtn: document.getElementById('aspectLockBtn'),
    aspectLockLabel: document.getElementById('aspectLockLabel'),
    lockIcon: document.getElementById('lockIcon'),
    unlockIcon: document.getElementById('unlockIcon'),
    formatSelect: document.getElementById('formatSelect'),
    qualityGroup: document.getElementById('qualityGroup'),
    qualitySlider: document.getElementById('qualitySlider'),
    qualityValue: document.getElementById('qualityValue'),
    downloadBar: document.getElementById('downloadBar'),
    downloadBtn: document.getElementById('downloadBtn'),
    originalSize: document.getElementById('originalSize'),
    resizedSize: document.getElementById('resizedSize'),
    sizePercent: document.getElementById('sizePercent'),
    statusRegion: document.getElementById('statusRegion'),
    presetGrid: document.getElementById('presetGrid'),
    scaleGrid: document.getElementById('scaleGrid'),
    footer: document.getElementById('appFooter'),
  };

  // ---- View Switching ----
  function showUploadView() {
    dom.uploadView.classList.add('active');
    dom.editorView.classList.remove('active');
    dom.downloadBar.classList.remove('active');
    dom.newImageBtn.classList.add('hidden');
    if (dom.footer) dom.footer.classList.remove('hidden');
    resetDropzone();
  }

  function showEditorView() {
    dom.uploadView.classList.remove('active');
    dom.editorView.classList.add('active');
    dom.downloadBar.classList.add('active');
    dom.newImageBtn.classList.remove('hidden');
    if (dom.footer) dom.footer.classList.add('hidden');
  }

  // ---- Dropzone ----
  function resetDropzone() {
    dom.dropzoneHeading.textContent = 'Drop your image here';
    dom.dropzoneSubtext.textContent = 'or click to browse';
    dom.dropzoneError.textContent = '';
    dom.dropzone.classList.remove('dropzone--hover', 'dropzone--error');
  }

  function showDropzoneError(message) {
    dom.dropzoneError.textContent = message;
    dom.dropzone.classList.add('dropzone--error');
    announce(message);
    setTimeout(function () {
      dom.dropzoneError.textContent = '';
      dom.dropzone.classList.remove('dropzone--error');
    }, 3000);
  }

  // ---- File Handling ----
  function handleFile(file) {
    const validation = Utils.validateFile(file);
    if (!validation.valid) {
      showDropzoneError(validation.error);
      return;
    }

    state.originalFile = file;
    state.originalSize = file.size;
    state.originalFormat = Utils.mimeToFormat(file.type);

    // Load image
    const img = new Image();
    const objectUrl = URL.createObjectURL(file);

    img.onload = function () {
      state.originalImage = img;
      state.originalWidth = img.naturalWidth;
      state.originalHeight = img.naturalHeight;
      state.width = img.naturalWidth;
      state.height = img.naturalHeight;
      state.aspectRatio = img.naturalWidth / img.naturalHeight;

      // Set format to match original
      state.format = state.originalFormat;
      dom.formatSelect.value = state.format;
      updateQualityVisibility();

      // Populate controls
      dom.widthInput.value = state.width;
      dom.heightInput.value = state.height;

      // Populate preview info (textContent only)
      dom.previewFilename.textContent = file.name;
      dom.previewDimensions.textContent = state.originalWidth + ' × ' + state.originalHeight;
      dom.previewSize.textContent = Utils.formatFileSize(state.originalSize);

      // Show original size in download bar
      dom.originalSize.textContent = Utils.formatFileSize(state.originalSize);

      // Show editor
      showEditorView();

      // Revoke the object URL used for loading
      URL.revokeObjectURL(objectUrl);

      // Trigger initial preview
      updatePreview();

      announce('Image loaded: ' + file.name + ', ' + state.originalWidth + ' by ' + state.originalHeight + ' pixels.');
    };

    img.onerror = function () {
      URL.revokeObjectURL(objectUrl);
      showDropzoneError("Something went wrong loading that image. Try a different file?");
    };

    img.src = objectUrl;
  }

  // ---- Drag and Drop ----
  function setupDropzone() {
    dom.dropzone.addEventListener('dragover', function (e) {
      e.preventDefault();
      e.stopPropagation();
      dom.dropzone.classList.add('dropzone--hover');
      dom.dropzoneHeading.textContent = 'Let go to resize!';
      dom.dropzoneSubtext.textContent = '';
    });

    dom.dropzone.addEventListener('dragleave', function (e) {
      e.preventDefault();
      e.stopPropagation();
      dom.dropzone.classList.remove('dropzone--hover');
      dom.dropzoneHeading.textContent = 'Drop your image here';
      dom.dropzoneSubtext.textContent = 'or click to browse';
    });

    dom.dropzone.addEventListener('drop', function (e) {
      e.preventDefault();
      e.stopPropagation();
      dom.dropzone.classList.remove('dropzone--hover');
      dom.dropzoneHeading.textContent = 'Drop your image here';
      dom.dropzoneSubtext.textContent = 'or click to browse';

      var files = e.dataTransfer.files;
      if (files.length > 0) {
        handleFile(files[0]);
      }
    });

    // Click to browse
    dom.dropzone.addEventListener('click', function () {
      dom.fileInput.click();
    });

    // Keyboard: Enter/Space to open file picker
    dom.dropzone.addEventListener('keydown', function (e) {
      if (e.key === 'Enter' || e.key === ' ') {
        e.preventDefault();
        dom.fileInput.click();
      }
    });

    dom.fileInput.addEventListener('change', function () {
      if (dom.fileInput.files.length > 0) {
        handleFile(dom.fileInput.files[0]);
      }
    });
  }

  // ---- Aspect Ratio Lock ----
  function toggleAspectLock() {
    state.aspectLocked = !state.aspectLocked;
    dom.aspectLockBtn.setAttribute('aria-pressed', String(state.aspectLocked));

    if (state.aspectLocked) {
      dom.aspectLockBtn.classList.add('aspect-lock__btn--locked');
      dom.lockIcon.classList.remove('hidden');
      dom.unlockIcon.classList.add('hidden');
      dom.aspectLockLabel.textContent = 'Ratio locked';
      // Recalculate aspect ratio from current dimensions
      state.aspectRatio = state.width / state.height;
    } else {
      dom.aspectLockBtn.classList.remove('aspect-lock__btn--locked');
      dom.lockIcon.classList.add('hidden');
      dom.unlockIcon.classList.remove('hidden');
      dom.aspectLockLabel.textContent = 'Ratio unlocked';
    }
  }

  // ---- Preset & Scale Buttons ----
  function clearActivePreset() {
    state.activePreset = null;
    var allBtns = document.querySelectorAll('.preset-btn');
    for (var i = 0; i < allBtns.length; i++) {
      allBtns[i].classList.remove('preset-btn--active');
    }
  }

  function setActivePreset(btn) {
    clearActivePreset();
    state.activePreset = btn;
    btn.classList.add('preset-btn--active');
  }

  function onPresetClick(e) {
    var btn = e.target.closest('.preset-btn');
    if (!btn || !state.originalImage) return;

    var w = parseInt(btn.getAttribute('data-w'), 10);
    var h = parseInt(btn.getAttribute('data-h'), 10);
    if (isNaN(w) || isNaN(h)) return;

    setActivePreset(btn);

    // Presets override aspect lock — unlock so dimensions are exact
    if (state.aspectLocked) {
      toggleAspectLock();
    }

    state.width = w;
    state.height = h;
    dom.widthInput.value = w;
    dom.heightInput.value = h;

    debouncedPreview();
    announce('Preset selected: ' + btn.textContent.trim() + ', ' + w + ' by ' + h + ' pixels.');
  }

  function onScaleClick(e) {
    var btn = e.target.closest('.preset-btn');
    if (!btn || !state.originalImage) return;

    var scale = parseInt(btn.getAttribute('data-scale'), 10);
    if (isNaN(scale) || scale < 1) return;

    setActivePreset(btn);

    var w = Math.max(1, Math.round(state.originalWidth * scale / 100));
    var h = Math.max(1, Math.round(state.originalHeight * scale / 100));

    state.width = w;
    state.height = h;
    dom.widthInput.value = w;
    dom.heightInput.value = h;

    // Re-lock aspect ratio since scale preserves proportions
    if (!state.aspectLocked) {
      toggleAspectLock();
    }
    state.aspectRatio = state.originalWidth / state.originalHeight;

    debouncedPreview();
    announce('Scaled to ' + scale + '%: ' + w + ' by ' + h + ' pixels.');
  }

  // ---- Dimension Inputs ----
  function onWidthChange() {
    var w = parseInt(dom.widthInput.value, 10);
    if (isNaN(w) || w < 1) return;
    state.width = w;

    // Custom input deselects any active preset
    clearActivePreset();

    if (state.aspectLocked) {
      state.height = Math.max(1, Math.round(w / state.aspectRatio));
      dom.heightInput.value = state.height;
    }

    debouncedPreview();
  }

  function onHeightChange() {
    var h = parseInt(dom.heightInput.value, 10);
    if (isNaN(h) || h < 1) return;
    state.height = h;

    // Custom input deselects any active preset
    clearActivePreset();

    if (state.aspectLocked) {
      state.width = Math.max(1, Math.round(h * state.aspectRatio));
      dom.widthInput.value = state.width;
    }

    debouncedPreview();
  }

  // ---- Format & Quality ----
  function onFormatChange() {
    state.format = dom.formatSelect.value;
    updateQualityVisibility();
    debouncedPreview();
  }

  function updateQualityVisibility() {
    if (state.format === 'png') {
      dom.qualityGroup.classList.add('hidden');
    } else {
      dom.qualityGroup.classList.remove('hidden');
    }
  }

  function onQualityChange() {
    state.quality = parseInt(dom.qualitySlider.value, 10) / 100;
    dom.qualityValue.textContent = dom.qualitySlider.value + '%';
    debouncedPreview();
  }

  // ---- Live Preview ----
  var debouncedPreview = Utils.debounce(updatePreview, 300);

  function updatePreview() {
    if (!state.originalImage || state.isProcessing) return;

    var w = parseInt(dom.widthInput.value, 10);
    var h = parseInt(dom.heightInput.value, 10);
    if (isNaN(w) || isNaN(h) || w < 1 || h < 1) return;

    state.width = w;
    state.height = h;
    state.isProcessing = true;

    // Show loading
    dom.previewLoading.classList.add('active');
    dom.previewImage.style.display = 'none';
    dom.downloadBtn.disabled = true;
    setDownloadButtonText('Processing...');
    dom.progressFill.style.width = '30%';

    Resizer.resize(state.originalImage, {
      width: state.width,
      height: state.height,
      format: state.format,
      quality: state.quality,
    })
      .then(function (blob) {
        // Revoke old URL
        if (state.resizedBlobUrl) {
          URL.revokeObjectURL(state.resizedBlobUrl);
        }

        state.resizedBlob = blob;
        state.resizedBlobUrl = URL.createObjectURL(blob);

        dom.progressFill.style.width = '100%';

        // Show preview
        dom.previewImage.src = state.resizedBlobUrl;
        dom.previewImage.style.display = '';
        dom.previewImage.style.animation = 'fadeIn 0.2s ease';

        // Update download bar info
        dom.resizedSize.textContent = Utils.formatFileSize(blob.size);
        var change = Utils.calcSizeChange(state.originalSize, blob.size);
        dom.sizePercent.textContent = '(' + change.label + ')';
        dom.sizePercent.className = 'download-bar__percent ' +
          (change.smaller ? 'download-bar__percent--smaller' : 'download-bar__percent--larger');

        // Hide loading
        setTimeout(function () {
          dom.previewLoading.classList.remove('active');
          dom.progressFill.style.width = '0%';
        }, 150);

        // Enable download with pulse animation
        dom.downloadBtn.disabled = false;
        setDownloadButtonText('Download');
        dom.downloadBtn.classList.add('download-bar__btn--pulse');
        setTimeout(function () {
          dom.downloadBtn.classList.remove('download-bar__btn--pulse');
        }, 300);

        state.isProcessing = false;
      })
      .catch(function (err) {
        dom.previewLoading.classList.remove('active');
        dom.progressFill.style.width = '0%';
        dom.downloadBtn.disabled = false;
        setDownloadButtonText('Download');
        state.isProcessing = false;

        announce(err.message || 'Something went wrong processing the image.');
      });
  }

  function setDownloadButtonText(text) {
    // The button has an SVG icon + text. We only change the text node.
    var btn = dom.downloadBtn;
    var lastChild = btn.lastChild;
    if (lastChild && lastChild.nodeType === Node.TEXT_NODE) {
      lastChild.textContent = ' ' + text;
    }
  }

  // ---- Download ----
  function downloadImage() {
    if (!state.resizedBlob || state.isProcessing) return;

    var filename = Utils.buildDownloadFilename(state.originalFile.name, state.format);
    var url = URL.createObjectURL(state.resizedBlob);

    var a = document.createElement('a');
    a.href = url;
    a.download = filename;
    document.body.appendChild(a);
    a.click();
    document.body.removeChild(a);

    // Revoke after a short delay to ensure download starts
    setTimeout(function () {
      URL.revokeObjectURL(url);
    }, 1000);

    // Success feedback
    setDownloadButtonText('Downloaded!');
    dom.downloadBtn.classList.add('download-bar__btn--success');
    announce('Image downloaded as ' + filename);

    setTimeout(function () {
      setDownloadButtonText('Download');
      dom.downloadBtn.classList.remove('download-bar__btn--success');
    }, 2000);
  }

  // ---- New Image / Reset ----
  function resetApp() {
    // Revoke blob URLs
    if (state.resizedBlobUrl) {
      URL.revokeObjectURL(state.resizedBlobUrl);
    }

    // Reset state
    state.originalFile = null;
    state.originalImage = null;
    state.originalWidth = 0;
    state.originalHeight = 0;
    state.originalSize = 0;
    state.originalFormat = 'jpeg';
    state.width = 0;
    state.height = 0;
    state.aspectRatio = 1;
    state.aspectLocked = true;
    state.format = 'jpeg';
    state.quality = 0.85;
    state.resizedBlob = null;
    state.resizedBlobUrl = null;
    state.isProcessing = false;

    // Reset controls
    dom.widthInput.value = '';
    dom.heightInput.value = '';
    dom.formatSelect.value = 'jpeg';
    dom.qualitySlider.value = 85;
    dom.qualityValue.textContent = '85%';
    dom.qualityGroup.classList.remove('hidden');
    dom.previewImage.src = '';
    dom.fileInput.value = '';

    // Reset aspect lock to locked
    dom.aspectLockBtn.classList.add('aspect-lock__btn--locked');
    dom.aspectLockBtn.setAttribute('aria-pressed', 'true');
    dom.lockIcon.classList.remove('hidden');
    dom.unlockIcon.classList.add('hidden');
    dom.aspectLockLabel.textContent = 'Ratio locked';

    // Clear active preset
    clearActivePreset();

    // Switch view
    showUploadView();
    announce('Ready for a new image.');
  }

  // ---- Accessibility ----
  function announce(message) {
    dom.statusRegion.textContent = message;
  }

  // ---- Event Binding ----
  function init() {
    setupDropzone();

    dom.widthInput.addEventListener('input', onWidthChange);
    dom.heightInput.addEventListener('input', onHeightChange);
    dom.aspectLockBtn.addEventListener('click', toggleAspectLock);
    dom.formatSelect.addEventListener('change', onFormatChange);
    dom.qualitySlider.addEventListener('input', onQualityChange);
    dom.downloadBtn.addEventListener('click', downloadImage);
    dom.newImageBtn.addEventListener('click', resetApp);

    // Preset buttons (event delegation)
    if (dom.presetGrid) {
      dom.presetGrid.addEventListener('click', onPresetClick);
    }

    // Scale buttons (event delegation)
    if (dom.scaleGrid) {
      dom.scaleGrid.addEventListener('click', onScaleClick);
    }

    // Prevent page-level drag/drop
    document.addEventListener('dragover', function (e) { e.preventDefault(); });
    document.addEventListener('drop', function (e) { e.preventDefault(); });

    // Start in upload view
    showUploadView();
  }

  // Go
  init();
})();
