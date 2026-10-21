import assert from "node:assert/strict";
import { describe, it } from "node:test";
import nextConfig from "../next.config.ts";

describe("next.config", () => {
  it("does not suppress TypeScript build errors", () => {
    assert.notEqual(nextConfig.typescript?.ignoreBuildErrors, true);
  });

  it("has no custom webpack hook, so Next 16 builds with Turbopack", () => {
    assert.equal((nextConfig as { webpack?: unknown }).webpack, undefined);
  });

  it("proxies the Go and Python APIs", async () => {
    const rewrites = await nextConfig.rewrites?.();
    const sources = (Array.isArray(rewrites) ? rewrites : []).map((r) => r.source);
    assert.deepEqual(sources, ["/api/go/:path*", "/api/python/:path*"]);
  });
});
