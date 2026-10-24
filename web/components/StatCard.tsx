import type { LucideIcon } from 'lucide-react';

type Props = {
  label: string;
  value: number | string;
  accent?: 'blue' | 'green' | 'red' | 'purple';
  Icon?: LucideIcon;
};

const ACCENTS = {
  blue:   { color: '#0A84FF', glow: 'rgba(10,132,255,0.15)' },
  green:  { color: '#30D158', glow: 'rgba(48,209,88,0.15)' },
  red:    { color: '#FF453A', glow: 'rgba(255,69,58,0.15)' },
  purple: { color: '#BF5AF2', glow: 'rgba(191,90,242,0.15)' },
};

export default function StatCard({ label, value, accent = 'blue', Icon }: Props) {
  const a = ACCENTS[accent];
  return (
    <div className="card animate-fade-in">
      <div className="flex items-start justify-between">
        <div className="text-[13px] font-medium" style={{ color: 'var(--label-secondary)' }}>{label}</div>
        {Icon && (
          <div className="w-8 h-8 rounded-lg flex items-center justify-center"
            style={{ background: a.glow, color: a.color }}>
            <Icon size={16} strokeWidth={2.2} />
          </div>
        )}
      </div>
      <div className="mt-3 text-3xl font-semibold tracking-tight" style={{ color: a.color }}>{value}</div>
    </div>
  );
}