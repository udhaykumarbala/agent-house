export type CategoryId = 'length' | 'weight' | 'temperature';

export interface Unit {
  id: string;
  label: string;
  symbol: string;
  toBase: number | null; // null for temperature (formula-based)
}

export interface UnitCategory {
  id: CategoryId;
  label: string;
  units: Unit[];
  defaultFrom: string;
  defaultTo: string;
}

export interface ConversionState {
  category: CategoryId;
  fromUnit: string;
  toUnit: string;
  fromValue: string;
  toValue: string;
  inputDirection: 'from' | 'to';
}
