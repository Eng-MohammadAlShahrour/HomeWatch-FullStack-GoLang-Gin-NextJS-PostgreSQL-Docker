'use client';
import { useEffect, useState } from 'react';
import { api, getToken } from '@/lib/api';
import { AreaChart, Area, XAxis, YAxis, Tooltip, ResponsiveContainer, CartesianGrid } from 'recharts';

type Point = { time: string; avg_ms: number; up_count: number; down_count: number };
type Data = {
  service_id: number;
  hours: number;
  total: number;
  up: number;
  down: number;
  uptime_pct: number;
  series: Point[];
};

export default function UptimeChart({ serviceId, serviceName }: { serviceId: number; serviceName: string }) {
  const [data, setData] = useState<Data | null>(null);
  const [hours, setHours] = useState(24);

  useEffect(() => {
    if (!getToken()) return;
    api<Data>(`/services/${serviceId}/uptime?hours=${hours}`)
      .then(setData)
      .catch(console.error);
  }, [serviceId, hours]);

  const chartData = (data?.series ?? []).map(p => ({
    time: new Date(p.time).toLocaleTimeString('en-US', { hour: '2-digit', minute: '2-digit' }),
    ms: Math.round(p.avg_ms),
  }));

  const uptimeColor = (data?.uptime_pct ?? 0) >= 99 ? '#30D158' : (data?.uptime_pct ?? 0) >= 95 ? '#FF9F0A' : '#FF453A';

  return (
    <div className="card mt-2 animate-fade-in">
      <div className="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 mb-4">
        <div>
          <h3 className="font-semibold text-[14px] md:text-[15px]">{serviceName} — Uptime</h3>
          <p className="text-[12px] md:text-[13px]" style={{ color: 'var(--label-secondary)' }}>
            Last {hours}h · {data?.total ?? 0} checks
          </p>
        </div>
        <div className="flex items-center justify-between sm:justify-end gap-3">
          <div className="text-right">
            <div className="text-xl md:text-2xl font-semibold" style={{ color: uptimeColor }}>
              {(data?.uptime_pct ?? 0).toFixed(2)}%
            </div>
            <div className="text-[10px] uppercase tracking-wide" style={{ color: 'var(--label-secondary)' }}>Uptime</div>
          </div>
          <select value={hours} onChange={e => setHours(Number(e.target.value))} className="input !w-auto !py-2 text-sm">
            <option value={1}>1h</option>
            <option value={6}>6h</option>
            <option value={24}>24h</option>
            <option value={168}>7d</option>
          </select>
        </div>
      </div>

      <div style={{ width: '100%', height: 200 }}>
        <ResponsiveContainer>
          <AreaChart data={chartData} margin={{ top: 5, right: 5, left: -25, bottom: 0 }}>
            <defs>
              <linearGradient id={`g-${serviceId}`} x1="0" y1="0" x2="0" y2="1">
                <stop offset="0%" stopColor="#0A84FF" stopOpacity={0.4} />
                <stop offset="100%" stopColor="#0A84FF" stopOpacity={0} />
              </linearGradient>
            </defs>
            <CartesianGrid stroke="rgba(255,255,255,0.05)" vertical={false} />
            <XAxis dataKey="time" stroke="rgba(235,235,245,0.4)" tick={{ fontSize: 10 }} />
            <YAxis stroke="rgba(235,235,245,0.4)" tick={{ fontSize: 10 }} unit="ms" />
            <Tooltip
              contentStyle={{ background: 'rgba(28,28,30,0.95)', border: '1px solid rgba(255,255,255,0.1)', borderRadius: 12, fontSize: 12 }}
              labelStyle={{ color: 'rgba(235,235,245,0.6)' }}
            />
            <Area type="monotone" dataKey="ms" stroke="#0A84FF" strokeWidth={2} fill={`url(#g-${serviceId})`} />
          </AreaChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}