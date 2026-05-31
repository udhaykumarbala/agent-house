import type { NextConfig } from "next";

// Static export, embedded into the Go binary and served by handleStatic.
// trailingSlash keeps asset/route paths stable as /mission/ -> mission/index.html.
const nextConfig: NextConfig = {
  output: "export",
  trailingSlash: true,
  images: { unoptimized: true },
  // The Go server can serve the bundle under a sub-path; override at build time.
  basePath: process.env.NEXT_BASE_PATH || "",
  assetPrefix: process.env.NEXT_ASSET_PREFIX || undefined,
};

export default nextConfig;
