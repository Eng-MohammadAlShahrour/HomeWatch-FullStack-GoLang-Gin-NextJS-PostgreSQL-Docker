'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { api, getToken } from '@/lib/api';
import Sidebar from '@/components/Sidebar';
import StatusDot from '@/components/StatusDot';
import UptimeChart from '@/components/UptimeChart';
import { ChevronDown, ChevronRight } from 'lucide-react';

type Service = {
  id: number; name: string; url: string; device_id: number;
  check_type: string; check_interval: number; status: string;
  last_checked_at: string | null; response_ms: number;
};

type Summary = {
  service_id: number;
  name: string;
  total: number;
  up: number;
  uptime_pct: number;
  avg_ms: number;
};

export default function ServicesPage() {
  const router = useRouter();
  const [services, setServices] = useState<Service[]>([]);
  const [summaries, setSummaries] = useState<Record<number, Summary>>({});
  const [expanded, setExpanded] = useState<number | null>(null);
  const [form, setForm] = useState({ name: '', url: '', device_id: 0, check_interval: 30 });
  const [err, setErr] = useState('');

  async function load() {
    try {
      const [s, sum] = await Promise.all([
        api<Service[]>('/services'),
        api<Summary[]>('/uptime/summary'),
      ]);
      setServices(s);
      const map: Record<number, Summary> = {};
      sum.forEach(x => { map[x.service_id] = x; });
      setSummaries(map);
    } catch { router.replace('/login'); }
  }

  useEffect(() => {
    if (!getToken()) return router.replace('/login');
    load();
    const t = setInterval(load, 5000);
    return () => clearInterval(t);
  }, [router]);

  async function add(e: React.FormEvent) {
    e.preventDefault();
    setErr('');
    try {
      await api('/services', { method: 'POST', body: JSON.stringify({ ...form, device_id: Number(form.device_id), check_type: 'http' }) });
      setForm({ name: '', url: '', device_id: 0, check_interval: 30 });
      load();
    } catch (e: any) { setErr(e.message); }
  }

  async function del(id: number) {
    if (!confirm('Delete this service?')) return;
    await api('/services/' + id, { method: 'DELETE' });
    load();
  }

  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <main className="flex-1 p-4 md:p-8 pt-20 md:pt-8 max-w-6xl min-w-0">
        <div className="animate-fade-in">
          <h1 className="text-2xl md:text-3xl font-semibold tracking-tight">Services</h1>
          <p className="mt-1 text-sm" style={{ color: 'var(--label-secondary)' }}>Monitored endpoints · tap to view uptime</p>
        </div>

        <form onSubmit={add} className="card mt-6 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-5 gap-2 items-stretch">
          <input className="input" placeholder="Name" value={form.name} onChange={e => setForm({...form, name: e.target.value})} required />
          <input className="input md:col-span-2" placeholder="https://example.com" value={form.url} onChange={e => setForm({...form, url: e.target.value})} />
          <input className="input" placeholder="Device ID" type="number" value={form.device_id || ''} onChange={e => setForm({...form, device_id: Number(e.target.value)})} />
          <button className="btn btn-primary py-3">Add Service</button>
        </form>
        {err && <div className="mt-2 text-sm" style={{ color: '#FF453A' }}>{err}</div>}

        <div className="mt-6 space-y-2">
          {services.map((s, i) => {
            const sum = summaries[s.id];
            const uptimeColor = sum ? (sum.uptime_pct >= 99 ? '#30D158' : sum.uptime_pct >= 95 ? '#FF9F0A' : '#FF453A') : '#8E8E93';
            const isOpen = expanded === s.id;
            return (
              <div key={s.id} className="animate-fade-in" style={{ animationDelay: i * 30 + 'ms' }}>
                <div className="card">
                  <div className="flex items-start sm:items-center gap-3">
                    <button onClick={() => setExpanded(isOpen ? null : s.id)} className="text-slate-400 hover:text-white transition-colors mt-2 sm:mt-0 shrink-0">
                      {isOpen ? <ChevronDown size={18} /> : <ChevronRight size={18} />}
                    </button>
                    <div className="w-10 h-10 rounded-xl flex items-center justify-center shrink-0"
                      style={{ background: s.status === 'up' ? 'rgba(48,209,88,0.15)' : s.status === 'down' ? 'rgba(255,69,58,0.15)' : 'rgba(142,142,147,0.15)' }}>
                      <StatusDot status={s.status} />
                    </div>
                    <div className="flex-1 min-w-0">
                      <div className="font-semibold text-[15px] flex items-center gap-2 flex-wrap">
                        <span className="truncate">{s.name}</span>
                        <span className="text-[10px] md:text-[11px] px-2 py-0.5 rounded-full uppercase tracking-wider font-medium shrink-0"
                          style={{ background: s.status === 'up' ? 'rgba(48,209,88,0.15)' : s.status === 'down' ? 'rgba(255,69,58,0.15)' : 'rgba(142,142,147,0.15)', color: s.status === 'up' ? '#30D158' : s.status === 'down' ? '#FF453A' : '#8E8E93' }}>{s.status}</span>
                      </div>
                      <div className="text-[11px] md:text-[13px] font-mono truncate" style={{ color: 'var(--label-secondary)' }}>{s.url} {s.response_ms > 0 && '- ' + s.response_ms + 'ms'}</div>
                    </div>
                    <div className="flex flex-col sm:flex-row items-end sm:items-center gap-2 shrink-0">
                      {sum && sum.total > 0 && (
                        <div className="text-right">
                          <div className="text-sm md:text-lg font-semibold" style={{ color: uptimeColor }}>
                            {sum.uptime_pct.toFixed(1)}%
                          </div>
                          <div className="text-[9px] md:text-[10px] uppercase tracking-wider" style={{ color: 'var(--label-secondary)' }}>uptime</div>
                        </div>
                      )}
                      <button onClick={() => del(s.id)} className="btn btn-danger text-xs">Delete</button>
                    </div>
                  </div>
                </div>
                {isOpen && <UptimeChart serviceId={s.id} serviceName={s.name} />}
              </div>
            );
          })}
          {services.length === 0 && (
            <div className="card text-center py-12" style={{ color: 'var(--label-secondary)' }}>
              No services yet. Add an endpoint to monitor.
            </div>
          )}
        </div>
      </main>
    </div>
  );
}