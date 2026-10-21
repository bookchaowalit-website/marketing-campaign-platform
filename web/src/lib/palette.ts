// Stat cards keep their colour as a theme path ("success.main"), which `sx.color` understands but
// colour functions such as alpha() and a hex suffix (`${color}20`) do not. Resolve it first.

const LITERAL = /^(#|rgb|hsl|color\()/;

export function paletteColor(palette: object, color: string): string {
  if (LITERAL.test(color)) return color;
  let value: unknown = palette;
  for (const key of color.split(".")) value = (value as Record<string, unknown> | undefined)?.[key];
  if (typeof value !== "string") throw new Error(`Unknown palette colour: ${color}`);
  return value;
}
