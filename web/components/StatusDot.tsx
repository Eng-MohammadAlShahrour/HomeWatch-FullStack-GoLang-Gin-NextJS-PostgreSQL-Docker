export default function StatusDot({ status }: { status: string }) {
  const color = status === 'up' ? '#30D158' : status === 'down' ? '#FF453A' : '#8E8E93';
  const glow = status === 'up' ? 'rgba(48,209,88,0.5)' : status === 'down' ? 'rgba(255,69,58,0.5)' : 'transparent';
  return <span className="inline-block w-2 h-2 rounded-full" style={{ background: color, boxShadow: '0 0 8px ' + glow }} />;
}