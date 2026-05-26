type TableData = { name: string; engine: string; row_count: number; columns?: { name: string; type: string }[] };
export type SchemaData = { [db: string]: { tables?: TableData[]; loading?: boolean } };

let cachedDatabases: string[] = [];
let cachedSchemaData: SchemaData = {};
let lastLoadTime = 0;

export function getCachedDatabases(): string[] {
  return cachedDatabases;
}

export function getCachedSchemaData(): SchemaData {
  return cachedSchemaData;
}

export function setCachedDatabases(dbs: string[]): void {
  cachedDatabases = dbs;
  if (dbs.length > 0) lastLoadTime = Date.now();
}

export function setCachedSchemaData(data: SchemaData): void {
  cachedSchemaData = data;
}

export function updateSchemaDb(db: string, entry: { tables?: TableData[]; loading?: boolean }): void {
  cachedSchemaData = { ...cachedSchemaData, [db]: entry };
}

export function invalidateSchema(): void {
  cachedDatabases = [];
  cachedSchemaData = {};
  lastLoadTime = 0;
}

export function isSchemaStale(): boolean {
  return lastLoadTime > 0 && Date.now() - lastLoadTime > 86400000;
}

export function extractDatabaseNames(sql: string): string[] {
  const dbs = new Set<string>();
  const re = /(?:FROM|JOIN|INTO|TABLE)\s+([a-zA-Z_][a-zA-Z0-9_]*)\./gi;
  let m: RegExpExecArray | null;
  while ((m = re.exec(sql)) !== null) {
    dbs.add(m[1]);
  }
  return Array.from(dbs);
}
