import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  reactStrictMode: true,
  output: "standalone",
  // Proxy backend requests so the browser never calls the APIs cross-origin.
  async rewrites() {
    const goApiUrl = process.env.NEXT_PUBLIC_API_GO_URL || "http://localhost:8080";
    const pythonApiUrl = process.env.NEXT_PUBLIC_API_PYTHON_URL || "http://localhost:8081";
    return [
      { source: "/api/go/:path*", destination: `${goApiUrl}/api/:path*` },
      { source: "/api/python/:path*", destination: `${pythonApiUrl}/api/:path*` },
    ];
  },
};

export default nextConfig;
