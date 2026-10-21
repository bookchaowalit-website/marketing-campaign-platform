import assert from "node:assert/strict";
import { describe, it } from "node:test";
import { paletteColor } from "../src/lib/palette.ts";

const palette = { success: { main: "#2e7d32" }, grey: { 200: "#eeeeee" } };

describe("paletteColor", () => {
  it("resolves a theme path such as success.main to the colour", () => {
    assert.equal(paletteColor(palette, "success.main"), "#2e7d32");
    assert.equal(paletteColor(palette, "grey.200"), "#eeeeee");
  });

  it("passes literal colours through unchanged", () => {
    assert.equal(paletteColor(palette, "#1976d2"), "#1976d2");
    assert.equal(paletteColor(palette, "rgb(0, 0, 0)"), "rgb(0, 0, 0)");
  });

  it("throws on an unknown path instead of producing invalid CSS", () => {
    assert.throws(() => paletteColor(palette, "nope.main"), /Unknown palette colour: nope.main/);
  });
});
