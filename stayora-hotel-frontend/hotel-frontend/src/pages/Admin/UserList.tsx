import { useMemo, useState } from "react";
import { RefreshCw, Search, ShieldCheck, UserRound, Users } from "lucide-react";
import { useAdminData } from "@/context/AdminDataContext";
import { AdminPageHeader } from "@/components/Admin/AdminPageHeader";
import {
  AdminEmptyState,
  AdminErrorState,
  AdminLoadingRows,
} from "@/components/Admin/AdminStates";

const ROLE_STYLES: Record<string, string> = {
  admin: "bg-ink text-white",
  manager: "bg-teal/10 text-teal",
  customer: "bg-bg text-muted",
};

export function UserList() {
  const { users, isLoading, errors, refresh } = useAdminData();
  const [search, setSearch] = useState("");

  const filtered = useMemo(() => {
    const query = search.trim().toLowerCase();
    if (!query) return users;
    return users.filter((u) =>
      [u.first_name, u.last_name, u.email, u.phone].join(" ").toLowerCase().includes(query)
    );
  }, [users, search]);

  return (
    <div className="mx-auto max-w-5xl">
      <AdminPageHeader
        title="Guests"
        subtitle="Everyone registered with NileStay."
        actions={
          <button
            onClick={refresh}
            disabled={isLoading}
            title="Refresh"
            aria-label="Refresh guests"
            className="flex h-10 w-10 items-center justify-center rounded-full border border-line text-muted transition-colors hover:border-teal/40 hover:text-teal disabled:cursor-not-allowed disabled:opacity-50 cursor-pointer"
          >
            <RefreshCw size={16} className={isLoading ? "animate-spin" : ""} />
          </button>
        }
      />

      <div className="mt-6 relative max-w-xs">
        <Search
          size={15}
          className="pointer-events-none absolute left-3.5 top-1/2 -translate-y-1/2 text-muted"
        />
        <input
          value={search}
          onChange={(e) => setSearch(e.target.value)}
          placeholder="Search by name, email, or phone..."
          className="w-full rounded-pill border border-line bg-white py-2.5 pl-10 pr-4 text-sm text-ink outline-none focus:border-teal focus:ring-4 focus:ring-teal/10"
        />
      </div>

      <div className="mt-6">
        {isLoading && <AdminLoadingRows count={5} />}

        {!isLoading && errors.users && (
          <AdminErrorState message={`Failed to load guests: ${errors.users}`} onRetry={refresh} />
        )}

        {!isLoading && !errors.users && filtered.length === 0 && (
          <AdminEmptyState
            icon={Users}
            title="No guests found"
            description={
              users.length === 0 ? "No one has registered yet." : "Try a different search."
            }
          />
        )}

        {!isLoading && !errors.users && filtered.length > 0 && (
          <div className="overflow-x-auto rounded-[14px] border border-line bg-white">
            <table className="w-full min-w-[640px] text-left text-sm">
              <thead>
                <tr className="border-b border-line text-xs font-semibold uppercase tracking-wide text-muted">
                  <th className="px-4 py-3 font-semibold">Guest</th>
                  <th className="px-4 py-3 font-semibold">Contact</th>
                  <th className="px-4 py-3 font-semibold">Role</th>
                  <th className="px-4 py-3 font-semibold">Status</th>
                  <th className="px-4 py-3 font-semibold">Joined</th>
                </tr>
              </thead>
              <tbody>
                {filtered.map((user) => (
                  <tr key={user.id} className="border-b border-line last:border-0">
                    <td className="px-4 py-3">
                      <div className="flex items-center gap-2.5">
                        <span className="flex h-8 w-8 shrink-0 items-center justify-center rounded-full bg-cream text-teal">
                          <UserRound size={14} />
                        </span>
                        <span className="font-medium text-ink">
                          {user.first_name} {user.last_name}
                        </span>
                      </div>
                    </td>
                    <td className="px-4 py-3 text-muted">
                      <div>{user.email}</div>
                      {user.phone && <div className="text-xs">{user.phone}</div>}
                    </td>
                    <td className="px-4 py-3">
                      <span
                        className={`inline-flex items-center rounded-pill px-2.5 py-1 text-xs font-semibold capitalize ${
                          ROLE_STYLES[user.role] ?? "bg-bg text-muted"
                        }`}
                      >
                        {user.role}
                      </span>
                    </td>
                    <td className="px-4 py-3">
                      {user.is_active ? (
                        <span className="inline-flex items-center gap-1 text-xs font-semibold text-emerald-600">
                          <ShieldCheck size={12} />
                          Active
                        </span>
                      ) : (
                        <span className="text-xs font-semibold text-muted">Inactive</span>
                      )}
                    </td>
                    <td className="px-4 py-3 text-muted">
                      {new Date(user.created_at).toLocaleDateString()}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </div>
    </div>
  );
}
