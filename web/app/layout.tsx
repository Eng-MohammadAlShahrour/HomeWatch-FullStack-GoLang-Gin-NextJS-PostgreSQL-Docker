import './globals.css';
import type { ReactNode } from 'react';

export const metadata = {
  title: 'HomeWatch - Monitor your home lab',
  description: 'Track your home lab devices and services',
  icons: {
    icon: '/favicon.svg',
  },
};

export default function RootLayout({ children }: { children: ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}