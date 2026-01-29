'use client';

import * as React from 'react';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { Select } from '@/components/ui/select';
import { POPULAR_TEMPLATES } from '@/lib/constants';
import {
  BILLING_CYCLES,
  CATEGORIES,
  type SubscriptionInput,
  type SubscriptionTemplate,
} from '@/types/subscription';

interface SubscriptionFormProps {
  initialData?: Partial<SubscriptionInput>;
  onSubmit: (data: SubscriptionInput) => void;
  onCancel: () => void;
  isLoading?: boolean;
  submitLabel?: string;
}

export function SubscriptionForm({
  initialData,
  onSubmit,
  onCancel,
  isLoading,
  submitLabel = 'Add Subscription',
}: SubscriptionFormProps) {
  const [formData, setFormData] = React.useState<Partial<SubscriptionInput>>({
    name: '',
    price: 0,
    billingCycle: 'monthly',
    nextBillingDate: new Date().toISOString().split('T')[0],
    category: 'other',
    ...initialData,
  });

  const [searchQuery, setSearchQuery] = React.useState('');

  const filteredTemplates = POPULAR_TEMPLATES.filter((t) =>
    t.name.toLowerCase().includes(searchQuery.toLowerCase())
  ).slice(0, 8);

  const handleTemplateSelect = (template: SubscriptionTemplate) => {
    setFormData({
      ...formData,
      name: template.name,
      price: template.defaultPrice,
      billingCycle: template.billingCycle,
      category: template.category,
      icon: template.icon,
      color: template.color,
    });
    setSearchQuery('');
  };

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (!formData.name || !formData.price || !formData.nextBillingDate) return;

    onSubmit({
      name: formData.name,
      price: formData.price,
      billingCycle: formData.billingCycle || 'monthly',
      nextBillingDate: formData.nextBillingDate,
      category: formData.category,
      icon: formData.icon,
      color: formData.color,
      description: formData.description,
      notes: formData.notes,
      reminderDays: formData.reminderDays,
    });
  };

  return (
    <form onSubmit={handleSubmit} className="space-y-4">
      {/* Template search */}
      {!initialData?.name && (
        <div>
          <Input
            placeholder="Search popular services..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
          />
          {searchQuery && filteredTemplates.length > 0 && (
            <div className="mt-2 grid grid-cols-4 gap-2">
              {filteredTemplates.map((template) => (
                <button
                  key={template.name}
                  type="button"
                  onClick={() => handleTemplateSelect(template)}
                  className="flex flex-col items-center gap-1 rounded-lg border border-border p-2 transition-colors hover:bg-surface-hover dark:border-dark-border dark:hover:bg-dark-surface-hover"
                >
                  <span className="text-2xl">{template.icon}</span>
                  <span className="text-xs text-text-secondary dark:text-dark-text-secondary truncate w-full text-center">
                    {template.name}
                  </span>
                </button>
              ))}
            </div>
          )}

          {!searchQuery && (
            <>
              <p className="mt-4 mb-2 text-sm text-text-secondary dark:text-dark-text-secondary">
                Popular:
              </p>
              <div className="flex gap-2 overflow-x-auto pb-2">
                {POPULAR_TEMPLATES.slice(0, 6).map((template) => (
                  <button
                    key={template.name}
                    type="button"
                    onClick={() => handleTemplateSelect(template)}
                    className="flex flex-shrink-0 flex-col items-center gap-1 rounded-lg border border-border p-3 transition-colors hover:bg-surface-hover dark:border-dark-border dark:hover:bg-dark-surface-hover"
                  >
                    <span className="text-2xl">{template.icon}</span>
                    <span className="text-xs text-text-secondary dark:text-dark-text-secondary">
                      {template.name}
                    </span>
                  </button>
                ))}
              </div>
            </>
          )}

          <div className="my-4 flex items-center gap-3 text-sm text-text-tertiary dark:text-dark-text-tertiary">
            <div className="h-px flex-1 bg-border dark:bg-dark-border" />
            <span>or enter manually</span>
            <div className="h-px flex-1 bg-border dark:bg-dark-border" />
          </div>
        </div>
      )}

      {/* Form fields */}
      <Input
        label="Service Name"
        placeholder="e.g., Netflix"
        value={formData.name || ''}
        onChange={(e) => setFormData({ ...formData, name: e.target.value })}
        required
      />

      <div className="grid grid-cols-2 gap-3">
        <Input
          label="Price"
          type="number"
          step="0.01"
          min="0"
          placeholder="0.00"
          value={formData.price || ''}
          onChange={(e) =>
            setFormData({ ...formData, price: parseFloat(e.target.value) || 0 })
          }
          required
        />
        <Select
          label="Billing Cycle"
          options={BILLING_CYCLES.map((c) => ({ value: c.value, label: c.label }))}
          value={formData.billingCycle || 'monthly'}
          onChange={(e) =>
            setFormData({
              ...formData,
              billingCycle: e.target.value as SubscriptionInput['billingCycle'],
            })
          }
        />
      </div>

      <Input
        label="Next Billing Date"
        type="date"
        value={formData.nextBillingDate || ''}
        onChange={(e) =>
          setFormData({ ...formData, nextBillingDate: e.target.value })
        }
        required
      />

      <Select
        label="Category"
        options={CATEGORIES.map((c) => ({
          value: c.value,
          label: `${c.icon} ${c.label}`,
        }))}
        value={formData.category || 'other'}
        onChange={(e) =>
          setFormData({
            ...formData,
            category: e.target.value as SubscriptionInput['category'],
          })
        }
      />

      {/* Actions */}
      <div className="flex gap-3 pt-4">
        <Button
          type="button"
          variant="secondary"
          className="flex-1"
          onClick={onCancel}
        >
          Cancel
        </Button>
        <Button type="submit" className="flex-1" disabled={isLoading}>
          {isLoading ? 'Saving...' : submitLabel}
        </Button>
      </div>
    </form>
  );
}
