import type { CategoryId } from './types';
import { getUnit } from './units';
import { formatResult } from '../utils/format';

function convertTemperature(value: number, fromId: string, toId: string): number {
  if (fromId === toId) return value;

  // Convert to Celsius first
  let celsius: number;
  switch (fromId) {
    case 'celsius':
      celsius = value;
      break;
    case 'fahrenheit':
      celsius = (value - 32) * 5 / 9;
      break;
    case 'kelvin':
      celsius = value - 273.15;
      break;
    default:
      return NaN;
  }

  // Convert from Celsius to target
  switch (toId) {
    case 'celsius':
      return celsius;
    case 'fahrenheit':
      return celsius * 9 / 5 + 32;
    case 'kelvin':
      return celsius + 273.15;
    default:
      return NaN;
  }
}

export function convert(
  categoryId: CategoryId,
  fromUnitId: string,
  toUnitId: string,
  value: number
): number {
  if (!Number.isFinite(value)) return NaN;
  if (fromUnitId === toUnitId) return value;

  if (categoryId === 'temperature') {
    return convertTemperature(value, fromUnitId, toUnitId);
  }

  const fromUnit = getUnit(categoryId, fromUnitId);
  const toUnit = getUnit(categoryId, toUnitId);

  if (fromUnit.toBase === null || toUnit.toBase === null) return NaN;

  return (value * fromUnit.toBase) / toUnit.toBase;
}

export function getQuickReference(
  categoryId: CategoryId,
  fromUnitId: string,
  toUnitId: string
): string {
  const fromUnit = getUnit(categoryId, fromUnitId);
  const toUnit = getUnit(categoryId, toUnitId);

  if (categoryId === 'temperature') {
    if (fromUnitId === 'celsius' && toUnitId === 'fahrenheit') return '0°C = 32°F';
    if (fromUnitId === 'fahrenheit' && toUnitId === 'celsius') return '32°F = 0°C';
    if (fromUnitId === 'celsius' && toUnitId === 'kelvin') return '0°C = 273.15 K';
    if (fromUnitId === 'kelvin' && toUnitId === 'celsius') return '273.15 K = 0°C';
    if (fromUnitId === 'fahrenheit' && toUnitId === 'kelvin') return '32°F = 273.15 K';
    if (fromUnitId === 'kelvin' && toUnitId === 'fahrenheit') return '273.15 K = 32°F';
    return '';
  }

  const result = convert(categoryId, fromUnitId, toUnitId, 1);
  const formatted = formatResult(result);
  return `1 ${fromUnit.symbol} = ${formatted} ${toUnit.symbol}`;
}