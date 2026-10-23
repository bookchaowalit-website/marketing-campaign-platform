import assert from "node:assert/strict";
import { readFileSync } from "node:fs";
import { describe, it } from "node:test";

// The UI must not show (or hard-code) sign-in credentials.
const files = ["src/app/login/page.tsx", "public/index.html"];
const secrets = ["demo@example.com", "test123", "demo123"];

describe("no demo credentials in the UI", () => {
  for (const file of files) {
    it(`${file} has no demo credentials`, () => {
      const source = readFileSync(new URL(`../${file}`, import.meta.url), "utf8");
      for (const secret of secrets) {
        assert.ok(!source.includes(secret), `${file} contains "${secret}"`);
      }
    });
  }
});
