/**
 * Clipboard copy logic and toast notifications
 */

var _toastTimer = null;

/**
 * Copy text to clipboard
 * Uses Clipboard API with execCommand fallback
 */
function copyToClipboard(text) {
  if (navigator.clipboard && navigator.clipboard.writeText) {
    navigator.clipboard.writeText(text).catch(function() {
      _fallbackCopy(text);
    });
  } else {
    _fallbackCopy(text);
  }
}

function _fallbackCopy(text) {
  var textarea = document.createElement('textarea');
  textarea.value = text;
  textarea.setAttribute('readonly', '');
  textarea.style.position = 'fixed';
  textarea.style.left = '-9999px';
  textarea.style.opacity = '0';
  document.body.appendChild(textarea);
  textarea.select();
  try {
    document.execCommand('copy');
  } catch (e) {
    // Silent fail
  }
  document.body.removeChild(textarea);
}

/**
 * Show toast notification — auto-dismisses after 1.5s
 * Replaces any existing toast (no stacking)
 */
function showToast(message) {
  var toast = document.getElementById('toast');
  if (!toast) return;

  if (_toastTimer) {
    clearTimeout(_toastTimer);
  }

  toast.textContent = message;
  toast.classList.add('toast-visible');

  _toastTimer = setTimeout(function() {
    toast.classList.remove('toast-visible');
    _toastTimer = null;
  }, 1500);
}
