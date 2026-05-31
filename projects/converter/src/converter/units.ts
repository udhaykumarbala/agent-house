import type { UnitCategory, CategoryId } from './types';

export const categories: UnitCategory[] = [
  {
    id: 'length',
    label: 'Length',
    units: [
      { id: 'mm', label: 'Millimeter', symbol: 'mm', toBase: 0.001 },
      { id: 'cm', label: 'Centimeter', symbol: 'cm', toBase: 0.01 },
      { id: 'm', label: 'Meter', symbol: 'm', toBase: 1 },
      { id: 'km', label: 'Kilometer', symbol: 'km', toBase: 1000 },
      { id: 'in', label: 'Inch', symbol: 'in', toBase: 0.0254 },
      { id: 'ft', label: 'Foot', symbol: 'ft', toBase: 0.3048 },
      { id: 'yd', label: 'Yard', symbol: 'yd', toBase: 0.9144 },
      { id: 'mi', label: 'Mile', symbol: 'mi', toBase: 1609.344 },
    ],
    defaultFrom: 'km',
    defaultTo: 'mi',
  },
  {
    id: 'weight',
    label: 'Weight',
    units: [
      { id: 'mg', label: 'Milligram', symbol: 'mg', toBase: 0.001 },
      { id: 'g', label: 'Gram', symbol: 'g', toBase: 1 },
      { id: 'kg', label: 'Kilogram', symbol: 'kg', toBase: 1000 },
      { id: 'oz', label: 'Ounce', symbol: 'oz', toBase: 28.3495 },
      { id: 'lb', label: 'Pound', symbol: 'lb', toBase: 453.592 },
      { id: 't', label: 'Metric Ton', symbol: 't', toBase: 1000000 },
    ],
    defaultFrom: 'kg',
    defaultTo: 'lb',
  },
  {
    id: 'temperature',
    label: 'Temperature',
    units: [
      { id: 'celsius', label: 'Celsius', symbol: '°C', toBase: null },
      { id: 'fahrenheit', label: 'Fahrenheit', symbol: '°F', toBase: null },
      { id: 'kelvin', label: 'Kelvin', symbol: 'K', toBase: null },
    ],
    defaultFrom: 'celsius',
    defaultTo: 'fahrenheit',
  },
];

export function getCategory(id: CategoryId): UnitCategory {
  return categories.find((c) => c.id === id)!;
}

export function getUnit(categoryId: CategoryId, unitId: string) {
  const cat = getCategory(categoryId);
  return cat.units.find((u) => u.id === unitId)!;
}
