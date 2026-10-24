'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { api, getToken } from '@/lib/api';
import Sidebar from '@/components/Sidebar';
import { Server } from 'lucide-react';

type Device = { id: number; name: string; type: string; ip: string; location: string; tags: string };

export default function DevicesPage() {
  const router = useRouter();
  const [devices, setDevices] = useState<Device[]>([]);
  const [form, setForm] = useState({ name: '', type: '', ip: '', location: '', tags: '' });
  const [err, setErr] = useState('');

  async function load() {
    try { setDevices(await api<Device[]>('/devices')); } catch { router.replace('/login'); }
  }

  useEffect(() => { if (!getToken()) return router.replace('/login'); load(); }, [router]);

  async function add(e: React.FormEvent) {
    e.preventDefault();
    setErr('');
    try {
      await api('/devices', { method: 'POST', body: JSON.stringify(form) });
      setForm({ name: '', type: '', ip: '', location: '', tags: '' });
      load();
    } catch (e: any) { setErr(e.message); }
  }

  async function del(id: number) {
    if (!confirm('Delete this device and its services?')) return;
    await api('/devices/' + id, { method: 'DELETE' });
    load();
  }

  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <main className="flex-1 p-4 md:p-8 pt-20 md:pt-8 max-w-6xl min-w-0">
        <div className="animate-fade-in">
          <h1 className="text-2xl md:text-3xl font-semibold tracking-tight">Devices</h1>
          <p className="mt-1 text-sm" style={{ color: 'var(--label-secondary)' }}>All hardware in your infrastructure</p>
        </div>
        <form onSubmit={add} className="card mt-6 grid grid-cols-1 sm:grid-cols-2 md:grid-cols-6 gap-2 items-stretch">
          <input className="input md:col-span-2" placeholder="Name" value={form.name} onChange={e => setForm({...form, name: e.target.value})} required />
          <input className="input" placeholder="Type" value={form.type} onChange={e => setForm({...form, type: e.target.value})} />
          <input className="input" placeholder="IP" value={form.ip} onChange={e => setForm({...form, ip: e.target.value})} />
          <input className="input" placeholder="Location" value={form.location} onChange={e => setForm({...form, location: e.target.value})} />
          <button className="btn btn-primary py-3">Add Device</button>
        </form>
        {err && <div className="mt-2 text-sm" style={{ color: '#FF453A' }}>{err}</div>}

        <div className="mt-6 space-y-2">
          {devices.map((d, i) => (
            <div key={d.id} className="card flex flex-col sm:flex-row sm:justify-between sm:items-center gap-3 animate-fade-in" style={{ animationDelay: i * 30 + 'ms' }}>
              <div className="flex items-center gap-4 min-w-0">
                <div className="w-10 h-10 rounded-xl flex items-center justify-center shrink-0" style={{ background: 'rgba(10,132,255,0.15)', color: '#0A84FF' }}>
                  <Server size={18} strokeWidth={2.2} />
                </div>
                <div className="min-w-0">
                  <div className="font-semibold text-[15px] truncate">{d.name}</div>
                  <div className="text-[13px] truncate" style={{ color: 'var(--label-secondary)' }}>{[d.type, d.ip, d.location].filter(Boolean).join(' - ') || 'No details'}</div>
                </div>
              </div>
              <button onClick={() => del(d.id)} className="btn btn-danger self-end sm:self-auto">Delete</button>
            </div>
          ))}
          {devices.length === 0 && <div className="card text-center py-12" style={{ color: 'var(--label-secondary)' }}>No devices yet. Add your first device above.</div>}
        </div>
      </main>
    </div>
  );
}