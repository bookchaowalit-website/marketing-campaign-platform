// CSV text for exports: RFC 4180 quoting plus a guard against spreadsheet formula
// injection (OWASP "CSV Injection"): text cells starting with = + - @ tab or CR get a leading '.

const FORMULA_START = /^[=+\-@\t\r]/;

function cell(value: unknown): string {
  if (value === null || value === undefined) return "";
  if (typeof value === "number" || typeof value === "boolean") return String(value);
  let text = String(value);
  if (FORMULA_START.test(text)) text = `'${text}`;
  return /[",\r\n]/.test(text) ? `"${text.replace(/"/g, '""')}"` : text;
}

export function toCSV(data: Record<string, unknown>[], columns?: string[]): string {
  if (!data || data.length === 0) throw new Error("No data to export");
  const headers = columns && columns.length > 0 ? columns : Object.keys(data[0]);
  const rows = [headers.map(cell).join(",")];
  for (const row of data) rows.push(headers.map((header) => cell(row[header])).join(","));
  return rows.join("\n");
}
