/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,
  output: 'standalone',
  // Reduce memory usage
  webpack: (config) => {
    config.optimization = {
      ...config.optimization,
      minimize: false,
    };
    return config;
  },
  // API rewrites to proxy backend requests (avoids CORS issues)
  async rewrites() {
    const goApiUrl = process.env.NEXT_PUBLIC_API_GO_URL || 'http://localhost:8080';
    const pythonApiUrl = process.env.NEXT_PUBLIC_API_PYTHON_URL || 'http://localhost:8081';
    
    return [
      // Proxy Go API requests
      {
        source: '/api/go/:path*',
        destination: `${goApiUrl}/api/:path*`,
      },
      // Proxy Python API requests
      {
        source: '/api/python/:path*',
        destination: `${pythonApiUrl}/api/:path*`,
      },
    ];
  },
}

module.exports = nextConfig
