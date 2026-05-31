/**
 * Color generation engine
 * Pure functions — no DOM access
 */

/**
 * Generate a random HSL color within constrained ranges
 * S: 45-90%, L: 35-65% for vibrant, readable colors
 */
function randomHSL() {
  var h = Math.floor(Math.random() * 360);
  var s = Math.floor(Math.random() * 46) + 45;
  var l = Math.floor(Math.random() * 31) + 35;
  return { h: h, s: s, l: l };
}

/**
 * Generate 5 hue values based on harmony mode
 */
function harmonize(baseHue, mode) {
  var hues = [];
  switch (mode) {
    case 'analogous':
      hues.push(
        baseHue,
        (baseHue - 30 + 360) % 360,
        (baseHue - 15 + 360) % 360,
        (baseHue + 15) % 360,
        (baseHue + 30) % 360
      );
      break;
    case 'complementary':
      hues.push(
        baseHue,
        (baseHue + 180) % 360,
        (baseHue + 15) % 360,
        (baseHue + 180 + 15) % 360,
        (baseHue - 15 + 360) % 360
      );
      break;
    case 'triadic':
      hues.push(
        baseHue,
        (baseHue + 120) % 360,
        (baseHue + 240) % 360,
        (baseHue + 135) % 360,
        (baseHue + 225) % 360
      );
      break;
    case 'split-complementary':
      hues.push(
        baseHue,
        (baseHue + 150) % 360,
        (baseHue + 210) % 360,
        (baseHue + 165) % 360,
        (baseHue + 195) % 360
      );
      break;
    case 'monochromatic':
      hues.push(baseHue, baseHue, baseHue, baseHue, baseHue);
      break;
    default:
      for (var i = 0; i < 5; i++) {
        hues.push(Math.floor(Math.random() * 360));
      }
      break;
  }
  return hues;
}

/**
 * Convert HSL to RGB
 * h: 0-360, s: 0-100, l: 0-100
 * Returns { r, g, b } with values 0-255
 */
function hslToRgb(h, s, l) {
  s /= 100;
  l /= 100;
  var c = (1 - Math.abs(2 * l - 1)) * s;
  var x = c * (1 - Math.abs((h / 60) % 2 - 1));
  var m = l - c / 2;
  var r, g, b;

  if (h < 60)       { r = c; g = x; b = 0; }
  else if (h < 120) { r = x; g = c; b = 0; }
  else if (h < 180) { r = 0; g = c; b = x; }
  else if (h < 240) { r = 0; g = x; b = c; }
  else if (h < 300) { r = x; g = 0; b = c; }
  else              { r = c; g = 0; b = x; }

  return {
    r: Math.round((r + m) * 255),
    g: Math.round((g + m) * 255),
    b: Math.round((b + m) * 255)
  };
}

/**
 * Convert RGB to hex string (uppercase, with #)
 */
function rgbToHex(r, g, b) {
  return '#' + [r, g, b]
    .map(function(v) { return v.toString(16).padStart(2, '0'); })
    .join('')
    .toUpperCase();
}

/**
 * Convert hex string to RGB object
 * Accepts with or without #, 3 or 6 chars
 */
function hexToRgb(hex) {
  var clean = hex.replace('#', '');
  if (clean.length === 3) {
    clean = clean[0] + clean[0] + clean[1] + clean[1] + clean[2] + clean[2];
  }
  return {
    r: parseInt(clean.substring(0, 2), 16),
    g: parseInt(clean.substring(2, 4), 16),
    b: parseInt(clean.substring(4, 6), 16)
  };
}

/**
 * Get contrast text color based on background luminance
 * Formula: 0.299*R + 0.587*G + 0.114*B, threshold 150
 * Returns #000000 for light backgrounds, #FFFFFF for dark
 */
function getContrastColor(hex) {
  var rgb = hexToRgb(hex);
  var luminance = 0.299 * rgb.r + 0.587 * rgb.g + 0.114 * rgb.b;
  return luminance > 150 ? '#000000' : '#FFFFFF';
}

/**
 * Strict hex validation (without #)
 */
function isValidHex(value) {
  return /^[0-9A-Fa-f]{3,6}$/.test(value);
}

/**
 * Format a color object as a display string for the given format
 * 'hex' → '#3A86FF', 'rgb' → 'rgb(58, 134, 255)', 'hsl' → 'hsl(220, 45%, 61%)'
 */
function formatColor(color, format) {
  switch (format) {
    case 'rgb':
      return 'rgb(' + color.rgb.r + ', ' + color.rgb.g + ', ' + color.rgb.b + ')';
    case 'hsl':
      return 'hsl(' + color.h + ', ' + color.s + '%, ' + color.l + '%)';
    default:
      return color.hex;
  }
}

/**
 * Create a color object from HSL values
 */
function createColorObject(h, s, l, locked) {
  var rgb = hslToRgb(h, s, l);
  var hex = rgbToHex(rgb.r, rgb.g, rgb.b);
  return { h: h, s: s, l: l, hex: hex, rgb: rgb, locked: locked || false };
}

/**
 * Generate a 5-color palette based on harmony mode
 * Preserves locked colors from existingPalette
 */
function generatePalette(mode, existingPalette) {
  mode = mode || 'random';
  existingPalette = existingPalette || [];

  var base = randomHSL();
  var hues = harmonize(base.h, mode);
  var palette = [];

  for (var i = 0; i < 5; i++) {
    if (existingPalette[i] && existingPalette[i].locked) {
      palette.push(existingPalette[i]);
    } else {
      var s, l;
      if (mode === 'monochromatic') {
        var satValues = [50, 60, 70, 80, 88];
        var lightValues = [36, 43, 50, 57, 64];
        s = satValues[i] + Math.floor(Math.random() * 6) - 3;
        l = lightValues[i] + Math.floor(Math.random() * 4) - 2;
        s = Math.min(90, Math.max(45, s));
        l = Math.min(65, Math.max(35, l));
      } else {
        s = Math.floor(Math.random() * 46) + 45;
        l = Math.floor(Math.random() * 31) + 35;
      }
      palette.push(createColorObject(hues[i], s, l, false));
    }
  }

  return palette;
}
