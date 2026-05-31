/**
 * Utility functions for the Image Resizer app.
 */

const Utils = {
  /**
   * Format bytes into a human-readable string.
   * @param {number} bytes
   * @returns {string} e.g. "4.2 MB", "380 KB"
   */
  formatFileSize(bytes) {
    if (bytes === 0) return '0 B';
    const units = ['B', 'KB', 'MB', 'GB'];
    const k = 1024;
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    const value = bytes / Math.pow(k, i);
    return value.toFixed(value >= 100 ? 0 : value >= 10 ? 1 : 1) + ' ' + units[i];
  },

  /**
   * Sanitize a filename for safe download.
   * Strips path traversal, control characters, null bytes.
   * @param {string} name
   * @returns {string}
   */
  sanitizeFilename(name) {
    return name
      .replace(/\.\./g, '')           // strip path traversal
      .replace(/[\/\\]/g, '')         // strip slashes
      .replace(/[\x00-\x1f\x7f]/g, '') // strip control chars
      .replace(/[<>:"|?*]/g, '')      // strip reserved chars
      .trim() || 'image';
  },

  /**
   * Get file extension from MIME type.
   * @param {string} format - 'jpeg', 'png', 'webp'
   * @returns {string}
   */
  formatToExtension(format) {
    const map = { jpeg: 'jpg', png: 'png', webp: 'webp' };
    return map[format] || 'jpg';
  },

  /**
   * Get MIME type string from format key.
   * @param {string} format - 'jpeg', 'png', 'webp'
   * @returns {string}
   */
  formatToMime(format) {
    return 'image/' + format;
  },

  /**
   * Detect format from a file's MIME type.
   * @param {string} mime
   * @returns {string} 'jpeg', 'png', or 'webp'
   */
  mimeToFormat(mime) {
    if (mime === 'image/png') return 'png';
    if (mime === 'image/webp') return 'webp';
    return 'jpeg';
  },

  /**
   * Get the base name of a file (without extension).
   * @param {string} filename
   * @returns {string}
   */
  getBaseName(filename) {
    const lastDot = filename.lastIndexOf('.');
    return lastDot > 0 ? filename.substring(0, lastDot) : filename;
  },

  /**
   * Build a sanitized download filename.
   * @param {string} originalName
   * @param {string} format - 'jpeg', 'png', 'webp'
   * @returns {string}
   */
  buildDownloadFilename(originalName, format) {
    const base = Utils.sanitizeFilename(Utils.getBaseName(originalName));
    const ext = Utils.formatToExtension(format);
    return base + '-resized.' + ext;
  },

  /**
   * Calculate size change percentage.
   * @param {number} original
   * @param {number} resized
   * @returns {{ percent: number, label: string, smaller: boolean }}
   */
  calcSizeChange(original, resized) {
    if (original === 0) return { percent: 0, label: '0%', smaller: true };
    const ratio = ((original - resized) / original) * 100;
    const smaller = ratio >= 0;
    const absRatio = Math.abs(ratio);
    const label = smaller
      ? Math.round(absRatio) + '% smaller'
      : Math.round(absRatio) + '% larger';
    return { percent: ratio, label, smaller };
  },

  /**
   * Debounce a function call.
   * @param {Function} fn
   * @param {number} delay - ms
   * @returns {Function}
   */
  debounce(fn, delay) {
    let timer;
    return function (...args) {
      clearTimeout(timer);
      timer = setTimeout(() => fn.apply(this, args), delay);
    };
  },

  /**
   * Validate if a file is an accepted image type.
   * @param {File} file
   * @returns {{ valid: boolean, error: string }}
   */
  validateFile(file) {
    const allowedTypes = ['image/jpeg', 'image/png', 'image/webp'];
    const maxSize = 50 * 1024 * 1024; // 50MB

    if (!allowedTypes.includes(file.type)) {
      return {
        valid: false,
        error: "That doesn't look like an image. Try JPG, PNG, or WebP."
      };
    }

    if (file.size > maxSize) {
      const sizeMB = (file.size / (1024 * 1024)).toFixed(1);
      return {
        valid: false,
        error: 'That file is too heavy (' + sizeMB + ' MB). Max size is 50MB.'
      };
    }

    return { valid: true, error: '' };
  }
};
