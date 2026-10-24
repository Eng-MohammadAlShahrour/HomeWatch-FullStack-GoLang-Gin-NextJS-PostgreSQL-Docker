'use client';
import { useState } from 'react';
import { useRouter } from 'next/navigation';
import { api, setToken } from '@/lib/api';
import { Boxes } from 'lucide-react';

export default function LoginPage() {
  const router = useRouter();
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [err, setErr] = useState('');
  const [mode, setMode] = useState<'login' | 'register'>('login');
  const [loading, setLoading] = useState(false);

  async function submit(e: React.FormEvent) {
    e.preventDefault();
    setErr('');
    setLoading(true);
    try {
      const data = await api<{ token: string }>('/' + mode, { method: 'POST', body: JSON.stringify({ email, password }) });
      setToken(data.token);
      router.push('/dashboard');
    } catch (e: any) { setErr(e.message); }
    finally { setLoading(false); }
  }

  return (
    <div className="min-h-screen flex items-center justify-center p-4 relative overflow-hidden">
      <div className="absolute inset-0 pointer-events-none">
        <div className="absolute top-1/4 left-1/4 w-96 h-96 rounded-full blur-3xl opacity-20" style={{ background: 'radial-gradient(circle, #0A84FF, transparent)' }} />
        <div className="absolute bottom-1/4 right-1/4 w-96 h-96 rounded-full blur-3xl opacity-20" style={{ background: 'radial-gradient(circle, #BF5AF2, transparent)' }} />
      </div>
      <form onSubmit={submit} className="card w-full max-w-sm p-6 md:p-8 space-y-5 animate-fade-in relative" style={{ background: 'rgba(28, 28, 30, 0.8)' }}>
        <div className="text-center space-y-1">
          <div className="w-14 h-14 mx-auto rounded-2xl flex items-center justify-center mb-3" style={{ background: 'linear-gradient(135deg, #0A84FF, #BF5AF2)', boxShadow: '0 8px 24px rgba(10, 132, 255, 0.4)' }}>
            <Boxes size={28} color="white" strokeWidth={2.2} />
          </div>
          <h1 className="text-xl md:text-2xl font-semibold tracking-tight">{mode === 'login' ? 'Welcome back' : 'Create account'}</h1>
          <p className="text-sm" style={{ color: 'var(--label-secondary)' }}>{mode === 'login' ? 'Sign in to HomeWatch' : 'Start tracking your infrastructure'}</p>
        </div>
        {err && <div className="text-sm px-3 py-2 rounded-lg" style={{ background: 'rgba(255,69,58,0.12)', color: '#FF453A' }}>{err}</div>}
        <div className="space-y-3">
          <input className="input" placeholder="Email" type="email" value={email} onChange={e => setEmail(e.target.value)} required />
          <input className="input" placeholder="Password (min 6)" type="password" value={password} onChange={e => setPassword(e.target.value)} minLength={6} required />
        </div>
        <button type="submit" disabled={loading} className="btn btn-primary w-full py-3 disabled:opacity-50">{loading ? 'Please wait...' : mode === 'login' ? 'Sign In' : 'Create Account'}</button>
        <button type="button" className="w-full text-sm" style={{ color: 'var(--label-secondary)' }} onClick={() => { setMode(mode === 'login' ? 'register' : 'login'); setErr(''); }}>
          {mode === 'login' ? "Don't have an account? Register" : 'Have an account? Sign in'}
        </button>
      </form>
    </div>
  );
}