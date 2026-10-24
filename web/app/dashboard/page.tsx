'use client';
import { useEffect, useState } from 'react';
import { useRouter } from 'next/navigation';
import { api, getToken } from '@/lib/api';
import Sidebar from '@/components/Sidebar';
import StatCard from '@/components/StatCard';
import { Server, Activity, CircleCheck, CircleX } from 'lucide-react';

type Summary = { devices: number; services: number; up: number; down: number };

export default function Dashboard() {
  const router = useRouter();
  const [data, setData] = useState<Summary | null>(null);

  useEffect(() => {
    if (!getToken()) return router.replace('/login');
    const load = () => api<Summary>('/dashboard').then(setData).catch(() => router.replace('/login'));
    load();
    const t = setInterval(load, 5000);
    return () => clearInterval(t);
  }, [router]);

  return (
    <div className="flex min-h-screen">
      <Sidebar />
      <main className="flex-1 p-4 md:p-8 pt-20 md:pt-8 max-w-6xl min-w-0">
        <div className="animate-fade-in">
          <h1 className="text-2xl md:text-3xl font-semibold tracking-tight">Dashboard</h1>
          <p className="mt-1 text-sm" style={{ color: 'var(--label-secondary)' }}>Live overview of your infrastructure</p>
        </div>
        <div className="grid grid-cols-2 md:grid-cols-4 gap-3 md:gap-4 mt-6 md:mt-8">
          <StatCard label="Devices"  value={data?.devices ?? 0}  accent="blue"   Icon={Server} />
          <StatCard label="Services" value={data?.services ?? 0} accent="purple" Icon={Activity} />
          <StatCard label="Online"   value={data?.up ?? 0}       accent="green"  Icon={CircleCheck} />
          <StatCard label="Offline"  value={data?.down ?? 0}     accent="red"    Icon={CircleX} />
        </div>
        <div className="card mt-6 md:mt-8 animate-fade-in">
          <div className="flex items-center gap-2">
            <div className="w-2 h-2 rounded-full animate-pulse" style={{ background: '#30D158', boxShadow: '0 0 12px #30D158' }} />
            <span className="text-sm font-medium">Auto-refresh every 5s</span>
          </div>
          <p className="mt-2 text-sm" style={{ color: 'var(--label-secondary)' }}>Health checks run every 15 seconds in the background.</p>
        </div>
      </main>
    </div>
  );
}