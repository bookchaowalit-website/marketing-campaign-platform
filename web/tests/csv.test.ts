import assert from "node:assert/strict";
import { describe, it } from "node:test";
import { toCSV } from "../src/lib/csv.ts";

describe("toCSV", () => {
  it("writes a header row from the first record, or from the given columns", () => {
    assert.equal(toCSV([{ a: 1, b: 2 }]), "a,b\n1,2");
    assert.equal(toCSV([{ a: 1, b: 2 }], ["b"]), "b\n2");
  });

  it("quotes commas, quotes and line breaks (RFC 4180)", () => {
    const csv = toCSV([{ name: 'Spring, "VIP"\nlaunch' }]);
    assert.equal(csv, 'name\n"Spring, ""VIP""\nlaunch"');
  });

  it("neutralises spreadsheet formulas in text cells", () => {
    for (const value of ["=HYPERLINK(\"x\")", "+1+1", "-2+3", "@SUM(A1)", "\tcmd"]) {
      const cell = toCSV([{ v: value }]).split("\n")[1];
      assert.ok(cell.startsWith("'") || cell.startsWith("\"'"), `${value} → ${cell}`);
    }
  });

  it("keeps numbers, including negative ones, as numbers", () => {
    assert.equal(toCSV([{ n: -5, m: 0 }]), "n,m\n-5,0");
  });

  it("writes missing values as empty cells", () => {
    assert.equal(toCSV([{ a: null, b: undefined }], ["a", "b", "c"]), "a,b,c\n,,");
  });

  it("refuses an empty export", () => {
    assert.throws(() => toCSV([]), /No data to export/);
  });
});
