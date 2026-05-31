/**
 * Image Resizer Engine — Canvas API based
 * Supports JPEG, PNG, WebP with quality control
 * Uses multi-step downscaling for reductions > 50%
 */

const Resizer = {
  /**
   * Resize an image to target dimensions
   * @param {HTMLImageElement} image - Source image element
   * @param {Object} options
   * @param {number} options.width - Target width in pixels
   * @param {number} options.height - Target height in pixels
   * @param {string} options.format - 'jpeg', 'png', or 'webp'
   * @param {number} options.quality - 0.0 to 1.0 (ignored for PNG)
   * @returns {Promise<Blob>}
   */
  async resize(image, { width, height, format = 'jpeg', quality = 0.85 }) {
    // Validate inputs
    width = Math.max(1, Math.round(width));
    height = Math.max(1, Math.round(height));
    quality = Math.max(0.01, Math.min(1.0, quality));

    if (!image || !image.naturalWidth || !image.naturalHeight) {
      throw new Error('Invalid image source');
    }

    const mime = 'image/' + format;

    // Check if multi-step downscaling is needed
    const scaleX = width / image.naturalWidth;
    const scaleY = height / image.naturalHeight;
    const needsMultiStep = scaleX < 0.5 || scaleY < 0.5;

    let blob;
    try {
      if (needsMultiStep) {
        blob = await this._multiStepResize(image, width, height, mime, quality);
      } else {
        blob = await this._singleStepResize(image, width, height, mime, quality);
      }
    } catch (err) {
      throw new Error('Failed to process image: ' + err.message);
    }

    return blob;
  },

  /**
   * Single-step resize for moderate changes
   */
  async _singleStepResize(image, width, height, mime, quality) {
    const canvas = document.createElement('canvas');
    canvas.width = width;
    canvas.height = height;

    const ctx = canvas.getContext('2d');
    ctx.imageSmoothingEnabled = true;
    ctx.imageSmoothingQuality = 'high';
    ctx.drawImage(image, 0, 0, width, height);

    return this._canvasToBlob(canvas, mime, quality);
  },

  /**
   * Multi-step downscaling — halves dimensions iteratively
   * until within 2x of target, then does final resize.
   * Produces much better quality for large reductions.
   */
  async _multiStepResize(image, targetWidth, targetHeight, mime, quality) {
    let currentWidth = image.naturalWidth;
    let currentHeight = image.naturalHeight;

    // Start with a canvas of the original size
    let srcCanvas = document.createElement('canvas');
    srcCanvas.width = currentWidth;
    srcCanvas.height = currentHeight;
    let srcCtx = srcCanvas.getContext('2d');
    srcCtx.imageSmoothingEnabled = true;
    srcCtx.imageSmoothingQuality = 'high';
    srcCtx.drawImage(image, 0, 0);

    // Iteratively halve until within 2x of target
    while (currentWidth / 2 > targetWidth && currentHeight / 2 > targetHeight) {
      const halfW = Math.round(currentWidth / 2);
      const halfH = Math.round(currentHeight / 2);

      const stepCanvas = document.createElement('canvas');
      stepCanvas.width = halfW;
      stepCanvas.height = halfH;
      const stepCtx = stepCanvas.getContext('2d');
      stepCtx.imageSmoothingEnabled = true;
      stepCtx.imageSmoothingQuality = 'high';
      stepCtx.drawImage(srcCanvas, 0, 0, halfW, halfH);

      srcCanvas = stepCanvas;
      currentWidth = halfW;
      currentHeight = halfH;
    }

    // Final step to exact target dimensions
    const finalCanvas = document.createElement('canvas');
    finalCanvas.width = targetWidth;
    finalCanvas.height = targetHeight;
    const finalCtx = finalCanvas.getContext('2d');
    finalCtx.imageSmoothingEnabled = true;
    finalCtx.imageSmoothingQuality = 'high';
    finalCtx.drawImage(srcCanvas, 0, 0, targetWidth, targetHeight);

    return this._canvasToBlob(finalCanvas, mime, quality);
  },

  /**
   * Convert canvas to blob
   * @param {HTMLCanvasElement} canvas
   * @param {string} mime
   * @param {number} quality
   * @returns {Promise<Blob>}
   */
  _canvasToBlob(canvas, mime, quality) {
    return new Promise((resolve, reject) => {
      // PNG ignores quality parameter
      const args = mime === 'image/png' ? [mime] : [mime, quality];
      canvas.toBlob(
        (blob) => {
          if (blob) {
            resolve(blob);
          } else {
            reject(new Error('Canvas toBlob returned null'));
          }
        },
        ...args
      );
    });
  }
};
