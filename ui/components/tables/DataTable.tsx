import { ReactNode } from "react";

interface Column<T> {
  header: string;
  render: (row: T) => ReactNode;
}

interface DataTableProps<T> {
  title: string;
  columns: Column<T>[];
  data: T[];
}

export default function DataTable<T>({
  title,
  columns,
  data,
}: DataTableProps<T>) {
  return (
    <div className="rounded-xl border border-zinc-800 bg-zinc-950 overflow-hidden">
      <div className="border-b border-zinc-800 px-5 py-4">
        <h2 className="font-semibold text-white">{title}</h2>
      </div>

      <table className="w-full">
        <thead className="bg-zinc-900">
          <tr>
            {columns.map((c) => (
              <th
                key={c.header}
                className="px-5 py-3 text-left text-xs uppercase text-zinc-400"
              >
                {c.header}
              </th>
            ))}
          </tr>
        </thead>

        <tbody>
          {data.map((row, index) => (
            <tr
              key={index}
              className="border-t border-zinc-800 hover:bg-zinc-900"
            >
              {columns.map((c) => (
                <td key={c.header} className="px-5 py-4 text-sm text-zinc-200">
                  {c.render(row)}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
}
