'use client';
import { useState } from 'react';
import Link from 'next/link';
import { usePathname, useRouter } from 'next/navigation';
import { clearToken } from '@/lib/api';
import { LayoutDashboard, Server, Activity, LogOut, Boxes, Menu, X } from 'lucide-react';

const NAV = [
  { href: '/dashboard', label: 'Dashboard', Icon: LayoutDashboard },
  { href: '/devices',   label: 'Devices',   Icon: Server },
  { href: '/services',  label: 'Services',  Icon: Activity },
];

function NavContent({ path, onNavigate, onSignOut }: { path: string; onNavigate: () => void; onSignOut: () => void }) {
  return (
    <>
      <div className="px-3 py-5 flex items-center gap-3">
        <div className="w-9 h-9 rounded-xl flex items-center justify-center"
          style={{ background: 'linear-gradient(135deg, #0A84FF, #BF5AF2)', boxShadow: '0 4px 16px rgba(10, 132, 255, 0.4)' }}>
          <Boxes size={20} color="white" strokeWidth={2.2} />
        </div>
        <div>
          <div className="font-semibold text-[15px] tracking-tight">HomeWatch</div>
          <div className="text-[11px]" style={{ color: 'var(--label-secondary)' }}>Monitor your home lab</div>
        </div>
      </div>

      <div className="px-3 mt-2 text-[11px] font-medium uppercase tracking-wider"
        style={{ color: 'var(--label-tertiary)' }}>
        Menu
      </div>

      <nav className="flex flex-col gap-1 mt-1">
        {NAV.map(({ href, label, Icon }) => {
          const active = path === href;
          return (
            <Link key={href} href={href} onClick={onNavigate}
              className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all"
              style={{
                background: active ? 'var(--sf-blue)' : 'transparent',
                color: active ? 'white' : 'var(--label-primary)',
                boxShadow: active ? '0 4px 16px rgba(10, 132, 255, 0.35)' : 'none',
              }}>
              <Icon size={18} strokeWidth={2} />
              {label}
            </Link>
          );
        })}
      </nav>

      <div className="flex-1" />

      <button onClick={onSignOut}
        className="flex items-center gap-3 px-3 py-2.5 rounded-xl text-sm font-medium transition-all"
        style={{ color: 'var(--sf-red)' }}>
        <LogOut size={18} strokeWidth={2} />
        Sign Out
      </button>
    </>
  );
}

export default function Sidebar() {
  const path = usePathname();
  const router = useRouter();
  const [open, setOpen] = useState(false);

  const signOut = () => { clearToken(); router.push('/login'); };

  return (
    <>
      {/* Desktop Sidebar */}
      <aside className="hidden md:flex w-64 min-h-screen p-4 flex-col gap-1 shrink-0"
        style={{ background: 'rgba(28, 28, 30, 0.6)', backdropFilter: 'saturate(180%) blur(30px)', WebkitBackdropFilter: 'saturate(180%) blur(30px)', borderRight: '1px solid rgba(255, 255, 255, 0.06)' }}>
        <NavContent path={path} onNavigate={() => {}} onSignOut={signOut} />
      </aside>

      {/* Mobile Top Bar */}
      <header className="md:hidden fixed top-0 left-0 right-0 z-40 flex items-center justify-between px-3 py-2.5"
        style={{ background: 'rgba(20, 20, 22, 0.85)', backdropFilter: 'saturate(180%) blur(20px)', WebkitBackdropFilter: 'saturate(180%) blur(20px)', borderBottom: '1px solid rgba(255, 255, 255, 0.06)' }}>
        <button onClick={() => setOpen(true)} className="w-9 h-9 rounded-lg flex items-center justify-center"
          style={{ background: 'rgba(255,255,255,0.08)' }}>
          <Menu size={18} />
        </button>
        <div className="flex items-center gap-2">
          <div className="w-6 h-6 rounded-md flex items-center justify-center"
            style={{ background: 'linear-gradient(135deg, #0A84FF, #BF5AF2)' }}>
            <Boxes size={12} color="white" strokeWidth={2.2} />
          </div>
          <span className="font-semibold text-[13px] tracking-tight">HomeWatch</span>
        </div>
        <div className="w-9" />
      </header>

      {/* Mobile Drawer */}
      {open && (
        <>
          <div onClick={() => setOpen(false)} className="md:hidden fixed inset-0 bg-black/70 z-40" />
          <aside className="md:hidden fixed top-0 left-0 bottom-0 w-64 p-4 flex flex-col gap-1 z-50 overflow-y-auto"
            style={{ background: 'rgba(18, 18, 20, 0.98)', backdropFilter: 'saturate(180%) blur(30px)', WebkitBackdropFilter: 'saturate(180%) blur(30px)', borderRight: '1px solid rgba(255, 255, 255, 0.08)' }}>
            <button onClick={() => setOpen(false)} className="absolute top-4 right-3 w-8 h-8 rounded-lg flex items-center justify-center"
              style={{ background: 'rgba(255,255,255,0.08)' }}>
              <X size={14} />
            </button>
            <NavContent path={path} onNavigate={() => setOpen(false)} onSignOut={signOut} />
          </aside>
        </>
      )}
    </>
  );
}