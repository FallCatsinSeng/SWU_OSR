/** @type {import('next').NextConfig} */
const nextConfig = {
  output: 'standalone',

  images: {
    remotePatterns: [
      {
        protocol: 'https',
        hostname: 'avatars.githubusercontent.com',
      },
      {
        protocol: 'https',
        hostname: '*.githubusercontent.com',
      },
      {
        protocol: 'https',
        hostname: 'github.com',
      },
    ],
  },

  // Proxy /uploads/* requests to the backend in development.
  // In production, nginx handles this directly and these rewrites are unused.
  async rewrites() {
    const apiUrl = process.env.NEXT_PUBLIC_API_URL || '';
    // Only apply rewrites when API URL points to a different origin (development)
    if (apiUrl.startsWith('http')) {
      const origin = new URL(apiUrl).origin;
      return [
        {
          source: '/uploads/:path*',
          destination: `${origin}/uploads/:path*`,
        },
      ];
    }
    return [];
  },
};

module.exports = nextConfig;
